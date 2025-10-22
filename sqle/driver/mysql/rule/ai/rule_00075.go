package ai

import (
	"strings"

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
In MySQL, you should check if the SQL violate the rule(SQLE00075): "Compare whether the character set of the specified table and columns in the statement are consistent".
You should follow the following logic:

1. For "create table ..." statement, check if the table's charset is consistent with each column's charset. If a column has specified charset and it's different from table's charset, add the column name to violation-list
2. For "alter table ... add column ..." statement, check if the table's charset is consistent with the new column's charset. If the column has specified charset and it's different from table's charset, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check if the table's charset is consistent with the modified column's charset. If the column has specified charset and it's different from table's charset, add the column name to violation-list
4. For "alter table ... change column ..." statement, check if the table's charset is consistent with the new column's charset. If the column has specified charset and it's different from table's charset, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00075(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	var tableCharset string
	var err error

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// Get table charset from CREATE TABLE statement
		if charsetOption := util.GetTableOption(stmt.Options, ast.TableOptionCharset); charsetOption != nil {
			tableCharset = charsetOption.StrValue
		} else {
			// If no table charset specified, get from schema default
			tableCharset, err = input.Ctx.GetSchemaCharacter(stmt.Table, "")
			if err != nil {
				return err
			}
		}

		// Check each column's charset against table charset
		for _, col := range stmt.Cols {
			if util.IsColumnHasSpecifiedCharset(col) {
				columnCharset := col.Tp.Charset
				if !strings.EqualFold(columnCharset, tableCharset) {
					violateColumns = append(violateColumns, col)
				}
			}
		}

	case *ast.AlterTableStmt:
		// Get table charset from ALTER TABLE statement
		tableCharset, err = input.Ctx.GetSchemaCharacter(stmt.Table, "")
		if err != nil {
			return err
		}

		// Check if table charset is being changed in this ALTER statement
		for _, spec := range stmt.Specs {
			// Use the last defined table character set
			if spec.Tp == ast.AlterTableOption {
				for _, option := range spec.Options {
					if option.Tp == ast.TableOptionCharset {
						tableCharset = option.StrValue
						break
					}
				}
			}
		}

		// Check columns in ALTER TABLE commands
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnHasSpecifiedCharset(col) {
					columnCharset := col.Tp.Charset
					if !strings.EqualFold(columnCharset, tableCharset) {
						violateColumns = append(violateColumns, col)
					}
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
