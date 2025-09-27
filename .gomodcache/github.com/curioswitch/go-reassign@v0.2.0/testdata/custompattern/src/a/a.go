package a

import (
	"b"
)

var st = struct {
	ErrSt error
}{}

func foo() {
	b.ErrB = nil // want "reassigning variable"

	b.NotErr = "is error" // want "reassigning variable"
}
