package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
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
			Desc:       "禁止修改大表字段类型",
			Annotation: "对于大型数据表，修改字段类型的DDL操作将导致显著的性能下降和可用性影响。此类操作通常需要复制整个表来更改数据类型，期间表将无法进行写操作，并且可能导致长时间的锁等待，对线上业务造成过长时间的影响。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "表大小(GB)",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "禁止修改大表字段类型，表大小阈值: %v GB",
		AllowOffline: false,
		Func:         RuleSQLE00059,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00059): "For MySQL DDL, modifying large table field types is prohibited.".
You should follow the following logic:
1. For "alter table... modify..." Statement, check whether the table size in the alter target table is greater than the rule threshold variable, if it is greater than the specified threshold, If it does, report a violation. The size of the table can be determined by 'select round((index_length+data_length)/1024/1024/1024) 'size_GB' from information_schema.tables where table_name=' table name ' 'to get, the table size needs to be obtained online,.
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
	maxSize, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s must be an integer", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table"

		if len(util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn)) > 0 {
			// "alter table... change/alter table... modify"

			// get the size of table
			size, err := util.GetTableSizeMB(input.Ctx, stmt.Table.Name.String())
			if err != nil {
				log.NewEntry().Errorf("get table size failed, sqle: %v, error: %v", input.Node.Text(), err)
				return nil
			}

			if size > int64(maxSize*1024) {
				// report rule violation
				rulepkg.AddResult(input.Res, input.Rule, SQLE00059, maxSize)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
