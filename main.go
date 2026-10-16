package main

import (
	"strconv"

	"honnef.co/go/js/dom/v2"
)

type UserMsg struct {
}

type UserState struct {
	counterValue int
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

func createElement[Msg any](v VNode[Msg], msg_chan chan Msg) dom.Node {
	doc := dom.GetWindow().Document()

	switch n := v.(type) {
	case VDiv[Msg]:
		d := doc.CreateElement("div")
		for _, c := range n.children {
			e := createElement(c, msg_chan)
			d.AppendChild(e)
		}
		return d
	case VP[Msg]:
		return doc.CreateTextNode(n.content)
	case VButton[Msg]:
		b := doc.CreateElement("button")
		b.AddEventListener("click", false, func(e dom.Event) {
			if n.onClick != nil {
				go func() { msg_chan <- n.onClick() }()
			}
		})
		return b
	default:

	}

	return doc.CreateTextNode("unreachable!")
}

func Launch[State any, Msg any](root dom.Element, view func(State) VNode[Msg], handler func(Msg, State) State, initialState State) {
	msg_chan := make(chan Msg)

	// メッセージが届く前の初回描画
	currentState := initialState
	currentView := view(initialState)
	dom := createElement(currentView, msg_chan)
	render(root, dom)

	for msg := range msg_chan {
		currentState = handler(msg, currentState)
		currentView := view(currentState)
		dom := createElement(currentView, msg_chan)
		render(root, dom)
	}
}

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VNode[UserMsg] {
		return h.Div(
			h.P("Counter: "+strconv.Itoa(state.counterValue)),
			h.Button(func() UserMsg { return UserMsg{} }, h.P("click me")),
		)
	}

	handler := func(msg UserMsg, state UserState) UserState {
		println(state.counterValue)

		return UserState{counterValue: state.counterValue + 1}
	}

	root := dom.GetWindow().Document().QuerySelector("#app")
	Launch(root, view, handler, UserState{counterValue: 0})
}
