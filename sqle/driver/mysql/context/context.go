package context

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
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

	lowerCaseTableNames, err := c.GetSystemVariable(SysVarLowerCaseTableNames)
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

// IsTableExist check table is exist or not.
func (c *Context) IsTableExist(stmt *ast.TableName) (bool, error) {
	schemaName := c.GetSchemaName(stmt)
	schemaExist, err := c.IsSchemaExist(schemaName)
	if err != nil {
		return schemaExist, err
	}
	if !schemaExist {
		return false, nil
	}

	if !c.HasLoadTables(schemaName) {
		if c.e == nil {
			return false, nil
		}

		tables, err := c.e.ShowSchemaTables(schemaName)
		if err != nil {
			return false, err
		}
		c.LoadTables(schemaName, tables)
	}

	lowerCaseTableNames, err := c.GetSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return false, err
	}

	if lowerCaseTableNames != "0" {
		capitalizedTable := make(map[string]struct{})
		schemaInfo, ok := c.GetSchema(schemaName)
		if !ok {
			return false, fmt.Errorf("schema %s not exist", schemaName)
		}

		for name := range schemaInfo.Tables {
			capitalizedTable[strings.ToUpper(name)] = struct{}{}
		}
		_, exist := capitalizedTable[strings.ToUpper(stmt.Name.String())]
		return exist, nil
	}
	return c.HasTable(schemaName, stmt.Name.String()), nil
}

const (
	SysVarLowerCaseTableNames = "lower_case_table_names"
)

// GetSystemVariable get system variable.
func (c *Context) GetSystemVariable(name string) (string, error) {
	v, exist := c.GetSysVar(name)
	if exist {
		return v, nil
	}

	if c.e == nil {
		return "", nil
	}

	results, err := c.e.Db.Query(fmt.Sprintf(`SHOW GLOBAL VARIABLES LIKE '%v'`, name))
	if err != nil {
		return "", err
	}
	if len(results) != 1 {
		return "", fmt.Errorf("unexpeted results when query system variable")
	}

	value := results[0]["Value"]
	c.AddSysVar(name, value.String)
	return value.String, nil
}

// GetCreateTableStmt get create table stmtNode for db by query; if table not exist, return null.
func (c *Context) GetCreateTableStmt(stmt *ast.TableName) (*ast.CreateTableStmt, bool, error) {
	exist, err := c.IsTableExist(stmt)
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, exist, nil
	}

	info, _ := c.GetTableInfo(stmt)
	if info.MergedTable != nil {
		return info.MergedTable, exist, nil
	}
	if info.OriginalTable != nil {
		return info.OriginalTable, exist, nil
	}

	if c.e == nil {
		return nil, false, nil
	}

	createTableSql, err := c.e.ShowCreateTable(util.GetTableNameWithQuote(stmt))
	if err != nil {
		return nil, exist, err
	}
	createStmt, err := util.ParseCreateTableStmt(createTableSql)
	if err != nil {
		return nil, exist, err
	}
	info.OriginalTable = createStmt
	return createStmt, exist, nil
}

// GetCollationDatabase get collation database.
func (c *Context) GetCollationDatabase(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = c.GetSchemaName(stmt)
	}
	schema, schemaExist := c.GetSchema(schemaName)
	if schemaExist && schema.collationLoad {
		return schema.DefaultCollation, nil
	}

	if c.e == nil {
		return "", nil
	}

	collation, err := c.e.ShowDefaultConfiguration("select @@collation_database", "@@collation_database")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultCollation = collation
		schema.collationLoad = true
	}
	return collation, nil
}

// GetMaxIndexOptionForTable get max index option column of table.
func (c *Context) GetMaxIndexOptionForTable(stmt *ast.TableName, columnNames []string) (string, error) {
	ti, exist := c.GetTableInfo(stmt)
	if !exist || !ti.isLoad {
		return "", nil
	}

	for _, columnName := range columnNames {
		if !util.TableExistCol(ti.OriginalTable, columnName) {
			return "", nil
		}
	}

	if c.e == nil {
		return "", nil
	}

	sqls := make([]string, 0, len(columnNames))
	for _, col := range columnNames {
		sqls = append(sqls, fmt.Sprintf("COUNT( DISTINCT ( %v ) ) / COUNT( * ) AS %v", col, col))
	}

	result, err := c.e.Db.Query(fmt.Sprintf("SELECT %v FROM %v", strings.Join(sqls, ","), stmt.Name))
	if err != nil {
		return "", fmt.Errorf("query max index option for table error: %v", err)
	}
	maxIndexOption := ""
	for _, r := range result {
		for _, value := range r {
			if maxIndexOption == "" {
				maxIndexOption = value.String
				continue
			}
			if strings.Compare(value.String, maxIndexOption) > 0 {
				maxIndexOption = value.String
			}
		}
	}
	return maxIndexOption, nil
}

// GetSchemaCharacter get schema default character.
func (c *Context) GetSchemaCharacter(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = c.GetSchemaName(stmt)
	}
	schema, schemaExist := c.GetSchema(schemaName)
	if schemaExist {
		if schema.characterLoad {
			return schema.DefaultCharacter, nil
		}
	}

	if c.e == nil {
		return "", nil
	}

	character, err := c.e.ShowDefaultConfiguration("select @@character_set_database", "@@character_set_database")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultCharacter = character
		schema.characterLoad = true
	}
	return character, nil
}

// GetSchemaEngine get schema default engine.
func (c *Context) GetSchemaEngine(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = c.GetSchemaName(stmt)
	}
	schema, schemaExist := c.GetSchema(schemaName)
	if schemaExist {
		if schema.engineLoad {
			return schema.DefaultEngine, nil
		}
	}

	if c.e == nil {
		return "", nil
	}

	engine, err := c.e.ShowDefaultConfiguration("select @@default_storage_engine", "@@default_storage_engine")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultEngine = engine
		schema.engineLoad = true
	}
	return engine, nil
}

// getTableSize get table size.
func (i *Inspect) getTableSize(stmt *ast.TableName) (float64, error) {
	exist, err := i.Ctx.IsTableExist(stmt)
	if err != nil {
		return 0, errors.Wrapf(err, "check table exist when get table size")
	}
	if !exist {
		return 0, nil
	}

	info, _ := i.Ctx.GetTableInfo(stmt)
	if !info.sizeLoad {
		conn, err := i.getDbConn()
		if err != nil {
			return 0, err
		}
		size, err := conn.ShowTableSizeMB(i.Ctx.GetSchemaName(stmt), stmt.Name.String())
		if err != nil {
			return 0, err
		}
		info.Size = size
	}
	return info.Size, nil
}
