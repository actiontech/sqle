package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00060 = "SQLE00060"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00060,
			Desc:       plocale.Rule00060Desc,
			Annotation: plocale.Rule00060Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00060Message,
		Func:    RuleSQLE00060,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00060): "在 MySQL 中，表建议添加注释."
您应遵循以下逻辑：
1、检查CREATE TABLE语法节点末尾是否都包含注释节点，否则，将该完整SQL语句加入到触发规则的SQL表表中。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00060(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// Check if the table has a comment at the end
		hasComment := false
		for _, opt := range stmt.Options {
			if opt.Tp == ast.TableOptionComment {
				hasComment = true
				break
			}
		}

		// If no comment is found, add the SQL to the result
		if !hasComment {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00060)
		}
	}

	return nil
}

// ==== Rule code end ====
