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

func applyPatch[M any](node dom.Node, diff []Patch[M], msgChan chan M) {
	for _, d := range diff {
		var target = node
		for _, i := range d.path {
			target = target.ChildNodes()[i]
		}

		newDom := createElement(d.newValue, msgChan)
		target.ParentNode().ReplaceChild(newDom, target)
	}
}
