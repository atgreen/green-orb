package sets

import (
	"flag"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Parallel()

	set := NewString("logr", "logr", "klog")
	assert.Equal(t, []string{"klog", "logr"}, set.List())
	assert.Equal(t, "klog,logr", set.String())
	assert.True(t, set.Has("logr"))
	assert.True(t, set.Has("klog"))
	assert.False(t, set.Has("zap"))
}

func TestString_Flag(t *testing.T) {
	testCases := []struct {
		name      string
		flagValue string
		want      []string
	}{
		{
			name:      "empty",
			flagValue: "",
			want:      nil,
		},
		{
			name:      "klog",
			flagValue: "klog",
			want:      []string{"klog"},
		},
		{
			name:      "klog-and-logr",
			flagValue: "logr,klog",
			want:      []string{"klog", "logr"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := StringSet{}
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			fs.Var(&f, "set", "")

			err := fs.Parse([]string{"-set=" + tc.flagValue})
			require.NoError(t, err)
			assert.Equal(t, tc.want, f.List())
		})
	}
}
