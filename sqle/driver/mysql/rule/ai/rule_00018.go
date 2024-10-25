package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00018 = "SQLE00018"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00018,
			Desc:       "对于MySQL的DDL, CHAR长度大于20时，必须使用VARCHAR类型",
			Annotation: "VARCHAR是变长字段，存储空间小，可节省存储空间，同时相对较小的字段检索效率显然也要高些",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, CHAR长度大于20时，必须使用VARCHAR类型. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00018,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00018): "In table definition, when the length of CHAR is greater than 20, VARCHAR type must be used".
You should follow the following logic:
1. For "create table ..." statement, for every column whose type is string-type and length is larger than 20, add the column name to violation-list
2. For "alter table ... add column ..." statement, if the column type is string-type and length is larger than 20, add the column name to violation-list
3. For "alter table ... modify column ..." statement, if the column type is string-type and length is larger than 20, add the column name to violation-list
4. For "alter table ... change column ..." statement, if the new column type is string-type and length is larger than 20, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00018(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeString) && util.GetColumnWidth(col) > 20 {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeString) && util.GetColumnWidth(col) > 20 {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00018, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
