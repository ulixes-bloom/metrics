// Staticlint is a static analysis tool that integrates multiple analyzers.
//
// The analyzers include:
//   - Standard Go analyzers like printf, shadow, structtag.
//   - Custom analyzers (e.g., OSExitCheckAnalyzer).
//   - Selective staticcheck analyzers based on specified IDs (e.g., "SA4006", "SA5000", "SA6000").
//   - Unchecked errors
//
// Usage:
//
//	go build .
//	./staticlint ./...
package main

import (
	"slices"
	"strings"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/analyzers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/kisielk/errcheck/errcheck"
)

// staticcheckAnalyzers defines a whitelist of staticcheck rules (by ID)
// that will be included in the analyzer suite.
var (
	staticcheckAnalyzers = []string{"SA4006", "SA5000", "SA6000"}
)

// main initializes and runs the multichecker with a combination of
// built-in, third-party, and custom analyzers to perform static code analysis.
//
// The analyzers include:
//   - Standard Go analyzers like printf, shadow, structtag.
//   - Custom analyzers (e.g., OSExitCheckAnalyzer).
//   - Selective staticcheck analyzers based on specified IDs.
//
// Usage:
//
//	Run this binary as part of CI/CD or as a standalone static analysis tool:
//	  go run main.go ./...
func main() {
	analyzers := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		analyzers.OSExitCheckAnalyzer,
		errcheck.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") || slices.Contains(staticcheckAnalyzers, v.Analyzer.Name) {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	multichecker.Main(
		analyzers...,
	)
}
