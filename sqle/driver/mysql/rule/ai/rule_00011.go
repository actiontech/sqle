package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00011 = "SQLE00011"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00011,
			Desc:       "在 MySQL 中, 存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
			Annotation: "避免对同一个表使用多条单独的ALTER语句，以减少数据库的锁定时间和执行开销，提高SQL语句的可读性和维护性。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
		AllowOffline: true,
		Func:         RuleSQLE00011,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00011): "在 MySQL 中，存在多条对同一个表的修改语句，建议合并成一个ALTER语句."
您应遵循以下逻辑：
1. 在一批SQL语句中，识别所有ALTER TABLE语句，以及提取其目标表的名称。
2. 如果发现多个ALTER TABLE语句针对同一个表，则报告违反规则，建议将这些语句合并为一个ALTER TABLE语句。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00011(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := input.Ctx.GetTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables)+1 > 1 {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00011)
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
