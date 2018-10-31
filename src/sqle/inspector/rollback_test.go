package inspector

import (
	"github.com/pingcap/tidb/ast"
	"sqle/model"
	"testing"
)

func TestAlterTableRollbackSql(t *testing.T) {
	type testCase struct {
		desc   string
		alter  string
		output string
	}

	runTest := func(tc *testCase) {
		node, err := parseOneSql(model.DB_TYPE_MYSQL, tc.alter)
		if err != nil {
			t.Errorf("%s test failled, error: %v", tc.desc, err)
			return
		}

		stmt, ok := node.(*ast.AlterTableStmt)
		if !ok {
			t.Errorf("%s test failled, error: \"alter\" query is invalid", tc.desc)
			return
		}
		inspect := DefaultMysqlInspect()
		output, _ := inspect.alterTableRollbackSql(stmt)
		if output != tc.output {
			t.Errorf("case: \"%s\" test failled\n sql:\n%s\nactual output:\n%s\nexpect output:\n%s\n",
				tc.desc, tc.alter, output, tc.output)
		}
	}

	runTest(&testCase{
		desc: "drop column need add",

		alter: `ALTER TABLE exist_db.exist_tb_1
DROP COLUMN v1;`,

		output: `ALTER TABLE exist_db.exist_tb_1
ADD COLUMN v1 varchar(255) DEFAULT NULL;`,
	})

	runTest(&testCase{
		desc: "add column need drop",

		alter: `ALTER TABLE exist_db.exist_tb_1
ADD COLUMN v3 varchar(255) DEFAULT NULL;`,

		output: `ALTER TABLE exist_db.exist_tb_1
DROP COLUMN v3;`,
	})

	runTest(&testCase{
		desc: "rename table",

		alter: `ALTER TABLE exist_db.exist_tb_1
RENAME AS exist_tb_2;`,

		output: `ALTER TABLE exist_db.exist_tb_2
RENAME AS exist_db.exist_tb_1;`,
	})
}
