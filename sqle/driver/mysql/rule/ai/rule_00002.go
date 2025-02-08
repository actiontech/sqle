package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00002 = "SQLE00002"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00002,
			Desc:       plocale.Rule00002Desc,
			Annotation: plocale.Rule00002Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "100",
				Desc:  plocale.Rule00002Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00002Message,
		Func:    RuleSQLE00002,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00002): "在 MySQL 中，SQL绑定的变量个数不建议超过阈值.默认参数描述: 绑定变量阈值, 默认参数值: 100"
您应遵循以下逻辑：
1. 针对DML语句中的SELECT子句（包括SELECT、INSERT、UPDATE、DELETE、UNION语句中的子SELECT语句），执行以下检查步骤：
   1. 获取SQL语句中占位符的数量，使用递归遍历语法树实现。
   2. 将统计出的占位符数量与规则配置的阈值（默认值为100）进行比较。
   3. 如果占位符数量超过配置的阈值，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00002(input *rulepkg.RuleHandlerInput) error {
	placeholdersCount := 0
	placeholdersLimit := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()

	var calculateSelectStmt func(node ast.Node) int
	calculateSelectStmt = func(node ast.Node) int {
		count := 0
		switch stmt := node.(type) {
		case *ast.SelectStmt:
			if whereStmt, ok := stmt.Where.(*ast.PatternInExpr); ok && stmt.Where != nil {
				for i := range whereStmt.List {
					item := whereStmt.List[i]
					if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
						count++
					}
				}
			}

			if stmt.Fields != nil {
				for i := range stmt.Fields.Fields {
					item := stmt.Fields.Fields[i]
					if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
						count++
					}
				}
			}

			if stmt.GroupBy != nil {
				for i := range stmt.GroupBy.Items {
					item := stmt.GroupBy.Items[i]
					if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
						count++
					}
				}
			}

			if stmt.Having != nil && stmt.Having.Expr != nil {
				item := stmt.Having.Expr
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					count++
				}
			}

			if stmt.OrderBy != nil {
				for i := range stmt.OrderBy.Items {
					item := stmt.OrderBy.Items[i]
					if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
						count++
					}
				}
			}
			if stmt.Limit != nil {
				if _, ok := stmt.Limit.Count.(*parserdriver.ParamMarkerExpr); ok && stmt.Limit.Count != nil {
					count++
				}
				if _, ok := stmt.Limit.Offset.(*parserdriver.ParamMarkerExpr); ok && stmt.Limit.Offset != nil {
					count++
				}
			}

		case *ast.UnionStmt:
			if stmt.OrderBy != nil {
				for i := range stmt.OrderBy.Items {
					item := stmt.OrderBy.Items[i]
					if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
						count++
					}
				}
			}
			for _, sel := range stmt.SelectList.Selects {
				count += calculateSelectStmt(sel)
			}
		}
		return count
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt:
		placeholdersCount = calculateSelectStmt(stmt)
	case *ast.InsertStmt:
		for i := range stmt.Lists {
			for j := range stmt.Lists[i] {
				item := stmt.Lists[i][j]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}
		for i := range stmt.Setlist {
			if _, ok := stmt.Setlist[i].Expr.(*parserdriver.ParamMarkerExpr); ok && stmt.Setlist[i].Expr != nil {
				placeholdersCount++
			}
		}
		for i := range stmt.OnDuplicate {
			if _, ok := stmt.OnDuplicate[i].Expr.(*parserdriver.ParamMarkerExpr); ok && stmt.OnDuplicate[i].Expr != nil {
				placeholdersCount++
			}
		}
		if stmt.Select != nil {
			placeholdersCount += calculateSelectStmt(stmt.Select)
		}

	case *ast.UpdateStmt:
		for i := range stmt.List {
			item := stmt.List[i]
			if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
				placeholdersCount++
			}
		}
		if whereStmt, ok := stmt.Where.(*ast.PatternInExpr); ok && stmt.Where != nil {
			for i := range whereStmt.List {
				item := whereStmt.List[i]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}
		if stmt.Order != nil {
			for i := range stmt.Order.Items {
				item := stmt.Order.Items[i]
				if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
					placeholdersCount++
				}
			}
		}

	case *ast.DeleteStmt:
		if whereStmt, ok := stmt.Where.(*ast.PatternInExpr); ok && stmt.Where != nil {
			for i := range whereStmt.List {
				item := whereStmt.List[i]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}
	}

	if placeholdersCount > placeholdersLimit {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00002)
	}

	return nil
}

// ==== Rule code end ====
