//go:build go1.18
// +build go1.18

package varnamelen

import "testing"

func TestVarNameLen_Run_FalsePositive(t *testing.T) {
	run(t, "falsepositive", nil)
}
