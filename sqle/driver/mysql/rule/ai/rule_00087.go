package ai

import (
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00087 = "SQLE00087"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00087,
			Desc:       plocale.Rule00087Desc,
			Annotation: plocale.Rule00087Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelError,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "500",
				Desc:  plocale.Rule00087Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00087Message,
		Func:    RuleSQLE00087,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00087): "在 MySQL 中，避免WHERE条件内IN语句中的参数值个数过多.默认参数描述: IN的参数值个数, 默认参数值: 500"
您应遵循以下逻辑：
1. 针对以下SQL语句类型进行检查：SELECT、WITH、INSERT ... SELECT、UPDATE、DELETE、UNION。
2. 对于每种语句中所有的WHERE条件：
   1. 检查WHERE条件中的IN、NOT IN列表：
      1. 计算IN列表中的元素数量。
      2. 如果元素数量超过当前规则的阈值，则报告违反规则。
   2.检查WHERE条件中IN、NOT IN子查询：
	  1. 使用辅助函数GetExecutionPlan获取子查询的执行计划，递归遍历执行计划树。
	  2. 获取子查询的扫描行数。
	  3. 如果扫描行数超过当前规则的阈值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00087(input *rulepkg.RuleHandlerInput) error {
	inListThresholdParam := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if inListThresholdParam == nil {
		return fmt.Errorf("param not found")
	}
	inListThreshold, err := strconv.Atoi(inListThresholdParam.Value)
	if err != nil {
		return fmt.Errorf("param should be an integer, got: %v", inListThresholdParam.Value)
	}

	// 内部匿名的辅助函数
	calculateScanRows := func(plans []*executor.ExplainRecord) int64 {
		totalRows := int64(0)
		for _, record := range plans {
			totalRows += record.Rows
		}
		return totalRows
	}

	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}
	whereList := util.GetWhereExprFromDMLStmt(input.Node)
	for _, where := range whereList {
		isViolate := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch e := expr.(type) {
			case *ast.PatternInExpr:
				if len(e.List) > inListThreshold {
					isViolate = true
					return true
				}
				if e.Sel != nil {
					if subExpr, ok := e.Sel.(*ast.SubqueryExpr); ok {
						executionPlan, err := util.GetExecutionPlan(input.Ctx, subExpr.Query.Text()) // SelectStmt、UnionStmt
						if err != nil {
							log.NewEntry().Errorf("Failed to get execution plan for subquery: %v", err)
							return false
						}
						scanRows := calculateScanRows(executionPlan.Plan)
						if scanRows > int64(inListThreshold) {
							isViolate = true
							return true
						}
					}
				}
			}
			return true
		}, where)

		if isViolate {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00087)
			return nil
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
