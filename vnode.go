package main

import (
	"strconv"
)

type Applier[M any] interface{ apply(*VNode[M]) }

type VNode[Msg any] struct {
	tagName    string
	attributes []Attribute[Msg]
	event      map[string]func() Msg
	children   []VNode[Msg]
	text       string
}

func (v VNode[M]) apply(p *VNode[M]) {
	p.children = append(p.children, v)
}

type Attribute[M any] struct {
	key   string
	value string
}

func (a Attribute[M]) apply(p *VNode[M]) {
	p.attributes = append(p.attributes, a)
}

type Event[M any] struct {
	key   string
	value func() M
}

func (e Event[M]) apply(p *VNode[M]) {
	p.event[e.key] = e.value
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

func (h HTML[Msg]) Div(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("div").Apply(applier...)
}

func (h HTML[Msg]) P(content VNode[Msg], applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("p").Apply(content).Apply(applier...)
}

func (h HTML[Msg]) Ul(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("ul").Apply(applier...)
}

func (h HTML[Msg]) Li(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("li").Apply(applier...)
}

func (h HTML[Msg]) Hn(n uint, applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("h" + strconv.Itoa(int(n))).Apply(applier...)
}

func (h HTML[Msg]) A(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("a").Apply(applier...)
}

func (h HTML[Msg]) Table(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("table").Apply(applier...)
}

func (h HTML[Msg]) Thead(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("thead").Apply(applier...)
}

func (h HTML[Msg]) Tbody(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("tbody").Apply(applier...)
}

func (h HTML[Msg]) Tr(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("tr").Apply(applier...)
}

func (h HTML[Msg]) Th(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("th").Apply(applier...)
}

func (h HTML[Msg]) Td(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("td").Apply(applier...)
}

func (h HTML[Msg]) Span(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("span").Apply(applier...)
}

func (h HTML[Msg]) Text(content string, applier ...Applier[Msg]) VNode[Msg] {
	return *CreateTextNode[Msg](content).Apply(applier...)
}

func (h HTML[Msg]) Button(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("button").Apply(applier...)
}

func (h HTML[Msg]) Img(srcUrl string) VNode[Msg] {
	return *CreateVNode[Msg]("img").SetAttribute("src", srcUrl)
}

func (h HTML[Msg]) Footer(applier ...Applier[Msg]) VNode[Msg] {
	return *CreateVNode[Msg]("footer").Apply(applier...)
}

func (h HTML[Msg]) Attr(key string, value string) Attribute[Msg] {
	return Attribute[Msg]{key, value}
}

func (h HTML[Msg]) Class(value string) Attribute[Msg] {
	return h.Attr("class", value)
}

func CreateVNode[M any](tagName string) *VNode[M] {
	return &VNode[M]{
		tagName:    tagName,
		children:   make([]VNode[M], 0),
		event:      map[string]func() M{},
		attributes: make([]Attribute[M], 0),
		text:       "",
	}
}

func CreateTextNode[M any](content string) *VNode[M] {
	el := CreateVNode[M]("")
	el.text = content
	return el
}

func (vnode *VNode[M]) SetAttribute(key string, value string) *VNode[M] {
	vnode.attributes = append(vnode.attributes, Attribute[M]{key, value})

	return vnode
}

func (vnode *VNode[M]) Apply(applier ...Applier[M]) *VNode[M] {
	for _, a := range applier {
		a.apply(vnode)
	}

	return vnode
}
