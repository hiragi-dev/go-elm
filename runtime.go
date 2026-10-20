package main

import (
	"io"
	"net/http"

	"honnef.co/go/js/dom/v2"
)

type Cmd[Env any, M any] func(r Env, send func(M))

func TaggedAsyncDo[E, M any](f func(e E) M, m M) Cmd[E, M] {
	c := AsyncDo(f)

	return func(r E, send func(M)) {
		// wrap(send)
		send(m)
		c(r, send)
	}
}

func AsyncDo[Env, M any](f func(e Env) M) Cmd[Env, M] {
	return func(r Env, send func(M)) {
		send(f(r))
	}
}

func None[Env, M any]() Cmd[Env, M] {
	return func(_ Env, _ func(M)) {

	}
}

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

func createElement[Msg any](v VNode[Msg], msgChan chan Msg) dom.Node {
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

func Launch[State any, Msg any, E any](
	root dom.Element,
	view func(State) VNode[Msg],
	update func(Msg, State) (State, Cmd[E, Msg]),
	makeEnv func(Std) E,
	init func() (State, Cmd[E, Msg]),
) {
	msgChan := make(chan Msg)
	currentState, initCmd := init()

	// メッセージが届く前の初回描画
	currentView := view(currentState)
	dom := createElement(currentView, msgChan)
	render(root, dom)

	env := makeEnv(DefaultStd())

	go initCmd(env, func(m Msg) { msgChan <- m })

	for msg := range msgChan {
		nextState, cmd := update(msg, currentState)
		currentState = nextState

		// pretty.Println(reflect.TypeOf(cmd))

		go cmd(env, func(m Msg) { msgChan <- m })

		nextView := view(nextState)
		render(root, createElement(nextView, msgChan))
		// patch := diff(currentView, nextView)

		// // pretty.Println(patch)

		// currentView = nextView

		// applyPatch(root, patch, msgChan)
	}
}
