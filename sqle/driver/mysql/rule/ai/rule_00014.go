package ai

import (
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00014 = "SQLE00014"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00014,
			Desc:       "不建议使用自定义函数",
			Annotation: "自定义函数和存储过程维护较困难，且依赖性高，可能导致SQL无法跨库使用。此外，它们在使用时存在一些限制，如无法使用事务相关语句、无法直接产生输出的语句，以及无法在函数体内使用USE语句指定数据库。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "不建议使用自定义函数",
		AllowOffline: true,
		Func:         RuleSQLE00014,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
