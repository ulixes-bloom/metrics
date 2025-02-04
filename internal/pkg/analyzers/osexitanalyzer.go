// Package analyzers provides custom static analysis checks for Go code.
//
// It includes analyzers designed to enforce specific coding standards,
// such as restricting the use of os.Exit in the main function.
package analyzers

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// OSExitCheckAnalyzer is a custom static analysis tool that checks
// for calls to os.Exit within the main function.
//
// Using os.Exit directly in main can cause abrupt termination without
// proper cleanup. It's recommended to handle errors gracefully instead.
//
// Example:
//
//	func main() {
//	    os.Exit(1) // This will trigger an analyzer warning
//	}
var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check for os.Exit call in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				ast.Inspect(fn, func(node ast.Node) bool {
					call, ok := node.(*ast.CallExpr)
					if !ok {
						return true
					}

					if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
						if pkgIdent, ok := sel.X.(*ast.Ident); ok {
							if pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
								pass.Reportf(call.Pos(), "os.Exit should not be called in main function")
							}
						}
					}

					return true
				})
			}
		}
	}
	return nil, nil
}
