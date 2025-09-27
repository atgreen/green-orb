package main

import "testing"

func TestIgnoreTests(t *testing.T) {
	t.Log("如果忽略测试就不应该报这一行")
}
