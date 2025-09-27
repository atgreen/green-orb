package f

import "testing"

func FuzzNothing(f *testing.F) {
	f.Add("nothing")
	f.Fuzz(func(t *testing.T, s string) {
	})
}
