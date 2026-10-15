func render(view View, handler func(PartialGenericMsg, GenericState) GenericState)

func handler(msg PartialGenericMsg, state GenericState) GenericState {

}

func view(state GenericState) View {

}

func main() {
	app.render(view, handler)

	select {}
}