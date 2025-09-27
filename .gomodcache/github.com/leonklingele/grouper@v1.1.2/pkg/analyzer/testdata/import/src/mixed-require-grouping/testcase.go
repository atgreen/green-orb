package testcase

import (
	"fmt"
)

import (
	"log"
)

import "sync" // want "should only use grouped 'import' declarations"
import "time"

func dummy() { fmt.Println("dummy"); log.Println("dummy"); _ = sync.Mutex{}; _ = time.Nanosecond }
