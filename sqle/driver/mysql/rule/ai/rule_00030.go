package ai

import (
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00030 = "SQLE00030"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00030,
			Desc:       "在 MySQL 中, 禁止使用触发器",
			Annotation: "触发器难以开发和维护，不能高效移植，且在复杂的逻辑以及高并发下，容易出现死锁影响业务。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeUsageSuggestion,
			Params:     params.Params{},
		},
		Message: "在 MySQL 中, 禁止使用触发器",
		Func:    RuleSQLE00030,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00030): "在 MySQL 中，禁止使用触发器."
您应遵循以下逻辑：
1. 对于 "CREATE ..." 语句，如果存在以下任何一项，则报告违反规则：
  1. 语法树中包含触发器定义节点
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00030(input *rulepkg.RuleHandlerInput) error {
	// 解析 SQL 语句
	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createTriggerReg1.MatchString(input.Node.Text()) ||
			createTriggerReg2.MatchString(input.Node.Text()) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00030)
		}
	}
	return nil
}

var createTriggerReg1 = regexp.MustCompile(`(?i)create[\s]+trigger[\s]+[\S\s]+(before|after)+`)
var createTriggerReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+trigger[\s]+[\S\s]+(before|after)+`)

// ==== Rule code end ====
