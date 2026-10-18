package main

import (
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

func (h HTML[Msg]) Typ(content string) VDom[Msg] {
	return h.P(
		h.Text(content),
		h.Class("text-base-content/70 text-xs"),
	)
}

func (h HTML[M]) Profile() VDom[M] {
	return h.Div(
		h.Hn(2,
			h.Text("ABOUT"),
			h.Class("mb-2 text-md font-semibold tracking-tight text-violet-600"),
		),
		h.Div(
			h.Class("card bg-base-100 shadow-lg"),
			h.Div(
				h.Class("card-body"),
				h.Div(
					h.Class("flex"),
					h.Div(
						h.Class("flex items-center"),
						h.Div(
							h.Class("avatar"),
							h.Div(
								h.Class("rounded-full"),
								h.Img("https://avatars.githubusercontent.com/u/298766645?v=4"),
							),
						),
					),
					h.Div(
						h.Class("ml-8 flex flex-col justify-center"),
						h.Hn(
							2,
							h.Text("Hi, I'm hiragi 👋 — a student and hobbyist engineer in Japan."),
							h.Class("text-2xl font-semibold text-gray-900"),
						),
						h.Div(
							h.Class("mt-2"),
							h.P(
								h.Text("I'm interested in low-level programming, and have built things like an SSH server in Rust to learn how they work. Lately I've been building an alarm clock that absolutely cannot be turned off until you walk to a location you registered in advance, running on a small M5Stack device. I'm also working as a long-term software engineering intern."),
								h.Class("text-base opacity-60"),
							),
						),
					),
				),
			),
		),
	)
}

func (h HTML[M]) Works() VDom[M] {
	listRow := func(header VDom[M], title string, description string) VDom[M] {
		return h.Li(
			h.Class("list-row ml-2"),
			header,
			h.Div(
				h.Class("list-col-grow ml-2"),
				h.Div(
					h.Class("text-md font-semibold"),
					h.Text(title),
				),
				h.P(
					h.Text(description),
					h.Class("opacity-60"),
				),
			),
		)
	}

	return h.Div(
		h.Hn(2,
			h.Text("WORKS"),
			h.Class("mb-2 text-md font-semibold tracking-tight text-violet-600"),
		),
		h.Ul(
			h.Class("list bg-base-100 rounded-box shadow-md"),
			listRow(
				h.Div(
					h.Class("size-10 rounded-box"),
					h.Img("https://cdn.simpleicons.org/github"),
				),
				"OneAlarm",
				"絶対に電源が切れないアラーム",
			),
			listRow(
				h.Div(
					h.Class("size-10 rounded-box"),
					h.Img("https://cdn.simpleicons.org/github"),
				),
				"ssh-server",
				"RFCを参照して実装したRust製SSHサーバ",
			),
		),
	)
}

func main() {
	h := HTML[UserMsg]{}
	// ui := UI[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.Class("min-h-screen bg-base-200 antialiased"),
			h.Div(
				h.Class("navbar bg-base-100 shadow-sm px-40"),
				h.Div(
					h.Class("flex-1"),
					h.P(h.Text("hiragi's portfolio"), h.Class("text-xl font-semibold")),
				),
				h.Div(
					h.Class("ml-auto"),
					h.Div(
						h.Class("btn btn-ghost btn-circle avatar "),
						h.Div(
							h.Class("w-10 rounded-full"),
							h.Img("https://avatars.githubusercontent.com/u/298766645?v=4"),
						),
					),
				),
			),
			h.Div(
				h.Class("flex flex-col px-40 gap-y-10 mt-10"),
				h.Profile(),
				h.Works(),
			),
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
	// root.SetAttribute("class", "h-screen")

	Launch(root, view, update, UserState{counterValue: 0})
}
