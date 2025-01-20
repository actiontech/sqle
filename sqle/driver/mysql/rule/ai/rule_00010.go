package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00010 = "SQLE00010"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00010,
			Desc:         plocale.Rule00010Desc,
			Annotation:   plocale.Rule00010Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00010Message,
		Func:    RuleSQLE00010,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00010): "In DDL, deleting primary key is prohibited".
You should follow the following logic:
1. For "alter table ... drop primary key ..." statement, report a violation
2. For "drop index ..." statement, if the index is primary index, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00010(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// "alter table"
		for range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableDropPrimaryKey) {
			// "alter table drop primary key"
			rulepkg.AddResult(input.Res, input.Rule, SQLE00010)
			return nil
		}
	case *ast.DropIndexStmt:
		// "drop index"
		if strings.EqualFold(stmt.IndexName, "primary") {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00010)
		}
	}
	return nil
}

// ==== Rule code end ====
