package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00153 = "SQLE00153"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00153,
			Desc:       "对于MySQL的DDL, 创建表建议添加索引",
			Annotation: "规划和设计表时，索引应根据业务需求和数据分布合理创建，无索引通常是不合理的情况",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 创建表建议添加索引.",
		AllowOffline: true,
		Func:    RuleSQLE00153,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00153): "In table definition, secondary index must be used when create table".
You should follow the following logic:
1. For "create table ..." statement, check if there is any secondary index definition, and if not, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00153(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		found := false

		// check if secondary index is defined
		constraints := util.GetTableConstraints(stmt.Constraints,
			ast.ConstraintIndex,
			ast.ConstraintUniqIndex,
			ast.ConstraintKey,
			ast.ConstraintUniq,
			ast.ConstraintUniqKey,
		)
		if len(constraints) > 0 {
			found = true
		}

		// check if index is defined in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) {
				found = true
			}
		}

		if !found {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00153)
		}
	}
	return nil
}

// ==== Rule code end ====
