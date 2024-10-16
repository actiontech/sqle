package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
)

const (
	SQLE00086 = "SQLE00086"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00086,
			Desc:       "在 MySQL 中, 禁止使用子字符串匹配或后缀匹配搜索",
			Annotation: "使用子字符串匹配搜索或后缀匹配搜索将导致查询无法使用索引，导致全表扫描",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止使用子字符串匹配或后缀匹配搜索",
		AllowOffline: true,
		Func:         RuleSQLE00086,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00086): "在 MySQL 中，禁止使用子字符串匹配或后缀匹配搜索."
您应遵循以下逻辑：
1. 对于"SELECT..."语句，检查SQL语句， 如果下面任意一项为真，则报告违反规则。
  1. 存在 like '%ab'或'_ab' 这样的后缀匹配模糊检索。
  2. 存在 like '%ab%'或'_ab_' 这样的子字符串匹配模糊检索。
2. 对于"INSERT...SELECT..."语句，递归检查所有嵌套的SELECT语句，执行与上面类似的检查。
3. 对于"UPDATE..."语句，递归检查所有嵌套的SELECT语句，执行与上面类似的检查。
4. 对于"DELETE..."语句，递归检查所有嵌套的SELECT语句，执行与上面类似的检查。
5. 对于"... UNION ALL ..."语句，递归检查所有嵌套的SELECT语句，执行与上面类似的检查。
6. 对于"WITH..."语句，递归检查所有嵌套的SELECT语句，执行与上面类似的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00086(input *rulepkg.RuleHandlerInput) error {
	whereList := util.GetWhereExprFromDMLStmt(input.Node)
	negative := false
	util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.PatternLikeExpr:
			switch pattern := x.Pattern.(type) {
			case *parserdriver.ValueExpr:
				datum := pattern.Datum.GetString()
				if strings.HasPrefix(datum, "%") ||
					strings.HasPrefix(datum, "_") {
					negative = true
					return true
				}
			}
		}
		return false
	}, whereList...)

	// 对于"WITH..."语句
	// TODO 待实现

	if negative {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00086)
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
