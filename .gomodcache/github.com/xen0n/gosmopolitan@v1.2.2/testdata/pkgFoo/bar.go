package main

import (
	. "time"
)

func testA() string {
	return Now().In(Local).String()
}
