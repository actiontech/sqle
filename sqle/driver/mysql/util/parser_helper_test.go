package util

import (
	"testing"

	"github.com/pingcap/parser"
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
	{
		Tp: ast.ColumnOptionAutoIncrement,
	},
	{
		Tp: ast.ColumnOptionPrimaryKey,
	},
	{
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
			//nolint:forcetypeassert
			assert.Panics(t, func() { _ = stmt.(*ast.UnparsedStmt) })
			stmt.Accept(c.visitor)
			restoredSQL, err := restoreToSqlWithFlag(0, stmt)
			assert.NoError(t, err)
			assert.Equal(t, c.expectSQLs[i], restoredSQL)
		}
	}
}

type FpCase struct {
	input  string
	expect string
}

func TestFingerprint(t *testing.T) {
	cases := []FpCase{
		{
			input:  `update  tb1 set a = "2" where a = "3" and b = 4`,
			expect: "UPDATE `tb1` SET `a`=? WHERE `a`=? AND `b`=?",
		},
		{
			input:  "select * from tb1 where a in (select a from tb2 where b = 2) and c = 100",
			expect: "SELECT * FROM `tb1` WHERE `a` IN (SELECT `a` FROM `tb2` WHERE `b`=?) AND `c`=?",
		},
		{
			input:  "REPLACE INTO `tb1` (a, b, c, d, e) VALUES (1, 1, '小明', 'F', 99)",
			expect: "REPLACE INTO `tb1` (`a`,`b`,`c`,`d`,`e`) VALUES (?,?,?,?,?)",
		},
		{
			input:  "CREATE TABLE `tb1` SELECT * FROM `tb2` WHERE a=1",
			expect: "CREATE TABLE `tb1` AS SELECT * FROM `tb2` WHERE `a`=?",
		},
		{
			input:  "CREATE TABLE `tb1` AS SELECT * FROM `tb2` WHERE a=1",
			expect: "CREATE TABLE `tb1` AS SELECT * FROM `tb2` WHERE `a`=?",
		},
		// newline
		{
			input:  "CREATE TABLE `tb1` (\n    a BIGINT NOT NULL AUTO_INCREMENT,\n    b BIGINT NOT NULL,\n    c DOUBLE NOT NULL,\n    PRIMARY KEY (a)\n)",
			expect: "CREATE TABLE `tb1` (`a` BIGINT NOT NULL AUTO_INCREMENT,`b` BIGINT NOT NULL,`c` DOUBLE NOT NULL,PRIMARY KEY(`a`))",
		},

		// whitespace
		{
			input:  "select * from `tb1` where a='my_db'  and  b='test1'",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},

		// comment
		{
			input:  "create database database_x -- this is a comment ",
			expect: "CREATE DATABASE `database_x`",
		},
		{
			input:  "select * from tb1 where a='my_db' and b='test1'/*this is a comment*/",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},
		{
			input:  "select * from tb1 where a='my_db' and b='test1'# this is a comment",
			expect: "SELECT * FROM `tb1` WHERE `a`=? AND `b`=?",
		},
	}
	for _, c := range cases {
		testFingerprint(t, c.input, c.expect)
	}
}

func testFingerprint(t *testing.T, input, expect string) {
	actual, err := Fingerprint(input, true)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, expect, actual)
}
