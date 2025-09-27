package main

import no_alias_in_config "fmt" // want `import "fmt" has alias "no_alias_in_config" which is not part of config`

func main() {
	no_alias_in_config.Println("test")
}
