package e

import (
	wrong_alias "fmt"      // want `import "fmt" imported as "wrong_alias" but must be "fff" according to config`
	"io"                   // want `import "io" imported without alias but must be with alias "iio" according to config`
	wrong_alias_again "os" // want `import "os" imported as "wrong_alias_again" but must be "stdos" according to config`
)

func ImportAsWrongAlias() {
	wrong_alias.Println("foo")
	wrong_alias_again.Stdout.WriteString("bar")
	io.Pipe()
}
