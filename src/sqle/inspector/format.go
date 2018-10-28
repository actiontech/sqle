package inspector

import (
	"bytes"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"strings"
)

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
