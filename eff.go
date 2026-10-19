package main

import (
	"io"
	"net/http"
	"net/url"
)

type Http interface {
	Get(url string, queries url.Values) string
	Put(url string, value string, headers ...string)
}

type Console interface {
	Log(content string)
}

type Std interface {
	Http
	Console
}

type browserStd struct{}

func DefaultStd() Std { return browserStd{} }

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
