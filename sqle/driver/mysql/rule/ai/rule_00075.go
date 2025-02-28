package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00075 = "SQLE00075"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00075,
			Desc:       plocale.Rule00075Desc,
			Annotation: plocale.Rule00075Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID, plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagSQLTablespace.ID, plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagCorrection.ID, plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00075Message,
		Func:    RuleSQLE00075,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00075): "In table definition, setting column-level specified charset or collation is prohibited".
You should follow the following logic:

1. For "create table ..." statement, check every column, if it has no charset setting and no collation setting, otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, if it has no charset setting and no collation setting, otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, if it has no charset setting and no collation setting, otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, if it has no charset setting and no collation setting, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00075(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		//"create table ..."
		for _, col := range stmt.Cols {
			//if the column has "CHARSET" or "COLLATE" specified, it is violate the rule
			if util.IsColumnHasSpecifiedCharset(col) || util.IsColumnHasOption(col, ast.ColumnOptionCollate) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			// "alter table ... add column ..." or "alter table ... modify column ..." or "alter table ... change column ..."
			for _, col := range spec.NewColumns {
				//if the column has "CHARSET" or "COLLATE" specified, it is violate the rule
				if util.IsColumnHasSpecifiedCharset(col) || util.IsColumnHasOption(col, ast.ColumnOptionCollate) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00075, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
