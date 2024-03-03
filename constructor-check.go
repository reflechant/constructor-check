// Package constructorcheck is a go/analysis linter
// that reports types constructed manually while a constructor is present.
// A constructor for type T (only structs are supported at the moment)
// is a function with name starting with "New"
// that returns a value of type T or *T.
package constructorcheck

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type ConstructorFact struct {
	ConstructorName string
	Pos             token.Pos
	End             token.Pos
}

func (f *ConstructorFact) AFact() {}

var Analyzer = &analysis.Analyzer{
	Name:      "constructor_check",
	Doc:       "check for types constructed manually ignoring constructor",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{(*ConstructorFact)(nil)},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		switch decl := node.(type) {
		case *ast.CompositeLit:
			var ident *ast.Ident

			// select between native types and imported
			switch id := decl.Type.(type) {
			case *ast.Ident:
				ident = id
			case *ast.SelectorExpr:
				ident = id.Sel
			}

			if ident == nil {
				break
			}

			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}

			fact := new(ConstructorFact)
			if !pass.ImportObjectFact(obj, fact) {
				break
			}
			// if composite literal is inside it's own constructor - ignore
			if node.Pos() >= fact.Pos &&
				node.Pos() < fact.End {
				break
			}

			pass.Reportf(
				node.Pos(),
				"use constructor %s for type %s instead of a composite literal",
				fact.ConstructorName,
				obj.Type())

		case *ast.FuncDecl:
			// check if it's a function not a method
			if decl.Recv != nil {
				break
			}

			// check if function name starts with "New"
			if !strings.HasPrefix(decl.Name.Name, "New") {
				break
			}

			// check if function name follows the NewT template
			// TODO: think about easing this requirement because often
			// they rename types and forget to rename constructors
			typeName, ok := strings.CutPrefix(decl.Name.Name, "New")
			if !ok {
				break
			}

			// check if type T extracted from function name exists
			obj := pass.Pkg.Scope().Lookup(typeName)
			if obj == nil {
				break
			}
			// check if supposed constructor returns exactly one value
			// TODO: implement other cases ?
			// (T, err), (*T, err), (T, bool), (*T, bool)
			returns := decl.Type.Results.List
			if len(returns) != 1 {
				break
			}
			// to be done later:
			// // check if supposed constructor returns a value of type T or *T
			// // declared in the same package and T equals extracted type name

			// assume we have a valid constructor
			fact := ConstructorFact{
				ConstructorName: decl.Name.Name,
				Pos:             decl.Pos(),
				End:             decl.End(),
			}
			pass.ExportObjectFact(obj, &fact)
		default:
			// fmt.Printf("%#v\n", node)
		}
	})

	return nil, nil
}
