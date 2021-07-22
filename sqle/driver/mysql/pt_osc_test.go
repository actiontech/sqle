package mysql

import (
	"fmt"
	"testing"

	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/stretchr/testify/assert"
)

func TestPTOSC(t *testing.T) {
	expect := "pt-online-schema-change D=exist_db,t=%s --alter='%s' --host=127.0.0.1 --user=root --port=3306 --ask-pass --print --execute"

	runOSCCase(t, "Add column not null no default",
		"alter table exist_tb_1 Add column v3 varchar(255) NOT NULL;",
		PTOSCAvoidNoDefaultValueOnNotNullColumn)

	runOSCCase(t, "not pk and unique key",
		"alter table exist_tb_3 Add column v3 varchar(255);",
		PTOSCNoUniqueIndexOrPrimaryKey)

	runOSCCase(t, "rename table",
		"alter table exist_tb_1 rename as not_exist_tb_1;",
		PTOSCAvoidRenameTable)

	runOSCCase(t, "Add unique index",
		"alter table exist_tb_1 Add unique index u_1 (v1) ",
		PTOSCAvoidUniqueIndex)

	runOSCCase(t, "Add column ok",
		"alter table exist_tb_1 Add column v3 varchar(255);",
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v3` varchar(255)"))

	runOSCCase(t, "drop foreign key",
		"alter table exist_tb_2 drop foreign key `pk_test_1`",
		fmt.Sprintf(expect, "exist_tb_2", "DROP FOREIGN KEY `pk_test_1`"))

	runOSCCase(t, "Add multi column(1)",
		`alter table exist_tb_1 Add column(v4 varchar(255),v5 varchar(255) not null default "1")`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN (`v4` varchar(255), `v5` varchar(255) NOT NULL DEFAULT \"1\")"))

	runOSCCase(t, "Add multi column(2)",
		`alter table exist_tb_1 Add column v4 varchar(255),Add column v5 varchar(255) not null default "1"`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v4` varchar(255),ADD COLUMN `v5` varchar(255) NOT NULL DEFAULT \"1\""))
}

func runOSCCase(t *testing.T, desc string, sql, expect string) {
	i := DefaultMysqlInspect()
	i.config.DDLOSCMinSize = 0
	stmt, err := parseOneSql(model.DBTypeMySQL, sql)
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
