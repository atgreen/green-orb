package comment_test

import (
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMaps_CommentsByPosLine(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path string
		want []string
	}{
		"single": {"testdata/Maps_CommentsByPosLine/single.go", []string{"a"}},
		"multi":  {"testdata/Maps_CommentsByPosLine/multi.go", []string{"b"}},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			ms := maps(t, fset, tt.path)
			p := pos(t, fset, tt.path)
			cgs := ms.CommentsByPosLine(fset, p)
			got := make([]string, len(cgs))
			for i := range cgs {
				got[i] = strings.TrimSpace(cgs[i].Text())
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestMaps_IgnoreLine(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path  string
		check string
		line  int
		want  bool
	}{
		"has_directive": {"testdata/Maps_IgnorePos/has_directive.go", "test-check", 4, true},
		"no_directive":  {"testdata/Maps_IgnorePos/no_directive.go", "test-check", 4, false},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			fset := token.NewFileSet()
			ms := maps(t, fset, tt.path)

			got := ms.IgnoreLine(fset, tt.line, tt.check)
			if tt.want != got {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}
