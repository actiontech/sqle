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
					"exist_tb_4": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          100,
						OriginalTable: getTestCreateTableStmt4(),
					},

					// used for test case problem
					"EXIST_TB_5": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt5(),
					},
					"exist_tb_6": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt6(),
					},
					"exist_tb_7": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt7(),
					},
					"exist_tb_8": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt8(),
					},
					"exist_tb_9": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt9(),
					},
					"exist_tb_10": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt10(),
					},
					"exist_tb_11": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt11(),
					},
				},
			},
			"myisam_utf8_db": {
				DefaultEngine:    "MyISAM",
				engineLoad:       true,
				DefaultCharacter: "utf8",
				characterLoad:    true,
				Tables: map[string]*TableInfo{
					"exist_tb_1": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt1(),
					},
				},
			},
		},
		historySqlInfo: &HistorySQLInfo{},
	}
}

func NewMockContextForTestLowerCaseTableNameOpen(e *executor.Executor) *Context {
	return &Context{
		e:             e,
		currentSchema: "exist_db",
		schemaHasLoad: true,
		executionPlan: map[string][]*executor.ExplainRecord{},
		sysVars: map[string]string{
			"lower_case_table_names": "1",
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
				},
			},
		},
		historySqlInfo: &HistorySQLInfo{},
	}
}

func NewMockContextForTestLowerCaseTableNameClose(e *executor.Executor) *Context {
	return &Context{
		e:             e,
		currentSchema: "exist_db",
		schemaHasLoad: true,
		executionPlan: map[string][]*executor.ExplainRecord{},
		sysVars: map[string]string{
			"lower_case_table_names": "0",
		},
		schemas: map[string]*SchemaInfo{
			"exist_db_1": {
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
				},
			},
			"EXIST_DB_2": {
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
					"EXIST_TB_1": {
						sizeLoad:      true,
						isLoad:        true,
						Size:          1,
						OriginalTable: getTestCreateTableStmt2(),
					},
				},
			},
		},
		historySqlInfo: &HistorySQLInfo{},
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

func getTestCreateTableStmt4() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_4 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
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

func getTestCreateTableStmt5() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.EXIST_TB_5 (
id bigint unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT "unit test",
v1 varchar(255) NOT NULL COMMENT "unit test",
v2 varchar(255) NOT NULL COMMENT "unit test"
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT="uint test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt6() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_1 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) NOT NULL DEFAULT "v1" COMMENT "unit test",
v2 varchar(255) COMMENT "unit test",
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1),
UNIQUE KEY uniq_1 (v1,v2),
KEY idx_100 (v2,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt7() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_7 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) character SET utf8,
v2 varchar(255) COMMENT "unit test" COLLATE utf8_bin,
v3 varchar(255),
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1),
UNIQUE KEY uniq_1 (v1,v2),
KEY idx_100 (v2,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt8() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_8 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 varchar(255) character SET utf8mb4 COLLATE utf8_bin,
v2 varchar(255) character SET utf8mb4,
v3 varchar(255),
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1),
UNIQUE KEY uniq_1 (v1,v2),
KEY idx_100 (v2,v1)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt9() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_9 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 int,
v2 varchar(255) character SET utf8mb4,
v3 int,
v4 int,
v5 int,
PRIMARY KEY (id) USING BTREE,
KEY idx_1 (v1,v2,v3, v4),
UNIQUE KEY uniq_1 (v2,v3),
KEY idx_100 (v3)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt10() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_10 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
v1 int,
v2 varchar(255) character SET utf8mb4,
v3 TEXT,
v4 JSON,
v5 int,
PRIMARY KEY (id) USING BTREE
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}

func getTestCreateTableStmt11() *ast.CreateTableStmt {
	baseCreateQuery := `
CREATE TABLE exist_db.exist_tb_11 (
id bigint(10) unsigned NOT NULL AUTO_INCREMENT COMMENT "unit test",
create_time datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
upgrade_time timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
year_time year(4) NOT NULL DEFAULT '2020',
data_time date NOT NULL DEFAULT '2020-01-01 00:00:00',
data_time2 TIME NOT NULL DEFAULT '12:00:00', 
PRIMARY KEY (id) USING BTREE
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT="unit test";
`
	node, err := util.ParseOneSql(baseCreateQuery)
	if err != nil {
		panic(err)
	}
	stmt, _ := node.(*ast.CreateTableStmt)
	return stmt
}
