package main

func If[V any](cond bool, t V, f V) V {
	if cond {
		return t
	} else {
		return f
	}
}
