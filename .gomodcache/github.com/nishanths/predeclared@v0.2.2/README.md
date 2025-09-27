# predeclared [![Build Status](https://travis-ci.org/nishanths/predeclared.svg?branch=master)](https://travis-ci.org/nishanths/predeclared) [![Godoc](https://godoc.org/github.com/nishanths/predeclared?status.svg)](http://godoc.org/github.com/nishanths/predeclared)


Find code that overrides one of Go's predeclared identifiers (`new`, `make`, `append` `uint`, etc.).

The list of predeclared identifiers can be found in the [spec](https://golang.org/ref/spec#Predeclared_identifiers).

```
go get github.com/nishanths/predeclared
```

See [godoc](https://godoc.org/github.com/nishanths/predeclared) or run `predeclared` without arguments to print usage.

## Test

```
cd passes/predeclared && go test
```

## Examples

Given a package with the file:

```go
package pkg // import "example.org/foo/pkg"

func copy()  {}
func print() {}

func foo() string {
	string := "x"
	return string
}

type int struct{}
```

running:

```
predeclared example.org/foo/pkg
```

prints:

```
example.go:3:6: function "copy" has same name as predeclared identifier
example.go:4:6: function "print" has same name as predeclared identifier
example.go:7:2: variable "string" has same name as predeclared identifier
example.go:11:6: type "int" has same name as predeclared identifier
```

Running the program on the standard library's `text` package and subpackages produces:

```sh
$ predeclared text/...
/usr/local/go/src/text/template/parse/node.go:554:33: param "true" has same name as predeclared identifier
/usr/local/go/src/text/template/parse/node.go:599:3: variable "rune" has same name as predeclared identifier
/usr/local/go/src/text/template/funcs.go:164:36: param "cap" has same name as predeclared identifier
/usr/local/go/src/text/template/funcs.go:233:3: variable "cap" has same name as predeclared identifier
/usr/local/go/src/text/template/template.go:218:30: param "new" has same name as predeclared identifier
/usr/local/go/src/text/template/exec_test.go:232:21: param "error" has same name as predeclared identifier
```
