package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00016_1 = "SQLE00016_1"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00016_1,
			Desc:       "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
			Annotation: "BLOB 和 TEXT 类型的字段无法指定默认值，如插入数据不指定字段默认为NULL，如果添加了 NOT NULL 限制，写入数据时又未对该字段指定值会导致写入失败",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL. 不符合规定的字段: %v",
		Func:    RuleSQLE00016_1,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00016_1): "In table definition, BLOB and TEXT fields cannot be set as NOT NULL".
You should follow the following logic:
1. For "create table ..." statement, for every column whose type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob), check if column constraints has no NOT-NULL; otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, for the column whose type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob), check if column constraints has no NOT-NULL; otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, if column type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob), check if column constraints has no NOT-NULL; otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, if column type is blob-type (including TinyBlob/MediumBlob/Blob/LongBlob) in the new column's definition, check if column constraints has no NOT-NULL; otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00016_1(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				//the column type is blob or text
				if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
					//the column type is blob or text
					if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
						violateColumns = append(violateColumns, col)
					}
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00016_1, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
