package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00033 = "SQLE00033"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00033,
			Desc:       plocale.Rule00033Desc,
			Annotation: plocale.Rule00033Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID, plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelError,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "UPDATE_TIME",
				Desc:  plocale.Rule00033Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00033Message,
		Func:    RuleSQLE00033,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00033): "In DDL, when creating table, table should have a field about update-timestamp, whose DEFAULT value and ON-UPDATE value should be both 'CURRENT_TIMESTAMP'", the update-timestamp column name is a parameter whose default value is 'UPDATE_TIME'.
You should follow the following logic:
1. For "create table ..." statement, check the following conditions, report violation if any condition is violated:
  1. The table should have a update-timestamp column whose type is datetime or timestamp, and column name is same as the parameter
  2. The update-timestamp column's DEFAULT value should be configured as 'CURRENT_TIMESTAMP'
  3. The update-timestamp column's ON-UPDATE value should be configured as 'CURRENT_TIMESTAMP'
2. For "ALTER TABLE..." Statement, if the added field column name is the same as the parameter, checks the following conditions, report violation if any condition is violated:
  1. its data type should be datetime or timestamp.
  2. The default value for this update time column should be set to 'CURRENT_TIMESTAMP'
  3. The ON-UPDATE value of the update time column should be configured to 'CURRENT_TIMESTAMP'
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00033(input *rulepkg.RuleHandlerInput) error {
	// get expected update_time column name in config
	updateTimeColumnName := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String()
	found := false

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table"
		for _, col := range stmt.Cols {
			if strings.EqualFold(util.GetColumnName(col), updateTimeColumnName) {
				// the column is update_time column
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) || util.IsColumnTypeEqual(col, mysql.TypeDatetime) {
					// the column is Timestamp-type or DateTime-type
					if c := util.GetColumnOption(col, ast.ColumnOptionDefaultValue); nil != c && util.IsOptionFuncCall(c, "current_timestamp") {
						// the column has "DEFAULT" constraint, the "DEFAULT" value is current_timestamp
						if c := util.GetColumnOption(col, ast.ColumnOptionOnUpdate); nil != c && util.IsOptionFuncCall(c, "current_timestamp") {
							// the column has "ON UPDATE" constraint, the "DEFAULT" value is current_timestamp
							found = true
						}
					}
				}
			}
		}

		if !found {
			//the column is not created by "create table..."
			rulepkg.AddResult(input.Res, input.Rule, SQLE00033, updateTimeColumnName)
		}
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			// "alter table add column" or "alter table change column" or "alter table add column"
			for _, col := range spec.NewColumns {
				if strings.EqualFold(util.GetColumnName(col), updateTimeColumnName) {
					violated := true
					// the column is update_time column
					if util.IsColumnTypeEqual(col, mysql.TypeTimestamp, mysql.TypeDatetime) {
						// the column type is timestamp or datetime
						if util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
							// the column has "DEFAULT" option
							option := util.GetColumnOption(col, ast.ColumnOptionDefaultValue)
							if util.IsOptionFuncCall(option, "current_timestamp") {
								// the "DEFAULT" value is current_timestamp
								if util.IsColumnHasOption(col, ast.ColumnOptionOnUpdate) {
									// the column has "ON UPDATE" option
									option := util.GetColumnOption(col, ast.ColumnOptionOnUpdate)
									if util.IsOptionFuncCall(option, "current_timestamp") {
										// the "ON UPDATE" value is current_timestamp
										violated = false
									}
								}
							}
						}
					}
					if violated {
						// the column is not created by "alter table ..."
						rulepkg.AddResult(input.Res, input.Rule, SQLE00033, updateTimeColumnName)
					}
				}
			}

		}

	default:
		return nil
	}

	return nil
}

// ==== Rule code end ====
