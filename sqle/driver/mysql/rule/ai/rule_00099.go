package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00099 = "SQLE00099"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00099,
			Desc:       "对于MySQL的DQL, 不建议使用SELECT FOR UPDATE",
			Annotation: "SELECT FOR UPDATE 会对查询结果集中每行数据都添加排他锁，其他线程对该记录的更新与删除操作都会阻塞，在高并发下，容易造成数据库大量锁等待，影响数据库查询性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "对于MySQL的DQL, 不建议使用SELECT FOR UPDATE",
		AllowOffline: true,
		Func:    RuleSQLE00099,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00099): "For MySQL DQL, SELECT FOR UPDATE is prohibited.".
You should follow the following logic:
1. For "select..." Statement, checks FOR the presence of an FOR UPDATE clause in the statement, If it does, report a violation.
2. For "insert... "Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For "union..." Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00099(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt:
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// "select..."
			if selectStmt.LockTp == ast.SelectLockForUpdate {
				//"select..." with "FOR UPDATE"
				rulepkg.AddResult(input.Res, input.Rule, SQLE00099)
				return nil
			}
		}

	}
	return nil
}

// ==== Rule code end ====
