package main

import "strconv"

type VDom[Msg any] struct {
	tagName    string
	attributes map[string]string
	event      map[string]func() Msg
	children   []VDom[Msg]
	text       string
}

type Msg interface{ isMsg() }

type EmptyMsg struct{}

func (EmptyMsg) isMsg() {}

// VNodeはメッセージの型がジェネリクスになっているが、goの型推論では各VNodeごとにMsgの型注釈をする必要があり煩雑なので
// HTML[Msg]なるビルダーっぽいものを用意して、型注釈を一箇所で済ませられるようにする
// ↓ ビルダーが無いとこうなってめんどい
// Div[UserMsg](Button(func() UserMsg { return UserMsg{} }, P[UserMsg]("click me!")))
type HTML[Msg any] struct {
}

func (h HTML[Msg]) Div(children ...VDom[Msg]) VDom[Msg] {
	return VDom[Msg]{tagName: "div", attributes: make(map[string]string), children: children}
}

func (h HTML[Msg]) P(content VDom[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)
	children = append(children, content)

	return VDom[Msg]{tagName: "p", children: children}
}

func (h HTML[Msg]) Hn(n uint, content string) VDom[Msg] {
	children := make([]VDom[Msg], 0)
	children = append(children, h.Text(content))

	return VDom[Msg]{tagName: "h" + strconv.Itoa(int(n)), children: children}
}

func (h HTML[Msg]) Text(content string) VDom[Msg] {
	return VDom[Msg]{text: content}
}

func (h HTML[Msg]) Button(onClick func() Msg, children ...VDom[Msg]) VDom[Msg] {
	event := make(map[string]func() Msg)
	event["click"] = onClick

	return VDom[Msg]{tagName: "button", children: children, event: event}
}
