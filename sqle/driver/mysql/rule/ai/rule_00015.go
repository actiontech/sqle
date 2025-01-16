package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00015 = "SQLE00015"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00015,
			Desc:       "避免库内出现多种数据库排序规则",
			Annotation: "建议库内使用一致的数据库排序规则，以确保查询性能和索引有效性，避免因排序规则不一致导致的全表扫描和数据一致性问题。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "避免库内出现多种数据库排序规则",
		AllowOffline: false,
		Func:         RuleSQLE00015,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00015): "在 MySQL 中，避免库内出现多种数据库排序规则."
您应遵循以下逻辑：
1. 对于 CREATE TABLE 语句：
   1. 检查是否在表级别指定了 COLLATION。
   2. 如果指定了表级 COLLATION，使用辅助函数 GetDatabaseOption 验证其是否与数据库默认 COLLATION 一致。
   3. 检查所有字符类型列（如 CHAR、VARCHAR、TEXT 等）是否指定了列级 COLLATION。
   4. 如果指定了列级 COLLATION，使用辅助函数 GetDatabaseOption 验证其是否与数据库默认 COLLATION 一致。
   5. 如果表级或任何列级 COLLATION 与数据库默认 COLLATION 不一致，报告规则违规。

2. 对于 ALTER TABLE 语句：
   1. 检查语句中是否包含 CONVERT TO CHARACTER SET 子句，并使用辅助函数 GetDatabaseOption 验证指定的 COLLATION 是否与数据库默认 COLLATION 一致。
   2. 如果添加或修改字符类型列（如 CHAR、VARCHAR、TEXT 等），检查是否指定了 COLLATION，并使用辅助函数 GetDatabaseOption 验证其是否与数据库默认 COLLATION 一致。
   3. 如果指定的 COLLATION 与数据库默认 COLLATION 不一致，报告规则违规。
==== Prompt end ====
*/

// ==== Rule code start ====

// 规则函数实现开始
func RuleSQLE00015(input *rulepkg.RuleHandlerInput) error {

	getDefaultCollation := func(defaultCollation string, table *ast.TableName) (string, error) {
		if defaultCollation != "" {
			return defaultCollation, nil
		}
		defaultCollation, err := input.Ctx.GetCollationDatabase(table, "")
		if err != nil {
			log.NewEntry().Errorf("GetCollationDatabase, fail err: %v", err)
			return "", err
		}
		return defaultCollation, nil
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		var defaultCollation string
		var err error

		// "create table ..."
		// get collate in column definition
		for _, col := range stmt.Cols {
			if option := util.GetColumnOption(col, ast.ColumnOptionCollate); nil != option {
				defaultCollation, err = getDefaultCollation(defaultCollation, stmt.Table)
				if err != nil {
					return err
				}
				if !strings.EqualFold(defaultCollation, option.StrValue) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00015)
					return nil
				}
			}
		}

		// get collate in table option
		if option := util.GetTableOption(stmt.Options, ast.TableOptionCollate); nil != option {
			defaultCollation, err = getDefaultCollation(defaultCollation, stmt.Table)
			if err != nil {
				return err
			}
			if !strings.EqualFold(defaultCollation, option.StrValue) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00015)
				return nil
			}
		}

	case *ast.AlterTableStmt:
		var defaultCollation string
		var err error
		// "alter table"
		for _, spec := range stmt.Specs {
			// get collate in column definition
			for _, col := range spec.NewColumns {
				if option := util.GetColumnOption(col, ast.ColumnOptionCollate); nil != option {
					defaultCollation, err = getDefaultCollation(defaultCollation, stmt.Table)
					if err != nil {
						return err
					}
					if !strings.EqualFold(defaultCollation, option.StrValue) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00015)
						return nil
					}
				}
			}

			// get collate in table option
			if option := util.GetTableOption(spec.Options, ast.TableOptionCollate); nil != option {
				defaultCollation, err = getDefaultCollation(defaultCollation, stmt.Table)
				if err != nil {
					return err
				}
				if !strings.EqualFold(defaultCollation, option.StrValue) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00015)
					return nil
				}
			}
		}
	default:
		return nil
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
