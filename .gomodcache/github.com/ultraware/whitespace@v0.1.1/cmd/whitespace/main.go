package main

import (
	"github.com/ultraware/whitespace"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(whitespace.NewAnalyzer(nil))
}
