package ai

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00175 = "SQLE00175"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00175,
			Desc:         plocale.Rule00175Desc,
			Annotation:   plocale.Rule00175Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00175Message,
		Func:    RuleSQLE00175,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00175): "For dml, index scan merges are prohibited".
You should follow the following logic:
1. For "select..." statement, check the output of the execution plan. If type is index_merge, report a violation. The execution plan should be the information obtained online.
2. For "insert..." Statement, perform the same check as above.
3. For "union..." Statement, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00175(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt:
		// "select..." "insert..." "union..."
		executionPlan, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}

		for _, record := range executionPlan.Plan {
			if record.Type == executor.ExplainRecordAccessTypeIndexMerge {
				// index merge
				rulepkg.AddResult(input.Res, input.Rule, SQLE00175)
				return nil
			}
		}
	}
	return nil
}

// ==== Rule code end ====
