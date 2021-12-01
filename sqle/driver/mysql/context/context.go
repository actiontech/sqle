package context

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
)

type TableInfo struct {
	Size     float64
	sizeLoad bool

	// isLoad indicate whether TableInfo load from database or not.
	isLoad bool

	// OriginalTable save parser object from db by query "show create table ...";
	// using in inspect and generate rollback sql
	OriginalTable *ast.CreateTableStmt

	//
	MergedTable *ast.CreateTableStmt

	// save alter table parse object from input sql;
	AlterTables []*ast.AlterTableStmt
}

type SchemaInfo struct {
	DefaultEngine    string
	engineLoad       bool
	DefaultCharacter string
	characterLoad    bool
	DefaultCollation string
	collationLoad    bool
	Tables           map[string]*TableInfo
}

// Context is a database information cache.
//
// Is provides many methods to get database information.
//
// It do lazy load and cache the information if executor
// provided. Otherwise, it only return from cache.
type Context struct {
	e *executor.Executor

	// currentSchema will change after sql "use database"
	currentSchema string

	schemas map[string]*SchemaInfo
	// if schemas info has collected, set true
	schemaHasLoad bool

	// executionPlan store batch SQLs' execution plan during one inspect context.
	executionPlan map[string][]*executor.ExplainRecord

	// sysVars keep some MySQL global system variables during one inspect context.
	sysVars map[string]string
}

func NewContext(parent *Context) *Context {
	ctx := &Context{
		schemas:       map[string]*SchemaInfo{},
		executionPlan: map[string][]*executor.ExplainRecord{},
		sysVars:       map[string]string{},
	}
	if parent == nil {
		return ctx
	}
	ctx.schemaHasLoad = parent.schemaHasLoad
	ctx.currentSchema = parent.currentSchema
	for schemaName, schema := range parent.schemas {
		newSchema := &SchemaInfo{
			Tables: map[string]*TableInfo{},
		}
		if schema == nil || schema.Tables == nil {
			continue
		}
		for tableName, table := range schema.Tables {
			newSchema.Tables[tableName] = &TableInfo{
				Size:          table.Size,
				sizeLoad:      table.sizeLoad,
				isLoad:        table.isLoad,
				OriginalTable: table.OriginalTable,
				MergedTable:   table.MergedTable,
				AlterTables:   table.AlterTables,
			}
		}
		ctx.schemas[schemaName] = newSchema
	}

	for k, v := range parent.sysVars {
		ctx.sysVars[k] = v
	}
	return ctx
}

func (c *Context) GetSysVar(name string) (string, bool) {
	v, exist := c.sysVars[name]
	return v, exist
}

func (c *Context) AddSysVar(name, value string) {
	c.sysVars[name] = value
	return
}

func (c *Context) HasLoadSchemas() bool {
	return c.schemaHasLoad
}

func (c *Context) SetSchemasLoad() {
	c.schemaHasLoad = true
}

func (c *Context) LoadSchemas(schemas []string) {
	if c.HasLoadSchemas() {
		return
	}
	for _, schema := range schemas {
		c.schemas[schema] = &SchemaInfo{}
	}
	c.SetSchemasLoad()
}

// Schemas return all schemas info in current context.
func (c *Context) Schemas() map[string]*SchemaInfo {
	return c.schemas
}

func (c *Context) GetSchema(schemaName string) (*SchemaInfo, bool) {
	schema, has := c.schemas[schemaName]
	return schema, has
}

func (c *Context) HasSchema(schemaName string) (has bool) {
	_, has = c.GetSchema(schemaName)
	return
}

func (c *Context) AddSchema(name string) {
	if c.HasSchema(name) {
		return
	}
	c.schemas[name] = &SchemaInfo{
		Tables: map[string]*TableInfo{},
	}
}

func (c *Context) DelSchema(name string) {
	delete(c.schemas, name)
}

func (c *Context) HasLoadTables(schemaName string) (hasLoad bool) {
	if schema, ok := c.GetSchema(schemaName); ok {
		if schema.Tables == nil {
			hasLoad = false
		} else {
			hasLoad = true
		}
	}
	return
}

func (c *Context) LoadTables(schemaName string, tablesName []string) {
	schema, ok := c.GetSchema(schemaName)
	if !ok {
		return
	}
	if c.HasLoadTables(schemaName) {
		return
	}
	schema.Tables = map[string]*TableInfo{}
	for _, name := range tablesName {
		schema.Tables[name] = &TableInfo{
			isLoad:      true,
			AlterTables: []*ast.AlterTableStmt{},
		}
	}
}

func (c *Context) GetTable(schemaName, tableName string) (*TableInfo, bool) {
	schema, SchemaExist := c.GetSchema(schemaName)
	if !SchemaExist {
		return nil, false
	}
	if !c.HasLoadTables(schemaName) {
		return nil, false
	}
	table, tableExist := schema.Tables[tableName]
	return table, tableExist
}

func (c *Context) HasTable(schemaName, tableName string) (has bool) {
	_, has = c.GetTable(schemaName, tableName)
	return
}

func (c *Context) AddTable(schemaName, tableName string, table *TableInfo) {
	schema, exist := c.GetSchema(schemaName)
	if !exist {
		return
	}
	if !c.HasLoadTables(schemaName) {
		return
	}
	schema.Tables[tableName] = table
}

func (c *Context) DelTable(schemaName, tableName string) {
	schema, exist := c.GetSchema(schemaName)
	if !exist {
		return
	}
	delete(schema.Tables, tableName)
}

func (c *Context) SetCurrentSchema(schema string) {
	c.currentSchema = schema
}

// CurrentSchema return current schema.
func (c *Context) CurrentSchema() string {
	return c.currentSchema
}

func (c *Context) AddExecutionPlan(sql string, records []*executor.ExplainRecord) {
	c.executionPlan[sql] = records
}

func (c *Context) GetExecutionPlan(sql string) ([]*executor.ExplainRecord, bool) {
	records, ok := c.executionPlan[sql]
	return records, ok
}

func (c *Context) UpdateContext(node ast.Node) {
	switch s := node.(type) {
	case *ast.UseStmt:
		// change current schema
		if c.HasSchema(s.DBName) {
			c.SetCurrentSchema(s.DBName)
		}
	case *ast.CreateDatabaseStmt:
		if c.HasLoadSchemas() {
			c.AddSchema(s.Name)
		}
	case *ast.CreateTableStmt:
		schemaName := c.GetSchemaName(s.Table)
		tableName := s.Table.Name.L
		if c.HasTable(schemaName, tableName) {
			return
		}
		c.AddTable(schemaName, tableName,
			&TableInfo{
				Size:          0, // table is empty after create
				sizeLoad:      true,
				isLoad:        false,
				OriginalTable: s,
				AlterTables:   []*ast.AlterTableStmt{},
			})
	case *ast.DropDatabaseStmt:
		if c.HasLoadSchemas() {
			c.DelSchema(s.Name)
		}
	case *ast.DropTableStmt:
		if c.HasLoadSchemas() {
			for _, table := range s.Tables {
				schemaName := c.GetSchemaName(table)
				tableName := table.Name.L
				if c.HasTable(schemaName, tableName) {
					c.DelTable(schemaName, tableName)
				}
			}
		}

	case *ast.AlterTableStmt:
		info, exist := c.GetTableInfo(s.Table)
		if exist {
			var oldTable *ast.CreateTableStmt
			var err error
			if info.MergedTable != nil {
				oldTable = info.MergedTable
			} else if info.OriginalTable != nil {
				oldTable, err = util.ParseCreateTableStmt(info.OriginalTable.Text())
				if err != nil {
					return
				}
			}
			info.MergedTable, _ = util.MergeAlterToTable(oldTable, s)
			info.AlterTables = append(info.AlterTables, s)
			// rename table
			if s.Table.Name.L != info.MergedTable.Table.Name.L {
				schemaName := c.GetSchemaName(s.Table)
				c.DelTable(schemaName, s.Table.Name.L)
				c.AddTable(schemaName, info.MergedTable.Table.Name.L, info)
			}
		}
	default:
	}
}

// GetSchemaName get schema name from AST or current schema.
func (c *Context) GetSchemaName(stmt *ast.TableName) string {
	if stmt.Schema.String() == "" {
		return c.currentSchema
	}

	return stmt.Schema.String()
}

// GetTableInfo get table info from context.
func (c *Context) GetTableInfo(stmt *ast.TableName) (*TableInfo, bool) {
	schema := c.GetSchemaName(stmt)
	table := stmt.Name.String()
	return c.GetTable(schema, table)
}

// IsSchemaExist check schema is exist or not.
func (c *Context) IsSchemaExist(schemaName string) (bool, error) {
	if !c.HasLoadSchemas() {
		if c.e == nil {
			return false, nil
		}

		schemas, err := c.e.ShowDatabases(false)
		if err != nil {
			return false, err
		}
		c.LoadSchemas(schemas)
	}

	lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return false, err
	}

	if lowerCaseTableNames != "0" {
		capitalizedSchema := make(map[string]struct{})
		for name := range c.Schemas() {
			capitalizedSchema[strings.ToUpper(name)] = struct{}{}
		}
		_, exist := capitalizedSchema[strings.ToUpper(schemaName)]
		return exist, nil
	}
	return c.HasSchema(schemaName), nil
}
