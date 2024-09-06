package mysql

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"golang.org/x/text/language"

	"testing"
)

func TestCheckInvalidOffline(t *testing.T) {
	testCheckInvalidCreateTableOffline(t)
	testCheckInvalidAlterTableOffline(t)
	testCheckInvalidCreateIndexOffline(t)
	testCheckInvalidInsertOffline(t)
}

func testCheckInvalidCreateTableOffline(t *testing.T) {
	runEmptyRuleInspectCase(t, "column name can't duplicated. f", DefaultMysqlInspectOffline(),
		`create table t (a int,b int, a int)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateColumnsMessage), "a"))
	runEmptyRuleInspectCase(t, "column name can't duplicated. t", DefaultMysqlInspectOffline(),
		`create table t (a int,b int)`,
		newTestResult())

	runEmptyRuleInspectCase(t, "pk can only be set once. f1", DefaultMysqlInspectOffline(),
		`create table t (a int primary key,b int primary key)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.MultiPrimaryKeyMessage)))
	runEmptyRuleInspectCase(t, "pk can only be set once. f2", DefaultMysqlInspectOffline(),
		"create table t (a int primary key,b int, PRIMARY KEY (`b`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.MultiPrimaryKeyMessage)))
	runEmptyRuleInspectCase(t, "pk can only be set once. f3", DefaultMysqlInspectOffline(),
		"create table t (a int primary key,b int, PRIMARY KEY (`a`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.MultiPrimaryKeyMessage)))
	runEmptyRuleInspectCase(t, "pk can only be set once. f4", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int, PRIMARY KEY (`a`), PRIMARY KEY (`b`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.MultiPrimaryKeyMessage)))
	runEmptyRuleInspectCase(t, "pk can only be set once. t", DefaultMysqlInspectOffline(),
		`create table t (a int,b int primary key)`,
		newTestResult())

	runEmptyRuleInspectCase(t, "index name can't be duplicated. f1", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , KEY `a` (`a`), KEY `a` (`b`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateIndexesMessage), "a"))
	runEmptyRuleInspectCase(t, "index name can't be duplicated. f2", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , unique `a`(`a`), KEY `a`(`b`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateIndexesMessage), "a"))
	runEmptyRuleInspectCase(t, "index name can't be duplicated. t", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , unique `a`(`a`), KEY `b`(`b`))",
		newTestResult())

	runEmptyRuleInspectCase(t, "index column must exist. f1", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , unique `a`(`c`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.KeyedColumnNotExistMessage), "c"))
	runEmptyRuleInspectCase(t, "index column must exist. f2", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , unique `a`(`a`,`c`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.KeyedColumnNotExistMessage), "c"))
	runEmptyRuleInspectCase(t, "index column must exist. t", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , unique `a`(`b`))",
		newTestResult())

	runEmptyRuleInspectCase(t, "index column can't duplicated. f", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , index `idx`(`a`,`a`))",
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateIndexedColumnMessage), "idx", "a"))
	runEmptyRuleInspectCase(t, "index column can't duplicated. t", DefaultMysqlInspectOffline(),
		"create table t (a int ,b int , index `idx`(`a`,`b`))",
		newTestResult())

}

func testCheckInvalidAlterTableOffline(t *testing.T) {
	runEmptyRuleInspectCase(t, "add pk, ok can only be set once. f1", DefaultMysqlInspectOffline(),
		`alter table t add (a int primary key, b int primary key)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PrimaryKeyExistMessage)))
	runEmptyRuleInspectCase(t, "add pk, ok can only be set once. f2", DefaultMysqlInspectOffline(),
		`alter table t add primary key (a), add primary key (b)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.PrimaryKeyExistMessage)))
	runEmptyRuleInspectCase(t, "add pk, ok can only be set once. t1", DefaultMysqlInspectOffline(),
		`alter table t add primary key (a), add index (b)`,
		newTestResult())
	runEmptyRuleInspectCase(t, "add pk, ok can only be set once. t2", DefaultMysqlInspectOffline(),
		`alter table t add (a int primary key , b int)`,
		newTestResult())

	runEmptyRuleInspectCase(t, "index column can't duplicated. f1", DefaultMysqlInspectOffline(),
		`alter table t add index b(a,a)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateIndexedColumnMessage), "b", "a"))
	runEmptyRuleInspectCase(t, "index column can't duplicated. f2", DefaultMysqlInspectOffline(),
		`alter table t add primary key a(a,a)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicatePrimaryKeyedColumnMessage), "a"))
	runEmptyRuleInspectCase(t, "index column can't duplicated. t", DefaultMysqlInspectOffline(),
		`alter table t add index a(a,b), add index b(b,c)`,
		newTestResult())

}

func testCheckInvalidCreateIndexOffline(t *testing.T) {
	runEmptyRuleInspectCase(t, "index column name can't be duplicated. f", DefaultMysqlInspectOffline(),
		`create index idx on t (a,a)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateIndexedColumnMessage), "idx", "a"))
	runEmptyRuleInspectCase(t, "index column name can't be duplicated. t", DefaultMysqlInspectOffline(),
		`create index idx on t (a,b)`,
		newTestResult())
}

func testCheckInvalidInsertOffline(t *testing.T) {
	runEmptyRuleInspectCase(t, "index column can't be duplicated. f1", DefaultMysqlInspectOffline(),
		`insert into t (a,a) value (1,1)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateColumnsMessage), "a"))
	runEmptyRuleInspectCase(t, "index column can't be duplicated. f2", DefaultMysqlInspectOffline(),
		`insert into t set a=1, a=1`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.DuplicateColumnsMessage), "a"))
	runEmptyRuleInspectCase(t, "index column can't be duplicated. t1", DefaultMysqlInspectOffline(),
		`insert into t set a=1, b=1`,
		newTestResult())
	runEmptyRuleInspectCase(t, "index column can't be duplicated. t2", DefaultMysqlInspectOffline(),
		`insert into t (a,b) value (1,1)`,
		newTestResult())

	runEmptyRuleInspectCase(t, "value length must match column length. f", DefaultMysqlInspectOffline(),
		`insert into t (a,b) value (1,1,1)`,
		newTestResult().add(driverV2.RuleLevelError, "", plocale.ShouldLocalizeMsgByLang(language.Chinese, plocale.ColumnsValuesNotMatchMessage)))
	runEmptyRuleInspectCase(t, "value length must match column length. t1", DefaultMysqlInspectOffline(),
		`insert into t (a,b) value (1,1)`,
		newTestResult())
	runEmptyRuleInspectCase(t, "value length must match column length. t1", DefaultMysqlInspectOffline(),
		`insert into t values (1,1)`,
		newTestResult())

}
