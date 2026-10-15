package main

import "github.com/kr/pretty"

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	"honnef.co/go/js/dom/v2"
// )

type UserMsg struct {
}

func main() {

	h := HTML[UserMsg]{}

	vnode := h.Button(func() UserMsg { return UserMsg{} }, h.P("click me"))

	pretty.Println(vnode)
}

// type Msg struct {
// 	Tag string
// }

// // type VNode struct {
// // 	Tag      string
// // 	Children []*VNode
// // 	Text     string
// // 	OnClick  func() Msg
// // }

// type Applier interface{ apply(*VNode) }

// func Div(children ...Applier) *VNode {
// 	n := &VNode{Tag: "Div", Text: ""}

// 	for _, c := range children {
// 		c.apply(n)
// 	}

// 	return n
// }

// func Text(content string) *VNode {
// 	return &VNode{Tag: "Text", Text: content}
// }

// func Button(onClick func() Msg, children ...Applier) *VNode {
// 	n := &VNode{Tag: "Button", OnClick: onClick}

// 	for _, c := range children {
// 		c.apply(n)
// 	}

// 	return n
// }

// func (c *VNode) apply(p *VNode) {
// 	p.Children = append(p.Children, c)
// }

// type State struct {
// 	value int
// }

// var event_handler_l func(Msg)
// var init_flag bool = true

// func createElement(v *VNode) dom.Node {
// 	var retVal dom.Node

// 	if v.Tag == "Text" {
// 		return dom.GetWindow().Document().CreateTextNode(v.Text)
// 	} else {
// 		if v.Tag == "Div" || v.Tag == "Button" {
// 			p := dom.GetWindow().Document().CreateElement(strings.ToLower(v.Tag))

// 			if v.Tag == "Button" {
// 				p.AddEventListener("click", false, func(e dom.Event) {
// 					if v.OnClick != nil {
// 						m := v.OnClick()
// 						event_handler_l(m)
// 						// dispatcher(m, state)
// 					}
// 				})
// 			}

// 			for _, c := range v.Children {
// 				e := createElement(c)
// 				p.AppendChild(e)
// 			}

// 			// return p
// 			retVal = p
// 		} else {
// 			// return dom.GetWindow().Document().CreateElement(v.Tag)
// 			retVal = dom.GetWindow().Document().CreateElement(v.Tag)
// 		}
// 	}

// 	// if init_flag {
// 	// 	init_flag = false
// 	// 	event_handler_l(Msg{Tag: "init"})
// 	// }

// 	return retVal
// }

// func clear(root dom.Element) {
// 	for _, c := range root.ChildNodes() {
// 		root.RemoveChild(c)
// 	}
// }

// type Var struct {
// 	uuid int
// }

// var vars map[int]int = make(map[int]int)
// var uuid_acc = 0

// func createVar(init int) *Var {
// 	uuid := uuid_acc
// 	vars[uuid_acc] = init
// 	uuid_acc += 1

// 	return &Var{uuid: uuid}
// }

// func setVar(va *Var, v int) {
// 	var_ref := vars[va.uuid]
// 	if var_ref != 0 {
// 		vars[va.uuid] = v
// 	}
// }

// func getVar(va *Var) int {
// 	return vars[va.uuid]
// }

// func main() {
// 	vnodes := func(appState State) *VNode {
// 		return Div(
// 			Button(func() Msg {
// 				return Msg{Tag: "onClick"}
// 			}, Text("Click Me!")),
// 			Text("Counter"+strconv.Itoa(appState.value)),
// 		)
// 	}

// 	b, _ := json.MarshalIndent(vnodes, "", "  ")
// 	fmt.Println(string(b))

// 	doc := dom.GetWindow().Document()
// 	root := doc.QuerySelector("#app")

// 	appState := State{value: 999}

// 	event_handler_l = func(m Msg) {
// 		println("event handle")

// 		if m.Tag == "onClick" {
// 			appState.value = appState.value + 1
// 		}

// 		nextVNode := vnodes(appState)
// 		clear(root)
// 		root.AppendChild(createElement(nextVNode))
// 	}

// 	event_handler_l(Msg{Tag: "init"})

// 	select {}
// }

// func render(appState State, vnode func(State) *VNode) {
// 	vnode(appState)
// }
