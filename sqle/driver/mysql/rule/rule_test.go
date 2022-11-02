package rule

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/stretchr/testify/assert"
)

func TestInspectResults(t *testing.T) {
	results := driver.NewInspectResults()
	handler := RuleHandlerMap[DDLCheckPKWithoutIfNotExists]
	results.Add(handler.Rule.Level, handler.Message)
	assert.Equal(t, driver.RuleLevelError, results.Level())
	assert.Equal(t, "[error]新建表必须加入 if not exists，保证重复执行不报错", results.Message())

	results.Add(driver.RuleLevelError, "表 %s 不存在", "not_exist_tb")
	assert.Equal(t, driver.RuleLevelError, results.Level())
	assert.Equal(t,
		`[error]新建表必须加入 if not exists，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.Message())

	results2 := driver.NewInspectResults()
	results2.Add(results.Level(), results.Message())
	results2.Add(driver.RuleLevelNotice, "test")
	assert.Equal(t, driver.RuleLevelError, results2.Level())
	assert.Equal(t,
		`[error]新建表必须加入 if not exists，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test`, results2.Message())

	results3 := driver.NewInspectResults()
	results3.Add(results2.Level(), results2.Message())
	results3.Add(driver.RuleLevelNotice, "[osc]test")
	assert.Equal(t, driver.RuleLevelError, results3.Level())
	assert.Equal(t,
		`[error]新建表必须加入 if not exists，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test
[osc]test`, results3.Message())

	results4 := driver.NewInspectResults()
	results4.Add(driver.RuleLevelNotice, "[notice]test")
	results4.Add(driver.RuleLevelError, "[osc]test")
	assert.Equal(t, driver.RuleLevelError, results4.Level())
	assert.Equal(t,
		`[osc]test
[notice]test`, results4.Message())

	results5 := driver.NewInspectResults()
	results5.Add(driver.RuleLevelWarn, "[warn]test")
	results5.Add(driver.RuleLevelNotice, "[osc]test")
	assert.Equal(t, driver.RuleLevelWarn, results5.Level())
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
