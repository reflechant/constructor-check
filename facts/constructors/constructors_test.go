package constructors

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConstructors(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "testpackage")
}
