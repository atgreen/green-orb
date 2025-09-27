package testfiles // want `^missing output for example, go test can't validate it$`

import "fmt"

// If file contains only one example and none test/benchmark/fuzz,
// the whole file is treated as an example.
func Example_wholeFileBad() {
	doBad("hello")
}

func doBad(s string) {
	fmt.Println(s)
}
