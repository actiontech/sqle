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
	SQLE00012 = "SQLE00012"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00012,
			Desc:       "在 MySQL 中, 建议使用BIGINT类型表示小数",
			Annotation: "在MySQL中，对于金额等需要高精度计算的小数，建议使用BIGINT类型表示，以避免浮点数精度问题。例如，可以用分来表示金额，1元在数据库中用整型表示为100。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议使用BIGINT类型表示小数",
		AllowOffline: true,
		Func:         RuleSQLE00012,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00012): "在 MySQL 中，建议使用BIGINT类型表示小数."
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句，执行以下检查：
   1. 检查语法节点中是否定义了 DECIMAL 类型字段（例如用于表示价格、金额、数量等）。
   如果发现 DECIMAL 类型字段，则标记为违反规则。

2. 针对 "ALTER TABLE..." 语句，执行以下检查：
   1. 检查语法节点中是否添加或修改为 DECIMAL 类型字段（例如用于表示价格、金额、数量等）。
   如果发现 DECIMAL 类型字段，则标记为违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00012(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	default:
		return nil
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00012, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
