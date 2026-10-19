package main

import (
	"reflect"
)

type AnyOp interface{ isOp() }

type Op[R any] interface {
	AnyOp
	isEff(R)
}

type HttpGet struct {
	url     string
	headers []string
}

func (HttpGet) isEff(AnyOp) {}

func (HttpGet) isOp() {}

type Ctl[R any] interface {
	isCtl(R)
}

type Resumption[R any] struct {
	retVal R
}

func (Resumption[R]) isCtl(R) {}

type Handlers map[reflect.Type]func(any) any

func (hs Handlers) On[E Op[R], R any](f func(E) R) {
	hs[reflect.TypeFor[E]()] = func(a any) any {
		return f(a.(E))
	}
}

func SystemHandler() Handlers {
	handlers := Handlers{}

	// handlers.On(func(arg HttpGet) string {
	// 	pretty.Println("handler got ", arg)
	// 	return "handler return value"
	// })

	return handlers
}

// func SystemHandler[R any](eff Op[R]) Ctl[R] {
// 	switch op := eff.(type) {
// 	case HttpGet:
// 		pretty.Println("SystemHandler Get ", op)
// 		return Resumption[string]{retVal: ""}
// 	default:
// 		return Resumption{}
// 	}
// }
