package warnings

type compositeLit struct {
	x int
}

func CompositeLiteral() {
	i := 123 // want `variable name 'i' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = compositeLit{
		x: i,
	}
}
