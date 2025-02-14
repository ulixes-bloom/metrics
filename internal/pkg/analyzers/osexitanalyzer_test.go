package analyzers_test

import (
	"testing"

	"github.com/ulixes-bloom/ya-metrics/internal/pkg/analyzers"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOSExitCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzers.OSExitCheckAnalyzer, "./...")
}
