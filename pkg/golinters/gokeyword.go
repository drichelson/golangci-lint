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
	goKeywordErrorMsg    = "detected use of `go` keyword: %s"
	goKeywordDescription = "detects presence of the `go` keyword"
	defaultDetails       = "no details provided"
	//detailsFlag          = "details"
)

func NewGoKeyword(cfg *config.GoKeywordSettings) *goanalysis.Linter {
	//gka2 := goKeywordAnalyzer{
	//	//details: details,
	//}
	//gka2.analyzer = &analysis.Analyzer{
	//	Name:     goKeywordName,
	//	Doc:      goKeywordDescription,
	//	Run:      run,
	//	Requires: []*analysis.Analyzer{inspect.Analyzer},
	//}

	//a :=
	//a.Flags.Init(goKeywordName, flag.ExitOnError)
	//a.Flags.Var(goKeywordAnalyzer, detailsFlag, "Documentation on why this linter is enabled")
	//return a
	//gka2.analyzer.Flags.String(detailsFlag, defaultDetails, "Documentation on why this linter is enabled")
	//gka := &gka2
	//
	//var cfgMap map[string]map[string]interface{}
	//if cfg != nil && cfg.Details != "" {
	//	cfgMap = map[string]map[string]interface{}{
	//		gka.analyzer.Name: {
	//			detailsFlag: cfg.Details,
	//		},
	//	}
	//}

	return goanalysis.NewLinter(
		goKeywordName,
		goKeywordDescription,
		[]*analysis.Analyzer{{
			Name: goKeywordName,
			Doc:  goKeywordDescription,
			Run: func(pass *analysis.Pass) (interface{}, error) {
				details := defaultDetails
				if cfg != nil && cfg.Details != "" {
					details = cfg.Details
				}
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
						pass.Reportf(node.Pos(), goKeywordErrorMsg, details)
					}
				})
				return nil, nil
			},
			Requires: []*analysis.Analyzer{inspect.Analyzer},
		}},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}

//func newGoKeywordAnalyzer() *goKeywordAnalyzer {
//	gka := goKeywordAnalyzer{
//		//details: details,
//	}
//	gka.analyzer = &analysis.Analyzer{
//		Name:     goKeywordName,
//		Doc:      goKeywordDescription,
//		Run:      fun,
//		Requires: []*analysis.Analyzer{inspect.Analyzer},
//	}
//
//	//a :=
//	//a.Flags.Init(goKeywordName, flag.ExitOnError)
//	//a.Flags.Var(goKeywordAnalyzer, detailsFlag, "Documentation on why this linter is enabled")
//	//return a
//	gka.analyzer.Flags.String(detailsFlag, defaultDetails, "Documentation on why this linter is enabled")
//
//	return &gka
//}

//type goKeywordAnalyzer struct {
//	//details  string
//	analyzer *analysis.Analyzer
//}

//func (a *goKeywordAnalyzer) String() string {
//	return a.details
//}
//
//func (a *goKeywordAnalyzer) Set(details string) error {
//	fmt.Println("Setting details to: ", details)
//	a.details = details
//	return nil
//}

//func run(pass *analysis.Pass) (interface{}, error) {
//	details := pass.Analyzer.Flags.Lookup(detailsFlag).Value.String()
//	i, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
//	if !ok {
//		return nil, errors.New("analyzer is not type *inspector.Inspector")
//	}
//
//	nodeFilter := []ast.Node{
//		(*ast.GoStmt)(nil),
//	}
//
//	i.Preorder(nodeFilter, func(node ast.Node) {
//		foundGo := false
//		switch node.(type) {
//		case *ast.GoStmt:
//			foundGo = true
//		}
//		if foundGo {
//			pass.Reportf(node.Pos(), goKeywordErrorMsg, details)
//		}
//	})
//	return nil, nil
//}
