package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00064 = "SQLE00064"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00064,
			Desc:       "对于MySQL的DDL, 定义VARCHAR 长度时不建议大于阈值",
			Annotation: "MySQL建立索引时没有限制索引的大小，索引长度会默认采用的该字段的长度，VARCHAR 定义长度越长建立的索引存储大小越大；具体规则阈值可以根据业务需求调整，默认值：1024",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "1024",
					Desc:  "VARCHAR最大长度",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "对于MySQL的DDL, 定义VARCHAR 长度时不建议大于阈值. 阈值: %v",
		Func:    RuleSQLE00064,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00064): "In table definition, the length of VARCHAR column should be kept within the threshold", the threshold should be a parameter whose default value is 1024.
You should follow the following logic:

1. For "create table ..." statement, check every column whose type is VARCHAR if it's type length is within threshold, otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, when the type is VARCHAR, if it's type length is within threshold, otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, when the type is VARCHAR, if it's type length is within threshold, otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, when the type is VARCHAR, if it's type length is within threshold, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00064(input *rulepkg.RuleHandlerInput) error {
	// get expected length of VARCHAR
	expectedLength := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeVarchar) {
				//the column type is varchar
				if util.GetColumnWidth(col) > expectedLength {
					// the column width exceeds the expected length
					rulepkg.AddResult(input.Res, input.Rule, SQLE00064, expectedLength)
					return nil
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			//"alter table ... add column ..." or "alter table ... change column ..." or "alter table ... add column ..."
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeVarchar) {
					//the column type is varchar
					if util.GetColumnWidth(col) > expectedLength {
						// the column width exceeds the expected length
						rulepkg.AddResult(input.Res, input.Rule, SQLE00064, expectedLength)
						return nil
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
