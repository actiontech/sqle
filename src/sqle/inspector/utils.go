package inspector

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/parser"
	"sqle/model"
	"strings"
)

type Result struct {
	Level   string
	Message string
}

type InspectResults struct {
	results []*Result
}

func newInspectResults(results ...*Result) *InspectResults {
	ir := &InspectResults{
		results: []*Result{},
	}
	ir.results = append(ir.results, results...)
	return ir
}

// level find highest level in result
func (rs *InspectResults) level() string {
	level := model.RULE_LEVEL_NORMAL
	for _, result := range rs.results {
		if model.RuleLevelMap[level] < model.RuleLevelMap[result.Level] {
			level = result.Level
		}
	}
	return level
}

func (rs *InspectResults) message() string {
	messages := make([]string, len(rs.results))
	for n, result := range rs.results {
		messages[n] = fmt.Sprintf("[%s]%s", result.Level, result.Message)
	}
	return strings.Join(messages, "\n")
}

func (rs *InspectResults) add(level, message string) {
	rs.results = append(rs.results, &Result{
		Level:   level,
		Message: message,
	})
}

func parseSql(dbType, sql string) ([]ast.StmtNode, error) {
	switch dbType {
	case model.DB_TYPE_MYSQL:
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

func parseOneSql(dbType, sql string) (ast.StmtNode, error) {
	switch dbType {
	case model.DB_TYPE_MYSQL:
		p := parser.New()
		stmt, err := p.ParseOneStmt(sql, "", "")
		if err != nil {
			fmt.Printf("parse error: %v\nsql: %v", err, sql)
			return nil, err
		}
		return stmt, nil
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

func SplitSql(dbType, sql string) ([]string, error) {
	stmts, err := parseSql(dbType, sql)
	if err != nil {
		return nil, err
	}
	sqlArray := make([]string, len(stmts))
	for n, stmt := range stmts {
		sqlArray[n] = stmt.Text()
	}
	return sqlArray, nil
}

func RemoveArrayRepeat(input []string) (output []string) {
	for _, i := range input {
		repeat := false
		for _, j := range output {
			if i == j {
				repeat = true
				break
			}
		}
		if !repeat {
			output = append(output, i)
		}
	}
	return output
}
