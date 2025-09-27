package asciicheck

import (
	"fmt"
	"testing"
	"unicode"
)

func TestIsASCII(t *testing.T) {
	t.Run("ascii", func(t *testing.T) {
		s := generateString(unicode.MaxASCII + 1)
		ch, ok := isASCII(s)

		if !ok {
			t.Error("expected that string contains only ASCII symbols")
		}

		if ch != 0 {
			t.Error("expected that ch is zero")
		}
	})

	t.Run("non-ascii", func(t *testing.T) {
		s := "привет!"
		ch, ok := isASCII(s)

		if ok {
			t.Error("expected that string contains non-ASCII symbols")
		}

		if ch != 'п' {
			t.Error("expected that ch is equal to first letter")
		}
	})
}

func BenchmarkIsASCII(b *testing.B) {
	// We are usually check small strings, that represents identifiers.
	sizes := []int{
		1, 8, 16, 32,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Len=%d", size), func(b *testing.B) {
			s := generateString(size)
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if _, ok := isASCII(s); !ok {
					b.Fatal("unexpected result")
				}
			}
		})

	}

}

func generateString(l int) string {
	s := make([]byte, l)
	for i := range s {
		s[i] = byte(i)
	}
	return string(s)
}
