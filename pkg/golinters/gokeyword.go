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
)

var details = defaultDetails

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {
	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil && cfg.Details != "" {
		cfgMap[goKeywordName] = map[string]interface{}{"details": cfg.Details}
	}

	return goanalysis.NewLinter(
		goKeywordName,
		goKeywordDescription,
		[]*analysis.Analyzer{newGoKeywordAnalyzer()},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func newGoKeywordAnalyzer() *analysis.Analyzer {
	goKeywordAnalyzer := &goKeywordAnalyzer{Details: defaultDetails}

	a := &analysis.Analyzer{
		Name:     goKeywordName,
		Doc:      goKeywordDescription,
		Run:      goKeywordAnalyzer.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Flags.Init(goKeywordName, flag.ExitOnError)
	a.Flags.Var(goKeywordAnalyzer, "details", "Documentation on why this linter is enabled")
	return a
}

type goKeywordAnalyzer struct {
	Details string
}

func (a *goKeywordAnalyzer) String() string {
	return a.Details
}

func (a *goKeywordAnalyzer) Set(details string) error {
	a.Details = details
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
			pass.Reportf(node.Pos(), goKeywordErrorMsg+": "+details)
		}
	})
	return nil, nil
}
