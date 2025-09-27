package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/moricho/tparallel"
)

func main() { unitchecker.Main(tparallel.Analyzer) }
