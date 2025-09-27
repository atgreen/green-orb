package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		dir     string
		pattern string
	}{
		{
			dir: "defaults",
		},
		{
			dir:     "custompattern",
			pattern: `.*`,
		},
		{
			dir:     "defaultclient",
			pattern: `^(DefaultClient|DefaultTransport)$`,
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	for _, tc := range tests {
		tt := tc

		t.Run(tt.dir, func(t *testing.T) {
			td := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata", tt.dir)
			a := New()
			if tt.pattern != "" {
				_ = a.Flags.Set("pattern", tt.pattern)
			}
			analysistest.Run(t, td, a, "./...")
		})
	}
}
