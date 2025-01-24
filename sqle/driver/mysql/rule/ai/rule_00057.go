package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00057 = "SQLE00057"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00057,
			Desc:         plocale.Rule00057Desc,
			Annotation:   plocale.Rule00057Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00057Message,
		Func:    RuleSQLE00057,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00057): "在 MySQL 中，必须使用INNODB数据库引擎."
您应遵循以下逻辑：
1. 对于 CREATE TABLE 语句，执行以下检查：
   1. 使用辅助函数util.GetTableOption检查语法树中是否包含 ENGINE 节点。
   2. 如果未包含 ENGINE 节点，使用函数 input.Ctx.GetSchemaEngine 获取参数 default_storage_engine 值是否为innodb，若不是innodb, ，则报告违反规则。
   3. 如果包含 ENGINE 节点，判断是否为innodb，若不是innodb ，则报告违反规则。

2. 对于 ALTER TABLE 语句，执行以下检查：
   1. 使用辅助函数util.GetTableOption检查语法树中是否包含 ENGINE 节点，包含则进入下一步。
   2. 如果ENGINE不等于InnoDB，而是其他引擎（如 MyISAM、CSV、ARCHIVE 等），则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00057(input *rulepkg.RuleHandlerInput) error {
	expectEngine := "InnoDB"
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE TABLE 语句中是否包含 ENGINE 选项
		engineOption := util.GetTableOption(stmt.Options, ast.TableOptionEngine)
		if engineOption == nil {
			// 若未包含 ENGINE 选项，获取 default_storage_engine 参数
			defaultEngine, err := input.Ctx.GetSchemaEngine(stmt.Table, stmt.Table.Schema.L)
			if err != nil {
				log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
				return err
			}
			// 验证 default_storage_engine 是否为 InnoDB
			if !strings.EqualFold(defaultEngine, expectEngine) {
				// default_storage_engine 不是 InnoDB，记录该 SQL
				rulepkg.AddResult(input.Res, input.Rule, SQLE00057)
				return nil
			}
		} else {
			// ENGINE 选项存在，检查是否为 InnoDB
			if !strings.EqualFold(engineOption.StrValue, expectEngine) {
				// ENGINE 不是 InnoDB 或包含其他引擎节点，记录该 SQL
				rulepkg.AddResult(input.Res, input.Rule, SQLE00057)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中是否包含 ENGINE 选项
		tableOptions := []*ast.TableOption{}
		for _, spec := range stmt.Specs {
			if len(spec.Options) > 0 {
				tableOptions = append(tableOptions, spec.Options...)
			}
		}
		engineOption := util.GetTableOption(tableOptions, ast.TableOptionEngine)
		if engineOption != nil {
			// ENGINE 选项存在，检查是否为 InnoDB
			if !strings.EqualFold(engineOption.StrValue, expectEngine) {
				// ENGINE 不是 InnoDB 或包含其他引擎节点，记录该 SQL
				rulepkg.AddResult(input.Res, input.Rule, SQLE00057)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
