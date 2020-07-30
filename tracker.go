package acn

import (
	"go/ast"
)

func findObjectCreationInFunction(pkgName, filePath, objectType string, decl *ast.FuncDecl) (ast.Expr, error) {
	for _, stmt := range decl.Body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			// Lhs are the new variables assigned
			expr, err := findObjectCreationInAssignStmt(pkgName, filePath, objectType, s)
			if err == nil {
				return expr, nil
			}
		}
	}
	return nil, ErrNotFound
}

func findObjectCreationInAssignStmt(pkgName, filePath, objectType string, assignStmt *ast.AssignStmt) (ast.Expr, error) {
	for _, expr := range assignStmt.Lhs {
		describeResponse := GuruDescribeResponse{}
		err := RunGuruDescribeQuery(filePath, int(expr.Pos()), &describeResponse)
		if err != nil {
			switch err {
			case ErrGuruAmbiguous:
				// ignore this error
			default:
				return nil, err
			}
		} else {
			if describeResponse.Value != nil {
				if describeResponse.Value.Type == objectType {
					return expr, nil
				}
				if pkgName+"/"+describeResponse.Value.Type == objectType {
					return expr, nil
				}
			}
		}
	}
	return nil, ErrNotFound
}
