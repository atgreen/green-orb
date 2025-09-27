package analyzer_test

import (
	"log"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/kulti/thelper/pkg/analyzer"
)

//go:generate go run github.com/kulti/thelper/scripts/generator --name t --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name f --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name b --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name tb --interface --path testdata/src

func TestAllChecks(t *testing.T) {
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	testdata := analysistest.TestData()

	t.Run("without check flag", func(t *testing.T) {
		t.Parallel()

		a := analyzer.NewAnalyzer()
		analysistest.Run(t, testdata, a, "t", "f", "b", "tb")
	})

	t.Run("empty checks flag", func(t *testing.T) {
		t.Parallel()

		a := analyzer.NewAnalyzer()
		err := a.Flags.Set("checks", "")
		if err != nil {
			t.Fatalf("failed to set checks empty value: %v", err)
		}
		analysistest.Run(t, testdata, a, "t", "f", "b", "tb")
	})
}

//go:generate go run github.com/kulti/thelper/scripts/generator --name t --check begin --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name f --check begin --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name b --check begin --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name tb --interface --check begin --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name t --check first --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name f --check first --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name b --check first --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name tb --interface --check first --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name t --check name --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name f --check name --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name b --check name --path testdata/src
//go:generate go run github.com/kulti/thelper/scripts/generator --name tb --interface --check name --path testdata/src

func TestSingleCheck(t *testing.T) {
	t.Parallel()

	checks := []string{"t_begin", "t_first", "t_name", "f_begin", "f_first", "f_name", "b_begin", "b_first", "b_name", "tb_begin", "tb_first", "tb_name"}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	testdata := analysistest.TestData()

	for _, tc := range checks {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()

			a := analyzer.NewAnalyzer()
			err := a.Flags.Set("checks", tc)
			if err != nil {
				t.Fatalf("failed to set checks into %q: %v", tc, err)
			}
			analysistest.Run(t, testdata, a, tc)
		})
	}
}

func TestShadowTesting(t *testing.T) {
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(), "shadow_testing")
}
