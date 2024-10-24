package mysql

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

// ==== Rule test code start ====
func TestRuleSQLE00057(t *testing.T) {
	ruleName := ai.SQLE00057
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	runAIRuleCase(rule, t, "case 0: CREATE TABLE 未指定 ENGINE，默认存储引擎为 MyISAM",
		"CREATE TABLE db1.user_data (id INT PRIMARY KEY, name VARCHAR(100));",
		session.NewAIMockContext().WithSQL("create database db1;use db1;"),
		[]*AIMockSQLExpectation{
			{
				Query: "select @@default_storage_engine",
				Rows:  sqlmock.NewRows([]string{"@@default_storage_engine"}).AddRow("MyISAM"),
			},
		}, newTestResult().addResult(ruleName))

	// exist_db 是  InnoDB
	runAIRuleCase(rule, t, "case 1: CREATE TABLE 未指定 ENGINE，默认存储引擎为 ENGINE",
		"CREATE TABLE exist_db.user_data (id INT PRIMARY KEY, name VARCHAR(100));",
		nil,
		nil, newTestResult())

	runAIRuleCase(rule, t, "case 2: CREATE TABLE 使用 MyISAM 引擎",
		"CREATE TABLE archive_data (id INT PRIMARY KEY, archive_date DATE) ENGINE=MyISAM;",
		nil,
		nil,
		newTestResult().addResult(ruleName),
	)

	runAIRuleCase(rule, t, "case 5: ALTER TABLE 将存储引擎修改为 InnoDB",
		"ALTER TABLE user_data ENGINE=InnoDB;",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil,
		newTestResult(),
	)

	runAIRuleCase(rule, t, "case 6: ALTER TABLE 将存储引擎修改为 MyISAM",
		"ALTER TABLE user_data ENGINE=MyISAM;",
		session.NewAIMockContext().WithSQL("CREATE TABLE user_data (id INT PRIMARY KEY, name VARCHAR(100));"),
		nil,
		newTestResult().addResult(ruleName),
	)

}

// ==== Rule test code end ====
