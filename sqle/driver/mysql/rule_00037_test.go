package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00037(t *testing.T) {
	ruleName := ai.SQLE00037
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 1: 创建一个没有二级索引的表",
		"CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 2: 创建一个包含3个二级索引的表",
		"CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 3: 创建一个包含5个二级索引的表",
		"CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, address VARCHAR(255), city VARCHAR(100), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address), INDEX idx_city (city));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 4: 创建一个包含6个二级索引的表，超过限制",
		"CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, address VARCHAR(255), city VARCHAR(100), country VARCHAR(100), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address), INDEX idx_city (city), INDEX idx_country (country));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 5: 向已有4个二级索引的表中添加1个新二级索引，不超过限制",
		"ALTER TABLE user_data ADD INDEX idx_phone (phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 6: 向已有4个二级索引的表中添加2个新二级索引，导致总数超过限制",
		"ALTER TABLE user_data ADD INDEX idx_phone (phone), ADD INDEX idx_country (country);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), country VARCHAR(100), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 7: 向已有5个二级索引的表中添加1个新二级索引，导致总数超过限制",
		"ALTER TABLE user_data ADD INDEX idx_phone (phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address), INDEX idx_city (city));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 8: 向已有5个二级索引的表中不添加任何新索引",
		"ALTER TABLE user_data MODIFY COLUMN name VARCHAR(150);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address), INDEX idx_city (city));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 9: 向已有3个二级索引的表中添加2个新二级索引，总数达到5",
		"ALTER TABLE user_data ADD INDEX idx_phone (phone), ADD INDEX idx_country (country);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), country VARCHAR(100), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 10: 在已有4个二级索引的表中创建一个新索引，总数达到5",
		"CREATE INDEX idx_phone ON user_data (phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 11: 在已有5个二级索引的表中创建一个新索引，导致总数超过限制",
		"CREATE INDEX idx_phone ON user_data (phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age), INDEX idx_address (address), INDEX idx_city (city));"),
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 12: 在已有3个二级索引的表中创建一个新索引，总数增加到4",
		"CREATE INDEX idx_phone ON user_data (phone);",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100), age INT, phone VARCHAR(20), INDEX idx_name (name), INDEX idx_email (email), INDEX idx_age (age));"),
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 13: 创建一个没有二级索引的表(从xml中补充)",
		"CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(32) NOT NULL, sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 14: 创建一个包含5个二级索引的表(从xml中补充)",
		"CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(32) NOT NULL, sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100), INDEX idx_name_customers (name), INDEX idx_age_customers (age), INDEX idx_sex_customers (sex), INDEX idx_mark1_customers (mark1), INDEX idx_mark2_customers (mark2));",
		nil,
		nil,
		newTestResult())

	runAIRuleCase(rule, t, "case 15: 创建一个包含6个二级索引的表，超过限制(从xml中补充)",
		"CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(32) NOT NULL, sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100), INDEX idx_name_customers (name), INDEX idx_age_customers (age), INDEX idx_sex_customers (sex), INDEX idx_mark1_customers (mark1), INDEX idx_mark2_customers (mark2), INDEX idx_mark3_customers (mark3));",
		nil,
		nil,
		newTestResult().addResult(ruleName))

	runAIRuleCase(rule, t, "case 16: 向已有5个二级索引的表中添加1个新二级索引，导致总数超过限制(从xml中补充)",
		"ALTER TABLE customers ADD INDEX idx_mark5_customers (mark5);",
		session.NewAIMockContext().WithSQL("CREATE TABLE customers (id INT PRIMARY KEY, name VARCHAR(32) NOT NULL, sex INT, age INT, mark1 VARCHAR(20), mark2 VARCHAR(30), mark3 VARCHAR(40), mark4 VARCHAR(50), mark5 VARCHAR(100), INDEX idx_name_customers (name), INDEX idx_age_customers (age), INDEX idx_sex_customers (sex), INDEX idx_mark1_customers (mark1), INDEX idx_mark2_customers (mark2));"),
		nil,
		newTestResult().addResult(ruleName))
}

// ==== Rule test code end ====
