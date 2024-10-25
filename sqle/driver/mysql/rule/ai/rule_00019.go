package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00019 = "SQLE00019"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00019,
			Desc:       "对于MySQL的DDL, 不建议使用 ENUM 类型",
			Annotation: "ENUM类型不是SQL标准，移植性较差，后期如修改或增加枚举值需重建整张表，代价较大，且无法通过字面量值进行排序",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 不建议使用 ENUM 类型. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00019,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00019): "In table definition, ENUM-type columns are prohibited".
You should follow the following logic:
1. For "create table ..." statement, check every column, if its type is not ENUM-type, otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, if its type is not ENUM-type, otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, if its type is not ENUM-type, otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, if its type is not ENUM-type, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00019(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeEnum) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeEnum) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00019, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
