package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00022 = "SQLE00022"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00022,
			Desc:       plocale.Rule00022Desc,
			Annotation: plocale.Rule00022Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelNotice,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "40",
				Desc:  plocale.Rule00022Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00022Message,
		Func:    RuleSQLE00022,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00022): "In table definition, the count of table columns should be within threshold", the threshold should be a parameter whose default value is 40.
You should follow the following logic:
1. For "create table ..." statement, check column count should be within threshold, otherwise, report a violation.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00022(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxColumnCount, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be a number", param.Value)
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."
		if len(stmt.Cols) > maxColumnCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00022, maxColumnCount)
		}
	default:
		return nil
	}
	return nil
}

// ==== Rule code end ====
