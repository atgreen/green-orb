package main

import (
	"4d63.com/gocheckcompilerdirectives/checkcompilerdirectives"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checkcompilerdirectives.Analyzer())
}
