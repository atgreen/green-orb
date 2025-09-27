package b

import (
	fff "fmt"
	stdos "os"
	"io"
)

func ImportAsWrongAlias() {
	fff.Println("foo")
	stdos.Stdout.WriteString("bar")
	io.Pipe()
}
