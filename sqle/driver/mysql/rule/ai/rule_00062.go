package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	dry "github.com/ungerik/go-dry"
)

const (
	SQLE00062 = "SQLE00062"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00062,
			Desc:       "建议事务隔离级别设置成RC",
			Annotation: "RC 虽然没有解决幻读的问题，但是没有间隙锁，从而每次在做更新操作时影响的行数比默认RR要小很多；默认的RR隔离级别虽然解决了幻读问题，但是增加了间隙锁，导致加锁的范围扩大，性能比RC要低，增加死锁的概率；在大多数情况下，出现幻读的几率较小，所以建议使用RC。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "建议事务隔离级别设置成RC",
		AllowOffline: true,
		Func:         RuleSQLE00062,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
