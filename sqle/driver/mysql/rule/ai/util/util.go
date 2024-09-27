package util

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/tidb/sessionctx/stmtctx"
	driver "github.com/pingcap/tidb/types/parser_driver"
	parser "github.com/pingcap/tidb/types/parser_driver"
)

// a helper function to join column names
func JoinColumnNames(columns []*ast.ColumnDef) string {
	names := make([]string, len(columns))
	for i := range columns {
		names[i] = columns[i].Name.OrigColName()
	}
	return strings.Join(names, ",")
}

// a helper function to check if a string is in list
func IsStrInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// a helper function to get MySQL alter-table commands by command types
func GetAlterTableCommandsByTypes(alterTableStmt *ast.AlterTableStmt, ts ...ast.AlterTableType) []*ast.AlterTableSpec {
	s := []*ast.AlterTableSpec{}

	if alterTableStmt == nil || alterTableStmt.Specs == nil {
		return nil
	}

	for _, spec := range alterTableStmt.Specs {
		for _, tp := range ts {
			if spec.Tp == tp {
				s = append(s, spec)
			}
		}
	}
	return s
}

// a helper function to get the MySQL table option whose type is the _targetOption_ (Engine/Auto Increment/...)
func GetTableOption(options []*ast.TableOption, targetOption ast.TableOptionType) *ast.TableOption {
	for _, option := range options {
		if option.Tp == targetOption {
			return option
		}
	}
	return nil
}

// a helper function to get the MySQL database option whose type is the _targetOption_ (Charset/Collate/...)
func GetDatabaseOption(options []*ast.DatabaseOption, targetOption ast.DatabaseOptionType) *ast.DatabaseOption {
	for _, option := range options {
		if option.Tp == targetOption {
			return option
		}
	}
	return nil
}

// a helper function to get the MySQL column option whose type is the _targetOption_ (NOT NULL/DEFAULT/...)
func GetColumnOption(columnDef *ast.ColumnDef, targetOption ast.ColumnOptionType) *ast.ColumnOption {
	for _, option := range columnDef.Options {
		if option.Tp == targetOption {
			return option
		}
	}
	return nil
}

// a helper function to check if MySQL column has the target option (MySQL's column option, such as NOT NULL, AUTO_INCREMENT...)
func IsColumnHasOption(columnDef *ast.ColumnDef, targetOption ast.ColumnOptionType) bool {
	for _, option := range columnDef.Options {
		if option.Tp == targetOption {
			return true
		}
	}
	return false
}

// a helper function to get the MySQL table constraint whose type is the _targetConstraint_ (PRIMARY/FOREIGN/...)
func GetTableConstraint(constraints []*ast.Constraint, targetConstraint ast.ConstraintType) *ast.Constraint {
	for _, constraint := range constraints {
		if constraint.Tp == targetConstraint {
			return constraint
		}
	}
	return nil
}

// a helper function to get the MySQL table constraints whose type is the _targetConstraint_ (PRIMARY/FOREIGN/...)
func GetTableConstraints(constraints []*ast.Constraint, targetConstraint ...ast.ConstraintType) []*ast.Constraint {
	c := []*ast.Constraint{}
	for _, constraint := range constraints {
		for _, target := range targetConstraint {
			if constraint.Tp == target {
				c = append(c, constraint)
			}
		}
	}
	return c
}

// a helper function to get the MySQL index constraint types
func GetIndexConstraintTypes() []ast.ConstraintType {
	return []ast.ConstraintType{
		ast.ConstraintIndex,
		ast.ConstraintUniqIndex,
		ast.ConstraintKey,
		ast.ConstraintUniq,
		ast.ConstraintUniqKey,
		ast.ConstraintPrimaryKey,
	}
}

// a helper function to check if MySQL column has specified character set
func IsColumnHasSpecifiedCharset(columnDef *ast.ColumnDef) bool {
	return columnDef.Tp.Charset != ""
}

// a helper function to get MySQL column width
func GetColumnWidth(columnDef *ast.ColumnDef) int {
	return columnDef.Tp.Flen
}

// a helper function to check if alter table is target type
func IsAlterTableCommand(spec *ast.AlterTableSpec, expectedType ast.AlterTableType) bool {
	return spec.Tp == expectedType
}

// a helper function to check if the alter table is altering table option _expectOption_ in MySQL
func IsAlterTableCommandAlterOption(spec *ast.AlterTableSpec, expectOption ast.TableOptionType) bool {
	if spec.Tp != ast.AlterTableOption {
		return false
	}
	for _, option := range spec.Options {
		if option.Tp == expectOption {
			return true
		}
	}
	return false
}

// a helper function to check if the column type is in the given data types
func IsColumnTypeEqual(columnDef *ast.ColumnDef, targetType ...byte) bool {
	if columnDef == nil {
		return false
	}
	for _, dbType := range targetType {
		if columnDef.Tp.Tp == dbType {
			return true
		}
	}
	return false
}

// a helper function to get the MySQL Blob column types
func GetBlobDbTypes() []byte {
	return []byte{
		mysql.TypeBlob,
		mysql.TypeTinyBlob,
		mysql.TypeMediumBlob,
		mysql.TypeLongBlob,
	}
}

// a helper function to check if the column has auto increment constraint in MySQL
func IsColumnAutoIncrement(columnDef *ast.ColumnDef) bool {
	return IsColumnHasOption(columnDef, ast.ColumnOptionAutoIncrement)
}

// a helper function to check if the column is a primary key constraint in MySQL
func IsColumnPrimaryKey(columnDef *ast.ColumnDef) bool {
	return IsColumnHasOption(columnDef, ast.ColumnOptionPrimaryKey)
}

// a helper function to get the column name in MySQL
func GetColumnName(columnDefNode *ast.ColumnDef) string {
	return columnDefNode.Name.Name.O
}

// a helper function to check if the column has Null option in MySQL
func IsOptionValIsNull(opt *ast.ColumnOption) bool {
	return opt.Expr.GetType().Tp == mysql.TypeNull
}

// a helper function to check if the option is target function in MySQL
func IsOptionFuncCall(option *ast.ColumnOption, expectedFuncCall string) bool {
	funcCallExpr, ok := option.Expr.(*ast.FuncCallExpr)
	if !ok {
		return false
	}
	return strings.EqualFold(funcCallExpr.FnName.L, expectedFuncCall)
}

// a helper function to get ValueExpr string
func GetValueExprStr(expr ast.ExprNode) string {
	if stmt, ok := expr.(*parser.ValueExpr); ok {
		return stmt.GetDatumString()
	}
	return ""
}

// a helper function to get the offset value of the limit in Select
func GetLimitOffsetValue(stmt *ast.SelectStmt) int64 {
	if stmt.Limit != nil && stmt.Limit.Offset != nil {
		offsetVal, ok := stmt.Limit.Offset.(*parser.ValueExpr)
		if !ok {
			return -2
		}
		return offsetVal.Datum.GetInt64()
	}
	return -1
}

// a helper function to get the offset value of the limit in the Union query
func GetLimitOffsetValueByUnionStmt(stmt *ast.UnionStmt) int64 {
	if stmt.Limit != nil && stmt.Limit.Offset != nil {
		offsetVal, ok := stmt.Limit.Offset.(*parser.ValueExpr)
		if !ok {
			return -2
		}
		return offsetVal.Datum.GetInt64()
	}
	return -1
}

// a helper function to get index column name
func GetIndexColName(index *ast.IndexPartSpecification) string {
	return index.Column.Name.String()
}

// a helper function to get create table stmt by query
func GetCreateTableStmt(context *session.Context, table *ast.TableName) (*ast.CreateTableStmt, error) {
	stmt, exist, err := context.GetCreateTableStmt(table)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("failed to get create table stmt, table is not exist")
	}

	return stmt, nil

}

// a helper function to get schema name from AST or current schema.
func GetSchemaName(context *session.Context, schemaName string) string {
	if schemaName == "" {
		return context.CurrentSchema()
	}

	return schemaName
}

// a helper function to indexes info for the given table
func GetTableIndexes(context *session.Context, tableName, schemaName string) ([]*executor.TableIndexesInfo, error) {
	schemaName = GetSchemaName(context, schemaName)
	return context.GetExecutor().GetTableIndexesInfo(supplementalQuotationMarks(schemaName), supplementalQuotationMarks(tableName))
}

// a helper function to get index expression for the given table
func GetIndexExpressionsForTables(ctx *session.Context, tables []*ast.TableName) ([]string, error) {
	existIndexExpr := []string{}
	for _, table := range tables {
		indexesInfo, err := GetTableIndexes(ctx, table.Name.String(), table.Schema.String())
		if err != nil {
			return nil, err
		}
		for _, indexInfo := range indexesInfo {
			existIndexExpr = append(existIndexExpr, indexInfo.Expression)
		}
	}

	return existIndexExpr, nil
}

// a helper function to extract all select stmt from a given AST Node
func GetSelectStmt(stmt ast.Node) []*ast.SelectStmt {
	if stmt == nil {
		return nil
	}
	selectStmtExtractor := selectStmtExtractor{}
	stmt.Accept(&selectStmtExtractor)
	return selectStmtExtractor.SelectStmts
}

// a helper function to extract all table name from a given AST Node
func GetTableNames(stmt ast.Node) []*ast.TableName {
	if stmt == nil {
		return nil
	}
	e := tableNameExtractor{}
	stmt.Accept(&e)
	return e.tableNames
}

// a helper function to extract all subquery from a given AST Node
func GetSubquery(stmt ast.Node) []*ast.SubqueryExpr {
	if stmt == nil {
		return nil
	}
	e := SubqueryExprExtractor{}
	stmt.Accept(&e)
	return e.expr
}

// a helper function to get the default table from a given select statement
func GetDefaultTable(stmt *ast.SelectStmt) *ast.TableName {
	if stmt == nil {
		return nil
	}
	if stmt.From == nil {
		return nil
	}
	if stmt.From.TableRefs == nil {
		return nil
	}
	if t, ok := stmt.From.TableRefs.Left.(*ast.TableSource); ok && t != nil {
		if name, ok := t.Source.(*ast.TableName); ok && name != nil {
			return name
		}
	}
	return nil
}

// a helper function to get the table source from join node
func GetTableSourcesFromJoin(join *ast.Join) []*ast.TableSource {
	sources := []*ast.TableSource{}
	if join == nil {
		return sources
	}
	if n := join.Left; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			sources = append(sources, t)
		case *ast.Join:
			sources = append(sources, GetTableSourcesFromJoin(t)...)
		}
	}
	if n := join.Right; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			sources = append(sources, t)
		case *ast.Join:
			sources = append(sources, GetTableSourcesFromJoin(t)...)
		}
	}
	return sources
}

// a helper function to get the table alias info from join node
func GetTableAliasInfoFromJoin(join *ast.Join) []*TableAliasInfo {
	tableAlias := make([]*TableAliasInfo, 0)
	tableSources := GetTableSourcesFromJoin(join)
	for _, tableSource := range tableSources {
		if tableName, ok := tableSource.Source.(*ast.TableName); ok {
			tableAlias = append(tableAlias, &TableAliasInfo{
				TableAliasName: tableSource.AsName.String(),
				TableName:      tableName.Name.L,
				SchemaName:     tableName.Schema.L,
			})
		}
	}
	return tableAlias
}

// a helper function to extract all function name from a given expr Node of a SQL statement
func GetFuncName(expr ast.ExprNode) []string {
	extractor := funcExtractor{}
	expr.Accept(&extractor)

	return extractor.funcNames
}

// a helper function to extract function expressions from a given expr node of a SQL statement
func GetFuncExpr(expr ast.ExprNode) []string {
	extractor := funcExtractor{}
	expr.Accept(&extractor)

	return extractor.expr
}

// a helper function to extract math op expressions from a given expr node of a SQL statement
func GetMathOpExpr(expr ast.ExprNode) []string {
	extractor := mathOpExtractor{}
	expr.Accept(&extractor)

	return extractor.expr
}

// a helper function to get the column expr from a given expr node of a SQL statement
func GetColumnNameInExpr(expr ast.ExprNode) []*ast.ColumnNameExpr {
	if expr == nil {
		return nil
	}
	extractor := columnNameExprExtractor{}
	expr.Accept(&extractor)
	return extractor.columnExpr
}

// a helper function to converts an AST (Abstract Syntax Tree) expression node into its string representation.
func ExprFormat(node ast.ExprNode) string {
	switch node.(type) {
	case *ast.DefaultExpr:
		return "DEFAULT"
	default:
		writer := bytes.NewBufferString("")
		node.Format(writer)
		return writer.String()
	}
}

// a helper function to calculate index discrimination in MySQL
func CalculateIndexDiscrimination(context *session.Context, table *ast.TableName, colNames []string) (map[string]float64, error) {
	return context.GetSelectivityOfColumnsV2(table, colNames)
}

// a helper function to get the execution plan of a SQL statement in MySQL
func GetExecutionPlan(context *session.Context, sql string) (*executor.ExplainWithWarningsResult, error) {
	return context.GetExecutionPlanWithWarnings(sql)
}

// a helper function to get table row count in MySQL
func GetTableRowCount(context *session.Context, table *ast.TableName) (int, error) {
	return context.GetTableRowCount(table)
}

// a helper function to get the MySQL table size of a table
func GetTableSizeMB(context *session.Context, tableName string) (int64, error) {
	size, err := context.GetTableSize(&ast.TableName{Name: model.NewCIStr(tableName)})
	if err != nil {
		return 0, err
	}
	return int64(math.Floor(size)), nil
}

// a helper function to get if the table is a temporary table in MySQL
func IsTemporaryTable(context *session.Context, tableName *ast.TableName) (bool, error) {
	t, exist := context.GetTableInfo(tableName)
	if !exist {
		return false, fmt.Errorf("table %s not exist", tableName.Name.String())
	}
	return t.OriginalTable.IsTemporary, nil
}

// a helper function to check if the where clause is a constant true
func IsExprConstTrue(where ast.ExprNode) bool {
	notAlwaysTrue := false
	ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.FuncCallExpr:
			notAlwaysTrue = true
			return true
		case *ast.ColumnNameExpr:
			notAlwaysTrue = true
			return true
		case *ast.ExistsSubqueryExpr:
			notAlwaysTrue = true
			return true
		case *ast.BinaryOperationExpr:
			compareResult, err := getBinaryExprCompareResult(x)
			if err == nil && !compareResult {
				notAlwaysTrue = true
				return true
			}
			col1, ok := x.R.(*ast.ColumnNameExpr)
			if !ok {
				return false
			}
			col2, ok := x.L.(*ast.ColumnNameExpr)
			if !ok {
				return false
			}
			if col1.Name.String() == col2.Name.String() {
				return true
			}
		}
		return false
	}, where)
	return !notAlwaysTrue
}

// a helper function to get the where clause from a DML statement
func GetWhereExprFromDMLStmt(node ast.Node) (whereList []ast.ExprNode) {
	switch stmt := node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		for _, selectStmt := range GetSelectStmt(stmt) {
			// "select..."
			if selectStmt.Where != nil {
				whereList = append(whereList, selectStmt.Where)
			}
		}
	}

	switch stmt := node.(type) {
	case *ast.DeleteStmt:
		// "delete..."
		if stmt.Where != nil {
			whereList = append(whereList, stmt.Where)
		}
	case *ast.UpdateStmt:
		// "update..."
		if stmt.Where != nil {
			whereList = append(whereList, stmt.Where)
		}
	}
	return whereList
}

// a helper function to scan the where clause of a SQL statement and apply the given function to each expression node
func ScanWhereStmt(fn func(expr ast.ExprNode) (skip bool), exprs ...ast.ExprNode) {
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		// skip all children node
		if fn(expr) {
			continue
		}
		switch x := expr.(type) {
		case *ast.ColumnNameExpr:
		case *ast.SubqueryExpr:
		case *ast.BinaryOperationExpr:
			ScanWhereStmt(fn, x.L, x.R)
		case *ast.UnaryOperationExpr:
			ScanWhereStmt(fn, x.V)
			// boolean_primary is true|false
		case *ast.IsTruthExpr:
			ScanWhereStmt(fn, x.Expr)
			// boolean_primary is (not) null
		case *ast.IsNullExpr:
			ScanWhereStmt(fn, x.Expr)
			// boolean_primary comparison_operator {ALL | ANY} (subquery)
		case *ast.CompareSubqueryExpr:
			ScanWhereStmt(fn, x.L, x.R)
		case *ast.ExistsSubqueryExpr:
			ScanWhereStmt(fn, x.Sel)
			// boolean_primary IN (expr,...)
		case *ast.PatternInExpr:
			es := []ast.ExprNode{}
			es = append(es, x.Expr)
			es = append(es, x.Sel)
			es = append(es, x.List...)
			ScanWhereStmt(fn, es...)
			// boolean_primary Between expr and expr
		case *ast.BetweenExpr:
			ScanWhereStmt(fn, x.Expr, x.Left, x.Right)
			// boolean_primary (not) like expr
		case *ast.PatternLikeExpr:
			ScanWhereStmt(fn, x.Expr, x.Pattern)
			// boolean_primary (not) regexp expr
		case *ast.PatternRegexpExpr:
			ScanWhereStmt(fn, x.Expr, x.Pattern)
		case *ast.RowExpr:
			ScanWhereStmt(fn, x.Values...)
		case *ast.ParenthesesExpr:
			ScanWhereStmt(fn, x.Expr)
		}
	}
}

func utilGetTableName(tableName *ast.TableName) string {
	// 假设这是一个简单的实现
	return tableName.Name.String()
}

func utilGetTargetObjectType(cmd *ast.AlterTableSpec) string {
	// 假设这是一个简单的实现
	return "TEMPORARY_TABLE"
}

// end helper function file. this line which used for ai scanner should be at the end of the file, please do not delete it

// If there are no quotation marks (', ", `) at the beginning and end of the string, the string will be wrapped with "`"
// Need to be wary of the presence of "`" in the string
// do nothing if s is an empty string
func supplementalQuotationMarks(s string) string {
	if s == "" {
		return ""
	}
	end := len(s) - 1
	if s[0] != s[end] {
		return fmt.Sprintf("`%s`", s)
	}
	if string(s[0]) != "'" && s[0] != '"' && s[0] != '`' {
		return fmt.Sprintf("`%s`", s)
	}
	return s
}

// compare binary.L to binary.R
func getBinaryExprCompareResult(binary *ast.BinaryOperationExpr) (bool, error) {
	col1, ok := binary.L.(*driver.ValueExpr)
	if !ok {
		return false, errors.New("binary.L is not driver.ValueExpr")
	}
	col2, ok := binary.R.(*driver.ValueExpr)
	if !ok {
		return false, errors.New("binary.R is not driver.ValueExpr")
	}
	// 暂时只判断相同类型数据的比值，不考虑隐式转换
	if col1.Datum.Kind() != col2.Datum.Kind() {
		return false, nil
	}
	sc := &stmtctx.StatementContext{}

	// col1 < col2; return -1
	// col1 == col2; return 0
	// col1 > col2; return 1
	result, err := col1.CompareDatum(sc, &col2.Datum)
	if err != nil {
		return false, err
	}
	switch binary.Op {
	case opcode.GE:
		if result == 1 || result == 0 {
			return true, nil
		}
	case opcode.GT:
		if result == 1 {
			return true, nil
		}
	case opcode.LE:
		if result == 0 || result == -1 {
			return true, nil
		}
	case opcode.LT:
		if result == -1 {
			return true, nil
		}
	case opcode.EQ:
		if result == 0 {
			return true, nil
		}
	case opcode.NE:
		if result != 0 {
			return true, nil
		}

	}

	return false, nil
}

type TableAliasInfo struct {
	TableName      string
	SchemaName     string
	TableAliasName string
}
