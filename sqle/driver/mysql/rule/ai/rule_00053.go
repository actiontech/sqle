package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

const (
	SQLE00053 = "SQLE00053"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00053,
			Desc:       "不建议使用SELECT *",
			Annotation: "当表结构变更时，使用*通配符选择所有列将导致查询行为会发生更改，与业务期望不符；同时SELECT * 中的无用字段会带来不必要的磁盘I/O，以及网络开销，且无法覆盖索引进而回表，大幅度降低查询效率。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "不建议使用SELECT *",
		AllowOffline: true,
		Func:         RuleSQLE00053,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
