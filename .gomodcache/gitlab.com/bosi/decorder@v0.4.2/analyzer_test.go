package decorder

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	tcs := []struct {
		name string
		dir  string
		opts func(opts *options)
	}{
		{
			name: "default conf",
			dir:  "a",
			opts: func(opts *options) {},
		},
		{
			name: "custom dec order",
			dir:  "customDecOrder",
			opts: func(opts *options) { opts.decOrder = "const,var" },
		},
		{
			name: "custom dec order all",
			dir:  "customDecOrderAll",
			opts: func(opts *options) { opts.decOrder = "func ,const,   var ,type" },
		},
		{
			name: "disabled dec order check",
			dir:  "disabledDecOrderCheck",
			opts: func(opts *options) { opts.disableDecOrderCheck = true },
		},
		{
			name: "disabled init func first check",
			dir:  "disabledInitFuncFirstCheck",
			opts: func(opts *options) { opts.disableInitFuncFirstCheck = true },
		},
		{
			name: "disabled dec num check",
			dir:  "disabledDecNumCheck",
			opts: func(opts *options) { opts.disableDecNumCheck = true },
		},
		{
			name: "disabled type dec num check",
			dir:  "disabledTypeDecNumCheck",
			opts: func(opts *options) { opts.disableTypeDecNumCheck = true },
		},
		{
			name: "disabled const dec num check",
			dir:  "disabledConstDecNumCheck",
			opts: func(opts *options) { opts.disableConstDecNumCheck = true },
		},
		{
			name: "disabled var dec num check",
			dir:  "disabledVarDecNumCheck",
			opts: func(opts *options) { opts.disableVarDecNumCheck = true },
		},
		{
			name: "ignore underscore var at dec order check",
			dir:  "ignoreUnderscoreVars",
			opts: func(opts *options) { opts.ignoreUnderscoreVars = true },
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			opts = options{decOrder: defaultDecOrder}
			tc.opts(&opts)
			analysistest.Run(t, testdata(t), Analyzer, tc.dir)
		})
	}
}

func testdata(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(wd, "testdata")
}
