package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00001 = "SQLE00001"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00001,
			Desc:       "在 MySQL 中, 禁止SQL语句不带WHERE条件或者WHERE条件为永真",
			Annotation: "使用有效的WHERE条件能够避免全表扫描，提高SQL执行效率；而恒为TRUE的WHERE条件，如where 1=1、where true=true等，在执行时会进行全表扫描产生额外开销。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止SQL语句不带WHERE条件或者WHERE条件为永真",
		AllowOffline: false,
		Func:         RuleSQLE00001,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00001): "在 MySQL 中，禁止SQL语句不带WHERE条件或者WHERE条件为永真."
您应遵循以下逻辑：
1. 对于DML语句的SELECT子句（包括SELECT、INSERT、UPDATE、DELETE、UNION语句中的子SELECT语句），如果满足以下任一条件，则报告违反规则：
   1. SQL语句未包含WHERE条件。
   2. 使用辅助函数IsExprConstTrue检查WHERE条件是否为永真表达式。
   3. WHERE条件的最外层使用OR，并且OR条件中包含恒真表达式。

2. 对于"WITH.."语句，执行与上述相同的检查。

3. 对于DML语句的SELECT子句（包括SELECT、INSERT、UPDATE、DELETE、UNION语句中的子SELECT语句），如果满足以下任一条件，则报告违反规则：
   1. WHERE条件为column IS NOT NULL，且column是不可为空的字段。
   2. WHERE条件的最外层使用OR，并且OR条件中包含上述条件。

4. 对于"WITH.."语句，执行与上述相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====

// 规则函数实现开始
func RuleSQLE00001(input *rulepkg.RuleHandlerInput) error {

	// 检索出所有的select语句中的 where
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		selectList := util.GetSelectStmt(stmt)
		for _, sel := range selectList {
			if sel.From != nil {
				aliasInfo := util.GetTableAliasInfoFromJoin(sel.From.TableRefs)
				if sel.Where != nil {
					isConst := util.IsExprConstTrue(input.Ctx, sel.Where, aliasInfo)
					if isConst {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
						return nil
					}
				} else {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
					return nil
				}
			}
		}
	}

	// 单独处理 delete、update的 where
	switch stmt2 := input.Node.(type) {
	case *ast.DeleteStmt:
		if stmt2.Where == nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
			return nil
		} else {
			aliasInfos := util.GetTableAliasInfoFromJoin(stmt2.TableRefs.TableRefs)
			isConst := util.IsExprConstTrue(input.Ctx, stmt2.Where, aliasInfos)
			if isConst {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
				return nil
			}
		}
	case *ast.UpdateStmt:
		if stmt2.Where == nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
			return nil
		} else {
			aliasInfos := util.GetTableAliasInfoFromJoin(stmt2.TableRefs.TableRefs)
			isConst := util.IsExprConstTrue(input.Ctx, stmt2.Where, aliasInfos)
			if isConst {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
				return nil
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
