package session

import (
	"database/sql"
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
	// OriginalTableError save the error about getting original table
	OriginalTableError error // todo #1630 临时缓存错误，方便跳过解析建表语句的错误

	//
	MergedTable *ast.CreateTableStmt

	// save alter table parse object from input sql;
	AlterTables []*ast.AlterTableStmt

	Selectivity map[string] /*column name or index name*/ float64 /*selectivity*/
}

type SchemaInfo struct {
	DefaultEngine    string
	engineLoad       bool
	DefaultCharacter string
	characterLoad    bool
	DefaultCollation string
	collationLoad    bool
	IsRealSchema     bool // issue #1832, 判断当前的 schema 是否真实存在于数据库中.
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
		c.schemas[schema] = &SchemaInfo{
			IsRealSchema: true,
		}
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
		schemaInfo, ok := c.getSchema(s.DBName)
		if !ok {
			return
		}
		if schemaInfo.IsRealSchema {
			// issue #1832
			err := c.UseSchema(s.DBName)
			if err != nil {
				log.Logger().Warnf("update sql context failed, error: %v", err)
			}
		}
		c.SetCurrentSchema(s.DBName)

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
			if info.MergedTable == nil || info.MergedTable.Table == nil {
				return
			}
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

type ParseShowCreateTableContentError struct { // todo #1630 临时返回一个指定的错误类型，方便跳过解析建表语句的错误
	Msg string
}

func (p *ParseShowCreateTableContentError) Error() string {
	return fmt.Sprintf("parse show create table content failed: %v", p.Msg)
}

func IsParseShowCreateTableContentErr(err error) bool {
	var target *ParseShowCreateTableContentError
	return errors.As(err, &target)
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

	if info.OriginalTableError != nil && IsParseShowCreateTableContentErr(info.OriginalTableError) { // todo #1630 临时减少解析失败时的调用次数
		return nil, false, info.OriginalTableError
	}

	createTableSql, err := c.e.ShowCreateTable(utils.SupplementalQuotationMarks(c.GetSchemaName(stmt)), utils.SupplementalQuotationMarks(stmt.Name.String()))
	if err != nil {
		return nil, exist, err
	}
	createStmt, errByMysqlParser := util.ParseCreateTableStmt(createTableSql)
	if errByMysqlParser != nil {
		//todo to be compatible with OceanBase-MySQL-Mode
		log.Logger().Warnf("parse create table stmt failed. try to parse it with compatible method. err:%v", errByMysqlParser)
		createStmt, err = c.parseCreateTableSqlCompatibly(createTableSql)
		if err != nil {
			info.OriginalTableError = &ParseShowCreateTableContentError{Msg: errByMysqlParser.Error()}
			return nil, exist, info.OriginalTableError
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

建表语句后半段是options，oceanbase mysql模式下的show create table结果返回的options中包含mysql不支持的options, 为了能解析, 方法将会倒着遍历建表语句, 每次找到右括号时截断后面的部分, 然后尝试解析一次, 直到解析成功, 此时剩余的建表语句将不在包含OB特有options
*/
func (c *Context) parseCreateTableSqlCompatibly(createTableSql string) (*ast.CreateTableStmt, error) {
	for i := len(createTableSql) - 1; i >= 0; i-- {
		if createTableSql[i] == ')' {
			stmt, err := util.ParseCreateTableStmt(createTableSql[0 : i+1])
			if err == nil {
				return stmt, nil
			}
		}
	}
	errMsg := "parse create table sql with compatible method failed"
	log.Logger().Errorf(errMsg)
	return nil, errors.New(errMsg)
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

type index struct {
	SchemaName string
	TableName  string
	IndexName  string
}

/*
示例：

	mysql> [透传语句]SELECT (s.CARDINALITY / t.TABLE_ROWS) * 100 AS INDEX_SELECTIVITY, s.INDEX_NAME FROM INFORMATION_SCHEMA.STATISTICS s JOIN INFORMATION_SCHEMA.TABLES t ON s.TABLE_SCHEMA = t.TABLE_SCHEMA AND s.TABLE_NAME = t.TABLE_NAME WHERE (s.TABLE_SCHEMA , s.TABLE_NAME , s.INDEX_NAME) IN (("db_name","table_name","idx_name_1"),("db_name","table_name","idx_name_2"));

									  ↓包含透传语句时会多出info列
	+-------------------+------------+--------------------+
	| INDEX_SELECTIVITY | INDEX_NAME | info               |
	+-------------------+------------+--------------------+
	|          100.0000 | idx_name_2 | set_1700620716_1   |
	|           28.5714 | idx_name_1 | set_1700620716_1   |
	|           28.5714 | idx_name_1 | set_1700620716_1   |
	+-------------------+------------+--------------------+
*/
func (c *Context) getSelectivityByIndex(indexes []index) (map[string] /*index name*/ float64, error) {
	if len(indexes) == 0 {
		return make(map[string]float64, 0), nil
	}
	if c.e == nil {
		return nil, nil
	}
	values := make([]string, 0, len(indexes))
	for _, index := range indexes {
		values = append(
			values,
			fmt.Sprintf("('%s', '%s', '%s')", index.SchemaName, index.TableName, index.IndexName),
		)
	}
	results, err := c.e.Db.Query(
		fmt.Sprintf(
			`SELECT (s.CARDINALITY / t.TABLE_ROWS) * 100 AS INDEX_SELECTIVITY,s.INDEX_NAME FROM INFORMATION_SCHEMA.STATISTICS s JOIN INFORMATION_SCHEMA.TABLES t ON s.TABLE_SCHEMA = t.TABLE_SCHEMA AND s.TABLE_NAME = t.TABLE_NAME WHERE (s.TABLE_SCHEMA , s.TABLE_NAME , s.INDEX_NAME) IN (%s);`,
			strings.Join(values, ","),
		),
	)
	if err != nil {
		return nil, err
	}

	var selectivityValue float64
	var indexSelectivityMap = make(map[string]float64, len(indexes))
	var indexSelectivity, indexName sql.NullString
	for _, resultMap := range results {
		indexSelectivity = resultMap["INDEX_SELECTIVITY"]
		indexName = resultMap["INDEX_NAME"]
		if indexSelectivity.String == "" {
			// 跳过选择性为空的列
			continue
		}
		selectivityValue, err = strconv.ParseFloat(indexSelectivity.String, 64)
		if err != nil {
			return nil, err
		}
		indexSelectivityMap[indexName.String] = selectivityValue
	}
	return indexSelectivityMap, nil
}

func (c *Context) getSelectivity(schema, table, name string) (float64, bool) {
	tableInfo, exist := c.getTable(schema, table)
	if !exist {
		return -1, false
	}
	if tableInfo.Selectivity == nil {
		// selectivity not cached
		return -1, false
	}
	if selectivity, ok := tableInfo.Selectivity[name]; ok {
		return selectivity, true
	}
	return -1, false
}

func (c *Context) addSelectivity(schema, table, name string, selectivity float64) {
	tableInfo, exist := c.getTable(schema, table)
	if !exist {
		return
	}
	if tableInfo.Selectivity == nil {
		tableInfo.Selectivity = make(map[string]float64)
	}
	tableInfo.Selectivity[name] = selectivity
}

func (c *Context) GetSelectivityOfIndex(stmt *ast.TableName, indexNames []string) (map[string]float64, error) {
	if len(indexNames) == 0 || stmt == nil {
		return nil, nil
	}
	if exist, _ := c.IsTableExist(stmt); !exist {
		// would not get selectivity if table not exist
		return nil, fmt.Errorf("table not exist")
	}
	schemaName := c.GetSchemaName(stmt)
	tableName := stmt.Name.L
	cachedIndexSelectivity := make(map[string]float64)
	indexes := make([]index, 0, len(indexNames))
	for _, indexName := range indexNames {
		if selectivity, ok := c.getSelectivity(schemaName, tableName, indexName); ok {
			cachedIndexSelectivity[indexName] = selectivity
		} else {
			indexes = append(indexes, index{
				SchemaName: schemaName,
				TableName:  tableName,
				IndexName:  indexName,
			})
		}
	}
	indexSelectivity, err := c.getSelectivityByIndex(indexes)
	if err != nil {
		return nil, fmt.Errorf("get selectivity by index error: %v", err)
	}

	for indexName, selectivity := range indexSelectivity {
		c.addSelectivity(schemaName, tableName, indexName, selectivity)
	}
	for indexName, selectivity := range cachedIndexSelectivity {
		indexSelectivity[indexName] = selectivity
	}
	return indexSelectivity, nil
}

type column struct {
	SchemaName string
	TableName  string
	ColumnName string
}

/*
示例：

	mysql> [TDSQL透传语句]SELECT COUNT( DISTINCT ( name ) ) / COUNT( * ) * 100 AS name,COUNT( DISTINCT ( age  ) ) / COUNT( * ) * 100 AS age FROM (SELECT name,age FROM test.test_table LIMIT 50000) t;
						 ↓包含透传语句时会多出info列
	+---------+---------+--------------------+
	| name    | age     | info               |
	+---------+---------+--------------------+
	| 50.0000 | 75.0000 | set_1700620716_1   |
	+---------+---------+--------------------+
*/
func (c *Context) getSelectivityByColumn(columns []column) (map[string] /*index name*/ float64, error) {
	if len(columns) == 0 {
		return make(map[string]float64, 0), nil
	}
	if c.e == nil {
		return nil, nil
	}
	var selectivityValue float64
	var columnSelectivityMap = make(map[string]float64, len(columns))

	sqls := make([]string, 0, len(columns))
	selectColumns := make([]string, 0, len(columns))
	for _, column := range columns {
		sqls = append(
			sqls,
			fmt.Sprintf("COUNT( DISTINCT ( `%v` ) ) / COUNT( * ) * 100 AS '%v'", column.ColumnName, column.ColumnName),
		)
		selectColumns = append(selectColumns, "`"+column.ColumnName+"`")
		columnSelectivityMap[column.ColumnName] = 0
	}

	results, err := c.e.Db.Query(
		fmt.Sprintf(
			"SELECT %v FROM (SELECT %v FROM `%v`.`%v` LIMIT 50000) t;",
			strings.Join(sqls, ","),
			strings.Join(selectColumns, ","),
			columns[0].SchemaName, columns[0].TableName,
		),
	)
	if err != nil {
		return nil, err
	}
	for _, resultMap := range results {
		for k, v := range resultMap {
			if _, ok := columnSelectivityMap[k]; !ok {
				continue
			}
			if v.String == "" {
				selectivityValue = -1
			} else {
				selectivityValue, err = strconv.ParseFloat(v.String, 64)
				if err != nil {
					return nil, err
				}
			}
			columnSelectivityMap[k] = selectivityValue
		}
	}
	return columnSelectivityMap, nil
}

func (c *Context) GetSelectivityOfColumns(stmt *ast.TableName, indexColumns []string) (map[string] /*column name*/ float64, error) {
	if stmt == nil || len(indexColumns) == 0 {
		return nil, nil
	}
	if exist, _ := c.IsTableExist(stmt); !exist {
		// would not get selectivity if table not exist
		return nil, fmt.Errorf("table not exist")
	}
	schemaName := c.GetSchemaName(stmt)
	tableName := stmt.Name.L
	cachedIndexSelectivity := make(map[string]float64)
	columns := make([]column, 0, len(indexColumns))
	for _, columnName := range indexColumns {
		if selectivity, ok := c.getSelectivity(schemaName, tableName, columnName); ok {
			cachedIndexSelectivity[columnName] = selectivity
		} else {
			columns = append(columns, column{
				SchemaName: schemaName,
				TableName:  tableName,
				ColumnName: columnName,
			})
		}
	}
	columnSelectivity, err := c.getSelectivityByColumn(columns)
	if err != nil {
		return nil, fmt.Errorf("get selectivity by column error: %v", err)
	}
	for indexName, selectivity := range columnSelectivity {
		c.addSelectivity(schemaName, tableName, indexName, selectivity)
	}
	for indexName, selectivity := range cachedIndexSelectivity {
		columnSelectivity[indexName] = selectivity
	}
	return columnSelectivity, nil
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

/*
Example:

	mysql> SELECT CHARACTER_SET_NAME FROM INFORMATION_SCHEMA.COLLATIONS WHERE COLLATION_NAME = "armscii8_bin";
	+--------------------+
	| CHARACTER_SET_NAME |
	+--------------------+
	| armscii8           |
	+--------------------+
	1 row in set (0.01 sec)
*/
func (c *Context) GetSchemaCharacterByCollation(collation string) (string, error) {
	if collation == "" || c.e == nil {
		return "", nil
	}
	return c.e.ShowDefaultConfiguration(
		fmt.Sprintf("SELECT CHARACTER_SET_NAME FROM INFORMATION_SCHEMA.COLLATIONS WHERE COLLATION_NAME = \"%s\"", collation), "CHARACTER_SET_NAME")
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
	key := fmt.Sprintf("%s.%s", c.currentSchema, sql)
	if ep, ok := c.executionPlan[key]; ok {
		return ep, nil
	}

	if c.e == nil {
		return nil, nil
	}

	records, err := c.e.GetExplainRecord(sql)
	if err != nil {
		return nil, err
	}

	c.executionPlan[key] = records
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
		query := fmt.Sprintf("show table status from `%s` where name = '%s'", c.GetSchemaName(tn), tn.Name.String())
		if c.IsLowerCaseTableName() {
			query = fmt.Sprintf("show table status from `%s` where lower(name) = '%s'", c.GetSchemaName(tn), tn.Name.L)
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

func (c *Context) UseSchema(schemaName string) error {
	_, err := c.e.Db.Exec(fmt.Sprintf("use %s", schemaName))
	if err != nil {
		return errors.Wrap(err, "exec use schema")
	}
	return nil
}

func (c *Context) GetExecutor() *executor.Executor {
	return c.e
}

func (c *Context) GetTableIndexesInfo(schema, tableName string) ([]*executor.TableIndexesInfo, error) {
	return c.e.GetTableIndexesInfo(utils.SupplementalQuotationMarks(schema), utils.SupplementalQuotationMarks(tableName))
}

func (c *Context) GetTableNameCreateTableStmtMap(joinStmt *ast.Join) map[string] /*table name or alias table name*/ *ast.CreateTableStmt {
	tableNameCreateTableStmtMap := make(map[string]*ast.CreateTableStmt)
	tableSources := util.GetTableSources(joinStmt)
	for _, tableSource := range tableSources {
		if tableNameStmt, ok := tableSource.Source.(*ast.TableName); ok {
			tableName := tableNameStmt.Name.L
			if tableSource.AsName.L != "" {
				// 如果使用别名，则需要用别名引用
				tableName = tableSource.AsName.L
			}

			createTableStmt, exist, err := c.GetCreateTableStmt(tableNameStmt)
			if err != nil || !exist {
				continue
			}
			// TODO: 跨库的 JOIN 无法区分
			tableNameCreateTableStmtMap[tableName] = createTableStmt
		}
	}
	return tableNameCreateTableStmtMap
}
