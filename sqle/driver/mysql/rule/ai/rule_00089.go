package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00089 = "SQLE00089"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00089,
			Desc:       "在 MySQL 中, 禁止INSERT ... SELECT",
			Annotation: "使用 INSERT ... SELECT 在默认事务隔离级别下，可能会导致对查询的表施加表级锁",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止INSERT ... SELECT",
		AllowOffline: true,
		Func:         RuleSQLE00089,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00089): "在 MySQL 中，禁止INSERT ... SELECT."
您应遵循以下逻辑：
1. 针对 "INSERT ... SELECT" 语句，则报告违反规则
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00089(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		if stmt.Select != nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00089)
		}
	}
	// TODO INSERT ... WITH
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
