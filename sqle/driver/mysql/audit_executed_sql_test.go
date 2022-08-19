package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/sirupsen/logrus"
)

func NewSQLExecutedInspect(e *executor.Executor) *MysqlDriverImpl {
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
			DDLGhostMinSize:    -1,
			DMLRollbackMaxRows: 1000,
			isExecutedSQL:      true,
		},
	}
}

func TestAuditExecutedSQL(t *testing.T) {

	{ // 完全屏蔽的规则

		// DDLCheckAlterTableNeedMerge
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckAlterTableNeedMerge].Rule,
			t,
			"DDLCheckAlterTableNeedMerge",
			NewSQLExecutedInspect(nil),
			`
ALTER TABLE exist_db.exist_tb_1 Add column v5 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
ALTER TABLE exist_db.exist_tb_1 Add column v6 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test";
`,
			newTestResult(),
			newTestResult(),
		)

		// DDLCheckTableSize
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableSize].Rule,
			t,
			"DDLCheckTableSize",
			NewSQLExecutedInspect(nil),
			`drop table exist_db.exist_tb_4;`,
			newTestResult(),
		)

		// DDLCheckIndexesExistBeforeCreateConstraints
		runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexesExistBeforeCreateConstraints].Rule,
			t,
			"DDLCheckIndexesExistBeforeCreateConstraints",
			NewSQLExecutedInspect(nil),
			`alter table exist_db.exist_tb_3 Add unique uniq_test(v2);`,
			newTestResult(),
		)

	}

	{ // 部分屏蔽的规则 详见: https://github.com/actiontech/sqle/issues/716

		{ // 只检查建表语句

			// DDLCheckIndexedColumnWithBlob
			runDefaultRulesInspectCase(
				t,
				"DDLCheckIndexedColumnWithBlob",
				NewSQLExecutedInspect(nil),
				`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
b1 blob UNIQUE KEY COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
CREATE INDEX idx_1 ON exist_db.not_exist_tb_1(b1);
ALTER TABLE exist_db.not_exist_tb_1 ADD INDEX idx_2(b1);
ALTER TABLE exist_db.not_exist_tb_1 ADD COLUMN b2 blob UNIQUE KEY COMMENT "unit test";
ALTER TABLE exist_db.not_exist_tb_1 MODIFY COLUMN b1 blob UNIQUE KEY COMMENT "unit test";
`,
				newTestResult().addResult(rulepkg.DDLCheckIndexedColumnWithBlob).
					add(driver.RuleLevelWarn, "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
					add(driver.RuleLevelWarn, "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
				newTestResult(),
				newTestResult(),
				newTestResult(),
				newTestResult(),
			)

			// DDLCheckIndexTooMany
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexTooMany].Rule,
				t,
				"DDLCheckIndexTooMany",
				NewSQLExecutedInspect(nil),
				`
CREATE TABLE if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id),
INDEX idx_1 (v1,id),
INDEX idx_2 (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
AlTER TABLE exist_db.not_exist_tb_1 ADD INDEX idx_1(id), ADD INDEX idx_2(id), ADD INDEX idx_3(id);
`,
				newTestResult().addResult(rulepkg.DDLCheckIndexTooMany, "id", 2),
				newTestResult(),
			)

			// DDLCheckIndexCount
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexCount].Rule,
				t,
				"DDLCheckIndexCount",
				NewSQLExecutedInspect(nil),
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
AlTER TABLE exist_db.not_exist_tb_1 ADD INDEX idx_1(id), ADD INDEX idx_2(id), ADD INDEX idx_3(id), ADD INDEX idx_4(id), ADD INDEX idx_5(id), ADD INDEX idx_6 (id);
`,
				newTestResult().addResult(rulepkg.DDLCheckIndexCount, 5),
				newTestResult(),
			)

			// DDLCheckCompositeIndexMax
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckCompositeIndexMax].Rule,
				t,
				"DDLCheckCompositeIndexMax",
				NewSQLExecutedInspect(nil),
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
ALTER TABLE exist_db.not_exist_tb_1 ADD INDEX idx_1 (id,v1,v2,v3,v4,v5);
			`,
				newTestResult().addResult(rulepkg.DDLCheckCompositeIndexMax, 3),
				newTestResult(),
			)

			// DDLCheckPKProhibitAutoIncrement
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckPKProhibitAutoIncrement].Rule,
				t,
				"DDLCheckPKProhibitAutoIncrement",
				NewSQLExecutedInspect(nil),
				`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL AUTO_INCREMENT DEFAULT "unit test" COMMENT "unit test" ,
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
ALTER TABLE exist_db.not_exist_tb_1 modify COLUMN id BIGINT auto_increment;
				`,
				newTestResult().addResult(rulepkg.DDLCheckPKProhibitAutoIncrement),
				newTestResult(),
			)

			// DDLCheckPKWithoutAutoIncrement
			runDefaultRulesInspectCase(t,
				"DDLCheckPKWithoutAutoIncrement",
				NewSQLExecutedInspect(nil),
				`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint unsigned NOT NULL KEY DEFAULT "unit test" COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
ALTER TABLE exist_db.exist_tb_1 Add primary key(v1);
			`,
				newTestResult().addResult(rulepkg.DDLCheckPKWithoutAutoIncrement).add(driver.RuleLevelWarn, "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
					add(driver.RuleLevelWarn, "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
				newTestResult(),
			)

			// DDLCheckPKWithoutBigintUnsigned
			runDefaultRulesInspectCase(t,
				"DDLCheckPKWithoutBigintUnsigned",
				NewSQLExecutedInspect(nil),
				`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
id bigint NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
ALTER TABLE exist_db.exist_tb_1 Add primary key(v1);
			`,
				newTestResult().addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned).add(driver.RuleLevelWarn, "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
					add(driver.RuleLevelWarn, "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
				newTestResult(),
			)

			// TODO 这个规则不允许离线运行, 手动测试保证
			// DDLCheckRedundantIndex
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckRedundantIndex].Rule,
				t,
				"DDLCheckRedundantIndex",
				NewSQLExecutedInspect(nil),
				`
			CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
			id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
			v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
			v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
			PRIMARY KEY (id),
			INDEX idx_1 (v1,id),
			INDEX idx_2 (id)
			)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
			alter table exist_db.exist_tb_1 add index idx_t (v1);
						`,
				newTestResult().addResult(rulepkg.DDLCheckRedundantIndex, "存在重复索引:(id); "),
				newTestResult(),
			)

		}
	}
}
