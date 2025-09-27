package testcase

const (
	a = "a"
	b = "b"

	_ = "underscore1"
)

func dummy() {
	const (
		_ = "ignore1"
		_ = "ignore2"
	)
}
