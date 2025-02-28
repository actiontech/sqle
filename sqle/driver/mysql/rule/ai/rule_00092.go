package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00092 = "SQLE00092"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00092,
			Desc:       plocale.Rule00092Desc,
			Annotation: plocale.Rule00092Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagSecurity.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00092Message,
		Func:    RuleSQLE00092,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00092): "在 MySQL 中，建议DELETE/UPDATE语句使用LIMIT子句控制影响行数."
您应遵循以下逻辑：
1. 对于"DELETE..."语句，检查以下条件，如果有任意一个条件不满足，则报告违反规则：
    1. 语法树中应该包含 LIMIT 节点。
2. 对于"UPDATE..."语句，进行与上述相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00092(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		// 检查 DELETE 语句
		if stmt.Limit == nil {
			// 如果没有 LIMIT 节点，报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00092)
			return nil
		}
	case *ast.UpdateStmt:
		// 检查 UPDATE 语句
		if stmt.Limit == nil {
			// 如果没有 LIMIT 节点，报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00092)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
