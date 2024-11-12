package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00060 = "SQLE00060"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00060,
			Desc:       "表建议添加注释",
			Annotation: "表添加注释能够使表的意义更明确，方便日后的维护",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params:     params.Params{},
		},
		Message:      "表建议添加注释",
		AllowOffline: true,
		Func:         RuleSQLE00060,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00060): "在 MySQL 中，表建议添加注释."
您应遵循以下逻辑：
1、检查CREATE TABLE语法节点末尾是否都包含注释节点，否则，将该完整SQL语句加入到触发规则的SQL表表中。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00060(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// Check if the table has a comment at the end
		hasComment := false
		for _, opt := range stmt.Options {
			if opt.Tp == ast.TableOptionComment {
				hasComment = true
				break
			}
		}

		// If no comment is found, add the SQL to the result
		if !hasComment {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00060)
		}
	}

	return nil
}

// ==== Rule code end ====
