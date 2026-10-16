package main

import (
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

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.P(h.Text("Counter: "+strconv.Itoa(state.counterValue))),
			h.Button(func() UserMsg { return AddCounter{} }, h.P(h.Text("click me")), h.Attr("class", "bg-blue-500 text-white font-bold py-2 px-4 rounded")),
			h.Button(func() UserMsg { return GetTimeNow{} }, h.Hn(3, strconv.Itoa(state.counterValue))),
			h.Hn(4, If(state.time != "", state.time, "none")),
		)
	}

	update := func(msg UserMsg, state UserState) (UserState, Cmd[UserMsg]) {
		switch m := msg.(type) {

		case AddCounter:
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

	Launch(root, view, update, UserState{counterValue: 0})
}
