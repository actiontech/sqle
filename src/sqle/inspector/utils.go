package inspector

import (
	"bytes"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"strings"
	"sqle/storage"
	"errors"
	"github.com/pingcap/tidb/parser"
)

func parseSql(dbType string, sql string) ([]ast.StmtNode, error) {
	switch dbType {
	case storage.DB_TYPE_MYSQL:
		p := parser.New()
		stmts, err := p.Parse(sql, "", "")
		if err != nil {
			fmt.Printf("parse error: %v\nsql: %v", err, sql)
			return nil, err
		}
		return stmts, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

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
	if stmt.Schema.String() == "" {
		return stmt.Name.String()
	} else {
		return fmt.Sprintf("%s.%s", stmt.Schema, stmt.Name)
	}
}

func alterTableStmtFormat(stmt *ast.AlterTableStmt) string {
	ops := []string{}
	for _, spec := range stmt.Specs {
		ops = append(ops, alterTableSpecFormat(spec))
	}
	return fmt.Sprintf("ALTER TABLE %s\n%s;", getTableName(stmt.Table), strings.Join(ops, ",\n"))
}

var ColumnOptionMap = map[ast.ColumnOptionType]string{
	ast.ColumnOptionNotNull:       "NOT NULL",
	ast.ColumnOptionNull:          "NULL",
	ast.ColumnOptionAutoIncrement: "AUTO_INCREMENT",
	ast.ColumnOptionPrimaryKey:    "PRIMARY KEY",
	ast.ColumnOptionUniqKey:       "UNIQUE KEY",
}

func alterTableSpecFormat(stmt *ast.AlterTableSpec) string {
	switch stmt.Tp {
	case ast.AlterTableRenameTable:
		return fmt.Sprintf("RENAME AS %s", getTableName(stmt.NewTable))
	case ast.AlterTableDropColumn:
		return fmt.Sprintf("DROP COLUMN %s", stmt.OldColumnName)
	case ast.AlterTableAddColumns:
		if len(stmt.NewColumns) == 1 {
			col := stmt.NewColumns[0]
			ops := []string{}
			for _, op := range col.Options {
				switch op.Tp {
				case ast.ColumnOptionDefaultValue:
					ops = append(ops, fmt.Sprintf("DEFAULT %s", exprFormat(op.Expr)))
				case ast.ColumnOptionGenerated:
					v := fmt.Sprintf("GENERATED ALWAYS AS (%s)", exprFormat(op.Expr))
					if op.Stored {
						v = fmt.Sprintf("%s STORED", v)
					}
					ops = append(ops, v)
				case ast.ColumnOptionComment:
					ops = append(ops, fmt.Sprintf("COMMENT %s", exprFormat(op.Expr)))
				default:
					if v, ok := ColumnOptionMap[op.Tp]; ok {
						ops = append(ops, v)
					}
				}
			}
			return fmt.Sprintf("ADD COLUMN %s %s %s", col.Name, col.Tp, strings.Join(ops, " "))
		}
	}
	return ""
}

func exprFormat(node ast.ExprNode) string {
	writer := bytes.NewBufferString("")
	node.Format(writer)
	return writer.String()
}
