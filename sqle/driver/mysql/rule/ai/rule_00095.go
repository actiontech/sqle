package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00095 = "SQLE00095"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00095,
			Desc:       plocale.Rule00095Desc,
			Annotation: plocale.Rule00095Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00095Message,
		Func:    RuleSQLE00095,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00095): "在 MySQL 中，建议使用'<>'代替'!='."
您应遵循以下逻辑：
1. 对于所有DML、DQL语句，如果以下任意一个为真，则报告违反规则：
  1. 语句里的WHERE 条件里存在'!='不等于操作符节点
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00095(input *rulepkg.RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
		// 获取 DML 语句中的 WHERE 条件
		whereList := util.GetWhereExprFromDMLStmt(input.Node)

		// 遍历 WHERE 条件中的每个表达式
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.BinaryOperationExpr:
				// 检查'!="'不等于操作符
				if x.Op == opcode.NE {
					if strings.Contains(input.Node.Text(), "!=") {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00095)
						return true
					}
				}
			}
			return false
		}, whereList...)
	}
	return nil
}

// ==== Rule code end ====
