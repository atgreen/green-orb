package testcase

import "fmt" // want "should only use grouped 'import' declarations"
import "log" // want "should only use a single 'import' declaration, 3 found"
import "sync"

func dummy() { fmt.Println("dummy"); log.Println("dummy"); _ = sync.Mutex{} }
