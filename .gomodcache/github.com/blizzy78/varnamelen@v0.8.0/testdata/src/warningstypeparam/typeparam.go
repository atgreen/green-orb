package warningstypeparam

func TypeParam[T any]() { // want `type parameter name 'T' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	var t T
	_ = t
}
