package main

import (
	"slices"

	"honnef.co/go/js/dom/v2"
)

type Patch[M any] struct {
	path     []int
	newValue VNode[M]
}

func diff[M any](oldVNode VNode[M], newVNode VNode[M]) []Patch[M] {
	out := make([]Patch[M], 0)
	_diff(oldVNode, newVNode, make([]int, 1), &out)

	return out
}

func _diff[M any](oldVNode VNode[M], newVNode VNode[M], path []int, out *[]Patch[M]) {
	if oldVNode.tagName != newVNode.tagName || len(oldVNode.children) != len(newVNode.children) {
		*out = append(*out, Patch[M]{path: path, newValue: newVNode})
		return
	}

	if oldVNode.text != newVNode.text {
		*out = append(*out, Patch[M]{path: path, newValue: newVNode})
		return
	}

	for i := range oldVNode.children {
		childPath := append(slices.Clone(path), i)

		_diff(oldVNode.children[i], newVNode.children[i], childPath, out)
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
