package main

import "log"

// welp, the order of the parameters is different from the builtin
// copy function!
func copy(src, dst []T) {
	for i := range dst {
		if i == len(src) {
			break
		}
		string := src[i].s
		dst[i].s = string
	}
}

// welp, not the builtin print.
func print(t *T) { log.Printf("{ x=%d, y=%d }", t.x, t.y) }
