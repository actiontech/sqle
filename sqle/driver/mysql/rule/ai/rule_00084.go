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
	SQLE00084 = "SQLE00084"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00084,
			Desc:       plocale.Rule00084Desc,
			Annotation: plocale.Rule00084Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00084Message,
		Func:    RuleSQLE00084,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00084): "在 MySQL 中，不建议使用临时表."
您应遵循以下逻辑：
1. 对于“CREATE ... ”语句，执行以下检查：
   1. 如果 CREATE 语句包含 TEMPORARY 关键词，则报告违反规则。

2. 对于 DML 语句（如 SELECT、INSERT...SELECT、UNION、UPDATE、DELETE），执行以下步骤：
   1. 使用辅助函数GetExecutionPlan查看 SQL 语句的执行计划。
   2. 如果执行计划中是否包含 Using temporary；如果包含，则报告违反规则。

3. 对于 WITH 语句，执行以下步骤：
   1. 使用辅助函数GetExecutionPlan查看 SQL 语句的执行计划。
   2. 如果执行计划中是否包含 Using temporary；如果包含，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00084(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE 语句是否包含 TEMPORARY 关键词
		if stmt.IsTemporary {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00084)
			return nil
		}
	default:
		if _, ok := input.Node.(ast.DMLNode); !ok {
			return nil
		}

		explain, err := util.GetExecutionPlan(input.Ctx, input.Node.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
			return err
		}
		for _, record := range explain.Plan {
			if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingTemporary) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00084)
				return nil
			}
		}
		return nil
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
