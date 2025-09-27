package disabledDecNumCheck

type (
	ba int
	bb int
)

const (
	bc = 1
	bd = 1
)

var (
	be = 1
	bf = 1
)

func bg() {}

func init() {} // want "init func must be the first function in file"

func (_ ba) init() {}
