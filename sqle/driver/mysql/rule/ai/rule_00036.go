package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00036 = "SQLE00036"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00036,
			Desc:         plocale.Rule00036Desc,
			Annotation:   plocale.Rule00036Annotation,
			Category:     plocale.RuleTypeIndexingConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00036Message,
		Func:    RuleSQLE00036,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00036): "For index, adding BLOB columns to the index is prohibited".
You should follow the following logic:
1. For the "CREATE TABLE . .." statements, check if the index which on the table option or on the column includes blob columns. If it does, report a violation.
2. For the  "CREATE INDEX ..." statements, check if the index includes blob columns. If it does, report a violation.
3. For the  "ALTER TABLE ... ADD INDEX ..." statements, check if the index includes blob columns. If it does, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00036(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:

		// get blob columns
		blobColNames := []string{}
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				blobColNames = append(blobColNames, util.GetColumnName(col))
			}
		}

		// check index in column definition
		for _, col := range stmt.Cols {
			// "create table ... "
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) || util.IsColumnPrimaryKey(col) {
				if util.IsStrInSlice(util.GetColumnName(col), blobColNames) {
					//the column is blob type
					rulepkg.AddResult(input.Res, input.Rule, SQLE00036)
					return nil
				}
			}
		}

		// check index in table constraint
		constraints := util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...)
		for _, constraint := range constraints {
			for _, col := range constraint.Keys {
				if util.IsStrInSlice(util.GetIndexColName(col), blobColNames) {
					//the column is blob type
					rulepkg.AddResult(input.Res, input.Rule, SQLE00036)
					return nil
				}
			}
		}

	case *ast.CreateIndexStmt:
		// "create index..."

		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// get blob columns
		blobColNames := []string{}
		for _, col := range createTableStmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				blobColNames = append(blobColNames, util.GetColumnName(col))
			}
		}

		for _, col := range stmt.IndexPartSpecifications {
			//"create index... column..."
			if util.IsStrInSlice(util.GetIndexColName(col), blobColNames) {
				//the column is blob type
				rulepkg.AddResult(input.Res, input.Rule, SQLE00036)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		// get blob columns
		blobColNames := []string{}
		for _, col := range createTableStmt.Cols {
			if util.IsColumnTypeEqual(col, util.GetBlobDbTypes()...) {
				blobColNames = append(blobColNames, util.GetColumnName(col))
			}
		}

		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			constraints := util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...)

			for _, constraint := range constraints {
				for _, col := range constraint.Keys {
					if util.IsStrInSlice(util.GetIndexColName(col), blobColNames) {
						//the column is blob type
						rulepkg.AddResult(input.Res, input.Rule, SQLE00036)
						return nil
					}
				}
			}

		}
	}
	return nil
}

// ==== Rule code end ====
