package ai

import (
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00014 = "SQLE00014"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00014,
			Desc:         plocale.Rule00014Desc,
			Annotation:   plocale.Rule00014Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00014Message,
		Func:    RuleSQLE00014,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00014): "在 MySQL 中，不建议使用自定义函数."
您应遵循以下逻辑：
1. 针对 SQL 语句进行解析，识别到 "CREATE FUNCTION ..." 语句，直接报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00014(input *rulepkg.RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createFunctionReg1.MatchString(input.Node.Text()) ||
			createFunctionReg2.MatchString(input.Node.Text()) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00014)
		}
	}
	return nil
}

var createFunctionReg1 = regexp.MustCompile(`(?i)create[\s]+function[\s]+[\S\s]+returns`)
var createFunctionReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+function[\s]+[\S\s]+returns`)

// ==== Rule code end ====
