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
	assert.Equal(t, "[error]新建表必须加入if not exists create，保证重复执行不报错", results.Message())

	results.Add(driver.RuleLevelError, "表 %s 不存在", "not_exist_tb")
	assert.Equal(t, driver.RuleLevelError, results.Level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.Message())

	results2 := driver.NewInspectResults()
	results2.Add(results.Level(), results.Message())
	results2.Add(driver.RuleLevelNotice, "test")
	assert.Equal(t, driver.RuleLevelError, results2.Level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test`, results2.Message())

	results3 := driver.NewInspectResults()
	results3.Add(results2.Level(), results2.Message())
	results3.Add(driver.RuleLevelNotice, "[osc]test")
	assert.Equal(t, driver.RuleLevelError, results3.Level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test
[osc]test`, results3.Message())

	results4 := driver.NewInspectResults()
	results4.Add(driver.RuleLevelNotice, "[notice]test")
	results4.Add(driver.RuleLevelError, "[osc]test")
	assert.Equal(t, driver.RuleLevelError, results4.Level())
	assert.Equal(t,
		`[notice]test
[osc]test`, results4.Message())

	results5 := driver.NewInspectResults()
	results5.Add(driver.RuleLevelWarn, "[warn]test")
	results5.Add(driver.RuleLevelNotice, "[osc]test")
	assert.Equal(t, driver.RuleLevelWarn, results5.Level())
	assert.Equal(t,
		`[warn]test
[osc]test`, results5.Message())
}
