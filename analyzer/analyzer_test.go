package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUsage(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "p")
}
