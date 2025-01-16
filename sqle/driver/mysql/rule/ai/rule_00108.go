package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00108 = "SQLE00108"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00108,
			Desc:       "避免子查询嵌套层数过多",
			Annotation: "子查询嵌套层数超过阈值，有些情况下，子查询并不能使用到索引。同时对于返回结果集比较大的子查询，会产生大量的临时表，消耗过多的CPU和IO资源，产生大量的慢查询",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "子查询嵌套层数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "避免子查询嵌套层数过多",
		AllowOffline: true,
		Func:         RuleSQLE00108,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00108): "在 MySQL 中，避免子查询嵌套层数过多.默认参数描述: 子查询嵌套层数, 默认参数值: 5"
您应遵循以下逻辑：
1. 针对所有DML语句：
   1. 识别并统计子查询的层数。
      子查询的定义：SQL语句中包含另一个SELECT查询。
      常见的嵌套子查询形式：
      1）WHERE子句中的子查询
      2）FROM子句中的子查询
      3）SELECT列表中的子查询
      4）GROUP BY子句中的子查询
      5）HAVING子句中的子查询
      6）ORDER BY子句中的子查询
      7）JOIN的ON条件中的子查询
      8）INSERT、UPDATE、DELETE语句中的子查询
      嵌套层数示例：SELECT... WHERE column IN (SELECT... WHERE column2 IN (SELECT ...))
      上述示例中有2层嵌套子查询，SELECT可替换为任何DML语句，IN可替换为=、ANY等。
   2. 如果嵌套层数超过设定阈值，报告规则违规。

2. 针对UNION语句：
   对每个SELECT子句进行与DML语句相同的子查询嵌套检查。

3. 针对WITH语句（CTE）：
   检查CTE中是否存在嵌套子查询，并按照DML语句的标准进行检查。
   递归检查所有嵌套的CTE。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00108(input *rulepkg.RuleHandlerInput) error {
	// 获取子查询嵌套阈值参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param should be an integer, got: %v", param.Value)
	}

	// 定义一个递归函数，用于计算子查询的最大嵌套层数
	var getMaxSubqueryDepth func(node ast.Node) int
	getMaxSubqueryDepth = func(node ast.Node) int {
		maxDepth := 0
		// 获取当前节点中的所有子查询
		subqueries := util.GetSubquery(node)
		for _, subquery := range subqueries {
			// 对每个子查询，递归计算其嵌套层数
			depth := 1 + getMaxSubqueryDepth(subquery.Query)
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return maxDepth
	}

	// 定义一个函数，用于检查给定节点的子查询嵌套层数是否超过阈值
	checkAndReport := func(node ast.Node) {
		depth := getMaxSubqueryDepth(node)
		if depth > threshold {
			// 记录规则违规
			rulepkg.AddResult(input.Res, input.Rule, SQLE00108)
		}
	}

	// 根据输入节点的类型，执行相应的检查
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt:
		// 对于SELECT、INSERT、UPDATE、DELETE语句，检查其子查询嵌套层数
		checkAndReport(stmt)

	case *ast.UnionStmt:
		// 对于UNION语句，分别检查每个SELECT子句的子查询嵌套层数
		for _, selectStmt := range stmt.SelectList.Selects {
			checkAndReport(selectStmt)
		}

	// Removed the case for *ast.WithStmt to fix the compilation error
	// TODO 针对WITH语句（CTE），解析器暂时不支持
	default:
		// 其他类型的语句不处理
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
