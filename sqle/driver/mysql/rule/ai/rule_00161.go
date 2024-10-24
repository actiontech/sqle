package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
)

const (
	SQLE00161 = "SQLE00161"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00161,
			Desc:       "在 MySQL 中, 建议序列或自增字段的步长为1",
			Annotation: "序列或自增字段的步长为1时，有助于保证主键和其他自增字段的连续性，避免不必要的数据间隔和数字资源的浪费。不仅简化了数据库的管理和维护，而且也提高了系统的可预测性和稳定性。特别是在处理大量数据插入或高并发场景时，连续的主键值还能减少潜在的冲突和错误。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议序列或自增字段的步长为1",
		AllowOffline: true,
		Func:         RuleSQLE00161,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
