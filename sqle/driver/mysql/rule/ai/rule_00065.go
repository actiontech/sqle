package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00065 = "SQLE00065"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00065,
			Desc:         plocale.Rule00065Desc,
			Annotation:   plocale.Rule00065Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00065Message,
		Func:    RuleSQLE00065,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00065): "在 MySQL 中，禁止修改表时指定或调整字段在表结构中的顺序."
您应遵循以下逻辑：
1. 对于 "ALTER TABLE...MODIFY ..."语句，如果存在以下任何一项，则报告违反规则：
  1. 检查是否有语法节点 AT_AddColumn 且包含 FIRST
  2. 检查是否有语法节点 AT_AddColumn 且包含 AFTER
2. 对于语句 "ALTER TABLE ... CHANGE ..."，执行与上述相同的检查步骤。
3. 对于语句 "ALTER TABLE ... ADD ..."，执行与上述相同的检查步骤。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00065(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn)

		for _, spec := range specs {
			if spec.Position == nil {
				continue
			}
			if spec.Position.Tp == ast.ColumnPositionFirst || spec.Position.Tp == ast.ColumnPositionAfter {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00065)
			}
		}
	}
	return nil
}

// 规则函数实现结束

// ==== Rule code end ====
