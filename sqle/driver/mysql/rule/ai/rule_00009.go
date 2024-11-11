package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00009 = "SQLE00009"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00009,
			Desc:       "在 MySQL 中, 避免对条件字段使用函数操作",
			Annotation: "对条件字段做函数操作，可能会破坏索引值的有序性，导致优化器选择放弃走索引，使查询性能大幅度降低",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免对条件字段使用函数操作",
		AllowOffline: false,
		Func:         RuleSQLE00009,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00009): "在 MySQL 中，避免对条件字段使用函数操作."
您应遵循以下逻辑：
1. 对于"SELECT..."语句，检查SQL语句是否包含 WHERE 子句。如果 WHERE 条件中的字段（如表中的列）被函数操作，并且没有使用相应的函数索引，则报告违反规则。函数索引信息需使用辅助函数GetIndexExpressionsForTables 获取。
2. 对于"INSERT...SELECT..."语句，检查其中的SELECT子句，执行与"SELECT..."语句相同的检查。
3. 对于"UPDATE..."语句，检查WHERE子句，执行与"SELECT..."语句相同的检查。
4. 对于"DELETE..."语句，检查WHERE子句，执行与"SELECT..."语句相同的检查。
5. 对于"UNION..."语句，检查其中每个SELECT子句，执行与"SELECT..."语句相同的检查。
6. 对于"WITH..."语句，递归检查其中的每个SELECT子句，执行与"SELECT..."语句相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00009(input *rulepkg.RuleHandlerInput) error {

	checkWhere := func(stmt ast.Node, whereClause ast.ExprNode) bool {
		if whereClause == nil {
			return false
		}

		// Check for function calls in the WHERE clause.
		funcExpr := util.GetFuncExpr(whereClause)
		// If there are no function expressions, then there's no issue.
		if len(funcExpr) == 0 {
			return false
		}

		// Retrieve the names of all tables involved in the statement.
		tables := util.GetTableNames(stmt)
		// Fetch the expressions for all indexes associated with these tables.
		indexExprs, err := util.GetIndexExpressionsForTables(input.Ctx, tables)
		if err != nil {
			// Log an error if fetching index expressions fails.
			log.NewEntry().Errorf("get table index failed, sqle: %v, error: %v", input.Node.Text(), err)
			return false
		}

		violate := false
		// Check if each function expression has a corresponding index.
		for _, e := range funcExpr {
			if !util.IsStrInSlice(e, indexExprs) {
				// Return true if any expression is not indexed, indicating a potential issue.
				violate = true
				break
			}
		}
		return violate
	}

	// 所有select语句
	selectList := util.GetSelectStmt(input.Node)
	for _, sel := range selectList {
		if checkWhere(sel, sel.Where) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00009)
			return nil
		}
	}
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		if checkWhere(stmt, stmt.Where) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00009)
			return nil
		}
	case *ast.DeleteStmt:
		if checkWhere(stmt, stmt.Where) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00009)
			return nil
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
