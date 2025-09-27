package main

import (
	"gitlab.com/bosi/decorder"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(decorder.Analyzer)
}
