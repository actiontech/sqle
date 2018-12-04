package inspector

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPTOSC(t *testing.T) {
	expect := "pt-online-schema-change D=exist_db,t=%s \\" + "\n" +
		"--alter=\"%s\" \\" + "\n" +
		"--host=127.0.0.1 \\" + "\n" +
		"--user=root \\" + "\n" +
		"--port=3306 \\" + "\n" +
		"--ask-pass \\" + "\n" +
		"--print \\" + "\n" +
		"--execute"

	runOSCCase(t, "add column not null no default",
		"alter table exist_tb_1 add column v3 varchar(255) NOT NULL;",
		OSC_AVOID_ADD_NOT_NULL_NO_DEFAULT_COLUMN)

	runOSCCase(t, "not pk and unique key",
		"alter table exist_tb_3 add column v3 varchar(255);",
		OSC_NO_UNIQUE_INDEX_AND_PRIMARY_KEY)

	runOSCCase(t, "rename table",
		"alter table exist_tb_1 rename as not_exist_tb_1;",
		OSC_AVOID_RENAME_TABLE)

	runOSCCase(t, "add unique index",
		"alter table exist_tb_1 add unique index u_1 (v1) ",
		OSC_AVOID_ADD_UNIQUE_INDEX)

	runOSCCase(t, "add column ok",
		"alter table exist_tb_1 add column v3 varchar(255);",
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v3` varchar(255)"))

	runOSCCase(t, "drop foreign key",
		"alter table exist_tb_2 drop foreign key `pk_test_1`",
		fmt.Sprintf(expect, "exist_tb_2", "DROP FOREIGN KEY `pk_test_1`"))

	runOSCCase(t, "add multi column(1)",
		`alter table exist_tb_1 add column(v4 varchar(255),v5 varchar(255) not null default "1")`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN (`v4` varchar(255), `v5` varchar(255) NOT NULL DEFAULT \"1\")"))

	runOSCCase(t, "add multi column(2)",
		`alter table exist_tb_1 add column v4 varchar(255),add column v5 varchar(255) not null default "1"`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v4` varchar(255),ADD COLUMN `v5` varchar(255) NOT NULL DEFAULT \"1\""))
}

func runOSCCase(t *testing.T, desc string, sql, expect string) {
	UpdateConfig(CONFIG_DDL_OSC_SIZE_LIMIT, "0")
	i := DefaultMysqlInspect()
	stmt, err := parseOneSql(i.Task.Instance.DbType, sql)
	if err != nil {
		t.Error(err)
		return
	}

	actual, err := i.generateOSCCommandLine(stmt)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, expect, actual, desc)
}
