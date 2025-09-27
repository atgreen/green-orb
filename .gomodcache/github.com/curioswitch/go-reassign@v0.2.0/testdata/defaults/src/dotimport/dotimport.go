package dotimport

import . "io"

func direct() {
	EOF = nil // want "reassigning variable"
}
