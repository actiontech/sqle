package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
)

const (
	SQLE00143 = "SQLE00143"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00143,
			Desc:       "对于MySQL的DML, 多表关联时，不建议在WHERE条件中对不同表的字段使用OR条件",
			Annotation: "多表关联时，在WHERE条件中对不同表的字段使用OR条件可能会导致SQL无法使用正确的索引",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeIndexInvalidation,
		},
		Message: "对于MySQL的DML, 多表关联时，不建议在WHERE条件中对不同表的字段使用OR条件.",
		AllowOffline: false,
		Func:    RuleSQLE00143,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, ensure compliance with the rule (SQLE00143), which states: "For DML operations involving multiple table joins, using an OR condition in the WHERE clause for fields from different tables is discouraged."
You should follow the following logic:
1. For "SELECT..." Statement, Identify if the statement involves joining multiple tables, If yes, then:
  1. Create a collection to store elements of the WHERE and JOIN conditions.
  2. Include all WHERE conditions in the collection.
  3. Add JOIN conditions to the collection.
  4. Extract and analyze the logical operators used within the collection.
  5. If any operator includes an 'OR' connecting fields from different tables, flag a rule violation.
2. For "INSERT..." Statement, and perform the same checks as above for each SELECT clause in the statement.
3. For "DELETE..." Statement, and perform the same checks as above for each SELECT clause in the statement.
4. For "UPDATE..." Statement, and perform the same checks as above for each SELECT clause in the statement.
5. For "UNION..." Statement, and perform the same checks as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00143(input *rulepkg.RuleHandlerInput) error {
	// Create a collection to store elements of the WHERE and JOIN conditions
	var collection []ast.ExprNode

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UnionStmt, *ast.UpdateStmt:
		// "SELECT...", "INSERT...", "DELETE...", "UNION...", "UPDATE..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				if len(util.GetTableSourcesFromJoin(selectStmt.From.TableRefs)) <= 1 {
					// The statement not involves joining multiple tables
					continue
				}
			}

			// Gather WHERE conditions
			if selectStmt.Where != nil {
				collection = append(collection, selectStmt.Where)
			}

			// Gather JOIN conditions
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil && selectStmt.From.TableRefs.On != nil {
				collection = append(collection, selectStmt.From.TableRefs.On.Expr)
			}
		}
	}

	// Check if there are any OR conditions involving columns from different tables
	for _, node := range collection {
		extractedExpr, ok := node.(*ast.BinaryOperationExpr)
		if !ok {
			continue
		}

		// Check if the expression is an OR
		if extractedExpr.Op != opcode.LogicOr {
			continue
		}
		// Check if conditions involving columns from different tables
		for _, node := range collection {
			if binaryExpr, ok := node.(*ast.BinaryOperationExpr); ok && binaryExpr.Op == opcode.LogicOr {
				leftColNames := util.GetColumnNameInExpr(binaryExpr.L)
				rightColNames := util.GetColumnNameInExpr(binaryExpr.R)
				t := make(map[string] /*table name*/ int)
				for _, col := range append(leftColNames, rightColNames...) {
					if _, ok := t[col.Name.Table.L]; ok {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00143)
						return nil
					}
					t[col.Name.Table.L] = 1
				}
			}
		}

	}
	return nil
}

// ==== Rule code end ====