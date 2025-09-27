package t

import "testing"

func subtestBuilderAnotherFile() func(*testing.T) {
	return func(t *testing.T) {}
}
