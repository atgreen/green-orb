package astequal_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"

	"github.com/go-toolsmith/astequal"
)

func Example() {
	const code = `
		package foo

		func main() {
			x := []int{1, 2, 3}
			x := []int{1, 2, 3}
		}`

	fset := token.NewFileSet()
	pkg, err := parser.ParseFile(fset, "string", code, 0)
	if err != nil {
		log.Fatalf("parse error: %+v", err)
	}

	fn := pkg.Decls[0].(*ast.FuncDecl)
	x := fn.Body.List[0]
	y := fn.Body.List[1]

	// Reflect DeepEqual will fail due to different Pos values.
	// astequal only checks whether two nodes describe AST.
	fmt.Println(reflect.DeepEqual(x, y)) // => false
	fmt.Println(astequal.Node(x, y))     // => true
	fmt.Println(astequal.Stmt(x, y))     // => true

	// Output:
	// false
	// true
	// true
}
