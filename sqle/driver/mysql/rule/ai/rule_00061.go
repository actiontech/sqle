package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00061 = "SQLE00061"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00061,
			Desc:       "建议新建表句子中包含表存在判断操作",
			Annotation: "新建表如果已经存在，不加 IF NOT EXISTS 会报错。新建表只在表不存在的前提下进行，避免SQL 实际执行报错。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeUsageSuggestion,
			Params:     params.Params{},
		},
		Message:      "建议新建表句子中包含表存在判断操作",
		AllowOffline: true,
		Func:         RuleSQLE00061,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00061): "在 MySQL 中，建议新建表句子中包含表存在判断操作."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE ..."语句，如果存在以下任何一项，则报告违反规则：
  1. 句子中不包含语法节点：IF NOT EXISTS
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00061(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 "CREATE TABLE ..." 语句中是否存在 "IF NOT EXISTS"
		if !stmt.IfNotExists {
			// 如果不存在 "IF NOT EXISTS"，则报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00061)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
