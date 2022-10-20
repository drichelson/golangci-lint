package golinters

import (
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
	goKeywordErrorMsg    = "detected use of `go` keyword"
	goKeywordDescription = "detects presence of the `go` keyword"
)

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {

	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil {

	}

	return goanalysis.NewLinter(
		goKeywordName,
		goKeywordDescription,
		[]*analysis.Analyzer{goKeywordFinder},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

var goKeywordFinder = &analysis.Analyzer{
	Name:     goKeywordName,
	Doc:      goKeywordDescription,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
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
			pass.Reportf(node.Pos(), goKeywordErrorMsg)
		}
	})
	return nil, nil
}
