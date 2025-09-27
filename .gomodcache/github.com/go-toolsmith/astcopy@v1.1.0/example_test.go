package astcopy_test

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/go-toolsmith/astcopy"
	"github.com/go-toolsmith/astequal"
	"github.com/go-toolsmith/strparse"
)

func Example() {
	x := strparse.Expr(`1 + 2`).(*ast.BinaryExpr)
	y := astcopy.BinaryExpr(x)
	fmt.Println(astequal.Expr(x, y)) // => true

	// Now modify x and make sure y is not modified.
	z := astcopy.BinaryExpr(y)
	x.Op = token.SUB
	fmt.Println(astequal.Expr(y, z)) // => true
	fmt.Println(astequal.Expr(x, y)) // => false

	// Output:
	// true
	// true
	// false
}
