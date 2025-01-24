package ai

import (
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00029 = "SQLE00029"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00029,
			Desc:         plocale.Rule00029Desc,
			Annotation:   plocale.Rule00029Annotation,
			Category:     plocale.RuleTypeUsageSuggestion,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00029Message,
		Func:    RuleSQLE00029,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00029): "在 MySQL 中，禁止使用存储过程."
您应遵循以下逻辑：
1. 对于 "CREATE ..." 语句，如果存在以下任何一项，则报告违反规则：
  1. 语法节点中包含 PROCEDURE。
2. 对于 "ALTER..." 语句，执行与上述同样检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00029(input *rulepkg.RuleHandlerInput) error {
	// 检查是否为 CREATE PROCEDURE 或 ALTER PROCEDURE 语句

	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createProcedureReg1.MatchString(input.Node.Text()) ||
			createProcedureReg2.MatchString(input.Node.Text()) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00029)
		} else if alterProcedureReg1.MatchString(input.Node.Text()) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00029)
		}
	}
	return nil
}

var createProcedureReg1 = regexp.MustCompile(`(?i)create[\s]+procedure[\s]+[\S\s]+`)
var createProcedureReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+procedure[\s]+[\S\s]+`)

var alterProcedureReg1 = regexp.MustCompile(`(?i)alter[\s]+procedure[\s]+[\S\s]+`)

// ==== Rule code end ====
