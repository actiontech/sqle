package mysql

import (
	"context"
	"testing"

	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/stretchr/testify/assert"
)

func runRollbackCase(t *testing.T, desc string, i *Inspect, sql string, results string) {
	stmts, err := parseSql(model.DBTypeMySQL, sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}

	for _, stmt := range stmts {
		sql, _, err := i.GenRollbackSQL(context.TODO(), stmt.Text())
		_, err = parseSql(model.DBTypeMySQL, sql)
		assert.NoError(t, err)
		assert.Equal(t, results, sql, desc)
	}
}

func TestAlterTableRollbackSql(t *testing.T) {
	runRollbackCase(t, "drop column need add", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP COLUMN v1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD COLUMN `v1` varchar(255) NOT NULL DEFAULT \"v1\" COMMENT \"unit test\";",
	)

	runRollbackCase(t, "add column need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD COLUMN v3 varchar(255) DEFAULT NULL COMMENT "unit test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP COLUMN `v3`;",
	)

	runRollbackCase(t, "rename table", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
RENAME AS exist_tb_2;`,
		"ALTER TABLE `exist_db`.`exist_tb_2`"+"\n"+
			"RENAME AS `exist_db`.`exist_tb_1`;",
	)

	runRollbackCase(t, "change column need change column", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
CHANGE COLUMN v1 v3 varchar(30) NOT NULL COMMENT "unit test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"CHANGE COLUMN `v3` `v1` varchar(255) NOT NULL DEFAULT \"v1\" COMMENT \"unit test\";",
	)

	runRollbackCase(t, "modify column need modify column", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
MODIFY COLUMN v1 varchar(30) NOT NULL COMMENT "unit test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"MODIFY COLUMN `v1` varchar(255) NOT NULL DEFAULT \"v1\" COMMENT \"unit test\";",
	)

	runRollbackCase(t, "alter column need alter column(1_1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v1 DROP DEFAULT;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v1` SET DEFAULT \"v1\";",
	)

	runRollbackCase(t, "alter column need alter column(1_2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v1 SET DEFAULT "test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v1` SET DEFAULT \"v1\";",
	)

	runRollbackCase(t, "alter column need alter column(2_1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v2 SET DEFAULT "test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v2` DROP DEFAULT;",
	)

	runRollbackCase(t, "alter column need alter column(2_2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v2 DROP DEFAULT;`,
		"",
	)

	runRollbackCase(t, "alter column add index need drop(1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD INDEX idx_2(v1);`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP INDEX `idx_2`;",
	)

	runRollbackCase(t, "alter column add index need drop(2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD KEY idx_2(v1);`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP INDEX `idx_2`;",
	)

	runRollbackCase(t, "alter column drop index need add(1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP INDEX idx_1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD INDEX `idx_1` (`v1`);",
	)

	runRollbackCase(t, "alter column drop index need add(2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP INDEX uniq_1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD UNIQUE INDEX `uniq_1` (`v1`,`v2`);",
	)

	runRollbackCase(t, "alter column add unique index need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD UNIQUE INDEX uniq_2(v1,v2);`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP INDEX `uniq_2`;",
	)

	runRollbackCase(t, "alter column drop unique index need add(1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP INDEX uniq_1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD UNIQUE INDEX `uniq_1` (`v1`,`v2`);",
	)

	runRollbackCase(t, "alter column add primary key need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD PRIMARY KEY (id) USING BTREE;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP PRIMARY KEY;",
	)

	runRollbackCase(t, "alter column drop primary key need add", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP PRIMARY KEY;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD PRIMARY KEY (`id`) USING BTREE;",
	)
	runRollbackCase(t, "alter column add foreign key need drop", DefaultMysqlInspect(),
		"ALTER TABLE exist_db.exist_tb_1"+"\n"+
			"ADD FOREIGN KEY pk_1 (user_id) REFERENCES exist_db.exist_tb_2 (id) ON DELETE NO ACTION;",
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP FOREIGN KEY `pk_1`;",
	)
	runRollbackCase(t, "alter column drop foreign key need add", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_2
DROP FOREIGN KEY pk_test_1;`,
		"ALTER TABLE `exist_db`.`exist_tb_2`"+"\n"+
			"ADD CONSTRAINT `pk_test_1` FOREIGN KEY (`user_id`) REFERENCES `exist_db`.`exist_tb_1` (`id`) ON DELETE NO ACTION;",
	)
	runRollbackCase(t, "alter column rename index", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
RENAME INDEX old_name TO new_name;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"RENAME INDEX `new_name` TO `old_name`;",
	)
}

func TestInsertRollbackSql(t *testing.T) {
	runRollbackCase(t, "insert into: need delete(1)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 (id,v1,v2) value (10,"v1","v2"),(11,"v1","v2");`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';"+
			"\nDELETE FROM `exist_db`.`exist_tb_1` WHERE id = '11';\n",
	)
	runRollbackCase(t, "insert into: need delete(2)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 value (10,"v1","v2"),(11,"v1","v2");`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';\n"+
			"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '11';\n",
	)
	runRollbackCase(t, "insert into: need delete(3)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 set id=10,v1="v1",v2="v2";`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';\n",
	)
}
