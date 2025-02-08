package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00042 = "SQLE00042"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00042,
			Desc:       plocale.Rule00042Desc,
			Annotation: plocale.Rule00042Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "tmp_",
				Desc:  plocale.Rule00042Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00042Message,
		Func:    RuleSQLE00042,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00042): "在 MySQL 中，临时表必须使用固定前缀.默认参数描述: 固定前缀, 默认参数值: tmp_"
您应遵循以下逻辑：
1、检查CREATE语法节点中是否存在TEMPORARY关键词节点，如果存在，则进入下一步检查。
2、检查CREATE语法节点中是否存在TABLE关键词节点，如果存在，则进入下一步检查。
3、检查目标表名是否包含固定前缀，如果不包含，报告违反规则。
1、检查ALTER语法节点中是否存在RENAME关键词节点，如果存在，则进入下一步检查。
2、通过上下文进行检查，检查RENAME的目标对象类型是否临时表，如果是，则进入下一步检查。
3、检查目标表名是否包含固定前缀，如果不包含，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00042(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	requiredPrefix := param.String()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查是否为临时表
		if stmt.IsTemporary {
			// 检查表名是否包含固定前缀
			if !strings.HasPrefix(stmt.Table.Name.String(), requiredPrefix) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00042)
			}
		}

	case *ast.AlterTableStmt:
		// 检查是否包含RENAME关键字
		if specs := util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableRenameTable); len(specs) > 0 {
			// 检查目标对象是否为临时表
			isTemp, err := util.IsTemporaryTable(input.Ctx, stmt.Table)
			if err != nil {
				return fmt.Errorf("failed to check if table is temporary: %s", err)
			}
			if isTemp {
				// 检查新表名是否包含固定前缀
				for _, spec := range specs {
					if !strings.HasPrefix(spec.NewTable.Name.String(), requiredPrefix) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00042)
					}
				}
			}
		}
	}

	return nil
}

// ==== Rule code end ====
