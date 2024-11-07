package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00021 = "SQLE00021"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00021,
			Desc:       "在 MySQL 中, 禁止表字段缺少NOT NULL约束",
			Annotation: "若数据库表字段缺少NOT NULL约束，则字段存储值可能是NULL，后期判断时，需要加上IS NULL判断，增加SQL编写的复杂度。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止表字段缺少NOT NULL约束",
		AllowOffline: true,
		Func:         RuleSQLE00021,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00021): "在 MySQL 中，禁止表字段缺少NOT NULL约束."
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句：
   - 检查表定义中的每个字段（如 INT、VARCHAR、DECIMAL 等）是否包含 NOT NULL 约束，使用辅助函数 IsColumnHasOption 检查字段是否具有 NOT NULL 约束。
   - 如果发现任何字段未指定 NOT NULL 约束，则记录为违反规则。

2. 针对 "ALTER TABLE..." 语句：
   1. 当添加新列时：
      - 检查新列定义是否包含 NOT NULL 约束，使用辅助函数 IsColumnHasOption 检查字段是否具有 NOT NULL 约束。
      - 如果未包含 NOT NULL 约束，则记录为违反规则。
   2. 当修改现有列时：
      - 检查修改后的列定义是否移除了原有的 NOT NULL 约束，使用辅助函数 IsColumnHasOption 检查字段是否具有 NOT NULL 约束。
      - 如果移除了 NOT NULL 约束或未添加 NOT NULL 约束，则记录为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00021(input *rulepkg.RuleHandlerInput) error {

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if !util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00021)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				if !util.IsColumnHasOption(col, ast.ColumnOptionNotNull) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00021)
					return nil
				}
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
