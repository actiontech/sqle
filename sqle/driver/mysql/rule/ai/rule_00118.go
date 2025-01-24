package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00118 = "SQLE00118"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00118,
			Desc:         plocale.Rule00118Desc,
			Annotation:   plocale.Rule00118Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00118Message,
		Func:    RuleSQLE00118,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
