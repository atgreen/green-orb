package _struct

import "net/http"

type HeaderStruct struct {
	header http.Header
}

func (h HeaderStruct) _() {
	h.header.Get("TT") // want `non-canonical header "TT", instead use: "Tt"`
}

func _() {
	HeaderStruct{}.header.Get("TT") // want `non-canonical header "TT", instead use: "Tt"`
}
