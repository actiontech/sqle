package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00096 = "SQLE00096"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00096,
			Desc:       "不建议参与连接操作的表数量过多",
			Annotation: "表关联越多，意味着各种驱动关系组合就越多，比较各种结果集的执行成本的代价也就越高，进而SQL查询性能会大幅度下降。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "3",
					Desc:  "参与表连接的表个数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "不建议参与连接操作的表数量过多",
		AllowOffline: true,
		Func:         RuleSQLE00096,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00096): "在 MySQL 中，不建议参与连接操作的表数量过多.默认参数描述: 参与表连接的表个数, 默认参数值: 3"
您应遵循以下逻辑：
1. 对于 DML 语句（包括 SELECT、UPDATE、DELETE、INSERT ... SELECT、UNION），执行以下步骤：
   a. 使用新创建的辅助函数GetAllJoinsFromNode获取涉及的表名节点。
   b. 去除重复的表名。
   c. 如果表的总数超过预设的阈值，则报告违反规则。

2. 对于 WITH 语句，执行以下步骤：
   a. 使用新创建的辅助函数GetAllJoinsFromNode获取涉及的表名节点。
   b. 去除重复的表名。
   c. 如果表的总数超过预设的阈值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00096(input *rulepkg.RuleHandlerInput) error {
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	threshold := param.Int()
	if threshold == 0 {
		return fmt.Errorf("param value should be greater than 0")
	}

	var tableNames []string

	// 内部匿名的辅助函数
	getTableNamesFromSQL := func(ctx *session.Context, node ast.Node) []string {
		var tableNames []string

		getTableNameStr := func(t *ast.TableSource) string {
			if table, ok := t.Source.(*ast.TableName); ok && table != nil {
				schemaName := util.GetSchemaName(ctx, table.Schema.L)
				return fmt.Sprintf("%s.%s", schemaName, table.Name.L)
			}
			return ""
		}

		// 获取所有的 join节点
		joins := util.GetAllJoinsFromNode(node)
		for _, join := range joins {
			if t, ok := join.Left.(*ast.TableSource); ok && t != nil {
				tName := getTableNameStr(t)
				if tName != "" {
					tableNames = append(tableNames, tName)
				}
			}
			if t, ok := join.Right.(*ast.TableSource); ok && t != nil {
				tName := getTableNameStr(t)
				if tName != "" {
					tableNames = append(tableNames, tName)
				}
			}
		}
		return tableNames
	}
	// 是否有表连接个数超过3个的
	checkViolate := func(tables []string) bool {
		tableMap := make(map[string]struct{})
		for _, name := range tables {
			tableMap[name] = struct{}{}
		}
		if len(tableMap) > threshold {
			return true
		}
		return false
	}

	// 先看子查询中
	switch node := input.Node.(type) {
	case ast.DMLNode:
		subs := util.GetSubquery(node)
		for _, sub := range subs {
			tableNames = getTableNamesFromSQL(input.Ctx, sub)
			if checkViolate(tableNames) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00096)
				return nil
			}
		}
	}

	// 再看sql本身语句中
	switch node := input.Node.(type) {
	case *ast.InsertStmt:
		tableNames = getTableNamesFromSQL(input.Ctx, node)
		tableNames = tableNames[:len(tableNames)-1] // 移除最左边的table，因为最左边的table 是insert 而不是 join
	case ast.DMLNode:
		tableNames = getTableNamesFromSQL(input.Ctx, node)
	// 注释掉未定义的 ast.WithStmt case
	// case *ast.WithStmt:
	// 	tableNames = getTableNamesFromSQL(input.Ctx, node)
	default:
		return nil
	}

	if checkViolate(tableNames) {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00096)
	}
	return nil
}

// ==== Rule code end ====
