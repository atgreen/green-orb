package copyloopvar

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		testdata := testutil.WithModules(t, analysistest.TestData(), nil)
		analysistest.Run(t, testdata, NewAnalyzer(), "basic")
	})

	t.Run("check-alias", func(t *testing.T) {
		analyzer := NewAnalyzer()
		if err := analyzer.Flags.Set("check-alias", "true"); err != nil {
			t.Error(err)
		}
		testdata := testutil.WithModules(t, analysistest.TestData(), nil)
		analysistest.Run(t, testdata, analyzer, "checkalias")
	})
}
