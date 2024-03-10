// Package analyzer is a linter that reports ignored constructors.
// It shows you places where someone is doing T{} or &T{}
// instead of using NewT declared in the same package as T.
// A constructor for type T (only structs are supported at the moment)
// is a function with name "NewT" that returns a value of type T or *T.
// Types returned by constructors are not checked right now,
// only that type T inferred from the function name exists in the same package.
// Standard library packages are excluded from analysis.
package analyzer

import (
	"go/ast"
	"go/token"
	"log"
	"os/exec"
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

var stdPackages = stdPackageNames()

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CompositeLit)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.CallExpr)(nil),
	}

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
			// select between native types and imported
			ident := typeIdent(decl.Args[0])
			if ident == nil {
				break
			}

			// check the type has a constructor
			typeObj := pass.TypesInfo.ObjectOf(ident)
			if typeObj == nil {
				break
			}
			fact := new(ConstructorFact)
			if !pass.ImportObjectFact(typeObj, fact) {
				break
			}

			// if new(T) is called inside T's constructor - ignore
			if node.Pos() >= fact.Pos &&
				node.Pos() < fact.End {
				break
			}

			pass.Reportf(
				node.Pos(),
				"nil value of type %s may be unsafe to use, use constructor %s instead",
				typeObj.Type(),
				fact.ConstructorName,
			)
		case *ast.CompositeLit:
			// select between native types and imported
			ident := typeIdent(decl.Type)
			if ident == nil {
				break
			}

			// check the type has a constructor
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
				obj.Type(),
			)
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

			// ignore standard library types
			if _, ok := stdPackages[obj.Pkg().Name()]; ok {
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

func stdPackageNames() map[string]struct{} {
	// inspired by https://pkg.go.dev/golang.org/x/tools/go/packages#Load
	cmd := exec.Command("go", "list", "std")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal("can't load standard library package names")
	}
	pkgs := strings.Fields(string(output))

	stdPkgNames := make(map[string]struct{}, len(pkgs))
	for _, pkg := range pkgs {
		stdPkgNames[pkg] = struct{}{}
	}
	return stdPkgNames
}