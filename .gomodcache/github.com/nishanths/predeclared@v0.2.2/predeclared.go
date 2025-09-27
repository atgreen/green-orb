// Command predeclared prints the names and locations of declarations in the
// given packages that have the same name as one of Go's predeclared
// identifiers (eg., int, string, delete, copy, append).
//
// Usage
//
// The command line usage is:
//
//  predeclared [flags] [packages...]
//
// Run 'predeclared' without arguments for help.
//
// Flags
//
// The '-q' boolean flag, if set, indicates to the command to check struct
// field names, interface methods, and method names — in addition to the
// default checks. (These checks aren't included by default since fields and
// method are always accessed by a qualifier—à la obj.Field—and hence are less
// likely to cause confusion when reading code even if they have the same name
// as a predeclared identifier.)
//
// The '-ignore' string flag can be used to specify predeclared identifiers to
// not report issues for. For example, to not report overriding of the
// predeclared identifiers 'new' and 'real', set the flag like so:
//
//  -ignore=new,real
//
package main

import (
	"github.com/nishanths/predeclared/passes/predeclared"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(predeclared.Analyzer)
}
