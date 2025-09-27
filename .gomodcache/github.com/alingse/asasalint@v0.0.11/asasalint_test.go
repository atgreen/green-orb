package asasalint

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

const filename = "<src>"

func makePkg(src string) (*types.Package, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	// use the package name as package path
	conf := types.Config{Importer: importer.Default()}
	return conf.Check(file.Name.Name, fset, []*ast.File{file}, nil)
}

func TestIsSliceAnyType(t *testing.T) {
	p, _ := makePkg("package p; var A []any\nvar B any = []any{1, 2, 3}\nvar C []int")
	aT := p.Scope().Lookup("A").Type()
	if !isSliceAnyType(aT) {
		t.Errorf("isSliceAnyType(%v) = false, want true", aT)
	}
	bT := p.Scope().Lookup("B").Type()
	if isSliceAnyType(bT) {
		t.Errorf("isSliceAnyType(%v) = true, want false", bT)
	}
	cT := p.Scope().Lookup("C").Type()
	if isSliceAnyType(cT) {
		t.Errorf("isSliceAnyType(%v) = true, want false", cT)
	}
}

func TestIsSliceAnyVariadicFuncType(t *testing.T) {
	p, _ := makePkg("package p; func hello(a int, b ...any) {}\nfunc hello2(a int, b int) {}")
	ft := p.Scope().Lookup("hello").Type()
	if !isSliceAnyVariadicFuncType(ft) {
		t.Errorf("isSliceAnyVariadicFuncType(%v) = false, want true", ft)
	}

	ft2 := p.Scope().Lookup("hello2").Type()
	if isSliceAnyVariadicFuncType(ft2) {
		t.Errorf("isSliceAnyVariadicFuncType(%v) = true, want false", ft2)
	}
}

func TestGetFuncName(t *testing.T) {
	testCases := []struct {
		desc     string
		src      string
		expected string
	}{
		{
			desc: "function",
			src: `package p
func hello(a int, b ...any) {}
func hello2(a int, b int) { hello(a, b)}
`,
			expected: "hello",
		},
		{
			desc: "method",
			src: `package p
type A struct {}
func (a *A) hello(a int, b ...any) {}
func (a *A) hello2(a int, b int) {
	a.hello(a, b)
}
`,
			expected: "a.hello",
		},
		{
			desc: "function inside a method",
			src: `package p
type A struct {}
func (a *A) hello(a int, b ...any) {
	hello(a, b...)
}
`,
			expected: "hello",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "src.go", test.src, 0)
			if err != nil {
				panic(err)
			}

			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.CallExpr:
					s, err := getFuncName(fset, x)
					if err != nil {
						t.Fatal(err)
					}

					if s != test.expected {
						t.Errorf("%s: got %s, want %s", fset.Position(x.Fun.Pos()), s, test.expected)
					}
				}
				return true
			})
		})
	}
}

func TestAnalyzer(t *testing.T) {
	testCases := []struct {
		desc     string
		settings LinterSetting
	}{
		{
			desc:     "basic",
			settings: LinterSetting{},
		},
		{
			desc: "nobuiltin",
			settings: LinterSetting{
				NoBuiltinExclusions: true,
			},
		},
		{
			desc: "ignoretest",
			settings: LinterSetting{
				IgnoreTest: true,
			},
		},
		{
			desc: "custom",
			settings: LinterSetting{
				Exclude: []string{"get.+"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			a, err := NewAnalyzer(test.settings)
			if err != nil {
				t.Fatal(err)
			}

			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), a, test.desc)
		})
	}
}
