package main

import (
	"fmt"
	"time"
)

func testC() {
	x := time.Local
	fmt.Println(time.Now().In(x))
}
