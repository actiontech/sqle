package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	dry "github.com/ungerik/go-dry"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00062 = "SQLE00062"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00062,
			Desc:       plocale.Rule00062Desc,
			Annotation: plocale.Rule00062Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagTransaction.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00062Message,
		Func:    RuleSQLE00062,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00062): "在 MySQL 中，建议事务隔离级别设置成RC."
您应遵循以下逻辑：
1. 对于 "SET ... TRANSACTION...ISOLATION LEVEL..."语句，如果不存在语法节点表示的隔离级别为 READ COMMITTED，则报告违反规则。
2. 对于 "SET ... transaction_isolation ... " 语句，如果不存在语法节点表示的隔离级别为 READ-COMMITTED，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00062(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SetStmt:
		for _, variable := range stmt.Variables {
			if dry.StringListContains([]string{"tx_isolation", "tx_isolation_one_shot", "transaction_isolation"}, variable.Name) {
				switch node := variable.Value.(type) {
				case *parserdriver.ValueExpr:
					if node.Datum.GetString() != ast.ReadCommitted {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00062)
						return nil
					}
				}
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
