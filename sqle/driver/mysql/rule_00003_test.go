package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00003(t *testing.T) {
	ruleName := ai.SQLE00003
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with indexed field emp_id having NOT NULL and DEFAULT",
		"CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, name VARCHAR(50), department_id INT, INDEX (emp_id));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with indexed field emp_id missing NOT NULL",
		"CREATE TABLE employees (emp_id INT DEFAULT 0, name VARCHAR(50), department_id INT, INDEX (emp_id));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE with indexed field emp_id missing DEFAULT",
		"CREATE TABLE employees (emp_id INT NOT NULL, name VARCHAR(50), department_id INT, INDEX (emp_id));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: CREATE TABLE with composite index on emp_id and department_id, both having NOT NULL and DEFAULT",
		"CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT NOT NULL DEFAULT 0, name VARCHAR(50), INDEX (emp_id, department_id));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE with composite index on emp_id and department_id, department_id missing NOT NULL",
		"CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT DEFAULT 0, name VARCHAR(50), INDEX (emp_id, department_id));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: CREATE INDEX on emp_id where emp_id has NOT NULL and DEFAULT",
		"CREATE INDEX idx_emp_id ON employees (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, name VARCHAR(50), department_id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE INDEX on emp_id where emp_id is missing NOT NULL",
		"CREATE INDEX idx_emp_id ON employees (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT DEFAULT 0, name VARCHAR(50), department_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: CREATE INDEX on emp_id where emp_id has NOT NULL but missing DEFAULT",
		"CREATE INDEX idx_emp_id ON employees (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL, name VARCHAR(50), department_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 9: CREATE INDEX on composite fields emp_id and department_id, both having NOT NULL and DEFAULT",
		"CREATE INDEX idx_emp_dept ON employees (emp_id, department_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT NOT NULL DEFAULT 0, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 10: CREATE INDEX on composite fields emp_id and department_id, department_id missing NOT NULL",
		"CREATE INDEX idx_emp_dept ON employees (emp_id, department_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT DEFAULT 0, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 11: ALTER TABLE ADD INDEX on emp_id where emp_id has NOT NULL and DEFAULT",
		"ALTER TABLE employees ADD INDEX idx_emp_id (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, name VARCHAR(50), department_id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 12: ALTER TABLE ADD INDEX on emp_id where emp_id is missing NOT NULL",
		"ALTER TABLE employees ADD INDEX idx_emp_id (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT DEFAULT 0, name VARCHAR(50), department_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 13: ALTER TABLE ADD INDEX on emp_id where emp_id has NOT NULL but DEFAULT NULL",
		"ALTER TABLE employees ADD INDEX idx_emp_id (emp_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT NULL, name VARCHAR(50), department_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: ALTER TABLE ADD composite INDEX on emp_id and department_id, both having NOT NULL and DEFAULT",
		"ALTER TABLE employees ADD INDEX idx_emp_dept (emp_id, department_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT NOT NULL DEFAULT 0, name VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 15: ALTER TABLE ADD composite INDEX on emp_id and department_id, department_id missing NOT NULL",
		"ALTER TABLE employees ADD INDEX idx_emp_dept (emp_id, department_id);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (emp_id INT NOT NULL DEFAULT 0, department_id INT DEFAULT 0, name VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: CREATE TABLE with indexed field age having NOT NULL and DEFAULT (从xml中补充)",
		"CREATE TABLE person (id INT PRIMARY KEY, name VARCHAR(50), age INT NOT NULL DEFAULT 0, INDEX idx_age (age));",
		nil, nil, newTestResult())

	runAIRuleCase(rule, t, "case 17: CREATE TABLE with indexed field age missing NOT NULL (从xml中补充)",
		"CREATE TABLE person (id INT PRIMARY KEY, name VARCHAR(50), age INT, INDEX idx_age (age));",
		nil, nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 18: ALTER TABLE ADD UNIQUE INDEX on name and age where age has NOT NULL and DEFAULT (从xml中补充)",
		"ALTER TABLE person ADD UNIQUE INDEX idx_name_age (name, age);",
		session.NewAIMockContext().WithSQL("CREATE TABLE person (id INT PRIMARY KEY, name VARCHAR(50), age INT NOT NULL DEFAULT 0);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 19: ALTER TABLE ADD UNIQUE INDEX on name and age where age is missing NOT NULL (从xml中补充)",
		"ALTER TABLE person ADD UNIQUE INDEX idx_name_age (name, age);",
		session.NewAIMockContext().WithSQL("CREATE TABLE person (id INT PRIMARY KEY, name VARCHAR(50), age INT);"),
		nil, newTestResult().addResult(ruleName))

}

// ==== Rule test code end ====
