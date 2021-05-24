package acn

import (
	"go/ast"
	"strings"
)

type ImportsParser func(astFile *ast.File) ([]GoImport, error)

type GoImport struct {
	Path        string
	PackageName string
	Alias       string
}

func ParseGoImports(astFile *ast.File) []GoImport {
	goImports := make([]GoImport, 0, len(astFile.Imports))
	for _, imp := range astFile.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		packageName := path[:]
		lastSlashIndex := strings.LastIndex(path, "/")
		if lastSlashIndex != -1 {
			packageName = packageName[lastSlashIndex+1:]
		}
		// handle gopkg.in/package.v#
		lastDotIndex := strings.LastIndex(packageName, ".")
		if lastDotIndex != -1 {
			packageName = packageName[:lastDotIndex]
		}
		goImport := GoImport{
			Path:        path,
			PackageName: packageName,
		}
		if imp.Name != nil {
			goImport.Alias = imp.Name.Name
		}
		goImports = append(goImports, goImport)
	}
	return goImports
}
