package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00016_1 = "SQLE00016_1"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00016_1,
			Desc:         plocale.Rule00016_1Desc,
			Annotation:   plocale.Rule00016_1Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00016_1Message,
		Func:    RuleSQLE00016_1,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
