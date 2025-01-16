package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00085 = "SQLE00085"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00085,
			Desc:       "不建议对表进行全索引扫描",
			Annotation: "MySQL需要单独维护重复的索引，冗余索引增加维护成本，影响更新性能",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "不建议对表进行全索引扫描",
		AllowOffline: false,
		Func:    RuleSQLE00085,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00085): "For dml, Full index scans of tables are prohibited".
You should follow the following logic:
1. For "select..." statement, check the output of the execution plan. If type is index, report a violation. The execution plan should be the information obtained online.
2. For "union ..." statement, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00085(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt:
		// "select..."
		executionPlan, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}

		for _, record := range executionPlan.Plan {
			if record.Type == executor.ExplainRecordAccessTypeIndex {
				// full index scan
				rulepkg.AddResult(input.Res, input.Rule, SQLE00085)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
