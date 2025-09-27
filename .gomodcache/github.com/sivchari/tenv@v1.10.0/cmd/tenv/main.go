package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/sivchari/tenv"
)

func main() { unitchecker.Main(tenv.Analyzer) }
