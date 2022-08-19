package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
)

func TestContext(t *testing.T) {
	handler := rule.RuleHandlerMap[rule.DDLCheckAlterTableNeedMerge]
	delete(rule.RuleHandlerMap, rule.DDLCheckAlterTableNeedMerge)
	defer func() {
		rule.RuleHandlerMap[rule.DDLCheckAlterTableNeedMerge] = handler
	}()

	runDefaultRulesInspectCase(t, "rename table and drop column: table not exists", DefaultMysqlInspect(),
		`
use exist_db;
create table if not exists not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
alter table not_exist_tb_1 rename as not_exist_tb_2;
alter table not_exist_tb_2 drop column v1;
alter table not_exist_tb_1 drop column v1;
`,
		newTestResult(),
		newTestResult(),
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, TableNotExistMessage, "exist_db.not_exist_tb_1"),
	)

	runDefaultRulesInspectCase(t, "drop column twice: column not exists(1)", DefaultMysqlInspect(),
		`
	use exist_db;
	alter table exist_tb_1 drop column v1;
	alter table exist_tb_1 drop column v1;
	`,
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v1"),
	)
	runDefaultRulesInspectCase(t, "drop column twice: column not exists(2)", DefaultMysqlInspect(),
		`
	use exist_db;
	create table if not exists not_exist_tb_1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
	v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
	create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
	alter table not_exist_tb_1 drop column v1;
	alter table not_exist_tb_1 drop column v1;
	`,
		newTestResult(),
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "change and drop column: column not exists", DefaultMysqlInspect(),
		`
	use exist_db;
	alter table exist_tb_1 change column v1 v11 varchar(255) DEFAULT "v11" COMMENT "uint test";
	alter table exist_tb_1 drop column v1;
	`,
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v1"),
	)

	runDefaultRulesInspectCase(t, "Add column twice: column exists", DefaultMysqlInspect(),
		`
	use exist_db;
	alter table exist_tb_1 add column v3 varchar(255) DEFAULT "v3" COMMENT "uint test";
	alter table exist_tb_1 add column v3 varchar(255) DEFAULT "v3" COMMENT "uint test";
	`,
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, ColumnExistMessage, "v3"),
	)

	runDefaultRulesInspectCase(t, "drop index twice: index not exists", DefaultMysqlInspect(),
		`
	use exist_db;
	alter table exist_tb_1 drop index idx_1;
	alter table exist_tb_1 drop index idx_1;
	`,
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, IndexNotExistMessage, "idx_1"),
	)
	runDefaultRulesInspectCase(t, "drop index, rename index: index not exists", DefaultMysqlInspect(),
		`
	use exist_db;
	alter table exist_tb_1 rename index idx_1 to idx_2;
	alter table exist_tb_1 drop index idx_1;
	`,
		newTestResult(),
		newTestResult(),
		newTestResult().add(driver.RuleLevelError, IndexNotExistMessage, "idx_1"),
	)
}

func TestParentContext(t *testing.T) {
	handler := rule.RuleHandlerMap[rule.DDLCheckAlterTableNeedMerge]
	delete(rule.RuleHandlerMap, rule.DDLCheckAlterTableNeedMerge)
	// It's trick :),
	// elegant method: unit test support MySQL.
	delete(rule.RuleHandlerMap, rule.DDLCheckTableDBEngine)
	delete(rule.RuleHandlerMap, rule.DDLCheckTableCharacterSet)
	defer func() {
		rule.RuleHandlerMap[rule.DDLCheckAlterTableNeedMerge] = handler
	}()

	inspect1 := DefaultMysqlInspect()
	runDefaultRulesInspectCase(t, "ddl 1: create table, ok", inspect1,
		`
use exist_db;
create table if not exists not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "unit test",
update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`,
		newTestResult(),
		newTestResult(),
	)

	inspect2 := DefaultMysqlInspect()
	inspect2.Ctx = session.NewContext(inspect1.Ctx)
	runDefaultRulesInspectCase(t, "ddl 2: drop column, ok", inspect2,
		`
alter table not_exist_tb_1 drop column v1;
`,
		newTestResult(),
	)

	inspect3 := DefaultMysqlInspect()
	inspect3.Ctx = session.NewContext(inspect2.Ctx)
	runDefaultRulesInspectCase(t, "ddl 3: drop column, column not exist", inspect3,
		`
alter table not_exist_tb_1 drop column v1;
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v1"),
	)

	inspect4 := DefaultMysqlInspect()
	inspect4.Ctx = session.NewContext(inspect2.Ctx)
	runDefaultRulesInspectCase(t, "ddl 4: add column, ok", inspect4,
		`
alter table not_exist_tb_1 add column v3 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
		newTestResult(),
	)

	inspect5 := DefaultMysqlInspect()
	inspect5.Ctx = session.NewContext(inspect4.Ctx)
	runDefaultRulesInspectCase(t, "dml 1: insert, column not exist", inspect5,
		`
insert into not_exist_tb_1 (id,v1,v2) values (1,"1","1");
`,
		newTestResult().add(driver.RuleLevelError, ColumnNotExistMessage, "v1"),
	)

	inspect6 := DefaultMysqlInspect()
	inspect6.Ctx = session.NewContext(inspect4.Ctx)
	runDefaultRulesInspectCase(t, "dml 2: insert, ok", inspect6,
		`
insert into not_exist_tb_1 (id,v2,v3) values (1,"1","1");
`,
		newTestResult(),
	)
}

// TODO: Add more test for relation audit, like create a database and create a table in it.
