package mysql

import (
	"context"
	"database/sql"
	_driver "database/sql/driver"
	"fmt"
	"strings"

	"actiontech.cloud/sqle/sqle/sqle/driver"
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/model"

	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
)

func init() {
	var allRules []*model.Rule
	for _, ruleHandler := range RuleHandlerMap {
		ruleHandler.Rule.DBType = model.DBTypeMySQL
		allRules = append(allRules, &ruleHandler.Rule)
	}

	driver.Register(model.DBTypeMySQL, newInspect, allRules)
}

// Inspect implements Inspector interface for MySQL.
type Inspect struct {
	// Ctx is SQL context.
	Ctx *Context
	// config is task config, config variables record in rules.
	config *Config
	// Results is inspect result for commit sql.
	Results *InspectResults
	// HasInvalidSql represent one of the commit sql base-validation failed.
	HasInvalidSql bool
	// currentRule is instance's rules.
	currentRule model.Rule

	inst *model.Instance

	log *logrus.Entry
	// dbConn is a SQL driver for MySQL.
	dbConn *Executor
	// isConnected represent dbConn has Connected.
	isConnected bool
	// counterDDL is a counter for all ddl sql.
	counterDDL uint
	// counterDML is a counter for all dml sql.
	counterDML uint

	// SqlArray and SqlAction is two list for Add-Do design.
	SqlArray  []ast.Node
	SqlAction []func(node ast.Node) error
}

func newInspect(log *logrus.Entry, inst *model.Instance, schema string) driver.Driver {
	ctx := NewContext(nil)
	ctx.UseSchema(schema)
	return &Inspect{
		Ctx:     ctx,
		Results: newInspectResults(),
		log:     log,
	}
}

func (i *Inspect) Exec(ctx context.Context, query string) (_driver.Result, error) {
	return nil, nil
}

func (i *Inspect) Tx(ctx context.Context, queries ...string) ([]_driver.Result, error) {
	return nil, nil
}

func (i *Inspect) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	return nil, nil
}

func (i *Inspect) Parse(sqlText string) ([]driver.Node, error) {
	return nil, nil
}

func (i *Inspect) Audit(rules []*model.Rule, baseSQLs []*model.BaseSQL, skipAudit func(node driver.Node) bool) ([]*model.ExecuteSQL, []*model.RollbackSQL, error) {
	for _, rule := range rules {
		if rule.Name == CONFIG_DML_ROLLBACK_MAX_ROWS {
			defaultRule := RuleHandlerMap[CONFIG_DML_ROLLBACK_MAX_ROWS].Rule
			i.config.DMLRollbackMaxRows = rule.GetValueInt(&defaultRule)
		} else {
			i.config.DMLRollbackMaxRows = -1
		}

		if rule.Name == CONFIG_DDL_OSC_MIN_SIZE {
			defaultRule := RuleHandlerMap[CONFIG_DDL_OSC_MIN_SIZE].Rule
			i.config.DDLOSCMinSize = rule.GetValueInt(&defaultRule)
		} else {
			i.config.DDLOSCMinSize = -1
		}
	}

	lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return nil, nil, err
	}

	var executeSQLs []*model.ExecuteSQL
	for _, baseSQL := range baseSQLs {
		nodes, err := i.ParseSql(baseSQL.Content)
		if err != nil {
			return nil, nil, err
		}

		node := &node{
			innerNode:       nodes[0],
			isCaseSensitive: lowerCaseTableNames == "0",
		}

		if skipAudit(node) {
			executeSQL := &model.ExecuteSQL{BaseSQL: *baseSQL}
			var results InspectResults
			results.add(model.RuleLevelNormal, "白名单")
			executeSQL.AuditStatus = model.SQLAuditStatusFinished
			executeSQL.AuditLevel = results.level()
			executeSQL.AuditResult = results.message()
			executeSQLs = append(executeSQLs, executeSQL)
			continue
		}

		executeSQL, err := i.audit(node, baseSQL, rules)
		if err != nil {
			return nil, nil, err
		}
		executeSQLs = append(executeSQLs, executeSQL)
	}

	var rollbackSQLs []*model.RollbackSQL
	// TODO(@wy) ignore rollback when audit Mybatis file
	if !i.HasInvalidSql {
		rollbackSQLs, err = i.GenerateAllRollbackSql()
		if err != nil {
			return nil, nil, err
		}
	}

	return executeSQLs, rollbackSQLs, nil
}

func (i *Inspect) Close() {

}

func (i *Inspect) Ping(ctx context.Context) error {
	return nil
}

type node struct {
	innerNode       ast.Node
	isCaseSensitive bool
}

func (n *node) Text() string {
	return n.innerNode.Text()
}

func (n *node) Type() string {
	switch n.innerNode.(type) {
	case ast.DDLNode:
		return model.SQL_TYPE_DDL
	case ast.DMLNode:
		return model.SQL_TYPE_DML
	}

	return ""
}

func (n *node) Fingerprint() (string, error) {
	return Fingerprint(n.innerNode.Text(), n.isCaseSensitive)
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
}

func NewInspect(entry *logrus.Entry, ctx *Context, task *model.Task,
	rules map[string]model.Rule) *Inspect {
	ctx.UseSchema(task.Schema)

	// load config
	config := &Config{}
	if rules != nil {
		if r, ok := rules[CONFIG_DML_ROLLBACK_MAX_ROWS]; ok {
			defaultRule := RuleHandlerMap[CONFIG_DML_ROLLBACK_MAX_ROWS].Rule
			config.DMLRollbackMaxRows = r.GetValueInt(&defaultRule)
		} else {
			config.DMLRollbackMaxRows = -1
		}

		if r, ok := rules[CONFIG_DDL_OSC_MIN_SIZE]; ok {
			defaultRule := RuleHandlerMap[CONFIG_DDL_OSC_MIN_SIZE].Rule
			config.DDLOSCMinSize = r.GetValueInt(&defaultRule)
		} else {
			config.DDLOSCMinSize = -1
		}
	}
	return &Inspect{
		Ctx:     ctx,
		config:  config,
		Results: newInspectResults(),
		Task:    task,
		log:     entry,
	}
}

func (i *Inspect) Context() *Context {
	return i.Ctx
}

func (i *Inspect) SqlType() string {
	hasDML := i.counterDML > 0
	hasDDL := i.counterDDL > 0

	if hasDML && hasDDL {
		return model.SQL_TYPE_MULTI
	}

	if hasDML {
		return model.SQL_TYPE_DML
	} else if hasDDL {
		return model.SQL_TYPE_DDL
	} else {
		return ""
	}
}

func (i *Inspect) ParseSqlType() error {
	for _, commitSql := range i.Task.ExecuteSQLs {
		nodes, err := i.ParseSql(commitSql.Content)
		if err != nil {
			return err
		}
		i.addNodeCounter(nodes)
	}
	return nil
}

func (i *Inspect) addNodeCounter(nodes []ast.Node) {
	for _, node := range nodes {
		switch node.(type) {
		case ast.DDLNode:
			i.counterDDL += 1
		case ast.DMLNode:
			i.counterDML += 1
		}
	}
}

func (i *Inspect) SqlInvalid() bool {
	return i.HasInvalidSql
}

func (i *Inspect) Add(sql *model.BaseSQL, action func(node ast.Node) error) error {
	nodes, err := i.ParseSql(sql.Content)
	if err != nil {
		return err
	}
	i.addNodeCounter(nodes)
	i.SqlArray = append(i.SqlArray, nodes[0])
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *Inspect) Do() error {
	defer i.closeDbConn()

	if i.SqlType() == model.SQL_TYPE_MULTI {
		i.Logger().Error(errors.ErrSQLTypeConflict)
		return errors.ErrSQLTypeConflict
	}
	for idx, node := range i.SqlArray {
		err := i.SqlAction[idx](node)
		if err != nil {
			return err
		}
		// update schema info
		i.updateContext(node)
	}
	return nil
}

func (i *Inspect) ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := parseSql(model.DBTypeMySQL, sql)
	if err != nil {
		i.Logger().Errorf("parse sql failed, error: %v, sql: %s", err, sql)
		return nil, err
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
		node := stmt.(ast.Node)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (i *Inspect) Logger() *logrus.Entry {
	return i.log
}

func (i *Inspect) addResult(ruleName string, args ...interface{}) {
	// if rule is not current rule, ignore save the message.
	if ruleName != i.currentRule.Name {
		return
	}
	level := i.currentRule.Level
	message := RuleHandlerMap[ruleName].Message
	i.Results.add(level, message, args...)
}

// getDbConn get db conn and just connect once.
func (i *Inspect) getDbConn() (*Executor, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := NewExecutor(i.log, i.inst, i.Ctx.currentSchema)
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

// closeDbConn close db conn and just close once.
func (i *Inspect) closeDbConn() {
	if i.isConnected {
		i.dbConn.Db.Close()
		i.isConnected = false
	}
}

// getSchemaName get schema name from TableName ast;
// if schema name is default, using current schema from SQL ctx.
func (i *Inspect) getSchemaName(stmt *ast.TableName) string {
	if stmt.Schema.String() == "" {
		return i.Ctx.currentSchema
	} else {
		return stmt.Schema.String()
	}
}

// isSchemaExist determine if the schema exists in the SQL ctx;
// and lazy load schema info from db to SQL ctx.
func (i *Inspect) isSchemaExist(schemaName string) (bool, error) {
	if !i.Ctx.HasLoadSchemas() {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		schemas, err := conn.ShowDatabases(false)
		if err != nil {
			return false, err
		}
		i.Ctx.LoadSchemas(schemas)
	}

	lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return false, err
	}

	if lowerCaseTableNames != "0" {
		capitalizedSchema := make(map[string]struct{})
		for name := range i.Ctx.schemas {
			capitalizedSchema[strings.ToUpper(name)] = struct{}{}
		}
		_, exist := capitalizedSchema[strings.ToUpper(schemaName)]
		return exist, nil
	}
	return i.Ctx.HasSchema(schemaName), nil
}

// getTableName get table name from TableName ast.
func (i *Inspect) getTableName(stmt *ast.TableName) string {
	schema := i.getSchemaName(stmt)
	if schema == "" {
		return fmt.Sprintf("%s", stmt.Name)
	}
	return fmt.Sprintf("%s.%s", schema, stmt.Name)
}

// getTableNameWithQuote get table name with quote.
func (i *Inspect) getTableNameWithQuote(stmt *ast.TableName) string {
	name := strings.Replace(i.getTableName(stmt), ".", "`.`", -1)
	return fmt.Sprintf("`%s`", name)
}

// isTableExist determine if the table exists in the SQL ctx;
// and lazy load table info from db to SQL ctx.
func (i *Inspect) isTableExist(stmt *ast.TableName) (bool, error) {
	schemaName := i.getSchemaName(stmt)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return schemaExist, err
	}
	if !schemaExist {
		return false, nil
	}
	if !i.Ctx.HasLoadTables(schemaName) {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		tables, err := conn.ShowSchemaTables(schemaName)
		if err != nil {
			return false, err
		}
		i.Ctx.LoadTables(schemaName, tables)
	}

	lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return false, err
	}

	if lowerCaseTableNames != "0" {
		capitalizedTable := make(map[string]struct{})
		for name := range i.Ctx.schemas[schemaName].Tables {
			capitalizedTable[strings.ToUpper(name)] = struct{}{}
		}
		_, exist := capitalizedTable[strings.ToUpper(stmt.Name.String())]
		return exist, nil
	}
	return i.Ctx.HasTable(schemaName, stmt.Name.String()), nil
}

// getTableInfo get table info if table exist.
func (i *Inspect) getTableInfo(stmt *ast.TableName) (*TableInfo, bool) {
	schema := i.getSchemaName(stmt)
	table := stmt.Name.String()
	return i.Ctx.GetTable(schema, table)
}

// getTableSize get table size.
func (i *Inspect) getTableSize(stmt *ast.TableName) (float64, error) {
	info, exist := i.getTableInfo(stmt)
	if !exist {
		return 0, nil
	}
	if !info.sizeLoad {
		conn, err := i.getDbConn()
		if err != nil {
			return 0, err
		}
		size, err := conn.ShowTableSizeMB(i.getSchemaName(stmt), stmt.Name.String())
		if err != nil {
			return 0, err
		}
		info.Size = size
	}
	return info.Size, nil
}

// getSchemaEngine get schema default engine.
func (i *Inspect) getSchemaEngine(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = i.getSchemaName(stmt)
	}
	schema, schemaExist := i.Ctx.GetSchema(schemaName)
	if schemaExist {
		if schema.engineLoad {
			return schema.DefaultEngine, nil
		}
	}
	conn, err := i.getDbConn()
	if err != nil {
		return "", err
	}

	engine, err := conn.ShowDefaultConfiguration("select @@default_storage_engine", "@@default_storage_engine")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultEngine = engine
		schema.engineLoad = true
	}
	return engine, nil
}

// getSchemaCharacter get schema default character.
func (i *Inspect) getSchemaCharacter(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = i.getSchemaName(stmt)
	}
	schema, schemaExist := i.Ctx.GetSchema(schemaName)
	if schemaExist {
		if schema.characterLoad {
			return schema.DefaultCharacter, nil
		}
	}
	conn, err := i.getDbConn()
	if err != nil {
		return "", err
	}
	character, err := conn.ShowDefaultConfiguration("select @@character_set_database", "@@character_set_database")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultCharacter = character
		schema.characterLoad = true
	}
	return character, nil
}

func (i *Inspect) getMaxIndexOptionForTable(stmt *ast.TableName, columnNames []string) (string, error) {
	ti, exist := i.getTableInfo(stmt)
	if !exist || !ti.isLoad {
		return "", nil
	}

	for _, columnName := range columnNames {
		if !tableExistCol(ti.OriginalTable, columnName) {
			return "", nil
		}
	}

	conn, err := i.getDbConn()
	if err != nil {
		return "", err
	}
	sqls := make([]string, 0, len(columnNames))
	for _, col := range columnNames {
		sqls = append(sqls, fmt.Sprintf("COUNT( DISTINCT ( %v ) ) / COUNT( * ) AS %v", col, col))
	}
	queryIndexOptionSql := fmt.Sprintf("SELECT %v FROM %v", strings.Join(sqls, ","), stmt.Name)

	result, err := conn.Db.Query(queryIndexOptionSql)
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

func (i *Inspect) getCollationDatabase(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = i.getSchemaName(stmt)
	}
	schema, schemaExist := i.Ctx.GetSchema(schemaName)
	if schemaExist && schema.collationLoad {
		return schema.DefaultCollation, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return "", err
	}

	collation, err := conn.ShowDefaultConfiguration("select @@collation_database", "@@collation_database")
	if err != nil {
		return "", err
	}
	if schemaExist {
		schema.DefaultCollation = collation
		schema.collationLoad = true
	}
	return collation, nil
}

// parseCreateTableStmt parse create table sql text to CreateTableStmt ast.
func (i *Inspect) parseCreateTableStmt(sql string) (*ast.CreateTableStmt, error) {
	t, err := parseOneSql(model.DBTypeMySQL, sql)
	if err != nil {
		i.Logger().Errorf("parse sql from show create failed, error: %v", err)
		return nil, err
	}
	createStmt, ok := t.(*ast.CreateTableStmt)
	if !ok {
		i.Logger().Error("parse sql from show create failed, not createTableStmt")
		return nil, fmt.Errorf("stmt not support")
	}
	return createStmt, nil
}

// getCreateTableStmt get create table stmtNode for db by query; if table not exist, return null.
func (i *Inspect) getCreateTableStmt(stmt *ast.TableName) (*ast.CreateTableStmt, bool, error) {
	exist, err := i.isTableExist(stmt)
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, exist, nil
	}

	info, _ := i.getTableInfo(stmt)
	if info.MergedTable != nil {
		return info.MergedTable, exist, nil
	}
	if info.OriginalTable != nil {
		return info.OriginalTable, exist, nil
	}

	// create from connection
	conn, err := i.getDbConn()
	if err != nil {
		return nil, exist, err
	}
	createTableSql, err := conn.ShowCreateTable(getTableNameWithQuote(stmt))
	if err != nil {
		return nil, exist, err
	}
	createStmt, err := i.parseCreateTableStmt(createTableSql)
	if err != nil {
		return nil, exist, err
	}
	info.OriginalTable = createStmt
	return createStmt, exist, nil
}

// getPrimaryKey get table's primary key.
func (i *Inspect) getPrimaryKey(stmt *ast.CreateTableStmt) (map[string]struct{}, bool, error) {
	pkColumnsName, hasPk := getPrimaryKey(stmt)
	if !hasPk {
		return pkColumnsName, hasPk, nil
	}
	return pkColumnsName, hasPk, nil
}

func (i *Inspect) getExecutionPlan(sql string) ([]*ExplainRecord, error) {
	if ep, ok := i.Ctx.GetExecutionPlan(sql); ok {
		return ep, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}

	records, err := conn.Explain(sql)
	if err != nil {
		return nil, err
	}
	i.Ctx.AddExecutionPlan(sql, records)
	return records, nil
}

const (
	SysVarLowerCaseTableNames = "lower_case_table_names"
)

func (i *Inspect) getSystemVariable(name string) (string, error) {
	v, exist := i.Ctx.GetSysVar(name)
	if exist {
		return v, nil
	}

	conn, err := i.getDbConn()
	if err != nil {
		return "", err
	}
	results, err := conn.Db.Query(fmt.Sprintf(`SHOW GLOBAL VARIABLES LIKE '%v'`, name))
	if err != nil {
		return "", err
	}
	if len(results) != 1 {
		return "", fmt.Errorf("unexpeted results when query system variable")
	}

	value := results[0]["Value"]
	i.Ctx.AddSysVar(name, value.String)
	return value.String, nil
}

func (i *Inspect) GetProcedureFunctionBackupSql(sql string) ([]string, error) {
	return nil, nil
}
