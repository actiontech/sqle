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
	SQLE00220 = "SQLE00220"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00220,
			Desc:       "在 MySQL 中, 避免不带where条件的count(*)或者count(1)",
			Annotation: "不带 where 条件的 count(*) 或者 count(1) 都是对表进行暴力扫描，极其耗费系统资源",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免不带where条件的count(*)或者count(1)",
		AllowOffline: true,
		Func:         RuleSQLE00220,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00220): "在 MySQL 中，避免不带where条件的count(*)或者count(1)."
您应遵循以下逻辑：
1. 对于所有 DQL 语句，检查是否存在 count(*) 或 count(1)。
2. 如果存在 count(*) 或 count(1)，进一步检查该语句是否不带 WHERE 条件。需要考虑 WHERE 条件可能出现在不同位置，如 SELECT 查询中（包括子查询）、UPDATE 语句中、DELETE 语句中。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00220(input *rulepkg.RuleHandlerInput) error {
	// Helper function to determine if any field contains count(*) or count(1)
	hasCountStarOrOne := func(fields []*ast.SelectField) bool {
		for _, field := range fields {
			funcCall, ok := field.Expr.(*ast.AggregateFuncExpr)
			if !ok {
				continue
			}
			// 检查函数名是否为 "count"
			if !strings.EqualFold(funcCall.F, "count") {
				continue
			}
			// 检查 count 函数的参数是否为 * 或 1
			if len(funcCall.Args) != 1 {
				continue
			}

			if _, ok := funcCall.Args[0].(*parserdriver.ValueExpr); ok {
				return true
			}
		}
		return false
	}

	// 处理 SELECT 语句
	for _, selectStmt := range util.GetSelectStmt(input.Node) {
		if hasCountStarOrOne(selectStmt.Fields.Fields) {
			if selectStmt.Where == nil {
				// 违反规则: 存在 count(*) 或 count(1) 且没有 WHERE 条件或 WHERE 恒真
				rulepkg.AddResult(input.Res, input.Rule, SQLE00220)
			}
		}
	}
	// TODO 解析器不支持WITH（CTE)语法
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
