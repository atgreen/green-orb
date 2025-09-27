package main

import "testing"

func TestSimple(t *testing.T) {
	var a = []any{1, 2, 3}

	getArgsLength(a) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.any\) int`
	getArgsLength(a...)
	getArgsLength(1, 2, 3)
	getArgsLength([]any{1, 2, 3})                  // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.any\) int`
	getArgsLength(append([]any{1, 2, 3}, 4, 5, 6)) // want `pass \[\]any as any to func getArgsLength func\(args \.\.\.any\) int`

	getOneOrMore(a)
	getOneOrMore(1, a) // want `pass \[\]any as any to func getOneOrMore func\(arg any, others \.\.\.any\) int`
}
