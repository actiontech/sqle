package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00120 = "SQLE00120"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00120,
			Desc:       "在 MySQL 中, 避免使用 IN (NULL) 或者 NOT IN (NULL)",
			Annotation: "使用 `IN(NULL)` 或 `NOT IN(NULL)` 会导致查询条件永远为假，从而使得查询无法返回任何结果。这不仅影响查询逻辑和结果的准确性，还可能导致性能问题和不必要的资源消耗。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免使用 IN (NULL) 或者 NOT IN (NULL)",
		AllowOffline: true,
		Func:         RuleSQLE00120,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00120): "在 MySQL 中，避免使用 IN (NULL) 或者 NOT IN (NULL)."
您应遵循以下逻辑：
1. 针对所有 DML 语句（包括 SELECT、UPDATE、DELETE、UNION、INSERT...SELECT、WITH），执行以下检查：
   1. 在 WHERE 子句的语法节点中查找 IN (NULL) 的使用。
   2. 在 WHERE 子句的语法节点中查找 NOT IN (NULL) 的使用。
   如果发现上述任一情况，立即报告为规则违规。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00120(input *rulepkg.RuleHandlerInput) error {
	// 获取所有 DML 语句的 WHERE 子句
	whereList := util.GetWhereExprFromDMLStmt(input.Node)

	containsNull := func(exprList []ast.ExprNode) bool {
		for _, expr := range exprList {
			if expr.GetType().Tp == mysql.TypeNull {
				return true
			}
		}
		return false
	}

	// 遍历每个 WHERE 条件表达式
	for _, whereExpr := range whereList {
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternInExpr:
				// 检查 IN、NOT IN 子句中的值是否包含 NULL
				if containsNull(x.List) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00120)
					return true
				}
			}
			return false
		}, whereExpr)
	}
	// TODO WHIT语法

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
