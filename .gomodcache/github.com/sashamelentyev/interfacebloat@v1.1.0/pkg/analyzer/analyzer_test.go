package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUseStdlibVars(t *testing.T) {
	pkgs := []string{
		"a",
	}

	analyzer := New()

	analysistest.Run(t, analysistest.TestData(), analyzer, pkgs...)
}
