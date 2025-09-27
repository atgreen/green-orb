package testcase

import (
	"fmt"
)

import (
	. "fmt"
)

import (
	_ "fmt"
)

import (
	format "fmt"
)

import (
	"log"
	logger "log"
)

import "sync" // want "should only use grouped 'import' declarations"
import . "sync"
import _ "sync"
import syncer "sync"

import "time"

func dummy() {
	Println("dummy")
	fmt.Println("dummy")
	format.Println("dummy")

	log.Println("dummy")
	logger.Println("dummy")

	_ = Mutex{}
	_ = sync.Mutex{}
	_ = syncer.Mutex{}

	_ = time.Nanosecond
}
