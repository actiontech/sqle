package inspector

import "github.com/pingcap/tidb/ast"

func getTables(stmt *ast.Join) []*ast.TableName {
	tables := []*ast.TableName{}
	if n := stmt.Right; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			tableName, ok := t.Source.(*ast.TableName)
			if ok {
				tables = append(tables, tableName)
			}
		case *ast.Join:
			tables = append(tables, getTables(t)...)
		}
	}
	if n := stmt.Left; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			tableName, ok := t.Source.(*ast.TableName)
			if ok {
				tables = append(tables, tableName)
			}
		case *ast.Join:
			tables = append(tables, getTables(t)...)
		}
	}
	return tables
}
