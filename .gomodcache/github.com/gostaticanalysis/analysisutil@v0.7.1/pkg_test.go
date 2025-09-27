package analysisutil_test

import (
	"fmt"
	"go/types"
	"path/filepath"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/packages/packagestest"
)

func TestRemoveVendor(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"vendor/path/to/pkg", "path/to/pkg"},
		{"path/to/pkg", "path/to/pkg"},
		{"", ""},
		{"a/vendor/path/to/pkg", "path/to/pkg"},
		{"pkg", "pkg"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			got := analysisutil.RemoveVendor(tt.path)
			if got != tt.want {
				t.Errorf("want %s but got %s", tt.want, got)
			}
		})
	}
}

func TestLookupFromImports(t *testing.T) {
	tests := []struct {
		path, name string
		found      bool
	}{
		{"fmt", "Println", true},
		{"b", "Msg", true},
		{"strings", "Join", false},
		{"b", "Func", false},
	}

	testdata := filepath.Join("testdata", "lookupfromimports")
	pkg := loadPkg(t, testdata, "a")
	for _, tt := range tests {
		tt := tt
		name := fmt.Sprintf("%s.%s", tt.path, tt.name)
		t.Run(name, func(t *testing.T) {
			obj := analysisutil.LookupFromImports(pkg.Imports(), tt.path, tt.name)
			switch {
			case obj == nil && tt.found:
				t.Error("not found")
			case obj != nil && !tt.found:
				t.Error("want not found but found:", obj)
			}
		})
	}
}

func loadPkg(t *testing.T, testdata, pkg string) *types.Package {
	t.Helper()

	exported := packagestest.Export(t, packagestest.GOPATH, []packagestest.Module{{
		Name:  pkg,
		Files: packagestest.MustCopyFileTree(testdata),
	}})
	defer exported.Cleanup()

	conf := exported.Config
	conf.Mode = packages.LoadAllSyntax
	pkgs, err := packages.Load(conf, pkg)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("cannot load package", pkg)
	}
	if len(pkgs[0].Errors) != 0 {
		t.Fatal("unexpected error:", pkgs[0].Errors)
	}
	return pkgs[0].Types
}

func TestImported(t *testing.T) {
	tests := []struct {
		path string
		used bool
	}{
		{"fmt", true},
		{"b", true},
		{"a", false},
		{"log", true},
	}

	run := func(pass *analysis.Pass) (interface{}, error) {
		for _, tt := range tests {
			t.Run(tt.path, func(t *testing.T) {
				used := analysisutil.Imported(tt.path, pass)
				if used && !tt.used {
					t.Error("not used")
				} else if !used && tt.used {
					t.Error("used")
				}
			})
		}
		return nil, nil
	}

	var analyzer = &analysis.Analyzer{
		Run:              run,
		RunDespiteErrors: true,
		Requires: []*analysis.Analyzer{
			buildssa.Analyzer,
		},
	}

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "pkgused")
}
