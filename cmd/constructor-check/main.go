package main

import (
	constructorcheck "github.com/reflechant/constructor-check"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(constructorcheck.Analyzer)
}
