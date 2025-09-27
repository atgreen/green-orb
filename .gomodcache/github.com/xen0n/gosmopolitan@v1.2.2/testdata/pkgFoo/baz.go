package main

import (
	stdTime "time"
)

func testB() int64 {
	return stdTime.Now().In(stdTime.Local).UnixNano()
}
