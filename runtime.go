package main

import (
	"io"
	"net/http"

	"honnef.co/go/js/dom/v2"
)

type HttpGet[Msg any] struct {
	onReceive func(string) Msg
	url       string
}

func (HttpGet[Msg]) isCmd() {}

type Cmd[Msg any] interface{ isCmd() }

type EmptyCmd struct{}

func (EmptyCmd) isCmd() {}

func render(root dom.Element, dom dom.Node) {
	clear(root)
	root.AppendChild(dom)
}

func clear(root dom.Element) {
	for _, c := range root.ChildNodes() {
		root.RemoveChild(c)
	}
}

func fetchUrl(url string) string {
	res, _ := http.Get(url)
	body, _ := io.ReadAll(res.Body)

	return string(body)
}

func handleCmd[Msg any](cmd Cmd[Msg], msgChan chan Msg) {
	switch c := cmd.(type) {
	case HttpGet[Msg]:
		go func() {
			r := fetchUrl(c.url)
			m := c.onReceive(r)
			msgChan <- m
		}()
	}
}

func createElement[Msg any](v VDom[Msg], msgChan chan Msg) dom.Node {
	doc := dom.GetWindow().Document()

	if v.text == "" {
		el := doc.CreateElement(v.tagName)

		for _, c := range v.children {
			e := createElement(c, msgChan)
			el.AppendChild(e)
		}

		// for key, value := range v.attributes {
		// println("set attribute", key, value)
		// el.SetAttribute(key, value)
		// }

		for _, attr := range v.attributes {
			el.SetAttribute(attr.key, attr.value)
		}

		for eventName, update := range v.event {
			el.AddEventListener(eventName, false, func(e dom.Event) {
				go func() { msgChan <- update() }()
			})
		}

		return el
	} else {
		t := doc.CreateTextNode(v.text)
		return t
	}
}

func Launch[State any, Msg any](root dom.Element, view func(State) VDom[Msg], update func(Msg, State) (State, Cmd[Msg]), initialState State) {
	msgChan := make(chan Msg)
	currentState := initialState

	// メッセージが届く前の初回描画
	currentView := view(initialState)
	dom := createElement(currentView, msgChan)
	render(root, dom)

	for msg := range msgChan {
		nextState, cmd := update(msg, currentState)
		currentState = nextState

		go handleCmd(cmd, msgChan)

		nextView := view(nextState)
		patch := diff(currentView, nextView)

		currentView = nextView

		applyPatch(root, patch, msgChan)
	}
}
