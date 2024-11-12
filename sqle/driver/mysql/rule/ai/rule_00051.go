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
	SQLE00051 = "SQLE00051"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00051,
			Desc:       "禁止主键使用自增",
			Annotation: "后期维护相对不便，过于依赖数据库自增机制达到全局唯一，不易拆分，容易造成主键冲突",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "禁止主键使用自增",
		AllowOffline: false,
		Func:         RuleSQLE00051,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00051): "在 MySQL 中，禁止主键使用自增."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句，检查以下条件：
   1. 如果任意字段被定义为 PRIMARY KEY 且同时使用 AUTO_INCREMENT，则报告违反规则。

2. 对于 "ALTER TABLE..." 语句，检查以下条件：
   1. 如果ADD/MODIFY/CHANGE操作涉及将某字段设置为 PRIMARY KEY 且同时使用 AUTO_INCREMENT，则报告违反规则。
   2. 如果MODIFY/CHANGE操作涉及将某字段设置为 PRIMARY KEY，但没有显式使用AUTO_INCREMENT，则获取操作表信息提取当前列中如果有使用AUTO_INCREMENT，则报告违反规则。
   3. 如果MODIFY/CHANGE操作涉及将某字段设置为 AUTO_INCREMENT，但没有显式设置PRIMARY KEY，则获取操作表信息提取当前列如果是PRIMARY KEY，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00051(input *rulepkg.RuleHandlerInput) error {

	checkViolation := func(stmt *ast.AlterTableStmt, currentCol string,
		createTableStmt *ast.CreateTableStmt, isPrimaryKey bool, isAutoIncrement bool) bool {
		var err error
		if isPrimaryKey && isAutoIncrement {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00051)
			return true
		} else if isPrimaryKey {
			// 当前设置为主键，在线检查原来是否有isAutoIncrement
			if createTableStmt == nil {
				createTableStmt, err = util.GetCreateTableStmt(input.Ctx, stmt.Table)
				if err != nil {
					log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
					return true
				}
			}

			for _, col := range createTableStmt.Cols {
				if strings.EqualFold(col.Name.Name.String(), currentCol) && util.IsColumnAutoIncrement(col) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00051)
					return true
				}
			}
		} else if isAutoIncrement {
			// 当前设置为isAutoIncrement，检查原来是不是主键
			if createTableStmt == nil {
				createTableStmt, err = util.GetCreateTableStmt(input.Ctx, stmt.Table)
				if err != nil {
					log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
					return true
				}
			}
			for _, col := range createTableStmt.Cols {
				if strings.EqualFold(col.Name.Name.String(), currentCol) && util.IsColumnPrimaryKey(col) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00051)
					return true
				}
			}
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnPrimaryKey(col) && util.IsColumnAutoIncrement(col) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00051)
				return nil
			}
		}
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				for _, key := range constraint.Keys {
					colName := util.GetIndexColName(key)
					for _, col := range stmt.Cols {
						if strings.EqualFold(util.GetColumnName(col), colName) && util.IsColumnAutoIncrement(col) {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00051)
							return nil
						}
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		var createTableStmt *ast.CreateTableStmt
		for _, spec := range stmt.Specs {

			isPrimaryKey := false
			isAutoIncrement := false

			// ADD col ..
			if util.IsAlterTableCommand(spec, ast.AlterTableAddColumns) {
				newCol := spec.NewColumns[0]
				for _, op := range newCol.Options {
					if op.Tp == ast.ColumnOptionAutoIncrement {
						isAutoIncrement = true
					} else if op.Tp == ast.ColumnOptionPrimaryKey {
						isPrimaryKey = true
					}
				}
				if checkViolation(stmt, newCol.Name.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement) {
					return nil
				}
				// Modify/Change col ..
			} else if util.IsAlterTableCommand(spec, ast.AlterTableModifyColumn) || util.IsAlterTableCommand(spec, ast.AlterTableChangeColumn) {
				newCol := spec.NewColumns[0]
				for _, op := range newCol.Options {
					if op.Tp == ast.ColumnOptionAutoIncrement {
						isAutoIncrement = true
					} else if op.Tp == ast.ColumnOptionPrimaryKey {
						isPrimaryKey = true
					}
				}
				if checkViolation(stmt, newCol.Name.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement) {
					return nil
				}
				// ADD PrimaryKey
			} else if util.IsAlterTableCommand(spec, ast.AlterTableAddConstraint) {
				if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
					isPrimaryKey = true
					for _, key := range spec.Constraint.Keys {
						if checkViolation(stmt, key.Column.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement) {
							return nil
						}
					}

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
