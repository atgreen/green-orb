package addon

//revive:disable
//TODO: remove foo (Line 3)
var foo = "Hello World" // line comment 2

type X struct {
	Name string //TODO: Rename field (Line 7)
}

/*TODO: get cat food (Line 10)*/
// This comment is associated with the main function.
func New() *X {
	return &X{}
	// ignored line
} //  todo  : todo comment (Line 15)
