// Package constructors analyzes type declarations for comments specifying type
// constructors. It reads comments like "// constructors: New, NewT" and exports
// facts connecting type names with lists of declared constructors. It's
// intended to be used as a dependency for the other analyzer in this project
// which actually enforces the usage of those constructors.
package constructors

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type ConstructorDecl struct {
	Name string
	Pos  token.Pos
	// End  token.Pos
}

// ConstructorFact represents the list of type constructors.
type ConstructorFact []ConstructorDecl

func (f *ConstructorFact) AFact() {}

var Analyzer = &analysis.Analyzer{
	Name:      "constructors_facts",
	Doc:       "exports facts about types with constructors declared in a special comment",
	Run:       constructorFactsAnalyzer,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{(*ConstructorFact)(nil)},
	// ResultType: reflect.TypeOf(ConstructorMap{}),
}

// // typeIdent returns either local or imported type ident or nil
// func typeIdent(expr ast.Expr) *ast.Ident {
// 	fmt.Printf("[debug] expr = %#v\n", expr)
// 	switch id := expr.(type) {
// 	case *ast.Ident:
// 		return id
// 	case *ast.SelectorExpr:
// 		// SelectorExpr represents dot-separated names like `bytes.Buffer` - they're used for imported types
// 		// this case returns the part after the last dot (similar to the file name in a path)
// 		return id.Sel
// 	}
// 	return nil
// }

// // ConstructorMap maps type names to their constructor names
// type ConstructorMap map[string][]ConstructorDecl

func constructorFactsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// constructorMap := make(ConstructorMap)
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// TODO: check if a simpler traversal mechanism gives a performance boost
	var err error
	inspector.Preorder(nodeFilter, func(genericNode ast.Node) {
		node, ok := genericNode.(*ast.GenDecl)
		if !ok {
			return
		}
		// only process type declarations
		if node.Tok != token.TYPE {
			return
		}
		// fmt.Printf("[debug] GenDecl = %#v\n", node)

		// "constructors" comments on block declarations with multiple types at once are not supported (see type_decl_block.go in testdata for examples)
		if len(node.Specs) > 1 {
			return
		}

		// if there are no types defined, ignore
		if len(node.Specs) == 0 {
			return
		}

		detectDeclaredConstructors(pass, node)
	})

	return nil, err
}

func detectDeclaredConstructors(pass *analysis.Pass, node *ast.GenDecl) {
	fmt.Printf("[debug] node = %#v\n", node)

	// get the only type declaration
	typeSpec, ok := node.Specs[0].(*ast.TypeSpec)
	if !ok {
		return
	}

	// // get base type if any
	// baseIdent := typeIdent(decl.Type)
	// fmt.Printf("[debug] baseIdent=%v\n", baseIdent)
	// if baseIdent == nil {
	// 	return
	// }

	// // get base type object
	// baseTypeObj := pass.TypesInfo.ObjectOf(baseIdent)
	// fmt.Printf("[debug] baseTypeObj=%v\n", baseTypeObj)
	// if baseTypeObj == nil {
	// 	return
	// }

	// get this type's object
	// pass.TypesInfo.ObjectOf()
	targetTypeObj := pass.TypesInfo.Defs[typeSpec.Name]
	fmt.Printf("[debug] targetTypeObj=%v\n", targetTypeObj)
	if targetTypeObj == nil {
		return
	}

	// // check the base type has a constructor
	// existingFact := new(ConstructorFact)
	// if pass.ImportObjectFact(baseTypeObj, existingFact) {
	// 	fmt.Printf("[debug] derived type has constructors\n")
	// 	// mark derived type as having constructor
	// 	pass.ExportObjectFact(targetTypeObj, existingFact)
	// 	return
	// }

	// // only process type declarations
	// if typeSpec.Tok != token.TYPE {
	// 	return
	// }

	// // "constructors" comments on block declarations with multiple types at once are not supported (see type_decl_block.go in testdata for examples)
	// if len(typeSpec.Specs) > 1 {
	// 	return
	// }

	// // if there are no types defined, ignore
	// if len(typeSpec.Specs) == 0 {
	// 	return
	// }

	// // get the only type declaration
	// targetTypeSpec, ok := typeSpec.Specs[0].(*ast.TypeSpec)
	// if !ok {
	// 	return
	// }

	fmt.Printf("[debug] node.Doc=%v\n", node.Doc)

	// extract constructor names from the comment
	constructorNames := parseConstructorComment(node.Doc)
	if len(constructorNames) == 0 {
		return
	}
	fmt.Printf("[debug] constructorNames=%v\n", constructorNames)

	// Validate that all specified constructors exist in the same package and are valid
	constructors := make([]ConstructorDecl, 0)
	for _, constructorName := range constructorNames {
		if c := constructor(pass, typeSpec, targetTypeObj, constructorName); c != nil {
			constructors = append(constructors, *c)
		}
	}

	if len(constructors) > 0 {
		var fact ConstructorFact = constructors
		pass.ExportObjectFact(targetTypeObj, &fact)
	}
}

func constructor(pass *analysis.Pass, decl *ast.TypeSpec, targetTypeObj types.Object, constructorName string) *ConstructorDecl {
	obj := pass.Pkg.Scope().Lookup(constructorName)
	if obj == nil {
		pass.Reportf(
			decl.Pos(),
			"`%s` marked as a constructor doesn't exist in the package",
			constructorName,
		)
		return nil
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		pass.Reportf(
			decl.Pos(),
			"`%s` marked as a constructor is not a function",
			constructorName,
		)
		return nil
	}

	// I don't check that constructor is exported - this is the problem of the user and they are fully entitled to use this linter whithin one package only

	sig := fn.Signature()
	if sig == nil {
		pass.Reportf(
			decl.Pos(),
			"can't get function signature for `%s` marked as a constructor",
			constructorName,
		)
		return nil
	}

	// check that it's a function and not a method
	if sig.Recv() != nil {
		pass.Reportf(
			decl.Pos(),
			"`%s` marked as a constructor is a method, not a function",
			constructorName,
		)
		return nil
	}

	fnResults := sig.Results() // nil value is valid
	if fnResults.Len() == 0 {
		pass.Reportf(
			decl.Pos(),
			"function `%s` marked as a constructor doesn't return anything",
			constructorName,
		)
		return nil
	}

	// check that a constructor returns T at least once
	// TODO: check it works for *T too
	for i := range fnResults.Len() {
		result := fnResults.At(i)
		if result == nil {
			pass.Reportf(
				decl.Pos(),
				"can't get the type of parameter %d of function `%s` marked as a constructor",
				i+1,
				constructorName,
			)
			return nil
		}

		// if constructor returns T
		if result.Type().String() == targetTypeObj.Type().String() {
			return &ConstructorDecl{
				Name: constructorName,
				Pos:  fn.Pos(),
			}
		}
	}

	pass.Reportf(
		decl.Pos(),
		"function `%s` marked as a constructor doesn't return type %s",
		constructorName,
		decl.Name,
	)

	return nil
}

// parseConstructorComment parses a comment like "//constructors: New, NewT" and
// returns the list of constructor names.
func parseConstructorComment(commentGroup *ast.CommentGroup) []string {
	if commentGroup == nil {
		return nil
	}

	text := commentGroup.Text()
	fmt.Printf("[debug] text = %v\n", text)

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("[debug] line = %v\n", line)

		// Look for "constructors:" comment
		if strings.HasPrefix(line, "constructors:") {
			constructorsListStr := strings.TrimSpace(strings.TrimPrefix(line, "constructors:"))

			// Split by comma and extract constructors names
			names := strings.Split(constructorsListStr, ",")
			var result []string
			for _, name := range names {
				name = strings.TrimSpace(name)
				if name != "" {
					result = append(result, name)
				}
			}
			return result
		}
	}

	return nil
}
