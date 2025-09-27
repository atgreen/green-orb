package test

import (
	"testing"
)

func TestWithoutSub(t *testing.T) { // OK
	t.Parallel()
	call("TestWithoutSub")
}
