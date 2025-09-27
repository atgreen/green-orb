package analysisutil_test

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"path/filepath"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestObjectOf(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		src   string
		pkg   string // blank means same as the map key
		name  string
		found bool
	}{
		"standard":    {`import _ "fmt"`, "fmt", "Println", true},
		"unimport":    {"", "fmt", "Println", false},
		"notexiststd": {`import _ "fmt"`, "fmt", "NOTEXIST", false},
		"typename":    {"type A int", "", "A", true},
		"unexportvar": {"var n int", "", "n", true},
		"exportvar":   {"var N int", "", "N", true},
		"notexist":    {"", "", "NOTEXIST", false},
		"vendored":    {`import _ "fmt"`, "vendor/fmt", "Println", true},
		"pointer":     {"type A int", "", "*A", false},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &analysis.Analyzer{
				Name: name + "Analyzer",
				Run: func(pass *analysis.Pass) (interface{}, error) {
					pkg := name
					if tt.pkg != "" {
						pkg = tt.pkg
					}
					obj := analysisutil.ObjectOf(pass, pkg, tt.name)
					switch {
					case tt.found && obj == nil:
						return nil, errors.New("expect found but not found")
					case !tt.found && obj != nil:
						return nil, fmt.Errorf("unexpected return value: %v", obj)
					}
					return nil, nil
				},
			}
			path := filepath.Join(name, name+".go")
			dir := WriteFiles(t, map[string]string{
				path: fmt.Sprintf("package %s\n%s", name, tt.src),
			})
			analysistest.Run(t, dir, a, name)
		})
	}

}

func TestUnder(t *testing.T) {
	t.Parallel()

	lookup := func(pass *analysis.Pass, n string) (types.Type, error) {
		_, obj := pass.Pkg.Scope().LookupParent(n, token.NoPos)
		if obj == nil {
			return nil, fmt.Errorf("does not find: %s", n)
		}
		return obj.Type(), nil
	}

	cases := map[string]struct {
		src  string
		typ  string
		want string
	}{
		"nonamed":  {"", "int", "int"},
		"named":    {"type A int", "A", "int"},
		"twonamed": {"type A int; type B A", "B", "int"},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &analysis.Analyzer{
				Name: name + "Analyzer",
				Run: func(pass *analysis.Pass) (interface{}, error) {
					typ, err := lookup(pass, tt.typ)
					if err != nil {
						return nil, err
					}
					want, err := lookup(pass, tt.want)
					if err != nil {
						return nil, err
					}
					got := analysisutil.Under(typ)
					if !types.Identical(want, got) {
						return nil, fmt.Errorf("want %v but got %v", want, got)
					}
					return nil, nil
				},
			}
			path := filepath.Join(name, name+".go")
			dir := WriteFiles(t, map[string]string{
				path: fmt.Sprintf("package %s\n%s", name, tt.src),
			})
			analysistest.Run(t, dir, a, name)
		})
	}
}

func TestField(t *testing.T) {
	t.Parallel()

	lookup := func(pass *analysis.Pass, n string) (types.Type, error) {
		_, obj := pass.Pkg.Scope().LookupParent(n, token.NoPos)
		if obj == nil {
			return nil, fmt.Errorf("does not find: %s", n)
		}
		return obj.Type(), nil
	}

	cases := map[string]struct {
		src   string
		typ   string
		field string
		want  int
	}{
		"nomarl":     {"type a struct{n int}", "a", "n", 0},
		"nofield":    {"type a struct{n int}", "a", "m", -1},
		"empty":      {"type a struct{}", "a", "n", -1},
		"two":        {"type a struct{n, m int}", "a", "m", 1},
		"nonamed":    {"var a struct{n, m int}", "a", "m", 1},
		"ptr":        {"var a *struct{n, m int}", "a", "m", 1},
		"namednamed": {"type a struct{n int}; type b a", "b", "n", 0},
		"alias":      {"type a struct{n int}; type b = a", "b", "n", 0},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &analysis.Analyzer{
				Name: name + "Analyzer",
				Run: func(pass *analysis.Pass) (interface{}, error) {
					typ, err := lookup(pass, tt.typ)
					if err != nil {
						return nil, err
					}

					got, _ := analysisutil.Field(typ, tt.field)
					if tt.want != got {
						return nil, fmt.Errorf("want %v but got %v", tt.want, got)
					}
					return nil, nil
				},
			}
			path := filepath.Join(name, name+".go")
			dir := WriteFiles(t, map[string]string{
				path: fmt.Sprintf("package %s\n%s", name, tt.src),
			})
			analysistest.Run(t, dir, a, name)
		})
	}
}
