package analysisutil_test

import (
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

var ssainspectAnalyzer = &analysis.Analyzer{
	Name: "test_ssainspect",
	Run:  ssainspectAnalyzerRun,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func TestSSAInspect(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ssainspectAnalyzer, "ssainspect")
}

func ssainspectAnalyzerRun(pass *analysis.Pass) (interface{}, error) {
	funcs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
	for _, f := range funcs {
		m := map[ssa.Instruction]bool{}
		if len(f.Blocks) == 0 {
			continue
		}
		analysisutil.InspectInstr(f.Blocks[0], 0, func(i int, instr ssa.Instruction) bool {
			m[instr] = true
			return true
		})

		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				if !m[instr] {
					pass.Reportf(instr.Pos(), "NG")
				}
			}
		}
	}
	return nil, nil
}
