package main

import (
	"io"
	"net/http"
	"net/url"

	"honnef.co/go/js/dom/v2"
)

type Http interface {
	Get(url string, queries url.Values) string
	Put(url string, value string, headers ...string)
}

type Console interface {
	Log(content string)
}

type Dom interface {
	GetPosition(id string) Pos
}

type Std interface {
	Http
	Console
	Dom
}

type browserStd struct {
	root dom.Element
}

func DefaultStd(root dom.Element) Std { return browserStd{root} }

func (r browserStd) GetPosition(id string) Pos {
	// pretty.Println("======", r.root)

	el := r.root.QuerySelector("#" + id)
	rect := el.GetBoundingClientRect()

	// pretty.Println(el, rect)

	return Pos{rect.Left(), rect.Top(), rect.Width(), rect.Height()}
}

func (r browserStd) Get(url string, queries url.Values) string {
	client := &http.Client{}

	req, _ := http.NewRequest(
		"GET",
		url,
		nil,
	)

	req.Header = http.Header{}

	// q := req.URL.Query()
	req.URL.RawQuery = queries.Encode()

	res, _ := client.Do(req)
	body, _ := io.ReadAll(res.Body)

	return string(body)
}

func (r browserStd) Put(url string, value string, headers ...string) {

}

func (r browserStd) Log(content string) {

}
