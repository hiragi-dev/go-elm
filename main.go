package main

import (
	"io"
	"net/http"
	"strconv"

	"honnef.co/go/js/dom/v2"
)

type UserMsg interface{ isUserMsg() }

type AddCounter struct{}

func (AddCounter) isUserMsg() {}

type GetTimeNow struct{}

func (GetTimeNow) isUserMsg() {}

type SetTimeNow struct{ time string }

func (SetTimeNow) isUserMsg() {}

type UserState struct {
	counterValue int
	time         string
}

func clear(root dom.Element) {
	for _, c := range root.ChildNodes() {
		root.RemoveChild(c)
	}
}

func render(root dom.Element, dom dom.Node) {
	clear(root)
	root.AppendChild(dom)
}

func createElement[Msg any](v VDom[Msg], msgChan chan Msg) dom.Node {
	doc := dom.GetWindow().Document()

	if v.text == "" {
		el := doc.CreateElement(v.tagName)

		for _, c := range v.children {
			e := createElement(c, msgChan)
			el.AppendChild(e)
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
	// switch n := v.(type) {
	// case VDiv[Msg]:
	// 	d := doc.CreateElement("div")
	// 	for _, c := range n.children {
	// 		e := createElement(c, msgChan)
	// 		d.AppendChild(e)
	// 	}
	// 	return d
	// case VP[Msg]:
	// 	return doc.CreateTextNode(n.content)
	// case VButton[Msg]:
	// 	b := doc.CreateElement("button")
	// 	b.AddEventListener("click", false, func(e dom.Event) {
	// 		if n.onClick != nil {
	// 			go func() { msgChan <- n.onClick() }()
	// 		}
	// 	})
	// 	return b
	// default:
	// 	return doc.CreateTextNode("unimplemented!")
	// }

	return doc.CreateTextNode("unreachable!")
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
		diffs := diff(currentView, nextView)

		currentView = nextView

		// pretty.Println(diffs)

		patch(root, diffs, msgChan)

		// dom := createElement(currentView, msgChan)
		// render(root, dom)
	}
}

type HttpGet[Msg any] struct {
	onReceive func(string) Msg
	url       string
}

func (HttpGet[Msg]) isCmd() {}

type Cmd[Msg any] interface{ isCmd() }

type EmptyCmd struct{}

func (EmptyCmd) isCmd() {}

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.P(h.Text("Counter: "+strconv.Itoa(state.counterValue))),
			h.Button(func() UserMsg { return AddCounter{} }, h.P(h.Text("click me"))),
			h.Button(func() UserMsg { return GetTimeNow{} }, h.Hn(3, strconv.Itoa(state.counterValue))),
			h.Hn(4, If(state.time != "", state.time, "none")),
		)
		// if state.counterValue == 0 {
		// 	return h.Button(func() UserMsg { return AddCounter{} }, h.Hn(3, "A"))
		// } else {
		// 	return h.Hn(3, "B")
		// }
	}

	// a := func(state UserState) VDom[UserMsg] {
	// 	return h.Div(h.P(h.Hn(1, "hp")))
	// }

	// b := func(state UserState) VDom[UserMsg] {
	// 	return h.Div(h.P(h.Hn(2, "henka")))
	// }

	// initState := UserState{}

	// diffU := diff(a(initState), b(initState))

	// pretty.Println(a(initState), b(initState), diffU)

	// c := make(chan UserMsg)

	// _ = createElement(a(initState), c)
	// bn := createElement(b(initState), c)

	// patch(bn, diffU, c)

	update := func(msg UserMsg, state UserState) (UserState, Cmd[UserMsg]) {
		// pretty.Println(msg)

		switch m := msg.(type) {

		case AddCounter:
			// pretty.Println(m)
			return UserState{counterValue: state.counterValue + 1}, EmptyCmd{}
		case GetTimeNow:
			return state, HttpGet[UserMsg]{url: "https://jsonplaceholder.typicode.com/todos/1", onReceive: func(r string) UserMsg { return SetTimeNow{time: r} }}
		case SetTimeNow:
			time := m.time
			return UserState{counterValue: state.counterValue, time: time}, EmptyCmd{}
		default:
			return UserState{counterValue: state.counterValue + 1}, EmptyCmd{}
		}
	}

	root := dom.GetWindow().Document().QuerySelector("#app")

	// root.AppendChild(bn)
	Launch(root, view, update, UserState{counterValue: 0})
}
