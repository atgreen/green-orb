package embedded

import "net/http"

type embedded struct {
	http.Header
}

func _() {
	embedded{}.Get("TT") // want `non-canonical header "TT", instead use: "Tt"`
}
