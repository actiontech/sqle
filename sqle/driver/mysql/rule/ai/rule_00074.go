package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00074 = "SQLE00074"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00074,
			Desc:       "对于MySQL的DDL, 禁止使用RENAME或CHANGE对表名字段名进行修改",
			Annotation: "RENAME/CHANGE 表名/列名会对线上业务不停机发布造成影响，如需这种操作应当DBA手工干预",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 禁止使用RENAME或CHANGE对表名字段名进行修改",
		Func:    RuleSQLE00074,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00074): "In DDL, renaming or changing table and column names is prohibited".
You should follow the following logic:
1. For "alter table ... rename table ..." statement, report a violation
2. For "alter table ... rename column ..." statement, report a violation
3. For "alter table ... change column ..." statement, if the new column name is different from old column name, report a violation
4. For "rename table ..." statement, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00074(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		for range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableRenameTable, ast.AlterTableRenameColumn) {
			//"alter table ... rename table ..." or "alter table ... rename column ..."
			rulepkg.AddResult(input.Res, input.Rule, SQLE00074)
			return nil
		}
		for _, cmd := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableChangeColumn) {
			//"alter table ... change column ..."
			if cmd.OldColumnName != cmd.NewColumns[0].Name {
				//the column name is changed
				rulepkg.AddResult(input.Res, input.Rule, SQLE00074)
				return nil
			}
		}
	case *ast.RenameTableStmt:
		//"rename table ..."
		rulepkg.AddResult(input.Res, input.Rule, SQLE00074)
	}
	return nil
}

// ==== Rule code end ====
