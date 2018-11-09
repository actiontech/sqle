package inspector

import (
	"errors"
	"fmt"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
	"github.com/pingcap/tidb/mysql"
	"github.com/pingcap/tidb/parser"
	"sqle/model"
	"strings"
)

type InspectResult struct {
	Level   string
	Message string
}

type InspectResults struct {
	results []*InspectResult
}

func newInspectResults() *InspectResults {
	return &InspectResults{
		results: []*InspectResult{},
	}
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

func (rs *InspectResults) add(level, rule string, args ...interface{}) {
	msg := model.RuleMessageMap[rule]
	rs.results = append(rs.results, &InspectResult{
		Level:   level,
		Message: fmt.Sprintf(msg, args...),
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

func getTableNameWithQuote(stmt *ast.TableName) string {
	if stmt.Schema.String() == "" {
		return fmt.Sprintf("`%s`", stmt.Name)
	} else {
		return fmt.Sprintf("`%s`.`%s`", stmt.Schema, stmt.Name)
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

func IsAllInOptions(Options []*ast.ColumnOption, opTp ...ast.ColumnOptionType) bool {
	exists := make(map[ast.ColumnOptionType]bool, len(opTp))
	for _, tp := range opTp {
		for _, op := range Options {
			if tp == op.Tp {
				exists[tp] = true
			}
		}
	}
	// has one no exists, return false
	for _, tp := range opTp {
		if _, exist := exists[tp]; !exist {
			return false
		}
	}
	return true
}

func HasOneInOptions(Options []*ast.ColumnOption, opTp ...ast.ColumnOptionType) bool {
	// has one exists, return true
	for _, tp := range opTp {
		for _, op := range Options {
			if tp == op.Tp {
				return true
			}
		}
	}
	return false
}

func MysqlDataTypeIsBlob(tp byte) bool {
	switch tp {
	case mysql.TypeBlob, mysql.TypeLongBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob:
		return true
	default:
		return false
	}
}

func scanWhereColumn(where ast.ExprNode) bool {
	switch x := where.(type) {
	case nil:
	case *ast.ColumnNameExpr:
		return true
	case *ast.BinaryOperationExpr:
		if scanWhereColumn(x.R) || scanWhereColumn(x.L) {
			return true
		} else {
			return false
		}
	case *ast.IsTruthExpr:
		return scanWhereColumn(x.Expr)
	case *ast.UnaryOperationExpr:
		return scanWhereColumn(x.V)
	case *ast.IsNullExpr:
		return scanWhereColumn(x.Expr)
	case *ast.CompareSubqueryExpr:
		return true
	}
	return false
}

func getAlterTableSpecByTp(specs []*ast.AlterTableSpec, tp ast.AlterTableType) []*ast.AlterTableSpec {
	s := []*ast.AlterTableSpec{}
	for _, spec := range specs {
		if spec.Tp == tp {
			s = append(s, spec)
		}
	}
	return s
}

func newTableName(schema, table string) *ast.TableName {
	return &ast.TableName{
		Name:   _model.NewCIStr(table),
		Schema: _model.NewCIStr(schema),
	}
}

func getPrimaryKey(stmt *ast.CreateTableStmt) (map[string]struct{}, bool) {
	hasPk := false
	pkColumnsName := map[string]struct{}{}
	for _, constraint := range stmt.Constraints {
		if constraint.Tp == ast.ConstraintPrimaryKey {
			hasPk = true
			for _, col := range constraint.Keys {
				pkColumnsName[col.Column.Name.L] = struct{}{}
			}
		}
	}
	return pkColumnsName, hasPk
}
