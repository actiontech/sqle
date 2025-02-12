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
	SQLE00026 = "SQLE00026"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00026,
			Desc:       plocale.Rule00026Desc,
			Annotation: plocale.Rule00026Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagCorrection.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelNotice,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00026Message,
		Func:    RuleSQLE00026,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00026): "在 MySQL 中，整数字段建议指定最大显示宽度."
您应遵循以下逻辑：
1. 对于"CREATE TABLE..."语句， 如果存在INT、INTEGER、TINYINT、SMALLINT、MEDIUMINT、BIGINT类型的字段，并且为其指定长度，并且不包含关键词：zerofill，则报告违反规则。
  - 解析语法树以识别CREATE TABLE语句。
  - 检查字段定义部分，确认字段类型是否为INT、INTEGER、TINYINT、SMALLINT、MEDIUMINT、BIGINT。
  - 确认这些字段是否指定了长度。
  - 检查字段定义中是否包含zerofill关键字。
2. 对于"ALTER TABLE..." 语句，执行上述相同的检查。
  - 解析语法树以识别ALTER TABLE语句。
  - 检查字段定义部分，确认字段类型是否为INT、INTEGER、TINYINT、SMALLINT、MEDIUMINT、BIGINT。
  - 确认这些字段是否指定了长度。
  - 检查字段定义中是否包含zerofill关键字。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00026(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	intTypes := []byte{mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong, mysql.TypeShort, mysql.TypeTiny}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, intTypes...) {
				if util.GetColumnWidth(col) > 0 {
					if !mysql.HasZerofillFlag(col.Tp.Flag) {
						violateColumns = append(violateColumns, col)
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, intTypes...) {
					if util.GetColumnWidth(col) > 0 {
						if !mysql.HasZerofillFlag(col.Tp.Flag) {
							violateColumns = append(violateColumns, col)
						}
					}
				}
			}
		}
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00026)
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
