package analysisutil_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestReportWithoutIgnore(t *testing.T) {
	testdata := analysistest.TestData()

	cases := []struct {
		pkg   string
		names []string
	}{
		{"withoutnames", nil},
		{"withnames", []string{"check1", "check2"}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.pkg, func(t *testing.T) {
			analyzer := &analysis.Analyzer{
				Name: "test",
				Run: func(pass *analysis.Pass) (interface{}, error) {
					pass.Report = analysisutil.ReportWithoutIgnore(pass, tt.names...)
					for _, f := range pass.Files {
						for _, decl := range f.Decls {
							decl, ok := decl.(*ast.GenDecl)
							if !ok || decl.Tok != token.VAR {
								continue
							}
							pass.Reportf(decl.Pos(), "NG")
							pass.ReportRangef(decl, "NG")
						}
					}
					return nil, nil
				},
			}
			analysistest.Run(t, testdata, analyzer, "diagnostic/"+tt.pkg)
		})
	}
}
