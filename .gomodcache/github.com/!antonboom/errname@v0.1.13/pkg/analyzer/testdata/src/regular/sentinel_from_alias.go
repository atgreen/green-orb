package regular

import (
	stderrors "errors"
	stdfmt "fmt"
)

var (
	ErrAlias = stderrors.New("err from alias")
	AliasErr = stderrors.New("err from alias") // want "the variable name `AliasErr` should conform to the `ErrXxx` format"
)

var (
	ErrOutOfSize2   = stdfmt.Errorf("out of size (max %d)", maxSize)
	OutOfSizeError2 = stdfmt.Errorf("out of size (max %d)", maxSize) // want "the variable name `OutOfSizeError2` should conform to the `ErrXxx` format"
)
