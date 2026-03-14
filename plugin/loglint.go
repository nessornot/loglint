package plugin

import (
	"golang.org/x/tools/go/analysis"
	"github.com/golangci/plugin-module-register/register"
	"github.com/nessornot/loglint/internal/analyzer"
)

func init() {
	register.Plugin("loglint", New)
}

func New(conf any) (register.LinterPlugin, error) {
	pluginInstance := &loglintPlugin{}
	return pluginInstance, nil
}

type loglintPlugin struct{}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzers := make([]*analysis.Analyzer, 1)
	analyzers[0] = analyzer.Analyzer
	return analyzers, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
