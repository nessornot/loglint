package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/nessornot/loglint/internal/analyzer"
)

// golangci-lint entry point
func New(conf any) ([]*analysis.Analyzer, error) {
	// TODO: handle custom regexp
	
	analyzers := make([] *analysis.Analyzer, 1)
	analyzers[0] = analyzer.Analyzer
	
	return analyzers, nil
}
