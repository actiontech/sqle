package session

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
)

// NewMockContext creates a new mock context for unit test.
func NewMockContext(e *executor.Executor) *Context {
	return &Context{
		e:             e,
		currentSchema: "exist_db",
		schemaHasLoad: true,
		executionPlan: map[string][]*executor.ExplainRecord{},
		sysVars: map[string]string{
			"lower_case_table_names": "0",
		},
		schemas: map[string]*SchemaInfo{
			"exist_db": {
				DefaultEngine:    "InnoDB",
				engineLoad:       true,
				DefaultCharacter: "utf8mb4",
				characterLoad:    true,
				Tables: map[string]*TableInfo{
					"exist_tb_1": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt1(),
					},
					"exist_tb_2": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt2(),
					},
					"exist_tb_3": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt3(),
					},
				},
			},
		},
	}
}

func getTestCreateTableStmt1() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_1 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "v1" COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1),
UNIQUE KEY uniq_1 (v1,v2)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt2() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_2 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
user_id bigint unsigned NOT NULL COMMENT "unit test",
UNIQUE KEY uniq_1(id),
CONSTRAINT pk_test_1 FOREIGN KEY (user_id) REFERENCES exist_db.exist_tb_1 (id) ON DELETE NO ACTION
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt3() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_3 (
id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
v3 int COMMENT "unit test"
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="uint test"
PARTITION BY LIST(v3) (
PARTITION p1 VALUES IN(1, 2, 3),
PARTITION p2 VALUES IN(4, 5, 6),
PARTITION p3 VALUES IN(7, 8, 9)
);
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}
