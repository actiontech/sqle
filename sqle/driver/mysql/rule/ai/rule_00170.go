package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00170 = "SQLE00170"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00170,
			Desc:       "在 MySQL 中, 避免缩短字段长度",
			Annotation: "修改字段长度值低于现有字段长度值，如果该字段现有数据超出设定后的长度，会造成语句执行报错",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免缩短字段长度. 字段 %s 的新长度 %d 小于当前最大长度 %d",
		AllowOffline: false,
		Func:         RuleSQLE00170,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00170): "在 MySQL 中，避免缩短字段长度."
您应遵循以下逻辑：
1. 对于 "ALTER TABLE... MODIFY ..." 语句：
   1. 提取字段名称及其类型和使用辅助函数GetColumnWidth获取其长度。
   2. 登录数据库。
   3. 使用SQL(select max(char_length(col_name)) "max_length" from table_name)获取字段当前存储值的最大长度。
   4. 比较新定义的字段长度与查询结果中的最大长度：
      - 如果新长度小于最大长度，则报告违反规则。

2. 对于 "ALTER TABLE ... CHANGE ..." 语句，执行与上述相同检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00170(input *rulepkg.RuleHandlerInput) error {
	// 确保输入节点为 ALTER TABLE 语句
	alterStmt, ok := input.Node.(*ast.AlterTableStmt)
	if !ok {
		return nil
	}

	// 获取所有 MODIFY 和 CHANGE 类型的 AlterTableSpec
	modifyChangeSpecs := util.GetAlterTableCommandsByTypes(alterStmt, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn)
	if len(modifyChangeSpecs) == 0 {
		return nil
	}

	// 遍历每个 MODIFY 或 CHANGE 规范
	for _, spec := range modifyChangeSpecs {
		for _, newCol := range spec.NewColumns {
			// 提取字段名称
			columnName := util.GetColumnName(newCol)

			// 提取新字段长度
			newLength := util.GetColumnWidth(newCol)

			// 检查是否成功提取新字段长度
			if newLength == 0 {
				log.NewEntry().Errorf("无法提取字段 %s 的新长度", columnName)
				continue
			}

			// 获取当前字段的最大存储长度
			currentMaxLength, err := util.GetCurrentMaxColumnWidth(input.Ctx, alterStmt.Table, columnName)
			if err != nil {
				// 记录日志并继续处理下一个字段
				log.NewEntry().Errorf("获取字段 %s 当前最大长度失败: %v", columnName, err)
				continue
			}

			// 比较新定义的字段长度与当前最大长度
			if newLength < currentMaxLength {
				// 报告违反规则
				rulepkg.AddResult(input.Res, input.Rule, SQLE00170, columnName, newLength, currentMaxLength)
			}
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====