// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/xen0n/gosmopolitan"
)

func main() {
	singlechecker.Main(gosmopolitan.DefaultAnalyzer)
}
