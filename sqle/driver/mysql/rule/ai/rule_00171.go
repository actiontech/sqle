package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00171 = "SQLE00171"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00171,
			Desc:       plocale.Rule00171Desc,
			Annotation: plocale.Rule00171Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			Level:      driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "CREATE_TIME",
				Desc:  plocale.Rule00171Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00171Message,
		Func:    RuleSQLE00171,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
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
