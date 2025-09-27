package testcase

var (
	a = "a"
	b = "b"

	_ = "underscore1"
)

var ( // want "should only use a single global 'var' declaration, 5 found"
	comment1 = "comment1" // some comment
)

const ignoreconst1 = "ignoreconst1"
const (
	ignoreconst2 = "ignoreconst2"
	ignoreconst3 = "ignoreconst3"
)

func dummy() {
	var (
		ignorea = "ignorea"
		ignoreb = "ignoreb"

		_ = "ignoreunderscore1"
	)

	var (
		ignorecomment1 = "ignorecomment1" // some comment
	)

	var ignorec = "ignorec"
	var _ = "ignoreunderscore2"

	var ()

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

var c = "c" // want "should only use grouped global 'var' declarations"

var comment2 = "comment2" // some comment

var ()
