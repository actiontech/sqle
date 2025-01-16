package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00010 = "SQLE00010"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00010,
			Desc:       "禁止进行删除主键的操作",
			Annotation: "在MySQL中删除已有主键代价高昂，极易引起业务阻塞、故障；开启该规则，SQLE将提醒删除主键为高危操作",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "禁止进行删除主键的操作",
		AllowOffline: true,
		Func:    RuleSQLE00010,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
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
