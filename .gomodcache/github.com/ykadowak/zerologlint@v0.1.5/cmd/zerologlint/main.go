package main

import (
	"github.com/ykadowak/zerologlint"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(zerologlint.Analyzer) }
