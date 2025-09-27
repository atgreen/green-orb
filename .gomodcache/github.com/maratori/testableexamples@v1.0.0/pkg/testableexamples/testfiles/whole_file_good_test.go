package testfiles

import "fmt"

// If file contains only one example and none test/benchmark/fuzz,
// the whole file is treated as an example.
func Example_wholeFileGood() {
	doGood("hello")
	// Output: hello
}

func doGood(s string) {
	fmt.Println(s)
}
