package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/sirupsen/logrus"
	"sqle/errors"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type Inspector interface {
	Context() *Context
	SqlType() string
	SqlInvalid() bool
	Add(sql *model.Sql, action func(sql *model.Sql) error) error
	Do() error
	Advise(rules []model.Rule) error
	GenerateAllRollbackSql() ([]*model.RollbackSql, error)
	Commit(sql *model.Sql) error
	ParseSql(sql string) ([]ast.Node, error)
	Logger() *logrus.Entry
}

func NewInspector(entry *logrus.Entry, ctx *Context, task *model.Task, relateTasks []model.Task,
	rules map[string]model.Rule) Inspector {
	if task.Instance.DbType == model.DB_TYPE_SQLSERVER {
		return NeSqlserverInspect(entry, ctx, task, relateTasks, rules)
	} else {
		return NewInspect(entry, ctx, task, relateTasks, rules)
	}
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
}

type Inspect struct {
	Ctx           *Context
	config        *Config
	Results       *InspectResults
	HasInvalidSql bool
	currentRule   model.Rule
	Task          *model.Task
	RelateTasks   []model.Task
	log           *logrus.Entry
	dbConn        *executor.Executor
	isConnected   bool
	counterDDL    uint
	counterDML    uint
	SqlArray      []*model.Sql
	SqlAction     []func(sql *model.Sql) error
}

func NewInspect(entry *logrus.Entry, ctx *Context, task *model.Task, relateTasks []model.Task,
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
		Ctx:         ctx,
		config:      config,
		Results:     newInspectResults(),
		Task:        task,
		RelateTasks: relateTasks,
		log:         entry,
		SqlArray:    []*model.Sql{},
		SqlAction:   []func(sql *model.Sql) error{},
	}
}

func (i *Inspect) Context() *Context {
	return i.Ctx
}

func (i *Inspect) SqlType() string {
	if i.counterDML > 0 && i.counterDDL > 0 {
		return model.SQL_TYPE_MULTI
	}
	if i.counterDML > 0 {
		return model.SQL_TYPE_DML
	} else {
		return model.SQL_TYPE_DDL
	}
}

func (i *Inspect) SqlInvalid() bool {
	return i.HasInvalidSql
}

func (i *Inspect) Add(sql *model.Sql, action func(sql *model.Sql) error) error {
	nodes, err := i.ParseSql(sql.Content)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		switch node.(type) {
		case ast.DDLNode:
			i.counterDDL += 1
		case ast.DMLNode:
			i.counterDML += 1
		}
	}
	sql.Stmts = nodes
	i.SqlArray = append(i.SqlArray, sql)
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *Inspect) Do() error {
	if i.SqlType() == model.SQL_TYPE_MULTI {
		i.Logger().Error(errors.SQL_STMT_CONFLICT_ERROR)
		return errors.SQL_STMT_CONFLICT_ERROR
	}
	for n, sql := range i.SqlArray {
		err := i.SqlAction[n](sql)
		if err != nil {
			return err
		}
		// update schema info
		for _, node := range sql.Stmts {
			i.updateContext(node)
		}
	}
	return nil
}

func (i *Inspect) ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := parseSql(i.Task.Instance.DbType, sql)
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

func (i *Inspect) getDbConn() (*executor.Executor, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := executor.NewExecutor(i.log, i.Task.Instance, i.Ctx.currentSchema)
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

func (i *Inspect) closeDbConn() {
	if i.isConnected {
		i.dbConn.Db.Close()
		i.isConnected = false
	}
}

func (i *Inspect) getSchemaName(stmt *ast.TableName) string {
	if stmt.Schema.String() == "" {
		return i.Ctx.currentSchema
	} else {
		return stmt.Schema.String()
	}
}

func (i *Inspect) isSchemaExist(schemaName string) (bool, error) {
	if !i.Ctx.HasLoadSchemas() {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		schemas, err := conn.ShowDatabases()
		if err != nil {
			return false, err
		}
		i.Ctx.LoadSchemas(schemas)
	}
	return i.Ctx.HasSchema(schemaName), nil
}

func (i *Inspect) getTableName(stmt *ast.TableName) string {
	schema := i.getSchemaName(stmt)
	if schema == "" {
		return fmt.Sprintf("%s", stmt.Name)
	}
	return fmt.Sprintf("%s.%s", schema, stmt.Name)
}

func (i *Inspect) getTableNameWithQuote(stmt *ast.TableName) string {
	name := strings.Replace(i.getTableName(stmt), ".", "`.`", -1)
	return fmt.Sprintf("`%s`", name)
}

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
	return i.Ctx.HasTable(schemaName, stmt.Name.String()), nil
}

func (i *Inspect) getTableInfo(stmt *ast.TableName) (*TableInfo, bool) {
	schema := i.getSchemaName(stmt)
	table := stmt.Name.String()
	return i.Ctx.GetTable(schema, table)
}

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

func (i *Inspect) parseCreateTableStmt(sql string) (*ast.CreateTableStmt, error) {
	t, err := parseOneSql(i.Task.Instance.DbType, sql)
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

func (i *Inspect) getPrimaryKey(stmt *ast.CreateTableStmt) (map[string]struct{}, bool, error) {
	var pkColumnsName = map[string]struct{}{}
	schemaName := i.getSchemaName(stmt.Table)
	tableName := stmt.Table.Name.String()

	pkColumnsName, hasPk := getPrimaryKey(stmt)
	if !hasPk {
		return pkColumnsName, hasPk, nil
	}
	// for mycat, while schema is a sharding schema, primary key is not a unique column
	// the primary key add the sharding column looks like a primary key
	if i.Task.Instance.DbType == model.DB_TYPE_MYCAT {
		mycatConfig := i.Task.Instance.MycatConfig
		ok, err := mycatConfig.IsShardingSchema(schemaName)
		if err != nil {
			return pkColumnsName, hasPk, err
		}
		if ok {
			shardingColumn, err := mycatConfig.GetShardingColumn(schemaName, tableName)
			if err != nil {
				return pkColumnsName, false, err
			}
			pkColumnsName[shardingColumn] = struct{}{}
		}
	}
	return pkColumnsName, hasPk, nil
}
