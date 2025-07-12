// Package constructors analyzes type declarations for comments specifying type
// constructors. It reads comments like "// constructors: New, NewT" and exports
// facts connecting type names with lists of declared constructors. It's
// intended to be used as a dependency for the other analyzer in this project
// which actually enforces the usage of those constructors.
package constructors

import (
	"bufio"
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

// func (c ConstructorDecl) String() string {
// 	return fmt.Sprintf("{%s %d}", c.Name, c.Pos)
// }

// ConstructorFact represents the list of type constructors.
type ConstructorFact []ConstructorDecl

func (f *ConstructorFact) AFact() {}

// func (f ConstructorFact) String() string {
// 	if len(f) == 1 {
// 		return f[0].String()
// 	}
// 	return fmt.Sprintf("%v", []ConstructorDecl(f))
// }

var Analyzer = &analysis.Analyzer{
	Name:      "constructors_facts",
	Doc:       "exports facts about types with constructors declared in a special comment",
	Run:       constructors,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{(*ConstructorFact)(nil)},
	// ResultType: reflect.TypeOf(ConstructorMap{}),
}

// // ConstructorMap maps type names to their constructor names
// type ConstructorMap map[string][]ConstructorDecl

func constructors(pass *analysis.Pass) (interface{}, error) {
	// constructorMap := make(ConstructorMap)
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// TODO: check if a simpler traversal mechanism gives a performance boost
	var err error
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return
		}

		// only process type declarations
		if genDecl.Tok != token.TYPE {
			return
		}

		// "constructors" comments on block declarations with multiple types at once are not supported (see type_decl_block.go in testdata for examples)
		if len(genDecl.Specs) > 1 {
			return
		}

		// if there are no types defined, ignore
		if len(genDecl.Specs) == 0 {
			return
		}

		// get the only type declaration
		targetTypeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return
		}

		targetTypeObj := pass.TypesInfo.ObjectOf(targetTypeSpec.Name)
		if targetTypeObj == nil {
			return
		}

		// extract constructor names from the comment
		constructorNames := parseConstructorComment(genDecl.Doc)
		if len(constructorNames) == 0 {
			return
		}

		// Validate that all specified constructors exist in the same package and are valid

		decls := make([]ConstructorDecl, 0)
	constructorLoop:
		for _, name := range constructorNames {
			obj := pass.Pkg.Scope().Lookup(name)
			if obj == nil {
				pass.Reportf(
					node.Pos(),
					"`%s` marked as a constructor doesn't exist in the package",
					name,
				)
				continue
			}

			fn, ok := obj.(*types.Func)
			if !ok {
				pass.Reportf(
					node.Pos(),
					"`%s` marked as a constructor is not a function",
					name,
				)
				continue
			}

			// I don't check that constructor is exported - this is the problem of the user and they are fully entitled to use this linter whithin one package only

			sig := fn.Signature()
			if sig == nil {
				pass.Reportf(
					node.Pos(),
					"can't get function signature for `%s` marked as a constructor",
					name,
				)
				continue
			}

			// check that is's a function and not a method
			if sig.Recv() != nil {
				pass.Reportf(
					node.Pos(),
					"`%s` marked as a constructor is a method, not a function",
					name,
				)
				continue
			}

			fnResults := sig.Results() // nil value is valid
			if fnResults.Len() == 0 {
				pass.Reportf(
					node.Pos(),
					"function `%s` marked as a constructor doesn't return anything",
					name,
				)
				continue
			}

			// check that a constructor returns T or *T at least once
			for i := range fnResults.Len() {
				result := fnResults.At(i)
				if result == nil {
					pass.Reportf(
						node.Pos(),
						"can't get the type of parameter %d of function `%s` marked as a constructor",
						i+1,
						name,
					)
					continue constructorLoop
				}

				// if constructor returns T
				if result.Type().String() == targetTypeObj.Type().String() {
					decls = append(decls, ConstructorDecl{
						Name: name,
						Pos:  fn.Pos(),
					})
					continue constructorLoop
				}
			}

			pass.Reportf(
				node.Pos(),
				"function `%s` marked as a constructor doesn't return type %s",
				name,
				targetTypeSpec.Name,
			)
		}

		if len(decls) > 0 {
			var fact ConstructorFact = decls
			pass.ExportObjectFact(targetTypeObj, &fact)
		}
	})

	return nil, err
}

// parseConstructorComment parses a comment like "//constructors: New, NewT" and
// returns the list of constructor names.
func parseConstructorComment(commentGroup *ast.CommentGroup) []string {
	if commentGroup == nil {
		return nil
	}

	text := commentGroup.Text()

	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()

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
