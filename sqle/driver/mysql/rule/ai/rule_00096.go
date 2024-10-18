package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00096 = "SQLE00096"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00096,
			Desc:       "对于MySQL的DML, 联合索引最左侧的字段必须出现在查询条件内",
			Annotation: "当查询条件包含联合索引的最左侧字段时，查询语句才能更好的利用索引的特性：有序性、过滤性等",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeIndexInvalidation,
		},
		Message: "对于MySQL的DML, 联合索引最左侧的字段必须出现在查询条件内. 不符合规范的字段: %v",
		Func:    RuleSQLE00096,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00096): "For dml, The leftmost field of the union index must appear within the query condition".
You should follow the following logic:
1. For SELECT... Statement, record the filter fields in the WHERE condition and the ON condition of the JOIN. When there is no WHERE condition or where condition is always True (for example, where 1=1 or where True), but there is grouping or sorting, note the grouped or sorted field as well. See if these fields are included in the union index list and are not the leftmost column of the union index, then the rule is violated. The acquisition of the table index needs to be obtained online.
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00096(input *rulepkg.RuleHandlerInput) error {

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

			// get column names in group by when there is no where condition
			if selectStmt.Where == nil || util.IsExprConstTrue(selectStmt.Where) {
				if selectStmt.GroupBy != nil {
					for _, item := range selectStmt.GroupBy.Items {
						for _, col := range util.GetColumnNameInExpr(item.Expr) {
							table2colNames[getTableName(col)] = append(table2colNames[getTableName(col)], col.Name)
						}
					}
				}
			}

			// get column names in order by when there is no where condition
			if selectStmt.Where == nil || util.IsExprConstTrue(selectStmt.Where) {
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

				for _, colName := range colNames {
					for _, index := range indexesInfo {
						// check if the column in the union index and is not the leftmost column
						if colName.Name.String() == index.ColumnName && index.SeqInIndex != "1" {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00096, colName.String())
						}
					}
				}
			}
		}

	}
	return nil
}

// ==== Rule code end ====
