package main

import (
	"time"

	"github.com/mono83/maybe"
	"honnef.co/go/js/dom/v2"
)

type UserMsg interface{ isUserMsg() }

type SetGitHubActivity struct {
	commits []GithubCommit
}

func (SetGitHubActivity) isUserMsg() {}

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

func (h HTML[Msg]) Activity(commits maybe.Maybe[[]GithubCommit]) VDom[Msg] {
	// commitsDisp := h.Div()
	// for _, commit := range commits {
	// 	commitsDisp.children = append(commitsDisp.children,
	// 		h.Typ(commit),
	// 	)
	// }

	tableEntry := func(key string, repoName string, value string) VDom[Msg] {
		return h.Tr(
			h.Th(
				h.Class("px-0"),
				h.Text(key),
			),
			h.Td(
				h.Div(
					h.Class("badge badge-soft badge-primary"),
					h.Text(repoName),
				),
			),
			h.Td(
				h.Text(value),
			),
		)
	}

	var activities VDom[Msg]
	if commits, ok := commits.Get(); ok {
		activities = h.Tbody()
		for i, commit := range commits {
			if i < 7 {
				activities.children = append(activities.children,
					tableEntry(
						commit.date.Format("2006/1/2 15:04"),
						commit.repoName,
						commit.message,
					),
				)
			} else {
				break
			}
		}
	} else {
		activities = h.Span(
			h.Class("loading loading-spinner loading-md text-violet-600"),
		)
	}

	return h.Content(
		"ACTIVITY",
		h.Div(
			h.Class("card bg-base-100 shadow-lg"),
			h.Div(
				h.Class("card-body"),
				h.Div(
					h.Class("flex"),
					h.Div(
						h.Class("ml-8 flex flex-col w-full"),
						// h.Hn(
						// 	2,
						// 	h.P(h.Text("Commits History")),
						// 	h.Class("text-3xl font-semibold text-gray-900"),
						// ),
						// h.Div(h.Class("my-1")),

						h.Table(
							h.Class("table"),
							h.Thead(
								h.Tr(
									h.Th(h.Text("Time")),
									h.Th(h.Text("Repository")),
									h.Th(h.Text("commit message")),
								),
							),

							activities,
						),
					),
				),
			),
		),
	)

}

type UserState struct {
	counterValue int
	time         string
	commits      maybe.Maybe[[]GithubCommit]
}

type GithubCommit struct {
	date     time.Time
	message  string
	repoName string
}

// maybe

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.Class("min-h-screen bg-base-200 antialiased"),
			h.NavBar(),
			h.Div(
				h.Class("flex flex-col px-40 gap-y-10 mt-10"),
				h.Profile(),
				h.Activity(state.commits),
				h.Works(),
			),
		)
	}

	update := func(msg UserMsg, state UserState) (UserState, Cmd[Std, UserMsg]) {
		switch m := msg.(type) {
		case SetGitHubActivity:
			state.commits = maybe.Just(m.commits)
		default:
			// pretty.Println("update: receive ", m)
		}

		return state, None[Std, UserMsg]()
	}

	init := func() (UserState, Cmd[Std, UserMsg]) {
		initAction := GetGitHubActivity[Std]()
		return UserState{commits: maybe.Nothing[[]GithubCommit]()}, initAction
	}

	root := dom.GetWindow().Document().QuerySelector("#app")

	Launch(
		root,
		view,
		update,
		func(std Std) Std { return DefaultStd() },
		init,
	)
}
