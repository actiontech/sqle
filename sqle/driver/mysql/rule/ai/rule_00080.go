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

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00080 = "SQLE00080"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00080,
			Desc:       plocale.Rule00080Desc,
			Annotation: plocale.Rule00080Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagSecurity.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "100",
				Desc:  plocale.Rule00080Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00080Message,
		Func:    RuleSQLE00080,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00080): "在 MySQL 中，建议单条SQL写入数据的行数不超过阈值.默认参数描述: 单条SQL写入行数上限, 默认参数值: 100"
您应遵循以下逻辑：
1. 对于 "INSERT...VALUES ..." 语句：
   1. 检查 VALUES 后的数据行数。
   2. 如果数据行数大于阈值，报告违反规则。

2. 对于 "REPLACE ... VALUES ..." 语句：
   1. 执行与 "INSERT...VALUES ..." 相同的检查。

3. 对于 "INSERT...SELECT ..." 语句：
   1. 连接数据库。
   2. 使用辅助函数GetExecutionPlan获取SELECT的执行计划估算数据行数。
   3. 如果估算行数大于阈值，报告违反规则。

4. 对于 "REPLACE ... SELECT ..." 语句：
   1. 执行与 "INSERT...SELECT ..." 相同的检查。

5. 对于 UNION 语句：
   1. 递归检查所有 SELECT 子句。
   2. 对每个 SELECT 子句执行与 "INSERT...SELECT ..." 相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00080(input *rulepkg.RuleHandlerInput) error {
	// 获取阈值参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be an integer", param.Value)
	}

	// 内部辅助函数
	getSelectRowCount := func(stmt *ast.InsertStmt) (int64, error) {
		switch stmt.Select.(type) {
		case *ast.SelectStmt, *ast.UnionStmt:
			explain, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
			if err != nil {
				return 0, fmt.Errorf("failed to get execution plan: %v", err)
			}

			// 假设 GetExecutionPlan 返回的执行计划中有一个估算的行数
			for _, record := range explain.Plan {
				if record.Rows > int64(threshold) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00080)
					return record.Rows, nil
				}
			}
		}
		return 0, nil
	}

	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		// 处理 INSERT,REPLACE 语句
		if stmt.Select != nil {
			// INSERT ... SELECT ...
			_, err := getSelectRowCount(stmt)
			if err != nil {
				log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
				return err
			}
		} else if len(stmt.Lists) > 0 {
			// INSERT ... VALUES ...
			rowCount := len(stmt.Lists)
			if rowCount > threshold {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00080)
			}
		}
	default:
		// 其他类型的语句不处理
		return nil
	}

	return nil
}

// ==== Rule code end ====
