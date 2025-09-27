package testcase

var a = "a" // want "should only use grouped global 'var' declarations"
var b = "b" // want "should only use a single global 'var' declaration, 3 found"

var _ = "underscore1"

func dummy() {
	var _ = "ignore1"
	var _ = "ignore2"

	println(a)
	println(b)
}
