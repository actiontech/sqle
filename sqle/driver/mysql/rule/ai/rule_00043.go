package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00043 = "SQLE00043"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00043,
			Desc:       plocale.Rule00043Desc,
			Annotation: plocale.Rule00043Annotation,
			Category:   plocale.RuleTypeIndexingConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID, plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "2",
				Desc:  plocale.Rule00043Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00043Message,
		Func:    RuleSQLE00043,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00043): "For table creation and index creation statements, the number of index on the same column should be within threshold, the threshold should be a parameter whose default value is 2.".
You should follow the following logic:
1. For the "CREATE TABLE ..." statements, count the number of times each column appears in the index，check if max number is more than 2. If it does, report a violation.
2. For the  "CREATE INDEX ..." statements, count the number of times each column appears in the index，check if max number is more than 2. If it does, report a violation.
3. For the  "ALTER TABLE ... ADD INDEX ..." statements, count the number of times each column appears in the index，check if max number is more than 2. If it does, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00043(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxIndexCount := param.Int()

	indexCounter := map[string] /*col name*/ int{}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."

		// count index column in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) || util.IsColumnPrimaryKey(col) {
				indexCounter[util.GetColumnName(col)]++
			}
		}

		constraints := util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...)
		// count index column in table constraint
		for _, constraint := range constraints {
			for _, key := range constraint.Keys {
				indexCounter[util.GetIndexColName(key)]++
			}
		}

	case *ast.CreateIndexStmt:
		// "create index..."

		// get create table stmt in context
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// count existed index column in table constraint
		for _, constraint := range util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...) {
			for _, key := range constraint.Keys {
				indexCounter[util.GetIndexColName(key)]++
			}
		}

		// count index in "create index ..."
		for _, col := range stmt.IndexPartSpecifications {
			//"create index... column..."
			indexCounter[util.GetIndexColName(col)]++
		}
	case *ast.AlterTableStmt:
		// "alter table"

		// get create table stmt in context
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// count existed index column in table constraint
		for _, constraint := range util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...) {
			for _, key := range constraint.Keys {
				indexCounter[util.GetIndexColName(key)]++
			}
		}

		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// count index column in table constraint
			for _, constraint := range util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...) {
				for _, key := range constraint.Keys {
					indexCounter[util.GetIndexColName(key)]++
				}
			}
		}

	}
	// check if the column counter in index is more than param
	for col, count := range indexCounter {
		if count > maxIndexCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00043, col, maxIndexCount)
		}
	}
	return nil
}

// ==== Rule code end ====
