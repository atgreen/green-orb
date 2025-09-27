package syntax

import (
	"fmt"
	"testing"
)

func TestParserErrorsPCRE(t *testing.T) {
	tests := []struct {
		pattern string
		want    string
	}{
		{``, `empty pattern: can't find delimiters`},
		{`aba`, `'a' is not a valid delimiter`},
		{` aa `, `whitespace is not a valid delimiter`},
		{`/abc`, `can't find '/' ending delimiter`},
		{`#abc`, `can't find '#' ending delimiter`},
	}

	p := NewParser(nil)
	for _, test := range tests {
		_, err := p.ParsePCRE(test.pattern)
		have := "<nil>"
		if err != nil {
			have = err.Error()
		}
		if have != test.want {
			t.Errorf("parse(%q):\nhave: %s\nwant: %s",
				test.pattern, have, test.want)
		}
	}
}

func TestParsePCRE(t *testing.T) {
	tests := []struct {
		source string

		wantPattern   string
		wantDelim     string
		wantModifiers string
	}{
		{`@@`, "", "@@", ""},
		{`//i`, "", "//", "i"},
		{`#hello#`, "hello", "##", ""},
		{`{pcre pattern}smi`, "pcre pattern", "{}", "smi"},
		{`<an[o]ther (example)!>ms`, "an[o]ther (example)!", "<>", "ms"},
		{`/clipFrom/([0-9]+)`, "clipFrom", "//", "([0-9]+)"},
	}

	p := NewParser(nil)
	for _, test := range tests {
		pcre, err := p.ParsePCRE(test.source)
		if err != nil {
			t.Fatalf("parse(%q): error: %v", test.source, err)
		}
		if pcre.Pattern != test.wantPattern {
			t.Fatalf("parse(%q): pattern mismatch:\nhave: `%s`\nwant: `%s`",
				test.source, pcre.Pattern, test.wantPattern)
		}
		haveDelim := fmt.Sprintf("%c%c", pcre.Delim[0], pcre.Delim[1])
		if haveDelim != test.wantDelim {
			t.Fatalf("parse(%q): delimiter mismatch:\nhave: `%s`\nwant: `%s`",
				test.source, haveDelim, test.wantDelim)
		}
		if pcre.Modifiers != test.wantModifiers {
			t.Fatalf("parse(%q): modifiers mismatch:\nhave: `%s`\nwant: `%s`",
				test.source, pcre.Modifiers, test.wantModifiers)
		}
	}
}
