package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestUnhandledExitCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), UnhandledExitCheckAnalyzer, "./...")
}
