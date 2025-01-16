package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00007 = "SQLE00007"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00007,
			Desc:       "建表时，自增字段只能设置一个",
			Annotation: "多个自增字段会造成表写入性能影响、可读性差、数据库设计不规范等缺点。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "建表时，自增字段只能设置一个",
		AllowOffline: true,
		Func:         RuleSQLE00007,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00007): "在 MySQL 中，建表时，自增字段只能设置一个."
您应遵循以下逻辑：
1. 针对每个 "CREATE TABLE..." 语句，执行以下检查：
   1. 使用辅助函数 IsColumnAutoIncrement 统计包含 auto_increment 属性的字段数量，当数量大于 1，则报告该表违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00007(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		found := 0

		for _, col := range stmt.Cols {
			if util.IsColumnAutoIncrement(col) {
				found++
			}
		}

		if found > 1 {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00007)
			return nil
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
