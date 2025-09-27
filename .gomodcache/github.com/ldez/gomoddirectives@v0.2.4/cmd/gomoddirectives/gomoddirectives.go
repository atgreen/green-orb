package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ldez/gomoddirectives"
)

type flagSlice []string

func (f flagSlice) String() string {
	return strings.Join(f, ":")
}

func (f *flagSlice) Set(s string) error {
	*f = append(*f, strings.Split(s, ",")...)
	return nil
}

type config struct {
	ReplaceAllowList          flagSlice
	ReplaceAllowLocal         bool
	ExcludeForbidden          bool
	RetractAllowNoExplanation bool
}

func main() {
	cfg := config{}

	flag.BoolVar(&cfg.ReplaceAllowLocal, "local", false, "Allow local replace directives")
	flag.Var(&cfg.ReplaceAllowList, "list", "List of allowed replace directives")
	flag.BoolVar(&cfg.RetractAllowNoExplanation, "retract-no-explanation", false, "Allow to use retract directives without explanation")
	flag.BoolVar(&cfg.ExcludeForbidden, "exclude", false, "Forbid the use of exclude directives")

	help := flag.Bool("h", false, "Show this help.")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	results, err := gomoddirectives.Analyze(gomoddirectives.Options{
		ReplaceAllowList:          cfg.ReplaceAllowList,
		ReplaceAllowLocal:         cfg.ReplaceAllowLocal,
		ExcludeForbidden:          cfg.ExcludeForbidden,
		RetractAllowNoExplanation: cfg.RetractAllowNoExplanation,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range results {
		fmt.Println(e)
	}

	if len(results) > 0 {
		os.Exit(1)
	}
}

func usage() {
	_, _ = os.Stderr.WriteString(`GoModDirectives

gomoddirectives [flags]

Flags:
`)
	flag.PrintDefaults()
	os.Exit(2)
}
