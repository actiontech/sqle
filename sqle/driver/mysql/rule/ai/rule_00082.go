package ai

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00082 = "SQLE00082"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00082,
			Desc:       "在 MySQL 中, 禁止使用文件排序",
			Annotation: "大数据量的情况下，文件排序意味着SQL性能较低，会增加OS的开销，影响数据库性能。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止使用文件排序",
		AllowOffline: false,
		Func:         RuleSQLE00082,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00082): "在 MySQL 中，禁止使用文件排序."
您应遵循以下逻辑：
1. 登录数据库。
2. 使用辅助函数GetExecutionPlan获取SQL语句的执行计划，选择适当的格式：
   1. 对于 explain format=traditional：
      - 检查执行计划中是否包含语法节点 "Using filesort"； 如果包含，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====

func RuleSQLE00082(input *rulepkg.RuleHandlerInput) error {
	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}

	explain, err := util.GetExecutionPlan(input.Ctx, input.Node.Text())
	if err != nil {
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
		return err
	}
	for _, record := range explain.Plan {
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingFilesort) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00082)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
