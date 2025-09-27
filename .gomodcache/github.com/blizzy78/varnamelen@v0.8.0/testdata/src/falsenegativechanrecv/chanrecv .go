package falsenegativechanrecv

func ChanRecvOK_Name() {
	_, o := <-make(chan int) // want `variable name 'o' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = o
}

func ChanRecvOK_Not2() {
	ok := <-make(chan int) // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
