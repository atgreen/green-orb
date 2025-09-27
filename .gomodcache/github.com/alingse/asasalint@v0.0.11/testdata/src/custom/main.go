package main

func main() {
	var a = []any{1, 2, 3}

	getArgsLength(a)
	getArgsLength(a...)
	getArgsLength(1, 2, 3)
	getArgsLength([]any{1, 2, 3})
	getArgsLength(append([]any{1, 2, 3}, 4, 5, 6))

	getOneOrMore(a)
	getOneOrMore(1, a)
}

func getArgsLength(args ...any) int {
	return len(args)
}

func getOneOrMore(arg any, others ...any) int {
	return len(append([]any{arg}, others...))
}
