package main

import (
	"io"
	"net/http"
	"reflect"

	"github.com/kr/pretty"
	"honnef.co/go/js/dom/v2"
)

// type HttpGet[Msg any] struct {
// 	onReceive func(string) Msg
// 	url       string
// 	headers   []string
// }

// func (HttpGet[Msg]) isCmd() {}

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

// func handleCmd[Msg any](cmd Cmd[Msg], msgChan chan Msg) {
// 	switch c := cmd.(type) {
// 	case HttpGet[Msg]:
// 		go func() {
// 			r := fetchUrl(c.url)
// 			m := c.onReceive(r)
// 			msgChan <- m
// 		}()
// 	}
// }

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

func Launch[State any, Msg any](
	root dom.Element,
	view func(State) VDom[Msg],
	update func(*Eff, Msg, State) State,
	init func(*Eff) State,
) {

	eff := &Eff{req: make(chan Op[any]), res: make(chan any)}

	msgChan := make(chan Msg)
	currentState := init(eff)

	// メッセージが届く前の初回描画
	currentView := view(currentState)
	dom := createElement(currentView, msgChan)
	render(root, dom)

	go func() {
		for msg := range msgChan {
			nextState := update(eff, msg, currentState)

			nextView := view(nextState)
			patch := diff(currentView, nextView)

			currentView = nextView

			applyPatch(root, patch, msgChan)
		}
	}()

	handlers := SystemHandler()

	for op := range eff.req {
		switch r := op.(type) {
		case HttpGet:
			handler := handlers[reflect.TypeOf(r)]
			pretty.Println("Runtime: ", handler)
			eff.res <- handler(r)
			pretty.Println("Runtime: end handle")
		}

	}
}

// for msg := range msgChan {
// 	nextState, cmd := update(msg, currentState)
// 	currentState = nextState

// 	go handleCmd(cmd, msgChan)

// 	nextView := view(nextState)
// 	patch := diff(currentView, nextView)

// 	currentView = nextView

// 	applyPatch(root, patch, msgChan)
// }
