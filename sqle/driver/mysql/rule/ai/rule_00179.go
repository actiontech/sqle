package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00179 = "SQLE00179"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00179,
			Desc:       "对于MySQL的DML, 避免隐式数据类型转换的SQL查询",
			Annotation: "确保WHERE子句中用于索引列的条件字段与索引列的数据类型一致。不一致的数据类型会导致执行计划存在隐式类型转换操作。这种转换不仅增加CPU负担，还可能使得原本高效的索引无法使用，导致查询性能显著下降。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 避免隐式数据类型转换的SQL查询",
		AllowOffline: false,
		Func:    RuleSQLE00179,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00179): "For dml, SQL queries with implicit data type conversions are prohibited".
You should follow the following logic:
1. For "select..." Statement, check the warning information of the execution plan of the SQL statement, and report the violation of the rule if the keyword "due to type or collation conversion on field" appears in the warning information. The warning information of the execution plan should be the information obtained online.
2. For "update..." Statement, performs the same checking process as above.
3. For "delete..." Statement, performs the same checking process as above.
4. For "insert..." Statement, performs the same checking process as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00179(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		// "select...", "insert...", "update...", "delete..."
		plan, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}
		for _, warning := range plan.Warnings {
			// "due to type or collation conversion on field"
			if strings.Contains(warning.Message, "due to type or collation conversion on field") {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00179)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
