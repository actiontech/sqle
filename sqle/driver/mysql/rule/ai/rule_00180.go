package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00180 = "SQLE00180"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00180,
			Desc:       plocale.Rule00180Desc,
			Annotation: plocale.Rule00180Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			Level:      driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "2",
				Desc:  plocale.Rule00180Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00180Message,
		Func:    RuleSQLE00180,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00180): "在 MySQL 中，避免执行计划中 filter 次数过多."
您应遵循以下逻辑：
1. 登录数据库。
2. 对于 "DML" 语句，使用辅助函数GetExecutionTreePlan获取SQL语句的执行计划，如果执行计划中包含 filter 语法节点，则记录 filter 语法节点的出现次数，并与规则阈值进行比较。如果出现次数超过阈值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00180(input *rulepkg.RuleHandlerInput) error {
	// 获取规则参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	threshold := param.Int()
	if threshold == 0 {
		return fmt.Errorf("param value should be greater than 0")
	}

	checkFilterViolation := func(plan string) bool {
		filterCount := strings.Count(plan, "Filter")

		return filterCount > threshold
	}

	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}

	plan, err := util.GetExecutionTreePlan(input.Ctx, input.Node.Text())
	if err != nil {
		return err
	}

	if violation := checkFilterViolation(plan); violation {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00180)
		return nil
	}

	return nil
}

// ==== Rule code end ====
