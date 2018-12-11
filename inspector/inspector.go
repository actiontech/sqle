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
	Add(sql *model.Sql, action func(sql *model.Sql) error) error
	Do() error
	Advise(rules []model.Rule) error
	GenerateAllRollbackSql() ([]*model.RollbackSql, error)
	Commit(sql *model.Sql) error
	ParseSql(sql string) ([]ast.Node, error)
	Logger() *logrus.Entry
}

func NewInspector(entry *logrus.Entry, task *model.Task, rules map[string]model.Rule) Inspector {
	if task.Instance.DbType == model.DB_TYPE_SQLSERVER {
		return NeSqlserverInspect(entry, task, rules)
	} else {
		return NewInspect(entry, task, rules)
	}
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
}

type Inspect struct {
	Ctx         *Context
	config      *Config
	Results     *InspectResults
	currentRule model.Rule
	Task        *model.Task
	log         *logrus.Entry
	dbConn      *executor.Executor
	isConnected bool

	SqlArray  []*model.Sql
	SqlAction []func(sql *model.Sql) error
}

func NewInspect(entry *logrus.Entry, task *model.Task, rules map[string]model.Rule) *Inspect {
	ctx := NewContext(task.Schema)

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
		Ctx:       ctx,
		config:    config,
		Results:   newInspectResults(),
		Task:      task,
		log:       entry,
		SqlArray:  []*model.Sql{},
		SqlAction: []func(sql *model.Sql) error{},
	}
}

func (i *Inspect) Add(sql *model.Sql, action func(sql *model.Sql) error) error {
	nodes, err := i.ParseSql(sql.Content)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		switch node.(type) {
		case ast.DDLNode:
			i.Ctx.AddDDL()
		case ast.DMLNode:
			i.Ctx.AddDML()
		}
	}
	sql.Stmts = nodes
	i.SqlArray = append(i.SqlArray, sql)
	i.SqlAction = append(i.SqlAction, action)
	return nil
}

func (i *Inspect) Do() error {
	if i.Ctx.GetDMLCounter() > 0 && i.Ctx.GetDDLCounter() > 0 {
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
	i.Results.add(i.currentRule, args...)
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
	if !i.Ctx.HasLoadSchema() {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		schemas, err := conn.ShowDatabases(i.Task.Instance.DbType)
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
	if !i.Ctx.HasLoadSchemaTables(schemaName) {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		tables, err := conn.ShowSchemaTables(schemaName)
		if err != nil {
			return false, err
		}
		i.Ctx.LoadSchemaTables(schemaName, tables)
	}
	_, exist := i.Ctx.GetTableInfo(schemaName, stmt.Name.String())
	return exist, nil
}

func (i *Inspect) getTableInfo(stmt *ast.TableName) (*TableInfo, bool) {
	schema := i.getSchemaName(stmt)
	table := stmt.Name.String()
	return i.Ctx.GetTableInfo(schema, table)
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
	if info.CreateTableStmt != nil {
		return info.CreateTableStmt, exist, nil
	}

	// create from connection
	conn, err := i.getDbConn()
	if err != nil {
		return nil, exist, err
	}
	sql, err := conn.ShowCreateTable(i.getTableName(stmt))
	if err != nil {
		return nil, exist, err
	}
	t, err := parseOneSql(i.Task.Instance.DbType, sql)
	if err != nil {
		i.Logger().Errorf("parse sql from show create failed, error: %v", err)
		return nil, exist, err
	}
	createStmt, ok := t.(*ast.CreateTableStmt)
	if !ok {
		i.Logger().Error("parse sql from show create failed, not createTableStmt")
		return nil, exist, fmt.Errorf("stmt not support")
	}
	info.CreateTableStmt = createStmt
	return createStmt, exist, nil
}

func (i *Inspect) updateContext(node ast.Node) {
	switch s := node.(type) {
	case *ast.UseStmt:
		// change current schema
		i.Ctx.UseSchema(s.DBName)
	case *ast.CreateDatabaseStmt:
		i.Ctx.CreateNewSchema(s.Name)
	case *ast.CreateTableStmt:
		i.Ctx.CreateNewTable(i.getSchemaName(s.Table), s.Table.Name.String(),
			&TableInfo{
				CreateTableStmt: s,
			})
	case *ast.DropDatabaseStmt:
		i.Ctx.DeleteSchema(s.Name)
	case *ast.DropTableStmt:
		for _, table := range s.Tables {
			i.Ctx.DeleteTable(i.getSchemaName(table), table.Name.String())
		}
	case *ast.AlterTableStmt:
		info, exist := i.getTableInfo(s.Table)
		if exist {
			info.alterTableStmts = append(info.alterTableStmts, s)
		}
	default:
	}
}
