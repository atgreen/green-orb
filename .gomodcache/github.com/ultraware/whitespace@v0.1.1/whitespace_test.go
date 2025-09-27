package whitespace

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWantMultiline(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(&Settings{
		Mode:      RunningModeNative,
		MultiIf:   true,
		MultiFunc: true,
	})

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "whitespace_multiline")
}

func TestNoMultiline(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := NewAnalyzer(&Settings{
		Mode:      RunningModeNative,
		MultiIf:   false,
		MultiFunc: false,
	})

	analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "whitespace_no_multiline")
}
