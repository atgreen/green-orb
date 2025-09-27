package astcast_test

import (
	"fmt"

	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/strparse"
)

func Example() {
	x := strparse.Expr(`(foo * bar) + 1`)

	// x type is ast.Expr, we want to access bar operand
	// that is a RHS of the LHS of the addition.
	// Note that addition LHS (X field) is has parenthesis,
	// so we have to remove them too.

	add := astcast.ToBinaryExpr(x)
	mul := astcast.ToBinaryExpr(astcast.ToParenExpr(add.X).X)
	bar := astcast.ToIdent(mul.Y)
	fmt.Printf("%T %s\n", bar, bar.Name) // => *ast.Ident bar

	// If argument has different dynamic type,
	// non-nil sentinel object of requested type is returned.
	// Those sentinel objects are exported so if you need
	// to know whether it was a nil interface value of
	// failed type assertion, you can compare returned
	// object with such a sentinel.

	y := astcast.ToCallExpr(strparse.Expr(`x`))
	if y == astcast.NilCallExpr {
		fmt.Println("it is a sentinel, type assertion failed")
	}

	// Output:
	// *ast.Ident bar
	// it is a sentinel, type assertion failed
}
