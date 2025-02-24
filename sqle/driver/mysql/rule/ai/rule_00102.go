package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00102 = "SQLE00102"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00102,
			Desc:       plocale.Rule00102Desc,
			Annotation: plocale.Rule00102Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00102Message,
		Func:    RuleSQLE00102,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00102): "在 MySQL 中，禁止UPDATE/DELETE语句使用ORDER BY操作 ."
您应遵循以下逻辑：
1. 对于 "UPDATE ..." 语句，检查以下条件：
   1. 句子中存在关键词： ORDER BY，则报告违反规则。

2. 对于 "DELETE ..." 语句，检查以下条件：
   1. 句子中存在关键词： ORDER BY，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00102(input *rulepkg.RuleHandlerInput) error {
	// 内部匿名的辅助函数
	getOrderByNode := func(node ast.Node) *ast.OrderByClause {
		switch stmt := node.(type) {
		case *ast.UpdateStmt:
			return stmt.Order
		case *ast.DeleteStmt:
			return stmt.Order
		case *ast.UnionStmt:
			return stmt.OrderBy
		case *ast.SelectStmt:
			return stmt.OrderBy
		default:
			return nil
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt, *ast.DeleteStmt:
		// 获取语句中的 ORDER BY 节点
		orderBy := getOrderByNode(stmt)
		if orderBy != nil && len(orderBy.Items) > 0 {
			// 如果存在 ORDER BY 节点，报告违反规则
			rulepkg.AddResult(input.Res, input.Rule, SQLE00102)
			return nil
		}
		// 子查询中
		subs := util.GetSubquery(stmt)
		for _, sub := range subs {
			orderBy := getOrderByNode(sub.Query)
			if orderBy != nil && len(orderBy.Items) > 0 {
				// 如果存在 ORDER BY 节点，报告违反规则
				rulepkg.AddResult(input.Res, input.Rule, SQLE00102)
				return nil
			}
		}
	}
	// 如果不符合条件，不报告任何违规
	return nil
}

// ==== Rule code end ====
