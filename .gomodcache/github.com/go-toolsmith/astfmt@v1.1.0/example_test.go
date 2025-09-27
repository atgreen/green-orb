package astfmt_test

import (
	"go/token"
	"os"

	"github.com/go-toolsmith/astfmt"
	"github.com/go-toolsmith/strparse"
)

func Example() {
	x := strparse.Expr(`foo(bar(baz(1+2)))`)
	astfmt.Println(x)                         // => foo(bar(baz(1 + 2)))
	astfmt.Fprintf(os.Stdout, "node=%s\n", x) // => node=foo(bar(baz(1 + 2)))

	// Can use specific file set with printer.
	fset := token.NewFileSet() // Suppose this fset is used when parsing
	pp := astfmt.NewPrinter(fset)
	pp.Println(x) // => foo(bar(baz(1 + 2)))

	// Output:
	// foo(bar(baz(1 + 2)))
	// node=foo(bar(baz(1 + 2)))
	// foo(bar(baz(1 + 2)))
}
