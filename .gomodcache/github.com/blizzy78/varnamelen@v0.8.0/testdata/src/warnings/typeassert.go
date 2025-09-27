package warnings

func TypeAssertOK() {
	_, ok := interface{}(1).(int) // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
