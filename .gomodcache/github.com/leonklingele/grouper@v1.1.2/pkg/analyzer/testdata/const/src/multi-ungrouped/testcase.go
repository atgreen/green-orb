package testcase

const a = "a" // want "should only use grouped global 'const' declarations"
const b = "b" // want "should only use a single global 'const' declaration, 3 found"

const _ = "underscore1"

func dummy() {
	const _ = "ignore1"
	const _ = "ignore2"
}
