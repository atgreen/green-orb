package regular

import (
	"errors"
	"fmt"
)

var (
	EOF          = errors.New("end of file")
	ErrEndOfFile = errors.New("end of file")
	errEndOfFile = errors.New("end of file")

	EndOfFileError = errors.New("end of file") // want "the variable name `EndOfFileError` should conform to the `ErrXxx` format"
	ErrorEndOfFile = errors.New("end of file") // want "the variable name `ErrorEndOfFile` should conform to the `ErrXxx` format"
	EndOfFileErr   = errors.New("end of file") // want "the variable name `EndOfFileErr` should conform to the `ErrXxx` format"
	endOfFileError = errors.New("end of file") // want "the variable name `endOfFileError` should conform to the `errXxx` format"
	errorEndOfFile = errors.New("end of file") // want "the variable name `errorEndOfFile` should conform to the `errXxx` format"
)

const maxSize = 256

var (
	ErrOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)
	errOutOfSize = fmt.Errorf("out of size (max %d)", maxSize)

	OutOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // want "the variable name `OutOfSizeError` should conform to the `ErrXxx` format"
	outOfSizeError = fmt.Errorf("out of size (max %d)", maxSize) // want "the variable name `outOfSizeError` should conform to the `errXxx` format"
)

func errInsideFuncIsNotSentinel() error {
	var lastErr error
	return lastErr
}
