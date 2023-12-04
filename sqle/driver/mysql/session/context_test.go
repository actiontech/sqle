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
			"exist_db": {},
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
