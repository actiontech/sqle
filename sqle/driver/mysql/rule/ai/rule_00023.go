package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00023 = "SQLE00023"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00023,
			Desc:       plocale.Rule00023Desc,
			Annotation: plocale.Rule00023Annotation,
			Category:   plocale.RuleTypeDDLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagIndex.ID, plocale.RuleTagColumn.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID, plocale.RuleTagIntegrity.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "2",
				Desc:  plocale.Rule00023Params1,
				Type:  params.ParamTypeInt,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
			Version:      2,
		},
		Message: plocale.Rule00023Message,
		Func:    RuleSQLE00023,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00023): "In table definition, the number of columns in a primary key should be kept within the threshold", the threshold should be a parameter whose default value is 2.
You should follow the following logic:
1. For "create table ..." statement, if the primary key constraint has columns more than threshold, report a violation
2. For "alter table ... add primary key ..." statement, if the primary key constraint has columns more than threshold, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00023(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxColumnCount := param.Int()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."
		constraint := util.GetTableConstraint(stmt.Constraints, ast.ConstraintPrimaryKey)
		if nil != constraint {
			//this is a table primary key definition
			if len(constraint.Keys) > maxColumnCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00023, maxColumnCount)
			}
		}
	case *ast.AlterTableStmt:
		// "alter table"
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table ... add constraint..."
			constraint := spec.Constraint
			if nil != spec.Constraint && spec.Constraint.Tp == ast.ConstraintPrimaryKey {
				//"alter table ... add primary key ..."
				if len(constraint.Keys) > maxColumnCount {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00023, maxColumnCount)
				}
			}
		}
	}
	return nil
}

// ==== Rule code end ====
