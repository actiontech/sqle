package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00034 = "SQLE00034"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00034,
			Desc:       plocale.Rule00034Desc,
			Annotation: plocale.Rule00034Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID, plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00034Message,
		Func:    RuleSQLE00034,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00034): "在 MySQL 中，字段约束为NOT NULL时必须带默认值."
您应遵循以下逻辑：
1. 对于"CREATE TABLE..."语句，检查语法树中的列定义节点，如果某个列定义包含NOT NULL约束但没有DEFAULT子节点，报告违反规则。
2. 对于"ALTER TABLE..."语句，检查语法树中的列修改节点，如果某个列修改包含NOT NULL约束但没有DEFAULT子节点，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00034(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			// if the column has "NOT NULL" constraint but no "DEFAULT" constraint
			if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) && !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				// if the column has "NOT NULL" constraint but no "DEFAULT" constraint
				if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) && !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00034)
		return nil
	}

	return nil
}

// ==== Rule code end ====
