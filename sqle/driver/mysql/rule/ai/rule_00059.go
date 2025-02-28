package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00059 = "SQLE00059"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00059,
			Desc:       plocale.Rule00059Desc,
			Annotation: plocale.Rule00059Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID, plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "5",
				Desc:  plocale.Rule00059Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00059Message,
		Func:    RuleSQLE00059,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00059): "For MySQL DDL, modifying large table field types is prohibited.".
You should follow the following logic:
1. For "alter table... modify..." Statement, check whether the table size in the alter target table is greater than the rule threshold variable, if it is greater than the specified threshold, If it does, report a violation. The size of the table can be determined by 'select round((index_length+data_length)/1024/1024/1024) 'size_GB' from information_schema.tables where table_name=' table name ' 'to get, the table size needs to be obtained online,.
2. For "alter table... change..." Statement to perform a check similar to the one described above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00059(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxSize, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s must be an integer", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table"

		if len(util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn)) > 0 {
			// "alter table... change/alter table... modify"

			// get the size of table
			size, err := util.GetTableSizeMB(input.Ctx, stmt.Table.Name.String())
			if err != nil {
				log.NewEntry().Errorf("get table size failed, sqle: %v, error: %v", input.Node.Text(), err)
				return nil
			}

			if size > int64(maxSize*1024) {
				// report rule violation
				rulepkg.AddResult(input.Res, input.Rule, SQLE00059, maxSize)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
