//go:build go1.18
// +build go1.18

package varnamelen

import "testing"

func TestVarNameLen_Run_Warnings_TypeParam(t *testing.T) {
	run(t, "warningstypeparam", map[string]string{
		"checkTypeParam": "true",
	})
}

func TestVarNameLen_Run_Distance_TypeParam(t *testing.T) {
	run(t, "distancetypeparam", map[string]string{
		"checkTypeParam": "true",
	})
}

func TestVarNameLen_Run_IgnoreNames_TypeParam(t *testing.T) {
	run(t, "ignorenamestypeparam", map[string]string{
		"ignoreNames":    "T",
		"checkTypeParam": "true",
	})
}

func TestVarNameLen_Run_IgnoreDecls_TypeParam(t *testing.T) {
	run(t, "ignoredeclstypeparam", map[string]string{
		"ignoreDecls":    "T any",
		"checkTypeParam": "true",
	})
}
