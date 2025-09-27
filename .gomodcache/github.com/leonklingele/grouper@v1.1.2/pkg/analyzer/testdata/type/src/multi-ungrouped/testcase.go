package testcase

type a string   // want "should only use grouped global 'type' declarations"
type b = string // want "should only use a single global 'type' declaration, 3 found"

type _ string

func dummy() {
	type _ string
	type _ = string
}
