package analysisutil_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func WriteFiles(t *testing.T, filemap map[string]string) string {
	t.Helper()
	dir, clean, err := analysistest.WriteFiles(filemap)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	t.Cleanup(clean)
	return dir
}
