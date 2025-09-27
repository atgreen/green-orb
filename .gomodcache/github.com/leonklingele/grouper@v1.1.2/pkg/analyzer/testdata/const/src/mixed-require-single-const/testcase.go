package testcase

const (
	a = "a"
	b = "b"

	_ = "underscore1"
)

const ( // want "should only use a single global 'const' declaration, 5 found"
	comment1 = "comment1" // some comment
)

const c = "c"

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
