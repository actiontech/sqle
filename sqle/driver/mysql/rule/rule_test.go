package rule

import (
	"testing"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/stretchr/testify/assert"
)

func TestInspectResults(t *testing.T) {
	results := driverV2.NewAuditResults()
	handler := RuleHandlerMap[DDLCheckPKWithoutIfNotExists]
	results.Add(handler.Rule.Level, handler.Rule.Name, plocale.Bundle.LocalizeAll(plocale.DDLCheckPKWithoutIfNotExistsMessage))
	assert.Equal(t, driverV2.RuleLevelError, results.Level())
	assert.Equal(t, "[error]新建表建议加入 IF NOT EXISTS，保证重复执行不报错", results.Message())

	results.Add(driverV2.RuleLevelError, "", plocale.Bundle.LocalizeAllWithArgs(plocale.TableNotExistMessage, "not_exist_tb"))
	assert.Equal(t, driverV2.RuleLevelError, results.Level())
	assert.Equal(t,
		`[error]新建表建议加入 IF NOT EXISTS，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.Message())

	results2 := driverV2.NewAuditResults()
	results2.Add(results.Level(), "", i18nPkg.ConvertStr2I18nAsDefaultLang(results.Message()))
	results2.Add(driverV2.RuleLevelNotice, "", i18nPkg.ConvertStr2I18nAsDefaultLang("test"))
	assert.Equal(t, driverV2.RuleLevelError, results2.Level())
	assert.Equal(t,
		`[error]新建表建议加入 IF NOT EXISTS，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test`, results2.Message())

	results3 := driverV2.NewAuditResults()
	results3.Add(results2.Level(), "", i18nPkg.ConvertStr2I18nAsDefaultLang(results2.Message()))
	results3.Add(driverV2.RuleLevelNotice, "", i18nPkg.ConvertStr2I18nAsDefaultLang("[osc]test"))
	assert.Equal(t, driverV2.RuleLevelError, results3.Level())
	assert.Equal(t,
		`[error]新建表建议加入 IF NOT EXISTS，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test
[osc]test`, results3.Message())

	results4 := driverV2.NewAuditResults()
	results4.Add(driverV2.RuleLevelNotice, "", i18nPkg.ConvertStr2I18nAsDefaultLang("[notice]test"))
	results4.Add(driverV2.RuleLevelError, "", i18nPkg.ConvertStr2I18nAsDefaultLang("[osc]test"))
	assert.Equal(t, driverV2.RuleLevelError, results4.Level())
	assert.Equal(t,
		`[osc]test
[notice]test`, results4.Message())

	results5 := driverV2.NewAuditResults()
	results5.Add(driverV2.RuleLevelWarn, "", i18nPkg.ConvertStr2I18nAsDefaultLang("[warn]test"))
	results5.Add(driverV2.RuleLevelNotice, "", i18nPkg.ConvertStr2I18nAsDefaultLang("[osc]test"))
	assert.Equal(t, driverV2.RuleLevelWarn, results5.Level())
	assert.Equal(t,
		`[warn]test
[osc]test`, results5.Message())
}

func TestCheckRedundantIndex(t *testing.T) {
	indexs1 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t3",
			Column: []string{"c3"},
		},
	}
	repeat, redundancy := checkRedundantIndex(indexs1)
	assert.Equal(t, repeat, []string{}, "indexs1,repeat")
	assert.Equal(t, len(redundancy), 1, "indexs1,redundancy")
	assert.Equal(t, redundancy["t2(c1)"], "t1(c1,c2,c3)", "indexs1,redundancy")

	indexs2 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t3",
			Column: []string{"c1", "c2"},
		},
	}
	repeat, redundancy = checkRedundantIndex(indexs2)
	assert.Equal(t, repeat, []string{}, "indexs2,repeat")
	assert.Equal(t, len(redundancy), 2, "indexs2,redundancy")
	assert.Equal(t, redundancy["t2(c1)"], "t1(c1,c2,c3)", "indexs2,redundancy")
	assert.Equal(t, redundancy["t3(c1,c2)"], "t1(c1,c2,c3)", "indexs2,redundancy")

	indexs3 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t3",
			Column: []string{"c1"},
		},
	}
	repeat, redundancy = checkRedundantIndex(indexs3)
	assert.Equal(t, repeat, []string{"t2(c1)"}, "indexs3,repeat")
	assert.Equal(t, len(redundancy), 1, "indexs3,redundancy")
	assert.Equal(t, redundancy["t3(c1)"], "t1(c1,c2,c3)", "indexs3,redundancy")

}

func TestCheckAlterTableRedundantIndex(t *testing.T) {
	newIndexs1 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
	}
	tableIndexs1 := []index{
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t3",
			Column: []string{"c1"},
		},
	}
	repeat, redundancy := checkAlterTableRedundantIndex(newIndexs1, tableIndexs1)
	assert.Equal(t, repeat, []string{}, "indexs1,repeat")
	assert.Equal(t, len(redundancy), 1, "indexs1,redundancy")
	assert.Equal(t, redundancy["t3(c1)"], "t1(c1,c2,c3)", "indexs1,redundancy")

	newIndexs2 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
	}
	tableIndexs2 := []index{
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t3",
			Column: []string{"c1"},
		},
	}
	repeat, redundancy = checkAlterTableRedundantIndex(newIndexs2, tableIndexs2)
	assert.Equal(t, repeat, []string{"t1(c1,c2,c3)"}, "indexs2,repeat")
	assert.Equal(t, len(redundancy), 1, "indexs2,redundancy")
	assert.Equal(t, redundancy["t3(c1)"], "t1(c1,c2,c3)", "indexs2,redundancy")

	newIndexs3 := []index{
		{
			Name:   "t1",
			Column: []string{"c1", "c2", "c3"},
		},
	}
	tableIndexs3 := []index{
		{
			Name:   "t2",
			Column: []string{"c1"},
		},
		{
			Name:   "t4",
			Column: []string{"c1", "c2", "c3"},
		},
		{
			Name:   "t3",
			Column: []string{"c1"},
		},
	}
	repeat, redundancy = checkAlterTableRedundantIndex(newIndexs3, tableIndexs3)
	assert.Equal(t, repeat, []string{"t1(c1,c2,c3)"}, "indexs3,repeat")
	assert.Equal(t, len(redundancy), 0, "indexs3,redundancy")

}
