package main

type VNode[Msg any] interface{ isVNode(Msg) }

type VDiv[Msg any] struct {
	children []VNode[Msg]
}

type VP[Msg any] struct {
	content string
}

type VButton[Msg any] struct {
	children []VNode[Msg]
	onClick  func() Msg
}

// VNodeはメッセージの型がジェネリクスになっているが、goの型推論では各VNodeごとにMsgの型注釈をする必要があり煩雑なので
// HTML[Msg]なるビルダーっぽいものを用意して、型注釈を一箇所で済ませられるようにする
// ↓ ビルダーが無いとこうなってめんどい
// Div[UserMsg](Button(func() UserMsg { return UserMsg{} }, P[UserMsg]("click me!")))
type HTML[Msg any] struct {
}

func (VDiv[Msg]) isVNode(Msg)    {}
func (VP[Msg]) isVNode(Msg)      {}
func (VButton[Msg]) isVNode(Msg) {}

func (h HTML[Msg]) Div(children ...VNode[Msg]) VNode[Msg] {
	return VDiv[Msg]{children}
}

func (h HTML[Msg]) P(content string) VNode[Msg] {
	return VP[Msg]{content}
}

func (h HTML[Msg]) Button(onClick func() Msg, children ...VNode[Msg]) VNode[Msg] {
	return VButton[Msg]{children, onClick}
}
