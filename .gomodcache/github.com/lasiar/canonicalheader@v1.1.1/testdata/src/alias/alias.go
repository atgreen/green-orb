package alias

import "net/http"

type myHeader = http.Header

func _() {
	myHeader{}.Get("TT") // want `non-canonical header "TT", instead use: "Tt"`
}
