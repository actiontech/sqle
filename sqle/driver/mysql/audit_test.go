package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/DATA-DOG/go-sqlmock"
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
	handler, ok := rulepkg.RuleHandlerMap[ruleName]
	if !ok {
		panic("should not enter here, it means that the uint test result is not expect")
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

func DefaultMysqlInspect() *MysqlDriverImpl {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &MysqlDriverImpl{
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
			DDLGhostMinSize:    -1,
			DMLRollbackMaxRows: 1000,
		},
	}
}

func NewMockInspect(e *executor.Executor) *MysqlDriverImpl {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &MysqlDriverImpl{
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

func runSingleRuleInspectCase(rule driver.Rule, t *testing.T, desc string, i *MysqlDriverImpl, sql string, results ...*testResult) {
	i.rules = []*driver.Rule{&rule}
	inspectCase(t, desc, i, sql, results...)
}

func runDefaultRulesInspectCase(t *testing.T, desc string, i *MysqlDriverImpl, sql string, results ...*testResult) {
	ptrRules := []*driver.Rule{}
	// this rule will be test in single rule
	filterRule := map[string]struct{}{
		rulepkg.DDLCheckObjectNameUseCN:                     struct{}{},
		rulepkg.DDLCheckRedundantIndex:                      struct{}{},
		rulepkg.DDLCheckPKProhibitAutoIncrement:             struct{}{},
		rulepkg.DDLCheckColumnBlobNotice:                    struct{}{},
		rulepkg.DDLCheckDatabaseCollation:                   struct{}{},
		rulepkg.DDLCheckIndexTooMany:                        struct{}{},
		rulepkg.DDLCheckIndexesExistBeforeCreateConstraints: struct{}{},
		rulepkg.DMLCheckInsertColumnsExist:                  struct{}{},
		rulepkg.DMLCheckLimitMustExist:                      struct{}{},
		rulepkg.DMLCheckWhereExistImplicitConversion:        struct{}{},
		rulepkg.DMLCheckSQLLength:                           {},
		rulepkg.DDLRecommendTableColumnCharsetSame:          {},
		rulepkg.DDLCheckAutoIncrement:                       {},
		rulepkg.DDLCheckColumnTypeInteger:                   {},
		rulepkg.DDLHintDropColumn:                           {},
		rulepkg.DMLHintDeleteTips:                           {},
		rulepkg.DMLHintUseTruncateInsteadOfDelete:           {},
		rulepkg.DDLCheckColumnQuantity:                      {},
		rulepkg.DMLHintInNullOnlyFalse:                      {},
		rulepkg.DMLNotRecommendIn:                           {},
		rulepkg.DMLCheckAlias:                               {},
	}
	for i := range rulepkg.RuleHandlers {
		handler := rulepkg.RuleHandlers[i]
		if _, ok := filterRule[handler.Rule.Name]; ok {
			continue
		}
		ptrRules = append(ptrRules, &handler.Rule)
	}

	i.rules = ptrRules
	inspectCase(t, desc, i, sql, results...)
}

func runEmptyRuleInspectCase(t *testing.T, desc string, i *MysqlDriverImpl, sql string, results ...*testResult) {
	i.rules = []*driver.Rule{}
	inspectCase(t, desc, i, sql, results...)
}

func inspectCase(t *testing.T, desc string, i *MysqlDriverImpl, sql string, results ...*testResult) {
	stmts, err := util.ParseSql(sql)
	if err != nil {
		t.Errorf("%s test failed, error: %v\n", desc, err)
		return
	}

	if len(stmts) != len(results) {
		t.Errorf("%s test failed, error: result is unknow\n", desc)
		return
	}

	for idx, stmt := range stmts {
		result, err := i.Audit(context.TODO(), stmt.Text())
		if err != nil {
			t.Error(err)
			return
		}
		if result.Level() != results[idx].level() || result.Message() != results[idx].message() {
			t.Errorf("%s test failed, \n\nsql:\n %s\n\nexpect level: %s\nexpect result:\n%s\n\nactual level: %s\nactual result:\n%s\n",
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
	inspect1.Ctx.AddSystemVariable(session.SysVarLowerCaseTableNames, "1")
	runDefaultRulesInspectCase(t, "", inspect1,
		"use EXIST_DB",
		newTestResult(),
	)
}

func TestCaseSensitive(t *testing.T) {
	runDefaultRulesInspectCase(t, "", DefaultMysqlInspect(),
		`
select id from exist_db.EXIST_TB_1 where id = 1 limit 1;
`,
		newTestResult().add(driver.RuleLevelError,
			TableNotExistMessage, "exist_db.EXIST_TB_1").add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"))

	inspect1 := DefaultMysqlInspect()
	inspect1.Ctx.AddSystemVariable(session.SysVarLowerCaseTableNames, "1")
	runDefaultRulesInspectCase(t, "", inspect1,
		`
select id from exist_db.EXIST_TB_1 where id = 1 limit 1;
`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"))
}

func TestDDLCheckTableSize(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableSize].Rule

	runSingleRuleInspectCase(rule, t, "drop_table: table1 oversized", DefaultMysqlInspect(),
		`drop table exist_db.exist_tb_1;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "alter_table: table1 oversized", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1;`, newTestResult())

	runSingleRuleInspectCase(rule, t, "drop_table: table4 oversized", DefaultMysqlInspect(),
		`drop table exist_db.exist_tb_4;`, newTestResult().addResult(rulepkg.DDLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "alter_table: table4 oversized", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_4;`, newTestResult().addResult(rulepkg.DDLCheckTableSize, "exist_tb_4", 16).addResult(rulepkg.ConfigDDLOSCMinSize, PTOSCNoUniqueIndexOrPrimaryKey))

}

func TestDMLCheckTableSize(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckTableSize].Rule
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "16")

	// TODO 'select from table1 , table2 ;' There is currently no single test, because this sql sqle cannot be supported as of the time of writing the comment
	runSingleRuleInspectCase(rule, t, "select: table1 oversized", DefaultMysqlInspect(),
		`select 1 from exist_db.exist_tb_1 where id = 1;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "update: table1 oversized", DefaultMysqlInspect(),
		`UPDATE exist_db.exist_tb_1 SET id = 0.8;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "insert: table1 oversized", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 VALUES(7500, 'A', 'SALESMAN');`, newTestResult())
	runSingleRuleInspectCase(rule, t, "delete: table1 oversized", DefaultMysqlInspect(),
		`DELETE id FROM exist_db.exist_tb_1;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "lock: table1 oversized", DefaultMysqlInspect(),
		`lock tables exist_db.exist_tb_1 read;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "selects: table1 oversized", DefaultMysqlInspect(),
		`select 1 from exist_db.exist_tb_1 join exist_db.exist_tb_2 where id = 1;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "updates: table1 oversized", DefaultMysqlInspect(),
		`UPDATE exist_db.exist_tb_1, exist_db.exist_tb_2 SET exist_db.exist_tb_2.id = exist_db.exist_tb_1.id * 0.8 WHERE exist_db.exist_tb_1.id= exist_db.exist_tb_2.id;`, newTestResult())
	runSingleRuleInspectCase(rule, t, "deletes: table1 oversized", DefaultMysqlInspect(),
		`DELETE id FROM exist_db.exist_tb_1 INNER JOIN exist_db.exist_tb_2 INNER JOIN exist_db.exist_tb_3;`, newTestResult())

	runSingleRuleInspectCase(rule, t, "select: table1 oversized", DefaultMysqlInspect(),
		`select 1 from exist_db.exist_tb_4 where id = 1;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "update: table1 oversized", DefaultMysqlInspect(),
		`UPDATE exist_db.exist_tb_4 SET id = 0.8;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "insert: table1 oversized", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_4 VALUES(7500, 'A', 'SALESMAN', 10);`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "delete: table1 oversized", DefaultMysqlInspect(),
		`DELETE id FROM exist_db.exist_tb_4;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "lock: table1 oversized", DefaultMysqlInspect(),
		`lock tables exist_db.exist_tb_4 read;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "selects: table1 oversized", DefaultMysqlInspect(),
		`select 1 from exist_db.exist_tb_4 join exist_db.exist_tb_2 where id = 1;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "updates: table1 oversized", DefaultMysqlInspect(),
		`UPDATE exist_db.exist_tb_4, exist_db.exist_tb_2 SET exist_db.exist_tb_2.id = exist_db.exist_tb_4.id * 0.8 WHERE exist_db.exist_tb_4.id= exist_db.exist_tb_2.id;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))
	runSingleRuleInspectCase(rule, t, "deletes: table1 oversized", DefaultMysqlInspect(),
		`DELETE id FROM exist_db.exist_tb_4 INNER JOIN exist_db.exist_tb_2 INNER JOIN exist_db.exist_tb_3;`, newTestResult().addResult(rulepkg.DMLCheckTableSize, "exist_tb_4", 16))

}

func TestCheckInvalidCreateTable(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: schema not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists not_exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
	handlerEngine := rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableDBEngine]
	handlerCharacter := rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableCharacterSet]
	handlerNotAllowRenaming := rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckTableDBEngine)
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckTableCharacterSet)
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLNotAllowRenaming)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableDBEngine] = handlerEngine
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableCharacterSet] = handlerCharacter
		rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming] = handlerNotAllowRenaming
	}()
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
select id from not_exist_db.not_exist_tb where id=1 limit 1;
`,
		newTestResult().add(driver.RuleLevelError, SchemaNotExistMessage, "not_exist_db").add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "select: table not exist", DefaultMysqlInspect(),
		`
select id from exist_db.not_exist_tb where id=1 limit 1;
`,
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb").add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
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
		"select id from exist_db.exist_tb_1 where id > 1 limit 1;",
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 limit 1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=1 and 2=2 limit 1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(3)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id=id limit 1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(4)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id limit 1;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
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

	// issue:691 https://github.com/actiontech/sqle/issues/691
	runDefaultRulesInspectCase(t, "where with () condition(1)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where (id = 1);",
		newTestResult())

	runDefaultRulesInspectCase(t, "where with () condition(2)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where (id = 1 and v1 = '2');",
		newTestResult())

	runDefaultRulesInspectCase(t, "where with () condition(3)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where (id = 1) and (v1 = '2');",
		newTestResult())
}

func TestCheckWhereInvalid_FP(t *testing.T) {
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id=? limit ?;",
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=? limit ?;",
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=? and 2=2 limit ?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where ?=? limit ?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
			"create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT \"unit test\","+
			"update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT \"unit test\","+
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_%s (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64).addResult(rulepkg.DDLNotAllowRenaming),
	)

	runDefaultRulesInspectCase(t, "alter_table:Add column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:change column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64).addResult(rulepkg.DDLNotAllowRenaming),
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

func TestCheckObjectNameIsUpperAndLowerLetterMixed(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckObjectNameIsUpperAndLowerLetterMixed].Rule

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db._Ab (
	Id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	NAME varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	A varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id),
	INDEX idx_ID_Name (id,name)
	)`, newTestResult().addResult(rule.Name, strings.Join([]string{"_Ab", "Id", "idx_ID_Name"}, ",")))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add column name varchar(255) NOT NULL DEFAULT "unit test"`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add column Name varchar(255) NOT NULL DEFAULT "unit test"`,
		newTestResult().addResult(rule.Name, "Name"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 rename test`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 rename Test`,
		newTestResult().addResult(rule.Name, "Test"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 change id id_test int unsigned NOT NULL AUTO_INCREMENT`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 change id id_Test int unsigned NOT NULL AUTO_INCREMENT`,
		newTestResult().addResult(rule.Name, "id_Test"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add constraint id_unique unique (v2)`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add constraint iD_unique unique (v2)`,
		newTestResult().addResult(rule.Name, "iD_unique"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 rename index idx_1 to idx_test`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 rename index idx_1 to idx_Test`,
		newTestResult().addResult(rule.Name, "idx_Test"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`create index i on exist_db.exist_tb_1 (v1)`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`create index Idx_test on exist_db.exist_tb_1 (v1)`,
		newTestResult().addResult(rule.Name, "Idx_test"))
}

func TestCheckFieldNotNUllMustContainDefaultValue(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue].Rule

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`create table exist_db.not_exist_tb_1(
			id int auto_increment not null,
			v1 varchar(255) not null default "unit test",
			v2 varchar(255) not null,
			V3 varchar(255) not null,
			primary key (id)
    )`, newTestResult().addResult(rule.Name, strings.Join([]string{"v2", "V3"}, ",")))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add column v4 int`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 add column v4 int not null `,
		newTestResult().addResult(rule.Name, strings.Join([]string{"v4"}, ",")))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 change v1 v1 int not null default 1`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 change column v1 v4 int not null`,
		newTestResult().addResult(rule.Name, strings.Join([]string{"v4"}, ",")))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 modify v1 int not null default 0`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`alter table exist_db.exist_tb_1 modify v1 int not null`,
		newTestResult().addResult(rule.Name, strings.Join([]string{"v1"}, ",")))
}

func TestCheckPrimaryKey(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: primary key exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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

func TestCheckDDLIndexTooMany(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexTooMany].Rule
	runSingleRuleInspectCase(rule, t, "create_table: index <= 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "create_table: index > 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,id),
INDEX idx_2 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexTooMany, "id", 2),
	)
}

func TestCheckDDLRedundantIndex(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckRedundantIndex].Rule
	runSingleRuleInspectCase(rule, t, "create_table: not redundant index", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "create_table: has repeat index", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,id),
INDEX idx_2 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "存在重复索引:(id); "),
	)

	runSingleRuleInspectCase(rule, t, "create_table: has redundant index", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id,v1),
INDEX idx_1 (id,v1,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "已存在索引 idx_1(id,v1,v2) , 索引 (id,v1) 为冗余索引; "),
	)

	runSingleRuleInspectCase(rule, t, "create_table: has repeat index 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id,v1),
INDEX idx_1 (id,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "存在重复索引:(id,v1); "),
	)

	runSingleRuleInspectCase(rule, t, "create_table: has repeat and redundant index", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (id,v1),
INDEX idx_2 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "存在重复索引:(id); 已存在索引 idx_1(id,v1) , 索引 idx_2(id) 为冗余索引; "),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: has repeat and redundant index", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_1 add index idx_t (v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "存在重复索引:idx_t(v1); "),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: has repeat and redundant index", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_6 add index idx_t (v2);
`,
		newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "已存在索引 idx_100(v2,v1) , 索引 idx_t(v2) 为冗余索引; "),
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

func TestCheckTableWithoutInnodb(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: table engine is innodb 1", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "create_table: table engine is innodb 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=Innodb AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "create_table: table engine is innodb 3", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=INNODB AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table engine is innodb 4", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists myisam_utf8_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=Innodb AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table engine not innodb 1", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableDBEngine, "Innodb"),
	)
	runDefaultRulesInspectCase(t, "create_table: table engine not innodb 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists myisam_utf8_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
) AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableDBEngine, "Innodb"),
	)
}

func TestCheckDatabaseWithoutUtf8mb4(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_database: character is utf8mb4", DefaultMysqlInspect(),
		`
CREATE DATABASE not_exist_db CHARACTER SET utf8mb4;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_database: character not utf8mb4", DefaultMysqlInspect(),
		`
CREATE DATABASE not_exist_db CHARACTER SET utf8;
`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)

	runDefaultRulesInspectCase(t, "alter_database: character not utf8mb4", DefaultMysqlInspect(),
		`
CREATE DATABASE not_exist_db CHARACTER SET utf8;
`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)
}

func TestCheckTableWithoutUtf8mb4(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: table charset is utf8mb4 1", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset is utf8mb4 2", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)CHARSET=utf8mb4 AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "create_table:table charset is utf8mb4 3", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)CHARSET=UTF8MB4 AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "create_table:table charset is utf8mb4 4", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists myisam_utf8_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB CHARSET=utf8mb4 AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset not utf8mb4 1", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1  COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset not utf8mb4 2", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists myisam_utf8_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)
	runDefaultRulesInspectCase(t, "create_table: column charset is utf8mb4", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
	v1 varchar(255) CHARACTER SET utf8mb4 NOT NULL DEFAULT "unit test" COMMENT "unit test",
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4  COMMENT="unit test";
	`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: column charset not utf8mb4 1", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
	v1 varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT "unit test" COMMENT "unit test",
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4  COMMENT="unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)

	runDefaultRulesInspectCase(t, "create_table: column charset has not utf8mb4 2", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
	v1 varchar(255) CHARACTER SET utf8mb4 NOT NULL DEFAULT "unit test" COMMENT "unit test",
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
	v2 varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT "unit test" COMMENT "unit test"
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4  COMMENT="unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)

	runDefaultRulesInspectCase(t, "alter_table: column charset has not utf8mb4 1", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 ADD column v3 varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT "unit test" COMMENT "unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckTableCharacterSet, "utf8mb4"),
	)

	runDefaultRulesInspectCase(t, "alter_table: column charset has not utf8mb4 3", DefaultMysqlInspect(),
		`
	ALTER TABLE exist_db.exist_tb_1 MODIFY column v2 varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT "unit test" COMMENT "unit test";
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_b1 (b1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
	)

	runDefaultRulesInspectCase(t, "create_table: disable index column blob (2)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob UNIQUE KEY COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
		newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob),
	)
}

func TestDisableForeignKey(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: has foreign key", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnTimestampWithoutDefault).
			addResult(rulepkg.DDLDisableTypeTimestamp),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 timestamp NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnTimestampWithoutDefault).
			addResult(rulepkg.DDLDisableTypeTimestamp),
	)
}

func TestCheckColumnBlobNotNull(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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

func TestDMLCheckSelectLimit(t *testing.T) {
	runDefaultRulesInspectCase(t, "success 1", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1000;
`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "success 2", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1;
`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "success 3", DefaultMysqlInspectOffline(),
		`
select 1;
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "success 4", DefaultMysqlInspectOffline(),
		`
select sleep(1);
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "failed big 1", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1001;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000).add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "failed big 2", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 2, 1001;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000).add(driver.RuleLevelNotice, "使用LIMIT分页时,避免使用LIMIT M,N").add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)

	runDefaultRulesInspectCase(t, "failed nil", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000),
	)
}

func TestDMLCheckSelectLimit_FP(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]success", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1 limit ?;
`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"),
	)
	runDefaultRulesInspectCase(t, "[fp]failed", DefaultMysqlInspect(),
		`
select id from exist_db.exist_tb_1 where id =1;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000),
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
	// default 5000,  unit testing :4
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
	// default 5000, unit testing :4
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

func Test_DMLCheckSelectWithOrderBy(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSelectWithOrderBy].Rule
	runSingleRuleInspectCase(rule, t, "",
		DefaultMysqlInspect(), `select id from exist_db.exist_tb_1 where v1 = '1' order by id`,
		newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))

	runSingleRuleInspectCase(rule, t, "",
		DefaultMysqlInspect(), `select id from exist_db.exist_tb_1 where v1 = '1' order by ?`,
		newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))

	runSingleRuleInspectCase(rule, t, "",
		DefaultMysqlInspect(), `select id from exist_db.exist_tb_1 where v1 = '1'`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select id from exist_db.exist_tb_1 where (select id from COLLATIONS order by id limit 1) = 1`,
		newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select id from exist_db.exist_tb_1 where (select id from COLLATIONS order by ? limit 1) = 1`,
		newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select id from exist_db.exist_tb_1 where (select id from COLLATIONS limit 1) = 1`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select id
from (select * from exist_db.exist_tb_1 order by id limit 10) as test
where id = 1;`, newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select id
from (select * from exist_db.exist_tb_1 limit 10) as test
where id = 1;`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select (select count(*) from exist_db.exist_tb_1 order by id limit 10) as count 
				from exist_db.exist_tb_1 where id = 1;`,
		newTestResult().addResult(rulepkg.DMLCheckSelectWithOrderBy))
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

func Test_DDLCheckCreateTimeColumn(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTimeColumn].Rule
	param := rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String()
	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
create table table_10
(
    id          int primary key,
    CREATe_TIME timestamp not null default CURRENT_TIMESTAMP,
    name        varchar(255)
)
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
create table table_10
(
    id          int primary key,
    create_time timestamp
)
`, newTestResult().addResult(rulepkg.DDLCheckCreateTimeColumn, param))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
create table table_10
(
    id          int primary key,
    create_time timestamp default current_timestamp
)
`, newTestResult())
}

func Test_DDLCheckSubQueryNestNum(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSubQueryNestNum]
	param := rule.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`select (select count(*) from users) as a
from exist_db.exist_tb_1
where (select count(*) from exist_db.exist_tb_1 where (select id from exist_db.exist_tb_1 limit 1) = 1)`, newTestResult())

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`select (select count(*) from users) as a
from exist_db.exist_tb_1
where (select count(*) from exist_db.exist_tb_2) > 1
  and (select count(*)
       from exist_db.exist_tb_1
       where (select id
              from exist_db.exist_tb_1
              where (select count(*) from exist_db.exist_tb_2 where (select count(*) from exist_db.exist_tb_2) = 1) =
                    1) = 1) = 1`,
		newTestResult().addResult(rulepkg.DMLCheckSubQueryNestNum, param))

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`select (select count(*)
        from users
        where (select count(*)
               from exist_db.exist_tb_1
               where (select id
                      from exist_db.exist_tb_1
                      where (select count(*)
                             from exist_db.exist_tb_2
                             where (select count(*) from exist_db.exist_tb_2) = 1) = 1) = 1) = 1) as a
from exist_db.exist_tb_1`,
		newTestResult().addResult(rulepkg.DMLCheckSubQueryNestNum, param))

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`delete
from exist_db.exist_tb_1
where (select count(*)
       from exist_db.exist_tb_1
       where (select id
              from exist_db.exist_tb_1
              where (select count(*) from exist_db.exist_tb_2 where exist_tb_2.id = 1) =
                    1) = 1) = 1`,
		newTestResult())

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`delete
from exist_db.exist_tb_1
where (select count(*)
       from exist_db.exist_tb_1
       where (select id
              from exist_db.exist_tb_1
              where (select count(*) from exist_db.exist_tb_2 where (select count(*) from exist_db.exist_tb_2) = 1) =
                    1) = 1) = 1`,
		newTestResult().addResult(rulepkg.DMLCheckSubQueryNestNum, param))

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`update exist_db.exist_tb_1,exist_db.exist_tb_2
set exist_tb_1.v1 = exist_tb_2.v1
where (select count(*) from exist_db.exist_tb_2) > 1
  and (select count(*)
       from exist_db.exist_tb_1
       where exist_tb_1.id = 1
         and (select id from exist_db.exist_tb_1 limit 1) = 1) > 1`, newTestResult())

	runSingleRuleInspectCase(rule.Rule, t, "", DefaultMysqlInspect(),
		`update exist_db.exist_tb_1,exist_db.exist_tb_2
set exist_tb_1.v1 = exist_tb_2.v1
where (select count(*)
       from exist_db.exist_tb_1
       where (select id
              from exist_db.exist_tb_1
              where (select count(*) from exist_db.exist_tb_2 where (select count(*) from exist_db.exist_tb_2) = 1) =
                    1) = 1) = 1;`,
		newTestResult().addResult(rulepkg.DMLCheckSubQueryNestNum, param))
}

func Test_DDLCheckUpdateTimeColumn(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckUpdateTimeColumn].Rule
	param := rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).String()
	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
	create table table_10
	(
	  id          int primary key,
	  update_time timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	  name        varchar(255)
	)
	`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
create table table_10
(
    id          int primary key,
    update_time timestamp,
    create_time timestamp
)
`, newTestResult().addResult(rulepkg.DDLCheckUpdateTimeColumn, param))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`
	create table table_10
	(
	   id          int primary key,
	   update_time timestamp default current_timestamp on    update current_timestamp
	)
	`, newTestResult())
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
		`create table`:                     `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8_general_ci;`,
		`alter table`:                      `ALTER TABLE exist_db.exist_tb_1 COLLATE utf8_general_ci;`,
		`create database`:                  `CREATE DATABASE db COLLATE utf8_general_ci;`,
		`create table with column collate`: `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10) COLLATE utf8_general_ci) COLLATE utf8mb4_0900_ai_ci;`,
		`alter table with column collate`:  `ALTER TABLE exist_db.exist_tb_1 modify column c1 varchar(255) COLLATE utf8_general_ci;`,
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
		`create table`:                     `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8mb4_0900_ai_ci;`,
		`alter table`:                      `ALTER TABLE exist_db.exist_tb_1 COLLATE utf8mb4_0900_ai_ci;`,
		`create database`:                  `CREATE DATABASE db COLLATE utf8mb4_0900_ai_ci;`,
		`create database upper case`:       `CREATE DATABASE db COLLATE UTF8MB4_0900_AI_CI;`,
		`create table with column collate`: `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10) COLLATE utf8mb4_0900_ai_ci) COLLATE utf8mb4_0900_ai_ci;`,
		`alter table with column collate`:  `ALTER TABLE exist_db.exist_tb_1 modify column c1 varchar(255) COLLATE utf8mb4_0900_ai_ci;`,
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
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnSetNotice].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckColumnSetNotice))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 INT;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB;`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnSetNotice].Rule,
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
		`(2)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, md5('aaa'), md5('bbb'));`,
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
		`(2)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, SHA1('aaa'), SHA1('bbb'));`,
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
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
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckObjectNameUseCN].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckObjectNameUseCN))
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
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckObjectNameUseCN].Rule,
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

func Test_DMLCheckJoinFieldType(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckJoinFieldType].Rule

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select * from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id left join exist_tb_3 t3 
    				on t3.id = t2.id where exist_tb_2.v2 = 'v1'`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select * from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id left join exist_tb_3 t3 
    				on t3.v1 = t2.id where exist_tb_2.v2 = 'v1'`,
		newTestResult().addResult(rulepkg.DMLCheckJoinFieldType))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select * from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id left join exist_tb_3 t3 
    				on t3.v1 = t2.id left join exist_tb_4 t4 on t4.id = t3.id where exist_tb_2.v2 = 'v1'`,
		newTestResult().addResult(rulepkg.DMLCheckJoinFieldType))

	// 不检测on condition表名不存在的情况
	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`select * from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id left join exist_tb_3 t3 
    				on t3.v1 = id  where exist_tb_2.v2 = 'v1'`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`update exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id left join exist_tb_3 t3 on t2.id=t3.id
set t1.id = 1
where t2.id = 2;`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`update exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.v1 left join exist_tb_3 t3 on t2.id=t3.id
set t1.id = 1
where t2.id = 2;`, newTestResult().addResult(rulepkg.DMLCheckJoinFieldType))

	// 不检测on condition表名不存在的情况
	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`update exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = v1 left join exist_tb_3 t3 on t2.id=t3.id
set t1.id = 1
where t2.id = 2;`, newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`delete exist_tb_1 , exist_tb_2 , exist_tb_3  from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id 
					left join exist_tb_3 t3 on t3.id = t2.id where t2.v2 = 'v1'`,
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`delete exist_tb_1 , exist_tb_2 , exist_tb_3  from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.v2 
					left join exist_tb_3 t3 on t3.id = t2.id where t2.v2 = 'v1'`,
		newTestResult().addResult(rulepkg.DMLCheckJoinFieldType))

	// 不检测on condition表名不存在的情况
	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		`delete exist_tb_1 , exist_tb_2 , exist_tb_3  from exist_tb_1 t1 left join exist_tb_2 t2 on t1.id = t2.id 
					left join exist_tb_3 t3 on t3.id = id where t2.v2 = 'v1'`, newTestResult())
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

func Test_DMLCheckInQueryLimit(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckInQueryNumber].Rule
	paramValue := "5"
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, paramValue)

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"select * from exist_tb_1 where id in (1,2,3,4,5,6)",
		newTestResult().addResult(rulepkg.DMLCheckInQueryNumber, 6, paramValue))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"select * from exist_tb_1 where id in (1,2,3,4,5)", newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"delete from exist_tb_1 where id in (1,2,3,4,5,6,7,8)",
		newTestResult().addResult(rulepkg.DMLCheckInQueryNumber, 8, paramValue))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"update exist_tb_1 set v1 = 'v1_next' where id in (1,2,3,4,5,6,7)",
		newTestResult().addResult(rulepkg.DMLCheckInQueryNumber, 7, paramValue))
}

func TestCheckIndexOption(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule
	e, handler, err := executor.NewMockExecutor()
	assert.NoError(t, err)

	inspect1 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT( DISTINCT ( v1 ) ) / COUNT( * ) * 100 AS v1 FROM exist_tb_3")).
		WillReturnRows(sqlmock.NewRows([]string{"v1"}).
			AddRow("100.0000"))
	runSingleRuleInspectCase(rule, t, "", inspect1, "alter table exist_tb_3 add primary key (v1);", newTestResult())

	inspect2 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT( DISTINCT ( v1 ) ) / COUNT( * ) * 100 AS v1 FROM exist_tb_3")).
		WillReturnRows(sqlmock.NewRows([]string{"v1"}).
			AddRow("100.0000"))
	runSingleRuleInspectCase(rule, t, "", inspect2, "alter table exist_tb_3 add unique(v1);", newTestResult())

	inspect3 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT( DISTINCT ( v2 ) ) / COUNT( * ) * 100 AS v2 FROM exist_tb_3")).
		WillReturnRows(sqlmock.NewRows([]string{"v2"}).
			AddRow("30.0000"))
	runSingleRuleInspectCase(rule, t, "", inspect3, "alter table exist_tb_3 add index idx_c2(v2);",
		newTestResult().addResult(rulepkg.DDLCheckIndexOption, "v2", 70))

	inspect4 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT( DISTINCT ( v3 ) ) / COUNT( * ) * 100 AS v3 FROM exist_tb_3")).
		WillReturnRows(sqlmock.NewRows([]string{"v3"}).
			AddRow("70.0000"))
	runSingleRuleInspectCase(rule, t, "", inspect4, "alter table exist_tb_3 add fulltext(v3);", newTestResult())

	inspect5 := NewMockInspect(e)
	handler.ExpectQuery(regexp.QuoteMeta("SELECT COUNT( DISTINCT ( v1 ) ) / COUNT( * ) * 100 AS v1,COUNT( DISTINCT ( v2 ) ) / COUNT( * ) * 100 AS v2 FROM exist_tb_3")).
		WillReturnRows(sqlmock.NewRows([]string{"v1", "v2"}).
			AddRow("100.0000", "30.0000"))
	runSingleRuleInspectCase(rule, t, "", inspect5, "alter table exist_tb_3 add index idx_c1_c2(v1,v2);", newTestResult())

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

func Test_DDLDisableAlterFieldUseFirstAndAfter(t *testing.T) {
	for _, sql := range []string{
		`alter table exist_db.exist_tb_2 Add column id_next int`,
		`alter table exist_db.exist_tb_2 change column v1 v1_next varchar(10)`,
		`alter table exist_db.exist_tb_2 modify column v1 varchar(10)`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLDisableAlterFieldUseFirstAndAfter].Rule, t, "",
			DefaultMysqlInspect(), sql, newTestResult())
	}

	for _, sql := range []string{
		`alter table exist_db.exist_tb_2 Add column id_next int after id`,
		`alter table exist_db.exist_tb_2 Add column id_next int first`,
		`alter table exist_db.exist_tb_2 change column id id_next int first`,
		`alter table exist_db.exist_tb_2 change column id id_next int after v1`,
		`alter table exist_db.exist_tb_2 modify column id varchar(3) first`,
		`alter table exist_db.exist_tb_2 modify column id varchar(3) after v1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLDisableAlterFieldUseFirstAndAfter].Rule, t, "",
			DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLDisableAlterFieldUseFirstAndAfter))
	}
}

func Test_DDLCheckBigintInsteadOfDecimal(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckBigintInsteadOfDecimal].Rule

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"create table t1(id_next decimal(10,2),id int,total_money decimal)",
		newTestResult().addResult(rulepkg.DDLCheckBigintInsteadOfDecimal, "id_next,total_money"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"create table t1(total_money decimal,remain_money decimal,id int)",
		newTestResult().addResult(rulepkg.DDLCheckBigintInsteadOfDecimal, "total_money,remain_money"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), "create table t1(remain_money bigint)",
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"alter table exist_db.exist_tb_2 modify column total_money decimal",
		newTestResult().addResult(rulepkg.DDLCheckBigintInsteadOfDecimal, "total_money"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"alter table exist_db.exist_tb_2 modify column total_money bigint",
		newTestResult())

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"alter table exist_db.exist_tb_2 add column remain_money decimal",
		newTestResult().addResult(rulepkg.DDLCheckBigintInsteadOfDecimal, "remain_money"))

	runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(),
		"alter table exist_db.exist_tb_2 change column id old_money decimal",
		newTestResult().addResult(rulepkg.DDLCheckBigintInsteadOfDecimal, "old_money"))
}

func Test_PerfectParse(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereIsInvalid].Rule, t, "", DefaultMysqlInspect(), `
SELECT * FROM exist_db.exist_tb_1;
OPTIMIZE TABLE exist_db.exist_tb_1;
SELECT * FROM exist_db.exist_tb_2;
`, newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
		newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性"),
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLCheckCreateTrigger))
	}

	for _, sql := range []string{
		`CREATE my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATEmy_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE trigger_1 BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER my_trigger BEEEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性"))
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLCheckCreateFunction))
	}

	for _, sql := range []string{
		`create function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`create123 function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER='sqle_op'@'localhost' hello (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性"))
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
			newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性").
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
			newTestResult().add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性"))
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
// 		if DefaultTemplateRules[i].Name == DDLCheckObjectNameUseCN {
// 			continue
// 		}

// 		ptrRules = append(ptrRules, &DefaultTemplateRules[i])
// 	}

// 	i.rules = ptrRules
// 	inspectCase(t, desc, i, wl, sql, results...)
// }

func Test_LowerCaseTableNameOpen(t *testing.T) {
	getLowerCaseOpenInspect := func() *MysqlDriverImpl {
		inspect := DefaultMysqlInspect()
		inspect.Ctx = session.NewMockContextForTestLowerCaseTableNameOpen(nil)
		return inspect
	}
	// check use
	{
		runEmptyRuleInspectCase(t, "test lower case table name open 1-1", getLowerCaseOpenInspect(),
			`use not_exist_db;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "not_exist_db"))

		runEmptyRuleInspectCase(t, "test lower case table name open 1-2", getLowerCaseOpenInspect(),
			`use NOT_EXIST_DB;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "NOT_EXIST_DB"))

		runEmptyRuleInspectCase(t, "test lower case table name open 1-3", getLowerCaseOpenInspect(),
			`use EXIST_DB;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 1-4", getLowerCaseOpenInspect(),
			`use EXIST_db;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 1-5", getLowerCaseOpenInspect(),
			`use exist_db;`,
			newTestResult())
	}
	// check schema
	{
		runEmptyRuleInspectCase(t, "test lower case table name open 2-1", getLowerCaseOpenInspect(),
			`create database EXIST_DB;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "EXIST_DB"))

		runEmptyRuleInspectCase(t, "test lower case table name open 2-2", getLowerCaseOpenInspect(),
			`create database exist_db;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "exist_db"))

		runEmptyRuleInspectCase(t, "test lower case table name open 2-3", getLowerCaseOpenInspect(),
			`create database not_exist_db;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 2-4", getLowerCaseOpenInspect(),
			`create database NOT_EXIST_DB;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 2-5", getLowerCaseOpenInspect(),
			`create database NOT_EXIST_DB;
create database NOT_EXIST_DB;`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "NOT_EXIST_DB"))

		runEmptyRuleInspectCase(t, "test lower case table name open 2-6", getLowerCaseOpenInspect(),
			`create database NOT_EXIST_DB;
create database not_exist_db;`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "not_exist_db"))
	}
	// check table
	{
		runEmptyRuleInspectCase(t, "test lower case table name open 3-1", getLowerCaseOpenInspect(),
			`create table EXIST_DB.exist_tb_1 (id int);`,
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "EXIST_DB.exist_tb_1"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-2", getLowerCaseOpenInspect(),
			`create table exist_db.exist_tb_1 (id int);`,
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "exist_db.exist_tb_1"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-3", getLowerCaseOpenInspect(),
			`create table EXIST_DB.EXIST_TB_1 (id int);`,
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "EXIST_DB.EXIST_TB_1"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-4", getLowerCaseOpenInspect(),
			`create table EXIST_DB.EXIST_TB_2 (id int);
create table EXIST_DB.exist_tb_2 (id int);`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "EXIST_DB.exist_tb_2"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-5", getLowerCaseOpenInspect(),
			`create table EXIST_DB.exist_tb_2 (id int);
create table EXIST_DB.EXIST_TB_2 (id int);`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "EXIST_DB.EXIST_TB_2"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-6", getLowerCaseOpenInspect(),
			`alter table exist_db.EXIST_TB_1 add column v3 varchar(255) COMMENT "unit test";`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 3-7", getLowerCaseOpenInspect(),
			`alter table exist_db.EXIST_TB_1 rename AS exist_tb_2;
alter table exist_db.EXIST_TB_1 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db.EXIST_TB_1"))

		runEmptyRuleInspectCase(t, "test lower case table name open 3-8", getLowerCaseOpenInspect(),
			`alter table exist_db.EXIST_TB_1 rename AS exist_tb_2;
alter table exist_db.exist_tb_2 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 3-9", getLowerCaseOpenInspect(),
			`alter table exist_db.EXIST_TB_1 rename AS exist_tb_2;
alter table exist_db.EXIST_TB_2 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult())
	}

	// check dml
	{
		runEmptyRuleInspectCase(t, "test lower case table name open 4-1", getLowerCaseOpenInspect(),
			`select id from exist_db.EXIST_TB_2 where id = 1;`,
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db.EXIST_TB_2"))

		runEmptyRuleInspectCase(t, "test lower case table name open 4-2", getLowerCaseOpenInspect(),
			`select id from exist_db.exist_tb_2 where id = 1;`,
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db.exist_tb_2"))

		runEmptyRuleInspectCase(t, "test lower case table name open 4-3", getLowerCaseOpenInspect(),
			`select id from exist_db.EXIST_TB_1 where id = 1;`, newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name open 4-4", getLowerCaseOpenInspect(),
			`select id from exist_db.exist_tb_1 where id = 1;`, newTestResult())
	}
}

func Test_LowerCaseTableNameClose(t *testing.T) {
	getLowerCaseCloseInspect := func() *MysqlDriverImpl {
		inspect := DefaultMysqlInspect()
		inspect.Ctx = session.NewMockContextForTestLowerCaseTableNameClose(nil)
		return inspect
	}
	// check use
	{
		runEmptyRuleInspectCase(t, "test lower case table name close 1-1", getLowerCaseCloseInspect(),
			`use not_exist_db;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "not_exist_db"))

		runEmptyRuleInspectCase(t, "test lower case table name close 1-2", getLowerCaseCloseInspect(),
			`use NOT_EXIST_DB;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "NOT_EXIST_DB"))

		runEmptyRuleInspectCase(t, "test lower case table name close 1-3", getLowerCaseCloseInspect(),
			`use exist_db_1;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 1-4", getLowerCaseCloseInspect(),
			`use EXIST_DB_1;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "EXIST_DB_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 1-5", getLowerCaseCloseInspect(),
			`use exist_DB_1;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "exist_DB_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 1-6", getLowerCaseCloseInspect(),
			`use EXIST_DB_2;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 1-7", getLowerCaseCloseInspect(),
			`use exist_db_2;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaNotExistMessage, "exist_db_2"))
	}
	// check schema
	{
		runEmptyRuleInspectCase(t, "test lower case table name close 2-1", getLowerCaseCloseInspect(),
			`create database exist_db_1;`,
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "exist_db_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 2-2", getLowerCaseCloseInspect(),
			`create database EXIST_DB_1;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 2-3", getLowerCaseCloseInspect(),
			`create database exist_DB_1;`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 2-4", getLowerCaseCloseInspect(),
			`create database NOT_EXIST_DB;
create database not_exist_db;`,
			newTestResult(),
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 2-5", getLowerCaseCloseInspect(),
			`create database NOT_EXIST_DB;
create database NOT_EXIST_DB;`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				SchemaExistMessage, "NOT_EXIST_DB"))
	}
	// check table
	{
		runEmptyRuleInspectCase(t, "test lower case table name close 3-1", getLowerCaseCloseInspect(),
			`create table exist_db_1.exist_tb_1 (id int);`,
			newTestResult().add(driver.RuleLevelError,
				TableExistMessage, "exist_db_1.exist_tb_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 3-2", getLowerCaseCloseInspect(),
			`create table exist_db_1.EXIST_TB_1 (id int);`,
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 3-3", getLowerCaseCloseInspect(),
			`alter table exist_db_1.EXIST_TB_1 rename AS exist_tb_2;
`,
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db_1.EXIST_TB_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 3-4", getLowerCaseCloseInspect(),
			`alter table exist_db_1.exist_tb_1 rename AS exist_tb_2;
alter table exist_db_1.exist_tb_1 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db_1.exist_tb_1"))

		runEmptyRuleInspectCase(t, "test lower case table name close 3-5", getLowerCaseCloseInspect(),
			`alter table exist_db_1.exist_tb_1 rename AS exist_tb_2;
alter table exist_db_1.exist_tb_2 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult())

		runEmptyRuleInspectCase(t, "test lower case table name close 3-6", getLowerCaseCloseInspect(),
			`alter table exist_db_1.exist_tb_1 rename AS exist_tb_2;
alter table exist_db_1.EXIST_TB_2 add column v3 varchar(255) COMMENT "unit test";
`,
			newTestResult(),
			newTestResult().add(driver.RuleLevelError,
				TableNotExistMessage, "exist_db_1.EXIST_TB_2"))
	}
}

// for issue 208
func TestInspect_CheckColumn(t *testing.T) {
	runDefaultRulesInspectCase(t, "check column 1", DefaultMysqlInspect(),
		`
	alter table exist_db.exist_tb_1 change column v1 v11 varchar(255) DEFAULT "v11" COMMENT "uint test";
	`,
		newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runDefaultRulesInspectCase(t, "check column 2", DefaultMysqlInspect(),
		`
	alter table exist_db.exist_tb_1 drop column v1;
	`,
		newTestResult())

	runDefaultRulesInspectCase(t, "check column 3", DefaultMysqlInspect(),
		`
	alter table exist_db.exist_tb_1 change column V1 v11 varchar(255) DEFAULT "v11" COMMENT "uint test";
	`,
		newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runDefaultRulesInspectCase(t, "check column 4", DefaultMysqlInspect(),
		`
	alter table exist_db.exist_tb_1 drop column V1;
	`,
		newTestResult())

	runDefaultRulesInspectCase(t, "check column 5", DefaultMysqlInspect(),
		`
	delete from exist_db.exist_tb_1 where id in (1, 2);
	`,
		newTestResult())

	runDefaultRulesInspectCase(t, "check column 6", DefaultMysqlInspect(),
		`
	delete from exist_db.exist_tb_1 where ID in (1, 2);
	`,
		newTestResult())

	runDefaultRulesInspectCase(t, "check column 7", DefaultMysqlInspect(),
		`
	select id, v1 from exist_db.exist_tb_1 where id in (1, 2) limit 1;
	`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"))

	runDefaultRulesInspectCase(t, "check column 8", DefaultMysqlInspect(),
		`
	select ID, V1 from exist_db.exist_tb_1 where ID in (1, 2) limit 1;
	`,
		newTestResult().add(driver.RuleLevelNotice, "未使用 ORDER BY 的 LIMIT 查询"))

	runDefaultRulesInspectCase(t, "check column 9", DefaultMysqlInspect(),
		`
	UPDATE exist_db.exist_tb_1 SET v1 = 1 WHERE id = 1;
	`,
		newTestResult())

	runDefaultRulesInspectCase(t, "check column 10", DefaultMysqlInspect(),
		`
	UPDATE exist_db.exist_tb_1 SET V1 = 1 WHERE ID = 1;
	`,
		newTestResult())
}

func Test_DDLDisableTypeTimestamp(t *testing.T) {
	for _, sql := range []string{
		`create table workflow_step_templates
		(
		   id                     int unsigned auto_increment
		       primary key,
		   created_at             datetime default CURRENT_TIMESTAMP null,
		   deleted_at             timestamp                           null
		);`,
		`alter table exist_tb_1
		   add column test_create_time timestamp;`,
		`alter table exist_tb_1
    modify column test_create_time timestamp;`,
		`alter table exist_tb_1
    change column v2 test_create_time timestamp;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLDisableTypeTimestamp].Rule, t, "",
			DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLDisableTypeTimestamp))
	}

	for _, sql := range []string{
		`create table workflow_step_templates
		(
		   id                     int unsigned auto_increment
		       primary key,
		   created_at             datetime default CURRENT_TIMESTAMP null
		);`,
		`alter table exist_tb_1
		   add column test_create_time datetime;`,
		`alter table exist_tb_1
    modify column test_create_time date;`,
		`alter table exist_tb_1
    change column v2 test_create_time date;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLDisableTypeTimestamp].Rule, t, "",
			DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDMLCheckAlias(t *testing.T) {
	for _, sql := range []string{
		"select id as a , a from exist_tb_1 where a = 1",
		//TODO　"select id from exist_tb_1 as exist_tb_1 where id = 1",
		//TODO　"select id from exist_tb_1 join exist_tb_2 as exist_tb_1 on id = 1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckAlias].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckAlias, "a"))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckAlias].Rule, t, "success", DefaultMysqlInspect(),
		"select id as a from exist_tb_1 as a1 join exist_tb_2 as a2 on id = 1",
		newTestResult())
}

func TestDDLHintUpdateTableCharsetWillNotUpdateFieldCharset(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE exist_tb_1 CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci;",
		`alter table exist_tb_1 default character set 'utf8';`,
		`alter table exist_tb_1 default character set='utf8';`,
		`ALTER TABLE exist_tb_1 CHANGE v1 v3 BIGINT NOT NULL, default character set utf8`,
		`ALTER TABLE exist_tb_1 CHANGE v1 v3 BIGINT NOT NULL, character set utf8`,
		//TODO　`alter table exist_tb_1 default collate = utf8_unicode_ci;`,
		`ALTER TABLE exist_tb_1 CHARACTER SET 'utf8';`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintUpdateTableCharsetWillNotUpdateFieldCharset].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLHintUpdateTableCharsetWillNotUpdateFieldCharset))
	}

	for _, sql := range []string{
		`ALTER TABLE exist_tb_1 MODIFY v1 TEXT CHARACTER SET utf8`,
		`ALTER TABLE exist_tb_1 CHANGE v1 v1 TEXT CHARACTER SET utf8;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintUpdateTableCharsetWillNotUpdateFieldCharset].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDDLHintDropColumn(t *testing.T) {
	for _, sql := range []string{
		`alter table exist_tb_1 drop column v2;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropColumn].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLHintDropColumn))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropColumn].Rule, t, "success", DefaultMysqlInspect(),
		"alter table exist_tb_1 drop index idx_1",
		newTestResult())
}

func TestDDLHintDropPrimaryKey(t *testing.T) {
	for _, sql := range []string{
		`alter table exist_tb_1 drop primary key`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropPrimaryKey].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLHintDropPrimaryKey))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropPrimaryKey].Rule, t, "success", DefaultMysqlInspect(),
		"alter table exist_tb_1 drop index idx_1",
		newTestResult())
}

func TestDDLHintDropForeignKey(t *testing.T) {
	for _, sql := range []string{
		`alter table exist_tb_1 drop foreign key v1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropForeignKey].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLHintDropForeignKey))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLHintDropForeignKey].Rule, t, "success", DefaultMysqlInspect(),
		"alter table exist_tb_1 drop index idx_1",
		newTestResult())
}

func TestDMLNotRecommendNotWildcardLike(t *testing.T) {
	for _, sql := range []string{
		`select * from exist_tb_1 where id like "a";`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendNotWildcardLike].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendNotWildcardLike))
	}

	for _, sql := range []string{
		`select * from exist_tb_1 where id like "%a";`,
		`select * from exist_tb_1 where id like "a%";`,
		`select * from exist_tb_1 where id like "%a%";`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendNotWildcardLike].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLHintInNullOnlyFalse(t *testing.T) {
	for _, sql := range []string{
		`SELECT * FROM exist_tb_1 WHERE v1 IN (NULL)`,
		`SELECT * FROM exist_tb_1 WHERE v1 NOT IN (NULL)`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintInNullOnlyFalse].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintInNullOnlyFalse))
	}
	for _, sql := range []string{
		`SELECT * FROM exist_tb_1 WHERE v1 IN ("1","2")`,
		`SELECT * FROM exist_tb_1 WHERE v1 NOT IN ("1","2")`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintInNullOnlyFalse].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDMLNotRecommendIn(t *testing.T) {
	for _, sql := range []string{
		`SELECT * FROM exist_tb_1 WHERE v1 IN (NULL)`,
		`SELECT * FROM exist_tb_1 WHERE v1 NOT IN (NULL)`,
		`SELECT * FROM exist_tb_1 WHERE v1 IN ("a")`,
		`SELECT * FROM exist_tb_1 WHERE v1 NOT IN ("a")`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendIn].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendIn))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendIn].Rule, t, "success", DefaultMysqlInspect(),
		`SELECT * FROM exist_tb_1 WHERE v1 <> "a"`,
		newTestResult())
}

func TestDMLCheckSpacesAroundTheString(t *testing.T) {
	for _, sql := range []string{
		`select ' 1'`,
		`select '1 '`,
		`select ' 1 '`,
		`select * from exist_tb_1 where id = ' 1'`,
		`select * from exist_tb_1 where id = '1 '`,
		`select * from exist_tb_1 where id = ' 1 '`,
		`insert into exist_tb_1 values (' 1','1','1')`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSpacesAroundTheString].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckSpacesAroundTheString))
	}
	for _, sql := range []string{
		`select '1'`,
		`select * from exist_tb_1 where id = '1'`,
		`insert into exist_tb_1 values ('1','1','1')`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSpacesAroundTheString].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDDLCheckFullWidthQuotationMarks(t *testing.T) {
	for _, sql := range []string{
		`alter table exist_tb_1 add column a int comment '”a“'`,
		`create table t (id int comment '”aaa“')`,
		//TODO　`alter table exist_tb_1 add column a int comment '‘'`,
		//TODO　`create table t (id int comment '’')`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckFullWidthQuotationMarks].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckFullWidthQuotationMarks))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckFullWidthQuotationMarks].Rule, t, "success", DefaultMysqlInspect(),
		`select "1"`,
		newTestResult())
}

func TestDMLNotRecommendOrderByRand(t *testing.T) {
	for _, sql := range []string{
		`select id from exist_tb_1 where id < 1000 order by rand(1)`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendOrderByRand].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendOrderByRand))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendOrderByRand].Rule, t, "success", DefaultMysqlInspect(),
		"select id from exist_tb_1 where id < 1000 order by v1",
		newTestResult())
}

func TestDMLNotRecommendGroupByConstant(t *testing.T) {
	for _, sql := range []string{
		`select v1,v2 from exist_tb_1 group by 1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendGroupByConstant].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendGroupByConstant))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendGroupByConstant].Rule, t, "success", DefaultMysqlInspect(),
		"select v1,v2 from exist_tb_1 group by v1",
		newTestResult())
}

func TestDMLCheckSortDirection(t *testing.T) {
	for _, sql := range []string{
		`select id,v1,v2 from exist_tb_1 where v1='foo' order by id desc, v2 asc`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortDirection].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckSortDirection))
	}

	for _, sql := range []string{
		`select id,v1,v2 from exist_tb_1 where v1='foo' order by id asc, v2 asc`,
		`select id,v1,v2 from exist_tb_1 where v1='foo' order by id desc, v2 desc`,
		`select id,v1,v2 from exist_tb_1 where v1='foo' order by id , v2 `,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortDirection].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDMLHintGroupByRequiresConditions(t *testing.T) {
	for _, sql := range []string{
		`select v1,v2 from exist_tb_1 group by 1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintGroupByRequiresConditions].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintGroupByRequiresConditions))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintGroupByRequiresConditions].Rule, t, "success", DefaultMysqlInspect(),
		"select v1,v2 from exist_tb_1 group by 1 order by v1",
		newTestResult())
}

func TestDMLNotRecommendGroupByExpression(t *testing.T) {
	for _, sql := range []string{
		"SELECT v1 FROM exist_tb_1 order by v1 - v2;",
		//TODO　"SELECT v1 - v2 a FROM exist_tb_1 order by a;",
		//TODO　"SELECT v1 FROM exist_tb_1 order by from_unixtime(v1);",
		//TODO　"SELECT from_unixtime(v1) a FROM exist_tb_1 order by a;",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendGroupByExpression].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendGroupByExpression))
	}

	for _, sql := range []string{
		`SELECT exist_tb_1.col FROM exist_tb_1 ORDER BY v1`,
		"SELECT sum(col) AS col FROM exist_tb_1 ORDER BY v1",
		"SELECT exist_tb_2.v1 FROM exist_tb_2, exist_tb_1 WHERE exist_tb_1.id = exist_tb_2.id ORDER BY exist_tb_1.v1",
		"SELECT col FROM exist_tb_1 order by `timestamp`;",
		"select col from exist_tb_1 where cl = 1 order by APPLY_TIME",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendGroupByExpression].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLCheckSQLLength(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSQLLength].Rule
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "64")
	for _, sql := range []string{
		"select * from exist_tb_1 where id = 'aaaaaaaaaaaaaaaaaaaaaaaaaaa'", // len = 65
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckSQLLength))
	}

	for _, sql := range []string{
		"select * from exist_tb_1 where id = 'aaaaaaaaaaaaaaaaaaaaaaaaaa'", // len = 64
		"select * from exist_tb_1 where id = 'aaaaaaaaaaaaaaaaaaaaaaaaa'",  // len = 63
	} {
		runSingleRuleInspectCase(rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLNotRecommendHaving(t *testing.T) {
	for _, sql := range []string{
		"SELECT exist_tb_1.id,count(exist_tb_1.id) FROM exist_tb_2 where id = 'test' GROUP BY exist_tb_1.id HAVING exist_tb_2.id <> '1660' AND exist_tb_2.id <> '2' order by exist_tb_2.id",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendHaving].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendHaving))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendHaving].Rule, t, "success", DefaultMysqlInspect(),
		"SELECT exist_tb_1.id,count(exist_tb_1.id) FROM exist_tb_2 where id = 'test' GROUP BY exist_tb_1.id",
		newTestResult())
}

func TestDMLHintUseTruncateInsteadOfDelete(t *testing.T) {
	for _, sql := range []string{
		"delete from exist_tb_1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintUseTruncateInsteadOfDelete].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintUseTruncateInsteadOfDelete))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintUseTruncateInsteadOfDelete].Rule, t, "success", DefaultMysqlInspect(),
		"truncate exist_tb_1",
		newTestResult())
}

func TestDMLNotRecommendUpdatePK(t *testing.T) {
	for _, sql := range []string{
		"update exist_tb_1 set id = '1'",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendUpdatePK].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendUpdatePK))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendUpdatePK].Rule, t, "success", DefaultMysqlInspect(),
		"update exist_tb_1 set v1 = 'a'",
		newTestResult())
}

func TestDDLCheckColumnQuantity(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnQuantity].Rule
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "5")

	for _, sql := range []string{
		"create table t(c1 int,c2 int,c3 int,c4 int,c5 int,c6 int);",
	} {
		runSingleRuleInspectCase(rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckColumnQuantity))
	}

	runSingleRuleInspectCase(rule, t, "success", DefaultMysqlInspect(),
		"create table t(c1 int,c2 int,c3 int,c4 int,c5 int);",
		newTestResult())
}

func TestDDLRecommendTableColumnCharsetSame(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` int(11) DEFAULT NULL, `col` char(10) CHARACTER SET utf8 DEFAULT NULL)",
		//TODO　"alter table exist_tb_1 change v1 v1 char(10) CHARACTER SET utf8 DEFAULT NULL;",
		"CREATE TABLE `t` ( `id` int(11) DEFAULT NULL, `col` char(10) CHARACTER SET utf8 DEFAULT NULL) DEFAULT CHARSET=utf8mb4",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLRecommendTableColumnCharsetSame].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLRecommendTableColumnCharsetSame))
	}

	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` int(10) )",
		//TODO　"CREATE TABLE `t` ( `id` varchar(10) CHARACTER SET utf8 ) DEFAULT CHARSET=utf8",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLRecommendTableColumnCharsetSame].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDDLCheckColumnTypeInteger(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` int(1) );",
		"CREATE TABLE `t` ( `id` bigint(1) );",
		//TODO　"alter TABLE `exist_tb_1` add column `v3` bigint(1);",
		//TODO　"alter TABLE `exist_tb_1` add column `v3` int(1);",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnTypeInteger].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckColumnTypeInteger))
	}

	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` int(10));",
		"CREATE TABLE `t` ( `id` bigint(20));",
		"alter TABLE `exist_tb_1` add column `v3` bigint(20);",
		"alter TABLE `exist_tb_1` add column `v3` int(10);",
		//TODO　"CREATE TABLE `t` ( `id` int);",
		//TODO　"alter TABLE `t` add column `id` bigint;",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnTypeInteger].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDDLCheckVarcharSize(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` varchar(1025) );",
		"alter TABLE `exist_tb_1` add column `v3` varchar(1025);",
		"alter table `exist_tb_1` modify column `v3` varchar(1025);",
		"alter table `exist_tb_1` change column `v2` `v3` varchar(1025);",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckVarcharSize].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckVarcharSize, 1024))
	}

	for _, sql := range []string{
		"CREATE TABLE `t` ( `id` varchar(1024));",
		"alter TABLE `exist_tb_1` add column `v3` varchar(1024);",
		"alter table `exist_tb_1` modify column `v3` varchar(1024);",
		"alter table `exist_tb_1` change column `v2` `v3` varchar(1024);",
		"alter table `exist_tb_1` drop column `v2`;",
		"alter table `exist_tb_1` rename column `v2` to `v3`;",
		"alter table `exist_tb_1` alter column `v2` drop default;",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckVarcharSize].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLNotRecommendFuncInWhere(t *testing.T) {
	for _, sql := range []string{
		`select id from exist_tb_1 where substring(v1,1,3)='abc';`,
		`SELECT * FROM exist_tb_1 WHERE UNIX_TIMESTAMP(v1) BETWEEN UNIX_TIMESTAMP('2018-11-16 09:46:00 +0800 CST') AND UNIX_TIMESTAMP('2018-11-22 00:00:00 +0800 CST')`,
		//TODO　`select id from exist_tb_1 where id/2 = 100`,
		//TODO　`select id from exist_tb_1 where id/2 < 100`,
		`SELECT * FROM exist_tb_1 WHERE DATE '2020-01-01'`,
		`DELETE FROM exist_tb_1 WHERE DATE '2020-01-01'`,
		`UPDATE exist_tb_1 SET id = 1 WHERE DATE '2020-01-01'`,
		`SELECT * FROM exist_tb_1 WHERE TIME '10:01:01'`,
		`SELECT * FROM exist_tb_1 WHERE TIMESTAMP '1587181360'`,
		`select * from exist_tb_1 where id = "root" and date '2020-02-01'`,
		`select id from exist_tb_1 where 'abc'=substring(v1,1,3);`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendFuncInWhere].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendFuncInWhere))
	}

	for _, sql := range []string{
		`select id from exist_tb_1 where v1 = (select 1)`,
		`select id from exist_tb_1 where v1 = 1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendFuncInWhere].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLNotRecommendSysdate(t *testing.T) {
	for _, sql := range []string{
		"select sysdate();",
		"select SYSDATE();",
		"select SysDate();",
		"select sysdate() from exist_tb_1;",
		"select SYSDATE() from exist_tb_1;",
		"select SysDate() from exist_tb_1;",
		"select * from exist_tb_1 where id = sysdate()",
		"select * from exist_tb_1 where id = SYSDATE()",
		"select * from exist_tb_1 where id = SysDate()",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendSysdate].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendSysdate))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendSysdate].Rule, t, "success", DefaultMysqlInspect(),
		"select * from exist_tb_1 where id =1 ",
		newTestResult())
}

func TestDMLHintSumFuncTips(t *testing.T) {
	for _, sql := range []string{
		"select sum(1);",
		"select SUM(1);",
		"select Sum(1);",
		"select * from exist_tb_1 where id = sum(1)",
		"select * from exist_tb_1 where id = SUM(1)",
		"select * from exist_tb_1 where id = Sum(1)",
		"select sum(1) from exist_tb_1",
		"select SUM(1) from exist_tb_1",
		"select Sum(1) from exist_tb_1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintSumFuncTips].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintSumFuncTips))
	}

	for _, sql := range []string{
		"select id from exist_tb_1 where id =1 ",
		"SELECT IF(ISNULL(SUM(v1)), 0, SUM(v1)) FROM exist_tb_1",
		"SELECT * FROM exist_tb_1 where id = IF(ISNULL(SUM(v1)), 0, SUM(v1))",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintSumFuncTips].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}
}

func TestDDLCheckColumnQuantityInPK(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t ( a int, b int, c int, PRIMARY KEY(`a`,`b`,`c`));",
		//TODO　"alter TABLE `exist_tb_1` add primary key (`id`,`v1`,`v2`);",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnQuantityInPK].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckColumnQuantityInPK))
	}
	for _, sql := range []string{
		"CREATE TABLE t ( a int, b int, c int, PRIMARY KEY(`a`,`b`));",
		//TODO　"alter TABLE `exist_tb_1` add primary key (`id`,`v1`);",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnQuantityInPK].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLHintLimitMustBeCombinedWithOrderBy(t *testing.T) {
	for _, sql := range []string{
		"select v1,v2 from exist_tb_1 where id =1 limit 10",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintLimitMustBeCombinedWithOrderBy].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintLimitMustBeCombinedWithOrderBy))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintLimitMustBeCombinedWithOrderBy].Rule, t, "success", DefaultMysqlInspect(),
		"select v1,v2 from exist_tb_1 where id =1 order by id limit 10",
		newTestResult())
}

func TestDMLHintTruncateTips(t *testing.T) {
	for _, sql := range []string{
		"TRUNCATE TABLE exist_tb_1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintTruncateTips].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintTruncateTips))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintTruncateTips].Rule, t, "success", DefaultMysqlInspect(),
		"delete from exist_tb_1",
		newTestResult())
}

func TestDMLHintDeleteTips(t *testing.T) {
	for _, sql := range []string{
		`delete from exist_tb_1 where v1 = v2;`,
		`truncate table exist_tb;`,
		`drop table exist_tb_1;`,
		//TODO　`drop database exist_db;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintDeleteTips].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLHintDeleteTips))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLHintDeleteTips].Rule, t, "success", DefaultMysqlInspect(),
		"select * from exist_tb_1 where id =1",
		newTestResult())
}

func TestDMLCheckSQLInjectionFunc(t *testing.T) {
	for _, sql := range []string{
		`select benchmark(10, rand())`,
		`select sleep(1)`,
		`select get_lock('lock_name', 1)`,
		`select release_lock('lock_name')`,
		`select id from exist_tb_1 where id = benchmark(10, rand())`,
		`select id from exist_tb_1 where id = sleep(1)`,
		`select id from exist_tb_1 where id = get_lock('lock_name', 1)`,
		`select id from exist_tb_1 where id = release_lock('lock_name')`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSQLInjectionFunc].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckSQLInjectionFunc))
	}

	for _, sql := range []string{
		`select sum(1)`,
		`select 1`,
		`select id from exist_tb_1 where id = sum(1)`,
		`select id from exist_tb_1 where id = 1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSQLInjectionFunc].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLCheckNotEqualSymbol(t *testing.T) {
	for _, sql := range []string{
		"select * from exist_tb_1 where id != 1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckNotEqualSymbol].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckNotEqualSymbol))
	}

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckNotEqualSymbol].Rule, t, "success", DefaultMysqlInspect(),
		"select * from exist_tb_1 where id <> 1",
		newTestResult())
}

func TestDMLNotRecommendSubquery(t *testing.T) {
	for _, sql := range []string{
		"select id,v1,v2 from exist_tb_1 where v1 in(select id from exist_tb_1)",
		"SELECT id,v1,v2 from exist_tb_1 where v1 =(SELECT id FROM `exist_tb_1` limit 1)",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendSubquery].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLNotRecommendSubquery))
	}

	for _, sql := range []string{
		"SELECT id,v1,v2 from exist_tb_1 where v1 = 1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLNotRecommendSubquery].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDMLCheckSubqueryLimit(t *testing.T) {
	for _, sql := range []string{
		"select id,v1,v2 from exist_tb_1 where v1 in(select id from exist_tb_1 limit 1)",
		"SELECT id,v1,v2 from exist_tb_1 where v1 =(SELECT id FROM `exist_tb_1` limit 1)",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSubqueryLimit].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DMLCheckSubqueryLimit))
	}
	for _, sql := range []string{
		"select id,v1,v2 from exist_tb_1 where v1 in(select id from exist_tb_1)",
		"SELECT id,v1,v2 from exist_tb_1 where v1 =(SELECT id FROM `exist_tb_1`)",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSubqueryLimit].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDDLCheckAutoIncrement(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE `tb` ( `id` int(10)) AUTO_INCREMENT=1",
		"CREATE TABLE `tb` ( `id` int(10)) AUTO_INCREMENT=2",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckAutoIncrement].Rule, t, "", DefaultMysqlInspect(), sql, newTestResult().addResult(rulepkg.DDLCheckAutoIncrement))
	}

	for _, sql := range []string{
		"CREATE TABLE `test1` ( `id` int(10))",
		"CREATE TABLE `test1` ( `id` int(10)) auto_increment = 0",
		"CREATE TABLE `test1` ( `id` int(10)) auto_increment = 0 DEFAULT CHARSET=latin1",
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckAutoIncrement].Rule, t, "success", DefaultMysqlInspect(), sql, newTestResult())
	}

}

func TestDDLNotAllowRenaming(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "success", DefaultMysqlInspect(), "ALTER TABLE exist_tb_1 MODIFY v1 CHAR(10);", newTestResult())

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "change 1", DefaultMysqlInspect(), "ALTER TABLE exist_tb_1 CHANGE v1 a BIGINT;", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "change 2", DefaultMysqlInspect(), "ALTER TABLE exist_tb_1 RENAME COLUMN v1 TO a", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "rename 1", DefaultMysqlInspect(), "RENAME TABLE exist_tb_1 TO test", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "rename 2", DefaultMysqlInspect(), "ALTER TABLE exist_tb_1 RENAME TO test", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

}
