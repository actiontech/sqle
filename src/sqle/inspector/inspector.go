package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"sqle/executor"
	"sqle/model"
	"strings"
)

type Inspector struct {
	Results     *InspectResults
	Rules       map[string]model.Rule
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

func NewInspector(rules map[string]model.Rule, db model.Instance, sqlArray []*model.CommitSql, Schema string) *Inspector {
	return &Inspector{
		Results:       newInspectResults(),
		Rules:         rules,
		Db:            db,
		currentSchema: Schema,
		SqlArray:      sqlArray,
		allSchema:     map[string]struct{}{},
		allTable:      map[string]map[string]struct{}{},
	}
}

func (i *Inspector) addResult(ruleName string, args ...interface{}) {

	// if rule is not exist, ignore save the message.
	rule, ok := i.Rules[ruleName]
	if !ok {
		return
	}
	level := rule.Level
	i.Results.add(level, ruleName, args...)
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

		stmt, err = parseOneSql(i.Db.DbType, sql.Sql)
		switch s := stmt.(type) {
		case *ast.SelectStmt:
			err = i.InspectSelectStmt(s)
		case *ast.AlterTableStmt:
			err = i.InspectAlterTableStmt(s)
		case *ast.UseStmt:
			err = i.InspectUseStmt(s)
		case *ast.CreateTableStmt:
			err = i.InspectCreateTableStmt(s)
		default:
		}
		if err != nil {
			sql.InspectStatus = model.TASK_ACTION_ERROR
			sql.InspectResult = err.Error()
			return i.SqlArray, err
		}
		sql.InspectStatus = model.TASK_ACTION_DONE
		sql.InspectLevel = i.Results.level()
		sql.InspectResult = i.Results.message()

		//clean up results
		i.Results = newInspectResults()
	}
	return i.SqlArray, nil
}

func (i *Inspector) InspectSelectStmt(stmt *ast.SelectStmt) error {
	i.DMLStmtCounter++

	// check schema, table must exist
	notExistSchemas := []string{}
	notExistTables := []string{}
	tableRefs := stmt.From.TableRefs
	for _, table := range getTables(tableRefs) {
		schema := table.Schema.String()
		tableName := getTableName(table)
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			notExistSchemas = append(notExistSchemas, schema)
			continue
		}
		// if schema not exist, table must not exist
		exist, err = i.isTableExist(tableName)
		if err != nil {
			return err
		}
		if !exist {
			notExistTables = append(notExistTables, tableName)
		}
	}
	if len(notExistSchemas) > 0 {
		i.addResult(model.SCHEMA_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistSchemas), ", "))
	}
	if len(notExistTables) > 0 {
		i.addResult(model.TABLE_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistTables), ", "))
	}

	// where
	return nil
}

func (i *Inspector) InspectCreateTableStmt(stmt *ast.CreateTableStmt) error {
	i.DDLStmtCounter++

	// check schema
	schema := stmt.Table.Schema.String()
	if schema == "" {
		schema = i.currentSchema
	}
	exist, err := i.isSchemaExist(schema)
	if err != nil {
		return err
	}
	if !exist {
		i.addResult(model.SCHEMA_NOT_EXIST, schema)

	} else {
		// check table
		tableName := getTableName(stmt.Table)
		exist, err = i.isTableExist(tableName)
		if err != nil {
			return err
		}
		if exist {
			i.addResult(model.TABLE_EXIST, tableName)
		}
	}

	// check `if not exists`
	if !stmt.IfNotExists {
		i.addResult(model.DDL_CREATE_TABLE_NOT_EXIST)
	}

	// check table length
	if len(stmt.Table.Name.String()) > 64 {
		i.addResult(model.DDL_CHECK_TABLE_NAME_LENGTH, stmt.Table.Name.String())
	}

	// check column length
	invalidColNames := []string{}
	for _, col := range stmt.Cols {
		colName := col.Name.Name.String()
		if len(colName) > 64 {
			invalidColNames = append(invalidColNames, colName)
		}
	}
	if len(invalidColNames) > 0 {
		i.addResult(model.DDL_CHECK_COLUMNS_NAME_LENGTH, strings.Join(invalidColNames, ", "))
	}

	// check primary key
	hasPk := false
	pkIsAutoIncrementUnsigned := false
	/*
		match sql like:
		CREATE TABLE  tb1 (
		a1.id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
		);
	*/
	for _, col := range stmt.Cols {
		if HasSpecialOption(col.Options, ast.ColumnOptionPrimaryKey) {
			hasPk = true
			if mysql.HasUnsignedFlag(col.Tp.Flag) && HasSpecialOption(col.Options, ast.ColumnOptionAutoIncrement) {
				pkIsAutoIncrementUnsigned = true
			}
		}

	}
	/*
		match sql like:
		CREATE TABLE  tb1 (
		a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
		PRIMARY KEY (id)
		);
	*/
	for _, constraint := range stmt.Constraints {
		if constraint.Tp == ast.ConstraintPrimaryKey {
			hasPk = true
			if len(constraint.Keys) == 1 {
				columnName := constraint.Keys[0].Column.Name.String()
				for _, col := range stmt.Cols {
					if col.Name.Name.String() == columnName {
						if mysql.HasUnsignedFlag(col.Tp.Flag) && HasSpecialOption(col.Options, ast.ColumnOptionAutoIncrement) {
							pkIsAutoIncrementUnsigned = true
						}
					}
				}
			}
		}
	}
	if !hasPk {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_EXIST)
	}
	if hasPk && !pkIsAutoIncrementUnsigned {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_TYPE)
	}

	// if char length >20 using varchar.
	for _, col := range stmt.Cols {
		if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
			i.addResult(model.DDL_CHECK_TYPE_CHAR_LENGTH)
		}
	}

	// index
	//for _,constraint:=range stmt.Constraints {
	//	if constraint.Tp == ast.ConstraintIndex {
	//
	//	}
	//	constraint.Keys
	//}
	return nil
}

func (i *Inspector) InspectAlterTableStmt(stmt *ast.AlterTableStmt) error {
	i.DDLStmtCounter++
	return nil
}

func (i *Inspector) InspectUseStmt(stmt *ast.UseStmt) error {
	exist, err := i.isSchemaExist(stmt.DBName)
	if err != nil {
		return err
	}
	if !exist {
		i.addResult(model.SCHEMA_NOT_EXIST, stmt.DBName)
	}
	// change current schema
	i.currentSchema = stmt.DBName
	return nil
}
