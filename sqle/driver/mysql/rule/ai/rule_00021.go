package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00021 = "SQLE00021"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00021,
			Desc:       "对于MySQL的DDL, 表字段必须有NOT NULL约束",
			Annotation: "表字段必须有 NOT NULL 约束可确保数据的完整性，防止插入空值，提升查询准确性",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 表字段必须有NOT NULL约束. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00021,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00021): "In table definition, table fields must have a NOT NULL constraint".
You should follow the following logic:
1. For "create table" statement, check every column if it has NOT NULL constraint, otherwise, add the column name to violation-list
2. For "alter table add column" statement, check the column if it has NOT NULL constraint, otherwise, add the column name to violation-list
3. For "alter table modify column" statement, check the modified column definition if it has NOT NULL constraint, otherwise, add the column name to violation-list
4. For "alter table change column" statement, check the new column's definition if it has NOT NULL constraint, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00021(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
				continue
			}
			violateColumns = append(violateColumns, col)
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
					continue
				}
				violateColumns = append(violateColumns, col)
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00021, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
