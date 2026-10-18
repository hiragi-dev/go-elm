package main

func OnClick[M any](h func() M) Event[M] {
	return Event[M]{key: "click", value: h}
}
