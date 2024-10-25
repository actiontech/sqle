package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00131 = "SQLE00131"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00131,
			Desc:       "对于MySQL的DML, 避免使用 ORDER BY RAND() 进行随机排序",
			Annotation: "使用 ORDER BY RAND() 会导致 MySQL 生成临时表并进行完整的表扫描和排序，这在处理大数据量时会显著增加查询时间和服务器负载。建议采用更高效的随机数据检索方法，如利用主键或其他索引实现快速随机访问。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 避免使用 ORDER BY RAND() 进行随机排序",
		AllowOffline: true,
		Func:    RuleSQLE00131,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00131): "For dml, using ORDER BY RAND() is prohibited".
You should follow the following logic:
1. For SELECT... Order BY... Statement, checks whether the Order By clause contains a RAND function, and if so, report a violation.
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00131(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt:
		// "SELECT...", "INSERT...", "UNION..."
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// Apply the check to each order by items
			orderBy := selectStmt.OrderBy
			if orderBy != nil {
				for _, item := range orderBy.Items {
					if item.Expr == nil {
						continue
					}

					// Check for function calls in the expr.
					funcNames := util.GetFuncName(item.Expr)
					for _, name := range funcNames {
						// Check if the function is the one we want to check.
						if name == "rand" {
							// Add a rule violation result if an issue is found.
							rulepkg.AddResult(input.Res, input.Rule, SQLE00131)
						}
					}
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
