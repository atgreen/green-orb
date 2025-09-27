package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/timonwong/loggercheck"
)

func main() {
	singlechecker.Main(loggercheck.NewAnalyzer())
}
