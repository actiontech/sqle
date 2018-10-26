package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
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
	currentSchema  string
	allSchema      map[string] /*schema*/ struct{}
	schemaHasLoad  bool
	allTable       map[string] /*schema*/ map[string] /*table*/ struct{}
	DDLStmtCounter int
	DMLStmtCounter int
	createTable    string
	alterTable     string
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

func (i *Inspector) isTableExist(tableName string) (bool, error) {
	var schema = i.currentSchema
	var table = ""
	if strings.Contains(tableName, ".") {
		splitStrings := strings.SplitN(tableName, ".", 2)
		schema = splitStrings[0]
		table = splitStrings[1]
	} else {
		table = tableName
	}

	tables, hasLoad := i.allTable[schema]
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
	_, exist := tables[table]
	return exist, nil
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
			results, err = i.InspectSelectStmt(s)
		case *ast.AlterTableStmt:
			results, err = i.InspectAlterTableStmt(s)
		case *ast.UseStmt:
			results, err = i.InspectUseStmt(s)
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

func (i *Inspector) InspectSelectStmt(stmt *ast.SelectStmt) (*InspectResults, error) {
	i.DMLStmtCounter++
	results := newInspectResults()

	// check schema, table must exist
	notExistSchemas := []string{}
	notExistTables := []string{}
	tableRefs := stmt.From.TableRefs
	for _, table := range getTables(tableRefs) {
		schema := table.Schema.String()
		tableName := getTableName(table)
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return results, err
		}
		if !exist {
			notExistSchemas = append(notExistSchemas, schema)
			continue
		}
		// if schema not exist, table must not exist
		exist, err = i.isTableExist(tableName)
		if err != nil {
			return results, err
		}
		if !exist {
			notExistTables = append(notExistTables, tableName)
		}
	}
	if len(notExistSchemas) > 0 {
		msg := fmt.Sprintf("schema %s not exist", strings.Join(RemoveArrayRepeat(notExistSchemas), ", "))
		results.add(model.TASK_ACTION_ERROR, msg)
	}
	if len(notExistTables) > 0 {
		msg := fmt.Sprintf("table %s not exist", strings.Join(RemoveArrayRepeat(notExistTables), ", "))
		results.add(model.TASK_ACTION_ERROR, msg)
	}

	// where
	return results, nil
}

func (i *Inspector) InspectAlterTableStmt(stmt *ast.AlterTableSpec) (*InspectResults, error) {
	i.DDLStmtCounter++
}

func (i *Inspector) InspectUseStmt(stmt *ast.UseStmt) (*InspectResults, error) {
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
