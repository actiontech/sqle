package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00109 = "SQLE00109"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00109,
			Desc:       plocale.Rule00109Desc,
			Annotation: plocale.Rule00109Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00109Message,
		Func:    RuleSQLE00109,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00109): "在 MySQL 中，禁止在子查询中使用LIMIT."
您应遵循以下逻辑：
1. 对于含有子查询的SQL语句（UPDATE、DELETE、SELECT、INSERT ... SELECT、SELECT ... UNION ALL SELECT），如果子查询中包含LIMIT子句，则报告违反规则。
   1. 检查UPDATE语句中的子查询。递归检查所有嵌套的SELECT节点。
   2. 检查DELETE语句中的子查询。递归检查所有嵌套的SELECT节点。
   3. 检查SELECT语句中的子查询。递归检查所有嵌套的SELECT节点。
   4. 检查INSERT ... SELECT语句中的子查询。递归检查所有嵌套的SELECT节点。
   5. 检查SELECT ... UNION ALL SELECT语句中的子查询。递归检查所有嵌套的SELECT节点。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00109(input *rulepkg.RuleHandlerInput) error {
	// 定义一个递归函数，用于检查SELECT节点中是否包含LIMIT子句
	var checkSelectForLimit func(stmt *ast.SelectStmt) bool
	checkSelectForLimit = func(stmt *ast.SelectStmt) bool {
		// 如果当前SELECT语句包含LIMIT子句，则违反规则
		if stmt.Limit != nil {
			return true
		}

		// 检查当前SELECT语句中的所有子查询
		subqueries := util.GetSubquery(stmt)
		for _, subquery := range subqueries {
			// 获取子查询中的SELECT语句
			selectStmts := util.GetSelectStmt(subquery.Query)
			for _, subSelect := range selectStmts {
				// 递归检查子查询中的SELECT语句
				if checkSelectForLimit(subSelect) {
					return true
				}
			}
		}

		return false
	}

	// 定义一个函数，用于遍历并检查不同类型的SQL语句
	traverseAndCheck := func(node ast.Node) bool {
		switch stmt := node.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt, *ast.SelectStmt, *ast.InsertStmt:
			// 获取当前语句中的所有SELECT语句
			selectStmts := util.GetSelectStmt(stmt)
			for _, selectStmt := range selectStmts {
				// 如果当前SELECT语句或其子查询中包含LIMIT子句，则违反规则
				if checkSelectForLimit(selectStmt) {
					return true
				}
			}
		case *ast.UnionStmt:
			// 递归检查UNION中的每个SELECT部分
			for _, selectStmt := range stmt.SelectList.Selects {
				if checkSelectForLimit(selectStmt) {
					return true
				}
			}
		}
		return false
	}

	// 执行遍历和检查
	if traverseAndCheck(input.Node) {
		// 如果发现违反规则的情况，记录结果
		rulepkg.AddResult(input.Res, input.Rule, SQLE00109)
		return nil
	}

	// 如果没有发现任何违规，返回nil
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
