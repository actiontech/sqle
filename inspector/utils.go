package inspector

import (
	"errors"
	"fmt"
	"github.com/pingcap/tidb/ast"
	_model "github.com/pingcap/tidb/model"
	"github.com/pingcap/tidb/mysql"
	"github.com/pingcap/tidb/parser"
	"regexp"
	"sqle/model"
	"strconv"
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

func (rs *InspectResults) add(rule model.Rule, args ...interface{}) {
	rs.results = append(rs.results, &InspectResult{
		Level:   rule.Level,
		Message: fmt.Sprintf(RuleHandlerMap[rule.Name].Message, args...),
	})
}

func parseSql(dbType, sql string) ([]ast.StmtNode, error) {
	switch dbType {
	case model.DB_TYPE_MYSQL, model.DB_TYPE_MYCAT:
		p := parser.New()
		stmts, err := p.Parse(sql, "", "")
		if err != nil {
			return nil, err
		}
		return stmts, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

func parseOneSql(dbType, sql string) (ast.StmtNode, error) {
	switch dbType {
	case model.DB_TYPE_MYSQL, model.DB_TYPE_MYCAT:
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
	if stmt == nil {
		return tables
	}
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

func whereStmtHasOneColumn(where ast.ExprNode) bool {
	return scanWhereStmtColumn(where, func(expr *ast.ColumnNameExpr) bool {
		return true
	})
}

func whereStmtHasSpecificColumn(where ast.ExprNode, columnName string) bool {
	return scanWhereStmtColumn(where, func(expr *ast.ColumnNameExpr) bool {
		if expr.Name.Name.L == strings.ToLower(columnName) {
			return true
		}
		return false
	})
}

func scanWhereStmtColumn(where ast.ExprNode, fn func(expr *ast.ColumnNameExpr) bool) bool {
	switch x := where.(type) {
	case nil:
	case *ast.ColumnNameExpr:
		return fn(x)
	case *ast.BinaryOperationExpr:
		if scanWhereStmtColumn(x.R, fn) || scanWhereStmtColumn(x.L, fn) {
			return true
		} else {
			return false
		}
	case *ast.UnaryOperationExpr:
		return scanWhereStmtColumn(x.V, fn)
	// boolean_primary is true|false
	case *ast.IsTruthExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	// boolean_primary is (not) null
	case *ast.IsNullExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	// boolean_primary comparison_operator {ALL | ANY} (subquery)
	case *ast.CompareSubqueryExpr:
		return scanWhereStmtColumn(x.L, fn)
	// boolean_primary IN (expr,...)
	case *ast.PatternInExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	// boolean_primary Between expr and expr
	case *ast.BetweenExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	// boolean_primary (not) like expr
	case *ast.PatternLikeExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	// boolean_primary (not) regexp expr
	case *ast.PatternRegexpExpr:
		return scanWhereStmtColumn(x.Expr, fn)
	case *ast.RowExpr:
		if x.Values != nil {
			ok := false
			for _, expr := range x.Values {
				ok = ok || scanWhereStmtColumn(expr, fn)
				if ok {
					return ok
				}
			}
			return ok
		}
	default:
		return false
	}
	return false
}

func getAlterTableSpecByTp(specs []*ast.AlterTableSpec, ts ...ast.AlterTableType) []*ast.AlterTableSpec {
	s := []*ast.AlterTableSpec{}
	for _, spec := range specs {
		for _, tp := range ts {
			if spec.Tp == tp {
				s = append(s, spec)
			}
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

func hasPrimaryKey(stmt *ast.CreateTableStmt) bool {
	_, hasPk := getPrimaryKey(stmt)
	return hasPk
}

func hasUniqIndex(stmt *ast.CreateTableStmt) bool {
	for _, constraint := range stmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintUniq:
			return true
		}
	}
	return false
}

func ReplaceTableName(node ast.Node) string {
	var schema string
	var table string
	query := node.Text()
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		schema = stmt.Table.Schema.String()
		table = stmt.Table.Name.String()
	case *ast.AlterTableStmt:
		schema = stmt.Table.Schema.String()
		table = stmt.Table.Name.String()
	}
	if schema != "" {
		return replaceTableName(query, schema, table)
	}
	return query
}

func replaceTableName(query, schema, table string) string {
	re := regexp.MustCompile(fmt.Sprintf("%s\\.%s|`%s`\\.`%s`|`%s`\\.%s|%s\\.`%s`",
		schema, table, schema, table, schema, table, schema, table))
	return re.ReplaceAllString(query, fmt.Sprintf("`%s`", table))
}

func getLimitCount(limit *ast.Limit, _default int64) (int64, error) {
	if limit == nil {
		return _default, nil
	}
	return strconv.ParseInt(exprFormat(limit.Count), 0, 64)
}

func getDuplicate(c []string) []string {
	d := []string{}
	for i, v1 := range c {
		for j, v2 := range c {
			if i >= j {
				continue
			}
			if v1 == v2 {
				d = append(d, v1)
			}
		}
	}
	return d
}