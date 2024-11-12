package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00022 = "SQLE00022"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00022,
			Desc:       "表的列数不建议超过阈值",
			Annotation: "避免在OLTP系统上做宽表设计，后期对性能影响很大；具体规则阈值可根据业务需求调整，默认值：40",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "40",
					Desc:  "最大列数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "表的列数不建议超过阈值. 阈值: %v",
		AllowOffline: false,
		Func:    RuleSQLE00022,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00022): "In table definition, the count of table columns should be within threshold", the threshold should be a parameter whose default value is 40.
You should follow the following logic:
1. For "create table ..." statement, check column count should be within threshold, otherwise, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00022(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxColumnCount, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be a number", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."
		if len(stmt.Cols) > maxColumnCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00022, maxColumnCount)
		}
	default:
		return nil
	}
	return nil
}

// ==== Rule code end ====
