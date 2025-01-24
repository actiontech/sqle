package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00088 = "SQLE00088"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00088,
			Desc:         plocale.Rule00088Desc,
			Annotation:   plocale.Rule00088Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00088Message,
		Func:    RuleSQLE00088,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00088): "在 MySQL 中，INSERT 语句必须指定COLUMN."
您应遵循以下逻辑：
1. 对于 insert... 语句，检查 INSERT 语句的目标表中是否显式指定了列名。如果未指定列名，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00088(input *rulepkg.RuleHandlerInput) error {
	// 记录违反规则的表名
	violateTables := []string{}

	// 检查语法树节点类型
	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		// 检查是否显式指定了列名
		if len(stmt.Columns) == 0 {
			// 获取插入的表名
			tableNames := util.GetTableNames(stmt)
			for _, tableName := range tableNames {
				violateTables = append(violateTables, tableName.Name.String())
			}
		}
	}

	// 如果存在违反规则的表，则将表名连接起来，并添加到检查结果中
	if len(violateTables) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00088)
	}

	return nil
}

// 规则函数实现结束

// ==== Rule code end ====
