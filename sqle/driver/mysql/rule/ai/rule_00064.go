package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00064 = "SQLE00064"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00064,
			Desc:       plocale.Rule00064Desc,
			Annotation: plocale.Rule00064Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID, plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "767",
				Desc:  plocale.Rule00064Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00064Message,
		Func:    RuleSQLE00064,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00064): "For MySQL DDL, it is not recommended that the length of the index field is greater than the threshold when it is of type VARCHAR".the threshold should be a parameter whose default value is 767.
You should follow the following logic:
1. For "CREATE TABLE..." Statement, if a defining index operation (such as key or index) exists, checks the following:
    1. Define a collection
    2. Get the index name of the table from the column definition and table constraint in the table construction clause, and store the index name into the collection
    3. Iterate over all field definitions in the list clause, and if the field name is in the set, determine whether the field is of type varchar, and the varchar size is greater than or equal to the threshold, then report the rule violation
2. For "CREATE INDEX..." The statement,
    1. Define a collection
    2. Put the names of all indexed fields into the set
    3. Get the create table sentence
    4. Iterate over all field definitions in the list clause, and if the field name is in the set, determine whether the field is of type varchar, and the varchar size is greater than or equal to the threshold, then report the rule violation
3. For "ALTER TABLE... ADD INDEX..." Statement to perform the same checks as above.
==== Prompt end ====
*/

// ==== Rule code start ====

func RuleSQLE00064(input *rulepkg.RuleHandlerInput) error {
	// get expected length of varchar
	expectedLength := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()

	var indexColName []string
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."

		// count index column in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) || util.IsColumnPrimaryKey(col) {
				indexColName = append(indexColName, util.GetColumnName(col))
			}
		}

		constraints := util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...)
		// count index column in table constraint
		for _, constraint := range constraints {
			for _, key := range constraint.Keys {
				indexColName = append(indexColName, util.GetIndexColName(key))
			}
		}

		// check if the column type is varchar and the length of varchar is greater than expectedLength
		for _, col := range stmt.Cols {
			if util.IsStrInSlice(util.GetColumnName(col), indexColName) {
				// the column is index column
				if util.IsColumnTypeEqual(col, mysql.TypeVarchar) {
					// the column type is varchar
					if util.GetColumnWidth(col) >= expectedLength {
						// the column width exceeds the expected length
						rulepkg.AddResult(input.Res, input.Rule, SQLE00064, util.GetColumnName(col))
						return nil
					}
				}
			}
		}

	case *ast.CreateIndexStmt:
		// "create index..."

		var indexColName []string
		// get index column in "create index ..."
		for _, col := range stmt.IndexPartSpecifications {
			indexColName = append(indexColName, util.GetIndexColName(col))
		}

		// get create table stmt in context
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// check if the column type is varchar and the length of varchar is greater than expectedLength
		for _, col := range createTableStmt.Cols {
			if util.IsStrInSlice(util.GetColumnName(col), indexColName) {
				// the column is index column
				if util.IsColumnTypeEqual(col, mysql.TypeVarchar) {
					// the column type is varchar
					if util.GetColumnWidth(col) >= expectedLength {
						// the column width exceeds the expected length
						rulepkg.AddResult(input.Res, input.Rule, SQLE00064, util.GetColumnName(col))
						return nil
					}
				}
			}
		}

	case *ast.AlterTableStmt:
		// "alter table"

		var indexColName []string
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// get index column in alter table command
			for _, constraint := range util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...) {
				for _, key := range constraint.Keys {
					indexColName = append(indexColName, util.GetIndexColName(key))
				}
			}
		}

		// get create table stmt in context
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// check if the column type is varchar and the length of varchar is greater than expectedLength
		for _, col := range createTableStmt.Cols {
			if util.IsStrInSlice(util.GetColumnName(col), indexColName) {
				// the column is index column
				if util.IsColumnTypeEqual(col, mysql.TypeVarchar) {
					// the column type is varchar
					if util.GetColumnWidth(col) >= expectedLength {
						// the column width exceeds the expected length
						rulepkg.AddResult(input.Res, input.Rule, SQLE00064, util.GetColumnName(col))
						return nil
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
