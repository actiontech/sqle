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
	SQLE00016 = "SQLE00016"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00016,
			Desc:       "在 MySQL 中, 存储大数据类型（如长文本、图片等）的字段只能设置为NULL",
			Annotation: "在MySQL中，存储大数据类型的内容常用BLOB和TEXT、GEOMETRY以及JSON类型，但它们无法指定默认值；写入数据时，如未对该字段指定值会导致写入失败。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 存储大数据类型（如长文本、图片等）的字段只能设置为NULL",
		AllowOffline: true,
		Func:         RuleSQLE00016,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00016): "在 MySQL 中，存储大数据类型（如长文本、图片等）的字段只能设置为NULL."
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句，执行以下检查：
   1. 确认表中是否存在 BLOB、TEXT、GEOMETRY 或 JSON 类型的字段，使用辅助函数IsColumnTypeEqual进行检查。
   2. 检查这些字段是否被设置了 NOT NULL 约束，使用辅助函数IsColumnHasOption进行检查。
   如果上述两个条件同时满足，则标记为违反规则。

2. 针对 "ALTER TABLE..." 语句，执行以下检查：
   1. 确认修改中是否涉及 BLOB、TEXT、GEOMETRY 或 JSON 类型的字段，使用辅助函数IsColumnTypeEqual进行检查。
   2. 检查这些字段是否被设置了 NOT NULL 约束，使用辅助函数IsColumnHasOption进行检查。
   如果上述两个条件同时满足，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00016(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 遍历 CREATE TABLE 语句中的所有列
		for _, col := range stmt.Cols {
			// 检查列类型是否为 BLOB、TEXT、GEOMETRY 或 JSON
			if util.IsColumnTypeEqual(col, append(util.GetBlobDbTypes(), mysql.TypeJSON, mysql.TypeGeometry)...) {
				// 检查列是否设置了 NOT NULL 约束
				if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
					violateColumns = append(violateColumns, col)
				}
			}
		}

	case *ast.AlterTableStmt:
		// 获取 ALTER TABLE 语句中涉及的 MODIFY 和 CHANGE 操作
		alterSpecs := util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn)
		for _, spec := range alterSpecs {
			for _, col := range spec.NewColumns {
				// 检查列类型是否为 BLOB、TEXT、GEOMETRY 或 JSON
				if util.IsColumnTypeEqual(col, append(util.GetBlobDbTypes(), mysql.TypeJSON, mysql.TypeGeometry)...) {
					// 检查列是否设置了 NOT NULL 约束
					if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
						violateColumns = append(violateColumns, col)
					}
				}
			}
		}

	default:
		// 其他类型的 SQL 语句不处理
		return nil
	}

	// 如果存在违反规则的列，记录违规结果
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00016)
	}
	return nil
}

// ==== Rule code end ====
