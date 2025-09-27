package analysisutil_test

import (
	"reflect"
	"testing"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

func Test_IfInstr(t *testing.T) {
	type I = ssa.Instruction
	var ifinst ssa.If
	tests := []struct {
		name   string
		instrs []ssa.Instruction
		hasIf  bool
	}{
		{"empty", []I{}, false},
		{"hasIf", []I{&ifinst}, true},
		{"hasIf and others", []I{nil, nil, &ifinst}, true},
		{"not tail", []I{&ifinst, nil}, false}, // actually not occur
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b ssa.BasicBlock
			b.Instrs = tt.instrs
			got := analysisutil.IfInstr(&b)
			switch {
			case got == nil && tt.hasIf:
				t.Errorf("want *ssa.If but got nil")
			case got != nil && !tt.hasIf:
				t.Errorf("want nil but got *ssa.If")
			}
		})
	}
}

func Test_Phi(t *testing.T) {
	type I = ssa.Instruction
	var phi1, phi2 ssa.Phi
	tests := []struct {
		name   string
		instrs []ssa.Instruction
		want   []*ssa.Phi
	}{
		{"empty", []I{}, nil},
		{"P", []I{&phi1}, []*ssa.Phi{&phi1}},
		{"P_", []I{&phi1, nil}, []*ssa.Phi{&phi1}},
		{"_P", []I{nil, &phi1}, nil}, // actually not occur
		{"PP", []I{&phi1, &phi2}, []*ssa.Phi{&phi1, &phi2}},
		{"P_P", []I{&phi1, nil, &phi2}, []*ssa.Phi{&phi1}}, // actually not occur
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b ssa.BasicBlock
			b.Instrs = tt.instrs
			got := analysisutil.Phi(&b)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %#v but got %#v", tt.want, got)
			}
		})
	}
}

func Test_BinOp(t *testing.T) {
	type I = ssa.Instruction
	var bo1, bo2 ssa.BinOp
	tests := []struct {
		name   string
		instrs []ssa.Instruction
		want   []*ssa.BinOp
	}{
		{"empty", []I{}, nil},
		{"B", []I{&bo1}, []*ssa.BinOp{&bo1}},
		{"B_", []I{&bo1, nil}, []*ssa.BinOp{&bo1}},
		{"_B", []I{nil, &bo1}, []*ssa.BinOp{&bo1}},
		{"BB", []I{&bo1, &bo2}, []*ssa.BinOp{&bo1, &bo2}},
		{"B_If", []I{&bo1, nil, &ssa.If{
			Cond: &bo1,
		}}, []*ssa.BinOp{&bo1}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var b ssa.BasicBlock
			b.Instrs = tt.instrs
			got := analysisutil.BinOp(&b)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %#v but got %#v", tt.want, got)
			}
		})
	}
}

func TestUsed(t *testing.T) {
	a := &analysis.Analyzer{
		Name:     "used",
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			srcFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs
			for _, f := range srcFuncs {
				if len(f.Params) != 1 {
					continue
				}
				v := f.Params[0]
				for _, b := range f.Blocks {
					if analysisutil.Used(v, b.Instrs) != nil {
						pass.Reportf(v.Pos(), "used")
					}
				}
			}
			return nil, nil
		},
	}
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, a, "used")
}
