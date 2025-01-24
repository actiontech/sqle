package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00161 = "SQLE00161"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00161,
			Desc:         plocale.Rule00161Desc,
			Annotation:   plocale.Rule00161Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00161Message,
		Func:    RuleSQLE00161,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00161): "在 MySQL 中，建议序列或自增字段的步长为1."
您应遵循以下逻辑：
1. 对于 "SET..." 语句，执行以下检查：
   1. 确认目标对象为系统参数 auto_increment_increment，其参数值不等于 1，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00161(input *rulepkg.RuleHandlerInput) error {
	// 确认语句是否为SET类型
	setStmt, ok := input.Node.(*ast.SetStmt)
	if !ok {
		return nil
	}

	// 遍历所有设置的变量
	for _, variable := range setStmt.Variables {
		// 获取设置的变量名
		varName := variable.Name

		// 确认目标对象为'auto_increment_increment'
		if strings.EqualFold(varName, "auto_increment_increment") {
			if v, ok := variable.Value.(*parserdriver.ValueExpr); ok {
				if v.Datum.GetInt64() != 1 {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00161)
					return nil
				}
			}
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
