package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00067 = "SQLE00067"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00067,
			Desc:       plocale.Rule00067Desc,
			Annotation: plocale.Rule00067Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00067Message,
		Func:    RuleSQLE00067,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00067): "在 MySQL 中，表不建议使用外键."
您应遵循以下逻辑：
1. 检查CREATE TABLE 语句或者 ALTER TABLE 语句的语法节点，查看有无外键定义，如果存在外键定义，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00067(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE TABLE 语句中的外键约束
		if constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintForeignKey); constraint != nil {
			// 存在外键约束，报告违规
			rulepkg.AddResult(input.Res, input.Rule, SQLE00067)
			return nil
		}
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中的外键约束
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			if constraint := util.GetTableConstraint([]*ast.Constraint{spec.Constraint}, ast.ConstraintForeignKey); constraint != nil {
				// 存在外键约束，报告违规
				rulepkg.AddResult(input.Res, input.Rule, SQLE00067)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
