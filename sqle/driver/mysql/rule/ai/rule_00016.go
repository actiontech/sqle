package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00016 = "SQLE00016"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00016,
			Desc:       "对于MySQL的DDL, BLOB 和 TEXT 类型的字段如果定义了默认值, 那默认值应为NULL",
			Annotation: "在SQL_MODE严格模式下BLOB 和 TEXT 类型无法设置默认值，如插入数据不指定值，字段会被设置为NULL",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, BLOB 和 TEXT 类型的字段如果定义了默认值, 那默认值应为NULL. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00016,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00016): "In table definition, the default value for BLOB-type and TEXT-type columns, if defined, should be NULL".
You should follow the following logic:
1. For "create table ..." statement, for every column whose type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob) and whose DEFAULT value is defined, check if the DEFAULT value is NULL; otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, for the column whose type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob) and whose DEFAULT value is defined, check if the DEFAULT value is NULL; otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, if column type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob) and the column DEFAULT value is defined, check if the DEFAULT value is NULL; otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, if column type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob) and the column DEFAULT value is defined in the new column's definition, check if the DEFAULT value is NULL; otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00016(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				//the column type is blob or text
				if util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
					//the column has "DEFAULT" option
					option := util.GetColumnOption(col, ast.ColumnOptionDefaultValue)

					//the "DEFAULT" value is not NULL
					if !util.IsOptionValIsNull(option) {
						violateColumns = append(violateColumns, col)
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
					//the column type is blob or text
					if util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
						//the column has "DEFAULT" option
						option := util.GetColumnOption(col, ast.ColumnOptionDefaultValue)

						//the "DEFAULT" value is not NULL
						if !util.IsOptionValIsNull(option) {
							violateColumns = append(violateColumns, col)
						}
					}
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00016, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
