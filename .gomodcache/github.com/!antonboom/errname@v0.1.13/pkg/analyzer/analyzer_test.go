package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrName(t *testing.T) {
	pkgs := []string{
		"regular",
		"unusual/errortype",
		"unusual/generics",
		"unusual/newfunc",
	}
	analysistest.Run(t, analysistest.TestData(), New(), pkgs...)
}

func Test_getPkgFromPath(t *testing.T) {
	for _, tt := range []struct {
		in, expected string
	}{
		{},
		{in: "errors", expected: "errors"},
		{in: "github/pkg/errors", expected: "errors"},
		{in: "github.com/cockroachdb/errors", expected: "errors"},
	} {
		pkg := getPkgFromPath(tt.in)
		if pkg != tt.expected {
			t.Errorf("%q != %q", pkg, tt.expected)
		}
	}
}
