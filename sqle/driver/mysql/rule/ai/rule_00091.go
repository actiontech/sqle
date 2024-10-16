package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
)

const (
	SQLE00091 = "SQLE00091"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00091,
			Desc:       "在 MySQL 中, 建议表连接时有连接条件",
			Annotation: "为了确保连接操作的正确性和可靠性，应该始终指定连接条件，定义正确的关联关系。缺少连接条件，可能导致连接操作失败，最终数据库会使用笛卡尔积的方式进行处理，产生不正确的连接结果，并导致性能问题，消耗大量的CPU和内存资源。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议表连接时有连接条件",
		AllowOffline: true,
		Func:         RuleSQLE00091,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00091): "在 MySQL 中，建议表连接时有连接条件."
您应遵循以下逻辑：
1. 对于所有DML语句中的SELECT子句，识别所有JOIN操作（包括隐式和显式）。
   1. 检查每个JOIN操作的ON子句或WHERE条件节点中是否存在连接条件（例如：table1.column1 =、>、<、!=、<> table2.column1）。
   2. 如果任何JOIN操作缺少连接条件，则报告违反规则。

2. 对于UNION语句，检查其中所有SELECT子句，遵循DML语句的相同规则。

3. 对于WITH语句，检查其中所有SELECT子句，遵循DML语句的相同规则。

4. 对于INSERT...SELECT语句，识别SELECT子句中的所有JOIN操作（包括隐式和显式）。
   1. 检查每个JOIN操作的ON子句或WHERE条件节点中是否存在连接条件。
   2. 如果任何JOIN操作缺少连接条件，则报告违反规则。

5. 对于UPDATE语句，识别所有JOIN操作（包括隐式和显式）。
   1. 检查每个JOIN操作的ON子句或WHERE条件节点中是否存在连接条件。
   2. 如果任何JOIN操作缺少连接条件，则报告违反规则。

6. 对于DELETE语句，执行与UPDATE语句相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00091(input *rulepkg.RuleHandlerInput) error {
	// 内部函数: 检查dml语句的JOIN
	checkJoin := func(stmt ast.Node) bool {
		getWhereExpr := func(node ast.Node) (where ast.ExprNode) {
			switch stmt := node.(type) {
			case *ast.SelectStmt:
				if stmt.From == nil { //If from is null skip check. EX: select 1;select version
					return nil
				}
				where = stmt.Where
			case *ast.UpdateStmt:
				where = stmt.Where
			case *ast.DeleteStmt:
				where = stmt.Where

			}
			return
		}

		doesNotJoinTables := func(tableRefs *ast.Join) bool {
			return tableRefs == nil || tableRefs.Left == nil || tableRefs.Right == nil
		}

		checkJoinConditionByOn := func(expr ast.ExprNode) (hasCondition bool) {
			util.ScanWhereStmt(func(node ast.ExprNode) bool {
				binExpr, ok := node.(*ast.BinaryOperationExpr)
				if !ok {
					return false
				}

				// 检查操作符是否为连接条件相关的运算符
				if binExpr.Op != opcode.EQ && binExpr.Op != opcode.GT &&
					binExpr.Op != opcode.LT && binExpr.Op != opcode.NE &&
					binExpr.Op != opcode.GE && binExpr.Op != opcode.LE {
					return false
				}

				// 获取左右操作数的列名
				leftCols := util.GetColumnNameInExpr(binExpr.L)
				rightCols := util.GetColumnNameInExpr(binExpr.R)

				if len(leftCols) == 0 || len(rightCols) == 0 {
					return false
				}

				leftTable := leftCols[0].Name.Table.L
				rightTable := rightCols[0].Name.Table.L

				// 确保比较的列来自不同的表
				if leftTable != "" && rightTable != "" && leftTable != rightTable {
					hasCondition = true
					return true // 找到一个连接条件，停止扫描
				}
				return false
			}, expr)
			return hasCondition
		}

		var checkJoinConditionInJoinNode func(ctx *session.Context, whereStmt ast.ExprNode, joinNode *ast.Join) (joinTables, hasCondition bool)
		checkJoinConditionInJoinNode = func(ctx *session.Context, whereStmt ast.ExprNode, joinNode *ast.Join) (joinTables, hasCondition bool) {
			if joinNode == nil {
				return false, false
			}
			if doesNotJoinTables(joinNode) {
				// 非JOIN两表的JOIN节点 一般是叶子节点 不检查
				return false, false
			}

			// 深度遍历左子树类型为ast.Join的节点 一旦有节点是JOIN两表的节点，并且没有连接条件，则返回
			if l, ok := joinNode.Left.(*ast.Join); ok {
				joinTables, hasCondition = checkJoinConditionInJoinNode(ctx, whereStmt, l) // 递归检测
				if joinTables && !hasCondition {
					return joinTables, hasCondition
				}
			}

			// 判断
			if joinNode.Using != nil {
				return true, true
			} else if joinNode.On == nil || !checkJoinConditionByOn(joinNode.On.Expr) {
				if (whereStmt == nil) || !checkJoinConditionByOn(whereStmt) {
					return true, false
				} else {
					return true, true
				}
			} else {
				return true, true
			}
		}

		joinNode := util.GetJoinNodeFromNode(stmt)
		whereStmt := getWhereExpr(stmt)
		if joinNode == nil {
			return true
		}
		joinTables, hasCondition := checkJoinConditionInJoinNode(input.Ctx, whereStmt, joinNode)
		if joinTables && !hasCondition {
			return true
		}
		return false
	}

	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}

	// dml中所有的select语句
	selectStmts := util.GetSelectStmt(input.Node)
	for _, selectStmt := range selectStmts {
		if checkJoin(selectStmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00091)
			return nil
		}
	}

	// 特殊处理：update join, delete join
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt, *ast.DeleteStmt:
		if checkJoin(stmt) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00091)
			return nil
		}

	// TODO 针对WITH语句（CTE），解析器暂时不支持

	default:
		return nil
	}

	return nil
}

// ==== Rule code end ====
