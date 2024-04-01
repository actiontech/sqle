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
	SQLE00033 = "SQLE00033"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00033,
			Desc:       "对于MySQL的DDL, 建表DDL必须包含更新时间字段, 默认值为CURRENT_TIMESTAMP, ON UPDATE值为CURRENT_TIMESTAMP",
			Annotation: "使用更新时间字段，有利于问题查找跟踪和检索数据，同时避免后期对数据生命周期管理不便 ，默认值为UPDATE_TIME可保证时间的准确性",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "UPDATE_TIME",
					Desc:  "更新时间字段名",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message: "对于MySQL的DDL, 建表DDL必须包含更新时间字段, 默认值为CURRENT_TIMESTAMP, ON UPDATE值为CURRENT_TIMESTAMP. 更新时间字段名: %v",
		Func:    RuleSQLE00033,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00033): "In DDL, when creating table, table should have a field about update-timestamp, whose DEFAULT value and ON-UPDATE value should be both 'CURRENT_TIMESTAMP'", the update-timestamp column name is a parameter whose default value is 'UPDATE_TIME'.
You should follow the following logic:
For "create table ..." statement, check the following conditions, report violation if any condition is violated:
1. The table should have a update-timestamp column whose type is datetime or timestamp, and column name is same as the parameter
2. The update-timestamp column's DEFAULT value should be configured as 'CURRENT_TIMESTAMP'
3. The update-timestamp column's ON-UPDATE value should be configured as 'CURRENT_TIMESTAMP'
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00033(input *rulepkg.RuleHandlerInput) error {
	// get expected update_time column name in config
	updateTimeColumnName := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String()
	found := false

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table"
		for _, col := range stmt.Cols {
			if strings.EqualFold(util.GetColumnName(col), updateTimeColumnName) {
				// the column is update_time column
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) || util.IsColumnTypeEqual(col, mysql.TypeDatetime) {
					// the column is Timestamp-type or DateTime-type
					if c := util.GetColumnOption(col, ast.ColumnOptionDefaultValue); nil != c && util.IsOptionFuncCall(c, "current_timestamp") {
						// the column has "DEFAULT" constraint, the "DEFAULT" value is current_timestamp
						if c := util.GetColumnOption(col, ast.ColumnOptionOnUpdate); nil != c && util.IsOptionFuncCall(c, "current_timestamp") {
							// the column has "ON UPDATE" constraint, the "DEFAULT" value is current_timestamp
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
		rulepkg.AddResult(input.Res, input.Rule, SQLE00033, updateTimeColumnName)
	}
	return nil
}

// ==== Rule code end ====
