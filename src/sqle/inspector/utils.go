package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
)

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

func getTableName(stmt *ast.TableName) string {
	fmt.Println("table name text:", stmt.Text())
	if stmt.Schema.String() == "" {
		return stmt.Name.String()
	} else {
		return fmt.Sprintf("%s.%s", stmt.Schema, stmt.Name)
	}
}
