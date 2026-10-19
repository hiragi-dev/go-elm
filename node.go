package main

import "strconv"

type Applier[M any] interface{ apply(*VDom[M]) }

type VDom[Msg any] struct {
	tagName    string
	attributes []Attribute[Msg]
	event      map[string]func() Msg
	children   []VDom[Msg]
	text       string
}

func (v VDom[M]) apply(p *VDom[M]) {
	p.children = append(p.children, v)
}

type Attribute[M any] struct {
	key   string
	value string
}

func (a Attribute[M]) apply(p *VDom[M]) {
	p.attributes = append(p.attributes, a)
}

type Event[M any] struct {
	key   string
	value func() M
}

func (e Event[M]) apply(p *VDom[M]) {
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

func (h HTML[Msg]) Div(applier ...Applier[Msg]) VDom[Msg] {
	event := make(map[string]func() Msg)
	children := make([]VDom[Msg], 0)
	attributes := make([]Attribute[Msg], 0)

	el := VDom[Msg]{tagName: "div", children: children, attributes: attributes, event: event}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) P(content VDom[Msg], applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)
	children = append(children, content)

	el := VDom[Msg]{tagName: "p", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Ul(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "ul", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Li(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "li", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Hn(n uint, applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "h" + strconv.Itoa(int(n)), children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) A(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "a", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Table(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "table", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Thead(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "thead", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Tbody(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "tbody", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Tr(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "tr", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Th(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "th", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Td(applier ...Applier[Msg]) VDom[Msg] {
	children := make([]VDom[Msg], 0)

	el := VDom[Msg]{tagName: "td", children: children}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Span(applier ...Applier[Msg]) VDom[Msg] {
	el := VDom[Msg]{tagName: "span"}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Text(content string, applier ...Applier[Msg]) VDom[Msg] {
	el := VDom[Msg]{text: content}

	for _, a := range applier {
		a.apply(&el)
	}

	return el
}

func (h HTML[Msg]) Button(applier ...Applier[Msg]) VDom[Msg] {
	event := make(map[string]func() Msg)
	children := make([]VDom[Msg], 0)
	attributes := make([]Attribute[Msg], 0)

	vdom := VDom[Msg]{tagName: "button", children: children, attributes: attributes, event: event}

	for _, a := range applier {
		a.apply(&vdom)
	}

	return vdom
}

func (h HTML[Msg]) Img(srcUrl string) VDom[Msg] {
	srcAttr := h.Attr("src", srcUrl)

	attributes := make([]Attribute[Msg], 0)
	attributes = append(attributes, srcAttr)

	vdom := VDom[Msg]{tagName: "img", attributes: attributes}

	return vdom
}

func (h HTML[Msg]) Attr(key string, value string) Attribute[Msg] {
	return Attribute[Msg]{key, value}
}

func (h HTML[Msg]) Class(value string) Attribute[Msg] {
	return h.Attr("class", value)
}
