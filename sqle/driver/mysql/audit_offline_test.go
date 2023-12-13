package mysql

import (
	"fmt"
	"testing"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/sirupsen/logrus"
)

func DefaultMysqlInspectOffline() *MysqlDriverImpl {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &MysqlDriverImpl{
		log:  log.NewEntry(),
		inst: nil,
		Ctx:  &session.Context{},
		cnf: &Config{
			DDLOSCMinSize:      -1,
			DMLRollbackMaxRows: -1,
		},
		isOfflineAudit: true,
	}
}

func TestCheckSelectAllOffline(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLDisableSelectAllColumn].Rule, t,
		"select_from: all columns",
		DefaultMysqlInspectOffline(),
		"select * from exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(rulepkg.DMLDisableSelectAllColumn),
	)
}

func TestCheckInvalidDropOffline(t *testing.T) {
	handler := rulepkg.RuleHandlerMap[rulepkg.DDLDisableDropStatement]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLDisableDropStatement)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLDisableDropStatement] = handler
	}()
	runDefaultRulesInspectCase(t, "drop_database: ok", DefaultMysqlInspectOffline(),
		`
DROP DATABASE if exists exist_db;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_database: schema not exist(1)", DefaultMysqlInspectOffline(),
		`
DROP DATABASE if exists not_exist_db;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_table: ok", DefaultMysqlInspectOffline(),
		`
DROP TABLE exist_db.exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_table: schema not exist(1)", DefaultMysqlInspectOffline(),
		`
DROP TABLE if exists not_exist_db.not_exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_index: ok", DefaultMysqlInspectOffline(),
		`
DROP INDEX idx_1 ON exist_db.exist_tb_1;
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "drop_index: if exists and index not exist", DefaultMysqlInspectOffline(),
		`
DROP INDEX IF EXISTS idx_2 ON exist_db.exist_tb_1;
`,
		newTestResult(),
	)
}

func TestCheckWhereInvalidOffline(t *testing.T) {
	// results in this unit test
	noResult := newTestResult()
	whereIsInvalid := newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid)
	// the rule this unit test test
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereIsInvalid].Rule

	testCases := []struct {
		testName string
		sql      string
		result   *testResult
	}{
		// WHERE
		{
			"select_from: has where condition",
			"select id from exist_db.exist_tb_1 where id > 1;",
			noResult,
		},
		{
			"select_from: has where condition",
			"select id from exist_db.exist_tb_1;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(1)",
			"select id from exist_db.exist_tb_1;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(2)",
			"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(3)",
			"select id from exist_db.exist_tb_1 where id=id;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(4)",
			"select id from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(5)",
			"select id from (select * from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id) t;",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(6)",
			"select id from (select * from exist_db.exist_tb_1 where exist_tb_1.id>1) t;",
			whereIsInvalid,
		},
		// UPDATE
		{
			"update: has where condition",
			"update exist_db.exist_tb_1 set v1='v1' where id = 1;",
			noResult,
		},
		{
			"update: no where condition(1)",
			"update exist_db.exist_tb_1 set v1='v1';",
			whereIsInvalid,
		},
		{
			"update: no where condition(2)",
			"update exist_db.exist_tb_1 set v1='v1' where 1=1 and 2=2;",
			whereIsInvalid,
		},
		{
			"update: no where condition(3)",
			"update exist_db.exist_tb_1 set v1='v1' where id=id;",
			whereIsInvalid,
		},
		{
			"update: no where condition(4)",
			"update exist_db.exist_tb_1 set v1='v1' where exist_tb_1.id=exist_tb_1.id;", whereIsInvalid,
		},
		{
			"update: has where condition(5)",
			"update exist_db.exist_tb_1 set v1=v1 = v1 * (SELECT AVG(id) FROM exist_db.exist_tb_1 WHERE v1=1)/100 where id = 1;",
			noResult,
		},
		{
			"update: has where condition(6)",
			"update exist_db.exist_tb_1 set v1=v1 = v1 * (SELECT AVG(id) FROM exist_db.exist_tb_1 WHERE exist_tb_1.id=exist_tb_1.id)/100 where id = 1;",
			whereIsInvalid,
		},
		// DELETE
		{
			"delete: has where condition",
			"delete from exist_db.exist_tb_1 where id = 1;",
			noResult,
		},
		{
			"delete: no where condition(1)",
			"delete from exist_db.exist_tb_1;",
			whereIsInvalid,
		},
		{
			"delete: no where condition(2)",
			"delete from exist_db.exist_tb_1 where 1=1 and 2=2;",
			whereIsInvalid,
		},
		{
			"update: no where condition(3)",
			"delete from exist_db.exist_tb_1 where 1=1 and id=id;",
			whereIsInvalid,
		},
		{
			"update: no where condition(4)",
			"delete from exist_db.exist_tb_1 where 1=1 and exist_tb_1.id=exist_tb_1.id;", whereIsInvalid,
		},
		{
			"delete: has where condition(5)",
			"delete from exist_db.exist_tb_1 USING (SELECT * FROM exist_db.exist_tb_1 WHERE v1='v1') t WHERE t.id > 10;",
			noResult,
		},
		{
			"delete: has where condition(6)",
			"delete from exist_db.exist_tb_1 USING (SELECT * FROM exist_db.exist_tb_1 WHERE exist_tb_1.id=exist_tb_1.id) t WHERE t.id > 10;",
			whereIsInvalid,
		},
		// exists
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2 t2 where t1.id=t2.id);",
			noResult,
		},
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2 t2 where 1=1);",
			whereIsInvalid,
		},
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2);",
			whereIsInvalid,
		},
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2 t2 where exists (select 1 from exist_db.exist_db_3 t3 where t1.id=t2.id and t2.id=t3.id));",
			noResult,
		},
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2 t2 where exists (select 1 from exist_db.exist_db_3 t3 where 1=1));",
			whereIsInvalid,
		},
		{
			"use exists",
			"select * from exist_db.exist_tb_1 t1 where exists (select 1 from exist_db.exist_tb_2 t2 where exists (select 1 from exist_db.exist_db_3));",
			whereIsInvalid,
		},
		{
			"use not exists",
			"select * from exist_db.exist_tb_1 t1 where not exists (select 1 from exist_db.exist_tb_2 t2 where t1.id=t2.id);",
			noResult,
		},
		{
			"use not exists",
			"select * from exist_db.exist_tb_1 t1 where not exists (select 1 from exist_db.exist_tb_2 t2 where 1=1);",
			whereIsInvalid,
		},
		{
			"use not exists",
			"select * from exist_db.exist_tb_1 t1 where not exists (select 1 from exist_db.exist_tb_2);",
			whereIsInvalid,
		},
		{
			"select_from: no where condition(2)",
			"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
			whereIsInvalid,
		},
		// value compare
		// int
		{
			"int compare(1)",
			"select * from exist_db.exist_tb_1 where 1 > 0",
			whereIsInvalid,
		},
		{
			"int compare(2)",
			"select * from exist_db.exist_tb_1 where 1 < 0",
			noResult,
		},
		{
			"int compare(3)",
			"select * from exist_db.exist_tb_1 where 1 >= 0",
			whereIsInvalid,
		},
		{
			"int compare(4)",
			"select * from exist_db.exist_tb_1 where 1 = 0",
			noResult,
		},
		{
			"int compare(5)",
			"select * from exist_db.exist_tb_1 where 1 <= 0",
			noResult,
		},
		{
			"int compare(6)",
			"select * from exist_db.exist_tb_1 where 1 != 0",
			whereIsInvalid,
		},
		{
			"int compare(7)",
			"select * from exist_db.exist_tb_1 where 1 = '1'",
			noResult,
		},
		// str
		{
			"str compare(1)",
			"select * from exist_db.exist_tb_1 where '1' = '1'",
			whereIsInvalid,
		},
		{
			"str compare(2)",
			"select * from exist_db.exist_tb_1 where '1' > '1'",
			noResult,
		},
		{
			"str compare(3)",
			"select * from exist_db.exist_tb_1 where '1' > '0'",
			whereIsInvalid,
		},
		{
			"str compare(3)",
			"select * from exist_db.exist_tb_1 where '1' != '0'",
			whereIsInvalid,
		},
		{
			"str compare(4)",
			"select * from exist_db.exist_tb_1 where '1' >= '1'",
			whereIsInvalid,
		},
		{
			"str compare(5)",
			"select * from exist_db.exist_tb_1 where '1' < '1'",
			noResult,
		},
		// float
		{
			"float compare(1)",
			"select * from exist_db.exist_tb_1 where 1.6 = 1.6",
			whereIsInvalid,
		},
		{
			"float compare(2)",
			"select * from exist_db.exist_tb_1 where 1.6 > 1.2",
			whereIsInvalid,
		},
		{
			"float compare(3)",
			"select * from exist_db.exist_tb_1 where 1.6 < 1.2",
			noResult,
		},
		{
			"float compare(4)",
			"select * from exist_db.exist_tb_1 where 1.6 >= 1.2",
			whereIsInvalid,
		},
	}
	offlineInspect := DefaultMysqlInspectOffline()
	for _, testCase := range testCases {
		runSingleRuleInspectCase(rule, t, testCase.testName, offlineInspect, testCase.sql, testCase.result)
	}
}

func TestCheckWhereInvalid_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereIsInvalid].Rule

	runSingleRuleInspectCase(rule, t, "[pf]select_from: has where condition(1)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where id=?;",
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[pf]select_from: has where condition(2)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=?;",
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[pf]select_from: no where condition(1)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where 1=? and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)
	runSingleRuleInspectCase(rule, t, "[pf]select_from: no where condition(2)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where ?=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
	)

	runSingleRuleInspectCase(rule, t, "[pf]update: has where condition", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where id = ?;",
		newTestResult())

	runSingleRuleInspectCase(rule, t, "[pf]update: no where condition(1)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runSingleRuleInspectCase(rule, t, "[pf]update: no where condition(2)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where 1=1 and 2=2;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runSingleRuleInspectCase(rule, t, "[pf]update: no where condition(3)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runSingleRuleInspectCase(rule, t, "[pf]update: no where condition(4)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runSingleRuleInspectCase(rule, t, "[pf]delete: no where condition(1)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=? and ?=?;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))

	runSingleRuleInspectCase(rule, t, "[pf]delete: no where condition(2)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=? and id=id;",
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))
}

func TestCheckCreateTableWithoutIfNotExistsOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: need \"if not exists\"", DefaultMysqlInspectOffline(),
		`
CREATE TABLE exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT= "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckPKWithoutIfNotExists),
	)
}

func TestCheckObjectNameUsingKeywordOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: using keyword", DefaultMysqlInspectOffline(),
		"CREATE TABLE if not exists exist_db.`select` ("+
			"id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT \"unit test\","+
			"v1 varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"`create` varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT \"unit test\","+
			"update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT \"unit test\","+
			"PRIMARY KEY (id),"+
			"INDEX `show` (v1)"+
			")ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT=\"unit test\";",
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckObjectNameUsingKeyword, "select, create, show").
			addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)
}

func TestCheckObjectNameLengthOffline(t *testing.T) {
	length64 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcd"
	length65 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcde"

	runDefaultRulesInspectCase(t, "create_table: table length <= 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length64),
		newTestResult().addResult(rulepkg.DDLCheckPKName),
	)

	runDefaultRulesInspectCase(t, "create_table: table length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
%s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "create_table: index length > 64", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64).addResult(rulepkg.DDLNotAllowRenaming),
	)

	runDefaultRulesInspectCase(t, "alter_table:Add column length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:change column length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64).addResult(rulepkg.DDLNotAllowRenaming),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD index idx_%s (v1);`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)

	runDefaultRulesInspectCase(t, "alter_table:rename index length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME index idx_1 TO idx_%s;`, length65),
		newTestResult().addResult(rulepkg.DDLCheckObjectNameLength, 64),
	)
}

func TestCheckPrimaryKeyOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: primary key exist", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: primary key not exist", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKNotExist).addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "id"),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(1)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL KEY DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(2)", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(1)", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(2)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned),
	)
}

func TestCheckColumnCharLengthOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: check char(20)", DefaultMysqlInspectOffline(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
	update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",v1 char(20) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult().addResult(rulepkg.DDLCheckPKName),
	)

	runDefaultRulesInspectCase(t, "create_table: check char(21)", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckColumnCharLength),
	)
}

func TestCheckIndexCountOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexCount].Rule
	runSingleRuleInspectCase(rule, t, "create_table: index <= 5", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: index > 5", DefaultMysqlInspectOffline(),
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

func TestCheckCompositeIndexMaxOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckCompositeIndexMax].Rule
	runSingleRuleInspectCase(rule, t, "create_table: composite index columns <= 3", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: composite index columns > 3", DefaultMysqlInspectOffline(),
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

func TestCheckTableWithoutInnodbUtf8mb4Offline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: ok", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: table engine not innodb", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table charset not utf8mb4", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1  COMMENT="unit test";
`,
		newTestResult(),
	)
}

func TestCheckIndexColumnWithBlobOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: disable index column blob (1)", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckIndexedColumnWithBlob).add(driverV2.RuleLevelWarn, rulepkg.DDLCheckIndexNotNullConstraint, "这些索引字段(b1)需要有非空约束"),
	)

	runDefaultRulesInspectCase(t, "create_table: disable index column blob (2)", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckIndexedColumnWithBlob).add(driverV2.RuleLevelWarn, rulepkg.DDLCheckIndexNotNullConstraint, "这些索引字段(b1)需要有非空约束"),
	)

	handler := rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckAlterTableNeedMerge)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge] = handler
	}()

}

func TestDisableForeignKeyOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: has foreign key", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLDisableFK),
	)
}

func TestCheckTableCommentOffline(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableWithoutComment].Rule, t, "create_table: table without comment", DefaultMysqlInspectOffline(),
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

func TestCheckColumnCommentOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutComment].Rule
	runSingleRuleInspectCase(rule, t, "create_table: column without comment", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "alter_table: column without comment(1)", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL DEFAULT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutComment),
	)

	runSingleRuleInspectCase(rule, t, "alter_table: column without comment(2)", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v2 v3 varchar(255) NOT NULL DEFAULT "unit test" ;
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutComment),
	)
}

func TestCheckDDLRedundantIndexOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckRedundantIndex].Rule
	runSingleRuleInspectCase(rule, t, "create_table: not redundant index", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: has repeat index", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: has redundant index", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: has repeat index 2", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "create_table: has repeat and redundant index", DefaultMysqlInspectOffline(),
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

}

func TestCheckIndexPrefixOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: index prefix not idx_", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

	runDefaultRulesInspectCase(t, "alter_table: index prefix not idx_", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD INDEX index_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckIndexPrefix, "idx_"),
	)

	runDefaultRulesInspectCase(t, "create_index: index prefix not idx_", DefaultMysqlInspectOffline(),
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexPrefix].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckUniqueIndexPrefixOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndexPrefix].Rule

	runSingleRuleInspectCase(rule, t, "create_table: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
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

	runSingleRuleInspectCase(rule, t, "alter_table: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX index_1(v1);
`,
		newTestResult().addResult(rulepkg.DDLCheckUniqueIndexPrefix, "uniq_"),
	)

	runSingleRuleInspectCase(rule, t, "create_index: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndexPrefix].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckColumnDefaultOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v1 varchar(255) COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckColumnWithoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnWithoutDefault).
			addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "v3"),
	)

	runDefaultRulesInspectCase(t, "alter_table: auto increment column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "v3"),
	)

	runDefaultRulesInspectCase(t, "alter_table: blob column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob COMMENT "unit test";
`,
		newTestResult(),
	)
}

func TestCheckColumnTimestampDefaultOffline(t *testing.T) {
	handler := rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutDefault]
	delete(rulepkg.RuleHandlerMap, rulepkg.DDLCheckColumnWithoutDefault)
	defer func() {
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckColumnWithoutDefault] = handler
	}()

	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v1 timestamp COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckColumnTimestampWithoutDefault).addResult(rulepkg.DDLDisableTypeTimestamp),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 timestamp NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnTimestampWithoutDefault).
			addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "v3").
			addResult(rulepkg.DDLDisableTypeTimestamp),
	)
}

func TestCheckColumnBlobNotNullOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
v1 blob NOT NULL COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckColumnBlobWithNotNull).
			addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "v1"),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobWithNotNull).addResult(rulepkg.DDLCheckFieldNotNUllMustContainDefaultValue, "v3"),
	)
}

func TestCheckColumnBlobDefaultNullOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 blob DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckPKName).addResult(rulepkg.DDLCheckColumnBlobDefaultIsNotNull),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().addResult(rulepkg.DDLCheckColumnBlobDefaultIsNotNull),
	)
}

func TestCheckDMLWithLimitOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)
}

func TestCheckDMLWithLimit_FPOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]update: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "[fp]delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithLimit),
	)
}

func TestDMLCheckSelectLimitOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "success 1", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1000;
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "success 2", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1;
`,
		newTestResult(),
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

	runDefaultRulesInspectCase(t, "failed big 1", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 1001;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000),
	)

	runDefaultRulesInspectCase(t, "failed big 2", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1 limit 2, 1001;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000).
			add(driverV2.RuleLevelNotice, "", "使用分页查询时，避免使用偏移量"),
	)

	runDefaultRulesInspectCase(t, "failed nil", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000),
	)
}

func TestDMLCheckSelectLimit_FPOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]success", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1 limit ?;
`,
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "[fp]failed", DefaultMysqlInspectOffline(),
		`
select id from exist_db.exist_tb_1 where id =1;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectLimit, 1000),
	)

}

func TestCheckDMLWithOrderByOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWithOrderBy].Rule
	runSingleRuleInspectCase(rule, t, "update: with order by", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)

	runSingleRuleInspectCase(rule, t, "delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)
}

func TestCheckDMLWithOrderBy_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWithOrderBy].Rule
	runSingleRuleInspectCase(rule, t, "[fp]update: with order by", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)

	runSingleRuleInspectCase(rule, t, "[fp]delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1=? where id=? order by ?;
`,
		newTestResult().addResult(rulepkg.DMLCheckWithOrderBy),
	)
}

func TestCheckInsertColumnsExistOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "insert: check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 values (1,"1","1"),(2,"2","2");
`,
		newTestResult().addResult(rulepkg.DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2");
`,
		newTestResult(),
	)
}

func TestCheckInsertColumnsExist_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]insert: check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 values (?,?,?),(?,?,?);
`,
		newTestResult().addResult(rulepkg.DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMaxOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckBatchInsertListsMax].Rule
	// default 5000,  unit testing :4
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "4")
	runSingleRuleInspectCase(rule, t, "insert:check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4"),(5,"5","5");
`,
		newTestResult().addResult(rulepkg.DMLCheckBatchInsertListsMax, 4),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4");
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMax_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckBatchInsertListsMax].Rule
	// default 5000, unit testing :4
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "4")
	runSingleRuleInspectCase(rule, t, "[fp]insert:check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult().addResult(rulepkg.DMLCheckBatchInsertListsMax, 4),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckPkProhibitAutoIncrementOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKProhibitAutoIncrement].Rule
	runSingleRuleInspectCase(rule, t, "create_table: primary key not auto increment", DefaultMysqlInspectOffline(),
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

}

func TestCheckWhereExistFuncOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistFunc].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist func", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where nvl(v2,"0") = "3";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "select: passing the check where exist func", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 = "3"
`,
		newTestResult(),
	)
}

func TestCheckWhereExistFunc_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistFunc].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: check where exist func", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where nvl(v2,?) = ?;
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check where exist func", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 = ?
`,
		newTestResult(),
	)
}

func TestCheckWhereExistNotOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistNot].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist <> ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 <> "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist <> ", DefaultMysqlInspectOffline(),
		`
		select v1 from (select * from exist_db.exist_tb_1 where v2 <> "3") t;
		`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not like ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 not like "%3%";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist != ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 != "3";
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not null ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 is not null;
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistNot),
	)
}

func TestCheckWhereExistImplicitConversionOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistImplicitConversion].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 = 3;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 = "3";
`,
		newTestResult(),
	)

	runSingleRuleInspectCase(rule, t, "select: check where exist implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where id = "3";
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where id = 3;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistImplicitConversion_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistImplicitConversion].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: unable to check implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 = ?;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: unable to check implicit conversion", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where id = ?;
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExistOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "delete: check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "delete: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1 limit 10 ;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "update: check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 ="1";
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "update: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 ="1" limit 10 ;
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExist_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]delete: check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "[fp]delete: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1 limit ? ;
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "[fp]update: check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 =?;
`,
		newTestResult().addResult(rulepkg.DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "[fp]update: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 =? limit ? ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueriesOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from  exist_db.exist_tb_2);
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
	select v1 from (select v1 from exist_db.exist_tb_1 where v1 in (select v1 from  exist_db.exist_tb_2)) t;
	`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select a.v1 from exist_db.exist_tb_1 a, exist_db.exist_tb_2 b  where a.v1 = b.v1 ;
`,
		newTestResult(),
	)
	// FIXME 子查询 (SELECT COUNT(*) FROM orders WHERE customer_id = customers.customer_id) 返回了每个客户的订单数量，并作为查询结果集中的一个列使用。这个子查询是一个标量子查询，因为它只返回一个值，即每个客户的订单数量。
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
		SELECT customer_name, (SELECT COUNT(*) FROM orders WHERE customer_id = customers.customer_id) AS order_count FROM customers;
		`,
		newTestResult(),
		// newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	// FIXME same with above
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
		SELECT customer_name, (SELECT MAX(age) FROM students) AS student_count FROM customers;
		`,
		newTestResult(),
		// newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	// FIXME same with above
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
		SELECT EXISTS(SELECT 1 FROM customers WHERE customer_name = 'John Doe');
		`,
		newTestResult(),
		// newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
}

func TestCheckWhereExistScalarSubQueries_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from exist_db.exist_tb_2 where v1 = ?);
`,
		newTestResult().addResult(rulepkg.DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (?);
`,
		newTestResult(),
	)
}

func TestCheckIndexesExistBeforeCreatConstraintsOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexesExistBeforeCreateConstraints].Rule
	runSingleRuleInspectCase(rule, t, "Add unique: check indexes exist before creat constraints", DefaultMysqlInspectOffline(),
		`
alter table exist_db.exist_tb_3 Add unique uniq_test(v2);
`, /*not exist index*/
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "Add unique: passing the check indexes exist before creat constraints", DefaultMysqlInspectOffline(),
		`
alter table exist_db.exist_tb_1 Add unique uniq_test(v1); 
`, /*exist index*/
		newTestResult(),
	)
}

func TestCheckSelectForUpdateOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "select : check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 for update;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectForUpdate),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1; 
`,
		newTestResult(),
	)
}

func TestCheckSelectForUpdate_FPOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select : check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 = ? for update;
`,
		newTestResult().addResult(rulepkg.DMLCheckSelectForUpdate),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1= ?; 
`,
		newTestResult(),
	)
}

func TestCheckCollationDatabaseOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`:    `CREATE TABLE exist_db.not_exist_tb_4 (v1 varchar(10)) COLLATE utf8_general_ci;`,
		`alter table`:     `ALTER TABLE exist_db.exist_tb_1 COLLATE utf8_general_ci;`,
		`create database`: `CREATE DATABASE db COLLATE utf8_general_ci;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseCollation].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckDecimalTypeColumnOffline(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckDecimalTypeColumn].Rule
	runSingleRuleInspectCase(rule, t, "create table: check decimal type column", DefaultMysqlInspectOffline(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 float(10));
`,
		newTestResult().addResult(rulepkg.DDLCheckDecimalTypeColumn),
	)
	runSingleRuleInspectCase(rule, t, "alter table: check decimal type column", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 FLOAT ( 10, 0 );
`,
		newTestResult().addResult(rulepkg.DDLCheckDecimalTypeColumn),
	)
	runSingleRuleInspectCase(rule, t, "create table: passing the check decimal type column", DefaultMysqlInspectOffline(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 DECIMAL);
`,
		newTestResult(),
	)
	runSingleRuleInspectCase(rule, t, "alter table: passing the check decimal type column", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 DECIMAL;
`,
		newTestResult(),
	)

}

func TestCheckColumnTypeBlobTextOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckColumnTypeSetOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckColumnTypeEnumOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckUniqueIndexOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`create table`: `CREATE TABLE t1(id INT, c1 INT, UNIQUE INDEX unique_idx (c1));`,
		`alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX unique_id(id);`,
		`create index`: `CREATE UNIQUE INDEX i_u_id ON exist_db.exist_tb_1(id);`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckUniqueIndex].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckWhereExistNullOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckWhereExistNull_FPOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp]select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id = ?;`,
		`[fp]update table`: `UPDATE exist_db.exist_tb_1 SET id = 10 WHERE id = ?;`,
		`[fp]delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id = ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLWhereExistNull].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckNeedlessFuncOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`(1)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, MD5('aaa'), MD5('bbb'));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckNeedlessFunc_FPOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp]INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(?, MD5(?), MD5(?));`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckDatabaseSuffixOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`create database`: `CREATE DATABASE app_service;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DDLCheckDatabaseSuffix].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckTransactionIsolationLevelOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckFuzzySearchOffline(t *testing.T) {
	for _, sql := range []string{
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '_a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a';`,
		`SELECT * FROM (SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a') t;`,
		`SELECT * FROM (SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a') t;`,

		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '%a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '%a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE '_a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '%a';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '%a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 NOT LIKE '_a';`,
		`UPDATE exist_db.exist_tb_1 SET v1 = v1 * (SELECT AVG(id) FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a')/100;`,
		`UPDATE exist_db.exist_tb_1 SET v1 = v1 * (SELECT AVG(id) FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a')/100;`,

		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE '_a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '%a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE '_a';`,
		`DELETE FROM exist_db.exist_tb_1 USING (SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE '%a%') t WHERE t.id = exist_db.exist_tb_1.id;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(rulepkg.DMLCheckFuzzySearch))
	}

	for _, sql := range []string{
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a___';`,
		`SELECT * FROM (SELECT * FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE 'a_') t;`,

		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a___';`,
		`UPDATE exist_db.exist_tb_1 SET v1 = v1 * (SELECT AVG(id) FROM exist_db.exist_tb_1 WHERE v1 NOT LIKE 'a_')/100;`,

		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a____';`,
		`DELETE FROM exist_db.exist_tb_1 USING (SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%') t WHERE t.id = exist_db.exist_tb_1.id;`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckFuzzySearch_FPOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp] "select" unable to check fuzzy search`: `SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
		`[fp] "update" unable to check fuzzy search`: `UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE ?;`,
		`[fp] "delete" unable to check fuzzy search`: `DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckFuzzySearch].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckTablePartitionOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckNumberOfJoinTablesOffline(t *testing.T) {
	// create table for JOIN test
	inspector := DefaultMysqlInspectOffline()
	{
		parent := DefaultMysqlInspectOffline()
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
			newTestResult().addResult(rulepkg.DDLCheckPKName),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckNumberOfJoinTables_FPOffline(t *testing.T) {
	// create table for JOIN test
	inspector := DefaultMysqlInspectOffline()
	{
		parent := DefaultMysqlInspectOffline()
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
			newTestResult().addResult(rulepkg.DDLCheckPKName),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckIsAfterUnionDistinctOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT 1, 2 UNION SELECT 'a', 'b';`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckIsAfterUnionDistinct_FPOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT ?, ? UNION SELECT ?, ?;`,
	} {
		runSingleRuleInspectCase(
			rulepkg.RuleHandlerMap[rulepkg.DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckIsExistLimitOffsetOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func Test_DDLCheckNameUseENAndUnderline_ShouldErrorOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(rulepkg.DDLCheckObjectNameUseCN))
	}
}

func Test_DDLCheckNameUseENAndUnderline_ShouldNotErrorOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckIndexOption_ShouldNot_QueryDBOffline(t *testing.T) {
	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(1)index on new db new column`,
		DefaultMysqlInspectOffline(),
		`CREATE TABLE t1(id int, name varchar(100), INDEX idx_name(name))`,
		newTestResult())

	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(2)index on new db new column`,
		DefaultMysqlInspectOffline(),
		`CREATE TABLE t1(id int, name varchar(100));
ALTER TABLE t1 ADD INDEX idx_name(name);
`,
		newTestResult(), newTestResult())

	runSingleRuleInspectCase(
		rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexOption].Rule,
		t,
		`(3)index on old db new column`,
		DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(100);
ALTER TABLE exist_db.exist_tb_1 ADD INDEX idx_v3(v3);
`,
		newTestResult(), newTestResult())
}

func Test_DDL_CHECK_PK_NAMEOffline(t *testing.T) {
	for _, sql := range []string{
		`create table t1(id int, primary key pk_t1(id))`,
		`create table t1(id int, primary key PK_T1(id))`,
		`alter table exist_db.exist_tb_2 Add primary key PK_EXIST_TB_2(id)`} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKName].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}

	for _, sql := range []string{
		`create table t1(id int, primary key wrongPK(id))`,
		`alter table exist_db.exist_tb_2 Add primary key wrongPK(id)`,
		`create table t1(id int, primary key(id))`,
		`alter table exist_db.exist_tb_2 Add primary key(id)`} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKName].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(rulepkg.DDLCheckPKName))
	}
}

func Test_PerfectParseOffline(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckWhereIsInvalid].Rule, t, "", DefaultMysqlInspectOffline(), `
SELECT * FROM exist_db.exist_tb_1;
OPTIMIZE TABLE exist_db.exist_tb_1;
SELECT * FROM exist_db.exist_tb_2;
`, newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid),
		newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"),
		newTestResult().addResult(rulepkg.DMLCheckWhereIsInvalid))
}

func Test_DDLCheckCreateViewOffline(t *testing.T) {
	for _, sql := range []string{
		`create view v as select * from t1`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateView].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(rulepkg.DDLCheckCreateView))
	}

	for _, sql := range []string{
		`create table t1(id int)`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateView].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func Test_DDLCheckCreateTriggerOffline(t *testing.T) {
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLCheckCreateTrigger))
	}

	for _, sql := range []string{
		`CREATE my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATEmy_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE trigger_1 BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER my_trigger BEEEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}
}

func Test_DDLCheckCreateFunctionOffline(t *testing.T) {
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
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLCheckCreateFunction))
	}

	for _, sql := range []string{
		`create function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`create123 function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER='sqle_op'@'localhost' hello (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
	} {
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}
}

func Test_DDLCheckCreateProcedureOffline(t *testing.T) {
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
			DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").
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
			DefaultMysqlInspectOffline(), sql,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性"))
	}
}

func TestDDLNotAllowRenamingOffline(t *testing.T) {
	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "success", DefaultMysqlInspectOffline(), "ALTER TABLE exist_tb_1 MODIFY v1 CHAR(10);", newTestResult())

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "change 1", DefaultMysqlInspectOffline(), "ALTER TABLE exist_tb_1 CHANGE v1 a BIGINT;", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "change 2", DefaultMysqlInspectOffline(), "ALTER TABLE exist_tb_1 RENAME COLUMN v1 TO a", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "rename 1", DefaultMysqlInspectOffline(), "RENAME TABLE exist_tb_1 TO test", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

	runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLNotAllowRenaming].Rule, t, "rename 2", DefaultMysqlInspectOffline(), "ALTER TABLE exist_tb_1 RENAME TO test", newTestResult().addResult(rulepkg.DDLNotAllowRenaming))

}

func TestDMLCheckLimitOffsetNum(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckLimitOffsetNum].Rule
	rule.Params.SetParamValue(rulepkg.DefaultSingleParamKeyName, "4")
	runSingleRuleInspectCase(
		rule,
		t,
		`(1)select with limit offset`,
		DefaultMysqlInspectOffline(),
		`SELECT * FROM tbl LIMIT 5,10`,
		newTestResult().addResult(rulepkg.DMLCheckLimitOffsetNum, 5, 4))

	runSingleRuleInspectCase(
		rule,
		t,
		`(2)select with limit explicit offset`,
		DefaultMysqlInspectOffline(),
		`SELECT * FROM tbl LIMIT 10 OFFSET 5`,
		newTestResult().addResult(rulepkg.DMLCheckLimitOffsetNum, 5, 4))

}

func TestDMLCheckUpdateOrDeleteHasWhere(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckUpdateOrDeleteHasWhere].Rule
	t.Run(`(1)update with where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`UPDATE t1 SET col1 = col1 + 1 WHERE a = 2`,
			newTestResult())
	})
	t.Run(`(2)update with where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`UPDATE t1 SET col1 = col1 + 1 WHERE a = 2`,
			newTestResult())
	})
	t.Run(`(3)update with where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`UPDATE t1 SET col1 = col1 + 1, col2 = col1 WHERE a = 2`,
			newTestResult())
	})
	t.Run(`(4)update without where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`UPDATE t1 SET col1 = col1 + 1;`,
			newTestResult().addResult(rulepkg.DMLCheckUpdateOrDeleteHasWhere))
	})
	t.Run(`(5)delete with where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`DELETE FROM t1 WHERE a = 2`,
			newTestResult())
	})
	t.Run(`(6)delete with where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`DELETE t1 FROM t1 LEFT JOIN t2 ON t1.id=t2.id WHERE t2.id IS NULL`,
			newTestResult())
	})
	t.Run(`(7)delete without where`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`DELETE FROM t1`,
			newTestResult().addResult(rulepkg.DMLCheckUpdateOrDeleteHasWhere))
	})
}

func TestDMLCheckJoinHasOn(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLCheckHasJoinCondition].Rule
	// this rule cover ON USING WHERE clause
	caseWithResult := map[string] /*case name*/ string{
		"no condition at all": `
			SELECT * 
			FROM t1 
			JOIN t2`,
		"no condition at all 2": `
			UPDATE employees
			JOIN merits
			SET salary = salary + salary * 0.015`,
		"on condition": `
			SELECT * 
			FROM t1
			JOIN t2 ON t2.a = t1.a
			JOIN t3
			JOIN t4 ON t4.a = t1.a`,
		"using condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 USING(id)
			INNER JOIN table3 t3`,
		"where condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2
			INNER JOIN table3 t3
			WHERE t1.id = t2.id`,
		"mix where and using condition": `
			SELECT *
			FROM table1 t1 
			JOIN table2 t2 USING(id)
			JOIN table3 t3
			WHERE t2.id = t1.id
			AND t3.id = 1`,
		"mix where and using condition2": `
			SELECT *
			FROM table1 t1 
			JOIN table2 t2 
			JOIN table3 t3 USING(id)
			WHERE t2.id = t2.id
			AND t3.id = 1`,
		"mix where and on condition": `
			SELECT *
			FROM table1 t1 
			JOIN table2 t2 ON t1.id = t2.id
			JOIN table3 t3
			WHERE t2.id = t2.id
			AND t3.id = 1`,
		"mix using and on condition": `
			SELECT *
			FROM table1 t1 
			JOIN table2 t2 ON t1.id = t2.id
			JOIN table3 t3 USING(id)
			JOIN table4 t4
			WHERE t2.id = t2.id
			AND t3.id = 1`,
	}
	caseWithoutResult := map[string] /*case name*/ string{
		"not join": `
			SELECT t1.id,t1.name FROM t1
		`,
		"on condition": `
			SELECT * FROM t1
			JOIN t2 ON t2.a = t1.a
			JOIN t3 ON t3.b = t2.b
			JOIN t4 ON t4.a = t1.a`,
		"on condition 2": `
			UPDATE employees
		    JOIN  merits ON employees.performance = merits.performance
			SET salary = salary + salary * 0.015`,
		"using condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 USING(id)
			INNER JOIN table3 t3 USING(id)`,
		"where condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2
			INNER JOIN table3 t3
			WHERE t1.id = t2.id
			AND t2.id = t3.id`,
		"mix using and where condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 USING(id)
			INNER JOIN table3 t3
			WHERE t2.id = t3.id`,
		"mix using and on condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 USING(id)
			INNER JOIN table3 t3 ON t2.id = t3.id`,
		"mix where and on condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 ON t1.id=t2.id
			INNER JOIN table3 t3 
			WHERE t2.id = t3.id`,
		"mix where,using and on condition": `
			SELECT *
			FROM table1 t1 
			INNER JOIN table2 t2 ON t1.id=t2.id
			INNER JOIN table3 t3 
			INNER JOIN table4 t4 USING(id)
			WHERE t2.id = t3.id`,
	}
	for name, sql := range caseWithoutResult {
		t.Run(name, func(t *testing.T) {
			runSingleRuleInspectCase(
				rule,
				t,
				name,
				DefaultMysqlInspectOffline(),
				sql,
				newTestResult(),
			)
		})
	}
	for name, sql := range caseWithResult {
		t.Run(name, func(t *testing.T) {
			runSingleRuleInspectCase(
				rule,
				t,
				name,
				DefaultMysqlInspectOffline(),
				sql,
				newTestResult().addResult(rulepkg.DMLCheckHasJoinCondition),
			)
		})
	}

}

func TestDMLHintCountFuncWithCol(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLHintCountFuncWithCol].Rule
	t.Run(`select count(col)`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b, COUNT(c) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult().addResult(rulepkg.DMLHintCountFuncWithCol))
	})
	t.Run(`select count(*)`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b, COUNT(*) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult())
	})
	t.Run(`select count(1)`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b, COUNT(1) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult())
	})
	t.Run(`select count(distinct(col))`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b, COUNT(distinct(col)) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult())
	})
	t.Run(`select count(distinct col)`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b, COUNT(distinct col) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult())
	})
	t.Run(`select fields contain different count(1) trigger rule`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b,COUNT(distinct(col)),COUNT(distinct col), COUNT(col) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult().addResult(rulepkg.DMLHintCountFuncWithCol))
	})
	t.Run(`select fields contain different count(2) `, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT a, b,COUNT(distinct col), COUNT(distinct(col)) AS t FROM test_table GROUP BY a,b ORDER BY a,t DESC;`,
			newTestResult())
	})
	t.Run(`select fields contain different count(3) `, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT t1.a, t1.b, COUNT(distinct t1.col) AS distinct_count, t2.col_count   
			FROM test_table AS t1   
			LEFT JOIN (SELECT c, d, COUNT(col_1) AS col_count FROM test_table_1) AS t2   
			ON t1.a = t2.c AND t1.b = t2.d;`,
			newTestResult().addResult(rulepkg.DMLHintCountFuncWithCol))
	})
	t.Run(`select fields contain different count(4) `, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`SELECT t1.a, t1.b, t2.col_count   
			FROM test_table AS t1   
			LEFT JOIN (SELECT c, COUNT(distinct col_1) AS distinct_count FROM test_table_1) AS t2   
			ON t1.a = t2.c AND t1.b = t2.d;`,
			newTestResult())
	})
}

func TestDDLCheckAutoIncrementFieldNum(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLCheckAutoIncrementFieldNum].Rule
	t.Run(`create table with one AUTO_INCREMENT field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE IF NOT EXISTS tbl1(
				id INT UNSIGNED AUTO_INCREMENT,
				title VARCHAR(100) NOT NULL,
				PRIMARY KEY ( id )
			 );`,
			newTestResult())
	})
	t.Run(`create table without AUTO_INCREMENT field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE IF NOT EXISTS tbl2(title VARCHAR(100) NOT NULL);`,
			newTestResult())
	})
	t.Run(`create table with two AUTO_INCREMENT fields`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			` CREATE TABLE IF NOT EXISTS tbl1(id INT UNSIGNED AUTO_INCREMENT primary key,id2 BIGINT UNSIGNED AUTO_INCREMENT);`,
			newTestResult().addResult(rulepkg.DDLCheckAutoIncrementFieldNum))
	})
}

func TestDDLAvoidText(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLAvoidText].Rule
	t.Run(`create table with text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE IF NOT EXISTS tbl1(
				id INT UNSIGNED PRIMARY KEY,
				product_code VARCHAR(10),
				title blob  NOT NULL
			 );`,
			newTestResult())
	})

	t.Run(`create table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE your_table_name (
				username VARCHAR(50),
				product_code VARCHAR(10),
				title blob  NOT NULL,
				title1 TINYBLOB,
				title2 MEDIUMBLOB,
				title3 LONGBLOB,
				name VARCHAR(50),
				PRIMARY KEY (username, product_code)
			);
			`,
			newTestResult())
	})

	t.Run(`create table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE your_table_name (
				username VARCHAR(50),
				product_code VARCHAR(10),
				blob_column blob  NOT NULL,
				title text not null,
				name VARCHAR(50),
				PRIMARY KEY (username, product_code)
			);
			`,
			newTestResult().addResult(rulepkg.DDLAvoidText, "title"))
	})

	t.Run(`create table with text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE IF NOT EXISTS tbl1(
				id INT UNSIGNED PRIMARY KEY,
				title text  NOT NULL
			 );`,
			newTestResult())
	})

	t.Run(`create table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE your_table_name (
				username VARCHAR(50),
				product_code VARCHAR(10),
				title text  NOT NULL,
				PRIMARY KEY (username, product_code)
			);
			`,
			newTestResult())
	})

	t.Run(`create table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE your_table_name (
				username VARCHAR(50),
				product_code VARCHAR(10),
				title text  NOT NULL,
				age int,
				PRIMARY KEY (username, product_code)
			);
			`,
			newTestResult().addResult(rulepkg.DDLAvoidText, "title"))
	})

	t.Run(`create table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE your_table_name (
				username VARCHAR(50),
				product_code VARCHAR(10),
				title text  NOT NULL,
				title1 TINYTEXT,
				title2 MEDIUMTEXT,
				title3 LONGTEXT,
				name VARCHAR(50),
				PRIMARY KEY (username, product_code)
			);
			`,
			newTestResult().addResult(rulepkg.DDLAvoidText, "title，title1，title2，title3"))
	})

	t.Run(`alter table with text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspect(),
			`ALTER TABLE exist_db.exist_tb_1
			ADD new_column_name text
			`,
			newTestResult().addResult(rulepkg.DDLAvoidText, "new_column_name"))
	})

	t.Run(`alter table with text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspect(),
			`ALTER TABLE exist_db.exist_tb_1
			ADD new_column_name blob
			`,
			newTestResult())
	})

	t.Run(`alter table with text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspect(),
			`ALTER TABLE exist_db.exist_tb_1
			ADD new_column_name blob,
			ADD title text not null
			`,
			newTestResult().addResult(rulepkg.DDLAvoidText, "title"))
	})

	t.Run(`alter table without text field`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER TABLE t1
			ADD new_column_name varchar(20);
			`,
			newTestResult())
	})
}

func TestDDLAvoidFullText(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLAvoidFullText].Rule
	// 全文索引
	t.Run(`create table without fulltext index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE example (
				id INT primary key,
				content TEXT,
				INDEX idx_id_content (id, content)
			);`,
			newTestResult())
	})

	t.Run(`create common table`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE example (
				id INT primary key,
				content TEXT
			);`,
			newTestResult())
	})

	t.Run(`create fulltext index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE FULLTEXT INDEX index_name
			ON table_name (column_name);`,
			newTestResult().addResult(rulepkg.DDLAvoidFullText))
	})

	t.Run(`create index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE INDEX index_name ON table_name (column_name);`,
			newTestResult())
	})

	t.Run(`create unique index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE unique INDEX index_name ON table_name (column_name);`,
			newTestResult())
	})

	t.Run(`create table with fulltext index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE example (
				id INT,
				content TEXT,
				FULLTEXT INDEX idx_content (content)
			);`,
			newTestResult().addResult(rulepkg.DDLAvoidFullText))
	})

	t.Run(`alter table add fulltext`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER TABLE your_table_name
			ADD FULLTEXT INDEX idx_name (name);`,
			newTestResult().addResult(rulepkg.DDLAvoidFullText))
	})

	t.Run(`alter table`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER TABLE example
			ADD INDEX idx_name (name),
			ADD UNIQUE INDEX idx_email (email);`,
			newTestResult())
	})
}

func TestDDLAvoidGeometry(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLAvoidGeometry].Rule
	// 空间字段
	t.Run(`create table with column point`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g POINT);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column GEOMETRY`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRY);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column LINESTRING`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g LINESTRING);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column POLYGON`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g POLYGON);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column MULTIPOINT`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g MULTIPOINT);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column MULTILINESTRING`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g MULTILINESTRING);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column MULTIPOLYGON`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g MULTIPOLYGON);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create table with column GEOMETRYCOLLECTION`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g GEOMETRYCOLLECTION);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	// alter add
	t.Run(`alter table with column point`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column t point;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column GEOMETRY`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g GEOMETRY;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column LINESTRING`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g LINESTRING;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column POLYGON`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g POLYGON;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column MULTIPOINT`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g MULTIPOINT;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column MULTILINESTRING`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g MULTILINESTRING;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column MULTIPOLYGON`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g MULTIPOLYGON;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with column GEOMETRYCOLLECTION`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column g GEOMETRYCOLLECTION;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	// 空间索引
	t.Run(`create table with GEOMETRY index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, g POINT, SPATIAL INDEX(g));`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`Create a normal index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE TABLE t (id INT PRIMARY KEY, name varchar(50), INDEX(name));`,
			newTestResult())
	})
	t.Run(`alter table a GEOMETRY index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER TABLE geom ADD SPATIAL INDEX(g);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table a normal index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER TABLE geom ADD INDEX(g);`,
			newTestResult())
	})
	t.Run(`create a GEOMETRY index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE SPATIAL INDEX g ON geom (g);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`create a normal index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`CREATE INDEX g ON geom (g);`,
			newTestResult())
	})
	t.Run(`alter table with geo index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter table table_1 add SPATIAL index index_2(g);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	// 添加多个字段
	t.Run(`alter table with column MULTIPOLYGON`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column name varchar(20),add column g MULTIPOLYGON;`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with normal columns`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter TABLE t add column name varchar(20),add column age int;`,
			newTestResult())
	})
	t.Run(`alter table with geo column and normal index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter table table_1 add column g point, add index index_1(column_name);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with normal column and geo index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter table table_1 add column name varchar(20), add SPATIAL INDEX(g);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
	t.Run(`alter table with normal column and geo index`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`alter table table_1 add SPATIAL index index_2(g);`,
			newTestResult().addResult(rulepkg.DDLAvoidGeometry))
	})
}

func TestDMLAvoidWhereEqualNull(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DMLAvoidWhereEqualNull].Rule
	t.Run(`select a = null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a = null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a is null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a is null;`,
			newTestResult())
	})
	t.Run(`select a is not null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a is not null;`,
			newTestResult())
	})
	t.Run(`select a != null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a != null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a <> null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a <> null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a >= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a >= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a > null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a > null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a <= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a <= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`select a < null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`select * from t1 where a < null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`update a = null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a = null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`update a != null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a != null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`update a <> null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a <> null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`update a >= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a >= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`update a is null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a is null;`,
			newTestResult())
	})
	t.Run(`update a is not null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`update t1 set name='v1' where a is not null;`,
			newTestResult())
	})
	t.Run(`delete a = null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a = null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a != null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a != null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a <> null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a <> null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a >= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a >= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a < null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a > null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a < null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a < null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`delete a is null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a is null;`,
			newTestResult())
	})
	t.Run(`delete a is not null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`delete from t1 where a is not null;`,
			newTestResult())
	})
	t.Run(`insert select v1 = null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 = null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 != null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 != null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 <> null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 <> null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 >= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 >= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 > null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 > null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 <= null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 <= null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 < null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 < null;`,
			newTestResult().addResult(rulepkg.DMLAvoidWhereEqualNull))
	})
	t.Run(`insert select v1 is null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 is null;`,
			newTestResult())
	})
	t.Run(`insert select v1 is not null`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`insert into t1(column1) select column1 from t2 where v1 is not null;`,
			newTestResult())
	})
}

func TestDDLAvoidEvent(t *testing.T) {
	rule := rulepkg.RuleHandlerMap[rulepkg.DDLAvoidEvent].Rule
	t.Run(`create event`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`create event my_event on schedule every 10 second do update schema.table set mycol = mycol + 1;`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
	t.Run(`create event with DEFINER`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`create DEFINER=user event my_event on schedule every 10 second do update schema.table set mycol = mycol + 1;`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
	t.Run(`alter event`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER EVENT your_event_name
			ON SCHEDULE
			  EVERY 1 DAY
			  STARTS '2023-01-01 00:00:00'
			DO
			  -- 修改事件的具体操作
			  UPDATE your_table SET your_column = your_value WHERE your_condition;
			`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
	t.Run(`alter event with DEFINER`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`ALTER DEFINER = user EVENT your_event_name
			ON SCHEDULE
			  EVERY 1 DAY
			  STARTS '2023-01-01 00:00:00'
			DO
			  -- 修改事件的具体操作
			  UPDATE your_table SET your_column = your_value WHERE your_condition;
			`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
	t.Run(`create event with blank line`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`

			
			create event my_event on schedule every 10 second do update schema.table set mycol = mycol + 1;`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
	t.Run(`create event with space`, func(t *testing.T) {
		runSingleRuleInspectCase(
			rule,
			t,
			``,
			DefaultMysqlInspectOffline(),
			`       create event my_event on schedule every 10 second do update schema.table set mycol = mycol + 1;`,
			newTestResult().add(driverV2.RuleLevelWarn, "", "语法错误或者解析器不支持，请人工确认SQL正确性").addResult(rulepkg.DDLAvoidEvent))
	})
}
