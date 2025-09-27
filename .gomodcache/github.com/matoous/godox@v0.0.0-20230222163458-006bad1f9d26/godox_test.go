package godox_test

import (
	"flag"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/matoous/godox"
)

//nolint // reason this is a unit test
func TestParse(t *testing.T) {
	t.Parallel()

	flag.Parse()

	tests := []struct {
		path         string
		result       []string
		includeTests bool
	}{
		{
			path: "./fixtures/00",
			result: []string{
				`fixtures/00/example1.go:3: Line contains TODO/BUG/FIXME: "TODO"`,
			},
		},
		{
			path: "./fixtures/01",
			result: []string{
				`fixtures/01/example1.go:14: Line contains TODO/BUG/FIXME: "TODO(fix): something (Line 13)"`,
				`fixtures/01/example1.go:21: Line contains TODO/BUG/FIXME: "todo compare apples to oranges on a supe..."`,
				`fixtures/01/example1.go:25: Line contains TODO/BUG/FIXME: "TODO: Multiline C1 (Line 24)"`,
				`fixtures/01/example1.go:26: Line contains TODO/BUG/FIXME: "TODO: Multiline C2 (Line 25)"`,
				`fixtures/01/example1.go:27: Line contains TODO/BUG/FIXME: "FIXME: Your attitude (Line 26)"`,
				`fixtures/01/example1.go:28: Line contains TODO/BUG/FIXME: "todo тут какой-то очень-очень-очень-очен..."`,
				`fixtures/01/example2.go:5: Line contains TODO/BUG/FIXME: "TODO: Add JSON tag (Line 4)"`,
				`fixtures/01/example2.go:6: Line contains TODO/BUG/FIXME: "toDO add more fields (Line 5)"`,
				`fixtures/01/example2.go:12: Line contains TODO/BUG/FIXME: "TODO: multiline todo 1 (Line 11)"`,
				`fixtures/01/example2.go:16: Line contains TODO/BUG/FIXME: "TOdo multiline todo 2 (Line 15)"`,
			},
		},
		{
			path: "./fixtures/02",
			result: []string{
				`fixtures/02/example3.go:4: Line contains TODO/BUG/FIXME: "TODO: remove foo (Line 3)"`,
				`fixtures/02/example3.go:8: Line contains TODO/BUG/FIXME: "TODO: Rename field (Line 7)"`,
				`fixtures/02/example3.go:11: Line contains TODO/BUG/FIXME: "TODO: get cat food (Line 10)"`,
				`fixtures/02/example3.go:16: Line contains TODO/BUG/FIXME: "todo  : todo comment (Line 15)"`,
				`fixtures/02/example3_test.go:9: Line contains TODO/BUG/FIXME: "TODO write test"`,
			},
			includeTests: true,
		},
		{
			path: "./fixtures/03",
			result: []string{
				`fixtures/03/main.go:1: Line contains TODO/BUG/FIXME: "TODO: Add package documentation"`,
				`fixtures/03/main.go:2: Line contains TODO/BUG/FIXME: "TODO: Write an actual application"`,
				`fixtures/03/main.go:9: Line contains TODO/BUG/FIXME: "FIXME: Spelling"`,
				`fixtures/03/main.go:14: Line contains TODO/BUG/FIXME: "TODO: Multi line 1"`,
				`fixtures/03/main.go:15: Line contains TODO/BUG/FIXME: "TODO: Multi line 2"`,
				`fixtures/03/main.go:16: Line contains TODO/BUG/FIXME: "FIXME: Mutli line 3"`,
			},
		},
		{
			path: "./fixtures/04",
		},
	}

	for _, tt := range tests {
		tt := tt //nolint // reason tt is ok on this context
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			var messages []godox.Message
			_ = filepath.Walk(tt.path, func(path string, info os.FileInfo, _ error) error {
				fset := token.NewFileSet()
				if info.IsDir() {
					return nil
				}

				f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
				if err != nil {
					panic(err)
				}

				res := godox.Run(f, fset)
				messages = append(messages, res...)

				return nil
			})

			a, b := len(messages), len(tt.result)

			switch {
			case b == 0 && a != b:
				t.Errorf("should expect no messages, instead got:\n%q", messages)
			case a > b:
				t.Errorf("should return less messages (got %d, expects %d)", a, b)
			case a < b:
				t.Errorf("should return more messages (got %d, expect %d)", a, b)
			}

			for i := range tt.result {
				if tt.result[i] != messages[i].Message {
					t.Errorf("not equal\nexpected: %s\nactual: %s", tt.result[i], messages[i])
				}
			}
		})
	}
}
