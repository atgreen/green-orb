package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/sivchari/containedctx"
)

func main() { unitchecker.Main(containedctx.Analyzer) }
