package nlreturn

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	analysistest.Run(t,
		filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata"),
		NewAnalyzer(),
		"p")
}

func TestBlockSize(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	analyzer := NewAnalyzer()
	err = analyzer.Flags.Set("block-size", "2")
	if err != nil {
		t.Fatalf("Failed to set flag: %s", err)
	}

	analysistest.Run(t,
		filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata"),
		analyzer,
		"bs")
}
