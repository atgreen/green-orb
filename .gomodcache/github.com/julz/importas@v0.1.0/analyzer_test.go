package importas

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	testCases := []struct {
		desc                 string
		pkg                  string
		aliases              stringMap
		disallowUnaliased    bool
		disallowExtraAliases bool
	}{
		{
			desc: "Invalid imports",
			pkg:  "a",
			aliases: stringMap{
				"fmt": "fff",
				"os":  "stdos",
				"io":  "iio",
			},
		},
		{
			desc: "Valid imports",
			pkg:  "b",
			aliases: stringMap{
				"fmt": "fff",
				"os":  "stdos",
			},
		},
		{
			desc: "external libs",
			pkg:  "c",
			aliases: stringMap{
				"knative.dev/serving/pkg/apis/autoscaling/v1alpha1": "autoscalingv1alpha1",
				"knative.dev/serving/pkg/apis/serving/v1":           "servingv1",
			},
		},
		{
			desc: "regexp",
			pkg:  "d",
			aliases: stringMap{
				"knative.dev/serving/pkg/apis/(\\w+)/(v[\\w\\d]+)": "$1$2",
			},
		},
		{
			desc: "disallow unaliased mode",
			pkg:  "e",
			aliases: stringMap{
				"fmt": "fff",
				"os":  "stdos",
				"io":  "iio",
			},
			disallowUnaliased: true,
		},
		{
			desc:                 "disallow extra alias mode",
			pkg:                  "f",
			disallowExtraAliases: true,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join(testdata, "src", test.pkg)

			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				cmd := exec.Command("go", "mod", "vendor")
				cmd.Dir = dir

				t.Cleanup(func() {
					_ = os.RemoveAll(filepath.Join(testdata, "src", test.pkg, "vendor"))
				})

				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatal(err, string(output))
				}
			}

			cnf := Config{
				RequiredAlias: make(map[string]string),
			}
			flgs := flags(&cnf)
			a := &analysis.Analyzer{
				Flags: flgs,
				Run: func(pass *analysis.Pass) (interface{}, error) {
					return runWithConfig(&cnf, pass)
				},
				Requires: []*analysis.Analyzer{inspect.Analyzer},
			}

			flg := a.Flags.Lookup("alias")
			for k, v := range test.aliases {
				err := flg.Value.Set(fmt.Sprintf("%s:%s", k, v))
				if err != nil {
					t.Fatal(err)
				}
			}

			noUnaliasedFlg := a.Flags.Lookup("no-unaliased")
			if err := noUnaliasedFlg.Value.Set(strconv.FormatBool(test.disallowUnaliased)); err != nil {
				t.Fatal(err)
			}

			noExtraAlisesFlg := a.Flags.Lookup("no-extra-aliases")
			if err := noExtraAlisesFlg.Value.Set(strconv.FormatBool(test.disallowExtraAliases)); err != nil {
				t.Fatal(err)
			}

			analysistest.RunWithSuggestedFixes(t, testdata, a, test.pkg)
		})
	}
}
