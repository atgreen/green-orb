package canonicalheader_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/lasiar/canonicalheader"
)

const testValue = "hello_world"

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := [...]string{
		"alias",
		"assigned",
		"common",
		"const",
		"embedded",
		"global",
		"initialism",
		"struct",
		"underlying",
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt, func(t *testing.T) {
			t.Parallel()

			analysistest.RunWithSuggestedFixes(
				t,
				analysistest.TestData(),
				canonicalheader.Analyzer,
				tt,
			)
		})
	}

	t.Run("are_test_cases_complete", func(t *testing.T) {
		t.Parallel()

		dirs, err := os.ReadDir(filepath.Join(analysistest.TestData(), "src"))
		require.NoError(t, err)
		require.Len(t, testCases, len(dirs))

		require.EqualValues(
			t,
			transform(dirs, func(d os.DirEntry) string {
				return d.Name()
			}),
			testCases,
		)
	})
}

func transform[S ~[]E, E any, T any](sl S, f func(E) T) []T {
	out := make([]T, len(sl))
	for i, t := range sl {
		out[i] = f(t)
	}

	return out
}

func BenchmarkCanonical(b *testing.B) {
	v := http.Header{
		"Canonical-Header": []string{testValue},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := v.Get("Canonical-Header")
		if s != testValue {
			b.Fatal()
		}
	}
}

func BenchmarkNonCanonical(b *testing.B) {
	v := http.Header{
		"Canonical-Header": []string{testValue},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := v.Get("CANONICAL-HEADER")
		if s != testValue {
			b.Fatal()
		}
	}
}
