package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00018(t *testing.T) {
	ruleName := ai.SQLE00018
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with CHAR column length exactly 20",
		"CREATE TABLE employees (name CHAR(20), id INT);",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with CHAR column length greater than 20",
		"CREATE TABLE departments (dept_name CHAR(25), dept_id INT);",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 3: CREATE TABLE with multiple columns, no CHAR types",
		"CREATE TABLE projects (project_id INT, project_budget DECIMAL(10,2));",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 4: CREATE TABLE with multiple CHAR columns all lengths <= 20",
		"CREATE TABLE locations (city CHAR(15), country CHAR(20), code CHAR(10));",
		nil, /*mock context*/
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE with multiple CHAR columns, one length > 20",
		"CREATE TABLE clients (client_name CHAR(22), client_id INT, address VARCHAR(50));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 6: ALTER TABLE to add CHAR column with length 20",
		"ALTER TABLE employees ADD COLUMN middle_name CHAR(20);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (name CHAR(20), id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 7: ALTER TABLE to add CHAR column with length greater than 20",
		"ALTER TABLE departments ADD COLUMN description CHAR(30);",
		session.NewAIMockContext().WithSQL("CREATE TABLE departments (dept_name CHAR(25), dept_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: ALTER TABLE to modify existing column to CHAR with length 20",
		"ALTER TABLE projects MODIFY COLUMN city CHAR(20);",
		session.NewAIMockContext().WithSQL("CREATE TABLE projects (project_id INT, project_budget DECIMAL(10,2), city VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 9: ALTER TABLE to modify existing column to CHAR with length greater than 20",
		"ALTER TABLE clients MODIFY COLUMN client_name CHAR(25);",
		session.NewAIMockContext().WithSQL("CREATE TABLE clients (client_name CHAR(22), client_id INT, address VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 10: ALTER TABLE to add multiple CHAR columns all lengths <= 20",
		"ALTER TABLE employees ADD COLUMN nickname CHAR(15), ADD COLUMN suffix CHAR(10);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (name CHAR(20), id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 11: ALTER TABLE to add multiple columns with one CHAR column length > 20",
		"ALTER TABLE departments ADD COLUMN short_name CHAR(10), ADD COLUMN full_description CHAR(25);",
		session.NewAIMockContext().WithSQL("CREATE TABLE departments (dept_name CHAR(25), dept_id INT);"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: ALTER TABLE to modify multiple columns to CHAR with lengths <= 20",
		"ALTER TABLE projects MODIFY COLUMN city CHAR(18), MODIFY COLUMN state CHAR(20);",
		session.NewAIMockContext().WithSQL("CREATE TABLE projects (project_id INT, project_budget DECIMAL(10,2), city VARCHAR(50), state VARCHAR(50));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 13: ALTER TABLE to modify multiple columns with one CHAR column length > 20",
		"ALTER TABLE clients MODIFY COLUMN client_name CHAR(22), MODIFY COLUMN contact CHAR(18);",
		session.NewAIMockContext().WithSQL("CREATE TABLE clients (client_name CHAR(22), client_id INT, address VARCHAR(50), contact VARCHAR(50));"),
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 14: ALTER TABLE to add non-CHAR columns",
		"ALTER TABLE employees ADD COLUMN hire_date DATE, ADD COLUMN salary DECIMAL(10,2);",
		session.NewAIMockContext().WithSQL("CREATE TABLE employees (name CHAR(20), id INT);"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 15: ALTER TABLE to modify non-CHAR columns",
		"ALTER TABLE departments MODIFY COLUMN dept_id BIGINT, MODIFY COLUMN budget DECIMAL(15,2);",
		session.NewAIMockContext().WithSQL("CREATE TABLE departments (dept_name CHAR(25), dept_id INT, budget DECIMAL(10,2));"),
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 16: CREATE TABLE with CHAR column length greater than 20 (从xml中补充)",
		"CREATE TABLE test_char (id INT PRIMARY KEY, char_col CHAR(30), varchar_col VARCHAR(30));",
		nil, /*mock context*/
		nil, newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 17: CREATE TABLE with VARCHAR columns (从xml中补充)",
		"CREATE TABLE employee_improved (id INT PRIMARY KEY, name VARCHAR(50), department VARCHAR(30));",
		nil, /*mock context*/
		nil, newTestResult())
}

// ==== Rule test code end ====
