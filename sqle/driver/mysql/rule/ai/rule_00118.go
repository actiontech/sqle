package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00118 = "SQLE00118"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00118,
			Desc:       "在 MySQL 中, 建议在执行DROP/TRUNCATE等操作前进行备份",
			Annotation: "DROP/TRUNCATE是DDL，操作立即生效，不会写入日志，所以无法回滚，在执行高危操作之前对数据进行备份是很有必要的",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议在执行DROP/TRUNCATE等操作前进行备份",
		AllowOffline: true,
		Func:         RuleSQLE00118,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00118): "在 MySQL 中，建议在执行DROP/TRUNCATE等操作前进行备份."
您应遵循以下逻辑：
1. 对于提供的SQL语句，执行以下检查，任何条件满足则报告违反规则：
    1. 语句中包含有 "DROP TABLE" 子句。
    2. 语句中包含有 "TRUNCATE TABLE" 子句。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00118(input *rulepkg.RuleHandlerInput) error {

	switch input.Node.(type) {
	case *ast.DropTableStmt, *ast.TruncateTableStmt:
		rulepkg.AddResult(input.Res, input.Rule, SQLE00118)
	}

	return nil
}

// ==== Rule code end ====
