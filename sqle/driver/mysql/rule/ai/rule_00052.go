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
	SQLE00052 = "SQLE00052"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00052,
			Desc:       plocale.Rule00052Desc,
			Annotation: plocale.Rule00052Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID, plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00052Message,
		Func:    RuleSQLE00052,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00052): "在 MySQL 中，建议主键使用自增."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句，检查以下条件：
   1. 如果任意字段被定义为 PRIMARY KEY 但不使用AUTO_INCREMENT，则报告违反规则。

2. 对于 "ALTER TABLE..." 语句，检查以下条件：
   1. 如果是ADD操作涉及将某字段设置为 PRIMARY KEY 但不使用AUTO_INCREMENT，则报告违反规则。
   2. 如果是MODIFY/CHANGE操作涉及将某字段设置为 PRIMARY KEY 但不使用AUTO_INCREMENT，且获取操作表信息提取当前列也没有AUTO_INCREMENT，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00052(input *rulepkg.RuleHandlerInput) error {
	checkViolation := func(stmt *ast.AlterTableStmt, currentCol string,
		createTableStmt *ast.CreateTableStmt, isPrimaryKey bool, isAutoIncrement bool, isOnline bool) bool {
		var err error
		if isPrimaryKey {
			if isOnline {
				// 在线查询, 原来是否有AutoIncrement
				if createTableStmt == nil {
					createTableStmt, err = util.GetCreateTableStmt(input.Ctx, stmt.Table)
					if err != nil {
						log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
						return true
					}
				}

				for _, col := range createTableStmt.Cols {
					if strings.EqualFold(col.Name.Name.String(), currentCol) && !util.IsColumnAutoIncrement(col) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00052)
						return true
					}
				}
			} else if !isAutoIncrement {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00052)
				return true
			}
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnPrimaryKey(col) && !util.IsColumnAutoIncrement(col) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00052)
				return nil
			}
		}
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				for _, key := range constraint.Keys {
					colName := util.GetIndexColName(key)
					for _, col := range stmt.Cols {
						if strings.EqualFold(util.GetColumnName(col), colName) && !util.IsColumnAutoIncrement(col) {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00052)
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
				if checkViolation(stmt, newCol.Name.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement, false) {
					return nil
				}
				// Modify/Change col ..
			} else if util.IsAlterTableCommand(spec, ast.AlterTableModifyColumn) || util.IsAlterTableCommand(spec, ast.AlterTableChangeColumn) {
				newCol := spec.NewColumns[0]
				isOnline := true
				for _, op := range newCol.Options {
					if op.Tp == ast.ColumnOptionAutoIncrement {
						isAutoIncrement = true
						isOnline = false
					} else if op.Tp == ast.ColumnOptionPrimaryKey {
						isPrimaryKey = true
					}
				}
				if checkViolation(stmt, newCol.Name.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement, isOnline) {
					return nil
				}
				// ADD PrimaryKey
			} else if util.IsAlterTableCommand(spec, ast.AlterTableAddConstraint) {
				if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
					isPrimaryKey = true
					for _, key := range spec.Constraint.Keys {
						if checkViolation(stmt, key.Column.Name.String(), createTableStmt, isPrimaryKey, isAutoIncrement, true) {
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

// ==== Rule code end ====
