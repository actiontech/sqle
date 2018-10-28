package inspector

import (
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

//func (rs *InspectResults) add(level, message string) {
//	rs.results = append(rs.results, &Result{
//		Level:   level,
//		Message: message,
//	})
//}

func (rs *InspectResults) add(level, rule string, args ...interface{}) *InspectResults {
	msg := model.RuleMessageMap[rule]
	rs.results = append(rs.results, &Result{
		Level:   level,
		Message: fmt.Sprintf(msg, args...),
	})
	return rs
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

func HasSpecialOption(Options []*ast.ColumnOption, opTp ...ast.ColumnOptionType) bool {
	exists := make(map[ast.ColumnOptionType]bool, len(opTp))
	for _, op := range Options {
		for _, tp := range opTp {
			if op.Tp == tp {
				exists[tp] = true
			}
		}
	}
	for _, tp := range opTp {
		if _, exist := exists[tp]; !exist {
			return false
		}
	}
	return true
}
