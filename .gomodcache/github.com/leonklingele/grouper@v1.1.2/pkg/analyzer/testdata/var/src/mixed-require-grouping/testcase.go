package testcase

var (
	a = "a"
	b = "b"

	_ = "underscore1"
)

var (
	comment1 = "comment1" // some comment
)

var c = "c" // want "should only use grouped global 'var' declarations"

var comment2 = "comment2" // some comment

var ()

func dummy() {
	var (
		_ = "ignore1"
		_ = "ignore2"
	)

	var (
		_ = "ignore3"
	)

	var d = "d"
	var comment3 = "comment3" // some comment

	println(d)
	println(comment3)

	var ()
}
