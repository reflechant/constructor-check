// Package usage analyzes type usage to check if constructors are being used. Constructors should be marked at the type declaration with a special comment like `// constructors: NewT, NewWithOptions, DefaultConfig`
//
// It consumes facts from the constructors analyzer and reports when types
// are created without using their declared constructors.
package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/reflechant/constructor-check/facts/constructors"
)

var Analyzer = &analysis.Analyzer{
	Name:     "constructorcheck",
	Doc:      "checks if types are created using their declared constructors",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer, constructors.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.TypeSpec)(nil),
	}

	zeroValues := make(map[token.Pos]types.Object)
	nilValues := make(map[token.Pos]types.Object)
	compositeLiterals := make(map[token.Pos]types.Object)
	typeAliases := make(map[types.Object]types.Object)

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		switch decl := node.(type) {
		case *ast.CallExpr:
			// check if it's a new call
			fn, ok := decl.Fun.(*ast.Ident)
			if !ok {
				break
			}
			if fn.Name != "new" {
				break
			}
			// check we have only one argument (the type)
			if len(decl.Args) != 1 {
				break
			}

			ident := typeIdent(decl.Args[0])
			if ident == nil {
				break
			}

			typeObj := pass.TypesInfo.ObjectOf(ident)
			if typeObj == nil {
				break
			}
			zeroValues[node.Pos()] = typeObj
		case *ast.ValueSpec:
			// check it's a pointer value
			starExpr, ok := decl.Type.(*ast.StarExpr)
			if !ok {
				break
			}
			// check it's using a named type
			ident := typeIdent(starExpr.X)
			if ident == nil {
				break
			}
			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}
			nilValues[node.Pos()] = obj
		case *ast.CompositeLit:
			ident := typeIdent(decl.Type)
			if ident == nil {
				break
			}

			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}
			// if it's a zero value literal
			if decl.Elts == nil {
				zeroValues[node.Pos()] = obj
				break
			}
			compositeLiterals[node.Pos()] = obj
		case *ast.TypeSpec:
			// get base type if any
			baseIdent := typeIdent(decl.Type)
			if baseIdent == nil {
				break
			}
			// get base type object
			baseTypeObj := pass.TypesInfo.ObjectOf(baseIdent)
			if baseTypeObj == nil {
				break
			}

			// get this type's object
			typeObj := pass.TypesInfo.ObjectOf(decl.Name)
			if typeObj == nil {
				break
			}
			typeAliases[typeObj] = baseTypeObj
		}
	})

	// Handle type aliases - propagate constructor facts from base types
	for typeObj, baseTypeObj := range typeAliases {
		// check the base type has constructors
		existingFact := new(constructors.ConstructorFact)
		if !pass.ImportObjectFact(baseTypeObj, existingFact) {
			continue
		}

		// mark derived type as having constructors
		newFact := &constructors.ConstructorFact{}
		*newFact = *existingFact
		pass.ExportObjectFact(typeObj, newFact)
	}

	// Report violations
	for pos, obj := range nilValues {
		if constructorNames, ok := getConstructorNames(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"nil value of type %s may be unsafe, use constructor %s instead",
				obj.Type(),
				formatConstructorNames(constructorNames),
			)
		}
	}
	for pos, obj := range zeroValues {
		if constructorNames, ok := getConstructorNames(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"zero value of type %s may be unsafe, use constructor %s instead",
				obj.Type(),
				formatConstructorNames(constructorNames),
			)
		}
	}
	for pos, obj := range compositeLiterals {
		if constructorNames, ok := getConstructorNames(pass, obj, pos); ok {
			pass.Reportf(
				pos,
				"use constructor %s for type %s instead of a composite literal",
				formatConstructorNames(constructorNames),
				obj.Type(),
			)
		}
	}

	return nil, nil
}

// getConstructorNames returns the constructor names for a type if they exist
// and the position is not inside a constructor function.
func getConstructorNames(pass *analysis.Pass, obj types.Object, pos token.Pos) ([]string, bool) {
	// Get constructor facts
	fact := &constructors.ConstructorFact{}
	found := pass.ImportObjectFact(obj, fact)

	if !found {
		return nil, false
	}

	// Extract constructor names from facts
	var constructorNames []string
	for _, decl := range *fact {
		constructorNames = append(constructorNames, decl.Name)
	}

	// Check if we're inside any of the constructor functions
	for _, decl := range *fact {
		if constructorObj := pass.Pkg.Scope().Lookup(decl.Name); constructorObj != nil {
			if fn, ok := constructorObj.(*types.Func); ok {
				// Get the position range of the constructor function
				// This is a simplified check - in practice, we'd need to track
				// the actual function boundaries more precisely
				constructorPos := fn.Pos()
				if pos >= constructorPos && pos < constructorPos+token.Pos(len(decl.Name)) {
					return nil, false
				}
			}
		}
	}

	return constructorNames, true
}

// formatConstructorNames formats a list of constructor names for display.
func formatConstructorNames(names []string) string {
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names, " or ")
}

// typeIdent returns either local or imported type ident or nil
func typeIdent(expr ast.Expr) *ast.Ident {
	switch id := expr.(type) {
	case *ast.Ident:
		return id
	case *ast.SelectorExpr:
		return id.Sel
	}
	return nil
}
