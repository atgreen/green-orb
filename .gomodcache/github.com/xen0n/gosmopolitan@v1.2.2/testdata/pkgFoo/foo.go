package main

import (
	"fmt"
	"time"
)

type escapeHatch = string

type someTable struct {
	// ID 主键
	ID uint64 `gorm:"comment:'此行不该被报告'"`
}

func pri18ntln(a ...any) (n int, err error) {
	return fmt.Println(a...)
}

type i18nMessage struct {
	ID       string
	Fallback string
}

var MsgHelloTest = &i18nMessage{
	ID:       "测试消息 ID",
	Fallback: "这两个字符串都不该被报告",
}

func main() {
	_ = "नमस्ते दुनिया" // should only get reported if configured to watch for Devanagari
	fmt.Println("当前系统时间:", time.Now().In(time.Local))
	fmt.Println(escapeHatch("XXX 不应该报告这个"), 123)
	_, _ = pri18ntln("XXX 也不应该报告这个字符串，但应该报出 time.Local", time.Now().In(time.Local))
}
