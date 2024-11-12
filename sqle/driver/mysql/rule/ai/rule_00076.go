package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00076 = "SQLE00076"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00076,
			Desc:       "UPDATE/DELETE操作影响行数不建议超过阈值",
			Annotation: "在数据库中，进行修改或删除等数据变更操作时，一次性操作的数据量过大，会消耗大量的系统资源，产生长事务，会导致查询性能下降，影响其他事务或查询的执行。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "10000",
					Desc:  "影响行数上限",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "UPDATE/DELETE操作影响行数不建议超过阈值",
		AllowOffline: false,
		Func:         RuleSQLE00076,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00076): "在 MySQL 中，UPDATE/DELETE操作影响行数不建议超过阈值.默认参数描述: 影响行数上限, 默认参数值: 10000"
您应遵循以下逻辑：
1. 对于 "UPDATE ..." 语句，连接到数据库，使用辅助函数GetExecutionPlan递归检查所有嵌套的 SELECT 语句，获取操作类型为 UPDATE 的估算影响行数。如果估算行数超过预设阈值，则标记为违反规则。
2. 对于 "DELETE ..." 语句，连接到数据库，使用辅助函数GetExecutionPlan递归检查所有嵌套的 SELECT 语句，获取操作类型为 DELETE 的估算影响行数。如果估算行数超过预设阈值，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始

func RuleSQLE00076(input *rulepkg.RuleHandlerInput) error {
	// 获取规则参数中的预设阈值
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param should be an integer, got: %v", param.Value)
	}

	// 确认输入的 SQL 语句类型为 UPDATE 或 DELETE
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt, *ast.DeleteStmt:
		// 获取 SQL 语句文本
		sqlText := stmt.Text()

		// 连接到数据库并获取执行计划
		explain, err := util.GetExecutionPlan(input.Ctx, sqlText)
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return err
		}
		for _, record := range explain.Plan {
			if record.Rows > int64(threshold) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00076)
				return nil
			}
		}
	}

	// 如果所有检查均通过，返回 nil
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
