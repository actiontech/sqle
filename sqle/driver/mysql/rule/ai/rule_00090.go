package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00090 = "SQLE00090"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00090,
			Desc:       plocale.Rule00090Desc,
			Annotation: plocale.Rule00090Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00090Message,
		Func:    RuleSQLE00090,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00090): "在 MySQL 中，建议使用UNION ALL替代UNION."
您应遵循以下逻辑：
1. 对于所有DML语句，
    1. 含有UNION语法节点，其中UNION语法节点类型不是UNION ALL，则报告违反规则
2. 对于"WITH ..."语句，执行与上述同样检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00090(input *rulepkg.RuleHandlerInput) error {
	// 内部匿名的辅助函数

	hasUnionWithoutAll := func(uStmt *ast.UnionStmt) bool {
		for _, ss := range uStmt.SelectList.Selects {
			if ss.IsAfterUnionDistinct {
				return true
			}
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.UnionStmt:
		if hasUnionWithoutAll(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00090)
			return nil
		}
	case *ast.InsertStmt:
		if uStmt, ok := stmt.Select.(*ast.UnionStmt); ok {
			if hasUnionWithoutAll(uStmt) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00090)
				return nil
			}
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UnionStmt:
		// 获取所有相关的 SELECT 语句
		subs := util.GetSubquery(stmt)
		for _, sub := range subs {
			// 检查 SELECT 语句中是否存在 UNION 而没有 UNION ALL
			if uStmt, ok := sub.Query.(*ast.UnionStmt); ok {
				if hasUnionWithoutAll(uStmt) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00090)
					return nil
				}
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
