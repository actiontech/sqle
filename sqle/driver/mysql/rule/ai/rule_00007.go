package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00007 = "SQLE00007"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00007,
			Desc:         plocale.Rule00007Desc,
			Annotation:   plocale.Rule00007Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00007Message,
		Func:    RuleSQLE00007,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
