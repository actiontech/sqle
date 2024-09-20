package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00034 = "SQLE00034"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00034,
			Desc:       "在 MySQL 中, 字段约束为NOT NULL时必须带默认值",
			Annotation: "如存在NOT NULL且不带默认值的字段，对字段进行写入时不包含该字段，会导致插入报错",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params:     params.Params{},
		},
		Message: "在 MySQL 中, 字段约束为NOT NULL时必须带默认值",
		Func:    RuleSQLE00034,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00034): "在 MySQL 中，字段约束为NOT NULL时必须带默认值."
您应遵循以下逻辑：
1. 对于"CREATE TABLE..."语句，检查语法树中的列定义节点，如果某个列定义包含NOT NULL约束但没有DEFAULT子节点，报告违反规则。
2. 对于"ALTER TABLE..."语句，检查语法树中的列修改节点，如果某个列修改包含NOT NULL约束但没有DEFAULT子节点，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00034(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			// if the column has "NOT NULL" constraint but no "DEFAULT" constraint
			if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) && !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				fmt.Println("666666")
				// if the column has "NOT NULL" constraint but no "DEFAULT" constraint
				if util.IsColumnHasOption(col, ast.ColumnOptionNotNull) && !util.IsColumnHasOption(col, ast.ColumnOptionDefaultValue) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00034)
		return nil
	}

	return nil
}

// ==== Rule code end ====
