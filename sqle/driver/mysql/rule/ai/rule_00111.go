package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00111 = "SQLE00111"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00111,
			Desc:       "避免对条件字段使用表达式操作",
			Annotation: "对条件字段做表达式操作，可能会破坏索引值的有序性，导致优化器选择放弃走索引，使查询性能大幅度降低。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "避免对条件字段使用表达式操作",
		AllowOffline: false,
		Func:    RuleSQLE00111,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00111): "For dml, Mathematical operations and the use of functions on indexed columns is prohibited".
You should follow the following logic:
1. For "select..." statement, checks if the field in the WHERE condition is an function call or math operation, if so, checks if the field has a corresponding function index, if not, report a violation.
2. For "update..." statement, perform the same check as above.
3. For "delete..." statement, perform the same check as above.
4. For "union ..." statement, perform the same check as mentioned above for each SELECT statement within the UNION.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00111(input *rulepkg.RuleHandlerInput) error {

	// Define an inner function to check if expressions in the WHERE clause have corresponding indexes.
	checkExprAndIndex := func(stmt ast.Node, whereClause ast.ExprNode, input *rulepkg.RuleHandlerInput) bool {
		// Return false immediately if there's no WHERE clause to check.
		if whereClause == nil {
			return false
		}
		// Check for mathematical operations in the WHERE clause. If any are found, return true to indicate a potential issue.
		mathExpr := util.GetMathOpExpr(whereClause)
		if len(mathExpr) > 0 {
			return true
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

		// Check if each function expression has a corresponding index.
		for _, e := range funcExpr {
			if !util.IsStrInSlice(e, indexExprs) {
				// Return true if any expression is not indexed, indicating a potential issue.
				return true
			}
		}
		return false
	}

	// Handle different types of SQL statements.
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt:
		// Iterate through all SELECT statements (including those within UNION statements).
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each SELECT statement's WHERE clause.
			if checkExprAndIndex(stmt, selectStmt.Where, input) {
				// Add a rule violation result if an issue is found.
				rulepkg.AddResult(input.Res, input.Rule, SQLE00111)
			}
		}

	case *ast.UpdateStmt:
		// For UPDATE statements, apply the check directly to the WHERE clause.
		if checkExprAndIndex(stmt, stmt.Where, input) {
			// Add a rule violation result if an issue is found.
			rulepkg.AddResult(input.Res, input.Rule, SQLE00111)
		}
	case *ast.DeleteStmt:
		// For DELETE statements, apply the check directly to the WHERE clause.
		if checkExprAndIndex(stmt, stmt.Where, input) {
			// Add a rule violation result if an issue is found.
			rulepkg.AddResult(input.Res, input.Rule, SQLE00111)
		}
	default:
		// No action needed for other types of statements.
		return nil
	}
	return nil
}

// ==== Rule code end ====
