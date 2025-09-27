package testableexamples_test

import (
	"path/filepath"
	"testing"

	"github.com/maratori/testableexamples/pkg/testableexamples"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// dir is named `testfiles` not `testdata` to be able to run examples with `make test`
	testdata, err := filepath.Abs("testfiles")
	if err != nil {
		t.FailNow()
	}

	analysistest.Run(t, testdata, testableexamples.NewAnalyzer())
}
