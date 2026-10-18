// devserver は開発用のホットリロードサーバー。
//
// ルートの .go ファイルが変わったら wasm を再ビルドし、web/ 以下のファイルが変わったら
// そのまま、接続中のブラウザにリロードを指示する。index.html にはリロード用のスクリプトを
// 配信時に差し込むので、web/index.html 自体には何も書かなくてよい。
//
//	go run ./cmd/devserver            # http://localhost:8080
//	go run ./cmd/devserver -addr :3000
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

const reloadScript = `<script>
(() => {
  const es = new EventSource("/__devserver/events");
  es.addEventListener("reload", () => location.reload());
  es.addEventListener("build-error", (e) => console.error("[devserver] ビルドに失敗しました\n" + JSON.parse(e.data)));
})();
</script>
`

func main() {
	addr := flag.String("addr", ":8080", "待ち受けるアドレス")
	webDir := flag.String("web", "web", "配信するディレクトリ")
	flag.Parse()

	hub := newHub()

	if err := build(); err != nil {
		log.Printf("ビルドに失敗しました:\n%s", err)
	}

	go watch(".", *webDir, hub)

	mux := http.NewServeMux()
	mux.HandleFunc("/__devserver/events", hub.serveEvents)
	mux.Handle("/", noCache(serveWeb(*webDir)))

	log.Printf("http://localhost%s で配信中（Ctrl+C で終了）", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

// build は web/main.wasm をビルドする。失敗したらコンパイラの出力をエラーとして返す。
func build() error {
	start := time.Now()
	cmd := exec.Command("go", "build", "-o", "web/main.wasm", ".")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}
	log.Printf("ビルド完了（%s）", time.Since(start).Round(time.Millisecond))
	return nil
}

// watch は 500ms ごとに更新時刻を調べる。依存を増やさないよう fsnotify は使わない。
func watch(srcDir, webDir string, hub *hub) {
	goFiles := func() map[string]time.Time {
		m := map[string]time.Time{}
		for _, pattern := range []string{"*.go", "go.mod", "go.sum"} {
			matches, _ := filepath.Glob(filepath.Join(srcDir, pattern))
			for _, p := range matches {
				if info, err := os.Stat(p); err == nil {
					m[p] = info.ModTime()
				}
			}
		}
		return m
	}
	webFiles := func() map[string]time.Time {
		m := map[string]time.Time{}
		filepath.WalkDir(webDir, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || filepath.Base(p) == "main.wasm" {
				return nil
			}
			if info, err := d.Info(); err == nil {
				m[p] = info.ModTime()
			}
			return nil
		})
		return m
	}

	prevGo, prevWeb := goFiles(), webFiles()
	for range time.Tick(500 * time.Millisecond) {
		curGo, curWeb := goFiles(), webFiles()

		if changed(prevGo, curGo) {
			log.Println("Go のソースが変わったので再ビルドします")
			if err := build(); err != nil {
				log.Printf("ビルドに失敗しました:\n%s", err)
				hub.broadcast("build-error", err.Error())
			} else {
				hub.broadcast("reload", "")
			}
		} else if changed(prevWeb, curWeb) {
			log.Println("web/ のファイルが変わったのでリロードします")
			hub.broadcast("reload", "")
		}

		prevGo, prevWeb = curGo, curWeb
	}
}

func changed(prev, cur map[string]time.Time) bool {
	if len(prev) != len(cur) {
		return true
	}
	for p, t := range cur {
		if !prev[p].Equal(t) {
			return true
		}
	}
	return false
}

// serveWeb は web/ を配信する。index.html だけはリロード用スクリプトを差し込んで返す。
func serveWeb(dir string) http.Handler {
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			files.ServeHTTP(w, r)
			return
		}
		html, err := os.ReadFile(filepath.Join(dir, "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if i := bytes.LastIndex(html, []byte("</body>")); i >= 0 {
			html = append(html[:i:i], append([]byte(reloadScript), html[i:]...)...)
		} else {
			html = append(html, reloadScript...)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})
}

// noCache はブラウザに古い main.wasm を使わせないためのもの。
func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		h.ServeHTTP(w, r)
	})
}

// hub は Server-Sent Events で接続中のブラウザにイベントを配る。
type hub struct {
	mu      sync.Mutex
	clients map[chan string]struct{}
}

func newHub() *hub {
	return &hub{clients: map[chan string]struct{}{}}
}

func (h *hub) broadcast(event, data string) {
	// SSE の data は改行を含められないので、JSON 文字列にして送る
	encoded, _ := json.Marshal(data)
	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", event, encoded)

	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c <- msg:
		default:
		}
	}
}

func (h *hub) serveEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	flusher.Flush()

	c := make(chan string, 4)
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
	}()

	for {
		select {
		case msg := <-c:
			fmt.Fprint(w, msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
