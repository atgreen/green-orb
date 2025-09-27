package falsenegativetypeassert

func TypeAssertOK_Name() {
	_, o := interface{}(1).(int) // want `variable name 'o' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = o
}

func TypeAssertOK_Not2() {
	ok := interface{}(1).(int) // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
