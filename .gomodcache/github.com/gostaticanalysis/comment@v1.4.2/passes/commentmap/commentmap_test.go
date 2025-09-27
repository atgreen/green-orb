package commentmap_test

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/txtar"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"github.com/gostaticanalysis/testutil"
)

func Test_Maps_Ignore(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path  string
		found bool
	}{
		"ignore": {
			path:  "ignore",
			found: true,
		},
		"notignore": {
			path:  "notignore",
			found: false,
		},
		"havecomment": {
			path:  "havecomment",
			found: true,
		},
	}
	for name, tt := range tests {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testdata := parseTestdata(t, tt.path)
			analyzer := &analysis.Analyzer{
				Requires: []*analysis.Analyzer{
					inspect.Analyzer,
					commentmap.Analyzer,
				},
				Run: func(pass *analysis.Pass) (interface{}, error) {
					var found bool

					ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
					cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)
					ins.Preorder(nil, func(n ast.Node) {
						if cmaps.Ignore(n, "check") {
							found = true
						}
					})

					if found != tt.found {
						return nil, fmt.Errorf("%q not found", name)
					}

					return nil, nil
				},
			}

			analysistest.Run(t, testdata, analyzer, tt.path)
		})
	}
}

func parseTestdata(t *testing.T, name string) string {
	t.Helper()

	filemap := make(map[string]string)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if !strings.Contains(path, name) {
			return filepath.SkipDir
		}

		ar, err := txtar.ParseFile(path)
		if err != nil {
			return err
		}
		for _, file := range ar.Files {
			filemap[filepath.Join(name, file.Name)] = string(file.Data)
		}

		return nil
	}
	if err := filepath.Walk(analysistest.TestData(), walkFn); err != nil {
		t.Fatal(err)
	}

	return testutil.WriteFiles(t, filemap)
}
