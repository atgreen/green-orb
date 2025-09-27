// This package implements the diff.Differ interface using github.com/mb0/diff as a backend.
package mb0diff

import (
	"bufio"
	"fmt"
	"io"
	"log"

	mb0diff "github.com/mb0/diff"
)

type differ struct{}

// New returns an implementation of diff.Differ using mb0diff as the backend.
func New() *differ {
	return &differ{}
}

// Implements Data equality comparison on a per-line basis.
type lineDiffer struct {
	a []string
	b []string
}

func (d *lineDiffer) Equal(i, j int) bool {
	if len(d.a) <= i || len(d.b) <= j {
		return false
	}
	return d.a[i] == d.b[j]
}

func (d *lineDiffer) Diff() []mb0diff.Change {
	return mb0diff.Diff(len(d.a), len(d.b), d)
}

func (d *lineDiffer) WriteHunk(out io.Writer, change mb0diff.Change, context int) error {
	// TODO: Write header
	startContext := context
	if change.B-context < 0 {
		startContext = change.B
	}

	endContext := context
	if change.B+context >= len(d.b) {
		endContext = len(d.b) - 1 - change.B
	}
	if endContext < 0 {
		endContext = 0
	}

	totalContext := startContext + endContext

	// +1 to all the things because unified diff headers are 1-indexed
	fromLine, fromNum := change.A-startContext+1, change.Del+totalContext
	toLine, toNum := change.B-startContext+1, change.Ins-change.Del+totalContext

	log.Printf("startContext=%d endContext=%d", startContext, endContext)

	if change.Ins == 0 {
		fmt.Fprintf(out, "@@ -%d,%d +%d @@\n", fromLine, fromNum, toLine)
	} else if change.Del == 0 {
		fmt.Fprintf(out, "@@ -%d +%d,%d @@\n", fromLine, toLine, toNum)
	} else {
		fmt.Fprintf(out, "@@ -%d,%d +%d,%d @@\n", fromLine, fromNum, toLine, toNum)
	}

	// Start context
	for i := startContext; i > 0; i-- {
		fmt.Fprint(out, " ", d.b[change.B-i], "\n")
	}

	// Changes
	for i := 0; i < change.Del; i++ {
		fmt.Fprint(out, "-", d.a[change.A+i], "\n")
	}
	for i := 0; i < change.Ins; i++ {
		fmt.Fprint(out, "+", d.b[change.B+i], "\n")
	}

	// End context
	for i := 0; i < endContext; i++ {
		fmt.Fprint(out, " ", d.b[change.B+i], "\n")
	}
	return nil
}

func readData(a, b io.Reader) (*lineDiffer, error) {
	d := lineDiffer{}
	scanner := bufio.NewScanner(a)
	for scanner.Scan() {
		d.a = append(d.a, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	scanner = bufio.NewScanner(b)
	for scanner.Scan() {
		d.b = append(d.b, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &d, nil
}

// Diff consumes the entire reader streams into memory before generating a diff
// which then gets filled into the buffer. This implementation stores and
// manipulates all three values in memory.
func (diff *differ) Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error {
	d, err := readData(a, b)
	if err != nil {
		return err
	}

	changes := d.Diff()
	for _, change := range changes {
		err := d.WriteHunk(out, change, 3)
		if err != nil {
			return err
		}
	}
	return nil
}
