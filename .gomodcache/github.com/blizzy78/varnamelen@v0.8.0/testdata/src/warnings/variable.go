package warnings

func Variable_Assign() {
	x := 123 // want `variable name 'x' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = x
}

func Variable_Assign_Range() {
	for i := range []int{1, 2, 3} { // want `variable name 'i' is too short for the scope of its usage`
		// fill
		// fill
		// fill
		// fill
		// fill
		_ = i
	}
}

func Variable_Assign_Range_2() {
	for _, i := range []int{1, 2, 3} { // want `variable name 'i' is too short for the scope of its usage`
		// fill
		// fill
		// fill
		// fill
		// fill
		_ = i
	}
}

func Variable_ValueSpec() {
	var x = 123 // want `variable name 'x' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = x
}

func Variable_TypeSwitch() {
	switch x := interface{}(1).(type) { // want `variable name 'x' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	case int:
		_ = x
	}
}
