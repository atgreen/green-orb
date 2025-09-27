package strparse_test

import (
	"fmt"

	"github.com/go-toolsmith/astequal"
	"github.com/go-toolsmith/strparse"
)

func Example() {
	// Comparing AST strings for equallity (note different spacing):
	x := strparse.Expr(`1 + f(v[0].X)`)
	y := strparse.Expr(` 1+f( v[0].X ) `)
	fmt.Println(astequal.Expr(x, y)) // => true

	// Output:
	// true
}
