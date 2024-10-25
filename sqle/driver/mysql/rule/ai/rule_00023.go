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
	SQLE00023 = "SQLE00023"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00023,
			Desc:       "对于MySQL的DDL, 主键包含的列数不建议超过阈值",
			Annotation: "主建中的列过多，会导致二级索引占用更多的空间，同时增加索引维护的开销；具体规则阈值可根据业务需求调整，默认值：2",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "2",
					Desc:  "最大列数",
					Type:  params.ParamTypeInt,
				},
			},
		},

		Message: "对于MySQL的DDL, 主键包含的列数不建议超过阈值. 阈值: %v",
		AllowOffline: true,
		Func:    RuleSQLE00023,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00023): "In table definition, the number of columns in a primary key should be kept within the threshold", the threshold should be a parameter whose default value is 2.
You should follow the following logic:
1. For "create table ..." statement, if the primary key constraint has columns more than threshold, report a violation
2. For "alter table ... add primary key ..." statement, if the primary key constraint has columns more than threshold, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00023(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxColumnCount := param.Int()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."
		constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintPrimaryKey)
		if nil != constraint {
			//this is a table primary key definition
			if len(constraint.Keys) > maxColumnCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00023, maxColumnCount)
			}
		}
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table ... add constraint..."
			constraint := spec.Constraint
			if nil != spec.Constraint && spec.Constraint.Tp == ast.ConstraintPrimaryKey {
				//"alter table ... add primary key ..."
				if len(constraint.Keys) > maxColumnCount {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00023, maxColumnCount)
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
