package inspector

import (
	"bytes"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/log"
	"strings"
)

func alterTableStmtFormat(stmt *ast.AlterTableStmt) string {
	if len(stmt.Specs) <= 0 {
		return ""
	}
	ops := make([]string, 0, len(stmt.Specs))
	for _, spec := range stmt.Specs {
		ops = append(ops, alterTableSpecFormat(spec))
	}
	return fmt.Sprintf("ALTER TABLE %s\n%s;", getTableNameWithQuote(stmt.Table), strings.Join(ops, ",\n"))
}

func alterTableSpecFormat(stmt *ast.AlterTableSpec) string {
	switch stmt.Tp {
	case ast.AlterTableRenameTable:
		return fmt.Sprintf("RENAME AS %s", getTableNameWithQuote(stmt.NewTable))
	case ast.AlterTableDropColumn:
		return fmt.Sprintf("DROP COLUMN `%s`", stmt.OldColumnName)
	case ast.AlterTableAddColumns:
		if stmt.NewColumns != nil {
			columns := []string{}
			for _, col := range stmt.NewColumns {
				columns = append(columns, columnDefFormat(col))
			}
			if len(columns) == 1 {
				return fmt.Sprintf("ADD COLUMN %s", columns[0])
			} else if len(columns) > 1 {
				return fmt.Sprintf("ADD COLUMN (%s)", strings.Join(columns, ", "))
			}
		}
	case ast.AlterTableChangeColumn:
		if stmt.NewColumns != nil {
			return fmt.Sprintf("CHANGE COLUMN `%s` %s",
				stmt.OldColumnName.Name.String(), columnDefFormat(stmt.NewColumns[0]))
		}
	case ast.AlterTableAlterColumn:
		if stmt.NewColumns != nil {
			col := stmt.NewColumns[0]
			if col.Options != nil {
				return fmt.Sprintf("ALTER COLUMN `%s` SET DEFAULT %s",
					col.Name.Name.String(), exprFormat(col.Options[0].Expr))
			} else {
				return fmt.Sprintf("ALTER COLUMN `%s` DROP DEFAULT",
					col.Name.Name.String())
			}
		}
	case ast.AlterTableAddConstraint:
		var format = ""
		constraint := stmt.Constraint
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			format = "ADD PRIMARY KEY"
		case ast.ConstraintIndex, ast.ConstraintKey:
			format = "ADD INDEX"
		case ast.ConstraintUniqIndex, ast.ConstraintUniqKey, ast.ConstraintUniq:
			format = "ADD UNIQUE INDEX"
		case ast.ConstraintFulltext:
			format = "ADD FULLTEXT INDEX"
		case ast.ConstraintForeignKey:
			format = "ADD FOREIGN KEY"
		default:
			log.NewEntry().Errorf("constraint tp %d not support on format alterTableStmt", constraint.Tp)
		}
		if constraint.Name != "" {
			format = fmt.Sprintf("%s `%s`", format, constraint.Name)
		}
		if indexColums := indexColumnsFormat(constraint.Keys); indexColums != "" {
			format = fmt.Sprintf("%s %s", format, indexColums)
		}
		// if refer is not nil, this is add foreign key stmt.
		if constraint.Refer != nil {
			format = fmt.Sprintf("%s %s", format, referDefFormat(constraint.Refer))
		}
		// if option is not nil, this is add index/primary key stmt.
		if constraint.Option != nil {
			format = fmt.Sprintf("%s %s", format, indexOptionFormat(constraint.Option))
		}
		return format

	case ast.AlterTableDropIndex:
		return fmt.Sprintf("DROP INDEX `%s`", stmt.Name)
	case ast.AlterTableDropPrimaryKey:
		return fmt.Sprintf("DROP PRIMARY KEY")
	case ast.AlterTableDropForeignKey:
		return fmt.Sprintf("DROP FOREIGN KEY `%s`", stmt.Name)
	case ast.AlterTableRenameIndex:
		return fmt.Sprintf("RENAME INDEX `%s` TO `%s`", stmt.FromKey, stmt.ToKey)
	}
	return ""
}

var ColumnOptionMap = map[ast.ColumnOptionType]string{
	ast.ColumnOptionNotNull:       "NOT NULL",
	ast.ColumnOptionNull:          "NULL",
	ast.ColumnOptionAutoIncrement: "AUTO_INCREMENT",
	ast.ColumnOptionPrimaryKey:    "PRIMARY KEY",
	ast.ColumnOptionUniqKey:       "UNIQUE KEY",
}

func columnDefFormat(col *ast.ColumnDef) string {
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
			ops = append(ops, fmt.Sprintf("COMMENT '%s'", exprFormat(op.Expr)))
		default:
			if v, ok := ColumnOptionMap[op.Tp]; ok {
				ops = append(ops, v)
			}
		}
	}
	return fmt.Sprintf("`%s` %s %s", col.Name, col.Tp, strings.Join(ops, " "))
}

func exprFormat(node ast.ExprNode) string {
	switch node.(type) {
	case *ast.DefaultExpr:
		return "DEFAULT"
	default:
		writer := bytes.NewBufferString("")
		node.Format(writer)
		return writer.String()
	}
}

func indexOptionFormat(op *ast.IndexOption) string {
	if op == nil {
		return ""
	}
	ops := make([]string, 0, 3)
	if op.Tp.String() != "" {
		ops = append(ops, fmt.Sprintf("USING %s", op.Tp.String()))
	}
	if op.KeyBlockSize != 0 {
		ops = append(ops, fmt.Sprintf("KEY_BLOCK_SIZE=%d", op.KeyBlockSize))
	}
	if op.Comment != "" {
		ops = append(ops, fmt.Sprintf("COMMENT '%s'", op.Comment))
	}
	if len(ops) > 0 {
		return fmt.Sprintf("%s", strings.Join(ops, " "))
	}
	return ""
}

func indexColumnsFormat(keys []*ast.IndexColName) string {
	if keys == nil {
		return ""
	}
	columnsName := make([]string, 0, len(keys))
	for _, key := range keys {
		columnsName = append(columnsName, fmt.Sprintf("`%s`", key.Column.Name.String()))
	}
	if len(columnsName) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(columnsName, ","))
	}
	return ""
}

func referDefFormat(refer *ast.ReferenceDef) string {
	if refer == nil {
		return ""
	}
	tableName := getTableNameWithQuote(refer.Table)
	indexColumns := indexColumnsFormat(refer.IndexColNames)
	format := fmt.Sprintf("REFERENCES %s %s", tableName, indexColumns)
	if refer.OnDelete.ReferOpt != ast.ReferOptionNoOption {
		format = fmt.Sprintf("%s ON DELETE %s", format, refer.OnDelete.ReferOpt)
	}
	if refer.OnUpdate.ReferOpt != ast.ReferOptionNoOption {
		format = fmt.Sprintf("%s ON UPDATE %s", format, refer.OnUpdate.ReferOpt)
	}
	return format
}
