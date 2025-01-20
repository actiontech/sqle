package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00094 = "SQLE00094"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00094,
			Desc:       plocale.Rule00094Desc,
			Annotation: plocale.Rule00094Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			Level:      driverV2.RuleLevelNotice,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET",
				Desc:  plocale.Rule00094Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00094Message,
		Func:    RuleSQLE00094,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00094): "在 MySQL 中，避免使用不必要的内置函数.默认参数描述: 函数名, 默认参数值: JSON_ARRAY,GROUP_CONCAT,CONCAT_WS,FIND_IN_SET"
您应遵循以下逻辑：
1. 对于所有DML、DQL语句，
  1. 获取该规则指定的函数名列表，通过英文逗号拆分成一个函数集合
  2. 如果SQL语句的语法节点中出现了函数集合中的任意一个函数，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00094(input *rulepkg.RuleHandlerInput) error {

	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	// 获取规则指定的函数名列表，通过英文逗号拆分成一个函数集合
	violationsFuncs := strings.Split(param.String(), ",")
	if len(violationsFuncs) < 1 {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	// 判断是否有违规函数的方法
	checkViolationFunc := func(checkExpr ast.ExprNode, funcs []string) bool {
		isExist := false
		util.ScanWhereStmt(func(expr ast.ExprNode) bool {
			switch pattern := expr.(type) {
			case *ast.FuncCallExpr:
				for _, name := range funcs {
					if strings.EqualFold(pattern.FnName.O, name) {
						isExist = true
					}
				}
			case *ast.AggregateFuncExpr:
				for _, name := range funcs {
					if strings.EqualFold(pattern.F, name) {
						isExist = true
					}
				}
			}
			return false
		}, checkExpr)

		return isExist
	}

	// select中： 查询列、where、having、group by、order by中可能涉及到的违规函数
	// select xx(col1) from table where xx(col1)=? having xx(col1)=? order by xx(col1)
	checkViolationFuncBySelect := func(selectNode *ast.SelectStmt, funcs []string) bool {
		// select col1、col2...
		for _, field := range selectNode.Fields.Fields {
			if checkViolationFunc(field.Expr, funcs) {
				return true
			}
		}
		// where
		if selectNode.Where != nil {
			if checkViolationFunc(selectNode.Where, funcs) {
				return true
			}
		}

		// group by
		if selectNode.GroupBy != nil {
			for _, groupby := range selectNode.GroupBy.Items {
				if checkViolationFunc(groupby.Expr, funcs) {
					return true
				}
			}
		}
		// having
		if selectNode.Having != nil {
			if checkViolationFunc(selectNode.Having.Expr, funcs) {
				return true
			}
		}

		// order by
		if selectNode.OrderBy != nil {
			for _, orderby := range selectNode.OrderBy.Items {
				if checkViolationFunc(orderby.Expr, funcs) {
					return true
				}
			}
		}
		return false
	}

	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}

	selectList := util.GetSelectStmt(input.Node)

	// 提取dml中所有的select语句（包括子查询
	for _, selectNode := range selectList {
		if checkViolationFuncBySelect(selectNode, violationsFuncs) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00094, param)
			return nil
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		// 上面已处理
	case *ast.DeleteStmt:
		if whereList := util.GetWhereExprFromDMLStmt(stmt); whereList != nil {
			for _, where := range whereList {
				if checkViolationFunc(where, violationsFuncs) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00094, param)
					return nil
				}
			}
		}
	case *ast.UpdateStmt:
		// set ...
		for _, setItem := range stmt.List {
			if checkViolationFunc(setItem.Expr, violationsFuncs) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00094, param)
				return nil
			}
		}
		// where ...
		if whereList := util.GetWhereExprFromDMLStmt(stmt); whereList != nil {
			for _, where := range whereList {
				if checkViolationFunc(where, violationsFuncs) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00094, param)
					return nil
				}
			}
		}
	case *ast.InsertStmt:
		if stmt.Lists != nil {
			for i := range stmt.Lists {
				for j := range stmt.Lists[i] {
					item := stmt.Lists[i][j]
					if checkViolationFunc(item, violationsFuncs) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00094, param)
						return nil
					}
				}
			}
		}

	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
