// This is the package comment.
package main

import "fmt"

//revive:disable
// This comment is associated with the hello constant.x
const hello = "Hello, World!" // line comment 1

// This comment is associated with the foo variable.
var foo = hello // line comment 2

type X struct {
	// TODO(fix): something (Line 13)
	Name string
}

// This comment is associated with the main function.
func main() {
	fmt.Println(hello) // line comment 3
	//  todo compare apples to oranges on a super super mega mega long long long unsigned line with one big comment (Line 20)
	// something
}

//TODO: Multiline C1 (Line 24)
//TODO: Multiline C2 (Line 25)
//FIXME: Your attitude (Line 26)
// todo тут какой-то очень-очень-очень-очень длинный комментарий про utf-8
