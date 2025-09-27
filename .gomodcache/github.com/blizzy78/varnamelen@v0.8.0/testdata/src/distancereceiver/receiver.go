package warningsreceiver

type foo struct{}

func (f *foo) Receiver() {
	_ = f
}
