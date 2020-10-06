package acn

import (
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

// packageFindFunction finds a function using a directory, package name and function name
func packageFindFunction(rootDir string, pkgPattern []string, funcName string, funcDeclFound *ast.FuncDecl) bool {
	var fset = token.NewFileSet()
	cfg := &packages.Config{Fset: fset, Mode: packages.LoadAllSyntax, Dir: rootDir}
	pkgs, err := packages.Load(cfg, pkgPattern...)
	if err != nil {
		log.Fatal(err)
	}
	for _, pkg := range pkgs {
		if findFunctionInPackage(funcName, pkg, fset, funcDeclFound) {
			return true
		}
	}
	return false
}

func findFunctionInPackage(funcName string, pkg *packages.Package, fset *token.FileSet, funcDeclFound *ast.FuncDecl) bool {
	found := false
	for _, fileAst := range pkg.Syntax {
		ast.Inspect(fileAst, func(n ast.Node) bool {
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == funcName {
					funcDeclFound = funcDecl
					found = true
				}
			}
			return !found
		})
		if !found {
			break
		}
	}
	return found
}
