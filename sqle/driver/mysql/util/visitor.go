package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

// FingerprintVisitor implements ast.Visitor interface.
type FingerprintVisitor struct{}

func (f *FingerprintVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		v.Type.Charset = ""
		v.SetValue([]byte("?"))
	}
	return n, false
}

func (f *FingerprintVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

type ParamMarkerChecker struct {
	HasParamMarker bool
}

func (p *ParamMarkerChecker) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	if _, ok := in.(*driver.ParamMarkerExpr); ok {
		p.HasParamMarker = true
		return in, true
	}
	return in, false
}

func (p *ParamMarkerChecker) Leave(in ast.Node) (node ast.Node, skipChildren bool) {
	return in, true
}

type HasVarChecker struct {
	HasVar bool
}

func (v *HasVarChecker) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	if _, ok := in.(*ast.VariableExpr); ok {
		v.HasVar = true
		return in, true
	}

	return in, false
}

func (v *HasVarChecker) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// OB建表语句中要剥离的关键词和对应的正则表达式
//
//	OBMySQL  v4.4.0 SQL型
var (
	obOptionsRegex = map[string]*regexp.Regexp{
		// 存储属性  https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000003382263
		"BLOCK_SIZE":  regexp.MustCompile(`(?i)BLOCK_SIZE\s*(=\s*\d+|\s+\d+\s*(LOCAL)?)`), // 兼容 BLOCK_SIZE=16384 和 BLOCK_SIZE 16384 [LOCAL]
		"COMPRESSION": regexp.MustCompile(`(?i)COMPRESSION\s*=\s*'[^']*'`),
		"PCTFREE":     regexp.MustCompile(`(?i)PCTFREE\s*=\s*\d+`),
		"TABLET_SIZE": regexp.MustCompile(`(?i)TABLET_SIZE\s*=\s*\d+`),
		// 分布式属性 https://www.oceanbase.com/docs/common-oceanbase-database-cn-1000000003382263
		"PRIMARY_ZONE":       regexp.MustCompile(`(?i)PRIMARY_ZONE\s*=\s*'[^']*'`),
		"ZONE_LIST":          regexp.MustCompile(`(?i)ZONE_LIST\s*=\s*'[^']*'`),
		"REPLICA_NUM":        regexp.MustCompile(`(?i)REPLICA_NUM\s*=\s*\d+`),
		"TABLEGROUP":         regexp.MustCompile(`(?i)TABLEGROUP\s*=\s*'[^']*'`),
		"DEFAULT_TABLEGROUP": regexp.MustCompile(`(?i)DEFAULT_TABLEGROUP\s*=\s*'[^']*'`),
	}

	// MySQL兼容的table options，需要保留
	mysqlCompatibleOptionsRegex = map[string]*regexp.Regexp{
		"DEFAULT_CHARSET": regexp.MustCompile(`(?i)DEFAULT\s+CHARSET\s*=\s*\w+`),
		"COLLATE":         regexp.MustCompile(`(?i)COLLATE\s*=\s*\w+`),
		"COMMENT":         regexp.MustCompile(`(?i)COMMENT\s*=\s*'[^']*'`),
	}
)

// ParseCreateTableSqlCompatibly 解析并清理OceanBase的建表语句
// 1. 先删除SQL尾部的options和分区等OceanBase特有的内容，只保留到最后一个右括号
// 2. 再做特有语法点的清理
func ParseCreateTableSqlCompatibly(createTableSql string) (*ast.CreateTableStmt, error) {
	// 先尝试完整解析
	if stmt, err := ParseCreateTableStmt(createTableSql); err == nil {
		return stmt, nil
	}

	/*
		建表语句可能如下:
		CREATE TABLE `__all_server_event_history` (

			`gmt_create` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
			`svr_ip` varchar(46) NOT NULL,
			`svr_port` bigint(20) NOT NULL,
			`module` varchar(64) NOT NULL,
			`event` varchar(64) NOT NULL,
			`name1` varchar(256) DEFAULT '',
			`value1` varchar(256) DEFAULT '',
			`name2` varchar(256) DEFAULT '',
			`value2` longtext DEFAULT NULL,
			`name3` varchar(256) DEFAULT '',
			`value3` varchar(256) DEFAULT '',
			`name4` varchar(256) DEFAULT '',
			`value4` varchar(256) DEFAULT '',
			`name5` varchar(256) DEFAULT '',
			`value5` varchar(256) DEFAULT '',
			`name6` varchar(256) DEFAULT '',
			`value6` varchar(256) DEFAULT '',
			`extra_info` varchar(512) DEFAULT '',
			PRIMARY KEY (`gmt_create`, `svr_ip`, `svr_port`)

		) DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'none' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 10 TABLEGROUP = 'oceanbase'

			partition by key_v2(svr_ip, svr_port)

		(partition p0,
		partition p1,
		partition p2,
		partition p3,
		partition p4,
		partition p5,
		partition p6,
		partition p7,
		partition p8,
		partition p9,
		partition p10,
		partition p11,
		partition p12,
		partition p13,
		partition p14,
		partition p15)

		建表语句后半段是options，oceanbase mysql模式下的show create table结果返回的options中包含mysql不支持的options, 为了能解析, 方法将会倒着遍历建表语句, 每次找到右括号时截断后面的部分, 检查截断部分是否包含MySQL兼容选项, 如果有则加回SQL末尾, 此时剩余的建表语句将不在包含OB特有options
	*/

	// 第1步：移除OB特有选项
	replacedSQL := createTableSql
	for _, regex := range obOptionsRegex {
		replacedSQL = regex.ReplaceAllString(replacedSQL, "")
	}

	// 第2步：尝试从后向前截断解析
	cleanedSQL := replacedSQL
	for i := len(replacedSQL) - 1; i >= 0; i-- {
		if replacedSQL[i] == ')' {
			truncatedPart := cleanedSQL[i+1:]
			cleanedSQL = cleanedSQL[:i+1]
			// 从截断部分提取MySQL兼容的选项
			var mysqlCompatibleOptions []string
			for _, regex := range mysqlCompatibleOptionsRegex {
				matches := regex.FindString(truncatedPart)
				if matches == "" {
					continue
				}
				mysqlCompatibleOptions = append(mysqlCompatibleOptions, matches)
			}
			// 如果有MySQL兼容选项，加回SQL末尾
			trySQL := cleanedSQL
			if len(mysqlCompatibleOptions) > 0 {
				trySQL += " " + strings.Join(mysqlCompatibleOptions, " ")
			}

			// 再次尝试解析
			if stmt, err := ParseCreateTableStmt(trySQL); err == nil {
				return stmt, nil
			}
		}
	}

	// 所有尝试都失败
	return nil, fmt.Errorf("failed to parse create table SQL with compatible method")
}

func ParseCreateTableStmt(sql string) (*ast.CreateTableStmt, error) {
	t, err := ParseOneSql(sql)
	if err != nil {
		return nil, err
	}
	createStmt, ok := t.(*ast.CreateTableStmt)
	if !ok {
		return nil, fmt.Errorf("stmt not support")
	}
	return createStmt, nil
}

// CapitalizeProcessor implements ast.Visitor interface.
//
// CapitalizeProcessor capitalize identifiers as needed.
//
// format.RestoreNameUppercase can not control name comparisons accurate.
// CASE:
// Database/Table/Table-alias names are case-insensitive when lower_case_table_names equals 1.
// Some identifiers, such as Tablespace names are case-sensitive which not affected by lower_case_table_names.
// ref: https://dev.mysql.com/doc/refman/5.7/en/identifier-case-sensitivity.html
type CapitalizeProcessor struct {
	capitalizeTableName      bool
	capitalizeTableAliasName bool
	capitalizeDatabaseName   bool
}

func (cp *CapitalizeProcessor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableSource:
		if cp.capitalizeTableAliasName {
			stmt.AsName.O = strings.ToUpper(stmt.AsName.O)
		}
	case *ast.TableName:
		if cp.capitalizeTableName {
			stmt.Name.O = strings.ToUpper(stmt.Name.O)
		}
		if cp.capitalizeDatabaseName {
			stmt.Schema.O = strings.ToUpper(stmt.Schema.O)
		}
	}

	if cp.capitalizeDatabaseName {
		switch stmt := in.(type) {
		case *ast.DropDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		case *ast.CreateDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		case *ast.AlterDatabaseStmt:
			stmt.Name = strings.ToUpper(stmt.Name)
		}
	}
	return in, false
}

func (cp *CapitalizeProcessor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// TableNameExtractor implements ast.Visitor interface.
type TableNameExtractor struct {
	TableNames map[string] /*origin table name without database name*/ *ast.TableName
}

func (te *TableNameExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableName:
		te.TableNames[stmt.Name.O] = stmt
	}
	return in, false
}

func (te *TableNameExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SelectStmtExtractor struct {
	SelectStmts []*ast.SelectStmt
}

func (se *SelectStmtExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		se.SelectStmts = append(se.SelectStmts, stmt)
	}
	return in, false
}

func (se *SelectStmtExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SubQueryMaxNestNumExtractor struct {
	MaxNestNum     *int
	CurrentNestNum int
}

func (se *SubQueryMaxNestNumExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	stmt, ok := in.(*ast.SubqueryExpr)
	if !ok {
		return in, false
	}

	if *se.MaxNestNum < se.CurrentNestNum {
		*se.MaxNestNum = se.CurrentNestNum
	}

	numExtractor := SubQueryMaxNestNumExtractor{MaxNestNum: se.MaxNestNum, CurrentNestNum: se.CurrentNestNum + 1}
	stmt.Query.Accept(&numExtractor)

	return in, true
}

func (se *SubQueryMaxNestNumExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type TableSourceExtractor struct {
	TableSources map[string] /*origin table name and as name without database name*/ *ast.TableSource
}

func (ts *TableSourceExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableSource:
		if stmt.AsName.L != "" {
			ts.TableSources[stmt.AsName.L] = stmt
		}
		if tableName, ok := stmt.Source.(*ast.TableName); ok {
			ts.TableSources[tableName.Name.O] = stmt
		}
	}
	return in, false
}

func (ts *TableSourceExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// SelectFieldExtractor
// 检测select的字段是否只包含count(*)函数
type SelectFieldExtractor struct {
	IsSelectOnlyIncludeCountFunc bool
}

func (se *SelectFieldExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		isOneFiled := len(stmt.Fields.Fields) == 1
		if !isOneFiled {
			return in, true
		}

		if aggregateFuncExpr, ok := stmt.Fields.Fields[0].Expr.(*ast.AggregateFuncExpr); ok {
			isOneArg := len(aggregateFuncExpr.Args) == 1
			if !isOneArg {
				return in, true
			}

			var arg interface{}
			if expr, ok := aggregateFuncExpr.Args[0].(ast.ValueExpr); ok {
				arg = expr.GetValue()
			}

			isDigitOne := arg.(int64) == 1
			isCountFunc := strings.ToLower(aggregateFuncExpr.F) == ast.AggFuncCount
			if isCountFunc && isDigitOne {
				se.IsSelectOnlyIncludeCountFunc = true
				return in, true
			}
		}
	}
	return in, true
}

func (se *SelectFieldExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type SelectVisitor struct {
	SelectList []*ast.SelectStmt
}

func (v *SelectVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		v.SelectList = append(v.SelectList, stmt)
	}
	return in, false
}

func (v *SelectVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type ColumnNameVisitor struct {
	ColumnNameList []*ast.ColumnNameExpr
}

func (v *ColumnNameVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.ColumnNameExpr:
		v.ColumnNameList = append(v.ColumnNameList, stmt)
	}
	return in, false
}

func (v *ColumnNameVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type WhereVisitor struct {
	WhereList         []ast.ExprNode
	WhetherContainNil bool // 是否需要包含空的where，例如select * from t1 该语句的where为空
}

func (v *WhereVisitor) append(where ast.ExprNode) {
	if where != nil {
		v.WhereList = append(v.WhereList, where)
	} else if v.WhetherContainNil {
		v.WhereList = append(v.WhereList, nil)
	}
}

func (v *WhereVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return in, false
		}
		v.append(stmt.Where)
	case *ast.UpdateStmt:
		v.append(stmt.Where)
	case *ast.DeleteStmt:
		v.append(stmt.Where)
	}
	return in, false
}

func (v *WhereVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type EqualColumns struct {
	Left  *ast.ColumnName
	Right *ast.ColumnName
}
type EqualConditionVisitor struct {
	ConditionList []EqualColumns
}

func (v *EqualConditionVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.BinaryOperationExpr:
		var tableNameL, tableNameR string
		var equalColumns EqualColumns
		if stmt.Op == opcode.EQ {
			switch t := stmt.L.(type) {
			case *ast.ColumnNameExpr:
				tableNameL = t.Name.Table.L
				equalColumns.Left = t.Name
			}
			switch t := stmt.R.(type) {
			case *ast.ColumnNameExpr:
				tableNameR = t.Name.Table.L
				equalColumns.Right = t.Name
			}
			if tableNameL != "" && tableNameR != "" && tableNameL != tableNameR {
				v.ConditionList = append(v.ConditionList, equalColumns)
			}
		}
	}
	return in, false
}

func (v *EqualConditionVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type FuncCallExprVisitor struct {
	FuncCallList []*ast.FuncCallExpr
}

func (v *FuncCallExprVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.FuncCallExpr:
		v.FuncCallList = append(v.FuncCallList, stmt)
	}
	return in, false
}

func (v *FuncCallExprVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type WhereWithTable struct {
	WhereStmt *ast.ExprNode
	TableRef  *ast.Join
}

type WhereWithTableVisitor struct {
	WhereStmts []*WhereWithTable
}

func (v *WhereWithTableVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null, skip check. EX: select 1;select version;
			return in, false
		}
		v.WhereStmts = append(v.WhereStmts, &WhereWithTable{WhereStmt: &stmt.Where, TableRef: stmt.From.TableRefs})
	case *ast.DeleteStmt:
		v.WhereStmts = append(v.WhereStmts, &WhereWithTable{WhereStmt: &stmt.Where, TableRef: stmt.TableRefs.TableRefs})
	case *ast.UpdateStmt:
		v.WhereStmts = append(v.WhereStmts, &WhereWithTable{WhereStmt: &stmt.Where, TableRef: stmt.TableRefs.TableRefs})
	}
	return in, false
}

func (v *WhereWithTableVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}
