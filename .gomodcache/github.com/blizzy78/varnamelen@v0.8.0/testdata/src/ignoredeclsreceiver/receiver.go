package ignoredeclsreceiver

type foo struct{}

func (f foo) Receiver_Struct() {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = f
}

func (f *foo) Receiver_Pointer() {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = f
}
