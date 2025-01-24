package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00017 = "SQLE00017"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00017,
			Desc:         plocale.Rule00017Desc,
			Annotation:   plocale.Rule00017Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00017Message,
		Func:    RuleSQLE00017,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00017): "In table definition, blob-type and text-type columns are prohibited".
You should follow the following logic:
1. For "create table ..." statement, check every column, if its type is not BLOB nor TEXT (including TinyBlob/MediumBlob/Blob/LongBlob), otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, if its type is not BLOB nor TEXT (including TinyBlob/MediumBlob/Blob/LongBlob), otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, if its type is not BLOB nor TEXT (including TinyBlob/MediumBlob/Blob/LongBlob), otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, if its type is not BLOB nor TEXT (including TinyBlob/MediumBlob/Blob/LongBlob), otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00017(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00017, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
