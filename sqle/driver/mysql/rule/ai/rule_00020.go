package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00020 = "SQLE00020"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00020,
			Desc:       "避免表中包含有太多的列",
			Annotation: "数据库表中字段过多会导致数据操作效率降低、数据完整性检查成本增加，以及索引维护与更新效率之间的权衡成本。对于追求事务响应和处理速度的OLTP系统，应尽量避免宽表设计，采用规范化数据模型以提升性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "40",
					Desc:  "表内列数上限",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "避免表中包含有太多的列",
		AllowOffline: false,
		Func:         RuleSQLE00020,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00020): "在 MySQL 中，避免表中包含有太多的列.默认参数描述: 表内列数上限, 默认参数值: 40"
您应遵循以下逻辑：
1. 对于“CREATE TABLE ...”语句：
   - 统计定义的字段个数。
   - 若字段个数超过预设阈值，则报告违反规则。

2. 对于“ALTER TABLE ...”语句：
   1. 使用辅助函数GetCreateTableStmt获取当前表的字段个数。
   2. 统计当前语句中的 DROP 和 ADD 操作的个数：
      - 计算字段个数的净变化（ADD 操作增加，DROP 操作减少）。
   3. 将当前表的字段个数与净变化相加。
   4. 如果结果超过预设阈值，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00020(input *rulepkg.RuleHandlerInput) error {
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param should be an integer, got: %v", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if len(stmt.Cols) > threshold {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00020)
			return nil
		}
	case *ast.AlterTableStmt:
		num := len(util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns)) - len(util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableDropColumn))
		createTable, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if err != nil {
			log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
			return err
		}
		if len(createTable.Cols)+num > threshold {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00020)
			return nil
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
