package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00027 = "SQLE00027"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00027,
			Desc:       "对于MySQL的DDL, 列定义建议加注释",
			Annotation: "MySQL列添加注释能够使列的意义更明确，方便日后的维护。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 列定义建议加注释. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00027,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00027): "In DDL, when creating table, it is recommended to add comments to column definitions."
You should follow the following logic:
1. For "create table ..." statement, check every column, if its option has no comment, add the column name to violation-list
2. For "create table ..." statement, check every column, if its comment option has only spaces or empty strings, add the column name to violation-list
3. For "alter table ... add column ..." statement, check the column, if its option has no comment, add the column name to violation-list
4. For "alter table ... add column ..." statement, check the column, if its comment option has only spaces or empty strings, add the column name to violation-list
5. For "alter table ... modify column ..." statement, check the modified column definition, if its option has no comment, add the column name to violation-list
6. For "alter table ... modify column ..." statement, check the modified column definition, if its comment option has only spaces or empty strings, add the column name to violation-list
7. For "alter table ... change column ..." statement, check the new column's definition, if its option has no comment, add the column name to violation-list
8. For "alter table ... change column ..." statement, check the new column's definition, if its comment option has only spaces or empty strings, add the column name to violation-list
9. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00027(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			//if the column has "COMMENT" option and its comment is not only spaces or empty strings, ignore it
			if c := util.GetColumnOption(col, ast.ColumnOptionComment); nil != c && len(strings.TrimSpace(util.GetValueExprStr(c.Expr))) > 0 {
				continue
			}

			violateColumns = append(violateColumns, col)
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				//if the column has "COMMENT" option and its comment is not only spaces or empty strings, ignore it
				if c := util.GetColumnOption(col, ast.ColumnOptionComment); nil != c && len(strings.TrimSpace(util.GetValueExprStr(c.Expr))) > 0 {
					continue
				}

				violateColumns = append(violateColumns, col)
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00027, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====