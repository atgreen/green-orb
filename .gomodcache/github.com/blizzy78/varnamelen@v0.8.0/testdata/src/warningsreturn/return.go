package warnings

func Return() (i int) { // want `return value name 'i' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	i = 123
	return
}
