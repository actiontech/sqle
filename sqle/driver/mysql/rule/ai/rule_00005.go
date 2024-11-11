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
			Desc:       "在 MySQL 中, 避免复合索引中包含过多字段",
			Annotation: "在设计复合索引过程中，每增加一个索引字段，都会使索引的大小线性增加，从而占用更多的磁盘空间，且增加索引维护的开销。尤其是在数据频繁变动的环境中，这会显著增加数据库的维护压力。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "复合索引内字段个数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 避免复合索引中包含过多字段",
		AllowOffline: true,
		Func:         RuleSQLE00005,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00005): "在 MySQL 中，避免复合索引中包含过多字段.默认参数描述: 复合索引内字段个数, 默认参数值: 5"
您应遵循以下逻辑：
1. 对于 CREATE TABLE 语句，检查以下条件：
   1. 定义一个字段集合。
   2. 解析语法树，提取语句中 key 或 index 定义的字段，并将其加入集合。
   3. 获取集合中字段的数量。
   4. 如果字段数量大于规则变量值，则报告违反规则。

2. 对于 ALTER TABLE 语句，检查以下条件：
   1. 定义一个字段集合。
   2. 解析语法树，提取语句中 key 或 index 定义的字段，并将其加入集合。
   3. 获取集合中字段的数量。
   4. 如果字段数量大于规则变量值，则报告违反规则。

3. 对于 CREATE INDEX 语句，检查以下条件：
   1. 定义一个字段集合。
   2. 解析语法树，提取语句中 index 定义的字段，并将其加入集合。
   3. 获取集合中字段的数量。
   4. 如果字段数量大于规则变量值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
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

// 规则函数实现结束
// ==== Rule code end ====
