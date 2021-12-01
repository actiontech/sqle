package mysql

import (
	_context "context"
	"database/sql"
	_driver "database/sql/driver"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/context"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/onlineddl"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	var allRules []*driver.Rule
	for i := range RuleHandlers {
		allRules = append(allRules, &RuleHandlers[i].Rule)
	}

	driver.Register(driver.DriverTypeMySQL, newInspect, allRules)

	if err := LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template"); err != nil {
		panic(err)
	}
}

// Inspect implements driver.Driver interface
type Inspect struct {
	// Ctx is SQL context.
	Ctx *context.Context
	// cnf is task cnf, cnf variables record in rules.
	cnf *Config

	rules []*driver.Rule

	// result keep inspect result for single audited SQL.
	// It refresh on every Audit.
	result *driver.AuditResult
	// HasInvalidSql represent one of the commit sql base-validation failed.
	HasInvalidSql bool
	// currentRule is instance's rules.
	currentRule driver.Rule

	inst *driver.DSN

	log *logrus.Entry
	// dbConn is a SQL driver for MySQL.
	dbConn *executor.Executor
	// isConnected represent dbConn has Connected.
	isConnected bool
	// isOfflineAudit represent Audit without instance.
	isOfflineAudit bool
}

func newInspect(log *logrus.Entry, cfg *driver.Config) (driver.Driver, error) {
	ctx := context.NewContext(nil)
	if cfg.DSN != nil {
		ctx.UseSchema(cfg.DSN.DatabaseName)
	}

	i := &Inspect{
		log: log,
		Ctx: ctx,
		cnf: &Config{
			DMLRollbackMaxRows: -1,
			DDLOSCMinSize:      -1,
			DDLGhostMinSize:    -1,
		},

		inst:           cfg.DSN,
		rules:          cfg.Rules,
		result:         driver.NewInspectResults(),
		isOfflineAudit: cfg.DSN == nil,
	}

	for _, rule := range cfg.Rules {
		if rule.Name == ConfigDMLRollbackMaxRows {
			defaultRule := RuleHandlerMap[ConfigDMLRollbackMaxRows].Rule
			i.cnf.DMLRollbackMaxRows = rule.GetValueInt(&defaultRule)
		}
		if rule.Name == ConfigDDLOSCMinSize {
			defaultRule := RuleHandlerMap[ConfigDDLOSCMinSize].Rule
			i.cnf.DDLOSCMinSize = rule.GetValueInt(&defaultRule)
		}
		if rule.Name == ConfigDDLGhostMinSize {
			defaultRule := RuleHandlerMap[ConfigDDLGhostMinSize].Rule
			i.cnf.DDLGhostMinSize = rule.GetValueInt(&defaultRule)
		}
	}

	return i, nil
}

func (i *Inspect) IsOfflineAudit() bool {
	return i.isOfflineAudit
}

func (i *Inspect) Exec(ctx _context.Context, query string) (_driver.Result, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}

	useGhost, err := i.onlineddlWithGhost(query)
	if err != nil {
		return nil, errors.Wrap(err, "check whether use ghost or not")
	}

	if useGhost {
		node, err := i.ParseSql(query)
		if err != nil {
			return nil, errors.Wrap(err, "parse SQL")
		}
		stmt := node[0].(*ast.AlterTableStmt)
		schema := i.Ctx.GetSchemaName(stmt.Table)

		run := func(dryRun bool) error {
			executor, err := onlineddl.NewExecutor(i.log, i.inst, schema, query)
			if err != nil {
				return err
			}

			err = executor.Execute(ctx, dryRun)
			if err != nil {
				return err
			}
			return nil
		}

		i.log.Infof("dry-run gh-ost")
		if err := run(true); err != nil {
			i.log.Errorf("dry-run gh-ost error:%v", err)
			return nil, errors.Wrap(err, "dry-run gh-ost")
		}
		i.log.Infof("dry-run OK!")

		i.log.Infof("run gh-ost")
		if err := run(false); err != nil {
			i.log.Errorf("run gh-ost error:%v", err)
			return nil, errors.Wrap(err, "run gh-ost")
		}
		i.log.Infof("run OK!")

		return _driver.ResultNoRows, nil
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Exec(query)
}

func (i *Inspect) onlineddlWithGhost(query string) (bool, error) {
	if i.cnf.DDLGhostMinSize == -1 {
		return false, nil
	}

	node, err := i.ParseSql(query)
	if err != nil {
		return false, errors.Wrap(err, "parse SQL")
	}

	stmt, ok := node[0].(*ast.AlterTableStmt)
	if !ok {
		return false, nil
	}

	tableSize, err := i.getTableSize(stmt.Table)
	if err != nil {
		return false, errors.Wrap(err, "get table size")
	}

	return int64(tableSize) > i.cnf.DDLGhostMinSize, nil
}

func (i *Inspect) Tx(ctx _context.Context, queries ...string) ([]_driver.Result, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Transact(queries...)
}

func (i *Inspect) Query(ctx _context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Query(query, args...)
}

func (i *Inspect) Parse(ctx _context.Context, sqlText string) ([]driver.Node, error) {
	nodes, err := i.ParseSql(sqlText)
	if err != nil {
		return nil, err
	}

	lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
	if err != nil {
		return nil, err
	}

	var ns []driver.Node
	for i := range nodes {
		n := driver.Node{}
		fingerprint, err := util.Fingerprint(nodes[i].Text(), lowerCaseTableNames == "0")
		if err != nil {
			return nil, err
		}
		n.Fingerprint = fingerprint
		n.Text = nodes[i].Text()
		switch nodes[i].(type) {
		case ast.DMLNode:
			n.Type = driver.SQLTypeDML
		default:
			n.Type = driver.SQLTypeDDL
		}

		ns = append(ns, n)
	}
	return ns, nil
}

func (i *Inspect) Audit(ctx _context.Context, sql string) (*driver.AuditResult, error) {
	i.result = driver.NewInspectResults()

	nodes, err := i.ParseSql(sql)
	if err != nil {
		return nil, err
	}
	if i.IsOfflineAudit() {
		err = i.CheckInvalidOffline(nodes[0])
	} else {
		err = i.CheckInvalid(nodes[0])
	}
	if err != nil {
		return nil, err
	}

	if i.result.Level() == driver.RuleLevelError {
		i.HasInvalidSql = true
		i.Logger().Warnf("SQL %s invalid, %s", nodes[0].Text(), i.result.Message())
	}

	for _, rule := range i.rules {
		i.currentRule = *rule
		handler, ok := RuleHandlerMap[rule.Name]
		if !ok || handler.Func == nil {
			continue
		}
		if i.IsOfflineAudit() && !handler.IsAllowOfflineRule(nodes[0]) {
			continue
		}
		if err := handler.Func(*rule, i, nodes[0]); err != nil {
			return nil, err
		}
	}

	// print osc
	oscCommandLine, err := i.generateOSCCommandLine(nodes[0])
	if err != nil {
		return nil, err
	}
	if oscCommandLine != "" {
		i.result.Add(driver.RuleLevelNotice, fmt.Sprintf("[osc]%s", oscCommandLine))
	}
	i.Ctx.UpdateContext(nodes[0])
	return i.result, nil
}

func (i *Inspect) GenRollbackSQL(ctx _context.Context, sql string) (string, string, error) {
	if i.IsOfflineAudit() {
		return "", "", nil
	}
	if i.HasInvalidSql {
		return "", "", nil
	}

	nodes, err := i.ParseSql(sql)
	if err != nil {
		return "", "", err
	}

	rollback, reason, err := i.GenerateRollbackSql(nodes[0])
	if err != nil {
		return "", "", err
	}

	i.Ctx.UpdateContext(nodes[0])

	return rollback, reason, nil
}

func (i *Inspect) Close(ctx _context.Context) {
	i.closeDbConn()
}

func (i *Inspect) Ping(ctx _context.Context) error {
	if i.IsOfflineAudit() {
		return nil
	}

	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	return conn.Db.Ping()
}

func (i *Inspect) Schemas(ctx _context.Context) ([]string, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.ShowDatabases(true)
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
	DDLGhostMinSize    int64
}

func (i *Inspect) Context() *context.Context {
	return i.Ctx
}

func (i *Inspect) ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := util.ParseSql(sql)
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
	i.result.Add(level, message, args...)
}

// getDbConn get db conn and just connect once.
func (i *Inspect) getDbConn() (*executor.Executor, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := executor.NewExecutor(i.log, i.inst, i.Ctx.CurrentSchema())
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
	schema := i.Ctx.GetSchemaName(stmt)
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
	schemaName := i.Ctx.GetSchemaName(stmt)
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

// getTableSize get table size.
func (i *Inspect) getTableSize(stmt *ast.TableName) (float64, error) {
	exist, err := i.isTableExist(stmt)
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

// getSchemaEngine get schema default engine.
func (i *Inspect) getSchemaEngine(stmt *ast.TableName, schemaName string) (string, error) {
	if schemaName == "" {
		schemaName = i.Ctx.GetSchemaName(stmt)
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
		schemaName = i.Ctx.GetSchemaName(stmt)
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
	ti, exist := i.Ctx.GetTableInfo(stmt)
	if !exist || !ti.isLoad {
		return "", nil
	}

	for _, columnName := range columnNames {
		if !util.TableExistCol(ti.OriginalTable, columnName) {
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
		schemaName = i.Ctx.GetSchemaName(stmt)
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

// getCreateTableStmt get create table stmtNode for db by query; if table not exist, return null.
func (i *Inspect) getCreateTableStmt(stmt *ast.TableName) (*ast.CreateTableStmt, bool, error) {
	exist, err := i.isTableExist(stmt)
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, exist, nil
	}

	info, _ := i.Ctx.GetTableInfo(stmt)
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
	createTableSql, err := conn.ShowCreateTable(util.GetTableNameWithQuote(stmt))
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

// getPrimaryKey get table's primary key.
func (i *Inspect) getPrimaryKey(stmt *ast.CreateTableStmt) (map[string]struct{}, bool, error) {
	pkColumnsName, hasPk := util.GetPrimaryKey(stmt)
	if !hasPk {
		return pkColumnsName, hasPk, nil
	}
	return pkColumnsName, hasPk, nil
}

func (i *Inspect) getExecutionPlan(sql string) ([]*executor.ExplainRecord, error) {
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
	if i.IsOfflineAudit() {
		return "", nil
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
