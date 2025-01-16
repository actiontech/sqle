package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00019 = "SQLE00019"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00019,
			Desc:       "不建议使用复合类型（SET和ENUM类型）数据",
			Annotation: "SET类型，ENUM类型不是SQL标准，移植性较差；后期如修改或增加枚举值需重建整张表，代价较大；且无法通过字面值进行排序；在插入数据时，必须带上引号，否则将写入枚举值的顺序值，造成不可预期的问题",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "不建议使用复合类型（SET和ENUM类型）数据",
		AllowOffline: true,
		Func:         RuleSQLE00019,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
