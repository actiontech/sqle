package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00031 = "SQLE00031"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00031,
			Desc:       plocale.Rule00031Desc,
			Annotation: plocale.Rule00031Annotation,
			Category:   plocale.RuleTypeUsageSuggestion,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagView.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID, plocale.RuleTagSQLView.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00031Message,
		Func:    RuleSQLE00031,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00031): "在 MySQL 中，禁止使用视图."
您应遵循以下逻辑：
1. 对于 "CREATE ..."语句，如果存在以下任何一项，则报告违反规则：
  1. 语法节点中包含视图定义（如CreateStmt中的View定义）
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00031(input *rulepkg.RuleHandlerInput) error {
	// 解析 SQL 语句
	switch input.Node.(type) {
	// 检查是否为 CREATE VIEW 语句
	case *ast.CreateViewStmt:
		// 对于“CREATE VIEW ...” 语句，直接报告违反规则
		rulepkg.AddResult(input.Res, input.Rule, SQLE00031)
		return nil
	}
	return nil
}

// ==== Rule code end ====
