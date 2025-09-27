package assigned

import (
	"fmt"
	"net/http"
)

func _() {
	h := http.Header{}

	i, g := 0, h.Del
	fmt.Println(i)
	g("TT") // want `non-canonical header "TT", instead use: "Tt"`

	f := h.Get
	f("TT") // want `non-canonical header "TT", instead use: "Tt"`
}
