package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00151 = "SQLE00151"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00151,
			Desc:       "避免CREATE TABLE/ALTER TABLE 使用禁止的表空间",
			Annotation: "不允许在系统表空间上创建用户对象，避免不必要的安全风险，方便维护",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "避免CREATE TABLE/ALTER TABLE 使用禁止的表空间.",
		AllowOffline: true,
		Func:    RuleSQLE00151,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00151): "In table definition, using system tablespace are prohibited.".
You should follow the following logic:
1. For "CREATE TABLE..." Statement, check whether the statement has the keyword tablespace innodb_system, if so, report the rule violation.
2. For "ALTER TABLE... TABLESPACE..."  Statement, check whether the statement has the keyword tablespace innodb_system, if so, report the rule violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00151(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		for _, option := range stmt.Options {
			// "create table..."
			if option.Tp == ast.TableOptionTablespace && strings.EqualFold(option.StrValue, "innodb_system") {
				//"create table... tablespace innodb_system..."
				rulepkg.AddResult(input.Res, input.Rule, SQLE00151)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableOption) {
			// "alter table... tablespace ..."
			if len(spec.Options) > 0 && strings.EqualFold(spec.Options[0].StrValue, "innodb_system") {
				//"alter table... tablespace innodb_system..."
				rulepkg.AddResult(input.Res, input.Rule, SQLE00151)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
