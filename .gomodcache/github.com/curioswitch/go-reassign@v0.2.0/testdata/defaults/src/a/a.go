package a

import (
	"b"

	"errors"
	"io"
)
import cc "c"

import (
	"d"
	"d/e"
	"d/e/f"
)

var st = struct {
	ErrSt error
}{}

func foo() {
	b.ErrB = nil // want "reassigning variable"

	cc.ErrC = nil // want "reassigning variable"

	d.ErrD = nil // want "reassigning variable"

	e.ErrE = nil // want "reassigning variable"

	f.ErrF = nil // want "reassigning variable"

	io.EOF = nil // want "reassigning variable"

	st.ErrSt = errors.New("foo")

	b.NotErr = "is error"
}
