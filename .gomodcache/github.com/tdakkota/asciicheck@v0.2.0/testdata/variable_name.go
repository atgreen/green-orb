package testdata

import "fmt"

func f() {
	var tеstVar int // want `identifier "tеstVar" contain non-ASCII character: U\+0435 'е'`
	tеstVar = 0
	fmt.Println(tеstVar)
}
