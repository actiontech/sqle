package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00174(t *testing.T) {
	ruleName := ai.SQLE00174
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: GRANT ALL 权限", "GRANT ALL ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 2: GRANT SUPER 权限", "GRANT SUPER ON *.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: GRANT WITH GRANT OPTION 权限", "GRANT SELECT ON database.* TO 'user'@'localhost' WITH GRANT OPTION;",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: GRANT SELECT 权限", "GRANT SELECT ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: GRANT INSERT 权限", "GRANT INSERT ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 6: GRANT UPDATE 权限", "GRANT UPDATE ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: GRANT DELETE 权限", "GRANT DELETE ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 8: GRANT CREATE 权限", "GRANT CREATE ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: GRANT ALTER 权限", "GRANT ALTER ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: GRANT DROP 权限", "GRANT DROP ON database.* TO 'user'@'localhost';",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: GRANT ALL 权限给 user1", "GRANT ALL ON *.* TO 'user1'@'localhost';",
		session.NewAIMockContext().WithSQL("CREATE USER 'user1'@'localhost' IDENTIFIED BY 'Root123@aB';"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: GRANT SELECT,UPDATE,DELETE,INSERT,ALTER,CREATE,DROP WITH GRANT OPTION 给 user1",
		"GRANT SELECT, UPDATE, DELETE, INSERT, ALTER, CREATE, DROP ON db_mysql.* TO 'user1'@'localhost' WITH GRANT OPTION;",
		session.NewAIMockContext().WithSQL("CREATE USER 'user1'@'localhost' IDENTIFIED BY 'Root123@aB';"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: GRANT SUPER WITH GRANT OPTION 给 user2",
		"GRANT SUPER ON *.* TO 'user2'@'localhost' WITH GRANT OPTION;",
		session.NewAIMockContext().WithSQL("CREATE USER 'user2'@'localhost' IDENTIFIED BY 'Root123@aB';"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: GRANT UPDATE,DELETE,INSERT 权限给 user1",
		"GRANT UPDATE, DELETE, INSERT ON db_mysql.* TO 'user1'@'localhost';",
		session.NewAIMockContext().WithSQL("CREATE USER 'user1'@'localhost' IDENTIFIED BY 'Root123@aB';"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 15: GRANT SELECT 权限给 user2",
		"GRANT SELECT ON db_mysql.* TO 'user2'@'localhost';",
		session.NewAIMockContext().WithSQL("CREATE USER 'user2'@'localhost' IDENTIFIED BY 'Root123@aB';"),
		nil, newTestResult())
}

// ==== Rule test code end ====
