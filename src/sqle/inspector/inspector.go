package inspector

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type Inspector struct {
	Config      map[string]*model.Rule
	Db          model.Instance
	SqlArray    []*model.CommitSql
	dbConn      *executor.Conn
	isConnected bool

	// currentSchema will change after sql "use database"
	currentSchema string
	allSchema     map[string] /*schema*/ struct{}
	schemaHasLoad bool
	allTable      map[string] /*schema*/ map[string] /*table*/ struct{}
	isDDLStmt     bool
	isDMLStmt     bool
	createTable   string
	alterTable    string
}

func NewInspector(config map[string]*model.Rule, db model.Instance, sqlArray []*model.CommitSql, Schema string) *Inspector {
	return &Inspector{
		Config:        config,
		Db:            db,
		currentSchema: Schema,
		SqlArray:      sqlArray,
		allSchema:     map[string]struct{}{},
		allTable:      map[string]map[string]struct{}{},
	}
}

func (i *Inspector) getDbConn() (*executor.Conn, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := executor.NewConn(i.Db.DbType, i.Db.User, i.Db.Password, i.Db.Host, i.Db.Port, i.currentSchema)
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

func (i *Inspector) closeDbConn() {
	if i.isConnected {
		i.dbConn.Close()
		i.isConnected = false
	}
}

func (i *Inspector) isSchemaExist(schema string) (bool, error) {
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

func (i *Inspector) Inspect() ([]*model.CommitSql, error) {
	defer i.closeDbConn()

	for _, sql := range i.SqlArray {
		var stmt ast.StmtNode
		var err error
		var results *InspectResults

		stmt, err = parseOneSql(i.Db.DbType, sql.Sql)
		switch s := stmt.(type) {
		case *ast.SelectStmt:
			results, err = i.inspectSelectStmt(s)
		case *ast.UseStmt:
			results, err = i.inspectUseStmt(s)
		default:
		}
		if err != nil {
			sql.InspectStatus = model.TASK_ACTION_ERROR
			sql.InspectResult = err.Error()
			return i.SqlArray, err
		}
		sql.InspectStatus = model.TASK_ACTION_DONE
		sql.InspectLevel = results.level()
		sql.InspectResult = results.message()
	}
	return i.SqlArray, nil
}

func (i *Inspector) inspectSelectStmt(stmt *ast.SelectStmt) (*InspectResults, error) {
	results := newInspectResults()

	// check table must exist
	tablerefs := stmt.From.TableRefs
	tables := getTables(tablerefs)
	tablesName := map[string]struct{}{}
	for _, t := range tables {
		tablesName[getTableName(t)] = struct{}{}
	}
	conn, err := i.getDbConn()
	if err != nil {
		return results, err
	}
	notExistTables := []string{}
	for name, _ := range tablesName {
		exist := conn.HasTable(name)
		if conn.Error != nil {
			return results, conn.Error
		}
		if !exist {
			notExistTables = append(notExistTables, name)
		}
	}
	if len(notExistTables) > 0 {
		msg := fmt.Sprintf("table %s not exist", strings.Join(notExistTables, ", "))
		results.add(model.TASK_ACTION_ERROR, msg)
	}
	return results, nil
}

func (i *Inspector) inspectAlterTableStmt(stmt *ast.AlterTableSpec) {

}

func (i *Inspector) inspectUseStmt(stmt *ast.UseStmt) (*InspectResults, error) {
	results := newInspectResults()
	exist, err := i.isSchemaExist(stmt.DBName)
	if err != nil {
		return results, err
	}
	if !exist {
		results.add(model.TASK_ACTION_ERROR, fmt.Sprintf("database %s not exist", stmt.DBName))
	}
	// change current schema
	i.currentSchema = stmt.DBName
	return results, nil
}
