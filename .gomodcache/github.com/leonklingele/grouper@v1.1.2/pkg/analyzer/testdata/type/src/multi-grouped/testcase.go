package testcase

type (
	a string
	b = string

	_ = string
)

func dummy() {
	type (
		_ string
		_ = string
	)
}
