package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/lasiar/canonicalheader"
)

func main() {
	singlechecker.Main(canonicalheader.Analyzer)
}
