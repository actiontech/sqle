package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00031 = "SQLE00031"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00031,
			Desc:       "在 MySQL 中, 禁止使用视图",
			Annotation: "视图的查询性能较差，同时基表结构变更，需要对视图进行维护。如果视图可读性差，且包含复杂的逻辑，会增加维护的成本。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeUsageSuggestion,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止使用视图",
		AllowOffline: true,
		Func:         RuleSQLE00031,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00031): "在 MySQL 中，禁止使用视图."
您应遵循以下逻辑：
1. 对于 "CREATE ..."语句，如果存在以下任何一项，则报告违反规则：
  1. 语法节点中包含视图定义（如CreateStmt中的View定义）
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00031(input *rulepkg.RuleHandlerInput) error {
	// 解析 SQL 语句
	switch input.Node.(type) {
	// 检查是否为 CREATE VIEW 语句
	case *ast.CreateViewStmt:
		// 对于“CREATE VIEW ...” 语句，直接报告违反规则
		rulepkg.AddResult(input.Res, input.Rule, SQLE00031)
		return nil
	}
	return nil
}

// ==== Rule code end ====
