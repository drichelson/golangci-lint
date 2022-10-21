package golinters

import (
	"flag"
	"go/ast"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	goKeywordName        = "gokeyword"
	goKeywordErrorMsg    = "detected use of `go` keyword: %s"
	goKeywordDescription = "detects presence of the `go` keyword"
	defaultDetails       = "no details provided"
	detailsFlag          = "details"
)

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {
	a := newGoKeywordAnalyzer()

	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil && cfg.Details != "" {
		cfgMap[a.Name] = map[string]interface{}{detailsFlag: cfg.Details}
	}

	return goanalysis.NewLinter(
		a.Name,
		goKeywordDescription,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func newGoKeywordAnalyzer() *analysis.Analyzer {
	goKeywordAnalyzer := &goKeywordAnalyzer{details: defaultDetails}

	a := &analysis.Analyzer{
		Name:     goKeywordName,
		Doc:      goKeywordDescription,
		Run:      goKeywordAnalyzer.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Flags.Init(goKeywordName, flag.ExitOnError)
	a.Flags.Var(goKeywordAnalyzer, detailsFlag, "Documentation on why this linter is enabled")
	return a
}

type goKeywordAnalyzer struct {
	details string
}

func (a *goKeywordAnalyzer) String() string {
	return a.details
}

func (a *goKeywordAnalyzer) Set(details string) error {
	a.details = details
	return nil
}

func (a *goKeywordAnalyzer) run(pass *analysis.Pass) (interface{}, error) {
	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.GoStmt)(nil),
	}

	i.Preorder(nodeFilter, func(node ast.Node) {
		foundGo := false
		switch node.(type) {
		case *ast.GoStmt:
			foundGo = true
		}
		if foundGo {
			pass.Reportf(node.Pos(), goKeywordErrorMsg, a.details)
		}
	})
	return nil, nil
}
