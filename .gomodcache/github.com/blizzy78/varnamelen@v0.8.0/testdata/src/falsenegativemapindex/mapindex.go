package falsenegativemapindex

func MapIndexOK_Name() {
	_, o := map[int]int{}[0] // want `variable name 'o' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = o
}

func MapIndexOK_Not2() {
	ok := map[int]int{}[0] // want `variable name 'ok' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ok
}
