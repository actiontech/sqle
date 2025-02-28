package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00068 = "SQLE00068"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00068,
			Desc:       plocale.Rule00068Desc,
			Annotation: plocale.Rule00068Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00068Message,
		Func:    RuleSQLE00068,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00068): "在 MySQL 中，禁止使用TIMESTAMP字段."
您应遵循以下逻辑：
1. 对于 "CREATE TABLE..." 语句，执行以下检查：
   1. 使用辅助函数 IsColumnTypeEqual 检查语法节点中是否包含 TIMESTAMP 字段定义。
   若包含，则报告违反规则。

2. 对于 "ALTER TABLE..." 语句，执行以下检查：
   1. 使用辅助函数 IsColumnTypeEqual 检查语法节点中是否包含 TIMESTAMP 字段定义。
   若包含，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00068(input *rulepkg.RuleHandlerInput) error {
	// 存储所有违反规则的列
	violateColumns := []*ast.ColumnDef{}

	switch stmt := input.Node.(type) {
	// 检查 "CREATE TABLE" 语句
	case *ast.CreateTableStmt:
		// 遍历所有列定义
		for _, col := range stmt.Cols {
			// 使用辅助函数检查列类型是否为 TIMESTAMP
			if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
				violateColumns = append(violateColumns, col)
			}
		}

	// 检查 "ALTER TABLE" 语句
	case *ast.AlterTableStmt:
		// 获取所有 ADD COLUMN、CHANGE COLUMN 和 MODIFY COLUMN 操作
		alterTypes := []ast.AlterTableType{
			ast.AlterTableAddColumns,
			ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn,
		}
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, alterTypes...) {
			// 遍历每个新的列定义
			for _, col := range spec.NewColumns {
				// 使用辅助函数检查列类型是否为 TIMESTAMP
				if util.IsColumnTypeEqual(col, mysql.TypeTimestamp) {
					violateColumns = append(violateColumns, col)
				}
			}
		}

	default:
		// 非 "CREATE TABLE" 或 "ALTER TABLE" 语句，不处理
		return nil
	}

	// 如果存在任何违反规则的列，则报告违规
	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00068, util.JoinColumnNames(violateColumns))
	}

	return nil
}

// ==== Rule code end ====
