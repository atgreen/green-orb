package testcase

const (
	a = "a"
	b = "b"

	_ = "underscore1"
)

const ( // want "should only use a single global 'const' declaration, 5 found"
	comment1 = "comment1" // some comment
)

var ignorevar1 string
var (
	ignorevar2 string
	ignorevar3 string
)

func dummy() {
	const (
		ignorea = "ignorea"
		ignoreb = "ignoreb"

		_ = "ignoreunderscore1"
	)

	const (
		ignorecomment1 = "ignorecomment1" // some comment
	)

	const ignorec = "ignorec"
	const _ = "ignoreunderscore2"

	const ()

	_ = a
	_ = b
	_ = c
	_ = comment1
	_ = comment2
	_ = ignorecomment1
	_ = ignorea
	_ = ignoreb
	_ = ignorec
}

const c = "c" // want "should only use grouped global 'const' declarations"

const comment2 = "comment2" // some comment

const ()
