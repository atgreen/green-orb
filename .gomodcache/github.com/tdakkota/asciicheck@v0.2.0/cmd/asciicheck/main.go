package main

import (
	"github.com/tdakkota/asciicheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(asciicheck.NewAnalyzer())
}
