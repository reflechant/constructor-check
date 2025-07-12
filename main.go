package main

import (
	"github.com/reflechant/constructor-check/facts/constructors"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(constructors.Analyzer)
}
