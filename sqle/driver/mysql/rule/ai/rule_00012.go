package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00012 = "SQLE00012"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00012,
			Desc:       plocale.Rule00012Desc,
			Annotation: plocale.Rule00012Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00012Message,
		Func:    RuleSQLE00012,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00012): "在 MySQL 中，建议使用BIGINT类型表示小数."
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句，执行以下检查：
   1. 检查语法节点中是否定义了 DECIMAL 类型字段（例如用于表示价格、金额、数量等）。
   如果发现 DECIMAL 类型字段，则标记为违反规则。

2. 针对 "ALTER TABLE..." 语句，执行以下检查：
   1. 检查语法节点中是否添加或修改为 DECIMAL 类型字段（例如用于表示价格、金额、数量等）。
   如果发现 DECIMAL 类型字段，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00012(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	default:
		return nil
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00012, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
