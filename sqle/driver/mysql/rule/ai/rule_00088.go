package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00088 = "SQLE00088"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00088,
			Desc:       "INSERT 语句必须指定COLUMN",
			Annotation: "当表结构发生变更，INSERT请求不明确指定列名，会发生插入数据不匹配的情况；建议开启此规则，避免插入结果与业务预期不符",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "INSERT 语句必须指定COLUMN",
		AllowOffline: true,
		Func:         RuleSQLE00088,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
