package main

import (
	"fmt"
	"math"
	"reflect"
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

type Pos struct {
	x, y, w, h float64
}

type SetNodesPosition struct {
	update, model, runtime, view, diagram Pos
}

func (SetNodesPosition) isUserMsg() {}

type NextAction struct {
}

func (NextAction) isUserMsg() {}

type AppendCmdOrMsg struct {
	cmdOrMsg CmdOrMsg
}

func (AppendCmdOrMsg) isUserMsg() {}

func (h HTML[Msg]) Typ(content string) VNode[Msg] {
	return h.P(
		h.Text(content),
		h.Class("opacity-60 text-xs"),
	)
}

func (h HTML[Msg]) Link(content string, url string) VNode[Msg] {
	return h.A(
		h.Class("link link-hover text-md font-semibold opacity-60 tracking-tight"),
		h.Text(
			content,
		),
	)
}

func (h HTML[M]) Profile() VNode[M] {
	tableEntry := func(key string, value string) VNode[M] {
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
		"プロフィール",
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
							h.Text("hiragi"),
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

func (h HTML[M]) Content(header string, headerSub string, inner VNode[M]) VNode[M] {
	return h.Div(
		h.P(
			h.Text(headerSub),
			h.Class("text-xs mb-1 font-mono font-semibold tracking-widest text-violet-600"),
		),
		h.Hn(2,
			h.Text(header),
			// h.Class("mb-2 text-md font-semibold tracking-tight text-violet-600"),
			h.Class("mb-4 text-md font-semibold tracking-tight"),
		),
		inner,
	)
}

func (h HTML[M]) Works() VNode[M] {
	listRow := func(header VNode[M], title string, description string, link string) VNode[M] {
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
		"作例",
		"WORKS",
		h.Ul(
			h.Class("list bg-base-100 rounded-box shadow-lg overflow-hidden"),
			listRow(
				h.Div(
					h.Class("size-10 rounded-box"),
					h.Img("https://cdn.simpleicons.org/github"),
				),
				"Go-Elm",
				"Go言語用Elmアーキテクチャベースのフロントエンドフレームワーク",
				"https://github.com/hiragi-dev/go-elm",
			),
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

func (h HTML[Msg]) NavBar() VNode[Msg] {
	return h.Div(
		h.Class("navbar bg-base-100 shadow-lg"),
		h.Div(
			h.Class("flex w-full lg:w-240 lg:px-0 px-4 mx-auto justify-between"),
			h.Div(
				h.P(h.Text("hiragi's portfolio"), h.Class("text-xl font-semibold")),
			),
			h.Div(
				h.Class("flex gap-x-4"),
				h.Link("PROFILE", "ok"),
				h.Link("ACTIVITY", "ok"),
				h.Link("WORKS", "ok"),
			),
		),
	)
}

func (h HTML[Msg]) Activity(commits maybe.Maybe[[]GithubCommit]) VNode[Msg] {
	type ColorAttr string

	tableEntry := func(key string, repoName string, value string, colorAttr ColorAttr) VNode[Msg] {
		return h.Tr(
			h.Td(
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

	var activities VNode[Msg]
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

		table := h.Tbody()
		for i, commit := range commits {
			// 表示するコミットは最大7件
			if i >= 7 {
				break
			}

			table.children = append(table.children,
				tableEntry(
					commit.date.Format("2006/1/2 15:04"),
					commit.repoName,
					commit.message,
					repos[commit.repoName],
				),
			)

		}

		activities = h.Table(
			h.Class("table"),
			h.Thead(
				h.Tr(
					h.Th(h.Text("Date"), h.Class("text-aligh; left pl-0")),
					h.Th(h.Text("Repository")),
					h.Th(h.Text("Commit Message")),
				),
			),
			table,
		)
	} else {
		activities = h.Div(
			h.Class("flex justify-center"),
			h.Span(
				h.Class("loading loading-spinner loading-md text-violet-600"),
			),
		)
	}

	return h.Content(
		"アクティビティ",
		"ACTIVITY",
		h.Div(
			h.Class("card bg-base-100 shadow-lg"),
			h.Div(
				h.Class("card-body"),
				h.Div(
					h.Class("flex"),
					h.Div(
						h.Class("ml-8 flex flex-col w-full"),
						activities,
					),
				),
			),
		),
	)

}

func (h HTML[Msg]) UserFooter() VNode[Msg] {
	return h.Footer(
		h.Class("footer bg-base-200 p-10 px-40 flex justify-center"),
		h.Nav(
			h.Typ("© 2026 hiragi. All rights reserved."),
		),
	)
}

type Action interface{ isAction() }

type NewMsg struct{}

func (NewMsg) isAction() {

}

func (h HTML[Msg]) arrow(x1, y1, x2, y2 float64, base Pos) VNode[Msg] {
	dx, dy := x2-x1, y2-y1
	length := math.Hypot(dx, dy)
	angle := math.Atan2(dy, dx) * 180 / math.Pi

	return h.Div(
		h.Class("absolute h-0.5 bg-primary origin-left pointer-events-none"),
		h.Attr("style", fmt.Sprintf(
			"left:%.1fpx;top:%.1fpx;width:%.1fpx;transform:rotate(%.2fdeg)",
			x1-base.x, y1-1-base.y, length, angle, // 太さ 2px なので、中心を y1 に合わせるため 1px 上げる
		)),
	)
}

func PlayGround(action CmdOrMsg, nodesPos NodesPos) VNode[UserMsg] {
	h := HTML[UserMsg]{}

	card := func(label string) VNode[UserMsg] {
		return h.Div(
			// h.Class("")
			h.Attr("id", label),
			h.Div(
				h.Class("card bg-base-200 w-24 h-8"),
				h.Div(
					h.Class("mx-auto my-auto"),
					h.P(
						h.Text(label),
						h.Class("font-semibold"),
					),
				),
				// h.Div(h.Class("h-0 w-0 border-x-8 border-x-transparent border-b-8 border-b-blue-600")),
			),
		)

	}

	// 斜めの矢印。rotate で向きを決め、active のときだけ矢印全体を進む向きに動かす。
	// rotate-* は個別の rotate プロパティ、animate-arrow-slide は transform なので両立し、動きは回転後の向きに沿う
	diagArrow := func(active bool, rotate string, label string) VNode[UserMsg] {
		return h.Span(
			// 外側は回転させない。ラベルを画面上の「上」に置くため
			h.Class("relative inline-flex items-center justify-center"),
			// アニメーション中だけ、矢印の上にラベルを出す。回転した矢印は箱の外にはみ出すので mb-4 で離す。
			// 非アクティブ時も空の span を置き、子の数を変えないようにしている（diff がノードを丸ごと作り直さないため）
			If(active,
				h.Span(
					h.Class("absolute bottom-full left-1/2 mb-4 -translate-x-1/2 badge badge-xs badge-neutral font-mono"),
					h.Text(label),
				),
				h.Span(),
			),
			// 文字の矢印を伸ばすと鏃まで潰れるので、軸（丸めた線）と鏃（2辺だけ枠線を付けて45度回した正方形）を別々に描く
			h.Span(
				h.Class("inline-flex items-center "+rotate+" "+
					If(active, "animate-arrow-slide [--arrow-dx:10px]", "opacity-20")),
				h.Span(h.Class("block h-0.5 w-10 rounded-full bg-black")),
				h.Span(h.Class("block size-2.5 -ml-2.5 rotate-45 rounded-tr-[2px] border-t-2 border-r-2 border-black")),
			),
		)
	}

	// var content VNode[UserMsg]
	var isActionMsg = false
	var isActionCmd = false
	switch a := action.(type) {

	case ActionMsg:
		isActionMsg = true
		_ = a
	case ActionCmd:
		isActionCmd = true
		_ = a

	}

	pretty.Println(isActionCmd, isActionMsg)

	return h.Div(
		h.Class("mockup-window bg-base-100 border border-base-300 shadow-lg"),
		h.Div(
			// h.Class("w-full"),
			h.Class("w-full grid place-content-center border-t border-base-300 h-80"),
			h.Button(
				h.Class("btn mb-5"),
				h.Text("next"),
				OnClick(func() UserMsg { return NextAction{} }),
			),

			//        1        2    3       4    5
			//   r1   .        .    Update  .    .
			//   r2   .        ↗    .       ↘    .
			//   r3   Runtime  .    (now)   .    Model
			//   r4   .        ↖    .       ↙    .
			//   r5   .        .    View    .    .
			h.Div(
				h.Class("grid grid-cols-[repeat(5,6rem)] grid-rows-[repeat(5,2.5rem)] gap-2 place-items-center"),
				h.Div(h.Class("col-start-3 row-start-1"), card("Update")),
				h.Div(h.Class("col-start-1 row-start-3"), card("Runtime")),
				h.Div(h.Class("col-start-5 row-start-3"), card("Model")),
				h.Div(h.Class("col-start-3 row-start-5"), card("View")),

				h.Div(h.Class("col-start-2 row-start-4"), diagArrow(isActionMsg, "rotate-[225deg]", "Msg")), // View → Runtime（↖）
				h.Div(h.Class("col-start-2 row-start-2"), diagArrow(isActionMsg, "-rotate-45", "Msg")),      // Runtime → Update（↗）

				h.Div(h.Class("col-start-4 row-start-2"), diagArrow(isActionCmd, "rotate-45", "Cmd")),       // Update → Model（↘）
				h.Div(h.Class("col-start-4 row-start-4"), diagArrow(isActionCmd, "rotate-[135deg]", "Cmd")), // Model → View（↙）
			),
		),
	)
}

type CmdOrMsg interface{ isCmdOrMsg() }

type ActionCmd struct {
	name string
}

func (ActionCmd) isCmdOrMsg() {}

type ActionMsg struct {
	name string
}

func (ActionMsg) isCmdOrMsg() {}

type NodesPos struct {
	update, model, view, runtime, diagram Pos
}

type UserState struct {
	counterValue  int
	time          string
	commits       maybe.Maybe[[]GithubCommit]
	action        CmdOrMsg
	actionHistory []CmdOrMsg
	nodesPos      NodesPos
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

// func (r UserEffSt) GetActivity(e Http) []GithubCommit {
// 	// r.trace("GetActivity")
// 	return GetGitHubActivity(e)
// }

type githubInfoMock struct {
}

func (r githubInfoMock) GetActivity(e Http) []GithubCommit {
	time.Sleep(time.Second * 1)

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

// func TaggedAsyncDo[E, M any](c Cmd[E, M], wrap func(env E, trace func(string)) E) Cmd[E, M] {
// 	return func(env E, send func(M)) {
// 		traced := wrap(env, func(s string) { send(SetGitHubActivity{}) })
// 	}
// }

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VNode[UserMsg] {
		return h.Div(
			h.Class("min-h-screen bg-base-200 antialiased flex flex-col"),
			h.NavBar(),
			h.Div(
				h.Class("flex grow justify-center"),
				h.Div(
					// h.Class("flex flex-col lg:w-1/2 sm:w-4/5 gap-y-10 my-10"),

					h.Class("flex flex-col gap-y-10 my-10 w-full lg:w-240 lg:px-0 px-4 mx-auto"),
					// lg:px-80
					// h.NavBar(),
					PlayGround(state.action, state.nodesPos),
					h.Profile(),
					h.Activity(state.commits),
					h.Works(),
				),
			),
			h.UserFooter(),
		)
	}

	update := func(msg UserMsg, state UserState) (UserState, Cmd[UserEff, UserMsg]) {
		// pretty.Println("2", state.actionHistory)
		// state.actionHistory = append(state.actionHistory, ActionMsg{name: reflect.TypeOf(msg).String()})

		switch m := msg.(type) {
		case SetGitHubActivity:
			state.commits = maybe.Just(m.commits)
			state.actionHistory = append(state.actionHistory, ActionMsg{name: reflect.TypeOf(msg).String()})

		case NextAction:
			// pretty.Println("next action required")
			state.actionHistory = append(state.actionHistory, ActionMsg{name: reflect.TypeOf(msg).String()})

			frontAction := state.actionHistory[0]
			// pretty.Println("front action", frontAction)
			// 先頭をポップ
			state.actionHistory = state.actionHistory[1:]

			state.action = frontAction

			cmd := AsyncDo(func(e UserEff) UserMsg {
				runtimePos := e.GetPosition("Runtime")
				updatePos := e.GetPosition("Update")
				modelPos := e.GetPosition("Model")
				viewPos := e.GetPosition("View")
				diagram := e.GetPosition("Diagram")
				// pretty.Println(x, y, w, h)

				return SetNodesPosition{runtime: runtimePos, update: updatePos, model: modelPos, view: viewPos, diagram: diagram}
			})

			return state, cmd

		case AppendCmdOrMsg:
			// pretty.Println("AppendCmdOrMsgOrMsg")
			state.actionHistory = append(state.actionHistory, m.cmdOrMsg)
		case SetNodesPosition:
			pretty.Println(m)

			state.nodesPos.model = m.model
			state.nodesPos.view = m.view
			state.nodesPos.update = m.update
			state.nodesPos.runtime = m.runtime
			state.nodesPos.diagram = m.diagram

		default:
			// pretty.Println("update: receive ", m)
		}

		// pretty.Println("1", state.actionHistory)

		return state, None[UserEff, UserMsg]()
	}

	init := func() (UserState, []Cmd[UserEff, UserMsg]) {
		initAction := TaggedAsyncDo[UserEff, UserMsg](
			func(e UserEff) UserMsg {
				commits := e.GetActivity(e)
				return SetGitHubActivity{commits}
			},
			AppendCmdOrMsg{cmdOrMsg: ActionCmd{"GetGitHubActivity"}},
		)

		// cmd := AsyncDo(func(e UserEff) UserMsg {
		// 	runtimePos := e.GetPosition("Runtime")
		// 	updatePos := e.GetPosition("Update")
		// 	modelPos := e.GetPosition("Model")
		// 	viewPos := e.GetPosition("View")
		// 	// pretty.Println(x, y, w, h)

		// 	return SetNodesPosition{runtime: runtimePos, update: updatePos, model: modelPos, view: viewPos}
		// })

		return UserState{
				commits:       maybe.Nothing[[]GithubCommit](),
				action:        ActionCmd{},
				actionHistory: []CmdOrMsg{},
				nodesPos:      NodesPos{},
			},
			[]Cmd[UserEff, UserMsg]{initAction}

	}

	root := dom.GetWindow().Document().QuerySelector("#app")
	go func() {
		time.Sleep(time.Second * 1)
		doc := dom.GetWindow().Document()
		// doc.
		pretty.Println("==============", doc.QuerySelector("#Runtime").GetBoundingClientRect().Left())
	}()

	Launch(
		root,
		view,
		update,
		func(std Std) UserEff { return DefaultUserEff(std) },
		init,
	)
}
