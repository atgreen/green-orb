package warnings

func Param(i int) { // want `parameter name 'i' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Param_FuncLit() {
	_ = func(i int) { // want `parameter name 'i' is too short for the scope of its usage`
		// fill
		// fill
		// fill
		// fill
		// fill
		_ = i
	}
}
