package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00018 = "SQLE00018"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00018,
			Desc:       "在 MySQL 中, CHAR长度大于20时，建议使用VARCHAR类型",
			Annotation: "VARCHAR是变长字段，存储空间小，可节省存储空间，同时相对较小的字段检索效率显然也要高些",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "20",
					Desc:  "CHAR最大长度",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, CHAR长度大于20时，建议使用VARCHAR类型",
		AllowOffline: true,
		Func:         RuleSQLE00018,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00018): "在 MySQL 中，CHAR长度大于20时，建议使用VARCHAR类型.默认参数描述: CHAR最大长度, 默认参数值: 20"
您应遵循以下逻辑：
1. 对于 CREATE TABLE 语句，执行以下检查：
   1. 检查列定义中的每个字段节点，确认是否为 CHAR 类型，使用辅助函数 IsColumnTypeEqual。
   2. 对于每个 CHAR 类型字段节点，确认其长度是否超过 20，使用辅助函数 GetColumnWidth。
   3. 如果存在长度超过 20 的 CHAR 类型字段，则报告违反规则。

2. 对于 ALTER TABLE 语句，执行以下检查：
   1. 检查新增或修改的字段节点，确认是否为 CHAR 类型，使用辅助函数 IsColumnTypeEqual。
   2. 对于每个新增或修改的 CHAR 类型字段节点，确认其长度是否超过 20，使用辅助函数 GetColumnWidth。
   3. 如果存在长度超过 20 的 CHAR 类型字段，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00018(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	threshold := param.Int()
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeString) && util.GetColumnWidth(col) > threshold {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeString) && util.GetColumnWidth(col) > threshold {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00018, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====
