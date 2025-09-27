package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/maratori/testableexamples/pkg/testableexamples"
)

func main() {
	singlechecker.Main(testableexamples.NewAnalyzer())
}
