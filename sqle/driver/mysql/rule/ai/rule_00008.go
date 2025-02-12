package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00008 = "SQLE00008"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00008,
			Desc:       plocale.Rule00008Desc,
			Annotation: plocale.Rule00008Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID, plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00008Message,
		Func:    RuleSQLE00008,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00008): "在 MySQL 中，表里必须存在主键."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句：
   - 使用辅助函数 GetCreateTableStmt 检查表定义中是否包含列级别或表级别的主键定义（PRIMARY KEY 关键字）。
   - 如果未包含主键定义，则报告违反规则。

2. 对于 "ALTER TABLE..." 语句：
     1. 执行了删除主键操作（DROP PRIMARY KEY）。
     2. 未同时添加新的主键定义。
   - 如果以上两种情况同时存在，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00008(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}

		found := false

		// check primary key in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnPrimaryKey(col) {
				found = true
				break
			}
		}

		// check primary key in table constraint
		constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintPrimaryKey)
		if nil != constraint {
			//this is a table primary key definition
			found = true
		}

		if !found {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00008)
			return nil
		}
	case *ast.AlterTableStmt:
		dropPrimary := len(util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableDropPrimaryKey)) > 0
		if dropPrimary {
			hasAddPrimary := false
			for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
				if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
					hasAddPrimary = true
				}
			}
			if !hasAddPrimary {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00008)
				return nil
			}
		}

	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
