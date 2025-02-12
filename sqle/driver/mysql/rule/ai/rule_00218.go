package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00218 = "SQLE00218"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00218,
			Desc:       plocale.Rule00218Desc,
			Annotation: plocale.Rule00218Annotation,
			Category:   plocale.RuleTypeIndexInvalidation,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00218Message,
		Func:    RuleSQLE00218,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00218): "For dml, The leftmost field of the union index must appear within the query condition".
You should follow the following logic:
1. For SELECT... Statement, record the filter fields in the WHERE condition and the ON condition of the JOIN. When there is no WHERE condition or where condition is always True (for example, where 1=1 or where True), but there is grouping or sorting, note the grouped or sorted field as well. See if these fields are included in the union index list and are not the leftmost column of the union index, then the rule is violated. The acquisition of the table index needs to be obtained online.
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00218(input *rulepkg.RuleHandlerInput) error {

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

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt:
		// "SELECT...", "INSERT...", "UNION..."
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
			// get column names in join condition
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil && selectStmt.From.TableRefs.On != nil {
				for _, col := range util.GetColumnNameInExpr(selectStmt.From.TableRefs.On.Expr) {
					table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
				}
			}

			// get column names in where condition
			for _, col := range util.GetColumnNameInExpr(selectStmt.Where) {
				table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
			}

			if selectStmt.Where == nil || util.IsExprConstTrue(input.Ctx, selectStmt.Where, alias) {
				// get column names in group by when there is no where condition
				if selectStmt.GroupBy != nil {
					for _, item := range selectStmt.GroupBy.Items {
						for _, col := range util.GetColumnNameInExpr(item.Expr) {
							table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
						}
					}
				}
				// get column names in order by when there is no where condition
				if selectStmt.OrderBy != nil {
					for _, item := range selectStmt.OrderBy.Items {
						for _, col := range util.GetColumnNameInExpr(item.Expr) {
							table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
						}
					}
				}
			}

			for table, colNames := range table2colNames {
				// get index of the table
				indexesInfo, err := util.GetTableIndexes(input.Ctx, table, colNames[0].Schema.L)
				if err != nil {
					log.NewEntry().Errorf("get table indexes failed, sqle: %v, error: %v", input.Node.Text(), err)
					return nil
				}

			LOOP:
				for _, colName := range colNames {
					// check if the column in the union index and is not the leftmost column
					for _, cols := range indexesInfo {
						for i, col := range cols {
							if colName.Name.String() == col && i == 0 {
								continue LOOP
							}
						}
					}
					rulepkg.AddResult(input.Res, input.Rule, SQLE00218, colName.Name.String())
				}
			}
		}

	}
	return nil
}

// ==== Rule code end ====
