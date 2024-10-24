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
	SQLE00068 = "SQLE00068"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00068,
			Desc:       "在 MySQL 中, 禁止使用TIMESTAMP字段",
			Annotation: "TIMESTAMP类型字段受制于2038年问题，其时间范围仅限于1970-01-01 00:00:01 UTC至2038-01-19 03:14:07 UTC。超过这个时间范围，TIMESTAMP将无法存储更晚的时间点，导致应用报错。此外，TIMESTAMP字段在存储时会根据数据库服务器的时区进行转换，这可能导致跨时区应用中的时间不一致问题。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止使用TIMESTAMP字段",
		AllowOffline: true,
		Func:         RuleSQLE00068,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00068): "在 MySQL 中，禁止使用TIMESTAMP字段."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句，执行以下检查：
   1. 使用辅助函数 IsColumnTypeEqual 检查语法节点中是否包含 TIMESTAMP 字段定义。
   若包含，则报告违反规则。

2. 对于 "ALTER TABLE..." 语句，执行以下检查：
   1. 使用辅助函数 IsColumnTypeEqual 检查语法节点中是否包含 TIMESTAMP 字段定义。
   若包含，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00068(input *rulepkg.RuleHandlerInput) error {
	// 存储所有违反规则的列
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	// 检查 "CREATE TABLE" 语句
	case *ast.CreateTableStmt:
		// 遍历所有列定义
		for _, col := range stmt.Cols {
			// 使用辅助函数检查列类型是否为 TIMESTAMP
			if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
				violateColumns = append(violateColumns, col)
			}
		}

	// 检查 "ALTER TABLE" 语句
	case *ast.AlterTableStmt:
		// 获取所有 ADD COLUMN、CHANGE COLUMN 和 MODIFY COLUMN 操作
		alterTypes := []ast.AlterTableType{
			ast.AlterTableAddColumns,
			ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn,
		}
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, alterTypes...) {
			// 遍历每个新的列定义
			for _, col := range spec.NewColumns {
				// 使用辅助函数检查列类型是否为 TIMESTAMP
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
					violateColumns = append(violateColumns, col)
				}
			}
		}

	default:
		// 非 "CREATE TABLE" 或 "ALTER TABLE" 语句，不处理
		return nil
	}

	// 如果存在任何违反规则的列，则报告违规
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00068, util.JoinColumnNames(violateColumns))
	}

	return nil
}

// ==== Rule code end ====
