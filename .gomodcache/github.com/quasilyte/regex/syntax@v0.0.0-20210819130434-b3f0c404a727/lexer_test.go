package syntax

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		tokens string
	}{
		{``, ``},

		{`x`, `Char`},
		{`xx`, `Char Concat Char`},
		{`xxx`, `Char Concat Char Concat Char`},
		{`..`, `. Concat .`},
		{`.x.`, `. Concat Char Concat .`},
		{`✓✓`, `Char Concat Char`},

		{`x|x`, `Char | Char`},
		{`x|x|x`, `Char | Char | Char`},
		{`x|xx|xxx`, `Char | Char Concat Char | Char Concat Char Concat Char`},

		{`()`, `( )`},
		{`(x)`, `( Char )`},
		{`((x))`, `( ( Char ) )`},
		{`(x)|x`, `( Char ) | Char`},
		{`x|(x)`, `Char | ( Char )`},
		{`(x)|(x)`, `( Char ) | ( Char )`},
		{`x(x)`, `Char Concat ( Char )`},
		{`(✓x✓x)`, `( Char Concat Char Concat Char Concat Char )`},

		{`(?<1>)`, `(?<name> )`},
		{`(?'1')`, `(?'name' )`},
		{`(?P<1>)`, `(?P<name> )`},
		{`(?P<foo>x)`, `(?P<name> Char )`},
		{`(?<foo>x)`, `(?<name> Char )`},
		{`(?'foo'x)`, `(?'name' Char )`},
		{`(?P<foo>xy)`, `(?P<name> Char Concat Char )`},
		{`a(?P<foo>x)b`, `Char Concat (?P<name> Char ) Concat Char`},
		{`a(?P<foo>xy)b`, `Char Concat (?P<name> Char Concat Char ) Concat Char`},
		{`a(?<foo>xy)b`, `Char Concat (?<name> Char Concat Char ) Concat Char`},
		{`a(?'foo'xy)b`, `Char Concat (?'name' Char Concat Char ) Concat Char`},

		{`(?#)`, `Comment`},
		{`a(?#test)(?#c2)b`, `Char Concat Comment Concat Comment Concat Char`},

		{`(?>)`, `(?> )`},
		{`a(?>xy)(?>z)`, `Char Concat (?> Char Concat Char ) Concat (?> Char )`},

		{`(?=)`, `(?= )`},
		{`(?!)`, `(?! )`},
		{`(?<=)`, `(?<= )`},
		{`(?<!)`, `(?<! )`},
		{`a(?=xy)(?=z)`, `Char Concat (?= Char Concat Char ) Concat (?= Char )`},
		{`a(?!xy)(?!z)`, `Char Concat (?! Char Concat Char ) Concat (?! Char )`},
		{`a(?<=xy)(?<=z)`, `Char Concat (?<= Char Concat Char ) Concat (?<= Char )`},
		{`a(?<!xy)(?<!z)`, `Char Concat (?<! Char Concat Char ) Concat (?<! Char )`},

		{`(?i)`, `(?flags )`},
		{`(?im)`, `(?flags )`},
		{`(?i-m)`, `(?flags )`},
		{`a(?i)b`, `Char Concat (?flags ) Concat Char`},
		{`a(?im)b`, `Char Concat (?flags ) Concat Char`},

		{`(?:)`, `(?flags )`},
		{`(?:xy)`, `(?flags Char Concat Char )`},
		{`(?i:xy)`, `(?flags Char Concat Char )`},
		{`(?im:xy)`, `(?flags Char Concat Char )`},
		{`a(?:)b`, `Char Concat (?flags ) Concat Char`},
		{`a(?:xy)b`, `Char Concat (?flags Char Concat Char ) Concat Char`},
		{`a(?i:xy)b`, `Char Concat (?flags Char Concat Char ) Concat Char`},
		{`a(?-im:xy)b`, `Char Concat (?flags Char Concat Char ) Concat Char`},

		{`\(\)`, `EscapeMeta Concat EscapeMeta`},
		{`\\`, `EscapeMeta`},
		{`\a`, `EscapeChar`},
		{`\\d`, `EscapeMeta Concat Char`},
		{`\d`, `EscapeChar`},
		{`\d\a`, `EscapeChar Concat EscapeChar`},
		{`\dd\a`, `EscapeChar Concat Char Concat EscapeChar`},
		{`\D`, `EscapeChar`},
		{`\s\S`, `EscapeChar Concat EscapeChar`},

		{`-`, `Char`},
		{`[\-]`, `[ EscapeMeta ]`},
		{`a[]a`, `Char Concat [ Char Concat Char`},
		{`[\^a]a`, `[ EscapeChar Char ] Concat Char`},
		{`[^a]a`, `[^ Char ] Concat Char`},
		{`a[^abc]a`, `Char Concat [^ Char Char Char ] Concat Char`},
		{`[[[]a`, `[ Char Char ] Concat Char`},
		{`[\[]a`, `[ EscapeChar ] Concat Char`},
		{`[\]]a`, `[ EscapeMeta ] Concat Char`},
		{`aa[\]1\]]`, `Char Concat Char Concat [ EscapeMeta Char EscapeMeta ]`},
		{`aa[1\]\]2]`, `Char Concat Char Concat [ Char EscapeMeta EscapeMeta Char ]`},
		{`[a-z0-9]a`, `[ Char - Char Char - Char ] Concat Char`},
		{`[0-9-]`, `[ Char - Char - ]`},
		{`[\d-\w]`, `[ EscapeChar - EscapeChar ]`},
		{`[\(-\)]`, `[ EscapeChar - EscapeChar ]`},
		{`[\[-\]]`, `[ EscapeChar - EscapeMeta ]`},

		{`[|]`, `[ Char ]`},
		{`[(-)]`, `[ Char - Char ]`},
		{`[$.+*^?]`, `[ Char Char Char Char Char Char ]`},
		{`[x{1}]`, `[ Char Char Char Char ]`},

		{`[^]`, `[^ Char`},
		{`[^^]`, `[^ Char ]`},

		{`[[:alpha:]]`, `[ PosixClass ]`},
		{`[[:alpha:]-[:blank:]]`, `[ PosixClass - PosixClass ]`},
		{`[[:^word:]]`, `[ PosixClass ]`},
		{`[[:bad:]]`, `[ PosixClass ]`},
		{`[:alpha:]`, `[ Char Char Char Char Char Char Char ]`},

		{`]`, `Char`},
		{`]]`, `Char Concat Char`},

		{`x+`, `Char +`},
		{`x+x+`, `Char + Concat Char +`},
		{`x+?`, `Char + ?`},
		{`x??`, `Char ? ?`},

		{`\pL`, `EscapeUni`},
		{`\pLL`, `EscapeUni Concat Char`},
		{`\p{Greek}`, `EscapeUniFull`},
		{`x\p{^Bad}y`, `Char Concat EscapeUniFull Concat Char`},
		{`\PL`, `EscapeUni`},
		{`\P{^L}`, `EscapeUniFull`},

		{`\0`, `EscapeOctal`},
		{`\01`, `EscapeOctal`},
		{`\012`, `EscapeOctal`},
		{`\777`, `EscapeOctal`},
		{`\78`, `EscapeOctal Concat Char`},
		{`\778`, `EscapeOctal Concat Char`},

		{`\xFF`, `EscapeHex`},
		{`\xab`, `EscapeHex`},
		{`\x10a`, `EscapeHex Concat Char`},
		{`\x1\x2`, `EscapeHex Concat EscapeHex`},

		{`\x{}a`, `EscapeHexFull Concat Char`},
		{`\x{f}a`, `EscapeHexFull Concat Char`},
		{`\x{F1}a`, `EscapeHexFull Concat Char`},

		{`x{10}y`, `Char Repeat Concat Char`},
		{`x{10,}y`, `Char Repeat Concat Char`},
		{`x{10,20}y`, `Char Repeat Concat Char`},
		{`x{1}{2}y`, `Char Repeat Repeat Concat Char`},
		{`ax{10}y`, `Char Concat Char Repeat Concat Char`},
		{`ax{10,}y`, `Char Concat Char Repeat Concat Char`},
		{`ax{10,20}y`, `Char Concat Char Repeat Concat Char`},
		{`ax{1}{2}y`, `Char Concat Char Repeat Repeat Concat Char`},

		{`{}`, `Char Concat Char`},
		{`x{}`, `Char Concat Char Concat Char`},
		{`x{a}`, `Char Concat Char Concat Char Concat Char`},
		{`x{-1}`, `Char Concat Char Concat Char Concat Char Concat Char`},
		{`x{1,b}`, `Char Concat Char Concat Char Concat Char Concat Char Concat Char`},
		{`x{1b}`, `Char Concat Char Concat Char Concat Char Concat Char`},

		{`x\Q`, `Char Concat \Q`},
		{`x\Q.`, `Char Concat \Q`},
		{`x\Q..`, `Char Concat \Q`},
		{`\Q\E`, `\Q`},
		{`\Q..\E`, `\Q`},
		{`x\Q\Ey`, `Char Concat \Q Concat Char`},
		{`x\Q..\Ey`, `Char Concat \Q Concat Char`},
		{`\Q\E\Q\E`, `\Q Concat \Q`},
	}

	removeBrackets := func(s string) string {
		return s[len("[") : len(s)-len("]")]
	}
	var l lexer
	for _, test := range tests {
		l.Init(test.input)
		want := test.tokens
		have := removeBrackets(fmt.Sprint(l.tokens))
		if have != want {
			t.Errorf("tokenize(%q):\nhave: %s\nwant: %s",
				test.input, have, want)
		}
	}
}
