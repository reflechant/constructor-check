// Package constructorcheck is a go/analysis linter
// that reports types constructed manually while a constructor is present.
// A constructor for type T (only structs are supported at the moment)
// is a function with name starting with "New"
// that returns a value of type T or *T.
package constructorcheck

import (
	"go/ast"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type HasConstructor struct {
	B bool
}

func (f *HasConstructor) AFact() {}

var Analyzer = &analysis.Analyzer{
	Name:     "constructor_check",
	Doc:      "check for types constructed manually ignoring constructor",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	// FactTypes: []analysis.Fact{(*HasConstructor)(nil)},
}

func debugRun(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	inspector.Preorder(nil, func(node ast.Node) {
		log.Printf("%T = %v", node, node)
	})

	return nil, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// nodeFilter := []ast.Node{
	// 	(*ast.CompositeLit)(nil),
	// 	// (*ast.StructType)(nil),
	// 	// (*ast.ValueSpec)(nil),
	// 	(*ast.FuncDecl)(nil),
	// }

	// fmt.Println("DEFS")
	// for ident, obj := range pass.TypesInfo.Defs {
	// 	fmt.Printf("%v: %v\n", ident, obj)
	// }
	// fmt.Println("TYPES")
	// for expr, typeAndValue := range pass.TypesInfo.Types {
	// 	fmt.Printf("%v: %v\n", expr, typeAndValue)
	// }
	// fmt.Println(pass.Pkg.Scope().Names())

	typeConstructors := make(map[types.Type]*ast.FuncDecl)
	typeLiteralNodes := make(map[types.Type][]ast.Node)

	inspector.Preorder(nil, func(node ast.Node) {
		switch decl := node.(type) {
		// case *ast.SelectorExpr:
		// 	obj := ssainfo.Pkg.Pkg.Scope().Lookup(decl.Sel.Name)
		// 	if obj == nil {
		// 		break
		// 	}
		// 	fmt.Println(obj.Type())
		// 	fmt.Printf("%v.%v\n", decl.X, decl.Sel)
		case *ast.CompositeLit:
			ident, ok := decl.Type.(*ast.Ident)
			if !ok {
				break
			}
			obj := pass.TypesInfo.ObjectOf(ident)
			if obj == nil {
				break
			}
			fact := new(HasConstructor)
			if !pass.ImportObjectFact(obj, fact) {
				break
			}
			typeLiteralNodes[obj.Type()] = append(typeLiteralNodes[obj.Type()], node)
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
			fact := HasConstructor{B: true}
			typeConstructors[obj.Type()] = decl
			pass.ExportObjectFact(obj, &fact)
		default:
			// fmt.Printf("%#v\n", node)
		}
	})

	for typ, nodes := range typeLiteralNodes {
		for _, node := range nodes {
			if constructor, ok := typeConstructors[typ]; ok {
				if node.Pos() >= constructor.Pos() &&
					node.Pos() < constructor.End() {
					continue
				}
				pass.Reportf(
					node.Pos(),
					"use constructor %s for type %s instead of a composite literal",
					constructor.Name.Name,
					typ.String())
			}
		}
	}

	return nil, nil
}
