package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/errname/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
