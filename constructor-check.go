// Package constructorcheck is a go/analysis linter
// that reports types constructed manually while a constructor is present.
// A constructor for type T (only structs are supported at the moment)
// is a function with name starting with "New"
// that returns a value of type T or *T.
package constructorcheck

import (
	"go/ast"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/ssa"
)

type HasConstructor bool

func (f *HasConstructor) AFact() {}

var Analyzer = &analysis.Analyzer{
	Name:      "constructor_check",
	Doc:       "check for types constructed manually ignoring constructor",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer, buildssa.Analyzer},
	FactTypes: []analysis.Fact{},
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
	ssainfo := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	// nodeFilter := []ast.Node{
	// 	(*ast.CompositeLit)(nil),
	// 	// (*ast.StructType)(nil),
	// 	// (*ast.ValueSpec)(nil),
	// 	(*ast.FuncDecl)(nil),
	// }

	typeConstructors := make(map[string]*ast.FuncDecl)
	typeLiteralNodes := make(map[string][]ast.Node)

	// inspector.Preorder(nodeFilter, func(node ast.Node) {
	inspector.Preorder(nil, func(node ast.Node) {
		switch decl := node.(type) {
		case *ast.CompositeLit:
			ident, ok := decl.Type.(*ast.Ident)
			if !ok {
				break
			}
			typeName := ident.Name
			// check if composite literal type exists in the same package,
			// ignore if not
			typ, ok := ssainfo.Pkg.Members[typeName]
			if !ok {
				break
			}
			_, ok = typ.(*ssa.Type)
			if !ok {
				break
			}
			typeLiteralNodes[typeName] = append(typeLiteralNodes[typeName], node)
		// case *ast.Ident:
		// 	if decl.Obj == nil {
		// 		break
		// 	}
		// 	// fmt.Printf("%#v\n", decl.Obj)
		// 	// fmt.Printf("obj %v = %#v\n", decl.Obj.Name, decl.Obj.Decl)
		// 	valueSpec, ok := decl.Obj.Decl.(*ast.ValueSpec)
		// 	if !ok {
		// 		break
		// 	}

		// 	// fmt.Printf("%#v\n", valueSpec)
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

			// check if type T extracted from function name
			// exists in the same package, ignore if not
			typ, ok := ssainfo.Pkg.Members[typeName]
			if !ok {
				break
			}
			_, ok = typ.(*ssa.Type)
			if !ok {
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
			// fact := HasConstructor(true)
			typeConstructors[typeName] = decl
			// pass.ExportObjectFact(ssaType.Object(), &fact)
		default:
			// fmt.Printf("%#v\n", node)
		}
	})

	for typeName, nodes := range typeLiteralNodes {
		for _, node := range nodes {
			if constructor, ok := typeConstructors[typeName]; ok {
				if node.Pos() >= constructor.Pos() &&
					node.Pos() < constructor.End() {
					continue
				}
				pass.Reportf(
					node.Pos(),
					"use constructor %s for type %s instead of a composite literal",
					constructor.Name.Name,
					typeName)
			}
		}
	}

	return nil, nil
}
