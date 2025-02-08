package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00053 = "SQLE00053"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00053,
			Desc:       plocale.Rule00053Desc,
			Annotation: plocale.Rule00053Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00053Message,
		Func:    RuleSQLE00053,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00053): "在 MySQL 中，不建议使用SELECT *."
您应遵循以下逻辑：
1. 针对所有 DML 和 DQL 语句，递归检查所有 SELECT 子句：
   1. 使用辅助函数GetSelectStmt获取SELECT子句。
   2. 如果 SELECT 子句中包含单独的 * 符号（表示选择所有列），则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00053(input *rulepkg.RuleHandlerInput) error {
	// 获取所有 SELECT 子句，包括嵌套的子查询
	selectStmts := util.GetSelectStmt(input.Node)

	// 检查是否成功获取 SELECT 子句
	if len(selectStmts) == 0 {
		return nil // 如果没有 SELECT 子句，则不违反规则
	}

	// 遍历每个 SELECT 子句
	for _, selectStmt := range selectStmts {
		// 检查 SELECT 子句是否为空
		if selectStmt.Fields == nil || len(selectStmt.Fields.Fields) == 0 {
			continue // 如果 SELECT 子句为空，则继续下一个
		}

		// 遍历 SELECT 列表中的每一项，包含 '*'，则违反规则
		for _, field := range selectStmt.Fields.Fields {
			if field.WildCard != nil {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00053)
			}

		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
