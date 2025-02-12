package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00025 = "SQLE00025"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00025,
			Desc:       plocale.Rule00025Desc,
			Annotation: plocale.Rule00025Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID, plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00025Message,
		Func:    RuleSQLE00025,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
