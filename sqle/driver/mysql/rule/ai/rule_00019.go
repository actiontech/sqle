package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00019 = "SQLE00019"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00019,
			Desc:         plocale.Rule00019Desc,
			Annotation:   plocale.Rule00019Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00019Message,
		Func:    RuleSQLE00019,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00019): "在 MySQL 中，不建议使用复合类型（SET和ENUM类型）数据."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句：
   - 解析语法树以识别字段定义。
   - 检查每个字段的数据类型，使用辅助函数 IsColumnTypeEqual 检查字段类型是否为 ENUM 或 SET。
   - 如果发现字段的数据类型为 ENUM 或 SET，则记录该字段并报告违反规则。

2. 对于 "ALTER TABLE..." 语句：
   - 解析语法树以识别字段变更或新增定义。
   - 检查变更或新增字段的数据类型，使用辅助函数 IsColumnTypeEqual 检查字段类型是否为 ENUM 或 SET。
   - 如果发现字段的数据类型为 ENUM 或 SET，则记录该字段并报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00019(input *rulepkg.RuleHandlerInput) error {

	// 确保输入的节点不为空
	if input.Node == nil {
		return nil
	}

	// 遍历所有SQL语句
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 处理 CREATE TABLE 语句
		for _, column := range stmt.Cols {
			if util.IsColumnTypeEqual(column, mysql.TypeEnum, mysql.TypeSet) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00019)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		// 处理 ALTER TABLE 语句
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn:
				for _, column := range spec.NewColumns {
					if util.IsColumnTypeEqual(column, mysql.TypeEnum, mysql.TypeSet) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00019)
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
