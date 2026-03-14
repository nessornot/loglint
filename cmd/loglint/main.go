package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"github.com/nessornot/loglint/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
