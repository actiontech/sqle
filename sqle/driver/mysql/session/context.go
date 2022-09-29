package session

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
)

type columnInfo struct {
	cardinality *int
}

type TableInfo struct {
	columns map[string]*columnInfo

	// SHOW TABLE STATUS
	tableStatus struct {
		rows *int
	}

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

type HistorySQLInfo struct {
	HasDML bool
	HasDDL bool
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

	// historySqlInfo historical sql information record
	historySqlInfo *HistorySQLInfo
}

type contextOption func(*Context)

func (o contextOption) apply(c *Context) {
	o(c)
}

// NewContext creates a new context.
func NewContext(parent *Context, opts ...contextOption) *Context {
	ctx := &Context{
		schemas:        map[string]*SchemaInfo{},
		executionPlan:  map[string][]*executor.ExplainRecord{},
		sysVars:        map[string]string{},
		historySqlInfo: &HistorySQLInfo{},
	}

	for _, opt := range opts {
		opt.apply(ctx)
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

func WithExecutor(e *executor.Executor) contextOption {
	return func(ctx *Context) {
		ctx.e = e
		e.SetLowerCaseTableNames(ctx.IsLowerCaseTableName())
	}
}

func (c *Context) GetHistorySQLInfo() *HistorySQLInfo {
	if c.historySqlInfo == nil {
		c.historySqlInfo = &HistorySQLInfo{}
	}
	return c.historySqlInfo
}

func (c *Context) hasLoadSchemas() bool {
	return c.schemaHasLoad
}

func (c *Context) setSchemasLoad() {
	c.schemaHasLoad = true
}

func (c *Context) loadSchemas(schemas []string) {
	if c.hasLoadSchemas() {
		return
	}
	isLowerCaseTableName := c.IsLowerCaseTableName()
	for _, schema := range schemas {
		if isLowerCaseTableName {
			schema = strings.ToLower(schema)
		}
		c.schemas[schema] = &SchemaInfo{}
	}
	c.setSchemasLoad()
}

// Schemas return all schemas info in current context.
func (c *Context) Schemas() map[string]*SchemaInfo {
	return c.schemas
}

func (c *Context) IsLowerCaseTableName() bool {
	lowerCaseTableNames, err := c.GetSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		log.NewEntry().Errorf("fail to load system variable lower_case_table_names, error: %v", err)
		// get system variable lower_case_table_names failed, default using false.
		return false
	}
	return lowerCaseTableNames != "0"
}

func (c *Context) getSchema(schemaName string) (*SchemaInfo, bool) {
	if c.IsLowerCaseTableName() {
		schemaName = strings.ToLower(schemaName)
	}
	schema, has := c.schemas[schemaName]
	return schema, has
}

func (c *Context) hasSchema(schemaName string) (has bool) {
	_, has = c.getSchema(schemaName)
	return
}

func (c *Context) addSchema(name string) {
	if c.hasSchema(name) {
		return
	}
	if c.IsLowerCaseTableName() {
		name = strings.ToLower(name)
	}
	c.schemas[name] = &SchemaInfo{
		Tables: map[string]*TableInfo{},
	}
}

func (c *Context) delSchema(name string) {
	if c.IsLowerCaseTableName() {
		name = strings.ToLower(name)
	}
	delete(c.schemas, name)
}

func (c *Context) hasLoadTables(schemaName string) (hasLoad bool) {
	if schema, ok := c.getSchema(schemaName); ok {
		if schema.Tables == nil {
			hasLoad = false
		} else {
			hasLoad = true
		}
	}
	return
}

func (c *Context) loadTables(schemaName string, tablesName []string) {
	schema, ok := c.getSchema(schemaName)
	if !ok {
		return
	}
	if c.hasLoadTables(schemaName) {
		return
	}
	schema.Tables = map[string]*TableInfo{}
	isLowerCaseTableName := c.IsLowerCaseTableName()
	for _, name := range tablesName {
		if isLowerCaseTableName {
			name = strings.ToLower(name)
		}
		schema.Tables[name] = &TableInfo{
			isLoad:      true,
			AlterTables: []*ast.AlterTableStmt{},
		}
	}
}

func (c *Context) getTable(schemaName, tableName string) (*TableInfo, bool) {
	schema, SchemaExist := c.getSchema(schemaName)
	if !SchemaExist {
		return nil, false
	}
	if !c.hasLoadTables(schemaName) {
		return nil, false
	}
	if c.IsLowerCaseTableName() {
		tableName = strings.ToLower(tableName)
	}
	table, tableExist := schema.Tables[tableName]
	return table, tableExist
}

func (c *Context) hasTable(schemaName, tableName string) (has bool) {
	_, has = c.getTable(schemaName, tableName)
	return
}

func (c *Context) addTable(schemaName, tableName string, table *TableInfo) {
	schema, exist := c.getSchema(schemaName)
	if !exist {
		return
	}
	if !c.hasLoadTables(schemaName) {
		return
	}
	if c.IsLowerCaseTableName() {
		tableName = strings.ToLower(tableName)
	}
	schema.Tables[tableName] = table
}

func (c *Context) delTable(schemaName, tableName string) {
	schema, exist := c.getSchema(schemaName)
	if !exist {
		return
	}
	if c.IsLowerCaseTableName() {
		tableName = strings.ToLower(tableName)
	}
	delete(schema.Tables, tableName)
}

func (c *Context) SetCurrentSchema(schema string) {
	if c.IsLowerCaseTableName() {
		schema = strings.ToLower(schema)
	}
	c.currentSchema = schema
}

// CurrentSchema return current schema.
func (c *Context) CurrentSchema() string {
	return c.currentSchema
}

func (c *Context) UpdateContext(node ast.Node) {
	// from a language type perspective
	switch node.(type) {
	case ast.DMLNode:
		c.GetHistorySQLInfo().HasDML = true
	case ast.DDLNode:
		c.GetHistorySQLInfo().HasDDL = true
	default:
	}
	// from the point of view of specific sql types
	switch s := node.(type) {
	case *ast.UseStmt:
		// change current schema
		if c.hasSchema(s.DBName) {
			c.SetCurrentSchema(s.DBName)
		}
	case *ast.CreateDatabaseStmt:
		if c.hasLoadSchemas() {
			c.addSchema(s.Name)
		}
	case *ast.CreateTableStmt:
		schemaName := c.GetSchemaName(s.Table)
		tableName := s.Table.Name.String()
		if c.hasTable(schemaName, tableName) {
			return
		}
		c.addTable(schemaName, tableName,
			&TableInfo{
				Size:          0, // table is empty after create
				sizeLoad:      true,
				isLoad:        false,
				OriginalTable: s,
				AlterTables:   []*ast.AlterTableStmt{},
			})
	case *ast.DropDatabaseStmt:
		if c.hasLoadSchemas() {
			c.delSchema(s.Name)
		}
	case *ast.DropTableStmt:
		if c.hasLoadSchemas() {
			for _, table := range s.Tables {
				schemaName := c.GetSchemaName(table)
				tableName := table.Name.String()
				if c.hasTable(schemaName, tableName) {
					c.delTable(schemaName, tableName)
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
			if s.Table.Name.String() != info.MergedTable.Table.Name.String() {
				schemaName := c.GetSchemaName(s.Table)
				c.delTable(schemaName, s.Table.Name.String())
				c.addTable(schemaName, info.MergedTable.Table.Name.String(), info)
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
	return c.getTable(schema, table)
}

// IsSchemaExist check schema is exist or not.
func (c *Context) IsSchemaExist(schemaName string) (bool, error) {
	if !c.hasLoadSchemas() {
		if c.e == nil {
			return false, nil
		}

		schemas, err := c.e.ShowDatabases(false)
		if err != nil {
			return false, err
		}
		c.loadSchemas(schemas)
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
	return c.hasSchema(schemaName), nil
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

	if !c.hasLoadTables(schemaName) {
		if c.e == nil {
			return false, nil
		}

		tables, err := c.e.ShowSchemaTables(schemaName)
		if err != nil {
			return false, err
		}
		c.loadTables(schemaName, tables)
	}

	lowerCaseTableNames, err := c.GetSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return false, err
	}

	if lowerCaseTableNames != "0" {
		capitalizedTable := make(map[string]struct{})
		schemaInfo, ok := c.getSchema(schemaName)
		if !ok {
			return false, fmt.Errorf("schema %s not exist", schemaName)
		}

		for name := range schemaInfo.Tables {
			capitalizedTable[strings.ToUpper(name)] = struct{}{}
		}
		_, exist := capitalizedTable[strings.ToUpper(stmt.Name.String())]
		return exist, nil
	}
	return c.hasTable(schemaName, stmt.Name.String()), nil
}

const (
	SysVarLowerCaseTableNames = "lower_case_table_names"
)

// GetSystemVariable get system variable.
func (c *Context) GetSystemVariable(name string) (string, error) {
	v, exist := c.sysVars[name]
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
		return "", fmt.Errorf("unexpected results when query system variable")
	}

	value := results[0]["Value"]
	c.AddSystemVariable(name, value.String)
	return value.String, nil
}

func (c *Context) AddSystemVariable(name, value string) {
	c.sysVars[name] = value
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

	createTableSql, err := c.e.ShowCreateTable(utils.SupplementalQuotationMarks(stmt.Schema.String()), utils.SupplementalQuotationMarks(stmt.Name.String()))
	if err != nil {
		return nil, exist, err
	}
	createStmt, err := util.ParseCreateTableStmt(createTableSql)
	if err != nil {
		//todo to be compatible with OceanBase-MySQL-Mode
		log.Logger().Warnf("parse create table stmt failed. try to parse it as OB-MySQL-Mode. err:%v", err)
		createStmt, err = c.parseObMysqlCreateTableSql(createTableSql)
		if err != nil {
			return nil, exist, err
		}
	}
	info.OriginalTable = createStmt
	return createStmt, exist, nil
}

/*
建表语句可能如下:
CREATE TABLE `__all_server_event_history` (
  `gmt_create` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `svr_ip` varchar(46) NOT NULL,
  `svr_port` bigint(20) NOT NULL,
  `module` varchar(64) NOT NULL,
  `event` varchar(64) NOT NULL,
  `name1` varchar(256) DEFAULT '',
  `value1` varchar(256) DEFAULT '',
  `name2` varchar(256) DEFAULT '',
  `value2` longtext DEFAULT NULL,
  `name3` varchar(256) DEFAULT '',
  `value3` varchar(256) DEFAULT '',
  `name4` varchar(256) DEFAULT '',
  `value4` varchar(256) DEFAULT '',
  `name5` varchar(256) DEFAULT '',
  `value5` varchar(256) DEFAULT '',
  `name6` varchar(256) DEFAULT '',
  `value6` varchar(256) DEFAULT '',
  `extra_info` varchar(512) DEFAULT '',
  PRIMARY KEY (`gmt_create`, `svr_ip`, `svr_port`)
) DEFAULT CHARSET = utf8mb4 ROW_FORMAT = COMPACT COMPRESSION = 'none' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 10 TABLEGROUP = 'oceanbase'
 partition by key_v2(svr_ip, svr_port)
(partition p0,
partition p1,
partition p2,
partition p3,
partition p4,
partition p5,
partition p6,
partition p7,
partition p8,
partition p9,
partition p10,
partition p11,
partition p12,
partition p13,
partition p14,
partition p15)

当左括号的数量和右括号的数量相同时, 以正好数量相等时的右括号为分界线, 后边是options，oceanbase mysql模式下的show create table结果返回的options中包含mysql不支持的options。为了能解析，临时处理方案是把options都截掉

假设建表语句如示例中所示, 左括号和右括号正好相等时最后一个右括号是 ') DEFAULT CHARSET = utf8mb4' 处第一个字符, 这个字符往前是建表语句的表结构声明部分, 后半部分是表参数, 需要截掉后半部分
*/
func (c *Context) parseObMysqlCreateTableSql(createTableSql string) (*ast.CreateTableStmt, error) {
	leftCount, rightCount := 0, 0
	for i, s := range createTableSql {
		if s == '(' {
			leftCount++
		}
		if s == ')' {
			rightCount++
		}
		if leftCount != 0 && leftCount == rightCount {
			return util.ParseCreateTableStmt(createTableSql[0 : i+1])
		}
	}

	return nil, fmt.Errorf("convert OB MySQL create table sql failed")
}

// GetCollationDatabase get collation database.
func (c *Context) GetCollationDatabase(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = c.GetSchemaName(stmt)
	}
	schema, schemaExist := c.getSchema(schemaName)
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
func (c *Context) GetMaxIndexOptionForTable(stmt *ast.TableName, columnNames []string) (float64, error) {
	ti, exist := c.GetTableInfo(stmt)
	if !exist || !ti.isLoad {
		return -1, nil
	}

	for _, columnName := range columnNames {
		if !util.TableExistCol(ti.OriginalTable, columnName) {
			return -1, nil
		}
	}

	if c.e == nil {
		return -1, nil
	}

	sqls := make([]string, 0, len(columnNames))
	for _, col := range columnNames {
		sqls = append(sqls, fmt.Sprintf("COUNT( DISTINCT ( %v ) ) / COUNT( * ) * 100 AS %v", col, col))
	}

	result, err := c.e.Db.Query(fmt.Sprintf("SELECT %v FROM %v", strings.Join(sqls, ","), stmt.Name))
	if err != nil {
		return -1, fmt.Errorf("query max index option for table error: %v", err)
	}
	maxIndexOption := -1.0
	for _, r := range result {
		for _, value := range r {
			// 当表里没数据时上面的SQL查出来的结果为Null
			if value.String == "" {
				value.String = "0"
			}
			v, err := strconv.ParseFloat(value.String, 64)
			if err != nil {
				return -1, err
			}
			if maxIndexOption == -1 {
				maxIndexOption = v
				continue
			}

			if v > maxIndexOption {
				maxIndexOption = v
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
	schema, schemaExist := c.getSchema(schemaName)
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
	schema, schemaExist := c.getSchema(schemaName)
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

// GetTableSize get table size.
func (c *Context) GetTableSize(stmt *ast.TableName) (float64, error) {
	exist, err := c.IsTableExist(stmt)
	if err != nil {
		return 0, errors.Wrapf(err, "check table exist when get table size")
	}
	if !exist {
		return 0, nil
	}

	info, _ := c.GetTableInfo(stmt)
	if !info.sizeLoad {
		if c.e == nil {
			return 0, nil
		}
		size, err := c.e.ShowTableSizeMB(c.GetSchemaName(stmt), stmt.Name.String())
		if err != nil {
			return 0, err
		}
		info.Size = size
	}
	return info.Size, nil
}

// GetExecutionPlan get execution plan of SQL.
func (c *Context) GetExecutionPlan(sql string) ([]*executor.ExplainRecord, error) {
	if ep, ok := c.executionPlan[sql]; ok {
		return ep, nil
	}

	if c.e == nil {
		return nil, nil
	}

	records, err := c.e.GetExplainRecord(sql)
	if err != nil {
		return nil, err
	}

	c.executionPlan[sql] = records
	return records, nil
}

// GetTableRowCount get table row count by show table status.
func (c *Context) GetTableRowCount(tn *ast.TableName) (int, error) {
	ti, exist := c.GetTableInfo(tn)
	if !exist {
		return 0, nil
	}
	if !ti.isLoad {
		return 0, nil
	}

	if ti.tableStatus.rows == nil {
		if c.e == nil {
			return 0, nil
		}
		query := fmt.Sprintf("show table status from %s where name = '%s'", c.GetSchemaName(tn), tn.Name.String())
		if c.IsLowerCaseTableName() {
			query = fmt.Sprintf("show table status from %s where lower(name) = '%s'", c.GetSchemaName(tn), tn.Name.L)
		}

		records, err := c.e.Db.Query(query)
		if err != nil {
			return 0, errors.Wrap(err, "get table row count error")
		}

		if len(records) != 1 {
			return 0, fmt.Errorf("get table row count error, records count: %v", len(records))
		}
		rows, err := strconv.Atoi(records[0]["Rows"].String)
		if err != nil {
			return 0, errors.Wrap(err, "get table row count error when parse rows")
		}
		ti.tableStatus.rows = &rows
	}

	return *ti.tableStatus.rows, nil
}

// IsTableExistInDatabase check table exist in database.
// Sometimes, we need explain on SQL, if table not exist, we will get error.
func (c *Context) IsTableExistInDatabase(tn *ast.TableName) (bool, error) {
	exist, err := c.IsTableExist(tn)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, nil
	}

	ti, _ := c.GetTableInfo(tn)
	return ti.isLoad, nil
}

func (c *Context) GetColumnCardinality(tn *ast.TableName, columnName string) (int, error) {
	exist, err := c.IsTableExist(tn)
	if err != nil {
		return 0, errors.Wrap(err, "check table exist when get column cardinality")
	}
	if !exist {
		return 0, nil
	}

	ti, _ := c.GetTableInfo(tn)
	if ti.columns == nil || ti.columns[columnName] == nil {
		if c.e == nil {
			return 0, nil
		}

		record, err := c.e.Db.Query(fmt.Sprintf("select count(distinct `%s`) as cardinality from `%s`.`%s`", columnName, c.GetSchemaName(tn), tn.Name.O))
		if err != nil {
			return 0, errors.Wrap(err, "get column cardinality error")
		}

		if len(record) != 1 {
			return 0, fmt.Errorf("get column cardinality error, records count: %v", len(record))
		}

		cardinality, err := strconv.Atoi(record[0]["cardinality"].String)
		if err != nil {
			return 0, errors.Wrap(err, "get column cardinality error when parse cardinality")
		}

		if ti.columns == nil {
			ti.columns = make(map[string]*columnInfo)
		}
		ti.columns[columnName] = &columnInfo{
			cardinality: &cardinality,
		}
	}

	return *ti.columns[columnName].cardinality, nil
}
