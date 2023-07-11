package session

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func Test_LoadSchemaLowerCaseTableNamesOpen(t *testing.T) {
	context := &Context{
		schemas:       map[string]*SchemaInfo{},
		schemaHasLoad: false,
		sysVars: map[string]string{
			"lower_case_table_names": "1",
		},
	}

	context.loadSchemas([]string{"a", "B", "aAa"})
	for schema := range context.schemas {
		for _, s := range schema {
			assert.True(t, unicode.IsLower(s))
		}
	}
}

func Test_LoadTablesLowerCaseTableNamesOpen(t *testing.T) {
	context := &Context{
		schemas: map[string]*SchemaInfo{
			"exist_db": &SchemaInfo{},
		},
		schemaHasLoad: true,
		sysVars: map[string]string{
			"lower_case_table_names": "1",
		},
	}

	context.loadTables("EXIST_DB", []string{"a", "B", "aAa"})
	for table := range context.schemas["exist_db"].Tables {
		for _, s := range table {
			assert.True(t, unicode.IsLower(s))
		}
	}
	context.schemas["exist_db"] = &SchemaInfo{}
	context.loadTables("exist_DB", []string{"a", "B", "aAa"})
	for table := range context.schemas["exist_db"].Tables {
		for _, s := range table {
			assert.True(t, unicode.IsLower(s))
		}
	}
	context.schemas["exist_db"] = &SchemaInfo{}
	context.loadTables("exist_db", []string{"a", "B", "aAa"})
	for table := range context.schemas["exist_db"].Tables {
		for _, s := range table {
			assert.True(t, unicode.IsLower(s))
		}
	}
}

func TestParseObMysqlCreateTableSql(t *testing.T) {
	args := []struct {
		sql     string
		wantSql string
	}{
		{
			`CREATE TABLE tb01 (
  ID int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 ID',
  NAME varchar(20) NOT NULL DEFAULT '' COMMENT '名字',
  PRIMARY KEY (ID)
) ENGINE=InnoDB DEFAULT anyOpts=anything COMMENT='测试表'`,
			`CREATE TABLE tb01 (
  ID int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 ID',
  NAME varchar(20) NOT NULL DEFAULT '' COMMENT '名字',
  PRIMARY KEY (ID)
)`,
		},
		{
			`CREATE TABLE tb01 (
  ID int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 ID'
) ENGINE=InnoDB DEFAULT anyOpts=anything COMMENT='测试表'`,
			`CREATE TABLE tb01 (
  ID int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增 ID'
)`,
		},
	}

	c := &Context{}
	for _, arg := range args {
		t.Run("test parse OB MySQL create table", func(t *testing.T) {
			createTbStmt, err := c.parseObMysqlCreateTableSql(arg.sql)
			assert.NoError(t, err)
			assert.Equal(t, arg.wantSql, createTbStmt.Text())
		})
	}
}
