package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00025 = "SQLE00025"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00025,
			Desc:       "TIMESTAMP 类型的列必须添加默认值",
			Annotation: "TIMESTAMP 类型的列添加默认值，可避免出现全为0的日期格式与业务预期不符",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "TIMESTAMP 类型的列必须添加默认值. 不符合规定的字段: %v",
		AllowOffline: false,
		Func:    RuleSQLE00025,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00025): "In table definition, time-type column must have a default value".
You should follow the following logic:
1. For "create table ..." statement, for every column whose type is time-type (Timestamp-type or Datetime-type), if the column has no DEFAULT value defined, add the column name to violation-list
2. For "alter table ... add column ..." statement, if the column type is time-type (Timestamp-type or Datetime-type), if the column has no DEFAULT value defined, add the column name to violation-list
3. For "alter table ... modify column ..." statement, if the column type is time-type (Timestamp-type or Datetime-type), if the column has no DEFAULT value defined, add the column name to violation-list
4. For "alter table ... change column ..." statement, if the new column type is time-type (Timestamp-type or Datetime-type), if the column has no DEFAULT value defined, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00025(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeTimestamp, mysql.TypeDatetime) {
				//the column type is timestamp or datetime
				if util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
					//the column has "DEFAULT" constraint
					continue
				}
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp, mysql.TypeDatetime) {
					//the column type is timestamp or datetime
					if util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
						//the column has "DEFAULT" constraint
						continue
					}
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00025, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
