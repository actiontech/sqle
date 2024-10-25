package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00004 = "SQLE00004"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00004,
			Desc:       "对于MySQL的DDL, 表的初始AUTO_INCREMENT值建议为0",
			Annotation: "创建表时AUTO_INCREMENT设置为0则自增从1开始，可以避免数据空洞。例如在导出表结构DDL时，表结构内AUTO_INCREMENT通常为当前的自增值，如果建表时没有把AUTO_INCREMENT设置为0，那么通过该DDL进行建表操作会导致自增值从一个无意义数字开始。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 表的初始AUTO_INCREMENT值建议为0",
		AllowOffline: true,
		Func:    RuleSQLE00004,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00004): "In table definition, the initial AUTO_INCREMENT value of a column should be 0".
You should follow the following logic:
1. For "create table ... auto_increment=..." statement, if the auto_increment value is specified other than 0, report a violation
2. For "alter table ... auto_increment=..." statement, if the auto_increment value is specified other than 0, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00004(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table"
		if option := util.GetTableOption(stmt.Options, ast.TableOptionAutoIncrement); nil != option {
			//"create table ... auto_increment=..."
			if option.UintValue != 0 {
				// the table option "auto increment" is other than 0
				rulepkg.AddResult(input.Res, input.Rule, SQLE00004)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableOption) {
			if option := util.GetTableOption(spec.Options, ast.TableOptionAutoIncrement); nil != option {
				//"alter table ... auto_increment=..."
				if option.UintValue != 0 {
					// the table option "auto increment" is other than 0
					rulepkg.AddResult(input.Res, input.Rule, SQLE00004)
					return nil
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
