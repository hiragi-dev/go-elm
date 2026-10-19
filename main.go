package main

import (
	"github.com/kr/pretty"
)

type UserMsg interface{ isUserMsg() }

type AddCounter struct{}

func (AddCounter) isUserMsg() {}

type GetTimeNow struct{}

func (GetTimeNow) isUserMsg() {}

type SetTimeNow struct{ time string }

func (SetTimeNow) isUserMsg() {}

type GetGitHubActivity struct{ response string }

func (GetGitHubActivity) isUserMsg() {}

type GotGitHubActivity struct{ response string }

func (GotGitHubActivity) isUserMsg() {}

type UserState struct {
	counterValue int
	time         string
	commits      []string
}

func (h HTML[Msg]) Typ(content string) VDom[Msg] {
	return h.P(
		h.Text(content),
		h.Class("opacity-60 text-xs"),
	)
}

func (h HTML[Msg]) Link(content string, url string) VDom[Msg] {
	return h.A(
		h.Class("link link-hover text-md font-semibold tracking-tight text-violet-600"),
		h.Text(
			content,
		),
	)
}

func (h HTML[M]) Profile() VDom[M] {
	tableEntry := func(key string, value string) VDom[M] {
		return h.Tr(
			h.Th(
				h.Class("px-0"),
				h.Text(key),
			),
			h.Td(
				h.Text(value),
			),
		)
	}

	return h.Content(
		"PROFILE",
		h.Div(
			h.Class("card bg-base-100 shadow-lg"),
			h.Div(
				h.Class("card-body"),
				h.Div(
					h.Class("flex"),
					h.Div(
						h.Class("flex items-center w-1/4"),
						h.Div(
							h.Class("avatar"),
							h.Div(
								h.Class("rounded-full"),
								h.Img("https://avatars.githubusercontent.com/u/298766645?v=4"),
							),
						),
					),
					h.Div(
						h.Class("ml-8 flex flex-col w-full"),
						h.Hn(
							2,
							h.Text("hiragi-dev"),
							h.Class("text-3xl font-semibold text-gray-900"),
						),
						h.Div(h.Class("my-1")),
						h.Table(
							h.Class("table"),
							h.Tbody(
								tableEntry("Currently", "Student / Long-term SWE intern"),
								tableEntry("Languages", "Rust / C / Python / TypeScript"),
								tableEntry("Interests", " Low-level programming / Algebraic effects / DIY keyboards"),
							),
						),
					),
				),
			),
		),
	)
}

func (h HTML[M]) Content(header string, inner VDom[M]) VDom[M] {
	return h.Div(
		h.Hn(2,
			h.Text(header),
			h.Class("mb-2 text-md font-semibold tracking-tight text-violet-600"),
		),
		inner,
	)
}

func (h HTML[M]) Works() VDom[M] {
	listRow := func(header VDom[M], title string, description string, link string) VDom[M] {
		return h.Li(
			h.Class("list-row transition-colors duration-200 hover:bg-base-300 "),
			header,
			h.A(
				h.Class("list-col-grow ml-2"),
				h.Attr("href", link),
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

	return h.Content(
		"WORKS",
		h.Ul(
			h.Class("list bg-base-100 rounded-box shadow-md overflow-hidden"),
			listRow(
				h.Div(
					h.Class("size-10 rounded-box"),
					h.Img("https://cdn.simpleicons.org/github"),
				),
				"OneAlarm",
				"絶対に電源が切れないアラーム",
				"https://github.com/hiragi-dev/onealarm-doc",
			),
			listRow(
				h.Div(
					h.Class("size-10 rounded-box"),
					h.Img("https://cdn.simpleicons.org/github"),
				),
				"ssh-server",
				"RFCを参照して実装したRust製SSHサーバ",
				"https://github.com/hiragi-dev/ssh-server",
			),
		),
	)
}

func (h HTML[Msg]) NavBar() VDom[Msg] {
	return h.Div(
		h.Class("navbar bg-base-100 shadow-sm px-40"),
		h.Div(
			h.Class("flex-1"),
			h.P(h.Text("hiragi's portfolio"), h.Class("text-xl font-semibold")),
		),
		h.Div(
			h.Class("ml-auto flex gap-x-4"),
			h.Link("PROFILE", "ok"),
			h.Link("ACTIVITY", "ok"),
			h.Link("WORKS", "ok"),
		),
	)
}

func (h HTML[Msg]) Activity(commits []string) VDom[Msg] {
	commitsDisp := h.Div()
	for _, commit := range commits {
		commitsDisp.children = append(commitsDisp.children,
			h.Typ(commit),
		)
	}

	pretty.Println(commits)

	return h.Content(
		"ACTIVITY",
		h.Div(
			h.Class("card bg-base-100 shadow-lg"),
			h.Div(
				h.Class("card-body"),
				h.Div(
					h.Class("flex flex-col"),
					commitsDisp,
				),
			),
		),
	)
}

type Continue[M any] interface{ AndThen(M) Cmd[M] }

func GET[M any](url string, headers ...string) HttpGet[M] {
	_headers := make([]string, 0)
	_headers = append(_headers, headers...)

	return HttpGet[M]{
		url: url,
		// onReceive: onReceive,
		headers: _headers,
	}
}

func (h HttpGet[M]) AndThen(c func(string) M) Cmd[M] {
	h.onReceive = c
	return h
}

type Eff[T, Op any] struct {
	op    chan Op
	opRes chan any
}

func Handle[T, Op any](program func(*Eff[T, Op])) any {
	eff := &Eff[T, Op]{op: make(chan Op), opRes: make(chan any)}

	go func() {
		pretty.Println("run program")
		program(eff)
	}()

	println("wait")
	for op := range eff.op {
		pretty.Println("handle op: ", op)
		eff.opRes <- "res"
	}

	return nil
}

func (eff *Eff[T, Op]) Perform(op any) any {
	eff.op <- op.(Op)
	res := <-eff.opRes

	return res
}

// func Do(func (eff *Eff[T, Op]))

func effectfull[T, Op any](eff *Eff[T, Op]) {
	pretty.Println("effectfull: start")
	res := eff.Perform("Operation")

	pretty.Println("effectfull: ", res)
}

func main() {
	// h := HTML[UserMsg]{}
	// // ui := UI[UserMsg]{}

	// effectfull()
	res := Handle[string, string](effectfull)
	pretty.Println("handle: ", res)

	// view := func(state UserState) VDom[UserMsg] {
	// 	return h.Div(
	// 		h.Class("min-h-screen bg-base-200 antialiased"),
	// 		h.NavBar(),
	// 		h.Div(
	// 			h.Class("flex flex-col px-40 gap-y-10 mt-10"),
	// 			h.Profile(),
	// 			h.Activity(state.commits),
	// 			h.Works(),
	// 		),
	// 	)
	// }

	// var update func(UserMsg, UserState) (UserState, Cmd[UserMsg])

	// update = func(msg UserMsg, state UserState) (UserState, Cmd[UserMsg]) {
	// 	switch m := msg.(type) {

	// 	case AddCounter:
	// 		return UserState{counterValue: state.counterValue + 1}, EmptyCmd{}
	// 	case GetTimeNow:
	// 		return state, HttpGet[UserMsg]{url: "https://jsonplaceholder.typicode.com/todos/1", onReceive: func(r string) UserMsg { return SetTimeNow{time: r} }}
	// 	case SetTimeNow:
	// 		time := m.time
	// 		return UserState{counterValue: state.counterValue, time: time}, EmptyCmd{}
	// 	case GotGitHubActivity:
	// 		var jsonMaps []map[string]any
	// 		_ = json.Unmarshal([]byte(m.response), &jsonMaps)

	// 		var commits []string

	// 		for _, jsonMap := range jsonMaps {
	// 			commitEntry, ok := jsonMap["commit"].(map[string]any)
	// 			if !ok {
	// 				println("ERROR:")
	// 			}
	// 			messageEntry, ok := commitEntry["message"].(string)
	// 			if !ok {
	// 				println("ERROR:")
	// 			}

	// 			pretty.Println(commitEntry)

	// 			commits = append(commits, messageEntry)
	// 		}

	// 		state.commits = commits

	// 		return state, EmptyCmd{}
	// 	default:
	// 		return UserState{counterValue: state.counterValue + 1}, EmptyCmd{}
	// 	}
	// }

	// root := dom.GetWindow().Document().QuerySelector("#app")
	// // root.SetAttribute("class", "h-screen")

	// init := func() (UserState, Cmd[UserMsg]) {

	// 	// HttpGet{
	// 	// 	url: "https://api.github.com/repos/hiragi-dev/go-elm/commits"
	// 	// }
	// 	initCmd :=
	// 		GET[UserMsg](
	// 			"https://api.github.com/repos/hiragi-dev/go-elm/commits",
	// 			"Accept: application/vnd.github+json",
	// 			"X-GitHub-Api-Version: 2026-03-10",
	// 		).
	// 			AndThen(func(r string) UserMsg { return GotGitHubActivity{r} })

	// 	return UserState{counterValue: 0, commits: make([]string, 0)}, initCmd
	// }

	// Launch(root, view, update, init)
}
