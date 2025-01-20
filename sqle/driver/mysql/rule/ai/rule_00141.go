package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00141 = "SQLE00141"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00141,
			Desc:       plocale.Rule00141Desc,
			Annotation: plocale.Rule00141Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			Level:      driverV2.RuleLevelNotice,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "3",
				Desc:  plocale.Rule00141Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00141Message,
		Func:    RuleSQLE00141,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00141): "在 MySQL 中，表关联嵌套循环的层次过多.默认参数描述: 表关联嵌套循环层数, 默认参数值: 3"
您应遵循以下逻辑：
1、检查句子中是否存在 FROM 子句，如果存在，则进一步检查。
2、使用辅助函数GetTableNames获取 FROM 子句中参与表连接的个数（如表 st1、st_ps、st_addr、st_dp），如果表连接个数超过阈值，报告违反规则。
3、检查当前句子的语法节点是否为 UPDATE 语句，如果是则进入下一步检查。
4、使用辅助函数GetTableNames获取 UPDATE 语句中 JOIN 操作的表连接个数（如表 st1、st_ps、st_addr、st_dp），如果表连接个数超过阈值，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00141(input *rulepkg.RuleHandlerInput) error {
	// 获取表连接的阈值参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxTableCount, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be an integer", param.Value)
	}

	// 检查 SQL 语句中是否存在 FROM 子句，并统计表的连接个数
	if stmt, ok := input.Node.(ast.DMLNode); ok {
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// 如果存在 FROM 子句
			if selectStmt.From != nil {
				// 获取 FROM 子句中所有表名
				tableNames := util.GetTableNames(selectStmt.From)
				tableCount := map[string]int{}
				for _, tableName := range tableNames {
					// 构建完整的表名：schema.table
					key := fmt.Sprintf("%s.%s", util.GetSchemaName(input.Ctx, tableName.Schema.L), tableName.Name.L)
					tableCount[key]++
				}
				// 检查表连接个数是否超过阈值
				if len(tableCount) > maxTableCount {
					// 报告规则违规
					rulepkg.AddResult(input.Res, input.Rule, SQLE00141)
				}
			}
		}
	}

	// 如果当前 SQL 语句是 UPDATE 语句，则进一步检查 JOIN 操作中的表连接个数
	if stmtUpdate, ok := input.Node.(*ast.UpdateStmt); ok {
		if stmtUpdate.TableRefs != nil {
			// 获取 UPDATE 语句中所有参与 JOIN 的表名
			joinTableNames := util.GetTableNames(stmtUpdate.TableRefs.TableRefs)
			joinTableCount := map[string]int{}
			for _, tableName := range joinTableNames {
				// 构建完整的表名：schema.table
				key := fmt.Sprintf("%s.%s", util.GetSchemaName(input.Ctx, tableName.Schema.L), tableName.Name.L)
				joinTableCount[key]++
			}
			// 检查 JOIN 表连接个数是否超过阈值
			if len(joinTableCount) > maxTableCount {
				// 报告规则违规
				rulepkg.AddResult(input.Res, input.Rule, SQLE00141)
			}
		}
	}

	// 检查 DELETE 语句中的 JOIN 操作的表连接个数
	if stmtDelete, ok := input.Node.(*ast.DeleteStmt); ok {
		if stmtDelete.TableRefs != nil {
			deleteJoinTableNames := util.GetTableNames(stmtDelete.TableRefs.TableRefs)
			deleteJoinTableCount := map[string]int{}
			for _, tableName := range deleteJoinTableNames {
				key := fmt.Sprintf("%s.%s", util.GetSchemaName(input.Ctx, tableName.Schema.L), tableName.Name.L)
				deleteJoinTableCount[key]++
			}
			if len(deleteJoinTableCount) > maxTableCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00141)
			}
		}
	}

	return nil
}

// ==== Rule code end ====
