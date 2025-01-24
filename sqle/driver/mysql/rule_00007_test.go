package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00007(t *testing.T) {
	ruleName := ai.SQLE00007
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: CREATE TABLE with no auto_increment columns",
		"CREATE TABLE employees (id INT, name VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE with one auto_increment column",
		"CREATE TABLE employees (id INT AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY(id));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: CREATE TABLE with two auto_increment columns",
		"CREATE TABLE employees (id INT AUTO_INCREMENT, user_id INT AUTO_INCREMENT, name VARCHAR(100), PRIMARY KEY(id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 4: CREATE TABLE with one auto_increment column not set as PRIMARY KEY",
		"CREATE TABLE employees (id INT AUTO_INCREMENT, name VARCHAR(100), email VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 5: CREATE TABLE with multiple columns but no auto_increment (从xml中补充)",
		"CREATE TABLE departments (dept_id INT, dept_name VARCHAR(100), location VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 6: CREATE TABLE with one auto_increment column in a different table (从xml中补充)",
		"CREATE TABLE projects (project_id INT AUTO_INCREMENT, project_name VARCHAR(100), PRIMARY KEY(project_id));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 7: CREATE TABLE with two auto_increment columns in a different table (从xml中补充)",
		"CREATE TABLE tasks (task_id INT AUTO_INCREMENT, subtask_id INT AUTO_INCREMENT, task_name VARCHAR(100), PRIMARY KEY(task_id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: CREATE TABLE with one auto_increment column and a composite primary key (从xml中补充)",
		"CREATE TABLE assignments (assignment_id INT AUTO_INCREMENT, employee_id INT, project_id INT, PRIMARY KEY(assignment_id, employee_id));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 9: CREATE TABLE with two auto_increment columns and a composite primary key (从xml中补充)",
		"CREATE TABLE schedules (schedule_id INT AUTO_INCREMENT, event_id INT AUTO_INCREMENT, event_name VARCHAR(100), PRIMARY KEY(schedule_id, event_id));",
		nil,
		nil,
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
