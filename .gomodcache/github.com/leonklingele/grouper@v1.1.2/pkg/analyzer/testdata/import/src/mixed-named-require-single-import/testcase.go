package testcase

import (
	"fmt"
)

import ( // want "should only use a single 'import' declaration, 10 found"
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

import "sync"
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
