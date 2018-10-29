package inspector

import (
	"fmt"
	"sqle/model"
	"testing"
)

var DefaultRules = model.GetRuleMapFromAllArray(model.DefaultRules)

type testResult struct {
	Results *InspectResults
	rules   map[string]model.Rule
}

func newTestResult() *testResult {
	return &testResult{
		Results: newInspectResults(),
		rules:   DefaultRules,
	}
}

func (t *testResult) addResult(ruleName string, args ...interface{}) *testResult {
	rule, ok := t.rules[ruleName]
	if !ok {
		return t
	}
	t.Results.add(rule.Level, ruleName, args...)
	return t
}

func (t *testResult) level() string {
	return t.Results.level()
}

func (t *testResult) message() string {
	return t.Results.message()
}

func DefaultMysqlInspect() *Inspector {
	return &Inspector{
		Results: newInspectResults(),
		Rules:   DefaultRules,
		Db: model.Instance{
			DbType: "mysql",
		},
		SqlArray:      []*model.CommitSql{},
		currentSchema: "exist_db",
		allSchema:     map[string]struct{}{"exist_db": struct{}{}},
		schemaHasLoad: true,
		allTable: map[string]map[string]struct{}{
			"exist_db": map[string]struct{}{
				"exist_tb_1": struct{}{},
				"exist_tb_2": struct{}{},
			}},
	}
}

func runInspectCase(t *testing.T, desc string, i *Inspector, sql string, results ...*testResult) {
	stmts, err := parseSql(i.Db.DbType, sql)
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	for n, stmt := range stmts {
		i.SqlArray = append(i.SqlArray, &model.CommitSql{
			Number: n + 1,
			Sql:    stmt.Text(),
		})
	}
	_, err = i.Inspect()
	if err != nil {
		t.Errorf("%s test failled, error: %v\n", desc, err)
		return
	}
	if len(i.SqlArray) != len(results) {
		t.Errorf("%s test failled, error: result is unknow\n", desc)
		return
	}
	for n, sql := range i.SqlArray {
		result := results[n]
		if sql.InspectLevel != result.level() || sql.InspectResult != result.message() {
			t.Errorf("%s test failled, \n\nsql:\n %s\n\nexpect level: %s\nexpect result:\n%s\n\nactual level: %s\nactual result:\n%s\n",
				desc, sql.Sql, result.level(), result.message(), sql.InspectLevel, sql.InspectResult)
		}
	}
}

func TestInspector_Inspect_Message(t *testing.T) {
	runInspectCase(t, "check inspect message", DefaultMysqlInspect(),
		"use no_exist_db",
		&testResult{
			Results: &InspectResults{
				[]*InspectResult{&InspectResult{
					Level:   "error",
					Message: "schema no_exist_db 不存在",
				}},
			},
		},
	)
}

func TestInspector_Inspect_UseDatabaseStmt(t *testing.T) {
	runInspectCase(t, "use_database: ok", DefaultMysqlInspect(),
		"use exist_db",
		newTestResult(),
	)
	runInspectCase(t, "use_database: database not exist", DefaultMysqlInspect(),
		"use no_exist_db",
		newTestResult().addResult(model.SCHEMA_NOT_EXIST, "no_exist_db"),
	)
}

func TestInspector_Inspect_SelectStmt(t *testing.T) {
	runInspectCase(t, "select_from: ok", DefaultMysqlInspect(),
		"select * from exist_db.exist_tb_1",
		newTestResult(),
	)
	runInspectCase(t, "select_from: schema not exist", DefaultMysqlInspect(),
		"select * from not_exist_db.exist_tb_1, not_exist_db.exist_tb_2",
		newTestResult().addResult(model.SCHEMA_NOT_EXIST, "not_exist_db"),
	)
	runInspectCase(t, "select_from: table not exist", DefaultMysqlInspect(),
		"select * from exist_db.exist_tb_1, exist_db.exist_tb_3",
		newTestResult().addResult(model.TABLE_NOT_EXIST, "exist_db.exist_tb_3"),
	)
}

func TestInspector_Inspect_CreateTableStmt(t *testing.T) {
	runInspectCase(t, "create_table: ok", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: need \"if not exists\"", DefaultMysqlInspect(),
		`
CREATE TABLE exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult().addResult(model.DDL_CREATE_TABLE_NOT_EXIST),
	)

	runInspectCase(t, "create_table: schema not exist", DefaultMysqlInspect(),
		`
CREATE TABLE if not exists not_exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult().addResult(model.SCHEMA_NOT_EXIST, "not_exist_db"),
	)

	length64 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcd"
	length65 := "aaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeeeffffffffffabcde"

	runInspectCase(t, "create_table: table length <= 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;`, length64),
		newTestResult(),
	)

	runInspectCase(t, "create_table: table length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.%s (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;`, length65),
		newTestResult().addResult(model.DDL_CHECK_TABLE_NAME_LENGTH, length65),
	)

	runInspectCase(t, "create_table: columns length > 64", DefaultMysqlInspect(),
		fmt.Sprintf(`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
%s varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;`, length65),
		newTestResult().addResult(model.DDL_CHECK_COLUMNS_NAME_LENGTH, length65),
	)

	runInspectCase(t, "create_table: primary key exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: primary key not exist", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult().addResult(model.DDL_CHECK_PRIMARY_KEY_EXIST),
	)

	runInspectCase(t, "create_table: primary key not auto increment", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) unsigned NOT NULL,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult().addResult(model.DDL_CHECK_PRIMARY_KEY_TYPE),
	)

	runInspectCase(t, "create_table: primary key not unsigned", DefaultMysqlInspect(),
		`
CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
a1.id int(10) NOT NULL AUTO_INCREMENT,
v1 varchar(255) DEFAULT NULL,
v2 varchar(255) DEFAULT NULL,
PRIMARY KEY (id)
)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
`,
		newTestResult().addResult(model.DDL_CHECK_PRIMARY_KEY_TYPE),
	)

	//	runInspectCase(t, "create_table: disable varchar(max)", DefaultMysqlInspect(),
	//		`
	//CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	//a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
	//v1 varchar(65535) DEFAULT NULL,
	//v2 varchar(255) DEFAULT NULL,
	//PRIMARY KEY (id)
	//)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
	//`,
	//		newInspectResults().add(model.RULE_LEVEL_ERROR, model.DDL_DISABLE_VARCHAR_MAX),
	//	)

	runInspectCase(t, "create_table: check char(20)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
	v1 char(20) DEFAULT NULL,
	v2 varchar(255) DEFAULT NULL,
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
	`,
		newTestResult(),
	)

	runInspectCase(t, "create_table: check char(21)", DefaultMysqlInspect(),
		`
	CREATE TABLE  if not exists exist_db.not_exist_tb_1 (
	a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
	v1 char(21) DEFAULT NULL,
	v2 varchar(255) DEFAULT NULL,
	PRIMARY KEY (id)
	)ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
	`,
		newTestResult().addResult(model.DDL_CHECK_TYPE_CHAR_LENGTH),
	)
}
