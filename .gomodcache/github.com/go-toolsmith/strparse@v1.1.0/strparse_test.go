package strparse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"testing"
)

func equalPrintedForm(x ast.Node, y string) bool {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), x); err != nil {
		panic(fmt.Sprintf("format error: %v", err))
	}
	b, err := format.Source([]byte(y))
	if err != nil {
		panic(fmt.Sprintf("format error: %v", err))
	}
	xString := buf.String()
	yString := string(b)
	return xString == yString
}

func TestParseExpr(t *testing.T) {
	const input = `1 * (2 << x)/s[1].f()(^y)`
	expr := Expr(input)
	if !equalPrintedForm(expr, input) {
		t.Error("expr printed form not match")
	}
}

func TestParseStmt(t *testing.T) {
	const input = `for i := 0; i < len(xs); i++ { fmt.Println(i) }`
	expr := Stmt(input)
	if !equalPrintedForm(expr, input) {
		t.Error("stmt printed form not match")
	}
}

func TestParseDecl(t *testing.T) {
	const input = `var (x int = 19; a, b = f(); c, d = "1", "2"+"3")`
	expr := Decl(input)
	if !equalPrintedForm(expr, input) {
		t.Error("decl printed form not match")
	}
}
