package varnamelen

import "testing"

func TestVarNameLen_Run_FalsePositive_Sixteen(t *testing.T) {
	run(t, "falsepositivesixteen", nil)
}
