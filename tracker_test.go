package acn

import (
	"go/token"
	"testing"
)

const modProjectHTTPGoFile = "testdata/mod-project/http/http.go"

func TestFindInAssignment(t *testing.T) {
	fset := token.NewFileSet()
	file, err := astForFile(modProjectHTTPGoFile, fset)
	if err != nil {
		t.Fatalf("error reading ast for file: %v", err)
	}
	funcDecl, err := findFunction(file, "Serve", paramCountOpt(2), isFunctionOpt())
	if err != nil {
		t.Fatalf("error finding function: %v", err)
	}
	expr, err := findObjectCreationInFunction("http", modProjectHTTPGoFile, "*github.com/labstack/echo/v4.Echo", funcDecl)
	if err != nil {
		t.Fatalf("error finding object allocation: %v", err)
	}
	if expr == nil {
		t.Fatal("allocation was not found")
	}
}
