package disabledDecNumCheck

type aa int
type ab int

const ac = 1
const ad = 1 // want "multiple \"const\" declarations are not allowed; use parentheses instead"

var ae = 1
var af = 1 // want "multiple \"var\" declarations are not allowed; use parentheses instead"
