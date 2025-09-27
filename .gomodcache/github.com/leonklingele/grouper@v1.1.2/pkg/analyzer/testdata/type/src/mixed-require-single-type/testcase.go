package testcase

type (
	a string
	b = string

	_ = string
)

type ( // want "should only use a single global 'type' declaration, 5 found"
	comment1 = string // some comment
)

type c = string

type comment2 string // some comment

type ()

func dummy() {
	type (
		_ string
		_ = string
	)

	type (
		_ string
	)

	type d string
	type comment3 = string // some comment

	type ()
}
