package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00067 = "SQLE00067"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00067,
			Desc:       "表不建议使用外键",
			Annotation: "外键在大量写入场景下性能较差，强烈禁止使用",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params:     params.Params{},
		},
		Message:      "表不建议使用外键",
		AllowOffline: true,
		Func:         RuleSQLE00067,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00067): "在 MySQL 中，表不建议使用外键."
您应遵循以下逻辑：
1. 检查CREATE TABLE 语句或者 ALTER TABLE 语句的语法节点，查看有无外键定义，如果存在外键定义，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00067(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE TABLE 语句中的外键约束
		if constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintForeignKey); constraint != nil {
			// 存在外键约束，报告违规
			rulepkg.AddResult(input.Res, input.Rule, SQLE00067)
			return nil
		}
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中的外键约束
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			if constraint := util.GetTableConstraint([]*ast.Constraint{spec.Constraint}, ast.ConstraintForeignKey); constraint != nil {
				// 存在外键约束，报告违规
				rulepkg.AddResult(input.Res, input.Rule, SQLE00067)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
