//golangcitest:args -Eireturn
//golangcitest:config_path testdata/ireturn_reject_stdlib.yml
package testdata

import (
	"bytes"
	"io"
)

func NewWriter() io.Writer { // want `NewWriter returns interface \(io.Writer\)`
	var buf bytes.Buffer
	return &buf
}

func TestError() error {
	return nil
}

type Foo interface {
	Foo()
}
type foo int

func (f foo) Foo() {}

func NewFoo() Foo {
	return foo(1)
}
