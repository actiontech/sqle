package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00171 = "SQLE00171"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00171,
			Desc:       "建表DDL必须包含创建时间字段且默认值为CURRENT_TIMESTAMP",
			Annotation: "使用CREATE_TIME字段，有利于问题查找跟踪和检索数据，同时避免后期对数据生命周期管理不便 ，默认值为CURRENT_TIMESTAMP可保证时间的准确性",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "CREATE_TIME",
					Desc:  "创建时间字段名",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message: "建表DDL必须包含创建时间字段且默认值为CURRENT_TIMESTAMP",
		Func:    RuleSQLE00171,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00171): "In DDL, when creating table, table should have a field about create-timestamp, whose DEFAULT value should be 'CURRENT_TIMESTAMP'", the create-timestamp column name is a parameter whose default value is 'CREATE_TIME'.
You should follow the following logic:
For "create table ..." statement, check the following conditions, report violation if any condition is violated:
1. The table should have a create-timestamp column whose type is datetime or timestamp, and column name is same as the parameter
2. The create-timestamp column's DEFAULT value should be configured as 'CURRENT_TIMESTAMP'
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00171(input *rulepkg.RuleHandlerInput) error {
	// get expected create_time field name in config
	createTimeFieldName := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String()
	found := false

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table"
		for _, col := range stmt.Cols {
			if strings.EqualFold(util.GetColumnName(col), createTimeFieldName) {
				// the column is create_time column
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) || util.IsColumnTypeEqual(col, mysql.TypeDatetime) {
					// the column is Timestamp-type or Datetime-type
					if c := util.GetColumnOption(col, ast.ColumnOptionDefaultValue); nil != c {
						// the column has "DEFAULT" option
						if util.IsOptionFuncCall(c, "current_timestamp") {
							// the "DEFAULT" value is current_timestamp
							found = true
						}
					}
				}
			}
		}
	default:
		return nil
	}

	if !found {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00024, createTimeFieldName)
	}

	return nil
}

// ==== Rule code end ====
