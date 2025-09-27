package warnings

func MapIndexOK() {
	_, ok := map[int]int{}[0] // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
