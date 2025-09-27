// This package implements the diff.Differ interface using diffmatchpatch as a backend
package dmp

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type differ struct {
	dmp *diffmatchpatch.DiffMatchPatch
}

// New returns an implementation of diff.Differ using diffmatchpatch as the backend.
func New() *differ {
	return &differ{
		dmp: diffmatchpatch.New(),
	}
}

// Diff consumes the entire reader streams into memory before generating a diff
// which then gets filled into the buffer. This implementation stores and
// manipulates all three values in memory.
//
// It's essentially an io wrapper around diffmatchpatch's PatchMake and PatchToText.
func (diff *differ) Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error {
	var src, dst []byte
	var err error

	if src, err = ioutil.ReadAll(a); err != nil {
		return err
	}
	if dst, err = ioutil.ReadAll(b); err != nil {
		return err
	}

	patch := diff.dmp.PatchMake(string(src), string(dst))
	diffText := strings.Replace(diff.dmp.PatchToText(patch), "%0A", "", -1)
	_, err = fmt.Fprint(out, diffText)
	return err
}
