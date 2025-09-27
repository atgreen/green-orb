package _const

import "net/http"

const (
	noUsage     = ""
	noCanonical = `TT`
	canonical   = "Tt"
)

const copiedFromNoCanonical = noCanonical

type myString string

const underlyingString myString = "TT"

func _() {
	var mstr myString = "Tt"
	http.Header{}.Get(string(mstr))

	http.Header{}.Get(string(underlyingString)) // want `const "underlyingString" used as a key at http.Header, but "TT" is not canonical, want "Tt"`
	http.Header{}.Get(string(underlyingString)) // want `const "underlyingString" used as a key at http.Header, but "TT" is not canonical, want "Tt"`
	http.Header{}.Get(noCanonical)              // want `const "noCanonical" used as a key at http.Header, but "TT" is not canonical, want "Tt"`
	http.Header{}.Get(copiedFromNoCanonical)    // want `const "copiedFromNoCanonical" used as a key at http.Header, but "TT" is not canonical, want "Tt"`
	http.Header{}.Get(canonical)
}
