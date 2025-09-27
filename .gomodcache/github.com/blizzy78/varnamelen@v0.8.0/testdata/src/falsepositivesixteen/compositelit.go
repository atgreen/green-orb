package falsepositive

import "testing"

func CompositeLiteral() {
	x := 123
	_ = x
	// fill
	// fill
	// fill
	// fill
	// fill
	i := 123
	_ = struct{ x int }{x: i}
}

func CompositeLiteral_Conventional() {
	t := &testing.T{}
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = struct{ t testing.TB }{t: t}
}
