package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/sashamelentyev/interfacebloat/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
