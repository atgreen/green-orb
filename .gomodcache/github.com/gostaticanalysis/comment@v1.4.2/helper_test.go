package comment_test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/txtar"

	"github.com/gostaticanalysis/comment"
)

func parse(t *testing.T, fset *token.FileSet, path string) []*ast.File {
	t.Helper()
	ar, err := txtar.ParseFile(path)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	files := make([]*ast.File, len(ar.Files))
	for i := range ar.Files {
		n, d := ar.Files[i].Name, ar.Files[i].Data
		f, err := parser.ParseFile(fset, n, d, parser.ParseComments)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		files[i] = f
	}

	return files
}

func maps(t *testing.T, fset *token.FileSet, path string) comment.Maps {
	t.Helper()
	files := parse(t, fset, path)
	return comment.New(fset, files)
}

// pos find position of `_` in source codes as a token.Pos.
func pos(t *testing.T, fset *token.FileSet, path string) token.Pos {
	t.Helper()

	ar, err := txtar.ParseFile(path)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	for i := range ar.Files {
		n, d := ar.Files[i].Name, ar.Files[i].Data
		index := bytes.Index(d, []byte("_"))
		if index == -1 {
			continue
		}

		var pos token.Pos
		fset.Iterate(func(f *token.File) bool {
			if n == f.Name() {
				pos = f.Pos(index)
				return false
			}
			return true
		})

		if pos != token.NoPos {
			return pos
		}
	}

	return token.NoPos
}
