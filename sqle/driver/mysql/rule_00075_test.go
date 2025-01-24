package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00075(t *testing.T) {
	ruleName := ai.SQLE00075
	rule := rulepkg.AIRuleHandlerMap[ruleName].Rule

	//create table, no charset, no collate
	runSingleRuleInspectCase(rule, t, "create table, no charset, no collate", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10),
	PRIMARY KEY (id)
	);
	`, newTestResult())

	//create table, with charset, no collate
	runSingleRuleInspectCase(rule, t, "create table, with charset, no collate", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10) CHARSET utf8mb4,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//create table, with charset, with collate
	runSingleRuleInspectCase(rule, t, "create table, with charset, with collate", DefaultMysqlInspect(), `
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned DEFAULT 100 AUTO_INCREMENT,
	a varchar(10) CHARSET utf8mb4 COLLATE utf8_general_ci,
	PRIMARY KEY (id)
	);
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add column, no charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table add column, no charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table add column, with charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table add column, with charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table add column, with charset, with collate
	runSingleRuleInspectCase(rule, t, "alter table add column, with charset, with collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 ADD COLUMN a varchar(10) CHARSET utf8mb4 COLLATE utf8_general_ci COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table modify column, no charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table modify column, no charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table modify column, with charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table modify column, with charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table modify column, with charset, with collate
	runSingleRuleInspectCase(rule, t, "alter table modify column, with charset, with collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 MODIFY v1 varchar(10) CHARSET utf8mb4 COLLATE utf8_general_ci COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "v1"))

	//alter table change column, no charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table change column, no charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) COMMENT "unit test";
	`, newTestResult())

	//alter table change column, with charset, no collate
	runSingleRuleInspectCase(rule, t, "alter table change column, with charset, no collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) CHARSET utf8mb4 COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))

	//alter table change column, with charset, with collate
	runSingleRuleInspectCase(rule, t, "alter table change column, with charset, with collate", DefaultMysqlInspect(), `
	ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 a varchar(10) CHARSET utf8mb4 COLLATE utf8_general_ci COMMENT "unit test";
	`, newTestResult().addResult(ruleName, "a"))
}

// ==== Rule test code end ====
