package warningstypeparam

func TypeParam[T any]() {
	var t T
	_ = t
}
