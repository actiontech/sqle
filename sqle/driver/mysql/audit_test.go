package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/log"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testResult struct {
	Results *driver.AuditResult
	rules   map[string]rulepkg.RuleHandler
}

func newTestResult() *testResult {
	return &testResult{
		Results: driver.NewInspectResults(),
		rules:   rulepkg.RuleHandlerMap,
	}
}

func (t *testResult) add(level driver.RuleLevel, message string, args ...interface{}) *testResult {
	t.Results.Add(level, message, args...)
	return t
}

func (t *testResult) addResult(ruleName string, args ...interface{}) *testResult {
	handler, ok := t.rules[ruleName]
	if !ok {
		return t
	}
	level := handler.Rule.Level
	message := handler.Message

	return t.add(level, message, args...)
}

func (t *testResult) level() driver.RuleLevel {
	return t.Results.Level()
}

func (t *testResult) message() string {
	return t.Results.Message()
}

func DefaultMysqlInspect() *Inspect {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &Inspect{
		log: log.NewEntry(),
		inst: &driver.DSN{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Password:     "123456",
			DatabaseName: "mysql",
		},
		Ctx: session.NewMockContext(nil),
		cnf: &Config{
			DDLOSCMinSize:      16,
			DDLGhostMinSize:    16,
			DMLRollbackMaxRows: 1000,
		},
	}
}

func NewMockInspect(e *executor.Executor) *Inspect {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &Inspect{
		log: log.NewEntry(),
		inst: &driver.DSN{
			Host:         "127.0.0.1",
			Port:         "3306",
			User:         "root",
			Password:     "123456",
			DatabaseName: "mysql",
		},
		Ctx: session.NewMockContext(e),
		cnf: &Config{
			DDLOSCMinSize:      16,
			DDLGhostMinSize:    16,
			DMLRollbackMaxRows: 1000,
		},
		dbConn: e,
	}
}

func runSingleRuleInspectCase(rule driver.Rule, t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	i.rules = []*driver.Rule{&rule}
	inspectCase(t, desc, i, sql, results...)
}

func runDefaultRulesInspectCase(t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	ptrRules := make([]*driver.Rule, len(rulepkg.DefaultTemplateRules))
	for i := range rulepkg.DefaultTemplateRules {
		// remove DDL_CHECK_OBJECT_NAME_USING_CN in default rules for init test.
		if rulepkg.DefaultTemplateRules[i].Name == rulepkg.DDLCheckOBjectNameUseCN {
			continue
		}

		ptrRules[i] = &rulepkg.DefaultTemplateRules[i]
	}

	i.rules = ptrRules
	inspectCase(t, desc, i, sql, results...)
}

func inspectCase(t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	stmts, err := util.ParseSql(sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}

	if len(stmts) != len(results) {
		t.Errorf("%s test failled, error: result is unknow\n", desc)
		return
	}

	for idx, stmt := range stmts {
		result, err := i.Audit(context.TODO(), stmt.Text())
		if err != nil {
			t.Error(err)
			return
		}
		if result.Level() != results[idx].level() || result.Message() != results[idx].message() {
			t.Errorf("%s test failled, \n\nsql:\n %s\n\nexpect level: %s\nexpect result:\n%s\n\nactual level: %s\nactual result:\n%s\n",
				desc, stmt.Text(), results[idx].level(), results[idx].message(), result.Level(), result.Message())
		} else {
			t.Log(fmt.Sprintf("\n\ncase:%s\nactual level: %s\nactual result:\n%s\n\n", desc, result.Level(), result.Message()))
		}
	}
}

func TestMessage(t *testing.T) {
	runDefaultRulesInspectCase(t, "check inspect message", DefaultMysqlInspect(),
		"use no_exist_db", newTestResult().add(driver.RuleLevelError, "schema no_exist_db 不存在"))
}

func TestCheckInvalidUse(t *testing.T) {
	runDefaultRulesInspectCase(t, "use_database: database not exist", DefaultMysqlInspect(),
		"use no_exist_db",
		newTestResult().add(driver.RuleLevelError,
			SchemaNotExistMessage, "no_exist_db"),
	)

	inspect1 := DefaultMysqlInspect()
	inspect1.Ctx.AddSysVar(session.SysVarLowerCaseTableNames, "1")
	runDefaultRulesInspectCase(t, "", inspect1,
		"use EXIST_DB",
		newTestResult(),
	)
}

func TestCaseSensitive(t *testing.T) {
	runDefaultRulesInspectCase(t, "", DefaultMysqlInspect(),
		`
select id from exist_db.EXIST_TB_1 where id = 1;
`,
		newTestResult().add(driver.RuleLevelError,
			TableNotExistMessage, "exist_db.EXIST_TB_1"))

	inspect1 := DefaultMysqlInspect()
	inspect1.Ctx.AddSysVar(session.SysVarLowerCaseTableNames, "1")
	runDefaultRulesInspectCase(t, "", inspect1,
		`
select id from exist_db.EXIST_TB_1 where id = 1;
`,
		newTestResult())
}

func TestCheckInvalidCreateTable(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: schema not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists not_exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError,
			SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "create_table: table is exist(1)", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)
	handler := rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKWithoutIfNotExists]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckPKWithoutIfNotExists)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKWithoutIfNotExists] = handler
	}()
	runDefaultRulesInspectCase(t, "create_table: table is exist(2)", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError,
			TableExistMessage, "exist_db.exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "create_table: refer table not exist", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_1 like exist_db.not_exist_tb_2;
`,
		newTestResult().add(driver.RuleLevelError,
			TableNotExistMessage, "exist_db.not_exist_tb_2"),
	)

	runDefaultRulesInspectCase(t, "create_table: multi pk(1)", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, MultiPrimaryKeyMessage))

	runDefaultRulesInspectCase(t, "create_table: multi pk(2)", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
PRIMARY KEY (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, MultiPrimaryKeyMessage))

	runDefaultRulesInspectCase(t, "create_table: duplicate column", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateColumnsMessage,
			"v1"))

	runDefaultRulesInspectCase(t, "create_table: duplicate index", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1),
INDEX idx_1 (v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexesMessage,
			"idx_1"))

	runDefaultRulesInspectCase(t, "create_table: key column not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v3),
INDEX idx_2 (v4,v5)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			"v3,v4,v5"))

	runDefaultRulesInspectCase(t, "create_table: pk column not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id11)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			"id11"))

	runDefaultRulesInspectCase(t, "create_table: pk column is duplicate", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id,id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
			"id"))

	runDefaultRulesInspectCase(t, "create_table: index column is duplicate", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_1",
			"v1"))

	runDefaultRulesInspectCase(t, "create_table: index column is duplicate(2)", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX (v1,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "(匿名)",
			"v1").addResult(rulepkg.DDLCheckIndexPrefix, "idx_"))

	runDefaultRulesInspectCase(t, "create_table: index column is duplicate(3)", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,v1),
INDEX idx_2 (v1,v2,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_1",
			"v1").add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_2", "v2"))
}

func TestCheckInvalidAlterTable(t *testing.T) {
	// It's trick :),
	// elegant method: unit test support MySQL.
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckTableDBEngine)
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckTableCharacterSet)
	runDefaultRulesInspectCase(t, "alter_table: schema not exist", DefaultMysqlInspect(),
		`ALTER TABLE not_exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage,
			"not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "alter_table: table not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.not_exist_tb_1 add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage,
			"exist_db.not_exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add a exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add column v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(driver.RuleLevelError, ColumnExistMessage,
			"v1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: drop a not exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 drop column v5;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage,
			"v5"),
	)

	runDefaultRulesInspectCase(t, "alter_table: alter a not exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 alter column v5 set default 'v5';
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage,
			"v5"),
	)

	runDefaultRulesInspectCase(t, "alter_table: change a exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 change column v1 v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "alter_table: change a not exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 change column v5 v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage,
			"v5"),
	)

	runDefaultRulesInspectCase(t, "alter_table: change column to a exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 change column v2 v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(driver.RuleLevelError, ColumnExistMessage,
			"v1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add pk ok", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 Add primary key(id);
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add pk but exist pk", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add primary key(v1);
`,
		newTestResult().add(driver.RuleLevelError, PrimaryKeyExistMessage).addResult(rulepkg.DDLCheckPKWithoutAutoIncrement).addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add pk but key column not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 Add primary key(id11);
`,
		newTestResult().add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			"id11"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add pk but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 Add primary key(id,id);
`,
		newTestResult().add(driver.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
			"id"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add a exist index", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add index idx_1 (v1);
`,
		newTestResult().add(driver.RuleLevelError, IndexExistMessage,
			"idx_1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: drop a not exist index", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 drop index idx_2;
`,
		newTestResult().add(driver.RuleLevelError, IndexNotExistMessage,
			"idx_2"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index but key column not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add index idx_2 (v3);
`,
		newTestResult().add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			"v3"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add index idx_2 (id,id);
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_2",
			"id"),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add index (id,id);
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "(匿名)",
			"id").addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)
}

func TestCheckInvalidCreateDatabase(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_database: schema exist(1)", DefaultMysqlInspect(),
		`
CREATE DATABASE if not exists exist_db;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_database: schema exist(2)", DefaultMysqlInspect(),
		`
CREATE DATABASE exist_db;
`,
		newTestResult().add(driver.RuleLevelError, SchemaExistMessage, "exist_db"),
	)
}

func TestCheckInvalidCreateIndex(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_index: schema not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON not_exist_db.not_exist_tb(v1);
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "create_index: table not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON exist_db.not_exist_tb(v1);
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "create_index: index exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().add(driver.RuleLevelError, IndexExistMessage, "idx_1"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(v3);
`,
		newTestResult().add(driver.RuleLevelError, KeyedColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column is duplicate", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(id,id);
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_2", "id"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column is duplicate", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(id,id,v1);
`,
		newTestResult().add(driver.RuleLevelError, DuplicateIndexedColumnMessage, "idx_2", "id"),
	)
}

func TestCheckInvalidDrop(t *testing.T) {
	handler := rulepkg.RuleHandlerMap[rulepkg.DDLDisableDropStatement]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLDisableDropStatement)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLDisableDropStatement] = handler
	}()
	runDefaultRulesInspectCase(t, "drop_database: ok", DefaultMysqlInspect(),
		`
DROP DATABASE if exists exist_db;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_database: schema not exist(1)", DefaultMysqlInspect(),
		`
DROP DATABASE if exists not_exist_db;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_database: schema not exist(2)", DefaultMysqlInspect(),
		`
DROP DATABASE not_exist_db;
`,
		newTestResult().add(driver.RuleLevelError,
			SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "drop_table: ok", DefaultMysqlInspect(),
		`
DROP TABLE exist_db.exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_table: schema not exist(1)", DefaultMysqlInspect(),
		`
DROP TABLE if exists not_exist_db.not_exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_table: schema not exist(2)", DefaultMysqlInspect(),
		`
DROP TABLE not_exist_db.not_exist_tb_1;
`,
		newTestResult().add(driver.RuleLevelError,
			SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "drop_table: table not exist", DefaultMysqlInspect(),
		`
DROP TABLE exist_db.not_exist_tb_1;
`,
		newTestResult().add(driver.RuleLevelError,
			TableNotExistMessage, "exist_db.not_exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "drop_index: ok", DefaultMysqlInspect(),
		`
DROP INDEX idx_1 ON exist_db.exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_index: index not exist", DefaultMysqlInspect(),
		`
DROP INDEX idx_2 ON exist_db.exist_tb_1;
`,
		newTestResult().add(driver.RuleLevelError, IndexNotExistMessage, "idx_2"),
	)

	runDefaultRulesInspectCase(t, "drop_index: if exists and index not exist", DefaultMysqlInspect(),
		`
DROP INDEX IF EXISTS idx_2 ON exist_db.exist_tb_1;
`,
		newTestResult(),
	)
}

func TestCheckInvalidInsert(t *testing.T) {
	runDefaultRulesInspectCase(t, "insert: schema not exist", DefaultMysqlInspect(),
		`
insert into not_exist_db.not_exist_tb values (1,"1","1");
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "insert: table not exist", DefaultMysqlInspect(),
		`
insert into exist_db.not_exist_tb values (1,"1","1");
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "insert: column not exist(1)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v3) values (1,"1","1");
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "insert: column not exist(2)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 set id=1,v1="1",v3="1";
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "insert: column is duplicate(1)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v1) values (1,"1","1");
`,
		newTestResult().add(driver.RuleLevelError, DuplicateColumnsMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "insert: column is duplicate(2)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 set id=1,v1="1",v1="1";
`,
		newTestResult().add(driver.RuleLevelError, DuplicateColumnsMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "insert: do not match values and columns", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2","2");
`,
		newTestResult().add(driver.RuleLevelError, ColumnsValuesNotMatchMessage),
	)
}

func TestCheckInvalidUpdate(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: ok", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1="2" where id=1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "update: ok", DefaultMysqlInspect(),
		`
update exist_tb_1 set v1="2" where exist_db.exist_tb_1.id=1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "update: schema not exist", DefaultMysqlInspect(),
		`
update not_exist_db.not_exist_tb set v1="2" where id=1;
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "update: table not exist", DefaultMysqlInspect(),
		`
update exist_db.not_exist_tb set v1="2" where id=1;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v3="2" where id=1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "update: where column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1="2" where v3=1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "update with alias: ok", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set t.v1 = "1" where t.id = 1;
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "update with alias: table not exist", DefaultMysqlInspect(),
		`
update exist_db.not_exist_tb as t set t.v3 = "1" where t.id = 1;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set t.v3 = "1" where t.id = 1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "t.v3"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set t.v1 = "1" where t.v3 = 1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "t.v3"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set exist_tb_1.v1 = "1" where t.id = 1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_1.v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: ok", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v1 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "multi-update: ok", DefaultMysqlInspect(),
		`
update exist_tb_1 inner join exist_tb_2 on exist_tb_1.id = exist_tb_2.id set exist_tb_1.v1 = "1" where exist_tb_1.id = 1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "multi-update: table not exist", DefaultMysqlInspect(),
		`
update exist_db.not_exist_tb set exist_tb_1.v2 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v3 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_1.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_2.v3 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_2.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v1 = "1" where exist_tb_1.v3 = exist_tb_2.v3;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_1.v3,exist_tb_2.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1,exist_db.exist_tb_2 set exist_tb_3.v1 = "1" where exist_tb_1.v1 = exist_tb_2.v1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_3.v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1,exist_db.exist_tb_2 set not_exist_db.exist_tb_1.v1 = "1" where exist_tb_1.v1 = exist_tb_2.v1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "not_exist_db.exist_tb_1.v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set user_id = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set v1 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(driver.RuleLevelError, ColumnIsAmbiguousMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set v1 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(driver.RuleLevelError, ColumnIsAmbiguousMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: where column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v1 = "1" where v1 = 1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnIsAmbiguousMessage, "v1"),
	)
}

func TestCheckInvalidDelete(t *testing.T) {
	runDefaultRulesInspectCase(t, "delete: ok", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where id=1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "delete: schema not exist", DefaultMysqlInspect(),
		`
delete from not_exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "delete: table not exist", DefaultMysqlInspect(),
		`
delete from exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where v3=1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where exist_tb_1.v3=1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_1.v3"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where exist_tb_2.id=1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "exist_tb_2.id"),
	)
}

func TestCheckInvalidSelect(t *testing.T) {
	runDefaultRulesInspectCase(t, "select: schema not exist", DefaultMysqlInspect(),
		`
select id from not_exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "select: table not exist", DefaultMysqlInspect(),
		`
select id from exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb"),
	)
}

func TestCheckSelectAll(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLDisableSelectAllColumn].Rule, t, "select_from: all columns", DefaultMysqlInspect(),
		"select * from exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(rulepkg.DMLDisableSelectAllColumn),
	)
}

func TestCheckWhereInvalid(t *testing.T) {
	runDefaultRulesInspectCase(t, "select_from: has where condition", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id > 1;",
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(3)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(4)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "update: has where condition", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "update: no where condition(1)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1';",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(2)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where 1=1 and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(3)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(4)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: has where condition", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "delete: no where condition(1)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(2)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(3)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(4)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))
}

func TestCheckWhereInvalid_FP(t *testing.T) {
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id=?;",
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=?;",
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=? and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where ?=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "[pf]update: has where condition", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where id = ?;",
		newTestResult())

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(1)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(2)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1=? where 1=1 and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(3)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1=? where id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(4)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1=? where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]delete: no where condition(1)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=? and ?=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]delete: no where condition(2)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=? and id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))
}

func TestCheckCreateTableWithoutIfNotExists(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: need \"if not exists\"", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutIfNotExists),
	)
}

func TestCheckObjectNameUsingKeyword(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: using keyword", DefaultMysqlInspect(),
		"CREATE TABLE if not exists exist_db.`select` ("+
			"id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT \"unit test\","+
			"v1 varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"`create` varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"PRIMARY KEY (id),"+
			"INDEX `show` (v1)"+
			")ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT=\"unit test\";",
		newTestResult().addResult(rulepkg.DDLCheckObjectNameUsingKeyword, "select, create, show").
			addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

}

func TestAlterTableMerge(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge].Rule, t, "alter_table: alter table need merge", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 Add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
ALTER TABLE exist_db.exist_tb_1 Add column v6 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult(),
		newTestResult().addResult(rulepkg.DDLCheckAlterTableNeedMerge),
	)
}

func TestCheckObjectNameLength(t *testing.T) {
	length64 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcd"
	length65 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcde"

	runDefaultRulesInspectCase(t, "create_table: table length <= 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length64),
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
%s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "create_table: index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_%s (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:Add column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:change column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD index idx_%s (v1);`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:rename index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME index idx_1 TO idx_%s;`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)
}

func TestCheckPrimaryKey(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: primary key exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKNotExist),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL KEY DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(2)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(2)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned),
	)
}

func TestCheckColumnCharLength(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: check char(20)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 char(20) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: check char(21)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 char(21) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckColumnCharLength),
	)
}

func TestCheckIndexCount(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexCount].Rule
	runSingleRuleInspectCase(rule, t, "create_table: index <= 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (id),
INDEX idx_2 (id),
INDEX idx_3 (id),
INDEX idx_4 (id),
INDEX idx_5 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "create_table: index > 5", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (id),
INDEX idx_2 (id),
INDEX idx_3 (id),
INDEX idx_4 (id),
INDEX idx_5 (id),
INDEX idx_6 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexCount, 5),
	)
}

func TestCheckCompositeIndexMax(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckCompositeIndexMax].Rule
	runSingleRuleInspectCase(rule, t, "create_table: composite index columns <= 3", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v3 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v4 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (id,v1,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "create_table: composite index columns > 3", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v3 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v4 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (id,v1,v2,v3,v4,v5)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckCompositeIndexMax, 3),
	)
}

func TestCheckTableWithoutInnodbUtf8mb4(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: ok", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table engine not innodb", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableDBEngine, "InnoDB"),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset not utf8mb4", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1  COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)
}

func TestCheckIndexColumnWithBlob(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: disable index column blob (1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_b1 (b1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
	)

	runDefaultRulesInspectCase(t, "create_table: disable index column blob (2)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob UNIQUE KEY COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
	)

	handler := rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckAlterTableNeedMerge)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge] = handler
	}()

	runDefaultRulesInspectCase(t, "create_table: disable index column blob (3)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
CREATE INDEX idx_1 ON exist_db.not_exist_tb_1(b1);
ALTER TABLE exist_db.not_exist_tb_1 ADD INDEX idx_2(b1);
ALTER TABLE exist_db.not_exist_tb_1 ADD COLUMN b2 blob UNIQUE KEY COMMENT "unit test";
ALTER TABLE exist_db.not_exist_tb_1 MODIFY COLUMN b1 blob UNIQUE KEY COMMENT "unit test";
`,
		newTestResult(),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBolb),
	)
}

func TestDisableForeignKey(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: has foreign key", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
FOREIGN KEY (id) REFERENCES exist_tb_1(id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLDisableFK),
	)
}

func TestCheckTableComment(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableWithoutComment].Rule, t, "create_table: table without comment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(rulepkg.DDLCheckTableWithoutComment),
	)
}

func TestCheckColumnComment(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutComment].Rule
	runSingleRuleInspectCase(rule, t, "create_table: column without comment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutComment),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: column without comment(1)", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL DEFAULT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutComment),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: column without comment(2)", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v2 v3 varchar(255) NOT NULL DEFAULT "unit test" ;
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutComment),
	)
}

func TestCheckIndexPrefix(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: index prefix not idx_", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

	runDefaultRulesInspectCase(t, "alter_table: index prefix not idx_", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD INDEX index_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

	runDefaultRulesInspectCase(t, "create_index: index prefix not idx_", DefaultMysqlInspect(),
		`
CREATE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

	for _, sql := range []string{
		`create table exist_db.t1(id int, c1 varchar(10), index idx_c1(c1))`,
		`create table exist_db.t1(id int, c1 varchar(10), index IDX_C1(c1))`,
		`create index IDX_v1 ON exist_db.exist_tb_1(v1);`,
		`create index idx_V1 ON exist_db.exist_tb_1(v1);`,
		`alter table exist_db.exist_tb_1 Add index idx_v1(v1);`,
		`alter table exist_db.exist_tb_1 Add index IDX_V1(v1);`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexPrefix].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestCheckUniqueIndexPrefix(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndexPrefix].Rule
	runSingleRuleInspectCase(rule, t, "create_table: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
UNIQUE INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckUniqueIndexPrefix, "uniq_"),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX index_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckUniqueIndexPrefix, "uniq_"),
	)

	runSingleRuleInspectCase(rule, t, "create_index: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
CREATE UNIQUE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckUniqueIndexPrefix, "uniq_"),
	)

	for _, sql := range []string{
		`create table exist_db.t1(id int, c1 varchar(10), unique index uniq_c1(c1))`,
		`create table exist_db.t1(id int, c1 varchar(10), unique index UNIQ_C1(c1))`,
		`create unique index uniq_v1 ON exist_db.exist_tb_1(v1);`,
		`create unique index UNIQ_V1 ON exist_db.exist_tb_1(v1);`,
		`alter table exist_db.exist_tb_1 Add unique index uniq_v1(v1);`,
		`alter table exist_db.exist_tb_1 Add unique index UNIQ_V1(v1);`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndexPrefix].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestCheckColumnDefault(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: auto increment column without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "alter_table: blob column without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob COMMENT "unit test";
`,
		newTestResult(),
	)
}

func TestCheckColumnTimestampDefault(t *testing.T) {
	handler := rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutDefault]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckColumnWithoutDefault)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutDefault] = handler
	}()

	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 timestamp COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnTimestampWitoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 timestamp NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnTimestampWitoutDefault),
	)
}

func TestCheckColumnBlobNotNull(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 blob NOT NULL COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobWithNotNull),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobWithNotNull),
	)
}

func TestCheckColumnBlobDefaultNull(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 blob DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobDefaultIsNotNull),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobDefaultIsNotNull),
	)
}

func TestCheckDMLWithLimit(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)
}

func TestCheckDMLWithLimit_FP(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]update: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "[fp]delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)
}

func TestCheckDMLWithOrderBy(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with order by", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)
}

func TestCheckDMLWithOrderBy_FP(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]update: with order by", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)

	runDefaultRulesInspectCase(t, "[fp]delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1=? where id=? order by ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)
}

func TestCheckInsertColumnsExist(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "insert: check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 values (1,"1","1"),(2,"2","2");
`,
		newTestResult().addResult(rulepkg.DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2");
`,
		newTestResult(),
	)
}

func TestCheckInsertColumnsExist_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]insert: check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 values (?,?,?),(?,?,?);
`,
		newTestResult().addResult(rulepkg.DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMax(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckBatchInsertListsMax].Rule
	// defult 5000,  unit testing :4
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "4")
	runSingleRuleInspectCase(rule, t, "insert:check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4"),(5,"5","5");
`,
		newTestResult().addResult(rulepkg.DMLCheckBatchInsertListsMax, 4),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4");
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMax_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckBatchInsertListsMax].Rule
	// defult 5000, unit testing :4
	//rule.Value = "4"
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "4")
	runSingleRuleInspectCase(rule, t, "[fp]insert:check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult().addResult(rulepkg.DMLCheckBatchInsertListsMax, 4),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckPkProhibitAutoIncrement(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKProhibitAutoIncrement].Rule
	runSingleRuleInspectCase(rule, t, "create_table: primary key not auto increment", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT DEFAULT "unit test" COMMENT "unit test" ,
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckPKProhibitAutoIncrement),
	)

	{
		inspect := DefaultMysqlInspect()
		runSingleRuleInspectCase(rule, t, "create_table: passing the primary key not auto increment", inspect,
			`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)

		inspect1 := DefaultMysqlInspect()
		inspect1.Ctx = inspect.Ctx
		runSingleRuleInspectCase(rule, t, "alter table", inspect1,
			`
ALTER TABLE exist_db.not_exist_tb_1 modify COLUMN id BIGINT auto_increment;
ALTER TABLE exist_db.not_exist_tb_1 change COLUMN id new_id bigint unsigned NOT NULL auto_increment;
`,
			newTestResult().addResult(rulepkg.DDLCheckPKProhibitAutoIncrement),
			newTestResult().addResult(rulepkg.DDLCheckPKProhibitAutoIncrement))
	}

	{
		inspect := DefaultMysqlInspect()
		runSingleRuleInspectCase(rule, t, "create_table", inspect,
			`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)
		inspect1 := DefaultMysqlInspect()
		inspect1.Ctx = inspect.Ctx
		runSingleRuleInspectCase(rule, t, "alter table", inspect1,
			`
ALTER TABLE exist_db.not_exist_tb_1 modify COLUMN id BIGINT;
ALTER TABLE exist_db.not_exist_tb_1 change COLUMN id new_id bigint unsigned NOT NULL;
`,
			newTestResult(),
			newTestResult())
	}

	{
		inspect := DefaultMysqlInspect()
		runSingleRuleInspectCase(rule, t, "create_table", inspect,
			`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult())

		inspect1 := DefaultMysqlInspect()
		inspect1.Ctx = inspect.Ctx
		runSingleRuleInspectCase(rule, t, "alter table: Add column should error", inspect1,
			`
ALTER TABLE exist_db.not_exist_tb_1 Add COLUMN id bigint unsigned PRIMARY KEY NOT NULL;
`,
			newTestResult())
	}

	{
		inspect := DefaultMysqlInspect()
		runSingleRuleInspectCase(rule, t, "create_table", inspect,
			`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult())

		inspect1 := DefaultMysqlInspect()
		inspect1.Ctx = inspect.Ctx
		runSingleRuleInspectCase(rule, t, "alter table: Add column should error", inspect1,
			`
ALTER TABLE exist_db.not_exist_tb_1 Add COLUMN id bigint unsigned PRIMARY KEY NOT NULL AUTO_INCREMENT;
`,
			newTestResult().addResult(rulepkg.DDLCheckPKProhibitAutoIncrement))
	}
}

func TestCheckWhereExistFunc(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistFunc].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where nvl(v2,"0") = "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistFunc),
	)

	runSingleRuleInspectCase(rule, t, "select: passing the check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 = "3"
`,
		newTestResult(),
	)
}

func TestCheckWhereExistFunc_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistFunc].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where nvl(v2,?) = ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistFunc),
	)

	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 = ?
`,
		newTestResult(),
	)
}

func TestCheckWhereExistNot(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistNot].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist <> ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 <> "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not like ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 not like "%3%";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist != ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 != "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not null ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 is not null;
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
}

func TestCheckWhereExistImplicitConversion(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistImplicitConversion].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = 3;
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistImplicitConversion),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = "3";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "select: check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where id = "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistImplicitConversion),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where id = 3;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistImplicitConversion_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistImplicitConversion].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: unable to check implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = ?;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: unable to check implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where id = ?;
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExist(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "delete: check limit must exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "delete: passing the check limit must exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 limit 10 ;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "update: check limit must exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1 ="1";
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "update: passing the check limit must exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1 ="1" limit 10 ;
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExist_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]delete: check limit must exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "[fp]delete: passing the check limit must exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 limit ? ;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[fp]update: check limit must exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1 =?;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "[fp]update: passing the check limit must exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1 =? limit ? ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueries(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from  exist_db.exist_tb_2);
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select a.v1 from exist_db.exist_tb_1 a, exist_db.exist_tb_2 b  where a.v1 = b.v1 ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueries_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from exist_db.exist_tb_2 where v1 = ?);
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (?);
`,
		newTestResult(),
	)
}

func TestCheckIndexesExistBeforeCreatConstraints(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexesExistBeforeCreateConstraints].Rule
	runSingleRuleInspectCase(rule, t, "Add unique: check indexes exist before creat constraints", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_3 Add unique uniq_test(v2);
`, /*not exist index*/
		newTestResult().addResult(rulepkg.DDLCheckIndexesExistBeforeCreateConstraints),
	)
	runSingleRuleInspectCase(rule, t, "Add unique: passing the check indexes exist before creat constraints", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_1 Add unique uniq_test(v1); 
`, /*exist index*/
		newTestResult(),
	)
}

func TestCheckSelectForUpdate(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "select : check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 for update;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectForUpdate),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1; 
`,
		newTestResult(),
	)
}

func TestCheckSelectForUpdate_FP(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select : check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = ? for update;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectForUpdate),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1= ?; 
`,
		newTestResult(),
	)
}

func TestCheckCollationDatabase(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`:    `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8_general_ci;`,
		`alter table`:     `ALTER TABLE exist_db.exist_tb_1 COLLATE utf8_general_ci;`,
		`create database`: `CREATE DATABASE db COLLATE utf8_general_ci;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseCollation].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckDatabaseCollation, "utf8mb4_0900_ai_ci"))
	}

	for desc, sql := range map[string]string{
		`create table`:    `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8mb4_0900_ai_ci;`,
		`alter table`:     `ALTER TABLE exist_db.exist_tb_1 COLLATE utf8mb4_0900_ai_ci;`,
		`create database`: `CREATE DATABASE db COLLATE utf8mb4_0900_ai_ci;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseCollation].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckDecimalTypeColumn(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckDecimalTypeColumn].Rule
	runSingleRuleInspectCase(rule, t, "create table: check decimal type column", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 float(10));
`,
		newTestResult().addResult(rulepkg.DDLCheckDecimalTypeColumn),
	)
	runSingleRuleInspectCase(rule, t, "alter table: check decimal type column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 FLOAT ( 10, 0 );
`,
		newTestResult().addResult(rulepkg.DDLCheckDecimalTypeColumn),
	)
	runSingleRuleInspectCase(rule, t, "create table: passing the check decimal type column", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 DECIMAL);
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "alter table: passing the check decimal type column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 DECIMAL;
`,
		newTestResult(),
	)

}

func TestCheckColumnTypeBlobText(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)create table`: `CREATE TABLE t1(id BLOB);`,
		`(2)create table`: `CREATE TABLE t1(id TINYBLOB);`,
		`(3)create table`: `CREATE TABLE t1(id MEDIUMBLOB);`,
		`(4)create table`: `CREATE TABLE t1(id LONGBLOB);`,
		`(5)create table`: `CREATE TABLE t1(id TEXT);`,
		`(6)create table`: `CREATE TABLE t1(id TINYTEXT);`,
		`(7)create table`: `CREATE TABLE t1(id MEDIUMTEXT);`,
		`(8)create table`: `CREATE TABLE t1(id LONGTEXT);`,
		`(1)alter table`:  `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
		`(2)alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 BLOB;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnBlobNotice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckColumnBlobNotice))
	}

	for desc, sql := range map[string]string{
		`(1)create table`: `CREATE TABLE t1(id INT);`,
		`(1)alter table`:  `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 VARCHAR(100);`,
		`(2)alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 SET('male', 'female');`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnBlobNotice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckColumnTypeSet(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT, sex SET("male", "female"));`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 SET("male", "female");`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 SET("male", "female");`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 SET("male", "female");`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnSetNitice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckColumnSetNitice))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 INT;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB;`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnSetNitice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckColumnTypeEnum(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT, sex ENUM("male", "female"));`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 ENUM("male", "female");`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 ENUM("male", "female");`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 ENUM("male", "female");`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnEnumNotice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckColumnEnumNotice))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 BLOB;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnEnumNotice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckUniqueIndex(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`: `CREATE TABLE t1(id INT, c1 INT, UNIQUE INDEX unique_idx (c1));`,
		`alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX unique_id(id);`,
		`create index`: `CREATE UNIQUE INDEX i_u_id ON exist_db.exist_tb_1(id);`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndex].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckUniqueIndex))
	}

	for desc, sql := range map[string]string{
		`create table`: `
CREATE TABLE t1(
id INT,
c1 INT,
c2 INT,
UNIQUE INDEX idx_uk_t1_c1 (c1),
UNIQUE INDEX IDX_UK_t1_c1_c2 (c1, c2),
INDEX idx_id(id)
);
`,
		`alter table`: `
ALTER TABLE exist_db.exist_tb_1
ADD UNIQUE INDEX idx_uk_exist_tb_1_v1(v1),
ADD UNIQUE INDEX IDX_UK_exist_tb_1_id_v1(id, v1),
ADD INDEX idx_v2(v2);
`,
		`(1)create index`: `CREATE UNIQUE INDEX idx_uk_exist_tb_1_id_v1 ON exist_db.exist_tb_1(id, v1);`,
		`(2)create index`: `CREATE UNIQUE INDEX IDX_UK_exist_tb_1_id ON exist_db.exist_tb_1(id);`,
		`(3)create index`: `CREATE INDEX idx_id ON exist_db.exist_tb_1(id);`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndex].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckWhereExistNull(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id IS NULL;`,
		`(2)select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id IS NOT NULL;`,
		`(1)update table`: `UPDATE exist_db.exist_tb_1 SET id = 1 WHERE id IS NULL;`,
		`(2)update table`: `UPDATE exist_db.exist_tb_1 SET id = 1 WHERE id IS NOT NULL;`,
		`(1)delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id IS NULL;`,
		`(2)delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id IS NOT NULL;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLWhereExistNull].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DMLWhereExistNull))
	}

	for desc, sql := range map[string]string{
		`select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id = 1;`,
		`update table`: `UPDATE exist_db.exist_tb_1 SET id = 10 WHERE id = 1;`,
		`delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id = 1;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLWhereExistNull].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckWhereExistNull_FP(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp]select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id = ?;`,
		`[fp]update table`: `UPDATE exist_db.exist_tb_1 SET id = 10 WHERE id = ?;`,
		`[fp]delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id = ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLWhereExistNull].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckNeedlessFunc(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, MD5('aaa'), MD5('bbb'));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DMLCheckNeedlessFunc, "sha(),sqrt(),md5()"))
	}

	for desc, sql := range map[string]string{
		`(1)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, sha1('aaa'), sha1('bbb'));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckNeedlessFunc_FP(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp]INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(?, MD5(?), MD5(?));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DMLCheckNeedlessFunc, "sha(),sqrt(),md5()"))
	}

	for desc, sql := range map[string]string{
		`[fp]INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(?, sha1(?), sha1(?));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckDatabaseSuffix(t *testing.T) {
	for desc, sql := range map[string]string{
		`create database`: `CREATE DATABASE app_service;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseSuffix].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckDatabaseSuffix, "_DB"))
	}

	for desc, sql := range map[string]string{
		`(0)create database`: `CREATE DATABASE app_service_db;`,
		`(1)create database`: `CREATE DATABASE APP_SERVICE_DB;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseSuffix].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckTransactionIsolationLevel(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)transaction isolation should notice`: `SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;`,
		`(2)transaction isolation should notice`: `SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;`,
		`(3)transaction isolation should notice`: `SET GLOBAL TRANSACTION ISOLATION LEVEL REPEATABLE READ;`,
		`(4)transaction isolation should notice`: `SET GLOBAL TRANSACTION READ ONLY, ISOLATION LEVEL SERIALIZABLE;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckTransactionIsolationLevel].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckTransactionIsolationLevel))
	}

	for desc, sql := range map[string]string{
		`(1)transaction isolation should not notice`: `SET TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(2)transaction isolation should not notice`: `SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(3)transaction isolation should not notice`: `SET GLOBAL TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(4)transaction isolation should not notice`: `SET GLOBAL TRANSACTION READ ONLY;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckTransactionIsolationLevel].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckFuzzySearch(t *testing.T) {
	for _, sql := range []string{
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '_a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a';`,

		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '%a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '%a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '_a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '%a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '%a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '_a';`,

		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '_a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a';`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckFuzzySearch))
	}

	for _, sql := range []string{
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a___';`,

		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a___';`,

		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a____';`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestCheckFuzzySearch_FP(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp] "select" unable to check fuzzy search`: `SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
		`[fp] "update" unable to check fuzzy search`: `UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE ?;`,
		`[fp] "delete" unable to check fuzzy search`: `DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckTablePartition(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table should error`: `
CREATE TABLE t1(
c1 INT,
c2 INT)
PARTITION BY LIST(c1)
(
PARTITION p1 VALUES IN(1, 2, 3),
PARTITION p2 VALUES IN(4, 5, 6),
PARTITION p3 VALUES IN(7, 8, 9)
)
`,
		`alter table should error`: `
ALTER TABLE exist_db.exist_tb_1
PARTITION BY LIST(v1)
(
PARTITION p1 VALUES IN(1, 2, 3),
PARTITION p2 VALUES IN(4, 5, 6),
PARTITION p3 VALUES IN(7, 8, 9)
)
`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckTablePartition].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckTablePartition))
	}

	for desc, sql := range map[string]string{
		`create table should not error`: `
CREATE TABLE t1(
c1 INT,
c2 INT)
`,
		`alter table should not error`: `
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 INT;
`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckTablePartition].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckNumberOfJoinTables(t *testing.T) {
	// create table for JOIN test
	inspector := DefaultMysqlInspect()
	{
		parent := DefaultMysqlInspect()
		runDefaultRulesInspectCase(t, "create table for JOIN test", parent,
			`
create table if not exists exist_db.exist_tb_4 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)
		inspector.Ctx = session.NewContext(parent.Ctx)
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT * FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON exist_db.exist_tb_1.id = exist_db.exist_tb_2.id 
JOIN exist_db.exist_tb_3 ON exist_db.exist_tb_2.id = exist_db.exist_tb_3.id
JOIN exist_db.exist_tb_4 ON exist_db.exist_tb_3.id = exist_db.exist_tb_4.id
`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			inspector,
			sql,
			newTestResult().addResult(rulepkg.DMLCheckNumberOfJoinTables, 3))
	}

	for desc, sql := range map[string]string{
		`(1)select table should not error`: `
		SELECT * FROM exist_db.exist_tb_1
		`,
		`(2)select table should not error`: `
SELECT * FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON exist_db.exist_tb_1.id = exist_db.exist_tb_2.id 
JOIN exist_db.exist_tb_3 ON exist_db.exist_tb_2.id = exist_db.exist_tb_3.id
		`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckNumberOfJoinTables_FP(t *testing.T) {
	// create table for JOIN test
	inspector := DefaultMysqlInspect()
	{
		parent := DefaultMysqlInspect()
		runDefaultRulesInspectCase(t, "create table for JOIN test", parent,
			`
create table if not exists exist_db.exist_tb_4 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)
		inspector.Ctx = session.NewContext(parent.Ctx)
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT * FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON exist_db.exist_tb_1.id = exist_db.exist_tb_2.id 
JOIN exist_db.exist_tb_3 ON exist_db.exist_tb_2.id = exist_db.exist_tb_3.id
JOIN exist_db.exist_tb_4 ON exist_db.exist_tb_3.id = exist_db.exist_tb_4.id
WHERE exist_db.exist_tb_1.v1 = ? AND exist_db.exist_tb_1.v2 = ?
`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			inspector,
			sql,
			newTestResult().addResult(rulepkg.DMLCheckNumberOfJoinTables, 3))
	}

	for desc, sql := range map[string]string{
		`(1)select table should not error`: `
		SELECT * FROM exist_db.exist_tb_1 WHERE exist_db.exist_tb_1.v1 = ?
		`,
		`(2)select table should not error`: `
SELECT * FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON exist_db.exist_tb_1.id = exist_db.exist_tb_2.id 
JOIN exist_db.exist_tb_3 ON exist_db.exist_tb_2.id = exist_db.exist_tb_3.id
WHERE exist_db.exist_tb_1.v1 = ? AND exist_db.exist_tb_1.v2 = ?
		`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckIsAfterUnionDistinct(t *testing.T) {
	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT 1, 2 UNION SELECT 'a', 'b';`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DMLCheckIfAfterUnionDistinct))
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT 1, 2 UNION ALL SELECT 'a', 'b';`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckIsAfterUnionDistinct_FP(t *testing.T) {
	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT ?, ? UNION SELECT ?, ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DMLCheckIfAfterUnionDistinct))
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT ?, ? UNION ALL SELECT ?, ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckIsExistLimitOffset(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)select table should error`: `
SELECT * FROM exist_db.exist_tb_1 LIMIT 5,6;`,
		`(2)select table should error`: `
SELECT * FROM exist_db.exist_tb_1 LIMIT 6 OFFSET 5;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckIsExistLimitOffset].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckIsExistLimitOffset))
	}

	for desc, sql := range map[string]string{
		`select table should not error`: `
SELECT * FROM exist_db.exist_tb_1 LIMIT 5`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckIsExistLimitOffset].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func Test_DDLCheckNameUseENAndUnderline_ShouldError(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)create database`: `CREATE DATABASE ®®;`,
		`(2)create database`: `CREATE DATABASE _app;`,
		`(3)create database`: `CREATE DATABASE 1_app;`,
		`(0)create table`:    `CREATE TABLE 应用1(字段1 int);`,
		`(1)create table`:    `CREATE TABLE ®®(®® int);`,
		`(2)create table`:    `CREATE TABLE _app(_col int);`,
		`(3)create table`:    `CREATE TABLE _app(col_ int);`,
		`(4)create table`:    `CREATE TABLE 1_app(col_ int);`,
		`(0)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN 字段 int;`,
		`(1)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN _col int;`,
		`(3)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN ®® int;`,
		`(4)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN 1_Col int;`,
		`(0)create index`:    `CREATE INDEX 索引1 ON exist_db.exist_tb_1(v1)`,
		`(1)create index`:    `CREATE INDEX _idx ON exist_db.exist_tb_1(v1)`,
		`(3)create index`:    `CREATE INDEX ®® ON exist_db.exist_tb_1(v1)`,
		`(4)create index`:    `CREATE INDEX 1_idx ON exist_db.exist_tb_1(v1)`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckOBjectNameUseCN].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckOBjectNameUseCN))
	}
}

func Test_DDLCheckNameUseENAndUnderline_ShouldNotError(t *testing.T) {
	for desc, sql := range map[string]string{
		`(0)create database`: `CREATE DATABASE db_app1;`,
		`(1)create database`: `CREATE DATABASE app_;`,
		`(0)create table`:    `CREATE TABLE tb_service1(pk_id int);`,
		`(0)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v4_col4 int;`,
		`(1)alter table`:     `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN col_ int;`,
		`(0)create index`:    `CREATE INDEX idx_v1 ON exist_db.exist_tb_1(v1)`,
		`(1)create index`:    `CREATE INDEX idx_ ON exist_db.exist_tb_1(v1)`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckOBjectNameUseCN].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func TestCheckIndexOption_ShouldNot_QueryDB(t *testing.T) {
	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(1)index on new db new column`,
		DefaultMysqlInspect(),
		`CREATE TABLE t1(id int, name varchar(100), INDEX idx_name(name))`,
		newTestResult())

	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(2)index on new db new column`,
		DefaultMysqlInspect(),
		`CREATE TABLE t1(id int, name varchar(100));
ALTER TABLE t1 ADD INDEX idx_name(name);
`,
		newTestResult(), newTestResult())

	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(3)index on old db new column`,
		DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(100);
ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_v3(v3);
`,
		newTestResult(), newTestResult())
}

func Test_CheckExplain_ShouldNotError(t *testing.T) {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect1 := NewMockInspect(e)

	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows"}).AddRow("ALL", "10"))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainAccessTypeAll].Rule, t, "", inspect1, "select * from exist_tb_1", newTestResult())

	inspect2 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows"}).AddRow("ALL", "10"))
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingFilesort].Rule, t, "", inspect2, "select * from exist_tb_1", newTestResult())

	inspect3 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows"}).AddRow("ALL", "10"))
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingFilesort].Rule, t, "", inspect3, "select * from exist_tb_1", newTestResult())

	assert.NoError(t, handler.ExpectationsWereMet())
}

func Test_CheckExplain_ShouldError(t *testing.T) {
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect1 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows"}).
			AddRow("ALL", "10001"))
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainAccessTypeAll].Rule, t, "", inspect1, "select * from exist_tb_1", newTestResult().addResult(rulepkg.DMLCheckExplainAccessTypeAll, 10001))

	inspect2 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows", "Extra"}).
			AddRow("ALL", "10", executor.ExplainRecordExtraUsingTemporary))
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingTemporary].Rule, t, "", inspect2, "select * from exist_tb_1", newTestResult().addResult(rulepkg.DMLCheckExplainExtraUsingTemporary))

	inspect3 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows", "Extra"}).
			AddRow("ALL", "10", executor.ExplainRecordExtraUsingFilesort))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingFilesort].Rule, t, "", inspect3, "select * from exist_tb_1", newTestResult().addResult(rulepkg.DMLCheckExplainExtraUsingFilesort))

	inspect4 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows", "Extra"}).
			AddRow("ALL", "100001", strings.Join([]string{executor.ExplainRecordExtraUsingFilesort, executor.ExplainRecordExtraUsingTemporary}, ";")))

	ruleDMLCheckExplainExtraUsingFilesort := rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingFilesort].Rule
	ruleDMLCheckExplainExtraUsingTemporary := rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainExtraUsingTemporary].Rule
	ruleDMLCheckExplainAccessTypeAll := rulepkg.RuleHandlerMap[rulepkg.DMLCheckExplainAccessTypeAll].Rule

	inspect4.rules = []*driver.Rule{
		&ruleDMLCheckExplainExtraUsingFilesort,
		&ruleDMLCheckExplainExtraUsingTemporary,
		&ruleDMLCheckExplainAccessTypeAll}

	inspectCase(t, "", inspect4, "select * from exist_tb_1",
		newTestResult().addResult(rulepkg.DMLCheckExplainExtraUsingFilesort).addResult(rulepkg.DMLCheckExplainExtraUsingTemporary).addResult(rulepkg.DMLCheckExplainAccessTypeAll, 100001))

	inspect5 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1")).
		WillReturnRows(sqlmock.NewRows([]string{"type", "rows"}).
			AddRow("ALL", "100001"))

	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1 where id = 1;")).
		WillReturnRows(sqlmock.NewRows([]string{"Extra"}).
			AddRow(executor.ExplainRecordExtraUsingFilesort))

	handler.ExpectQuery(regexp.QuoteMeta("select * from exist_tb_1 where id = 2;")).
		WillReturnRows(sqlmock.NewRows([]string{"Extra"}).
			AddRow(executor.ExplainRecordExtraUsingTemporary))

	inspect5.rules = []*driver.Rule{
		&ruleDMLCheckExplainExtraUsingFilesort,
		&ruleDMLCheckExplainExtraUsingTemporary,
		&ruleDMLCheckExplainAccessTypeAll}

	inspectCase(t, "", inspect5, "select * from exist_tb_1;select * from exist_tb_1 where id = 1;select * from exist_tb_1 where id = 2;",
		newTestResult().addResult(rulepkg.DMLCheckExplainAccessTypeAll, 100001), newTestResult().addResult(rulepkg.DMLCheckExplainExtraUsingFilesort), newTestResult().addResult(rulepkg.DMLCheckExplainExtraUsingTemporary))

	assert.NoError(t, handler.ExpectationsWereMet())
}

func Test_DDL_CHECK_PK_NAME(t *testing.T) {
	for _, sql := range []string{
		`create table t1(id int, primary key pk_t1(id))`,
		`create table t1(id int, primary key PK_T1(id))`,
		`create table t1(id int, primary key(id))`,
		`alter table exist_db.exist_tb_2 Add primary key(id)`,
		`alter table exist_db.exist_tb_2 Add primary key PK_EXIST_TB_2(id)`} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKName].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}

	for _, sql := range []string{
		`create table t1(id int, primary key wrongPK(id))`,
		`alter table exist_db.exist_tb_2 Add primary key wrongPK(id)`} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKName].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckPKName))
	}
}

func Test_PerfectParse(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereIsInvalid].Rule, t, "", DefaultMysqlInspect(), `
SELECT * FROM exist_db.exist_tb_1;
OPTIMIZE TABLE exist_db.exist_tb_1;
SELECT * FROM exist_db.exist_tb_2;
`, newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
		newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持"),
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))
}

func Test_DDLCheckCreateView(t *testing.T) {
	for _, sql := range []string{
		`create view v as select * from t1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateView].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckCreateView))
	}

	for _, sql := range []string{
		`create table t1(id int)`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateView].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func Test_DDLCheckCreateTrigger(t *testing.T) {
	for _, sql := range []string{
		`create trigger my_trigger before insert on t1 for each row insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE DEFINER='sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE DEFINER = 'sqle_op'@'localhost' TRIGGER my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`
CREATE
	DEFINER = 'sqle_op'@'localhost' 
	TRIGGER my_trigger 
	BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');
`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持").addResult(rulepkg.DDLCheckCreateTrigger))
	}

	for _, sql := range []string{
		`CREATE my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATEmy_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE trigger_1 BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER my_trigger BEEEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持"))
	}
}

func Test_DDLCheckCreateFunction(t *testing.T) {
	for _, sql := range []string{
		`create function hello_function (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE FUNCTION hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER='sqle_op'@'localhost' FUNCTION hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER = 'sqle_op'@'localhost' FUNCTION hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`
CREATE
	DEFINER = 'sqle_op'@'localhost' 
	FUNCTION hello_function (s CHAR(20)) 
	RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');
`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持").addResult(rulepkg.DDLCheckCreateFunction))
	}

	for _, sql := range []string{
		`create function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`create123 function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER='sqle_op'@'localhost' hello (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持"))
	}
}

func Test_DDLCheckCreateProcedure(t *testing.T) {
	for _, sql := range []string{
		`
CREATE DEFINER='sqle_op'@'localhost'
PROCEDURE proc1(OUT s int) COMMENT 'test'
BEGIN
SELECT * FROM t1;
SELECT * FROM t2;
END;`,
		`
create definer='sqle_op'@'localhost'
procedure proc1(out s int) comment 'test'
begin
select * from t1;
select * from t2;
end;`,
		`
create procedure proc1()
begin
select * from t1;
select * from t2;
end;`,
		`
create procedure proc1()
begin
end;`,
		`
create procedure proc1()
select * from t1;`,
		`
create 
procedure
proc1()
select * from t1;`,
		`
create 
	procedure
proc1()
select * from t1;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateProcedure].Rule, t, "",
			DefaultMysqlInspect(), sql,
			newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持").
				addResult(rulepkg.DDLCheckCreateProcedure))
	}

	for _, sql := range []string{
		`
CREATE DEFINER='sqle_op'@'localhost'PROCEDURE proc1(OUT s int) COMMENT 'test'
BEGIN
SELECT * FROM t1;
SELECT * FROM t2;
END;`,
		`
createdefiner='sqle_op'@'localhost' procedure proc1(out s int) comment 'test'
begin
select * from t1;
select * from t2;
end;`,
		`
create procedureproc1()
begin
select * from t1;
select * from t2;
end;`,
		`
createprocedure proc1()
begin
end;`,
		`
create123 procedure proc1()
begin
end;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateProcedure].Rule, t, "",
			DefaultMysqlInspect(), sql,
			newTestResult().add(driver.RuleLevelError, "语法错误或者解析器不支持"))
	}
}

// todo(@wy): move to auto test
func TestWhitelist(t *testing.T) {
	//	for _, sql := range []string{
	//		"select v1 from exist_tb_1 where id =2",
	//		"select v1 from exist_tb_1 where id =\"2\"",
	//		"select v1 from exist_tb_1 where id =100000",
	//	} {
	//		runDefaultRulesInspectCaseWithWL(t, "match fp whitelist", DefaultMysqlInspect(),
	//			[]driver.SqlWhitelist{
	//				{
	//					Value:     "select v1 from exist_tb_1 where id =2",
	//					MatchType: driver.SQLWhitelistFPMatch,
	//				},
	//			}, sql, newTestResult().add(driver.RuleLevelNormal, "白名单"))
	//	}
	//
	//	for _, sql := range []string{
	//		"select v1 from exist_tb_1 where ID =2",
	//		"select v1 from exist_tb_1 where id =2 and v2 = \"test\"",
	//	} {
	//		runDefaultRulesInspectCaseWithWL(t, "don't match fp whitelist", DefaultMysqlInspect(),
	//			[]driver.SqlWhitelist{
	//				{
	//					Value:     "select v1 from exist_tb_1 where id =2",
	//					MatchType: driver.SQLWhitelistFPMatch,
	//				},
	//			}, sql, newTestResult())
	//	}
	//
	//	for _, sql := range []string{
	//		"select v1 from exist_tb_1 where id = 1",
	//		"SELECT V1 FROM exist_tb_1 WHERE ID = 1",
	//	} {
	//		runDefaultRulesInspectCaseWithWL(t, "match exact whitelist", DefaultMysqlInspect(),
	//			[]driver.SqlWhitelist{
	//				driver.SqlWhitelist{
	//					CapitalizedValue: "SELECT V1 FROM EXIST_TB_1 WHERE ID = 1",
	//					MatchType:        driver.SQLWhitelistExactMatch,
	//				},
	//			}, sql,
	//			newTestResult().add(driver.RuleLevelNormal, "白名单"))
	//	}
	//
	//	for _, sql := range []string{
	//		"select v1 from exist_tb_1 where id = 2",
	//		"SELECT V1 FROM exist_tb_1 WHERE ID = 2",
	//	} {
	//		runDefaultRulesInspectCaseWithWL(t, "don't match exact whitelist", DefaultMysqlInspect(),
	//			[]driver.SqlWhitelist{
	//				driver.SqlWhitelist{
	//					CapitalizedValue: "SELECT V1 FROM EXIST_TB_1 WHERE ID = 1",
	//					MatchType:        driver.SQLWhitelistExactMatch,
	//				},
	//			}, sql,
	//			newTestResult())
	//	}
	//
	//	parentInspect := DefaultMysqlInspect()
	//	runDefaultRulesInspectCase(t, "", parentInspect, `
	//CREATE TABLE if not exists exist_db.t1 (
	//id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	//PRIMARY KEY (id) USING BTREE
	//)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	//`, newTestResult())
	//
	//	inspect1 := DefaultMysqlInspect()
	//	inspect1.Ctx = parentInspect.Ctx
	//
	//	runDefaultRulesInspectCaseWithWL(t, "4", inspect1,
	//		[]driver.SqlWhitelist{
	//			{
	//				Value:     "select * from t1 where id = 2",
	//				MatchType: driver.SQLWhitelistFPMatch,
	//			},
	//		}, `select id from T1 where id = 4`, newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.T1"))
	//
	//	inspect2 := DefaultMysqlInspect()
	//	inspect2.Ctx = parentInspect.Ctx
	//	inspect2.Ctx.AddSysVar("lower_case_table_names", "1")
	//	runDefaultRulesInspectCaseWithWL(t, "", inspect2,
	//		[]driver.SqlWhitelist{
	//			{
	//				Value:     "select * from t1 where id = 2",
	//				MatchType: driver.SQLWhitelistFPMatch,
	//			},
	//		}, `select * from T1 where id = 3`, newTestResult().add(driver.RuleLevelNormal, "白名单"))
	//
}

// func runDefaultRulesInspectCaseWithWL(t *testing.T, desc string, i *Inspect,
// 	wl []driver.SqlWhitelist, sql string, results ...*testResult) {
// 	var ptrRules []*driver.Rule
// 	for i := range DefaultTemplateRules {
// 		// remove DDL_CHECK_OBJECT_NAME_USING_CN in default rules for init test.
// 		if DefaultTemplateRules[i].Name == DDLCheckOBjectNameUseCN {
// 			continue
// 		}

// 		ptrRules = append(ptrRules, &DefaultTemplateRules[i])
// 	}

// 	i.rules = ptrRules
// 	inspectCase(t, desc, i, wl, sql, results...)
// }
