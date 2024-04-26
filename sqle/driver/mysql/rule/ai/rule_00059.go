package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00059 = "SQLE00059"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00059,
			Desc:       "对于MySQL的DDL, 禁止修改大表字段类型",
			Annotation: "对于大型数据表，修改字段类型的DDL操作将导致显著的性能下降和可用性影响。此类操作通常需要复制整个表来更改数据类型，期间表将无法进行写操作，并且可能导致长时间的锁等待，对线上业务造成过长时间的影响。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "1000000",
					Desc:  "目标表的最大行数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "对于MySQL的DDL，禁止修改大表字段类型，目标表的最大行数: %v",
		Func:    RuleSQLE00059,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00059): "For MySQL DDL, modifying large table field types is prohibited.".
You should follow the following logic:
1. For "alter table... modify..." Statement, check whether the number of rows in the alter target table is greater than the rule threshold variable min_rows, if it is greater than the specified threshold, If it does, report a violation. the number of rows in the alter target table should be the information obtained online.
2. For "alter table... change..." Statement to perform a check similar to the one described above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00059(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	min_rows, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s must be an integer", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table"
		// get the number of rows in the table
		rows, err := util.GetTableRowCount(input.Ctx, stmt.Table)
		if err != nil {
			return err
		}
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			// "alter table... change/alter table... modify"
			if rows > min_rows || len(spec.NewColumns) > min_rows {
				// the "alter table" is modifying large table field type
				rulepkg.AddResult(input.Res, input.Rule, SQLE00059, min_rows)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
