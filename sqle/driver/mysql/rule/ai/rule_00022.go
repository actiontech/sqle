package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00022 = "SQLE00022"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00022,
			Desc:       plocale.Rule00022Desc,
			Annotation: plocale.Rule00022Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelNotice,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "0.7",
				Desc:  plocale.Rule00022Params1,
				Type:  params.ParamTypeFloat64,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00022Message,
		Func:    RuleSQLE00022,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

// ==== Rule code start ====
func RuleSQLE00022(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold := param.Float64()
	if threshold > 1 || threshold <= 0 {
		return fmt.Errorf("param %s should be in range (0, 1]", rulepkg.DefaultSingleParamKeyName)
	}

	var tableName *ast.TableName
	indexColumns := make([]string, 0)
	switch stmt := input.Node.(type) {
	case *ast.CreateIndexStmt:
		// "create index..."

		tableName = stmt.Table
		for _, col := range stmt.IndexPartSpecifications {
			//"create index... column..."
			indexColumns = append(indexColumns, util.GetIndexColName(col))
		}

		if len(indexColumns) == 0 {
			// the table has no index
			return nil
		}

		skewness, err := util.CalculateIndexSkewness(input.Ctx, tableName, indexColumns)
		if err != nil {
			log.NewEntry().Errorf("get index skewness failed, sqle: %v, error: %v", input.Node.Text(), err)
			return nil
		}

		// the table has index, check the skewness
		for col, d := range skewness {
			if d > threshold {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00022, col, threshold)
			}
		}
	case *ast.AlterTableStmt:
		// "alter table"
		tableName = stmt.Table

		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			constraints := util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...)

			for _, constraint := range constraints {
				for _, col := range constraint.Keys {
					indexColumns = append(indexColumns, util.GetIndexColName(col))
				}
			}

		}
		if len(indexColumns) == 0 {
			// the table has no index
			return nil
		}

		skewness, err := util.CalculateIndexSkewness(input.Ctx, tableName, indexColumns)
		if err != nil {
			log.NewEntry().Errorf("get index skewness failed, sqle: %v, error: %v", input.Node.Text(), err)
			return nil
		}

		// the table has index, check the skewness
		for col, d := range skewness {
			if d > threshold {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00022, col, threshold)
			}
		}
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
		// "select..."

		var defaultTable string
		var alias []*util.TableAliasInfo
		getTableName := func(col *ast.ColumnNameExpr) string {
			if col.Name.Table.L != "" {
				for _, a := range alias {
					if a.TableAliasName == col.Name.Table.String() {
						return a.TableName
					}
				}
				return col.Name.Table.L
			}
			return defaultTable
		}

		for _, selectStmt := range util.GetSelectStmt(stmt) {

			// get default table name
			if t := util.GetDefaultTable(selectStmt); t != nil {
				defaultTable = t.Name.O
			}

			// get table alias info
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				alias = util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)
			}

			var (
				table2colNames = map[string] /*table name*/ []*ast.ColumnName /*col names*/ {}
			)

			// get column names in where condition
			for _, col := range util.GetColumnNameInExpr(selectStmt.Where) {
				table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
			}

			for table, colNames := range table2colNames {
				// get index of the table
				indexesInfo, err := util.GetTableIndexes(input.Ctx, table, colNames[0].Schema.L)
				if err != nil {
					log.NewEntry().Errorf("get table indexes failed, sqle: %v, error: %v", input.Node.Text(), err)
					return nil
				}

				// get index columns
				indexColumns := make([]string, 0)
				for _, colName := range colNames {

					for _, cols := range indexesInfo {
						// check if the column is index
						for _, col := range cols {
							if colName.Name.String() == col {
								indexColumns = append(indexColumns, colName.Name.String())
							}
						}
					}
				}

				if len(indexColumns) == 0 {
					// no index
					return nil
				}

				tableName = &ast.TableName{
					Schema: colNames[0].Schema,
					Name:   model.NewCIStr(table),
				}
				skewness, err := util.CalculateIndexSkewness(input.Ctx, tableName, indexColumns)
				if err != nil {
					log.NewEntry().Errorf("get index skewness failed, sqle: %v, error: %v", input.Node.Text(), err)
					return nil
				}

				// has index, check the skewness
				for col, d := range skewness {
					if d > threshold {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00022, col, threshold)
					}
				}

			}
		}

	}

	return nil
}

// ==== Rule code end ====
