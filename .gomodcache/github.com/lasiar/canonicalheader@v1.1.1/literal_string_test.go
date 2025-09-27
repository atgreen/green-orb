package canonicalheader

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLiteral(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		value string
		err   string
	}{
		{
			name:  "wrong_utf8",
			value: `"\xF4\x00"`,
			err:   `"\"\\xF4\\x00\"" is not a valid utf8 string`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLiteralString, err := newLiteralString(&ast.BasicLit{
				ValuePos: token.NoPos,
				Kind:     token.STRING,
				Value:    tt.value,
			})

			require.EqualError(t, err, tt.err)
			require.Zero(t, gotLiteralString)
		})
	}
}
