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
	SQLE00175 = "SQLE00175"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00175,
			Desc:       "对于MySQL的DML, 避免不必要的索引扫描合并",
			Annotation: "索引合并说明一个查询同时使用了多个索引，增加了更多IO操作，特别是在数据量大的情况下执行效率比复合索引明显更多。此外，索引合并操作可能消耗更多CPU和内存资源，以及较长的查询响应时间。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DML, 避免不必要的索引扫描合并",
		AllowOffline: false,
		Func:    RuleSQLE00175,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
