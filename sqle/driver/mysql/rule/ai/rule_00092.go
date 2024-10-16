package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00092 = "SQLE00092"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00092,
			Desc:       "在 MySQL 中, 建议DELETE/UPDATE语句使用LIMIT子句控制影响行数",
			Annotation: "在进行DELETE和UPDATE操作时，通过添加LIMIT子句可以明确限制操作影响的数据行数。这样做有助于减少由于执行错误而导致的数据损失风险，并可以有效地控制长事务的执行时间，降低对数据库性能的影响。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议DELETE/UPDATE语句使用LIMIT子句控制影响行数",
		AllowOffline: true,
		Func:         RuleSQLE00092,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00092): "在 MySQL 中，建议DELETE/UPDATE语句使用LIMIT子句控制影响行数."
您应遵循以下逻辑：
1. 对于"DELETE..."语句，检查以下条件，如果有任意一个条件不满足，则报告违反规则：
    1. 语法树中应该包含 LIMIT 节点。
2. 对于"UPDATE..."语句，进行与上述相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00092(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		// 检查 DELETE 语句
		if stmt.Limit == nil {
			// 如果没有 LIMIT 节点，报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00092)
			return nil
		}
	case *ast.UpdateStmt:
		// 检查 UPDATE 语句
		if stmt.Limit == nil {
			// 如果没有 LIMIT 节点，报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00092)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
