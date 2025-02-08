package ai

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00082 = "SQLE00082"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00082,
			Desc:       plocale.Rule00082Desc,
			Annotation: plocale.Rule00082Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00082Message,
		Func:    RuleSQLE00082,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
