package mysql

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"

	"github.com/sirupsen/logrus"
)

func NewSQLExecutedInspect(e *executor.Executor) *MysqlDriverImpl {
	log.Logger().SetLevel(logrus.ErrorLevel)
	return &MysqlDriverImpl{
		log: log.NewEntry(),
		inst: &driverV2.DSN{
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
		t.Run("DDLCheckAlterTableNeedMerge", func(t *testing.T) {
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
		})
		t.Run("DDLCheckAlterTableNeedMerge", func(t *testing.T) {
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
		})

		// DDLCheckTableSize
		t.Run("DDLCheckTableSize", func(t *testing.T) {
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckTableSize].Rule,
				t,
				"DDLCheckTableSize",
				NewSQLExecutedInspect(nil),
				`drop table exist_db.exist_tb_4;`,
				newTestResult(),
			)
		})

		// DDLCheckIndexesExistBeforeCreateConstraints
		t.Run("DDLCheckIndexesExistBeforeCreateConstraints", func(t *testing.T) {
			runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexesExistBeforeCreateConstraints].Rule,
				t,
				"DDLCheckIndexesExistBeforeCreateConstraints",
				NewSQLExecutedInspect(nil),
				`alter table exist_db.exist_tb_3 Add unique uniq_test(v2);`,
				newTestResult(),
			)
		})

	}

	{ // 部分屏蔽的规则 详见: https://github.com/actiontech/sqle/issues/716

		{ // 只检查建表语句

			// DDLCheckIndexedColumnWithBlob
			t.Run("DDLCheckIndexedColumnWithBlob", func(t *testing.T) {
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
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP").
						add(driverV2.RuleLevelWarn, "", "这些索引字段(b1)需要有非空约束"),
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "b1"),
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "b1"),
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "b2"),
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "b1"),
				)
			})

			// DDLCheckIndexTooMany
			t.Run("DDLCheckIndexTooMany", func(t *testing.T) {
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
			})

			// DDLCheckIndexCount
			t.Run("DDLCheckIndexCount", func(t *testing.T) {
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
			})

			// DDLCheckCompositeIndexMax
			t.Run("DDLCheckCompositeIndexMax", func(t *testing.T) {
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
			})

			// DDLCheckPKProhibitAutoIncrement
			t.Run("DDLCheckPKProhibitAutoIncrement", func(t *testing.T) {
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
			})

			// DDLCheckPKWithoutAutoIncrement
			t.Run("DDLCheckPKWithoutAutoIncrement", func(t *testing.T) {
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
					newTestResult().addResult(rulepkg.DDLCheckPKWithoutAutoIncrement).
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
					newTestResult(),
				)
			})

			// DDLCheckPKWithoutBigintUnsigned
			t.Run("DDLCheckPKWithoutBigintUnsigned", func(t *testing.T) {
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
					newTestResult().addResult(rulepkg.DDLCheckPKWithoutBigintUnsigned).
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP").
						add(driverV2.RuleLevelWarn, "", "建表DDL必须包含UPDATE_TIME字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"),
					newTestResult(),
				)
			})

			// TODO 这个规则不允许离线运行, 手动测试保证
			// DDLCheckRedundantIndex
			t.Run("DDLCheckRedundantIndex", func(t *testing.T) {
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
			})

			// DDLCheckIndexNotNullConstraint
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table index with not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE exist_db.not_exist_tb_1 (
			id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
			v1 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
			v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
			PRIMARY KEY (id),
			INDEX idx_1 (v1,id),
			INDEX idx_2 (id)
			)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`,
					newTestResult(),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table index without not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE exist_db.not_exist_tb_1 (
			id bigint unsigned AUTO_INCREMENT COMMENT "unit test",
			v1 varchar(255) DEFAULT "unit test" COMMENT "unit test",
			v2 varchar(255) NOT NULL DEFAULT "unit test" COMMENT "unit test",
			PRIMARY KEY (id),
			INDEX idx_1 (v1),
			INDEX idx_2 (id)
			)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";`,
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "id,v1"),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table unique with not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  username VARCHAR(50) NOT NULL,  
  email VARCHAR(100) NOT NULL,  
  UNIQUE KEY uq_username_email (username, email)  
); `,
					newTestResult(),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table unique without not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  username VARCHAR(50) NOT NULL,  
  email VARCHAR(100),  
  phone VARCHAR(100),  
  UNIQUE KEY uq_username_email (username, email, phone)  
); `,
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "email,phone"),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table unique key without not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  id INT UNIQUE KEY,  
  username VARCHAR(50),  
  email VARCHAR(100)  
); `,
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "id"),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table unique key with not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  id INT NOT NULL UNIQUE KEY,  
  username VARCHAR(50),  
  email VARCHAR(100)  
); `,
					newTestResult(),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table primary key without not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  username VARCHAR(50) NOT NULL,  
  email VARCHAR(100),  
  phone VARCHAR(100),  
  PRIMARY KEY uq_username_email (username, email, phone)  
); `,
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "email,phone"),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table primary key with not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  username VARCHAR(50) NOT NULL,  
  email VARCHAR(100) NOT NULL,  
  phone VARCHAR(100) NOT NULL,  
  PRIMARY KEY uq_username_email (username, email, phone)  
); `,
					newTestResult(),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table primary key without not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  id INT PRIMARY KEY,  
  username VARCHAR(50),  
  email VARCHAR(100)  
); `,
					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "id"),
				)
			})
			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
					t,
					"create table primary key with not null",
					NewSQLExecutedInspect(nil), `
CREATE TABLE users (  
  id INT NOT NULL PRIMARY KEY,  
  username VARCHAR(50),  
  email VARCHAR(100)  
); `,
					newTestResult(),
				)
			})

			// todo 添加alter table的单测，目前暂时手工测试
			//			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
			//					t,
			//					"alter table index with not null",
			//					NewSQLExecutedInspect(nil), `
			//alter table customers add index customer_id_index (customer_id);`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
			//					t,
			//					"alter table index without not null",
			//					NewSQLExecutedInspect(nil), `
			//alter table customers add index customer_id_index (customer_id);`,
			//					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "customer_id"),
			//				)
			//			})
			//			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
			//					t,
			//					"alter table unique with not null",
			//					NewSQLExecutedInspect(nil), `
			//ALTER TABLE table_name ADD UNIQUE KEY (column1, column2); `,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
			//					t,
			//					"alter table unique without not null",
			//					NewSQLExecutedInspect(nil), `
			//ALTER TABLE table_name ADD UNIQUE KEY (column1, column2); `,
			//					newTestResult().addResult(rulepkg.DDLCheckIndexNotNullConstraint, "column1", "column2"),
			//				)
			//			})
			// case: alter table add constraint constraint_name (col_name)

			// todo 添加create index的单测，目前暂时手工测试
			//			t.Run("DDLCheckIndexNotNullConstraint", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DDLCheckIndexNotNullConstraint].Rule,
			//					t,
			//					"create index with not null",
			//					NewSQLExecutedInspect(nil), `
			//CREATE INDEX part_of_name ON customer (name(10));`,
			//					newTestResult(),
			//				)
			//			})

			// DMLCheckSortColumnLength
			// todo 待支持线上单测，先手工测试
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"select order by",
			//					NewSQLExecutedInspect(nil), `
			//SELECT * FROM t1
			//  ORDER BY key_part1 DESC, key_part2 ASC;
			// `,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"select group by",
			//					NewSQLExecutedInspect(nil), `
			//				SELECT a, b, COUNT(c) AS t FROM test_table GROUP BY a,b
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"select distinct",
			//					NewSQLExecutedInspect(nil), `
			//SELECT DISTINCT c1, c2, c3 FROM t1
			//WHERE c1 > const;
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"UNION",
			//					NewSQLExecutedInspect(nil), `
			//SELECT 1, 2 UNION SELECT 'a', 'b';
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"UNION DISTINCT",
			//					NewSQLExecutedInspect(nil), `
			//SELECT 1, 2 UNION DISTINCT SELECT 'a', 'b';
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"UNION order by",
			//					NewSQLExecutedInspect(nil), `
			//SELECT name, age FROM table1 WHERE age > 20
			//UNION ALL
			//SELECT name, age FROM table2 WHERE age > 30
			//ORDER BY age DESC;
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"UNION ALL",
			//					NewSQLExecutedInspect(nil), `
			//SELECT 1, 2 UNION ALL SELECT 'a', 'b';
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"delete order by",
			//					NewSQLExecutedInspect(nil), `
			//DELETE FROM some WHERE user = 'f'
			//ORDER BY timestamp_column LIMIT 1;
			//				`,
			//					newTestResult(),
			//				)
			//			})
			//			t.Run("DMLCheckSortColumnLength", func(t *testing.T) {
			//				runSingleRuleInspectCase(rulepkg.RuleHandlerMap[rulepkg.DMLCheckSortColumnLength].Rule,
			//					t,
			//					"update order by",
			//					NewSQLExecutedInspect(nil), `
			//UPDATE t SET id = id + 1 ORDER BY id DESC;
			//				`,
			//					newTestResult(),
			//				)
			//			})
		}
	}
}
