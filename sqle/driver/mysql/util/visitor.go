package util

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
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

type TableSourceExtractor struct {
	TableSource *ast.TableSource
}

func (ts *TableSourceExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.TableSource:
		ts.TableSource = stmt
	}
	return in, false
}

func (ts *TableSourceExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}
