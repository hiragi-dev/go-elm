func render(view View, update func(PartialGenericMsg, GenericState) GenericState)

func update(msg PartialGenericMsg, state GenericState) GenericState {

}

func view(state GenericState) View {

}

func main() {
	app.render(view, update)

	select {}
}