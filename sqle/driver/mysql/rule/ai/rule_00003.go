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
	SQLE00003 = "SQLE00003"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00003,
			Desc:         plocale.Rule00003Desc,
			Annotation:   plocale.Rule00003Annotation,
			Category:     plocale.RuleTypeDMLConvention,
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00003Message,
		Func:    RuleSQLE00003,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00003): "在 MySQL 中，建议为组成索引的字段添加非空约束，并配置合理的default值."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句：
   1. 提取所有索引字段。
   2. 检查每个索引字段是否包含 NOT NULL 约束。
   3. 若任何索引字段缺少 NOT NULL 约束，则报告违反规则。
   4. 检查每个索引字段是否设置了合理的 DEFAULT 值，若设置值不合理则则报告违反规则。
2. 对于 "CREATE INDEX..." 语句：
   1. 提取所有索引字段。
   2. 使用辅助函数 GetCreateTableStmt 获取索引字段所在表的结构信息。
   3. 检查每个索引字段是否包含 NOT NULL 约束, 若缺少 NOT NULL 约束，则报告违反规则。
   4. 检查每个索引字段是否设置了合理的 DEFAULT 值，若设置值不合理则则报告违反规则。
3. 对于 "ALTER TABLE...ADD INDEX..." 语句：
   1. 提取所有新增索引字段。
   2. 使用辅助函数 GetCreateTableStmt 获取索引字段所在表的结构信息。
   3. 检查每个新增索引字段是否包含 NOT NULL 约束, 若缺少 NOT NULL 约束，则报告违反规则。
   4. 检查每个新增索引字段是否设置了合理的 DEFAULT 值，若设置值不合理则则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00003(input *rulepkg.RuleHandlerInput) error {

	// 定义匿名函数以提取表约束中的索引字段
	extractIndexesFromConstraints := func(constraints []*ast.Constraint) [][]string {
		indexes := [][]string{}

		for _, constraint := range constraints {
			indexCols := []string{}
			for _, key := range constraint.Keys {
				colName := util.GetIndexColName(key)
				if colName != "" {
					indexCols = append(indexCols, colName)
				}
			}
			if len(indexCols) > 0 {
				indexes = append(indexes, indexCols)
			}
		}

		return indexes
	}
	// 检查每个索引字段的约束 是否有NOT NULL 约束
	checkCol := func(createStmt *ast.CreateTableStmt, indexes [][]string) bool {
		for _, indexCols := range indexes {
			for _, colName := range indexCols {
				// 获取列定义
				var columnDef *ast.ColumnDef
				for _, col := range createStmt.Cols {
					if strings.EqualFold(util.GetColumnName(col), colName) {
						columnDef = col
						break
					}
				}
				if columnDef == nil {
					continue
				}

				// 检查 NOT NULL 约束
				if !util.IsColumnHasOption(columnDef, ast.ColumnOptionNotNull) {
					return true
				}

				// 检查合理的 DEFAULT 值
				defaultOption := util.GetColumnOption(columnDef, ast.ColumnOptionDefaultValue)
				if defaultOption != nil {
					// 这里假设合理的 DEFAULT 值是非空且具体化的值，可以根据实际需求调整
					if util.IsOptionValIsNull(defaultOption) { // DEFAULT NULL 时，不合理
						return true
					}
				} else { // 不设置DEFAULT
					return true
				}
			}
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 处理 "CREATE TABLE..." 语句
		indexes := [][]string{}

		// 从列定义中提取索引字段
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) { // 主键默认not null -- 则忽略：util.IsColumnPrimaryKey(col)
				indexes = append(indexes, []string{util.GetColumnName(col)})
			}
		}
		// 从表约束中提取索引字段
		indexes = append(indexes, extractIndexesFromConstraints(util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...))...)
		if checkCol(stmt, indexes) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00003)
			return nil
		}
	case *ast.CreateIndexStmt:
		// 处理 "CREATE INDEX..." 语句

		// 定义匿名函数以从 CreateIndexStmt 中提取列名
		extractColumnNamesFromIndexStmt := func(stmt *ast.CreateIndexStmt) []string {
			columns := []string{}
			for _, key := range stmt.IndexPartSpecifications {
				colName := util.GetIndexColName(key)
				if colName != "" {
					columns = append(columns, colName)
				}
			}
			return columns
		}

		// 提取新的索引字段
		newIndexCols := extractColumnNamesFromIndexStmt(stmt)
		indexes := [][]string{}
		indexes = append(indexes, newIndexCols)
		if len(indexes) > 0 {
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
			if err != nil {
				log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
				return err
			}
			if checkCol(createTableStmt, indexes) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00003)
				return nil
			}
		}

	case *ast.AlterTableStmt:
		// 处理 "ALTER TABLE...ADD INDEX..." 语句
		var constraints []*ast.Constraint
		for _, spec := range stmt.Specs {
			if spec.Tp != ast.AlterTableAddConstraint {
				continue
			}
			constraints = append(constraints, spec.Constraint)
		}
		indexes := [][]string{}
		indexes = append(indexes, extractIndexesFromConstraints(util.GetTableConstraints(constraints, util.GetIndexConstraintTypes()...))...)
		if len(indexes) > 0 {
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
			if err != nil {
				log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
				return err
			}
			if checkCol(createTableStmt, indexes) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00003)
				return nil
			}
		}
	default:
		// 其他类型的语句不处理
		return nil
	}

	return nil
}

// ==== Rule code end ====
