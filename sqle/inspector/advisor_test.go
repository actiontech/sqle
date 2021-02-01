package inspector

import (
	"fmt"
	"testing"

	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"actiontech.cloud/universe/sqle/v4/sqle/model"

	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

func getTestCreateTableStmt1() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_1 (
id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "v1" COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1),
UNIQUE KEY uniq_1 (v1,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := parseOneSql("mysql", baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt2() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_2 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
user_id bigint unsigned NOT NULL COMMENT "unit test",
UNIQUE KEY uniq_1(id),
CONSTRAINT pk_test_1 FOREIGN KEY (user_id) REFERENCES exist_db.exist_tb_1 (id) ON DELETE NO ACTION
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := parseOneSql("mysql", baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt3() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_3 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL COMMENT "unit test",
v2 varchar(255) COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="uint test";
`
	node, err := parseOneSql("mysql", baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

type testResult struct {
	Results *InspectResults
	rules   map[string]RuleHandler
}

func newTestResult() *testResult {
	return &testResult{
		Results: newInspectResults(),
		rules:   RuleHandlerMap,
	}
}

func (t *testResult) add(level, message string, args ...interface{}) *testResult {
	t.Results.add(level, message, args...)
	return t
}

func (t *testResult) addResult(ruleName string, args ...interface{}) *testResult {
	handler, ok := t.rules[ruleName]
	if !ok {
		return t
	}
	level := handler.Rule.Level
	message := handler.Message
	if len(args) == 0 && handler.Rule.Value != "" {
		message = fmt.Sprintf(message, handler.Rule.Value)
	}
	return t.add(level, message, args...)
}

func (t *testResult) level() string {
	return t.Results.level()
}

func (t *testResult) message() string {
	return t.Results.message()
}

func DefaultMysqlInspect() *Inspect {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &Inspect{
		log:     log.NewEntry(),
		Results: newInspectResults(),
		Task: &model.Task{
			Instance: &model.Instance{
				Host:     "127.0.0.1",
				Port:     "3306",
				User:     "root",
				Password: "123456",
				DbType:   model.DB_TYPE_MYSQL,
			},
			CommitSqls:   []*model.CommitSql{},
			RollbackSqls: []*model.RollbackSql{},
		},
		SqlArray: []*model.Sql{},
		Ctx: &Context{
			currentSchema: "exist_db",
			schemaHasLoad: true,
			schemas: map[string]*SchemaInfo{
				"exist_db": &SchemaInfo{
					DefaultEngine:    "InnoDB",
					engineLoad:       true,
					DefaultCharacter: "utf8mb4",
					characterLoad:    true,
					Tables: map[string]*TableInfo{
						"exist_tb_1": &TableInfo{
							sizeLoad:      true,
							Size:          1,
							OriginalTable: getTestCreateTableStmt1(),
						},
						"exist_tb_2": &TableInfo{
							sizeLoad:      true,
							Size:          1,
							OriginalTable: getTestCreateTableStmt2(),
						},
						"exist_tb_3": &TableInfo{
							sizeLoad:      true,
							Size:          1,
							OriginalTable: getTestCreateTableStmt3(),
						},
					},
				},
			},
		},
		config: &Config{
			DDLOSCMinSize:      16,
			DMLRollbackMaxRows: 1000,
		},
	}
}

func runSingleRuleInspectCase(rule model.Rule, t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	inspectCase([]model.Rule{rule}, t, desc, i, sql, results...)
}

func runDefaultRulesInspectCase(t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	inspectCase(DefaultTemplateRules, t, desc, i, sql, results...)
}

func inspectCase(rules []model.Rule, t *testing.T, desc string, i *Inspect, sql string, results ...*testResult) {
	stmts, err := parseSql(i.Task.Instance.DbType, sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	for n, stmt := range stmts {
		i.Task.CommitSqls = append(i.Task.CommitSqls, &model.CommitSql{
			Sql: model.Sql{
				Number:  uint(n + 1),
				Content: stmt.Text(),
			},
		})
	}
	err = i.Advise(rules)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	if len(i.SqlArray) != len(results) {
		t.Errorf("%s test failled, error: result is unknow\n", desc)
		return
	}
	for n, sql := range i.Task.CommitSqls {
		result := results[n]
		if sql.InspectLevel != result.level() || sql.InspectResult != result.message() {
			t.Errorf("%s test failled, \n\nsql:\n %s\n\nexpect level: %s\nexpect result:\n%s\n\nactual level: %s\nactual result:\n%s\n",
				desc, sql.Content, result.level(), result.message(), sql.InspectLevel, sql.InspectResult)
		} else {
			t.Log(fmt.Sprintf("\n\ncase:%s\nactual level: %s\nactual result:\n%s\n\n", desc, sql.InspectLevel, sql.InspectResult))
		}
	}
}

func TestMessage(t *testing.T) {
	runDefaultRulesInspectCase(t, "check inspect message", DefaultMysqlInspect(),
		"use no_exist_db",
		&testResult{
			Results: &InspectResults{
				[]*InspectResult{&InspectResult{
					Level:   "error",
					Message: "schema no_exist_db 不存在",
				}},
			},
		},
	)
}

func TestCheckInvalidUse(t *testing.T) {
	runDefaultRulesInspectCase(t, "use_database: database not exist", DefaultMysqlInspect(),
		"use no_exist_db",
		newTestResult().add(model.RULE_LEVEL_ERROR,
			SCHEMA_NOT_EXIST_MSG, "no_exist_db"),
	)
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
		newTestResult().add(model.RULE_LEVEL_ERROR,
			SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
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
	handler := RuleHandlerMap[DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST]
	delete(RuleHandlerMap, DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST)
	defer func() {
		RuleHandlerMap[DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST] = handler
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
		newTestResult().add(model.RULE_LEVEL_ERROR,
			TABLE_EXIST_MSG, "exist_db.exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "create_table: refer table not exist", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_1 like exist_db.not_exist_tb_2;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR,
			TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb_2"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_MULTI_ERROR_MSG))

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
		newTestResult().add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_MULTI_ERROR_MSG))

	runDefaultRulesInspectCase(t, "create_table: duplicate column", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_ERROR_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_PRIMARY_KEY_COLUMN_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_1",
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
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "(匿名)",
			"v1").addResult(DDL_CHECK_INDEX_PREFIX))

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
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_1",
			"v1").add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_2", "v2"))
}

func TestCheckInvalidAlterTable(t *testing.T) {
	runDefaultRulesInspectCase(t, "alter_table: schema not exist", DefaultMysqlInspect(),
		`ALTER TABLE not_exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG,
			"not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "alter_table: table not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.not_exist_tb_1 add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
			"exist_db.not_exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add a exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add column v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_EXIST_MSG,
			"v1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: drop a not exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 drop column v5;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			"v5"),
	)

	runDefaultRulesInspectCase(t, "alter_table: alter a not exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 alter column v5 set default 'v5';
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
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
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			"v5"),
	)

	runDefaultRulesInspectCase(t, "alter_table: change column to a exist column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 change column v2 v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_EXIST_MSG,
			"v1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add pk ok", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 add primary key(id);
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "alter_table: add pk but exist pk", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add primary key(v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_EXIST_MSG),
	)

	runDefaultRulesInspectCase(t, "alter_table: add pk but key column not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 add primary key(id11);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			"id11"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add pk but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_2 add primary key(id,id);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_PRIMARY_KEY_COLUMN_MSG,
			"id"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add a exist index", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add index idx_1 (v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, INDEX_EXIST_MSG,
			"idx_1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: drop a not exist index", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 drop index idx_2;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, INDEX_NOT_EXIST_MSG,
			"idx_2"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add index but key column not exist", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add index idx_2 (v3);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			"v3"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add index but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add index idx_2 (id,id);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_2",
			"id"),
	)

	runDefaultRulesInspectCase(t, "alter_table: add index but key column is duplicate", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add index (id,id);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "(匿名)",
			"id").addResult(DDL_CHECK_INDEX_PREFIX),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_EXIST_MSG, "exist_db"),
	)
}

func TestCheckInvalidCreateIndex(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_index: schema not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON not_exist_db.not_exist_tb(v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "create_index: table not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON exist_db.not_exist_tb(v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "create_index: index exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, INDEX_EXIST_MSG, "idx_1"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column not exist", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(v3);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG, "v3"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column is duplicate", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(id,id);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_2", "id"),
	)

	runDefaultRulesInspectCase(t, "create_index: key column is duplicate", DefaultMysqlInspect(),
		`
CREATE INDEX idx_2 ON exist_db.exist_tb_1(id,id,v1);
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, "idx_2", "id"),
	)
}

func TestCheckInvalidDrop(t *testing.T) {
	handler := RuleHandlerMap[DDL_DISABLE_DROP_STATEMENT]
	delete(RuleHandlerMap, DDL_DISABLE_DROP_STATEMENT)
	defer func() {
		RuleHandlerMap[DDL_DISABLE_DROP_STATEMENT] = handler
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
		newTestResult().add(model.RULE_LEVEL_ERROR,
			SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR,
			SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "drop_table: table not exist", DefaultMysqlInspect(),
		`
DROP TABLE exist_db.not_exist_tb_1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR,
			TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb_1"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, INDEX_NOT_EXIST_MSG, "idx_2"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "insert: table not exist", DefaultMysqlInspect(),
		`
insert into exist_db.not_exist_tb values (1,"1","1");
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "insert: column not exist(1)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v3) values (1,"1","1");
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "v3"),
	)

	runDefaultRulesInspectCase(t, "insert: column not exist(2)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 set id=1,v1="1",v3="1";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "v3"),
	)

	runDefaultRulesInspectCase(t, "insert: column is duplicate(1)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v1) values (1,"1","1");
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG, "v1"),
	)

	runDefaultRulesInspectCase(t, "insert: column is duplicate(2)", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 set id=1,v1="1",v1="1";
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG, "v1"),
	)

	runDefaultRulesInspectCase(t, "insert: do not match values and columns", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2","2");
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, NOT_MATCH_VALUES_AND_COLUMNS),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "update: table not exist", DefaultMysqlInspect(),
		`
update exist_db.not_exist_tb set v1="2" where id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v3="2" where id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "v3"),
	)

	runDefaultRulesInspectCase(t, "update: where column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1="2" where v3=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "v3"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set t.v3 = "1" where t.id = 1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "t.v3"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set t.v1 = "1" where t.v3 = 1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "t.v3"),
	)

	runDefaultRulesInspectCase(t, "update with alias: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1 as t set exist_tb_1.v1 = "1" where t.id = 1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_1.v1"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v3 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_1.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_2.v3 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_2.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v1 = "1" where exist_tb_1.v3 = exist_tb_2.v3;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_1.v3,exist_tb_2.v3"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1,exist_db.exist_tb_2 set exist_tb_3.v1 = "1" where exist_tb_1.v1 = exist_tb_2.v1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_3.v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1,exist_db.exist_tb_2 set not_exist_db.exist_tb_1.v1 = "1" where exist_tb_1.v1 = exist_tb_2.v1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "not_exist_db.exist_tb_1.v1"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_IS_AMBIGUOUS, "v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set v1 = "1" where exist_tb_1.id = exist_tb_2.id;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_IS_AMBIGUOUS, "v1"),
	)

	runDefaultRulesInspectCase(t, "multi-update: where column not ambiguous", DefaultMysqlInspect(),
		`
update exist_tb_1,exist_tb_2 set exist_tb_1.v1 = "1" where v1 = 1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_IS_AMBIGUOUS, "v1"),
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
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "delete: table not exist", DefaultMysqlInspect(),
		`
delete from exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where v3=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "v3"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where exist_tb_1.v3=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_1.v3"),
	)

	runDefaultRulesInspectCase(t, "delete: where column not exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1 where exist_tb_2.id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG, "exist_tb_2.id"),
	)
}

func TestCheckInvalidSelect(t *testing.T) {
	runDefaultRulesInspectCase(t, "select: schema not exist", DefaultMysqlInspect(),
		`
select id from not_exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, "not_exist_db"),
	)

	runDefaultRulesInspectCase(t, "select: table not exist", DefaultMysqlInspect(),
		`
select id from exist_db.not_exist_tb where id=1;
`,
		newTestResult().add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "exist_db.not_exist_tb"),
	)
}

func TestCheckSelectAll(t *testing.T) {
	runDefaultRulesInspectCase(t, "select_from: all columns", DefaultMysqlInspect(),
		"select * from exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(DML_DISABE_SELECT_ALL_COLUMN),
	)
}

func TestCheckWhereInvalid(t *testing.T) {
	runDefaultRulesInspectCase(t, "select_from: has where condition", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id > 1;",
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(1)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(2)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(3)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where id=id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(4)", DefaultMysqlInspect(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID),
	)

	runDefaultRulesInspectCase(t, "update: has where condition", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "update: no where condition(1)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1';",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "update: no where condition(2)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where 1=1 and 2=2;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "update: no where condition(3)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where id=id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "update: no where condition(4)", DefaultMysqlInspect(),
		"update exist_db.exist_tb_1 set v1='v1' where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "delete: has where condition", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "delete: no where condition(1)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "delete: no where condition(2)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "delete: no where condition(3)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and id=id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))

	runDefaultRulesInspectCase(t, "delete: no where condition(4)", DefaultMysqlInspect(),
		"delete from exist_db.exist_tb_1 where 1=1 and exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DML_CHECK_WHERE_IS_INVALID))
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
		newTestResult().addResult(DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST),
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
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_USING_KEYWORD, "select, create, show").
			addResult(DDL_CHECK_INDEX_PREFIX),
	)

}

func TestAlterTableMerge(t *testing.T) {
	runDefaultRulesInspectCase(t, "alter_table: alter table need merge", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
ALTER TABLE exist_db.exist_tb_1 add column v6 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult(),
		newTestResult().addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE),
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
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
%s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
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
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "alter_table:add column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "alter_table:change column length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "alter_table: add index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD index idx_%s (v1);`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
	)

	runDefaultRulesInspectCase(t, "alter_table:rename index length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME index idx_1 TO idx_%s;`, length65),
		newTestResult().addResult(DDL_CHECK_OBJECT_NAME_LENGTH),
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
		newTestResult().addResult(DDL_CHECK_PK_NOT_EXIST),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL KEY DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT),
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
		newTestResult().addResult(DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(1)", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED),
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
		newTestResult().addResult(DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED),
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
		newTestResult().addResult(DDL_CHECK_COLUMN_CHAR_LENGTH),
	)
}

func TestCheckIndexCount(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: index <= 5", DefaultMysqlInspect(),
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

	runDefaultRulesInspectCase(t, "create_table: index > 5", DefaultMysqlInspect(),
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
		newTestResult().addResult(DDL_CHECK_INDEX_COUNT),
	)
}

func TestCheckCompositeIndexMax(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: composite index columns <= 3", DefaultMysqlInspect(),
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

	runDefaultRulesInspectCase(t, "create_table: composite index columns > 3", DefaultMysqlInspect(),
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
		newTestResult().addResult(DDL_CHECK_COMPOSITE_INDEX_MAX),
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
		newTestResult().addResult(DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset not utf8mb4", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1  COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4),
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
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
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
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
	)

	handler := RuleHandlerMap[DDL_CHECK_ALTER_TABLE_NEED_MERGE]
	delete(RuleHandlerMap, DDL_CHECK_ALTER_TABLE_NEED_MERGE)
	defer func() {
		RuleHandlerMap[DDL_CHECK_ALTER_TABLE_NEED_MERGE] = handler
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
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
		newTestResult().addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB),
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
		newTestResult().addResult(DDL_DISABLE_FK),
	)
}

func TestCheckTableComment(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: table without comment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDL_CHECK_TABLE_WITHOUT_COMMENT),
	)
}

func TestCheckColumnComment(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column without comment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without comment(1)", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL DEFAULT "unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without comment(2)", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v2 v3 varchar(255) NOT NULL DEFAULT "unit test" ;
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT),
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
		newTestResult().addResult(DDL_CHECK_INDEX_PREFIX),
	)

	runDefaultRulesInspectCase(t, "alter_table: index prefix not idx_", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD INDEX index_1(v1);
`,
		newTestResult().addResult(DDL_CHECK_INDEX_PREFIX),
	)

	runDefaultRulesInspectCase(t, "create_index: index prefix not idx_", DefaultMysqlInspect(),
		`
CREATE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(DDL_CHECK_INDEX_PREFIX),
	)
}

func TestCheckUniqueIndexPrefix(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
UNIQUE INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_UNIQUE_INDEX_PRIFIX),
	)

	runDefaultRulesInspectCase(t, "alter_table: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX index_1(v1);
`,
		newTestResult().addResult(DDL_CHECK_UNIQUE_INDEX_PRIFIX),
	)

	runDefaultRulesInspectCase(t, "create_index: unique index prefix not uniq_", DefaultMysqlInspect(),
		`
CREATE UNIQUE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(DDL_CHECK_UNIQUE_INDEX_PRIFIX),
	)
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
		newTestResult().addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT),
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
	handler := RuleHandlerMap[DDL_CHECK_COLUMN_WITHOUT_DEFAULT]
	delete(RuleHandlerMap, DDL_CHECK_COLUMN_WITHOUT_DEFAULT)
	defer func() {
		RuleHandlerMap[DDL_CHECK_COLUMN_WITHOUT_DEFAULT] = handler
	}()

	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 timestamp COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 timestamp NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT),
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
		newTestResult().addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL),
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
		newTestResult().addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL),
	)
}

func TestCheckDMLWithLimit(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(DML_CHECK_WITH_LIMIT),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(DML_CHECK_WITH_LIMIT),
	)
}

func TestCheckDMLWithOrderBy(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with order by", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(DML_CHECK_WITH_ORDER_BY),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspect(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(DML_CHECK_WITH_ORDER_BY),
	)
}

func TestCheckInsertColumnsExist(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_INSERT_COLUMNS_EXIST].Rule
	runSingleRuleInspectCase(rule, t, "insert: check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 values (1,"1","1"),(2,"2","2");
`,
		newTestResult().addResult(DML_CHECK_INSERT_COLUMNS_EXIST),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check columns exist", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2");
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMax(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_BATCH_INSERT_LISTS_MAX].Rule
	// defult 5000,  unit testing :4
	rule.Value = "4"
	runSingleRuleInspectCase(rule, t, "insert:check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4"),(5,"5","5");
`,
		newTestResult().addResult(DML_CHECK_BATCH_INSERT_LISTS_MAX, rule.Value),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check batch insert lists max", DefaultMysqlInspect(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4");
`,
		newTestResult(),
	)
}

func TestCheckPkProhibitAutoIncrement(t *testing.T) {
	rule := RuleHandlerMap[DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT].Rule
	runSingleRuleInspectCase(rule, t, "create_table: primary key not auto increment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT DEFAULT "unit test" COMMENT "unit test" ,
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT),
	)
	runSingleRuleInspectCase(rule, t, "create_table: passing the primary key not auto increment", DefaultMysqlInspect(),
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
}

func TestCheckWhereExistFunc(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_WHERE_EXIST_FUNC].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where nvl(v2,"0") = "3";
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_FUNC),
	)

	runSingleRuleInspectCase(rule, t, "select: passing the check where exist func", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 = "3"
`,
		newTestResult(),
	)
}

func TestCheckWhereExistNot(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_WHERE_EXIST_NOT].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist <> ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 <> "3";
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_NOT),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not like ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 not like "%3%";
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_NOT),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist != ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 != "3";
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_NOT),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not null ", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v2 is not null;
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_NOT),
	)
}

func TestCheckWhereExistImplicitConversion(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = 3;
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist implicit conversion", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 = "3";
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExist(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_LIMIT_MUST_EXIST].Rule
	runSingleRuleInspectCase(rule, t, "delete: check limit must exist", DefaultMysqlInspect(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(DML_CHECK_LIMIT_MUST_EXIST),
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
		newTestResult().addResult(DML_CHECK_LIMIT_MUST_EXIST),
	)
	runSingleRuleInspectCase(rule, t, "update: passing the check limit must exist", DefaultMysqlInspect(),
		`
update exist_db.exist_tb_1 set v1 ="1" limit 10 ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueries(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from  exist_db.exist_tb_2);
`,
		newTestResult().addResult(DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspect(),
		`
select a.v1 from exist_db.exist_tb_1 a, exist_db.exist_tb_2 b  where a.v1 = b.v1 ;
`,
		newTestResult(),
	)
}

func TestCheckIndexesExistBeforeCreatConstraints(t *testing.T) {
	rule := RuleHandlerMap[DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS].Rule
	runSingleRuleInspectCase(rule, t, "add unique: check indexes exist before creat constraints", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_3 add unique uniq_test(v2);
`, /*not exist index*/
		newTestResult().addResult(DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS),
	)
	runSingleRuleInspectCase(rule, t, "add unique: passing the check indexes exist before creat constraints", DefaultMysqlInspect(),
		`
alter table exist_db.exist_tb_1 add unique uniq_test(v1); 
`, /*exist index*/
		newTestResult(),
	)
}

func TestCheckSelectForUpdate(t *testing.T) {
	rule := RuleHandlerMap[DML_CHECK_SELECT_FOR_UPDATE].Rule
	runSingleRuleInspectCase(rule, t, "select : check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1 for update;
`,
		newTestResult().addResult(DML_CHECK_SELECT_FOR_UPDATE),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check exist select_for_update", DefaultMysqlInspect(),
		`
select v1 from exist_db.exist_tb_1; 
`,
		newTestResult(),
	)
}
func TestCheckCollationDatabase(t *testing.T) {
	rule := RuleHandlerMap[DDL_CHECK_COLLATION_DATABASE].Rule
	runSingleRuleInspectCase(rule, t, "create table: collation database", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8_general_ci;
`,
		newTestResult().addResult(DDL_CHECK_COLLATION_DATABASE, rule.Value),
	)
	runSingleRuleInspectCase(rule, t, "create table: passing the check collation database", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8mb4_0900_ai_ci; 
`,
		newTestResult(),
	)
}

func TestCheckDecimalTypeColumn(t *testing.T) {
	rule := RuleHandlerMap[DDL_CHECK_DECIMAL_TYPE_COLUMN].Rule
	runSingleRuleInspectCase(rule, t, "create table: check decimal type column", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 float(10));
`,
		newTestResult().addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN),
	)
	runSingleRuleInspectCase(rule, t, "alter table: check decimal type column", DefaultMysqlInspect(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 FLOAT ( 10, 0 );
`,
		newTestResult().addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN),
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
			RuleHandlerMap[DDL_CHECK_COLUMN_BLOB_NOTICE].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(DDL_CHECK_COLUMN_BLOB_NOTICE))
	}

	for desc, sql := range map[string]string{
		`(1)create table`: `CREATE TABLE t1(id INT);`,
		`(1)alter table`:  `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 VARCHAR(100);`,
		`(2)alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 SET('male', 'female');`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDL_CHECK_COLUMN_BLOB_NOTICE].Rule,
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
			RuleHandlerMap[DDL_CHECK_COLUMN_SET_NOTICE].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(DDL_CHECK_COLUMN_SET_NOTICE))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 INT;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB;`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDL_CHECK_COLUMN_SET_NOTICE].Rule,
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
			RuleHandlerMap[DDL_CHECK_COLUMN_ENUM_NOTICE].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(DDL_CHECK_COLUMN_ENUM_NOTICE))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 BLOB;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDL_CHECK_COLUMN_ENUM_NOTICE].Rule,
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
			RuleHandlerMap[DDL_CHECK_UNIQUE_INDEX].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(DDL_CHECK_UNIQUE_INDEX))
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
			RuleHandlerMap[DDL_CHECK_UNIQUE_INDEX].Rule,
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
			RuleHandlerMap[DML_CHECK_WHERE_EXIST_NULL].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult().addResult(DML_CHECK_WHERE_EXIST_NULL))
	}

	for desc, sql := range map[string]string{
		`select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id = 1;`,
		`update table`: `UPDATE exist_db.exist_tb_1 SET id = 10 WHERE id = 1;`,
		`delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id = 1;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DML_CHECK_WHERE_EXIST_NULL].Rule,
			t,
			desc,
			DefaultMysqlInspect(),
			sql,
			newTestResult())
	}
}

func DefaultMycatInspect() *Inspect {
	return &Inspect{
		log:     log.NewEntry(),
		Results: newInspectResults(),
		Task: &model.Task{
			Instance: &model.Instance{
				DbType: model.DB_TYPE_MYCAT,
				MycatConfig: &model.MycatConfig{
					AlgorithmSchemas: map[string]*model.AlgorithmSchema{
						"multidb": &model.AlgorithmSchema{
							AlgorithmTables: map[string]*model.AlgorithmTable{
								"exist_tb_1": &model.AlgorithmTable{
									ShardingColumn: "v1",
								},
								"exist_tb_2": &model.AlgorithmTable{
									ShardingColumn: "v1",
								},
							},
						},
					},
				},
			},
			CommitSqls:   []*model.CommitSql{},
			RollbackSqls: []*model.RollbackSql{},
		},
		SqlArray: []*model.Sql{},
		Ctx: &Context{
			currentSchema: "multidb",
			schemaHasLoad: true,
			schemas: map[string]*SchemaInfo{
				"multidb": &SchemaInfo{
					Tables: map[string]*TableInfo{
						"exist_tb_1": &TableInfo{
							sizeLoad:      true,
							Size:          1,
							OriginalTable: getTestCreateTableStmt1(),
						},
						"exist_tb_2": &TableInfo{
							sizeLoad:      true,
							Size:          1,
							OriginalTable: getTestCreateTableStmt2(),
						},
					},
				},
			},
		},
		config: &Config{
			DDLOSCMinSize:      16,
			DMLRollbackMaxRows: 1000,
		},
	}
}
