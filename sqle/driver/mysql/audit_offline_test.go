package mysql

import (
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/sqle/sqle/sqle/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func DefaultMysqlInspectOffline() *Inspect {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &Inspect{
		log:  log.NewEntry(),
		inst: nil,
		Ctx: &Context{
			sysVars: map[string]string{
				"lower_case_table_names": "0",
			},
		},
		cnf: &Config{
			DDLOSCMinSize:      -1,
			DMLRollbackMaxRows: -1,
		},
		isOfflineAudit: true,
	}
}

func TestCheckSelectAllOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "select_from: all columns", DefaultMysqlInspectOffline(),
		"select * from exist_db.exist_tb_1 where id =1;",
		newTestResult().addResult(DMLDisableSelectAllColumn),
	)
}

func TestCheckInvalidDropOffline(t *testing.T) {
	handler := RuleHandlerMap[DDLDisableDropStatement]
	delete(RuleHandlerMap, DDLDisableDropStatement)
	defer func() {
		RuleHandlerMap[DDLDisableDropStatement] = handler
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
	runDefaultRulesInspectCase(t, "select_from: has where condition", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where id > 1;",
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(1)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(2)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(3)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where id=id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "select_from: no where condition(4)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "update: has where condition", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "update: no where condition(1)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1';",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(2)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where 1=1 and 2=2;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(3)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where id=id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "update: no where condition(4)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: has where condition", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where id = 1;",
		newTestResult())

	runDefaultRulesInspectCase(t, "delete: no where condition(1)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(2)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=1 and 2=2;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(3)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=1 and id=id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "delete: no where condition(4)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=1 and exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))
}

func TestCheckWhereInvalid_FPOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(1)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where id=?;",
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: has where condition(2)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where exist_tb_1.id=?;",
		newTestResult(),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(1)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where 1=? and 2=2;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)
	runDefaultRulesInspectCase(t, "[pf]select_from: no where condition(2)", DefaultMysqlInspectOffline(),
		"select id from exist_db.exist_tb_1 where ?=?;",
		newTestResult().addResult(DMLCheckWhereIsInvalid),
	)

	runDefaultRulesInspectCase(t, "[pf]update: has where condition", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1='v1' where id = ?;",
		newTestResult())

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(1)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=?;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(2)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where 1=1 and 2=2;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(3)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where id=id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]update: no where condition(4)", DefaultMysqlInspectOffline(),
		"update exist_db.exist_tb_1 set v1=? where exist_tb_1.id=exist_tb_1.id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]delete: no where condition(1)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=? and ?=?;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))

	runDefaultRulesInspectCase(t, "[pf]delete: no where condition(2)", DefaultMysqlInspectOffline(),
		"delete from exist_db.exist_tb_1 where 1=? and id=id;",
		newTestResult().addResult(DMLCheckWhereIsInvalid))
}

func TestCheckCreateTableWithoutIfNotExistsOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: need \"if not exists\"", DefaultMysqlInspectOffline(),
		`
CREATE TABLE exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKWithoutIfNotExists),
	)
}

func TestCheckObjectNameUsingKeywordOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: using keyword", DefaultMysqlInspectOffline(),
		"CREATE TABLE if not exists exist_db.`select` ("+
			"id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT \"unit test\","+
			"v1 varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"`create` varchar(255) NOT NULL DEFAULT \"unit test\" COMMENT \"unit test\","+
			"PRIMARY KEY (id),"+
			"INDEX `show` (v1)"+
			")ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT=\"unit test\";",
		newTestResult().addResult(DDLCheckObjectNameUsingKeyword, "select, create, show").
			addResult(DDLCheckIndexPrefix),
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
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length64),
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: table length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
%s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "create_table: index length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_%s (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "alter_table: table length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME %s;`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "alter_table:Add column length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "alter_table:change column length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 %s varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "alter_table: Add index length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 ADD index idx_%s (v1);`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)

	runDefaultRulesInspectCase(t, "alter_table:rename index length > 64", DefaultMysqlInspectOffline(),
		fmt.Sprintf(`
ALTER TABLE exist_db.exist_tb_1 RENAME index idx_1 TO idx_%s;`, length65),
		newTestResult().addResult(DDLCheckObjectNameLength),
	)
}

func TestCheckPrimaryKeyOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: primary key exist", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not exist", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKNotExist),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(1)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL KEY DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not auto increment(2)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKWithoutAutoIncrement),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(1)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKWithoutBigintUnsigned),
	)

	runDefaultRulesInspectCase(t, "create_table: primary key not bigint unsigned(2)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckPKWithoutBigintUnsigned),
	)
}

func TestCheckColumnCharLengthOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: check char(20)", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: check char(21)", DefaultMysqlInspectOffline(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 char(21) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult().addResult(DDLCheckColumnCharLength),
	)
}

func TestCheckIndexCountOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: index <= 5", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: index > 5", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(DDLCheckIndexCount),
	)
}

func TestCheckCompositeIndexMaxOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: composite index columns <= 3", DefaultMysqlInspectOffline(),
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

	runDefaultRulesInspectCase(t, "create_table: composite index columns > 3", DefaultMysqlInspectOffline(),
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
		newTestResult().addResult(DDLCheckCompositeIndexMax),
	)
}

func TestCheckTableWithoutInnodbUtf8mb4Offline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: ok", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
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
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_b1 (b1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckIndexedColumnWithBolb),
	)

	runDefaultRulesInspectCase(t, "create_table: disable index column blob (2)", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob UNIQUE KEY COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckIndexedColumnWithBolb),
	)

	handler := RuleHandlerMap[DDLCheckAlterTableNeedMerge]
	delete(RuleHandlerMap, DDLCheckAlterTableNeedMerge)
	defer func() {
		RuleHandlerMap[DDLCheckAlterTableNeedMerge] = handler
	}()

}

func TestDisableForeignKeyOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: has foreign key", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
FOREIGN KEY (id) REFERENCES exist_tb_1(id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLDisableFK),
	)
}

func TestCheckTableCommentOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: table without comment", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
`,
		newTestResult().addResult(DDLCheckTableWithoutComment),
	)
}

func TestCheckColumnCommentOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column without comment", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckColumnWithoutComment),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without comment(1)", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL DEFAULT "unit test";
`,
		newTestResult().addResult(DDLCheckColumnWithoutComment),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without comment(2)", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v2 v3 varchar(255) NOT NULL DEFAULT "unit test" ;
`,
		newTestResult().addResult(DDLCheckColumnWithoutComment),
	)
}

func TestCheckIndexPrefixOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: index prefix not idx_", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckIndexPrefix),
	)

	runDefaultRulesInspectCase(t, "alter_table: index prefix not idx_", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD INDEX index_1(v1);
`,
		newTestResult().addResult(DDLCheckIndexPrefix),
	)

	runDefaultRulesInspectCase(t, "create_index: index prefix not idx_", DefaultMysqlInspectOffline(),
		`
CREATE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(DDLCheckIndexPrefix),
	)

	for _, sql := range []string{
		`create table exist_db.t1(id int, c1 varchar(10), index idx_c1(c1))`,
		`create table exist_db.t1(id int, c1 varchar(10), index IDX_C1(c1))`,
		`create index IDX_v1 ON exist_db.exist_tb_1(v1);`,
		`create index idx_V1 ON exist_db.exist_tb_1(v1);`,
		`alter table exist_db.exist_tb_1 Add index idx_v1(v1);`,
		`alter table exist_db.exist_tb_1 Add index IDX_V1(v1);`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckIndexPrefix].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckUniqueIndexPrefixOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
UNIQUE INDEX index_1 (v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckUniqueIndexPrefix),
	)

	runDefaultRulesInspectCase(t, "alter_table: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD UNIQUE INDEX index_1(v1);
`,
		newTestResult().addResult(DDLCheckUniqueIndexPrefix),
	)

	runDefaultRulesInspectCase(t, "create_index: unique index prefix not uniq_", DefaultMysqlInspectOffline(),
		`
CREATE UNIQUE INDEX index_1 ON exist_db.exist_tb_1(v1);
`,
		newTestResult().addResult(DDLCheckUniqueIndexPrefix),
	)

	for _, sql := range []string{
		`create table exist_db.t1(id int, c1 varchar(10), unique index uniq_c1(c1))`,
		`create table exist_db.t1(id int, c1 varchar(10), unique index UNIQ_C1(c1))`,
		`create unique index uniq_v1 ON exist_db.exist_tb_1(v1);`,
		`create unique index UNIQ_V1 ON exist_db.exist_tb_1(v1);`,
		`alter table exist_db.exist_tb_1 Add unique index uniq_v1(v1);`,
		`alter table exist_db.exist_tb_1 Add unique index UNIQ_V1(v1);`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckUniqueIndexPrefix].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckColumnDefaultOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckColumnWithoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 varchar(255) NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDLCheckColumnWithoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: auto increment column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test";
`,
		newTestResult(),
	)

	runDefaultRulesInspectCase(t, "alter_table: blob column without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob COMMENT "unit test";
`,
		newTestResult(),
	)
}

func TestCheckColumnTimestampDefaultOffline(t *testing.T) {
	handler := RuleHandlerMap[DDLCheckColumnWithoutDefault]
	delete(RuleHandlerMap, DDLCheckColumnWithoutDefault)
	defer func() {
		RuleHandlerMap[DDLCheckColumnWithoutDefault] = handler
	}()

	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 timestamp COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckColumnTimestampWitoutDefault),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 timestamp NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDLCheckColumnTimestampWitoutDefault),
	)
}

func TestCheckColumnBlobNotNullOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 blob NOT NULL COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckColumnBlobWithNotNull),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob NOT NULL COMMENT "unit test";
`,
		newTestResult().addResult(DDLCheckColumnBlobWithNotNull),
	)
}

func TestCheckColumnBlobDefaultNullOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "create_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 blob DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult().addResult(DDLCheckColumnBlobDefaultIsNotNull),
	)

	runDefaultRulesInspectCase(t, "alter_table: column timestamp without default", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 blob DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult().addResult(DDLCheckColumnBlobDefaultIsNotNull),
	)
}

func TestCheckDMLWithLimitOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 limit 1;
`,
		newTestResult().addResult(DMLCheckWithLimit),
	)
}

func TestCheckDMLWithLimit_FPOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]update: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(DMLCheckWithLimit),
	)

	runDefaultRulesInspectCase(t, "[fp]delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=? limit ?;
`,
		newTestResult().addResult(DMLCheckWithLimit),
	)
}

func TestCheckDMLWithOrderByOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "update: with order by", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(DMLCheckWithOrderBy),
	)

	runDefaultRulesInspectCase(t, "delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by v1;
`,
		newTestResult().addResult(DMLCheckWithOrderBy),
	)
}

func TestCheckDMLWithOrderBy_FPOffline(t *testing.T) {
	runDefaultRulesInspectCase(t, "[fp]update: with order by", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1="2" where id=1 order by ?;
`,
		newTestResult().addResult(DMLCheckWithOrderBy),
	)

	runDefaultRulesInspectCase(t, "[fp]delete: with limit", DefaultMysqlInspectOffline(),
		`
UPDATE exist_db.exist_tb_1 Set v1=? where id=? order by ?;
`,
		newTestResult().addResult(DMLCheckWithOrderBy),
	)
}

func TestCheckInsertColumnsExistOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "insert: check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 values (1,"1","1"),(2,"2","2");
`,
		newTestResult().addResult(DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2");
`,
		newTestResult(),
	)
}

func TestCheckInsertColumnsExist_FPOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckInsertColumnsExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]insert: check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 values (?,?,?),(?,?,?);
`,
		newTestResult().addResult(DMLCheckInsertColumnsExist),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check columns exist", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMaxOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckBatchInsertListsMax].Rule
	// defult 5000,  unit testing :4
	rule.Value = "4"
	runSingleRuleInspectCase(rule, t, "insert:check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4"),(5,"5","5");
`,
		newTestResult().addResult(DMLCheckBatchInsertListsMax, rule.Value),
	)

	runSingleRuleInspectCase(rule, t, "insert: passing the check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (1,"1","1"),(2,"2","2"),(3,"3","3"),(4,"4","4");
`,
		newTestResult(),
	)
}

func TestCheckBatchInsertListsMax_FPOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckBatchInsertListsMax].Rule
	// defult 5000, unit testing :4
	rule.Value = "4"
	runSingleRuleInspectCase(rule, t, "[fp]insert:check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult().addResult(DMLCheckBatchInsertListsMax, rule.Value),
	)

	runSingleRuleInspectCase(rule, t, "[fp]insert: passing the check batch insert lists max", DefaultMysqlInspectOffline(),
		`
insert into exist_db.exist_tb_1 (id,v1,v2) values (?,?,?),(?,?,?),(?,?,?),(?,?,?);
`,
		newTestResult(),
	)
}

func TestCheckPkProhibitAutoIncrementOffline(t *testing.T) {
	rule := RuleHandlerMap[DDLCheckPKProhibitAutoIncrement].Rule
	runSingleRuleInspectCase(rule, t, "create_table: primary key not auto increment", DefaultMysqlInspectOffline(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT DEFAULT "unit test" COMMENT "unit test" ,
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	`,
		newTestResult().addResult(DDLCheckPKProhibitAutoIncrement),
	)

}

func TestCheckWhereExistFuncOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckWhereExistFunc].Rule
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
	rule := RuleHandlerMap[DMLCheckWhereExistFunc].Rule
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
	rule := RuleHandlerMap[DMLCheckWhereExistNot].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist <> ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 <> "3";
`,
		newTestResult().addResult(DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not like ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 not like "%3%";
`,
		newTestResult().addResult(DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist != ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 != "3";
`,
		newTestResult().addResult(DMLCheckWhereExistNot),
	)
	runSingleRuleInspectCase(rule, t, "select: check where exist not null ", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v2 is not null;
`,
		newTestResult().addResult(DMLCheckWhereExistNot),
	)
}

func TestCheckWhereExistImplicitConversionOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckWhereExistImplicitConversion].Rule
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
	rule := RuleHandlerMap[DMLCheckWhereExistImplicitConversion].Rule
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
	rule := RuleHandlerMap[DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "delete: check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(DMLCheckLimitMustExist),
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
		newTestResult().addResult(DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "update: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 ="1" limit 10 ;
`,
		newTestResult(),
	)
}

func TestCheckLimitMustExist_FPOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckLimitMustExist].Rule
	runSingleRuleInspectCase(rule, t, "[fp]delete: check limit must exist", DefaultMysqlInspectOffline(),
		`
delete from exist_db.exist_tb_1;
`,
		newTestResult().addResult(DMLCheckLimitMustExist),
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
		newTestResult().addResult(DMLCheckLimitMustExist),
	)
	runSingleRuleInspectCase(rule, t, "[fp]update: passing the check limit must exist", DefaultMysqlInspectOffline(),
		`
update exist_db.exist_tb_1 set v1 =? limit ? ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueriesOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "select: check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from  exist_db.exist_tb_2);
`,
		newTestResult().addResult(DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select a.v1 from exist_db.exist_tb_1 a, exist_db.exist_tb_2 b  where a.v1 = b.v1 ;
`,
		newTestResult(),
	)
}

func TestCheckWhereExistScalarSubQueries_FPOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckWhereExistScalarSubquery].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select: check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (select v1 from exist_db.exist_tb_2 where v1 = ?);
`,
		newTestResult().addResult(DMLCheckWhereExistScalarSubquery),
	)
	runSingleRuleInspectCase(rule, t, "[fp]select: passing the check where exist scalar sub queries", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 in (?);
`,
		newTestResult(),
	)
}

func TestCheckIndexesExistBeforeCreatConstraintsOffline(t *testing.T) {
	rule := RuleHandlerMap[DDLCheckIndexesExistBeforeCreateConstraints].Rule
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
	rule := RuleHandlerMap[DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "select : check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 for update;
`,
		newTestResult().addResult(DMLCheckSelectForUpdate),
	)
	runSingleRuleInspectCase(rule, t, "select: passing the check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1; 
`,
		newTestResult(),
	)
}

func TestCheckSelectForUpdate_FPOffline(t *testing.T) {
	rule := RuleHandlerMap[DMLCheckSelectForUpdate].Rule
	runSingleRuleInspectCase(rule, t, "[fp]select : check exist select_for_update", DefaultMysqlInspectOffline(),
		`
select v1 from exist_db.exist_tb_1 where v1 = ? for update;
`,
		newTestResult().addResult(DMLCheckSelectForUpdate),
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
			RuleHandlerMap[DDLCheckDatabaseCollation].Rule,
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
			RuleHandlerMap[DDLCheckDatabaseCollation].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckDecimalTypeColumnOffline(t *testing.T) {
	rule := RuleHandlerMap[DDLCheckDecimalTypeColumn].Rule
	runSingleRuleInspectCase(rule, t, "create table: check decimal type column", DefaultMysqlInspectOffline(),
		`
CREATE TABLE exist_db.not_exist_tb_4 (v1 float(10));
`,
		newTestResult().addResult(DDLCheckDecimalTypeColumn),
	)
	runSingleRuleInspectCase(rule, t, "alter table: check decimal type column", DefaultMysqlInspectOffline(),
		`
ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 FLOAT ( 10, 0 );
`,
		newTestResult().addResult(DDLCheckDecimalTypeColumn),
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
			RuleHandlerMap[DDLCheckColumnBlobNotice].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckColumnBlobNotice))
	}

	for desc, sql := range map[string]string{
		`(1)create table`: `CREATE TABLE t1(id INT);`,
		`(1)alter table`:  `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 VARCHAR(100);`,
		`(2)alter table`:  `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 SET('male', 'female');`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckColumnBlobNotice].Rule,
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
			RuleHandlerMap[DDLCheckColumnSetNitice].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckColumnSetNitice))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 INT;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB;`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckColumnSetNitice].Rule,
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
			RuleHandlerMap[DDLCheckColumnEnumNotice].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckColumnEnumNotice))
	}

	for desc, sql := range map[string]string{
		`create table`:   `CREATE TABLE t1(id INT);`,
		`(1)alter table`: `ALTER TABLE exist_db.exist_tb_1 ADD COLUMN v3 BLOB;`,
		`(2)alter table`: `ALTER TABLE exist_db.exist_tb_1 CHANGE COLUMN v1 v1 BLOB`,
		`(3)alter table`: `ALTER TABLE exist_db.exist_tb_1 MODIFY COLUMN v1 BLOB;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckColumnEnumNotice].Rule,
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
			RuleHandlerMap[DDLCheckUniqueIndex].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckUniqueIndex))
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
			RuleHandlerMap[DDLCheckUniqueIndex].Rule,
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
			RuleHandlerMap[DMLWhereExistNull].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DMLWhereExistNull))
	}

	for desc, sql := range map[string]string{
		`select table`: `SELECT * FROM exist_db.exist_tb_1 WHERE id = 1;`,
		`update table`: `UPDATE exist_db.exist_tb_1 SET id = 10 WHERE id = 1;`,
		`delete table`: `DELETE FROM exist_db.exist_tb_1 WHERE id = 1;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLWhereExistNull].Rule,
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
			RuleHandlerMap[DMLWhereExistNull].Rule,
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
			RuleHandlerMap[DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DMLCheckNeedlessFunc))
	}

	for desc, sql := range map[string]string{
		`(1)INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(1, sha1('aaa'), sha1('bbb'));`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckNeedlessFunc].Rule,
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
			RuleHandlerMap[DMLCheckNeedlessFunc].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DMLCheckNeedlessFunc))
	}

	for desc, sql := range map[string]string{
		`[fp]INSERT`: `INSERT INTO exist_db.exist_tb_1 VALUES(?, sha1(?), sha1(?));`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckNeedlessFunc].Rule,
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
			RuleHandlerMap[DDLCheckDatabaseSuffix].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckDatabaseSuffix))
	}

	for desc, sql := range map[string]string{
		`(0)create database`: `CREATE DATABASE app_service_db;`,
		`(1)create database`: `CREATE DATABASE APP_SERVICE_DB;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckDatabaseSuffix].Rule,
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
			RuleHandlerMap[DDLCheckTransactionIsolationLevel].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckTransactionIsolationLevel))
	}

	for desc, sql := range map[string]string{
		`(1)transaction isolation should not notice`: `SET TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(2)transaction isolation should not notice`: `SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(3)transaction isolation should not notice`: `SET GLOBAL TRANSACTION ISOLATION LEVEL READ COMMITTED;`,
		`(4)transaction isolation should not notice`: `SET GLOBAL TRANSACTION READ ONLY;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckTransactionIsolationLevel].Rule,
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
		runSingleRuleInspectCase(RuleHandlerMap[DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(DMLCheckFuzzySearch))
	}

	for _, sql := range []string{
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a___';`,

		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a%';`,
		`UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE 'a___';`,

		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a%';`,
		`DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE 'a____';`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DMLCheckFuzzySearch].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}
}

func TestCheckFuzzySearch_FPOffline(t *testing.T) {
	for desc, sql := range map[string]string{
		`[fp] "select" unable to check fuzzy search`: `SELECT * FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
		`[fp] "update" unable to check fuzzy search`: `UPDATE exist_db.exist_tb_1 SET id = 1 WHERE v1 LIKE ?;`,
		`[fp] "delete" unable to check fuzzy search`: `DELETE FROM exist_db.exist_tb_1 WHERE v1 LIKE ?;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckFuzzySearch].Rule,
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
			RuleHandlerMap[DDLCheckTablePartition].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckTablePartition))
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
			RuleHandlerMap[DDLCheckTablePartition].Rule,
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
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)
		inspector.Ctx = NewContext(parent.Ctx)
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT * FROM exist_db.exist_tb_1 JOIN exist_db.exist_tb_2 ON exist_db.exist_tb_1.id = exist_db.exist_tb_2.id 
JOIN exist_db.exist_tb_3 ON exist_db.exist_tb_2.id = exist_db.exist_tb_3.id
JOIN exist_db.exist_tb_4 ON exist_db.exist_tb_3.id = exist_db.exist_tb_4.id
`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			inspector,
			sql,
			newTestResult().addResult(DMLCheckNumberOfJoinTables))
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
			RuleHandlerMap[DMLCheckNumberOfJoinTables].Rule,
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
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
			newTestResult(),
		)
		inspector.Ctx = NewContext(parent.Ctx)
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
			RuleHandlerMap[DMLCheckNumberOfJoinTables].Rule,
			t,
			desc,
			inspector,
			sql,
			newTestResult().addResult(DMLCheckNumberOfJoinTables))
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
			RuleHandlerMap[DMLCheckNumberOfJoinTables].Rule,
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
			RuleHandlerMap[DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DMLCheckIfAfterUnionDistinct))
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT 1, 2 UNION ALL SELECT 'a', 'b';`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckIfAfterUnionDistinct].Rule,
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
			RuleHandlerMap[DMLCheckIfAfterUnionDistinct].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DMLCheckIfAfterUnionDistinct))
	}

	for desc, sql := range map[string]string{
		`select table should error`: `
SELECT ?, ? UNION ALL SELECT ?, ?;`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DMLCheckIfAfterUnionDistinct].Rule,
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
			RuleHandlerMap[DDLCheckIsExistLimitOffset].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckIsExistLimitOffset))
	}

	for desc, sql := range map[string]string{
		`select table should not error`: `
SELECT * FROM exist_db.exist_tb_1 LIMIT 5`,
	} {
		runSingleRuleInspectCase(
			RuleHandlerMap[DDLCheckIsExistLimitOffset].Rule,
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
			RuleHandlerMap[DDLCheckOBjectNameUseCN].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult().addResult(DDLCheckOBjectNameUseCN))
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
			RuleHandlerMap[DDLCheckOBjectNameUseCN].Rule,
			t,
			desc,
			DefaultMysqlInspectOffline(),
			sql,
			newTestResult())
	}
}

func TestCheckIndexOption_ShouldNot_QueryDBOffline(t *testing.T) {
	runSingleRuleInspectCase(
		RuleHandlerMap[DDLCheckIndexOption].Rule,
		t,
		`(1)index on new db new column`,
		DefaultMysqlInspectOffline(),
		`CREATE TABLE t1(id int, name varchar(100), INDEX idx_name(name))`,
		newTestResult())

	runSingleRuleInspectCase(
		RuleHandlerMap[DDLCheckIndexOption].Rule,
		t,
		`(2)index on new db new column`,
		DefaultMysqlInspectOffline(),
		`CREATE TABLE t1(id int, name varchar(100));
ALTER TABLE t1 ADD INDEX idx_name(name);
`,
		newTestResult(), newTestResult())

	runSingleRuleInspectCase(
		RuleHandlerMap[DDLCheckIndexOption].Rule,
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
		`create table t1(id int, primary key(id))`,
		`alter table exist_db.exist_tb_2 Add primary key(id)`,
		`alter table exist_db.exist_tb_2 Add primary key PK_EXIST_TB_2(id)`} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckPKName].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
	}

	for _, sql := range []string{
		`create table t1(id int, primary key wrongPK(id))`,
		`alter table exist_db.exist_tb_2 Add primary key wrongPK(id)`} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckPKName].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(DDLCheckPKName))
	}
}

func Test_PerfectParseOffline(t *testing.T) {
	runSingleRuleInspectCase(RuleHandlerMap[DMLCheckWhereIsInvalid].Rule, t, "", DefaultMysqlInspectOffline(), `
SELECT * FROM exist_db.exist_tb_1;
OPTIMIZE TABLE exist_db.exist_tb_1;
SELECT * FROM exist_db.exist_tb_2;
`, newTestResult().addResult(DMLCheckWhereIsInvalid),
		newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持"),
		newTestResult().addResult(DMLCheckWhereIsInvalid))
}

func Test_DDLCheckCreateViewOffline(t *testing.T) {
	for _, sql := range []string{
		`create view v as select * from t1`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateView].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().addResult(DDLCheckCreateView))
	}

	for _, sql := range []string{
		`create table t1(id int)`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateView].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult())
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
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持").addResult(DDLCheckCreateTrigger))
	}

	for _, sql := range []string{
		`CREATE my_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATEmy_trigger BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE trigger_1 BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER BEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
		`CREATE TRIGGER my_trigger BEEEFORE INSERT ON t1 FOR EACH ROW insert into t2(id, c1) values(1, '2');`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateTrigger].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持"))
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
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持").addResult(DDLCheckCreateFunction))
	}

	for _, sql := range []string{
		`create function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`create123 function_hello (s CHAR(20)) returns CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE hello_function (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
		`CREATE DEFINER='sqle_op'@'localhost' hello (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!');`,
	} {
		runSingleRuleInspectCase(RuleHandlerMap[DDLCheckCreateFunction].Rule, t, "", DefaultMysqlInspectOffline(), sql, newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持"))
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
			RuleHandlerMap[DDLCheckCreateProcedure].Rule, t, "",
			DefaultMysqlInspectOffline(), sql,
			newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持").
				addResult(DDLCheckCreateProcedure))
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
			RuleHandlerMap[DDLCheckCreateProcedure].Rule, t, "",
			DefaultMysqlInspectOffline(), sql,
			newTestResult().add(model.RuleLevelError, "语法错误或者解析器不支持"))
	}
}
