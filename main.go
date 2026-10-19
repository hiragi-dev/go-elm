package main

import (
	"github.com/kr/pretty"
	"honnef.co/go/js/dom/v2"
)

type UserMsg interface{ isUserMsg() }

type GetGitHubActivity struct{ response string }

func (GetGitHubActivity) isUserMsg() {}

type GotGitHubActivity struct{ response string }

func (GotGitHubActivity) isUserMsg() {}

type UserState struct {
	counterValue int
	time         string
	commits      []string
}

type Eff struct {
	req chan Op[AnyOp]
	res chan any
}

// func (Eff) None() {
// 	return Eff{req: make(chan Op), res: make(chan any)}
// }

func Handle[T, R any](program func(*Eff)) any {
	eff := &Eff{req: make(chan Op[AnyOp]), res: make(chan any)}

	go func() {
		pretty.Println("run program")
		program(eff)
	}()

	println("wait")
	for op := range eff.req {
		pretty.Println("handle op: ", op)
		eff.res <- "res"
	}

	return nil
}

func Perform[R any](eff *Eff, op Op[R]) R {
	pretty.Println("Perform: ", op)

	eff.req <- op
	res := <-eff.res

	return res
}

func main() {
	h := HTML[UserMsg]{}

	view := func(state UserState) VDom[UserMsg] {
		return h.Div(
			h.Button(
				h.P(
					h.Text("test"),
				),
				OnClick(func() UserMsg { return GetGitHubActivity{} }),
			),
		)
	}

	update := func(e *Eff, msg UserMsg, state UserState) UserState {
		switch m := msg.(type) {
		default:
			res := e.Perform(HttpGet{})

			pretty.Println("update: receive ", m)
			pretty.Println("and handled operation get ", res)
		}

		return state
	}

	root := dom.GetWindow().Document().QuerySelector("#app")

	init := func(e *Eff[any, any]) UserState {
		return UserState{counterValue: 0, commits: make([]string, 0)}
	}

	Launch(root, view, update, init)
}
