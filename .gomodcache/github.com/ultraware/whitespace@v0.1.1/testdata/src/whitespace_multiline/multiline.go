package whitespace

import "fmt"

// No comments, only whitespace.
func fn1() { // want "unnecessary leading newline"

	fmt.Println("Hello, World")

} // want "unnecessary trailing newline"

// Whitespace before leading and after trailing comment.
func fn2() { // want "unnecessary leading newline"

	// Some leading comments - this should yield error.
	fmt.Println("Hello, World")
	// And some trailing comments as well!

} // want "unnecessary trailing newline"

// Whitespace after leading and before trailing comment - no diagnostic.
func fn3() {
	// This is fine, newline is after comment

	fmt.Println("Hello, World")

	// This is fine, newline before comment
}

// MultiFunc (FuncDecl), MultiIf and MultiFunc(FuncLit) without newlinew.
func fn4(
	arg1 int,
	arg2 int,
) { // want "multi-line statement should be followed by a newline"
	fmt.Println("Hello, World")

	if true &&
		false { // want "multi-line statement should be followed by a newline"
		fmt.Println("Hello, World")
	}

	_ = func(
		arg1 int,
		arg2 int,
	) { // want "multi-line statement should be followed by a newline"
		fmt.Println("Hello, World")
	}

	_ = func(
		arg1 int,
		arg2 int,
	) { // want "multi-line statement should be followed by a newline"
		_ = func(
			arg1 int,
			arg2 int,
		) { // want "multi-line statement should be followed by a newline"
			fmt.Println("Hello, World")
		}
	}
}


// MultiFunc (FuncDecl), MultiIf and MultiFunc(FuncLit) with comments counting
// as newlines.
func fn5(
	arg1 int,
	arg2 int,
) {
	// Comment count as newline.
	fmt.Println("Hello, World")

	if true &&
		false {
		// Comment count as newline.
		fmt.Println("Hello, World")
	}

	_ = func(
		arg1 int,
		arg2 int,
	) {
		// Comment count as newline.
		fmt.Println("Hello, World")
	}
}
