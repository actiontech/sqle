package util

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
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
func GetSchemaName(context *session.Context, table *ast.TableName) string {
	return context.GetSchemaName(table)
}

// a helper function to indexes info for the given table
func GetTableIndexes(context *session.Context, table *ast.TableName) ([]*executor.TableIndexesInfo, error) {
	schemaName := GetSchemaName(context, table)
	tableName := table.Name.String()
	return context.GetExecutor().GetTableIndexesInfo(supplementalQuotationMarks(schemaName), supplementalQuotationMarks(tableName))
}

// a helper function to get index expression for the given table
func GetIndexExpressionsForTables(ctx *session.Context, tables []*ast.TableName) ([]string, error) {
	existIndexExpr := []string{}
	for _, table := range tables {
		indexesInfo, err := GetTableIndexes(ctx, table)
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

// a helper function to extract function expressions from a given WHERE clause of a SQL statement
func GetFuncExpr(whereClause ast.ExprNode) []string {
	extractor := funcExtractor{columnList: make([]*ast.ColumnName, 0)}
	whereClause.Accept(&extractor)

	return extractor.expr
}

// a helper function to extract math op expressions from a given WHERE clause of a SQL statement
func GetMathOpExpr(whereClause ast.ExprNode) []string {
	extractor := mathOpExtractor{columnList: make([]*ast.ColumnName, 0)}
	whereClause.Accept(&extractor)

	return extractor.expr
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

type funcExtractor struct {
	columnList []*ast.ColumnName
	expr       []string
}

func (fe *funcExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.FuncCallExpr:
		for _, columnNameExpr := range n.Args {
			col, ok := columnNameExpr.(*ast.ColumnNameExpr)
			if !ok {
				continue
			}
			fe.columnList = append(fe.columnList, col.Name)
		}

		fe.expr = append(fe.expr, ExprFormat(n))
		return in, true
	}
	return in, false
}

func (fe *funcExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

type mathOpExtractor struct {
	columnList []*ast.ColumnName
	expr       []string
}

func (me *mathOpExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.BinaryOperationExpr:
		// https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html
		if !isMathComputation(n) {
			return n, false
		}

		if col, ok := n.L.(*ast.ColumnNameExpr); ok {
			me.columnList = append(me.columnList, col.Name)
		}

		if col, ok := n.R.(*ast.ColumnNameExpr); ok {
			me.columnList = append(me.columnList, col.Name)
		}
		me.expr = append(me.expr, ExprFormat(n))
		return in, true

	case *ast.UnaryOperationExpr:
		if n.Op == opcode.Minus {
			col, ok := n.V.(*ast.ColumnNameExpr)
			if !ok {
				return n, false
			}
			me.columnList = append(me.columnList, col.Name)
		}
		me.expr = append(me.expr, ExprFormat(n))
		return in, true
	}

	return in, false
}

func (me *mathOpExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

func isMathComputation(stmt *ast.BinaryOperationExpr) bool {
	return stmt.Op == opcode.Plus || stmt.Op == opcode.Minus || stmt.Op == opcode.Mul || stmt.Op == opcode.Div || stmt.Op == opcode.IntDiv || stmt.Op == opcode.Mod
}

type selectStmtExtractor struct {
	SelectStmts []*ast.SelectStmt
}

func (se *selectStmtExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		se.SelectStmts = append(se.SelectStmts, stmt)
	}
	return in, false
}

func (se *selectStmtExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

// tableNameExtractor implements ast.Visitor interface.
type tableNameExtractor struct {
	tableNames []*ast.TableName
}

func (te *tableNameExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableName:
		te.tableNames = append(te.tableNames, stmt)
	}
	return in, false
}

func (te *tableNameExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}
