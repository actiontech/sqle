package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00068 = "SQLE00068"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00068,
			Desc:       "对于MySQL的DDL, 禁止使用TIMESTAMP类型字段",
			Annotation: "TIMESTAMP类型字段 有最大值限制（'2038-01-19 03:14:07' UTC），且会时区转换的问题",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 禁止使用TIMESTAMP类型字段. 不符合规则的字段有: %v",
		Func:    RuleSQLE00068,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00068): "In table definition, timestamp-type column is prohibited".
You should follow the following logic:
1. For "create table ..." statement, check every column, if its type is other than TIMESTAMP, otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, if its type is other than TIMESTAMP, otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, if its type is other than TIMESTAMP, otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, if its type is other than TIMESTAMP, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00068(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		//"create table ..."
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			// "alter table ... add column ..." or "alter table ... modify column ..." or "alter table ... change column ..."
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00068, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====