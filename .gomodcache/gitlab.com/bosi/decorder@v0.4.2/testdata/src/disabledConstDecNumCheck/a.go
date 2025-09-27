package disabledDecNumCheck

type aa int
type ab int // want "multiple \"type\" declarations are not allowed; use parentheses instead"

const ac = 1
const ad = 1

var ae = 1
var af = 1 // want "multiple \"var\" declarations are not allowed; use parentheses instead"
