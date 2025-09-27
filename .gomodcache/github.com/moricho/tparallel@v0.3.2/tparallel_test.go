package tparallel_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/moricho/tparallel"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, tparallel.Analyzer, "test")
}
