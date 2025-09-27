package loggercheck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/timonwong/loggercheck"
	"github.com/timonwong/loggercheck/internal/rules"
)

type dummyTestingErrorf struct {
	*testing.T
}

func (t dummyTestingErrorf) Errorf(format string, args ...interface{}) {}

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()

	testCases := []struct {
		name      string
		patterns  string
		flags     []string
		wantError string
	}{
		{
			name:     "all",
			patterns: "a/all",
			flags:    []string{"-disable="},
		},
		{
			name:     "require-string-key",
			patterns: "a/requirestringkey",
			flags:    []string{"-requirestringkey"},
		},
		{
			name:     "no-printf-like",
			patterns: "a/noprintflike",
			flags:    []string{"-noprintflike"},
		},
		{
			name:     "klogonly",
			patterns: "a/klogonly",
			flags:    []string{"-disable=logr,zap"},
		},
		{
			name:     "custom-only",
			patterns: "a/customonly",
			flags: []string{
				"-disable=logr",
				"-rulefile",
				"testdata/custom-rules.txt",
			},
		},
		{
			name:     "custom-generic",
			patterns: "a/custom-generic",
			flags: []string{
				"-rulefile",
				"testdata/custom-rules-generic.txt",
			},
		},
		{
			name:     "wrong-rules",
			patterns: "a/customonly",
			flags: []string{
				"-rulefile",
				"testdata/wrong-rules.txt",
			},
			wantError: rules.ErrInvalidRule.Error(),
		},
		{
			name:     "not-found-rules",
			patterns: "a/customonly",
			flags: []string{
				"-rulefile",
				"testdata/xxxxx-wrong-rules-xxxxx.txt",
			},
			wantError: "failed to open rule file",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := loggercheck.NewAnalyzer()
			err := a.Flags.Parse(tc.flags)
			require.NoError(t, err)

			var result []*analysistest.Result
			if tc.wantError != "" {
				result = analysistest.Run(&dummyTestingErrorf{t}, testdata, a, tc.patterns)
			} else {
				result = analysistest.Run(t, testdata, a, tc.patterns)
			}
			require.Len(t, result, 1)

			if tc.wantError != "" {
				assert.Error(t, result[0].Err)
				assert.ErrorContains(t, result[0].Err, tc.wantError)
			}
		})
	}
}

func TestOptions(t *testing.T) {
	testdata := analysistest.TestData()

	customRules := []string{
		"(*a/customonly.Logger).Debugw",
		"(*a/customonly.Logger).Infow",
		"(*a/customonly.Logger).Warnw",
		"(*a/customonly.Logger).Errorw",
		"(*a/customonly.Logger).With",

		"(a/customonly.Logger).XXXDebugw",

		"a/customonly.Debugw",
		"a/customonly.Infow",
		"a/customonly.Warnw",
		"a/customonly.Errorw",
		"a/customonly.With",
	}

	wrongCustomRules := []string{
		"# Wrong rule file",
		"(*a/wrong.Method.Rule",
	}

	testCases := []struct {
		name      string
		options   []loggercheck.Option
		patterns  string
		wantError string
	}{
		{
			name: "wrong-rules",
			options: []loggercheck.Option{
				loggercheck.WithRules(wrongCustomRules),
			},
			patterns:  "a/customonly",
			wantError: "failed to parse rules: ",
		},
		{
			name: "disable-all-then-enable-mylogger",
			options: []loggercheck.Option{
				loggercheck.WithDisable([]string{"klog", "logr", "zap"}),
				loggercheck.WithRules(customRules),
			},
			patterns: "a/customonly",
		},
		{
			name: "ignore-logr",
			options: []loggercheck.Option{
				loggercheck.WithDisable([]string{"logr"}),
				loggercheck.WithRules(customRules),
			},
			patterns: "a/customonly",
		},
		{
			name: "require-string-key",
			options: []loggercheck.Option{
				loggercheck.WithRequireStringKey(true),
			},
			patterns: "a/requirestringkey",
		},
		{
			name: "no-printf-like",
			options: []loggercheck.Option{
				loggercheck.WithNoPrintfLike(true),
			},
			patterns: "a/noprintflike",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := loggercheck.NewAnalyzer(tc.options...)

			var result []*analysistest.Result
			if tc.wantError != "" {
				result = analysistest.Run(&dummyTestingErrorf{t}, testdata, a, tc.patterns)
			} else {
				result = analysistest.Run(t, testdata, a, tc.patterns)
			}
			require.Len(t, result, 1)

			if tc.wantError != "" {
				assert.Error(t, result[0].Err)
				assert.ErrorContains(t, result[0].Err, tc.wantError)
			}
		})
	}
}
