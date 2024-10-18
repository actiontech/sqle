package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00176 = "SQLE00176"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00176,
			Desc:       "在 MySQL 中, 不建议SQL中包含hint指令",
			Annotation: "使用hint可能会导致数据库走错误的执行计划，从而影响执行效率，消耗系统资源。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 不建议SQL中包含hint指令",
		AllowOffline: true,
		Func:         RuleSQLE00176,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00176): "在 MySQL 中，不建议SQL中包含hint指令."
您应遵循以下逻辑：
1. 对所有 DML（INSERT、UPDATE、DELETE、REPLACE）和 DQL（SELECT）语句进行检查，如果以下任意一项为真，则报告违反规则：
  1. 存在 hit注释 形式的注释块。
  2. 表名后存在 FORCE INDEX 或 FORCE KEY 语法节点。
  3. 表名后存在 USE INDEX 或 USE KEY 语法节点。
  4. 表名后存在 IGNORE INDEX 或 IGNORE KEY 语法节点。
  5. 语句中存在 STRAIGHT_JOIN 语法节点。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00176(input *rulepkg.RuleHandlerInput) error {
	hasIndexHint := func(node ast.Node) bool {
		tableNames := util.GetTableNames(node)
		for _, tableName := range tableNames {
			if len(tableName.IndexHints) > 0 {
				return true
			}
		}
		return false
	}

	hasStraighJoin := func(node ast.Node) bool {
		joinNode := util.GetFirstJoinNodeFromStmt(input.Node)
		if joinNode == nil {
			return false
		}
		var checkJoinNode func(joinNode *ast.Join) bool
		checkJoinNode = func(joinNode *ast.Join) bool {
			if joinNode.StraightJoin {
				return true
			} else {
				if l, ok := joinNode.Left.(*ast.Join); ok {
					if checkJoinNode(l) {
						return true
					}
				}
			}
			return false
		}
		return checkJoinNode(joinNode)
	}
	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}
	for _, selectStmt := range util.GetSelectStmt(input.Node) {
		if selectStmt.SelectStmtOpts.StraightJoin || len(selectStmt.TableHints) > 0 || hasIndexHint(selectStmt) || hasStraighJoin(selectStmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00176)
			return nil
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		if len(stmt.TableHints) > 0 || hasIndexHint(stmt) || hasStraighJoin(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00176)
		}
	case *ast.InsertStmt:
		if len(stmt.TableHints) > 0 || hasIndexHint(stmt) || hasStraighJoin(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00176)
		}
	case *ast.UpdateStmt:
		if len(stmt.TableHints) > 0 || hasIndexHint(stmt) || hasStraighJoin(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00176)
		}
		// TODO 无法解析 delete/update/insert /*+ index(xx,xx) */...
		// TODO 无法解析 delete/update/insert /*+ set_var(xx) */...
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
