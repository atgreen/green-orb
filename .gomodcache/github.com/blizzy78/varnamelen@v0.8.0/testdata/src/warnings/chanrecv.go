package warnings

func ChanRecvOK() {
	_, ok := <-make(chan int) // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
