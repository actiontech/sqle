package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00107 = "SQLE00107"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00107,
			Desc:       plocale.Rule00107Desc,
			Annotation: plocale.Rule00107Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelNotice,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "1024",
				Desc:  plocale.Rule00107Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00107Message,
		Func:    RuleSQLE00107,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
