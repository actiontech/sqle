package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00107 = "SQLE00107"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00107,
			Desc:       "建议将过长的SQL分解成几个简单的SQL",
			Annotation: "过长的SQL可读性较差，难以维护，且容易引发性能问题。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "1024",
					Desc:  "句子长度限制",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "建议将过长的SQL分解成几个简单的SQL",
		AllowOffline: true,
		Func:         RuleSQLE00107,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00107): "在 MySQL 中，建议将过长的SQL分解成几个简单的SQL.默认参数描述: 句子长度限制, 默认参数值: 1024"
您应遵循以下逻辑：
1. 针对所有 DML 语句（包括 SELECT、UPDATE、DELETE、INSERT...SELECT、UNION、WITH），执行以下步骤：
   1. 计算语句的字符串长度，如果大于等于阈值，则报告违反规则
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00107(input *rulepkg.RuleHandlerInput) error {
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param should be an integer, got: %v", param.Value)
	}

	if stmt, ok := input.Node.(ast.DMLNode); ok {
		if len(stmt.Text()) > threshold {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00107)
		}
	}
	return nil
}

// ==== Rule code end ====
