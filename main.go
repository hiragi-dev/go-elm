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

func main() {
	h := HTML[UserMsg]{}
	// ui := UI[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.Class("bg-base-200 h-full antialiased"),
			h.Div(
				h.Class("navbar bg-base-100 shadow-sm px-6"),
				h.Div(
					h.Class("flex-1"),
					h.P(h.Text("hiragi's portfolio"), h.Class("text-xl font-bold")),
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
				h.Class("mx-8 my-8 flex justify-center"),
				// h.Hn(2,
				// 	h.Text("ABOUT"),
				// 	h.Class("text-sm mb-2 font-bold tracking-tight text-violet-600"),
				// ),
				h.Div(
					h.Class("card bg-base-100 shadow-lg w-1/2"),
					// h.Div(
					// 	h.Class("card-title ml-10 mt-5"),
					// 	h.Hn(2,
					// 		h.Text("ABOUT"),
					// 		h.Class("text-sm font-bold tracking-tight text-violet-600"),
					// 	),
					// ),
					h.Div(
						h.Class("card-body"),
						h.Div(
							h.Class("flex"),
							h.Div(
								h.Class("flex items-center"),
								h.Div(
									h.Class("avatar"),
									h.Div(
										h.Class("w-40 rounded-full"),
										h.Img("https://avatars.githubusercontent.com/u/298766645?v=4"),
									),
								),
							),
							h.Div(
								h.Class("ml-8 flex flex-1 flex-col"),
								h.Hn(2,
									h.Text("ABOUT"),
									h.Class("text-xs font-bold tracking-tight text-violet-600"),
								),
								h.Hn(
									2,
									h.Text("Hi, I'm hiragi 👋 — a student and hobbyist engineer in Japan."),
									h.Class("text-2xl font-semibold text-gray-900"),
								),
								h.Div(
									// h.Class("flex flex-1 items-center mb-2"),
									h.Class("mt-2"),
									h.P(
										h.Text("I'm interested in low-level programming, and have built things like an SSH server in Rust to learn how they work. Lately I've been building an alarm clock that absolutely cannot be turned off until you walk to a location you registered in advance, running on a small M5Stack device. I'm also working as a long-term software engineering intern."),
										h.Class("text-base-content/70"),
									),
								),
							),
						),
					),
				),
			),
			h.Div(
				h.Class("mx-8 my-8 flex justify-center"),
				h.Div(
					h.Class("card bg-base-100 shadow-lg w-1/2"),
					h.Div(
						h.Class("card-body"),
						h.Div(
							h.Class("flex"),
							h.Div(
								h.Class("flex items-center"),
								h.Div(
									h.Class("avatar"),
									h.Div(
										h.Class("w-40 rounded-full"),
										h.Img("https://avatars.githubusercontent.com/u/298766645?v=4"),
									),
								),
							),
							h.Div(
								h.Class("ml-8 flex flex-1 flex-col"),
								h.Hn(2,
									h.Text("WORKS"),
									h.Class("text-xs font-bold tracking-tight text-violet-600"),
								),
								h.Hn(
									2,
									h.Text("Hi, I'm hiragi 👋 — a student and hobbyist engineer in Japan."),
									h.Class("text-2xl font-semibold text-gray-900"),
								),
								h.Div(
									// h.Class("flex flex-1 items-center mb-2"),
									h.Class("mt-2"),
									h.P(
										h.Text("I'm interested in low-level programming, and have built things like an SSH server in Rust to learn how they work. Lately I've been building an alarm clock that absolutely cannot be turned off until you walk to a location you registered in advance, running on a small M5Stack device. I'm also working as a long-term software engineering intern."),
										h.Class("text-base-content/70"),
									),
								),
							),
						),
					),
				),
			),
		)
		// h.Button(h.Attr("class", "btn btn-primary")),

		// h.Div(
		// 	h.P(h.Text("Counter: ")),
		// 	h.Text(strconv.Itoa(state.counterValue)),
		// ),

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
	root.SetAttribute("class", "h-screen")

	Launch(root, view, update, UserState{counterValue: 0})
}
