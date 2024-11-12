package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00140 = "SQLE00140"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00140,
			Desc:       "建议对表、视图等对象进行操作时指定库名",
			Annotation: "对表、视图等对象进行创建、修改、查询、更新、删除等DDL、DML操作时，如未指定schema或者库名，会导致在不确定的数据库下执行，与实际业务预期不符合，而且会导致SQL语句执行错误。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "建议对表、视图等对象进行操作时指定库名",
		AllowOffline: true,
		Func:         RuleSQLE00140,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00140): "在 MySQL 中，建议对表、视图等对象进行操作时指定库名."
您应遵循以下逻辑：
1. 检查所有SQL语句，涵盖以下类型：
   - DDL语句：CREATE、ALTER、DROP等。
   - DML语句：SELECT、INSERT、UPDATE、DELETE等。
   - 存储过程调用：CALL。

2. 解析每条SQL语句，确认是否明确指定了库名（schema）。

3. 如果任何SQL语句中存在未指定库名的表、视图、存储过程、触发器、函数、事件或索引，标记为违反规则。

4. 对于UNION语句，递归检查所有SELECT子句，确保与独立的SELECT语句一致。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00140(input *rulepkg.RuleHandlerInput) error {
	// Helper function to check if an object has a schema specified
	hasSchema := func(schema string) bool {
		return schema != ""
	}

	// Helper function to check if any object lacks schema
	anyObjectMissingSchema := func(objects []*ast.TableName) bool {
		for _, obj := range objects {
			if !hasSchema(obj.Schema.O) {
				return true
			}
		}
		return false
	}

	// Handle different types of SQL statements
	switch stmt := input.Node.(type) {
	// DML Statements: SELECT, INSERT, UPDATE, DELETE
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		tables := util.GetTableNames(stmt)
		if anyObjectMissingSchema(tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	// UNION Statements
	case *ast.UnionStmt:
		selectStmts := util.GetSelectStmt(stmt)
		for _, selectStmt := range selectStmts {
			tables := util.GetTableNames(selectStmt)
			if anyObjectMissingSchema(tables) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
				return nil
			}
		}

	// DDL Statements: CREATE, ALTER, DROP
	case *ast.CreateTableStmt:
		// Check the table being created
		tableName := stmt.Table.Name
		if tableName.O != "" && (!hasSchema(stmt.Table.Schema.O)) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	case *ast.CreateViewStmt:
		// Check the view being created
		viewName := stmt.ViewName.Name
		if viewName.O != "" && (!hasSchema(stmt.ViewName.Schema.O)) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

		// Check tables referenced in the view definition
		tables := util.GetTableNames(stmt)
		if anyObjectMissingSchema(tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	case *ast.AlterTableStmt:
		tables := util.GetTableNames(stmt)
		if anyObjectMissingSchema(tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	case *ast.DropTableStmt:
		if anyObjectMissingSchema(stmt.Tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	case *ast.RenameTableStmt:
		// Check for all old and new names
		for _, tablePair := range stmt.TableToTables {
			if !hasSchema(tablePair.OldTable.Schema.O) || !hasSchema(tablePair.NewTable.Schema.O) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
				return nil
			}
		}

	// TODO: Stored Procedure Calls

	// Other DDL Statements: CREATE INDEX, ALTER INDEX, etc.
	case *ast.CreateIndexStmt:
		indexName := stmt.IndexName
		if indexName != "" {
			tableName := stmt.Table
			if tableName != nil && (!hasSchema(tableName.Schema.O)) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
				return nil
			}
		}

		tables := util.GetTableNames(stmt)
		if anyObjectMissingSchema(tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}

	// Handle other statement types as needed
	default:
		tables := util.GetTableNames(stmt)
		if anyObjectMissingSchema(tables) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00140)
			return nil
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
