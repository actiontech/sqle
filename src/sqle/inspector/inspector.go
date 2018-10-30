package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"sqle/executor"
	"sqle/model"
	"strings"
)

var (
	SQL_STMT_CONFLICT_ERROR = fmt.Errorf("不能同时提交 DDL 和 DML 语句")
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

	// save create table parser object from db by query "show create table tb_1";
	// using in inspect and generate rollback sql
	createTableStmts map[string]*ast.CreateTableStmt
	alterTable       string
}

func NewInspector(rules map[string]model.Rule, db model.Instance, sqlArray []*model.CommitSql, Schema string) *Inspector {
	return &Inspector{
		Results:          newInspectResults(),
		Rules:            rules,
		Db:               db,
		currentSchema:    Schema,
		SqlArray:         sqlArray,
		allSchema:        map[string]struct{}{},
		allTable:         map[string]map[string]struct{}{},
		createTableStmts: map[string]*ast.CreateTableStmt{},
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

func (i *Inspector) getTableName(stmt *ast.TableName) string {
	var schema string
	if stmt.Schema.String() == "" {
		schema = i.currentSchema
	} else {
		schema = stmt.Schema.String()
	}
	return fmt.Sprintf("%s.%s", schema, stmt.Name)
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

func (i *Inspector) getCreateTableStmt(tableName string) (*ast.CreateTableStmt, error) {

	// check local memory first, for uint test
	createStmt, ok := i.createTableStmts[tableName]
	if ok {
		return createStmt, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	sql, err := conn.ShowCreateTable(tableName)
	if err != nil {
		return nil, err
	}
	t, err := parseOneSql(i.Db.DbType, sql)
	if err != nil {
		return nil, err
	}
	createStmt, ok = t.(*ast.CreateTableStmt)
	if !ok {
		return nil, fmt.Errorf("stmt not support")
	}
	i.createTableStmts[tableName] = createStmt
	return createStmt, nil
}

func (i *Inspector) Inspect() ([]*model.CommitSql, error) {
	defer i.closeDbConn()

	for _, sql := range i.SqlArray {
		var stmt ast.StmtNode
		var err error

		stmt, err = parseOneSql(i.Db.DbType, sql.Sql)
		switch stmt.(type) {
		case ast.DDLNode:
			if i.DMLStmtCounter > 0 {
				return nil, SQL_STMT_CONFLICT_ERROR
			}
			i.DDLStmtCounter++
		case ast.DMLNode:
			if i.DDLStmtCounter > 0 {
				return nil, SQL_STMT_CONFLICT_ERROR
			}
		}

		// base check
		i.CheckObjectNameUsingKeyword(stmt)
		i.CheckEngineAndCharacterSet(stmt)
		i.DisableAddIndexForColumnsTypeBlob(stmt)
		i.CheckObjectNameLength(stmt)
		err = i.CheckPrimaryKey(stmt)
		if err != nil {
			return nil, err
		}
		// specific check
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
		tableName := i.getTableName(table)
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
		tableName := i.getTableName(stmt.Table)
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

	// if char length >20 using varchar.
	for _, col := range stmt.Cols {
		if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
			i.addResult(model.DDL_CHECK_TYPE_CHAR_LENGTH)
		}
	}

	// check foreign key
	hasRefer := false
	for _, constraint := range stmt.Constraints {
		if constraint.Refer != nil {
			hasRefer = true
			break
		}
	}
	if hasRefer {
		i.addResult(model.DDL_DISABLE_FOREIGN_KEY)
	}

	// check index
	// TODO: include keyword "KEY" "UNIQUE KEY"
	indexCounter := 0
	compositeIndexMax := 0
	for _, constraint := range stmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintIndex, ast.ConstraintUniqIndex:
			indexCounter++
			if compositeIndexMax < len(constraint.Keys) {
				compositeIndexMax = len(constraint.Keys)
			}
		}
	}
	if indexCounter > 5 {
		i.addResult(model.DDL_CHECK_INDEX_COUNT)
	}
	if compositeIndexMax > 5 {
		i.addResult(model.DDL_CHECK_COMPOSITE_INDEX_MAX)
	}
	return nil
}

func (i *Inspector) InspectAlterTableStmt(stmt *ast.AlterTableStmt) error {
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

func (i *Inspector) CheckObjectNameUsingKeyword(stmt ast.StmtNode) (err error) {
	names := []string{}
	invalidNames := []string{}

	// collect object name
	switch s := stmt.(type) {
	case *ast.CreateDatabaseStmt:
		// schema
		names = append(names, s.Name)
	case *ast.CreateTableStmt:
		// table
		names = append(names, s.Table.Name.String())
		for _, col := range s.Cols {
			names = append(names, col.Name.Name.String())
		}
		//index
		for _, constraint := range s.Constraints {
			if constraint.Name != "" {
				names = append(names, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		// table
		names = append(names, s.Table.Name.String())
		for _, spec := range s.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns, ast.AlterTableChangeColumn:
				for _, col := range spec.NewColumns {
					// column
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableRenameTable:
				// table
				names = append(names, spec.NewTable.Name.String())
			case ast.AlterTableRenameIndex:
				// index
				names = append(names, spec.ToKey.String())
			case ast.AlterTableAddConstraint:
				if spec.Constraint.Name != "" {
					names = append(names, spec.Constraint.Name)
				}
			}
		}
	case *ast.CreateIndexStmt:
		// index
		names = append(names, s.IndexName)
	}

	// filter object name
	for _, name := range names {
		if IsMysqlReservedKeyword(name) {
			invalidNames = append(invalidNames, name)
		}
	}

	if len(invalidNames) > 0 {
		i.addResult(model.DDL_DISABLE_USING_KEYWORD, strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func (i *Inspector) CheckEngineAndCharacterSet(node ast.StmtNode) error {
	var engine string
	var characterSet string
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, op := range stmt.Options {
			switch op.Tp {
			case ast.TableOptionEngine:
				engine = op.StrValue
			case ast.TableOptionCharset:
				characterSet = op.StrValue
			}
		}
	default:
		return nil
	}
	if strings.ToLower(engine) == "innodb" && strings.ToLower(characterSet) == "utf8mb4" {
		return nil
	}
	i.addResult(model.DDL_TABLE_USING_INNODB_UTF8MB4)
	return nil
}

func (i *Inspector) DisableAddIndexForColumnsTypeBlob(node ast.StmtNode) error {
	indexColumns := map[string]struct{}{}
	indexDataTypeIsBlob := false
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				for _, col := range constraint.Keys {
					indexColumns[col.Column.Name.String()] = struct{}{}
				}
			}
		}
		for _, col := range stmt.Cols {
			if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					indexDataTypeIsBlob = true
				}
			}
			if _, ok := indexColumns[col.Name.Name.String()]; ok {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					indexDataTypeIsBlob = true
				}
			}
		}
	}
	if indexDataTypeIsBlob {
		i.addResult(model.DDL_DISABLE_INDEX_DATA_TYPE_BLOB)
	}
	return nil
}

func (i *Inspector) CheckObjectNameLength(node ast.StmtNode) error {
	names := []string{}
	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		// schema
		names = append(names, stmt.Name)
	case *ast.CreateTableStmt:

		// table
		names = append(names, stmt.Table.Name.String())

		// column
		for _, col := range stmt.Cols {
			names = append(names, col.Name.Name.String())
		}
		// index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniqKey, ast.ConstraintKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
				names = append(names, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableRenameTable:
				// rename table
				names = append(names, spec.NewTable.Name.String())
			case ast.AlterTableAddColumns:
				// new column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableChangeColumn:
				// rename column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableAddConstraint:
				// if spec.Constraint.Name not index name, it will be null
				names = append(names, spec.Constraint.Name)
			}
		}
	case *ast.CreateIndexStmt:
		names = append(names, stmt.IndexName)
	}

	for _, name := range names {
		if len(name) > 64 {
			i.addResult(model.DDL_CHECK_OBJECT_NAME_LENGTH)
			return nil
		}
	}
	return nil
}

// CheckPrimaryKey used for "create table stmt" and "alter table stmt".
func (i *Inspector) CheckPrimaryKey(node ast.StmtNode) error {
	var hasPk = false
	var pkIsAutoIncrementBigIntUnsigned = false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// TODO: tidb parser not support keyword for SERIAL; it is a alias for "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE"
		/*
			match sql like:
			CREATE TABLE  tb1 (
			a1.id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
			);
		*/
		for _, col := range stmt.Cols {
			if IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
				hasPk = true
				if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
					IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					pkIsAutoIncrementBigIntUnsigned = true
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
							if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
								IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrementBigIntUnsigned = true
							}
						}
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		var newColumns = map[string]*ast.ColumnDef{}
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				for _, col := range spec.NewColumns {
					//
					newColumns[col.Name.Name.String()] = col
					if IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
						hasPk = true
						if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
							IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
							pkIsAutoIncrementBigIntUnsigned = true
						}
					}
				}
			case ast.AlterTableAddConstraint:
				if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
					hasPk = true
					for _, colName := range spec.Constraint.Keys {
						if col, ok := newColumns[colName.Column.Name.String()]; ok {
							if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
								IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrementBigIntUnsigned = true
							}
						} else {
							// if column not exist in new columns, column will exists in old table
							tableName := i.getTableName(stmt.Table)
							tableExist, err := i.isTableExist(tableName)
							if err != nil || !tableExist {
								return err
							}
							createTableStmt, err := i.getCreateTableStmt(tableName)
							for _, col := range createTableStmt.Cols {
								if colName.Column.Name.String() != col.Name.Name.String() {
									continue
								}
								if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
									IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
									pkIsAutoIncrementBigIntUnsigned = true
								}
							}
						}
					}
				}
			}
		}
	default:
		return nil
	}
	if !hasPk {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_EXIST)
	}
	if hasPk && !pkIsAutoIncrementBigIntUnsigned {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_TYPE)
	}
	return nil
}
