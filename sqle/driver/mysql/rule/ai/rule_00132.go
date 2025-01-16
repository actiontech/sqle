package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00132 = "SQLE00132"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00132,
			Desc:       "不推荐使用子查询",
			Annotation: "有些情况下，子查询并不能使用到索引，同时对于返回结果集比较大的子查询，会产生大量的临时表，消耗过多的CPU和IO资源，产生大量的慢查询",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "不推荐使用子查询.",
		AllowOffline: true,
		Func:    RuleSQLE00132,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00132): "For dml, Using subqueries are prohibited".
You should follow the following logic:
1. For "select..." The statement, checks if a SELECT subquery exists in the sentence, and if so, reports a rule violation
2. For "union..." Statement, perform the same checking process as above
3. For "update..." Statement, perform the same checking process as above
4. For "insert..." Statement, perform the same checking process as above
5. For "delete..." Statement, perform the same checking process as above
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00132(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		if len(util.GetSubquery(stmt)) > 0 {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00132)
		}
		return nil
	}
	return nil
}

// ==== Rule code end ====
