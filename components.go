package main

type UI[M any] struct{}

func (UI[M]) Btn(onClick func() M, applier ...Applier[M]) VDom[M] {
	h := HTML[M]{}
	all := append([]Applier[M]{OnClick(onClick), h.Attr("class", "inline-flex h-12 items-center justify-center rounded-md bg-neutral-950 px-6 font-medium text-neutral-50 transition active:scale-110")}, applier...)
	return h.Button(all...)
}
