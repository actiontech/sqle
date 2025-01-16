package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00124 = "SQLE00124"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00124,
			Desc:       "删除全表时建议使用 TRUNCATE 替代 DELETE",
			Annotation: "TRUNCATE TABLE 比 DELETE 速度快，且使用的系统和事务日志资源少，同时TRUNCATE后表所占用的空间会被释放，而DELETE后需要手工执行OPTIMIZE才能释放表空间",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message:      "删除全表时建议使用 TRUNCATE 替代 DELETE",
		AllowOffline: true,
		Func:         RuleSQLE00124,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00124): "For dml, using DELETE to delete an entire table is prohibited".
You should follow the following logic:
1. For "DELETE... Statement, checks if there is no WHERE condition or where condition is always True (for example, where 1=1 or where True), reports a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00124(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		// "delete"
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if stmt.Where == nil || util.IsExprConstTrue(input.Ctx, stmt.Where, aliasInfos) {
			// "delete...where..."
			rulepkg.AddResult(input.Res, input.Rule, SQLE00124)
			return nil
		}
	}
	return nil

}

// ==== Rule code end ====
