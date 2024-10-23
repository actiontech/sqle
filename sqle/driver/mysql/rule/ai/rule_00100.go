package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00100 = "SQLE00100"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00100,
			Desc:       "在 MySQL 中, 避免SELECT语句一次性返回的结果过多",
			Annotation: "如果查询的扫描行数很大，会导致IO、网络资源消耗过大，并且可能会导致优化器选择错误的执行计划而不走索引。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "1000",
					Desc:  "结果集返回行数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 避免SELECT语句一次性返回的结果过多",
		AllowOffline: false,
		Func:         RuleSQLE00100,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00100): "在 MySQL 中，避免SELECT语句一次性返回的结果过多.默认参数描述: 结果集返回行数, 默认参数值: 1000"
您应遵循以下逻辑：
1. 对于所有DML语句中的“SELECT ...”子句：
   1. 登录数据库。
   2. 使用辅助函数GetExecutionPlan获取SELECT子句的执行计划。
   3. 检查执行计划中最底层的行数节点。
   4. 如果行数节点的值大于设定的阈值，则报告违反规则。

2. 对于所有DML语句中的“SELECT ...”子句：
   1. 检查是否存在LIMIT语法节点。
   2. 如果存在，验证LIMIT节点后的行数。
   3. 如果行数大于设定的阈值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00100(input *rulepkg.RuleHandlerInput) error {
	// 获取数值类型的规则参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	threshold := param.Int()
	if threshold <= 0 {
		return fmt.Errorf("param value should be greater than 0")
	}

	checkLimit := func(limit *ast.Limit) (bool, error) {
		if limit != nil {
			if xx, ok := limit.Count.(ast.ValueExpr); ok {
				count, err := strconv.Atoi(fmt.Sprintf("%v", xx.GetValue()))
				if err != nil {
					return false, err
				}
				if count > threshold {
					return true, nil
				}
			}
		}
		return false, nil
	}

	checkExplain := func(node ast.Node) (bool, error) {
		switch stmt := node.(type) {
		case *ast.SelectStmt, *ast.UnionStmt:
			executionPlan, err := util.GetExecutionPlan(input.Ctx, stmt.Text()) // TODO [解析器缺陷] 当sql是insert ... select语句中的SelectStmt/UnionStmt的Text() 为''，因此无法对应 test code中Explain
			if err != nil {
				return false, fmt.Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)

			}
			for _, record := range executionPlan.Plan {
				if record.Rows > int64(threshold) {
					return true, nil
				}
			}
		}
		return false, nil
	}

	var processDML func(node ast.Node) (bool, error)
	processDML = func(node ast.Node) (bool, error) {
		switch stmt := node.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			subs := util.GetSubquery(stmt)
			for _, sub := range subs {
				return processDML(sub.Query)
			}
		case *ast.InsertStmt:
			if stmt.Select != nil {
				return processDML(stmt.Select)
			}
		case *ast.SelectStmt:
			isViolate, err := checkLimit(stmt.Limit)
			if err != nil {
				return false, err
			}
			if !isViolate {
				// check explain
				isViolate2, err := checkExplain(stmt)
				if err != nil {
					return false, err
				}
				return isViolate2, nil
			}
			return isViolate, nil
		case *ast.UnionStmt:
			isViolate, err := checkLimit(stmt.Limit)
			if err != nil {
				return false, err
			}
			if !isViolate {
				// check explain
				isViolate2, err := checkExplain(stmt)
				if err != nil {
					return false, err
				}
				return isViolate2, nil
			}
			return isViolate, nil
		}
		return false, nil
	}

	isViolate, err := processDML(input.Node)
	if err != nil {
		log.NewEntry().Errorf("%s", err)
	}
	if isViolate {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00100)
	}
	return nil
}

// 规则函数实现结束

// ==== Rule code end ====
