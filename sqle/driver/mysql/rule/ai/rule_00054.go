package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00054 = "SQLE00054"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00054,
			Desc:       "建议主键字段使用BIGINT时采用无符号的BIGINT",
			Annotation: "在设计主键时若选择BIGINT时，使用无符号类型，相对于有符号类型，可以使数据库的索引性能更加优化，因为它减少了负值处理的开销，并能在某些情况下提高查询速度。特别是在系统设计初期可能无法完全预见到未来数据量的情况下，无符号数值类型（BIGINT UNSIGNED）可以有效避免因数据增长导致的溢出问题。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "建议主键字段使用BIGINT时采用无符号的BIGINT",
		AllowOffline: true,
		Func:         RuleSQLE00054,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00054): "在 MySQL 中，建议主键字段使用BIGINT时采用无符号的BIGINT."
您应遵循以下逻辑：
1. 针对 "CREATE TABLE..." 语句，执行以下检查：
   1. 使用辅助函数IsColumnPrimaryKey确认主键字段。
   2. 使用辅助函数IsColumnTypeEqual确认主键字段的数据类型为 BIGINT。
   3. 使用mysql.HasUnsignedFlag检查主键字段是否未被定义为 UNSIGNED。
   如果以上条件同时成立，则标记为规则违规。

2. 针对 "ALTER TABLE..." 语句，执行以下检查：
   1. 使用辅助函数IsColumnPrimaryKey确认主键字段。
   2. 使用辅助函数IsColumnTypeEqual确认主键字段的数据类型为 BIGINT。
   3. 使用mysql.HasUnsignedFlag检查主键字段是否未被定义为 UNSIGNED。
   如果以上条件同时成立，则标记为规则违规。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00054(input *rulepkg.RuleHandlerInput) error {
	// 初始化违规列的列表
	violateColumns := []*ast.ColumnDef{}

	// 定义BIGINT的数据类型标识，假设mysql.TypeLonglong代表BIGINT
	bigintType := mysql.TypeLonglong

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 遍历所有列定义
		for _, col := range stmt.Cols {
			// 检查是否为主键字段
			if util.IsColumnPrimaryKey(col) {
				// 检查数据类型是否为 BIGINT
				if util.IsColumnTypeEqual(col, bigintType) {
					// 检查是否未被定义为 UNSIGNED
					if !mysql.HasUnsignedFlag(col.Tp.Flag) {
						// 如果满足所有条件，记录为违规列
						rulepkg.AddResult(input.Res, input.Rule, SQLE00054)
						return nil
					}
				}
			}
		}
		constantPrimaryKey := util.GetTableConstraints(stmt.Constraints, ast.ConstraintPrimaryKey)
		if len(constantPrimaryKey) > 0 {
			for _, key := range constantPrimaryKey[0].Keys {
				for _, col := range stmt.Cols {
					if key.Column.Name.L == col.Name.Name.L {
						if util.IsColumnTypeEqual(col, bigintType) {
							if !mysql.HasUnsignedFlag(col.Tp.Flag) {
								// 如果满足所有条件，记录为违规列
								rulepkg.AddResult(input.Res, input.Rule, SQLE00054)
								return nil
							}
						}
					}
				}
			}
		}

	case *ast.AlterTableStmt:
		// 获取所有涉及列添加、修改或更改的操作
		alterSpecs := util.GetAlterTableCommandsByTypes(
			stmt,
			ast.AlterTableAddColumns,
			ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn,
		)

		for _, spec := range alterSpecs {
			// 遍历所有新的列定义
			for _, newCol := range spec.NewColumns {
				// 检查是否为主键字段
				if util.IsColumnPrimaryKey(newCol) {
					// 检查数据类型是否为 BIGINT
					if util.IsColumnTypeEqual(newCol, bigintType) {
						// 检查是否未被定义为 UNSIGNED
						if !mysql.HasUnsignedFlag(newCol.Tp.Flag) {
							// 如果满足所有条件，记录为违规列
							rulepkg.AddResult(input.Res, input.Rule, SQLE00054)
							return nil
						}
					}
				}
			}
		}

	default:
		// 非 CREATE TABLE 或 ALTER TABLE 语句，不处理
		return nil
	}

	// 如果存在任何违规列，则报告规则违规
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00054, util.JoinColumnNames(violateColumns))
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
