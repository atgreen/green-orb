package varnamelen

import (
	"testing"

	"github.com/matryer/is"
)

func TestStringsValue_Set(t *testing.T) {
	is := is.New(t)
	v := stringsValue{}
	_ = v.Set("foo,bar,baz")
	is.Equal(v.Values, []string{"foo", "bar", "baz"})
}

func TestStringsValue_String(t *testing.T) {
	is := is.New(t)
	v := stringsValue{
		Values: []string{"foo", "bar", "baz"},
	}
	is.Equal(v.String(), "foo,bar,baz")
}
