package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00008 = "SQLE00008"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00008,
			Desc:       "对于MySQL的索引, 表必须有主键",
			Annotation: "主键使数据达到全局唯一，可提高数据检索效率。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeIndexingConvention,
		},
		Message: "对于MySQL的索引, 表必须有主键",
		AllowOffline: false,
		Func:    RuleSQLE00008,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00008): "In table definition, The primary key should exist".
You should follow the following logic:
1. For "create table ... primary key..." statement, if the primary key which on the table option or on the column is not exist, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00008(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}

		found := false

		// check primary key in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnPrimaryKey(col) {
				found = true
				break
			}
		}

		// check primary key in table constraint
		constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintPrimaryKey)
		if nil != constraint {
			//this is a table primary key definition
			found = true
		}

		if !found {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00008)
			return nil
		}
	}
	return nil
}

// ==== Rule code end ====
