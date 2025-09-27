package modinfo

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// NOTE: analysistest does not yet support modules;
	// see https://github.com/golang/go/issues/37054 for details.
	// The workspaces are also not really supported, we can't run the analyzer at the root of the workspace.

	testData := analysistest.TestData()

	testCases := []struct {
		desc     string
		dir      string
		patterns []string
		len      int
		expected []ModInfo
	}{
		{
			desc:     "simple",
			dir:      "a",
			patterns: []string{"a"},
			len:      1,
			expected: []ModInfo{{
				Path:      "github.com/golangci/modinfo/testdata/b",
				Dir:       filepath.Join(testData, "src", "a"),
				GoMod:     filepath.Join(testData, "src", "a", "go.mod"),
				GoVersion: "1.16",
				Main:      true,
			}},
		},
		{
			desc:     "module inside a workspace",
			dir:      "workspace",
			patterns: []string{"workspace/hello/..."},
			len:      2,
			expected: []ModInfo{
				{
					Path:      "example.com/world",
					Dir:       filepath.Join(testData, "src", "workspace", "world"),
					GoMod:     filepath.Join(testData, "src", "workspace", "world", "go.mod"),
					GoVersion: "1.20",
					Main:      true,
				},
				{
					Path:      "hello",
					Dir:       filepath.Join(testData, "src", "workspace", "hello"),
					GoMod:     filepath.Join(testData, "src", "workspace", "hello", "go.mod"),
					GoVersion: "1.20",
					Main:      true,
				},
			},
		},
		{
			desc:     "modules inside a workspace",
			dir:      "workspace",
			patterns: []string{"workspace/hello/...", "workspace/world/..."},
			len:      2,
			expected: []ModInfo{
				{
					Path:      "example.com/world",
					Dir:       filepath.Join(testData, "src", "workspace", "world"),
					GoMod:     filepath.Join(testData, "src", "workspace", "world", "go.mod"),
					GoVersion: "1.20",
					Main:      true,
				},
				{
					Path:      "hello",
					Dir:       filepath.Join(testData, "src", "workspace", "hello"),
					GoMod:     filepath.Join(testData, "src", "workspace", "hello", "go.mod"),
					GoVersion: "1.20",
					Main:      true,
				},
			},
		},
		{
			desc:     "bad module design",
			dir:      "badmodule",
			patterns: []string{"badmodule"},
			len:      1,
			expected: []ModInfo{{
				Path:      "example.com/hello",
				Dir:       filepath.Join(testData, "src", "badmodule"),
				GoMod:     filepath.Join(testData, "src", "badmodule", "go.mod"),
				GoVersion: "1.20",
				Main:      true,
			}},
		},
	}

	t.Setenv("MODINFO_DEBUG_DISABLE_ONCE", "true")

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			results := analysistest.Run(t, testData, Analyzer, test.patterns...)
			for _, result := range results {
				infos, ok := result.Result.([]ModInfo)
				require.True(t, ok)
				assert.Equal(t, test.expected, infos)
			}
		})
	}
}
