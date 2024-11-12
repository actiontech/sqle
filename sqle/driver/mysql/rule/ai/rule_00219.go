package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00219 = "SQLE00219"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00219,
			Desc:       "建表DDL必须包括创建时间字段，并应确保该字段能记录表记录的创建时间。",
			Annotation: "使用创建时间字段，有利于问题查找跟踪和检索数据，同时避免后期对数据生命周期管理不便 ，可保证时间的准确性",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "create_time",
					Desc:  "创建时间字段名",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "建表DDL必须包括创建时间字段，并应确保该字段能记录表记录的创建时间。",
		AllowOffline: true,
		Func:         RuleSQLE00219,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00219): "在 MySQL 中，建表DDL必须包括创建时间字段，并应确保该字段能记录表记录的创建时间。.默认参数描述: 创建时间字段名, 默认参数值: create_time"
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句，逐条验证以下条件，若任一条件不满足，则报告违反规则：
   1. 表中必须包含一个名为规则变量值（如 create_time）的列，且数据类型为 timestamp。
   2. 使用辅助函数GetColumnOption获取该列的默认值选项，并使用辅助函数IsOptionFuncCall检查该默认值是否设置为 `CURRENT_TIMESTAMP`。

2. 针对 "ALTER TABLE..." 语句，若新增或修改的列名为规则变量值（如 create_time），则逐条验证以下条件，若任一条件不满足，则报告违反规则：
   1. 该列的数据类型必须为 timestamp。
   2. 使用辅助函数GetColumnOption获取该列的默认值选项，并使用辅助函数IsOptionFuncCall检查该默认值是否设置为 `CURRENT_TIMESTAMP`。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00219(input *rulepkg.RuleHandlerInput) error {
	// 获取规则变量值，例如 "create_time"
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	targetColumnName := param.String()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 标记是否找到目标列
		found := false
		// 标记是否违反规则
		violated := false

		// 遍历所有列定义
		for _, col := range stmt.Cols {
			if strings.EqualFold(util.GetColumnName(col), targetColumnName) {
				found = true

				// 检查数据类型是否为 TIMESTAMP
				if !util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
					violated = true
				} else {
					// 获取并检查默认值是否为 CURRENT_TIMESTAMP
					defaultOption := util.GetColumnOption(col, ast.ColumnOptionDefaultValue)
					if defaultOption == nil || !util.IsOptionFuncCall(defaultOption, "current_timestamp") {
						violated = true
					}
				}

				// 找到目标列后无需继续遍历
				break
			}
		}

		// 如果未找到目标列或违反任一条件，则报告规则违规
		if !found || violated {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00219)
		}

	case *ast.AlterTableStmt:
		// 存储所有违反规则的列
		var violatedColumns []*ast.ColumnDef

		// 获取所有 ADD COLUMN 和 MODIFY COLUMN 操作
		addOrModifySpecs := util.GetAlterTableCommandsByTypes(
			stmt,
			ast.AlterTableAddColumns,
			ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn,
		)

		// 遍历每个操作
		for _, spec := range addOrModifySpecs {
			for _, col := range spec.NewColumns {
				if strings.EqualFold(util.GetColumnName(col), targetColumnName) {
					// 检查数据类型是否为 TIMESTAMP
					if !util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
						violatedColumns = append(violatedColumns, col)
					} else {
						// 获取并检查默认值是否为 CURRENT_TIMESTAMP
						defaultOption := util.GetColumnOption(col, ast.ColumnOptionDefaultValue)
						if defaultOption == nil || !util.IsOptionFuncCall(defaultOption, "current_timestamp") {
							violatedColumns = append(violatedColumns, col)
						}
					}

				}
			}
		}

		// 如果存在任何违反规则的列，则报告
		if len(violatedColumns) > 0 {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00219)
		}

	default:
		// 非 CREATE TABLE 或 ALTER TABLE 语句，不处理
		return nil
	}

	return nil
}

// ==== Rule code end ====
