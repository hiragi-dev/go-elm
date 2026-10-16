package main

import (
	"slices"

	"honnef.co/go/js/dom/v2"
)

type Patch[M any] struct {
	path     []int
	newValue VDom[M]
}

func diff[M any](oldVDom VDom[M], newVDom VDom[M]) []Patch[M] {
	out := make([]Patch[M], 0)
	_diff(oldVDom, newVDom, make([]int, 1), &out)

	return out
}

func _diff[M any](oldVDom VDom[M], newVDom VDom[M], path []int, out *[]Patch[M]) {
	if oldVDom.tagName != newVDom.tagName || len(oldVDom.children) != len(newVDom.children) {
		*out = append(*out, Patch[M]{path: path, newValue: newVDom})
		return
	}

	if oldVDom.text != newVDom.text {
		*out = append(*out, Patch[M]{path: path, newValue: newVDom})
		return
	}

	for i := range oldVDom.children {
		childPath := append(slices.Clone(path), i)

		_diff(oldVDom.children[i], newVDom.children[i], childPath, out)
	}
}

func patch[M any](node dom.Node, diff []Patch[M], msgChan chan M) {
	for _, d := range diff {
		var target = node
		// var parent dom.Node
		for _, i := range d.path {
			// parent = target
			// pretty.Println(i, target.ChildNodes())
			// println(i)
			// js.Global().Get("console").Call("log", "target", target.Underlying())
			target = target.ChildNodes()[i]
		}

		// pretty.Println("tareget", target.NodeName())
		// js.Global().Get("console").Call("log", "target", target.Underlying())
		// target.Underlying()
		// target = createElement(d.newValue, msgChan)
		// target.AppendChild()
		// if target.
		newDom := createElement(d.newValue, msgChan)
		// pretty.Println(newDom)
		// pretty.Println("hello, ")
		// js.Global().Get("console").Call("log", "target", newDom.Underlying())
		target.ParentNode().ReplaceChild(newDom, target)
		// parent.ReplaceChild(target, createElement(d.newValue, msgChan))
		// target.ParentNode().RemoveChild(createElement(d.newValue, msgChan))
	}
}
