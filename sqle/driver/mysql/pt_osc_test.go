package mysql

import (
	"fmt"
	"testing"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestPTOSC(t *testing.T) {
	expect := "[osc]pt-online-schema-change D=exist_db,t=%s --alter='%s' --host=127.0.0.1 --user=root --port=3306 --ask-pass --print --execute"

	runOSCCase(t, "add column not null no default",
		"alter table exist_tb_1 add column v3 varchar(255) NOT NULL;",
		plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PTOSCAvoidNoDefaultValueOnNotNullColumn))

	runOSCCase(t, "not pk and unique key",
		"alter table exist_tb_3 add column v3 varchar(255);",
		plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PTOSCNoUniqueIndexOrPrimaryKey))

	runOSCCase(t, "rename table",
		"alter table exist_tb_1 rename as not_exist_tb_1;",
		plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PTOSCAvoidRenameTable))

	runOSCCase(t, "add unique index",
		"alter table exist_tb_1 add unique index u_1 (v1) ",
		plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PTOSCAvoidUniqueIndex))

	runOSCCase(t, "add column ok",
		"alter table exist_tb_1 add column v3 varchar(255);",
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v3` varchar(255)"))

	runOSCCase(t, "drop foreign key",
		"alter table exist_tb_2 drop foreign key `pk_test_1`",
		fmt.Sprintf(expect, "exist_tb_2", "DROP FOREIGN KEY `pk_test_1`"))

	runOSCCase(t, "add multi column(1)",
		`alter table exist_tb_1 add column(v4 varchar(255),v5 varchar(255) not null default "1")`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN (`v4` varchar(255), `v5` varchar(255) NOT NULL DEFAULT \"1\")"))

	runOSCCase(t, "Add multi column(2)",
		`alter table exist_tb_1 Add column v4 varchar(255),Add column v5 varchar(255) not null default "1"`,
		fmt.Sprintf(expect, "exist_tb_1", "ADD COLUMN `v4` varchar(255),ADD COLUMN `v5` varchar(255) NOT NULL DEFAULT \"1\""))
}

func runOSCCase(t *testing.T, desc string, sql, expect string) {
	i := DefaultMysqlInspect()
	i.cnf.DDLOSCMinSize = 0
	stmt, err := util.ParseOneSql(sql)
	if err != nil {
		t.Error(err)
		return
	}
	actual, err := i.generateOSCCommandLine(stmt)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, expect, actual[locale.DefaultLang], desc)
}
