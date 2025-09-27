package disabledDecNumCheck

type aa int
type ab int // want "multiple \"type\" declarations are not allowed; use parentheses instead"

const ac = 1
const ad = 1 // want "multiple \"const\" declarations are not allowed; use parentheses instead"

var ae = 1
var af = 1
