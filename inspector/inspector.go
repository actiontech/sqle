package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/executor"
	"sqle/model"
	"strings"
)

var (
	SQL_STMT_CONFLICT_ERROR = fmt.Errorf("不能同时提交 DDL 和 DML 语句")
)

type Inspector struct {
	Results     *InspectResults
	Rules       []model.Rule
	currentRule model.Rule
	RulesFunc   map[string]func(stmt ast.StmtNode, rule string) error
	Instance    *model.Instance
	SqlArray    []*model.CommitSql
	dbConn      *executor.Executor
	isConnected bool

	// currentSchema will change after sql "use database"
	currentSchema string
	allSchema     map[string] /*schema*/ struct{}
	schemaHasLoad bool
	allTable      map[string] /*schema*/ map[string] /*table*/ struct{}
	isDDLStmt     bool
	isDMLStmt     bool

	// save create table parser object from db by query "show create table tb_1";
	// using in inspect and generate rollback sql
	createTableStmts map[string] /*schema.table*/ *ast.CreateTableStmt

	// save alter table parse object from input sql;
	alterTableStmts map[string] /*schema.table*/ []*ast.AlterTableStmt
	rollbackSqls    []string
}

func NewInspector(rules []model.Rule, instance *model.Instance, sqlArray []*model.CommitSql, Schema string) *Inspector {
	return &Inspector{
		Results:          newInspectResults(),
		Rules:            rules,
		Instance:         instance,
		currentSchema:    Schema,
		SqlArray:         sqlArray,
		allSchema:        map[string]struct{}{},
		allTable:         map[string]map[string]struct{}{},
		createTableStmts: map[string]*ast.CreateTableStmt{},
		alterTableStmts:  map[string][]*ast.AlterTableStmt{},
		rollbackSqls:     []string{},
	}
}

func (i *Inspector) addResult(ruleName string, args ...interface{}) {
	// if rule is not current rule, ignore save the message.
	if ruleName != i.currentRule.Name {
		return
	}
	i.Results.add(i.currentRule, args...)
}

func (i *Inspector) getDbConn() (*executor.Executor, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := executor.NewExecutor(i.Instance, i.currentSchema)
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

func (i *Inspector) closeDbConn() {
	if i.isConnected {
		i.dbConn.Db.Close()
		i.isConnected = false
	}
}

func (i *Inspector) getSchemaName(stmt *ast.TableName) string {
	if stmt.Schema.String() == "" {
		return i.currentSchema
	} else {
		return stmt.Schema.String()
	}
}

func (i *Inspector) isSchemaExist(schema string) (bool, error) {
	if schema == "" {
		schema = i.currentSchema
	}
	if !i.schemaHasLoad {
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		schemas, err := conn.ShowDatabases()
		if err != nil {
			return false, err
		}
		for _, schema := range schemas {
			i.allSchema[schema] = struct{}{}
		}
		i.schemaHasLoad = true
	}
	_, ok := i.allSchema[schema]
	return ok, nil
}

func (i *Inspector) getTableName(stmt *ast.TableName) string {
	var schema string
	if stmt.Schema.String() == "" {
		schema = i.currentSchema
	} else {
		schema = stmt.Schema.String()
	}
	if schema == "" {
		return fmt.Sprintf("%s", stmt.Name)
	}
	return fmt.Sprintf("%s.%s", schema, stmt.Name)
}

func (i *Inspector) getTableNameWithQuote(stmt *ast.TableName) string {
	name := strings.Replace(i.getTableName(stmt), ".", "`.`", -1)
	return fmt.Sprintf("`%s`", name)
}

func (i *Inspector) isTableExist(tableName string) (bool, error) {
	var schema = i.currentSchema
	var table = tableName
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		schema = splitStrings[0]
		table = splitStrings[1]
	}

	_, hasLoad := i.allTable[schema]
	if !hasLoad {
		schemaExist, err := i.isSchemaExist(schema)
		if err != nil {
			return schemaExist, err
		}
		if !schemaExist {
			return false, nil
		}
		conn, err := i.getDbConn()
		if err != nil {
			return false, err
		}
		tables, err := conn.ShowSchemaTables(schema)
		if err != nil {
			return false, err
		}
		i.allTable[schema] = make(map[string]struct{}, len(tables))
		for _, table := range tables {
			i.allTable[schema][table] = struct{}{}
		}
	}
	_, exist := i.allTable[schema][table]
	return exist, nil
}

// getCreateTableStmt get create table stmtNode for db by query; if table not exist, return null.
func (i *Inspector) getCreateTableStmt(tableName string) (*ast.CreateTableStmt, bool, error) {

	exist, err := i.isTableExist(tableName)
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, exist, nil
	}

	// check local memory first, for uint test
	createStmt, ok := i.createTableStmts[tableName]
	if ok {
		return createStmt, exist, nil
	}

	// create from connection
	conn, err := i.getDbConn()
	if err != nil {
		return nil, exist, err
	}
	sql, err := conn.ShowCreateTable(tableName)
	if err != nil {
		return nil, exist, err
	}
	t, err := parseOneSql(i.Instance.DbType, sql)
	if err != nil {
		return nil, exist, err
	}
	createStmt, ok = t.(*ast.CreateTableStmt)
	if !ok {
		return nil, exist, fmt.Errorf("stmt not support")
	}
	i.createTableStmts[tableName] = createStmt
	return createStmt, exist, nil
}

func (i *Inspector) updateSchemaCtx(node ast.StmtNode) {
	switch s := node.(type) {
	case *ast.UseStmt:
		// change current schema
		i.currentSchema = s.DBName
	case *ast.CreateDatabaseStmt:
		i.allSchema[s.Name] = struct{}{}
	case *ast.CreateTableStmt:
		i.createTableStmts[i.getTableName(s.Table)] = s
	case *ast.DropDatabaseStmt:
		delete(i.allSchema, s.Name)
	case *ast.DropTableStmt:
		for _, table := range s.Tables {
			delete(i.alterTableStmts, i.getTableName(table))
		}
	default:
	}
}
