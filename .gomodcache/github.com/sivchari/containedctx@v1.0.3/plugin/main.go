// This file can build as a plugin for golangci-lint by below command.
//    go build -buildmode=plugin -o path_to_plugin_dir containedctx/plugin/containedctx
// See: https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint

package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/sivchari/containedctx"
)

// flags for Analyzer.Flag.
// If you would like to specify flags for your plugin, you can put them via 'ldflags' as below.
// $ go build -buildmode=plugin -ldflags "-X 'main.flags=-opt val'" containedctx/plugin/containedctx
var flags string

// AnalyzerPlugin provides analyzers as a plugin.
// It follows golangci-lint style plugin.
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	if flags != "" {
		flagset := containedctx.Analyzer.Flags
		if err := flagset.Parse(strings.Split(flags, " ")); err != nil {
			panic("cannot parse flags of containedctx: " + err.Error())
		}
	}
	return []*analysis.Analyzer{
		containedctx.Analyzer,
	}
}
