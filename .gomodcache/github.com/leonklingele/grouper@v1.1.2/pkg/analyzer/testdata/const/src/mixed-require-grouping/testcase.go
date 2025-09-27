package testcase

const (
	a = "a"
	b = "b"

	_ = "underscore1"
)

const (
	comment1 = "comment1" // some comment
)

const c = "c" // want "should only use grouped global 'const' declarations"

const comment2 = "comment2" // some comment

const ()

func dummy() {
	const (
		_ = "ignore1"
		_ = "ignore2"
	)

	const (
		_ = "ignore3"
	)

	const d = "d"
	const comment3 = "comment3" // some comment

	const ()
}
