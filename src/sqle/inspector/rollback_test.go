package inspector

//import (
//	"github.com/pingcap/tidb/ast"
//	"sqle/storage"
//	"testing"
//)
//
//func TestAlterTableRollbackSql(t *testing.T) {
//
//	baseCreateQuery := `
//CREATE TABLE a1 (
//id int(10) unsigned NOT NULL AUTO_INCREMENT,
//v1 varchar(255) DEFAULT NULL,
//v2 varchar(255) DEFAULT NULL,
//PRIMARY KEY (id)
//)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
//`
//
//	type testCase struct {
//		desc   string
//		create string
//		alter  string
//		output string
//	}
//
//	runTest := func(tc *testCase) {
//		t1, err := parseSql(storage.DB_TYPE_MYSQL, tc.create)
//		if err != nil {
//			t.Errorf("%s test failled, error: %v", tc.desc, err)
//			return
//		}
//		t2, err := parseSql(storage.DB_TYPE_MYSQL, tc.alter)
//		if err != nil {
//			t.Errorf("%s test failled, error: %v", tc.desc, err)
//			return
//		}
//		t11, ok := t1[0].(*ast.CreateTableStmt)
//		if !ok {
//			t.Errorf("%s test failled, error: \"create\" query is invalid", tc.desc)
//			return
//		}
//		t22, ok := t2[0].(*ast.AlterTableStmt)
//		if !ok {
//			t.Errorf("%s test failled, error: \"alter\" query is invalid", tc.desc)
//			return
//		}
//		output, _ := alterTableRollbackSql(t11, t22)
//		if output != tc.output {
//			t.Errorf("case: \"%s\" test failled\nactual output:\n%s\nexpect output:\n%s\n", tc.desc, output, tc.output)
//		}
//	}
//
//	runTest(&testCase{
//		desc: "drop column need add",
//		create: `
//CREATE TABLE a1 (
//id int(10) unsigned NOT NULL AUTO_INCREMENT,
//v1 varchar(255) DEFAULT NULL,
//v2 varchar(255) DEFAULT NULL,
//PRIMARY KEY (id)
//)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
//`,
//		alter: `ALTER TABLE t1.a1
//DROP COLUMN v1;`,
//		output: `ALTER TABLE t1.a1
//ADD COLUMN v1 varchar(255) DEFAULT NULL;`,
//	})
//
//	runTest(&testCase{
//		desc:   "add column need drop",
//		create: baseCreateQuery,
//		alter: `ALTER TABLE t1.a1
//ADD COLUMN v3 varchar(255) DEFAULT NULL;`,
//		output: `ALTER TABLE t1.a1
//DROP COLUMN v3;`,
//	})
//
//	runTest(&testCase{
//		desc:   "rename table",
//		create: baseCreateQuery,
//		alter: `ALTER TABLE t1.a1
//RENAME AS a2;`,
//		output: `ALTER TABLE t1.a2
//RENAME AS a1;`,
//	})
//}
