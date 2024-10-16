package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
)

const (
	SQLE00095 = "SQLE00095"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00095,
			Desc:       "在 MySQL 中, 建议使用'<>'代替'!='",
			Annotation: "'<>' 是ANSI SQL标准中定义的不等于运算符。如果使用了!=运算符，数据库优化器会自动转换为SQL标准不等于运算符，增加了优化器的转换开销；另外，目前并非所有的SQL数据库系统都支持 !=，使用标准的运算符可以确保SQL在各数据库之间具有更高的兼容性。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议使用'<>'代替'!='",
		AllowOffline: true,
		Func:         RuleSQLE00095,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00095): "在 MySQL 中，建议使用'<>'代替'!='."
您应遵循以下逻辑：
1. 对于所有DML、DQL语句，如果以下任意一个为真，则报告违反规则：
  1. 语句里的WHERE 条件里存在'!='不等于操作符节点
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00095(input *rulepkg.RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
		// 获取 DML 语句中的 WHERE 条件
		whereList := util.GetWhereExprFromDMLStmt(input.Node)

		// 遍历 WHERE 条件中的每个表达式
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.BinaryOperationExpr:
				// 检查'!="'不等于操作符
				if x.Op == opcode.NE {
					if strings.Contains(input.Node.Text(), "!=") {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00095)
						return true
					}
				}
			}
			return false
		}, whereList...)
	}
	return nil
}

// ==== Rule code end ====
