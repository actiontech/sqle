package inspector

import (
	"github.com/pingcap/tidb/ast"
)

type TableInfo struct {
	Size     float64
	sizeLoad bool

	// save create table parser object from db by query "show create table tb_1";
	// using in inspect and generate rollback sql
	CreateTableStmt *ast.CreateTableStmt

	// save alter table parse object from input sql;
	alterTableStmts []*ast.AlterTableStmt
}

type SchemaInfo struct {
	Tables map[string]*TableInfo
}

type Context struct {
	// currentSchema will change after sql "use database"
	currentSchema string

	schemas map[string]*SchemaInfo
	// if schemas info has collected, set true
	schemaHasLoad bool

	counterDDL uint
	counterDML uint
}

func NewContext() *Context {
	ctx := &Context{
		schemas: map[string]*SchemaInfo{},
	}
	//if parent == nil {
	//	return ctx
	//}
	//ctx.schemas = parent.schemas
	return ctx
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
		c.schemas[schema] = nil
	}
	c.SetSchemasLoad()
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
	c.schemas[name] = nil
}

func (c *Context) DelSchema(name string) {
	delete(c.schemas, name)
}

func (c *Context) HasLoadTables(schemaName string) (hasLoad bool) {
	if schema, ok := c.GetSchema(schemaName); ok {
		if schema == nil {
			hasLoad = false
		} else {
			hasLoad = true
		}
	}
	return
}

func (c *Context) LoadTables(schemaName string, tablesName []string) {
	if !c.HasSchema(schemaName) {
		return
	}
	if c.HasLoadTables(schemaName) {
		return
	}
	schema := &SchemaInfo{
		Tables: map[string]*TableInfo{},
	}
	for _, name := range tablesName {
		schema.Tables[name] = &TableInfo{
			alterTableStmts: []*ast.AlterTableStmt{},
		}
	}
	c.schemas[schemaName] = schema
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

func (c *Context) UseSchema(schema string) {
	c.currentSchema = schema
}

func (c *Context) AddDDL() {
	c.counterDDL += 1
}

func (c *Context) GetDDLCounter() uint {
	return c.counterDDL
}

func (c *Context) AddDML() {
	c.counterDML += 1
}

func (c *Context) GetDMLCounter() uint {
	return c.counterDML
}

func (i *Inspect) updateContext(node ast.Node) {
	ctx := i.Ctx
	switch s := node.(type) {
	case *ast.UseStmt:
		// change current schema
		if ctx.HasSchema(s.DBName) {
			ctx.UseSchema(s.DBName)
		}
	case *ast.CreateDatabaseStmt:
		if ctx.HasLoadSchemas() {
			ctx.AddSchema(s.Name)
		}
	case *ast.CreateTableStmt:
		schemaName := i.getSchemaName(s.Table)
		tableName := s.Table.Name.L
		if ctx.HasTable(schemaName, tableName) {
			return
		}
		ctx.AddTable(schemaName, tableName,
			&TableInfo{
				Size:            0, // table is empty after create
				sizeLoad:        true,
				CreateTableStmt: s,
			})
	case *ast.DropDatabaseStmt:
		if ctx.HasLoadSchemas() {
			ctx.DelSchema(s.Name)
		}
	case *ast.DropTableStmt:
		if ctx.HasLoadSchemas() {
			for _, table := range s.Tables {
				schemaName := i.getSchemaName(table)
				tableName := table.Name.L
				if ctx.HasTable(schemaName, tableName) {
					ctx.DelTable(schemaName, tableName)
				}
			}
		}

	case *ast.AlterTableStmt:
		info, exist := i.getTableInfo(s.Table)
		if exist {
			info.CreateTableStmt, _ = mergeAlterToTable(info.CreateTableStmt, s)
			info.alterTableStmts = append(info.alterTableStmts, s)
			// rename table
			if s.Table.Name.L != info.CreateTableStmt.Table.Name.L {
				schemaName := i.getSchemaName(s.Table)
				i.Ctx.DelTable(schemaName, s.Table.Name.L)
				i.Ctx.AddTable(schemaName, info.CreateTableStmt.Table.Name.L, info)
			}
		}
	default:
	}
}
