package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsASCII(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty",
			input: "",
			want:  true,
		},
		{
			name:  "simple-ascii",
			input: "abcdefg",
			want:  true,
		},
		{
			name:  "cjk",
			input: "中文日a本b語ç日ð本Ê語þ日¥本¼語i日©",
			want:  false,
		},
		{
			name:  "emoji",
			input: "☺☻☹",
			want:  false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsASCII(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
