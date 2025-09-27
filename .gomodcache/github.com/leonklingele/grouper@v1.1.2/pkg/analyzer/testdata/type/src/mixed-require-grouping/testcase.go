package testcase

type (
	a string
	b = string

	_ = string
)

type (
	comment1 = string // some comment
)

type c = string // want "should only use grouped global 'type' declarations"

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
