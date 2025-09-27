package astp_test

import (
	"fmt"

	"github.com/go-toolsmith/astp"
	"github.com/go-toolsmith/strparse"
)

func Example() {
	if astp.IsIdent(strparse.Expr(`x`)) {
		fmt.Println("ident")
	}
	if astp.IsBlockStmt(strparse.Stmt(`{f()}`)) {
		fmt.Println("block stmt")
	}
	if astp.IsGenDecl(strparse.Decl(`var x int = 10`)) {
		fmt.Println("gen decl")
	}

	// Output:
	// ident
	// block stmt
	// gen decl
}
