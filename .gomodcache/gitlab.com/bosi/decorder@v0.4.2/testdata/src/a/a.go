package a

type aa int
type ab int // want "multiple \"type\" declarations are not allowed; use parentheses instead"

const ac = 1
const ad = 1 // want "multiple \"const\" declarations are not allowed; use parentheses instead"

var ae = 1
var af = 1 // want "multiple \"var\" declarations are not allowed; use parentheses instead"
var _ = 1  // want "multiple \"var\" declarations are not allowed; use parentheses instead"

func ag() {
	type h int
	const i = 1
	var j = 1
	_ = j
}
