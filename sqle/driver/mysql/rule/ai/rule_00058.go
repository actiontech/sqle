package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00058 = "SQLE00058"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00058,
			Desc:       "在 MySQL 中, 避免使用分区表相关功能",
			Annotation: "分区表在使用过程中存在诸多缺点，比如分区裁剪的不确定性、不支持全局分区索引、锁定粒度放大、分区前期规划较为繁杂等问题。如存在分区诉求，通常使用物理分表，即可避免分区表带来的缺点。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免使用分区表相关功能",
		AllowOffline: true,
		Func:         RuleSQLE00058,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00058): "在 MySQL 中，避免使用分区表相关功能."
您应遵循以下逻辑：
1. 检查 "CREATE TABLE ..." 语句：
   1. 句子中包含关键词：PARTITION，则报告违反规则。

2. 检查 "ALTER TABLE ..." 语句，执行与上述同样检查。

==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00058(input *rulepkg.RuleHandlerInput) error {
	// 确保 input.Node 是有效的语法树节点
	if input.Node == nil {
		return fmt.Errorf("input.Node is nil")
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE TABLE 语句中的分区定义
		if stmt.Partition != nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00058)
		}
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中的分区修改
		for _, spec := range stmt.Specs {
			if len(spec.PartitionNames) > 0 || len(spec.PartDefinitions) > 0 || spec.Partition != nil {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00058)
				break
			}
		}
	default:
		return nil
	}
	return nil
}

// ==== Rule code end ====
