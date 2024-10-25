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
	SQLE00005 = "SQLE00005"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00005,
			Desc:       "对于MySQL的DDL, 复合索引的列数量不建议超过阈值",
			Annotation: "复合索引会根据索引列数创建对应组合的索引，列数越多，创建的索引越多，每个索引都会增加磁盘空间的开销，同时增加索引维护的开销；具体规则阈值可以根据业务需求调整，默认值：3",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeIndexingConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "索引列数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "对于MySQL的DDL, 复合索引的列数量不建议超过阈值. 阈值: %v",
		AllowOffline: true,
		Func:    RuleSQLE00005,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00005): "For table creation and index creation statements, the number of columns included in the composite index should be within threshold, the threshold should be a parameter whose default value is 5.".
You should follow the following logic:
1. For the "CREATE TABLE ..." statements, check if the composite index includes more than three columns. If it does, report a violation.
2. For the  "CREATE INDEX ..." statements, check if the composite index includes more than three columns. If it does, report a violation.
3. For the  "ALTER TABLE ... ADD INDEX ..." statements, check if the composite index includes more than three columns. If it does, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00005(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxColumnCount := param.Int()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		constraints := util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...)

		for _, constraint := range constraints {
			// the table is created with composite index
			if len(constraint.Keys) > maxColumnCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00005, maxColumnCount)
				return nil
			}
		}
	case *ast.CreateIndexStmt:
		// "create index..."
		if len(stmt.IndexPartSpecifications) > maxColumnCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00005, maxColumnCount)
			return nil
		}
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			constraints := util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...)

			for _, constraint := range constraints {
				// the table is created with composite index
				if len(constraint.Keys) > maxColumnCount {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00005, maxColumnCount)
					return nil
				}
			}

		}
	}
	return nil
}

// ==== Rule code end ====
