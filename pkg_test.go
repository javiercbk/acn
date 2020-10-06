package acn

import (
	"go/ast"
	"testing"
)

func TestFindFunction(t *testing.T) {
	f := ast.FuncDecl{}
	functionFound := packageFindFunction("./testdata/mod-project", []string{"modproj/http"}, "customHTTPErrorHandler", &f)
	if !functionFound {
		t.Fatal("expected function to be found")
	}
}
