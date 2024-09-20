package ai

import (
	"fmt"
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00029 = "SQLE00029"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00029,
			Desc:       "在 MySQL 中, 禁止使用存储过程",
			Annotation: "存储过程在一定程度上能使程序难以调试和拓展，各种数据库端的存储过程语法相差很大，给将来的数据移植带来很大的困难，且会极大的出现BUG的几率",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeUsageSuggestion,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止使用存储过程",
		Func:         RuleSQLE00029,
		AllowOffline: false,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
			fmt.Print("Hello, ")
			rulepkg.AddResult(input.Res, input.Rule, SQLE00029)
		} else if alterProcedureReg1.MatchString(input.Node.Text()) ||
			alterProcedureReg2.MatchString(input.Node.Text()) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00029)
		}
	}
	return nil
}

var createProcedureReg1 = regexp.MustCompile(`(?i)create[\s]+procedure[\s]+[\S\s]+`)
var createProcedureReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+procedure[\s]+[\S\s]+`)

var alterProcedureReg1 = regexp.MustCompile(`(?i)alter[\s]+procedure[\s]+[\S\s]+`)
var alterProcedureReg2 = regexp.MustCompile(`(?i)alter[\s]+[\s\S]+[\s]+procedure[\s]+[\S\s]+`)

// ==== Rule code end ====
