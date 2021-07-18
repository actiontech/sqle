package inspector

import (
	"github.com/pingcap/parser"
	"testing"

	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/pingcap/parser/ast"
	"github.com/stretchr/testify/assert"
)

func TestRemoveArrayRepeat(t *testing.T) {
	input := []string{"a", "b", "c", "c", "a"}
	expect := []string{"a", "b", "c"}
	actual := RemoveArrayRepeat(input)
	assert.Equal(t, expect, actual)
}

var columnOptionsForTest = []*ast.ColumnOption{
	&ast.ColumnOption{
		Tp: ast.ColumnOptionAutoIncrement,
	},
	&ast.ColumnOption{
		Tp: ast.ColumnOptionPrimaryKey,
	},
	&ast.ColumnOption{
		Tp: ast.ColumnOptionNotNull,
	},
}

func TestIsAllInOptions(t *testing.T) {
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNotNull), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNull), false)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, IsAllInOptions(columnOptionsForTest, ast.ColumnOptionNull, ast.ColumnOptionPrimaryKey), false)
}

func TestHasOneInOptions(t *testing.T) {
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNotNull), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNull), false)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionAutoIncrement, ast.ColumnOptionPrimaryKey), true)
	assert.Equal(t, HasOneInOptions(columnOptionsForTest, ast.ColumnOptionNull, ast.ColumnOptionPrimaryKey), true)
}

func TestReplaceSchemaName(t *testing.T) {
	input := "alter table `db1`.tb1 drop column a1"
	output := "alter table `tb1` drop column a1"
	assert.Equal(t, replaceTableName(input, "db1", "tb1"), output)
}

func TestGetDuplicate(t *testing.T) {
	assert.Equal(t, []string{}, getDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"2"}, getDuplicate([]string{"1", "2", "2"}))
	assert.Equal(t, []string{"2", "3"}, getDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestRemoveDuplicate(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "2", "3"}))
	assert.Equal(t, []string{"1", "2", "3"}, removeDuplicate([]string{"1", "2", "2", "3", "3", "3"}))
}

func TestInspectResults(t *testing.T) {
	results := newInspectResults()
	handler := RuleHandlerMap[DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST]
	results.add(handler.Rule.Level, handler.Message)
	assert.Equal(t, "error", results.level())
	assert.Equal(t, "[error]新建表必须加入if not exists create，保证重复执行不报错", results.message())

	results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG, "not_exist_tb")
	assert.Equal(t, "error", results.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在`, results.message())

	results2 := newInspectResults()
	results2.add(results.level(), results.message())
	results2.add("notice", "test")
	assert.Equal(t, "error", results2.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test`, results2.message())

	results3 := newInspectResults()
	results3.add(results2.level(), results2.message())
	results3.add("notice", "[osc]test")
	assert.Equal(t, "error", results3.level())
	assert.Equal(t,
		`[error]新建表必须加入if not exists create，保证重复执行不报错
[error]表 not_exist_tb 不存在
[notice]test
[osc]test`, results3.message())

	results4 := newInspectResults()
	results4.add("notice", "[notice]test")
	results4.add("error", "[osc]test")
	assert.Equal(t, "error", results4.level())
	assert.Equal(t,
		`[notice]test
[osc]test`, results4.message())

	results5 := newInspectResults()
	results5.add("warn", "[warn]test")
	results5.add("notice", "[osc]test")
	assert.Equal(t, "warn", results5.level())
	assert.Equal(t,
		`[warn]test
[osc]test`, results5.message())
}

type VisitorTestCase struct {
	visitor      ast.Visitor
	inputSQLText string

	expectSQLs []string
}

func TestCapitalizeProcessor(t *testing.T) {
	for _, c := range []VisitorTestCase{
		{
			visitor: &CapitalizeProcessor{
				capitalizeTableName:      true,
				capitalizeTableAliasName: true,
				capitalizeDatabaseName:   true},
			inputSQLText: `INSERT INTO db1.t1 (id) VALUES (1);
DELETE FROM db1.t1 WHERE id=1;
UPDATE db1.t1 SET id=2 WHERE id=1;
SELECT * FROM db1.t1 AS t1_alias;
CREATE DATABASE db1;
DROP DATABASE db1;
ALTER DATABASE db1 COLLATE = utf8mb4_bin;
`,
			expectSQLs: []string{
				`INSERT INTO DB1.T1 (id) VALUES (1)`,
				`DELETE FROM DB1.T1 WHERE id=1`,
				`UPDATE DB1.T1 SET id=2 WHERE id=1`,
				`SELECT * FROM DB1.T1 AS T1_ALIAS`,
				`CREATE DATABASE DB1`,
				`DROP DATABASE DB1`,
				`ALTER DATABASE DB1 COLLATE = utf8mb4_bin`,
			},
		},

		{
			visitor: &CapitalizeProcessor{},
			inputSQLText: `INSERT INTO db1.t1 (id) VALUES (1);
DELETE FROM t1 WHERE id=1;
UPDATE t1 SET id=2 WHERE id=1;
SELECT * FROM t1 AS t1_alias;
CREATE DATABASE db1;
DROP DATABASE db1;
ALTER DATABASE db1 COLLATE = utf8mb4_bin;
`,
			expectSQLs: []string{
				`INSERT INTO db1.t1 (id) VALUES (1)`,
				`DELETE FROM t1 WHERE id=1`,
				`UPDATE t1 SET id=2 WHERE id=1`,
				`SELECT * FROM t1 AS t1_alias`,
				`CREATE DATABASE db1`,
				`DROP DATABASE db1`,
				`ALTER DATABASE db1 COLLATE = utf8mb4_bin`,
			},
		},

		{
			visitor: &CapitalizeProcessor{
				capitalizeTableName:      true,
				capitalizeTableAliasName: true,
				capitalizeDatabaseName:   false},
			inputSQLText: `INSERT INTO db1.t1 (id) VALUES (1);
DELETE FROM db1.t1 WHERE id=1;
UPDATE db1.t1 SET id=2 WHERE id=1;
SELECT * FROM db1.t1 AS t1_alias;
`,
			expectSQLs: []string{
				`INSERT INTO db1.T1 (id) VALUES (1)`,
				`DELETE FROM db1.T1 WHERE id=1`,
				`UPDATE db1.T1 SET id=2 WHERE id=1`,
				`SELECT * FROM db1.T1 AS T1_ALIAS`,
			},
		},
	} {
		stmts, _, err := parser.New().PerfectParse(c.inputSQLText, "", "")
		assert.NoError(t, err)

		for i, stmt := range stmts {
			assert.Panics(t, func() { _ = stmt.(*ast.UnparsedStmt) })
			stmt.Accept(c.visitor)
			restoredSQL, err := restoreToSqlWithFlag(0, stmt)
			assert.NoError(t, err)
			assert.Equal(t, c.expectSQLs[i], restoredSQL)
		}
	}
}
