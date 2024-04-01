package mysql

import (
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule/ai"
)

// ==== Rule test code start ====
func TestRuleSQLE00015(t *testing.T) {
	ruleName := ai.SQLE00015
	rule := rulepkg.RuleHandlerMap[ruleName].Rule

	//create table, without collate option
	runSingleRuleInspectCase(rule, t, "create table, without collate option", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	);
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//create table, with expected collate option
	runSingleRuleInspectCase(rule, t, "create table, with expected collate option", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	) COLLATE=utf8mb4_0900_ai_ci;
	`, newTestResult())

	//create table, with unexpected collate option
	runSingleRuleInspectCase(rule, t, "create table, with unexpected collate option", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test"
	) COLLATE=latin1_swedish_ci;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//create table, with expected column collate option
	runSingleRuleInspectCase(rule, t, "create table, with expected column collate option", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test" COLLATE utf8mb4_0900_ai_ci
	) ;
	`, newTestResult())

	//create table, with unexpected column collate option
	runSingleRuleInspectCase(rule, t, "create table, with unexpected column collate option", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test" COLLATE latin1_swedish_ci
	) ;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//alter table, without collate option
	runSingleRuleInspectCase(rule, t, "alter table, without collate option", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8mb4';
	`, newTestResult())

	//alter table, with expected collate option
	runSingleRuleInspectCase(rule, t, "alter table, with expected collate option", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8mb4' COLLATE utf8mb4_0900_ai_ci;
	`, newTestResult())

	//alter table, with unexpected collate option
	runSingleRuleInspectCase(rule, t, "alter table, with unexpected collate option", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 CHARACTER SET 'utf8mb4' COLLATE latin1_swedish_ci;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//alter table, with expected column collate option
	runSingleRuleInspectCase(rule, t, "alter table, with expected column collate option", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN c1 VARCHAR(255) COLLATE utf8mb4_0900_ai_ci;
	`, newTestResult())

	//alter table, with unexpected column collate option
	runSingleRuleInspectCase(rule, t, "alter table, with unexpected column collate option", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN c1 VARCHAR(255) COLLATE latin1_swedish_ci;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//create database, without collate option
	runSingleRuleInspectCase(rule, t, "create database, without collate option", DefaultMysqlInspect(),
		`
	CREATE DATABASE  if not exists exist_db_1;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//create database, with expected collate option
	runSingleRuleInspectCase(rule, t, "create database, with expected collate option", DefaultMysqlInspect(),
		`
	CREATE DATABASE  if not exists exist_db_1 COLLATE utf8mb4_0900_ai_ci;
	`, newTestResult())

	//create database, with expected collate option upper case
	runSingleRuleInspectCase(rule, t, "create database, with expected collate option upper case", DefaultMysqlInspect(),
		`
	CREATE DATABASE  if not exists exist_db_1 COLLATE UTF8MB4_0900_AI_CI;
	`, newTestResult())

	//create database, with unexpected collate option
	runSingleRuleInspectCase(rule, t, "create database, with unexpected collate option", DefaultMysqlInspect(),
		`
	CREATE DATABASE  if not exists exist_db_1 COLLATE latin1_swedish_ci;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))

	//alter database, without collate option
	runSingleRuleInspectCase(rule, t, "alter database, without collate option", DefaultMysqlInspect(),
		`
	ALTER DATABASE exist_db CHARACTER SET utf8mb4;
	`, newTestResult())

	//alter database, with expected collate option
	runSingleRuleInspectCase(rule, t, "alter database, with expected collate option", DefaultMysqlInspect(),
		`
	ALTER DATABASE exist_db COLLATE utf8mb4_0900_ai_ci;
	`, newTestResult())

	//alter database, with unexpected collate option
	runSingleRuleInspectCase(rule, t, "alter database, with unexpected collate option", DefaultMysqlInspect(),
		`
	ALTER DATABASE exist_db COLLATE latin1_swedish_ci;
	`, newTestResult().addResult(ruleName, "utf8mb4_0900_ai_ci"))
}

// ==== Rule test code end ====
