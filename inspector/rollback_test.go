package inspector

import (
	"github.com/stretchr/testify/assert"
	"sqle/model"
	"testing"
)

func runrollbackCase(t *testing.T, desc string, i *Inspector, sql string, results ...string) {
	stmts, err := parseSql(i.Instance.DbType, sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	for n, stmt := range stmts {
		i.SqlArray = append(i.SqlArray, &model.CommitSql{
			Number: n + 1,
			Sql:    stmt.Text(),
		})
	}
	rollbackSqls, err := i.GenerateRollbackSql()
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	assert.Equal(t, results, rollbackSqls, desc)
}

func TestAlterTableRollbackSql(t *testing.T) {
	runrollbackCase(t, "drop column need add", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP COLUMN v1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD COLUMN `v1` varchar(255) DEFAULT NULL;",
	)

	runrollbackCase(t, "add column need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD COLUMN v3 varchar(255) DEFAULT NULL;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP COLUMN `v3`;",
	)

	runrollbackCase(t, "rename table", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
RENAME AS exist_tb_2;`,
		"ALTER TABLE `exist_db`.`exist_tb_2`"+"\n"+
			"RENAME AS `exist_db`.`exist_tb_1`;",
	)

	runrollbackCase(t, "change column need change column", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
CHANGE COLUMN v1 v3 varchar(30) NOT NULL;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"CHANGE COLUMN `v3` `v1` varchar(255) DEFAULT NULL;",
	)

	runrollbackCase(t, "alter column need alter column(1_1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v1 DROP DEFAULT;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v1` SET DEFAULT NULL;",
	)

	runrollbackCase(t, "alter column need alter column(1_2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v1 SET DEFAULT "test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v1` SET DEFAULT NULL;",
	)

	runrollbackCase(t, "alter column need alter column(2_1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v2 SET DEFAULT "test";`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ALTER COLUMN `v2` DROP DEFAULT;",
	)

	runrollbackCase(t, "alter column need alter column(2_2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ALTER COLUMN v2 DROP DEFAULT;`,
		[]string{}...,
	)

	runrollbackCase(t, "alter column add index need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD INDEX v1(v1);`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP INDEX `v1`;",
	)

	runrollbackCase(t, "alter column drop index need add(1)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP INDEX v1;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD INDEX `v1` (`v1`);",
	)

	runrollbackCase(t, "alter column drop index need add(2)", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP INDEX v2;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD UNIQUE INDEX `v2` (`v1`,`v2`);",
	)

	runrollbackCase(t, "alter column add primary key need drop", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
ADD PRIMARY KEY (id) USING BTREE;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"DROP PRIMARY KEY;",
	)

	runrollbackCase(t, "alter column drop primary key need add", DefaultMysqlInspect(),
		`ALTER TABLE exist_db.exist_tb_1
DROP PRIMARY KEY;`,
		"ALTER TABLE `exist_db`.`exist_tb_1`"+"\n"+
			"ADD PRIMARY KEY (`id`) USING BTREE;",
	)
}

func TestInsertRollbackSql(t *testing.T) {
	runrollbackCase(t, "insert into: need delete(1)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 (id,v1,v2) value (10,"v1","v2"),(11,"v1","v2");`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';"+
			"\nDELETE FROM `exist_db`.`exist_tb_1` WHERE id = '11';\n",
	)
	runrollbackCase(t, "insert into: need delete(2)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 value (10,"v1","v2"),(11,"v1","v2");`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';\n"+
			"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '11';\n",
	)
	runrollbackCase(t, "insert into: need delete(3)", DefaultMysqlInspect(),
		`INSERT INTO exist_db.exist_tb_1 set id=10,v1="v1",v2="v2";`,
		"DELETE FROM `exist_db`.`exist_tb_1` WHERE id = '10';\n",
	)
}
