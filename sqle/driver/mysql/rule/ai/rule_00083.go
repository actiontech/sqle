package ai

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00083 = "SQLE00083"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00083,
			Desc:         plocale.Rule00083Desc,
			Annotation:   plocale.Rule00083Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00083Message,
		Func:    RuleSQLE00083,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00083): "在 MySQL 中，不建议对表进行索引跳跃扫描."
您应遵循以下逻辑：
1. 对于DML语句：
   1. 检查是否包含SELECT子句，若包含，继续。
   2. 检查是否存在GROUP BY或DISTINCT，若不存在，继续。
   3. 检查FROM子句是否仅涉及一张表，若是，继续。
   4. 连接数据库，验证SELECT子句的字段是否为该表联合索引的部分或全部字段，若是，继续。
   5. 使用辅助函数GetExecutionPlan获取SELECT语句的执行计划，检查是否包含索引跳跃扫描的节点，若包含，报告违反规则。

2. 对于WITH语句，执行与DML语句相同的检查。

3. 对于UNION语句，对每个SELECT子句执行与DML语句相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00083(input *rulepkg.RuleHandlerInput) error {

	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}
	if len(util.GetSelectStmt(input.Node)) > 0 {
		explain, err := util.GetExecutionPlan(input.Ctx, input.Node.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
			return err
		}
		for _, record := range explain.Plan {
			if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingIndexForSkipScan) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00083)
				return nil
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
