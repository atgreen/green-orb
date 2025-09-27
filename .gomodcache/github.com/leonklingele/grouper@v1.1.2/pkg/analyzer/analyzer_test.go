package analyzer_test

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/leonklingele/grouper/pkg/analyzer"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TODO(leon): Add fuzzing

func TestConst(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-const",
			flags: flags().
				withConstRequireSingleConst().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withConstRequireGrouping().
				build(),
		},

		{
			name: "mixed-named-with-vars",
			flags: flags().
				withConstRequireSingleConst().
				withConstRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "const")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

func TestImport(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withImportRequireSingleImport().
				withImportRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withImportRequireSingleImport().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-import",
			flags: flags().
				withImportRequireSingleImport().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},

		{
			name: "mixed-named-require-single-import",
			flags: flags().
				withImportRequireSingleImport().
				build(),
		},
		{
			name: "mixed-named-require-grouping",
			flags: flags().
				withImportRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "import")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

func TestType(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withTypeRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withTypeRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withTypeRequireSingleType().
				withTypeRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withTypeRequireSingleType().
				withTypeRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-type",
			flags: flags().
				withTypeRequireSingleType().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withTypeRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "type")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

func TestVar(t *testing.T) {
	t.Parallel()

	fixtures := []struct {
		name  string
		flags flag.FlagSet
	}{
		{
			name: "single-grouped",
			flags: flags().
				withVarRequireGrouping().
				build(),
		},
		{
			name: "single-ungrouped",
			flags: flags().
				withVarRequireGrouping().
				build(),
		},

		{
			name: "multi-grouped",
			flags: flags().
				withVarRequireSingleVar().
				withVarRequireGrouping().
				build(),
		},
		{
			name: "multi-ungrouped",
			flags: flags().
				withVarRequireSingleVar().
				withVarRequireGrouping().
				build(),
		},

		{
			name: "mixed-require-single-var",
			flags: flags().
				withVarRequireSingleVar().
				build(),
		},
		{
			name: "mixed-require-grouping",
			flags: flags().
				withVarRequireGrouping().
				build(),
		},

		{
			name: "mixed-named-with-consts",
			flags: flags().
				withVarRequireSingleVar().
				withVarRequireGrouping().
				build(),
		},
		{
			name: "mixed-named-with-var-shorthand",
			flags: flags().
				withVarRequireSingleVar().
				withVarRequireGrouping().
				build(),
		},
	}

	for _, f := range fixtures {
		f := f

		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			a := analyzer.New()
			a.Flags = f.flags

			testdata := filepath.Join(analysistest.TestData(), "var")
			_ = analysistest.Run(t, testdata, a, f.name)
		})
	}
}

type flagger struct {
	fs *flag.FlagSet
}

func (f *flagger) withConstRequireSingleConst() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameConstRequireSingleConst).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withConstRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameConstRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withImportRequireSingleImport() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameImportRequireSingleImport).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withImportRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameImportRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withTypeRequireSingleType() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameTypeRequireSingleType).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withTypeRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameTypeRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withVarRequireSingleVar() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameVarRequireSingleVar).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) withVarRequireGrouping() *flagger {
	if err := f.fs.Lookup(analyzer.FlagNameVarRequireGrouping).Value.Set("true"); err != nil {
		panic(err)
	}

	return f
}

func (f *flagger) build() flag.FlagSet {
	return *f.fs
}

func flags() *flagger {
	fs := analyzer.Flags()

	return &flagger{
		fs: &fs,
	}
}
