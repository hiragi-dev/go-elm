GOROOT := $(shell go env GOROOT)
OUT    := web

.PHONY: build serve clean

build: $(OUT)/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o $(OUT)/main.wasm .

$(OUT)/wasm_exec.js:
	cp "$(GOROOT)/lib/wasm/wasm_exec.js" $@

serve: build
	python3 -m http.server 8080 --directory $(OUT)

clean:
	rm -f $(OUT)/main.wasm $(OUT)/wasm_exec.js
