package main

import (
	"fmt"
	"time"

	"github.com/kr/pretty"
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
	type ColorAttr string
	// }

	tableEntry := func(key string, repoName string, value string, colorAttr ColorAttr) VDom[Msg] {
		return h.Tr(
			h.Th(
				h.Class("px-0"),
				h.Text(key),
			),
			h.Td(
				h.Div(
					h.Class(fmt.Sprintf("badge badge-soft badge-%s", colorAttr)),
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

		// リポジトリごとに表示色を分ける
		repos := make(map[string]ColorAttr)
		colorAttrs := []ColorAttr{"primary", "secondary", "accent", "info", "success", "warning", "error"}
		acc := 0
		for _, commit := range commits {
			if _, ok := repos[commit.repoName]; ok {
				continue
			} else {
				repos[commit.repoName] = colorAttrs[acc%len(colorAttrs)]
				acc++
			}
		}

		activities = h.Tbody()
		for i, commit := range commits {
			if i < 7 {
				activities.children = append(activities.children,
					tableEntry(
						commit.date.Format("2006/1/2 15:04"),
						commit.repoName,
						commit.message,
						repos[commit.repoName],
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

type GitHubInfo interface {
	GetActivity(Http) []GithubCommit
}

type githubInfo struct {
}

func (r githubInfo) GetActivity(e Http) []GithubCommit {
	return GetGitHubActivity(e)

}

type githubInfoMock struct {
}

func (r githubInfoMock) GetActivity(e Http) []GithubCommit {
	time.Sleep(time.Second * 1)

	// 表示確認用のダミーデータ。リポジトリは全件異なる名前にして、badge の7色が1つずつ当たるようにしている
	now := time.Now()
	return []GithubCommit{
		{date: now.Add(-25 * time.Minute), message: "impl keyed diff for list children", repoName: "go-elm"},
		{date: now.Add(-3 * time.Hour), message: "support SSH_MSG_CHANNEL_REQUEST exec (RFC 4254 §6.5)", repoName: "ssh-server"},
		{date: now.Add(-26 * time.Hour), message: "docs: add wiring diagram for M5Stack GPS module", repoName: "onealarm-doc"},
		{date: now.Add(-2 * 24 * time.Hour), message: "read GPS fix and compare with registered location", repoName: "onealarm-firmware"},
		{date: now.Add(-4 * 24 * time.Hour), message: "add layout for 40% ortholinear with rotary encoder", repoName: "keyboard-layout"},
		{date: now.Add(-9 * 24 * time.Hour), message: "chapter 3: effect handlers and one-shot continuations", repoName: "notes"},
		{date: now.Add(-20 * 24 * time.Hour), message: "initial commit", repoName: "dotfiles"},
	}
}

type UserEff interface {
	Std
	GitHubInfo
}

type UserEffSt struct {
	Std
	githubInfoMock
}

func DefaultUserEff(std Std) UserEff {
	return UserEffSt{std, githubInfoMock{}}
}

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

	update := func(msg UserMsg, state UserState) (UserState, Cmd[UserEff, UserMsg]) {
		switch m := msg.(type) {
		case SetGitHubActivity:
			state.commits = maybe.Just(m.commits)
		default:
			// pretty.Println("update: receive ", m)
		}

		return state, None[UserEff, UserMsg]()
	}

	init := func() (UserState, Cmd[UserEff, UserMsg]) {
		initAction := AsyncDo(func(e UserEff) UserMsg {
			commits := e.GetActivity(e)
			// commits := []GithubCommit{{message: "test", repoName: "reporepo"}}
			pretty.Println(commits)
			// return SetGitHubActivity{commits: commits}
			return SetGitHubActivity{commits}
		})

		return UserState{commits: maybe.Nothing[[]GithubCommit]()}, initAction
	}

	root := dom.GetWindow().Document().QuerySelector("#app")

	Launch(
		root,
		view,
		update,
		func(std Std) UserEff { return DefaultUserEff(std) },
		init,
	)
}
