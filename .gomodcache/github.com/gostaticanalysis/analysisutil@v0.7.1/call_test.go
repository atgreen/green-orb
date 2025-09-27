package analysisutil_test

import (
	"errors"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var callAnalyzer = &analysis.Analyzer{
	Name: "test_call",
	Run:  callAnalyzerRun,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func TestCall(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, callAnalyzer, "call")
}

func callAnalyzerRun(pass *analysis.Pass) (interface{}, error) {
	resTyp := analysisutil.TypeOf(pass, "call", "*res")
	if resTyp == nil {
		return nil, errors.New("analyzer does not find *call.res type")
	}

	close := analysisutil.MethodOf(resTyp, "close")
	if close == nil {
		return nil, errors.New("analyzer does not find (call.res).close method")
	}

	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		instrs := analysisutil.NotCalledIn(f, resTyp, close)
		for _, instr := range instrs {
			pass.Reportf(instr.Pos(), "NG")
		}
	}

	return nil, nil
}
