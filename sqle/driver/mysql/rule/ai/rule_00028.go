package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00028 = "SQLE00028"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00028,
			Desc:       "除了自增列和blob/text类型的列, 每个列都必须添加默认值",
			Annotation: "列添加默认值，可避免列为NULL值时对查询的影响",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "除了自增列和blob/text类型的列, 每个列都必须添加默认值. 不符合规定的字段: %v",
		Func:    RuleSQLE00028,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule: "In table definition, Every column, except auto-incrementing columns and blob/text columns, should have a default value".
You should follow the following logic:
1. For "CREATE TABLE ... " statement, check every column, if its type is not "blob" or "text", and its type is not "serial"/"smallserial"/"bigserial", and it has no DEFAULT value defined, add the column name to violation-list.
2. For "ALTER TABLE ... ADD COLUMN ... " statement, check new column, if its type is not "blob" or "text", and its type is not "serial"/"smallserial"/"bigserial", and it has no DEFAULT value defined, add the column name to violation-list.
3. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00028(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnAutoIncrement(col) {
				//the column is auto-increment
				continue
			}

			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				//the column is blob/text type
				continue
			}

			if !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
				//the column has no "DEFAULT" value
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnAutoIncrement(col) {
					//the column is auto-increment
					continue
				}

				if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
					//the column is blob/text type
					continue
				}

				if !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
					//the column has no "DEFAULT" value
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00028, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
