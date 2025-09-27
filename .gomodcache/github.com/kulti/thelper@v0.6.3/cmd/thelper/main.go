package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/kulti/thelper/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
