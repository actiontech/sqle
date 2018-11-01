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
	createTableStmts map[string] /*schema.table*/ *ast.CreateTableStmt

	// save alter table parse object from input sql;
	alterTableStmts map[string] /*schema.table*/ []*ast.AlterTableStmt
	rollbackSqls    []string
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
		alterTableStmts:  map[string][]*ast.AlterTableStmt{},
		rollbackSqls:     []string{},
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
	t, err := parseOneSql(i.Db.DbType, sql)
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

func (i *Inspector) Inspect() ([]*model.CommitSql, error) {
	defer i.closeDbConn()

	var inspectFns = []func(stmt ast.StmtNode) error{
		i.checkEngineAndCharacterSet,
		i.disableAddIndexForColumnsTypeBlob,
		i.checkNewObjectName,
		i.checkForeignKey,
		i.checkIndex,
		i.checkStringType,
	}

	for _, sql := range i.SqlArray {
		var node ast.StmtNode
		var err error

		node, err = parseOneSql(i.Db.DbType, sql.Sql)
		switch node.(type) {
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

		err = i.InspectSpecificStmt(node)
		if err != nil {
			return nil, err
		}

		// base check
		for _, fn := range inspectFns {
			if err := fn(node); err != nil {
				return nil, err
			}
		}

		sql.InspectStatus = model.TASK_ACTION_DONE
		sql.InspectLevel = i.Results.level()
		sql.InspectResult = i.Results.message()

		//clean up results
		i.Results = newInspectResults()
	}
	return i.SqlArray, nil
}

func (i *Inspector) InspectSpecificStmt(node ast.StmtNode) error {
	// specific check
	switch s := node.(type) {
	case *ast.SelectStmt:
		return i.inspectSelectStmt(s)
	case *ast.AlterTableStmt:
		return i.inspectAlterTableStmt(s)
	case *ast.UseStmt:
		return i.inspectUseStmt(s)
	case *ast.CreateTableStmt:
		return i.inspectCreateTableStmt(s)
	case *ast.CreateDatabaseStmt:
		return i.inspectCreateSchemaStmt(s)
	case *ast.DropDatabaseStmt:
		delete(i.allSchema, s.Name)
		i.addResult(model.DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		for _, table := range s.Tables {
			delete(i.alterTableStmts, i.getTableName(table))
		}
		i.addResult(model.DDL_DISABLE_DROP_STATEMENT)
	default:
		return nil
	}
	return nil
}

func (i *Inspector) inspectSelectStmt(stmt *ast.SelectStmt) error {
	// check schema, table must exist
	notExistSchemas := []string{}
	notExistTables := []string{}
	tableRefs := stmt.From.TableRefs
	for _, table := range getTables(tableRefs) {
		schema := i.getSchemaName(table)
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			notExistSchemas = append(notExistSchemas, schema)
			continue
		}
		// if schema not exist, table must not exist
		tableName := i.getTableName(table)
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

	// check select all column
	if stmt.Fields != nil && stmt.Fields.Fields != nil {
		for _, field := range stmt.Fields.Fields {
			if field.WildCard != nil {
				i.addResult(model.DML_DISABE_SELECT_ALL_COLUMN)
			}
		}
	}
	// where condition
	if stmt.Where == nil || !scanWhereColumn(stmt.Where) {
		i.addResult(model.DML_CHECK_INVALID_WHERE_CONDITION)
	}
	return nil
}

func (i *Inspector) inspectCreateSchemaStmt(stmt *ast.CreateDatabaseStmt) error {
	schemaName := stmt.Name
	exist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if exist {
		i.addResult(model.SCHEMA_EXIST)
	}
	i.allSchema[schemaName] = struct{}{}
	return nil
}

func (i *Inspector) inspectCreateTableStmt(stmt *ast.CreateTableStmt) error {
	// check schema
	schema := i.getSchemaName(stmt.Table)
	exist, err := i.isSchemaExist(schema)
	if err != nil {
		return err
	}
	if !exist {
		i.addResult(model.SCHEMA_NOT_EXIST, schema)

	} else {
		// check table need not exist
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

	// check primary key
	var hasPk = false
	var pkIsAutoIncrementBigIntUnsigned = false
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
	if !hasPk {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_EXIST)
	}
	if hasPk && !pkIsAutoIncrementBigIntUnsigned {
		i.addResult(model.DDL_CHECK_PRIMARY_KEY_TYPE)
	}
	return nil
}

func (i *Inspector) inspectAlterTableStmt(stmt *ast.AlterTableStmt) error {
	// check schema
	schema := i.getSchemaName(stmt.Table)
	tableName := i.getTableName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schema)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.addResult(model.SCHEMA_NOT_EXIST, schema)

	} else {
		// check table, need exist
		tableExist, err := i.isTableExist(tableName)
		if err != nil {
			return err
		}
		if !tableExist {
			i.addResult(model.TABLE_NOT_EXIST, tableName)
		}
	}

	// merge alter table
	_, ok := i.alterTableStmts[tableName]
	if ok {
		i.addResult(model.DDL_CHECK_ALTER_TABLE_NEED_MERGE)
		i.alterTableStmts[tableName] = append(i.alterTableStmts[tableName], stmt)
	} else {
		i.alterTableStmts[tableName] = []*ast.AlterTableStmt{stmt}
	}
	return nil
}

func (i *Inspector) inspectUseStmt(stmt *ast.UseStmt) error {
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

func (i *Inspector) checkEngineAndCharacterSet(node ast.StmtNode) error {
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

func (i *Inspector) disableAddIndexForColumnsTypeBlob(node ast.StmtNode) error {
	indexColumns := map[string]struct{}{}
	isTypeBlobCols := map[string]bool{}
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
					break
				}
			}
			if _, ok := indexColumns[col.Name.Name.String()]; ok {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					indexDataTypeIsBlob = true
					break
				}
			}
		}
	case *ast.AlterTableStmt:
		// collect index column
		for _, spec := range stmt.Specs {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
					indexColumns[col.Name.Name.String()] = struct{}{}
				}
			}
			if spec.Constraint != nil {
				switch spec.Constraint.Tp {
				case ast.ConstraintKey, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
					for _, col := range spec.Constraint.Keys {
						indexColumns[col.Column.Name.String()] = struct{}{}
					}
				}
			}
		}
		if len(indexColumns) <= 0 {
			return nil
		}

		// collect columns type
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		for _, spec := range stmt.Specs {
			if spec.NewColumns != nil {
				for _, col := range spec.NewColumns {
					if MysqlDataTypeIsBlob(col.Tp.Tp) {
						isTypeBlobCols[col.Name.Name.String()] = true
					} else {
						isTypeBlobCols[col.Name.Name.String()] = false
					}
				}
			}
		}
		// check index columns string type
		for colName, _ := range indexColumns {
			if isTypeBlobCols[colName] {
				indexDataTypeIsBlob = true
				break
			}
		}
	case *ast.CreateIndexStmt:
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil || !exist {
			return err
		}
		for _, col := range createTableStmt.Cols {
			if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) && MysqlDataTypeIsBlob(col.Tp.Tp) {
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, indexColumns := range stmt.IndexColNames {
			if isTypeBlobCols[indexColumns.Column.Name.String()] {
				indexDataTypeIsBlob = true
				break
			}
		}
	default:
		return nil
	}
	if indexDataTypeIsBlob {
		i.addResult(model.DDL_DISABLE_INDEX_DATA_TYPE_BLOB)
	}
	return nil
}

func (i *Inspector) checkNewObjectName(node ast.StmtNode) error {
	names := []string{}
	invalidNames := []string{}

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
			case ast.AlterTableRenameIndex:
				names = append(names, spec.ToKey.String())
			}
		}
	case *ast.CreateIndexStmt:
		names = append(names, stmt.IndexName)
	default:
		return nil
	}

	// check length
	for _, name := range names {
		if len(name) > 64 {
			i.addResult(model.DDL_CHECK_OBJECT_NAME_LENGTH)
			break
		}
	}
	// check keyword
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

func (i *Inspector) checkForeignKey(node ast.StmtNode) error {
	hasFk := false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	default:
		return nil
	}
	if hasFk {
		i.addResult(model.DDL_DISABLE_FOREIGN_KEY)
	}
	return nil
}

func (i *Inspector) checkIndex(node ast.StmtNode) error {
	indexCounter := 0
	compositeIndexMax := 0

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(constraint.Keys) {
					compositeIndexMax = len(constraint.Keys)
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint == nil {
				continue
			}
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(spec.Constraint.Keys) {
					compositeIndexMax = len(spec.Constraint.Keys)
				}
			}
		}
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}

	case *ast.CreateIndexStmt:
		indexCounter++
		if compositeIndexMax < len(stmt.IndexColNames) {
			compositeIndexMax = len(stmt.IndexColNames)
		}
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}
	default:
		return nil
	}
	if indexCounter > 5 {
		i.addResult(model.DDL_CHECK_INDEX_COUNT)
	}
	if compositeIndexMax > 5 {
		i.addResult(model.DDL_CHECK_COMPOSITE_INDEX_MAX)
	}
	return nil
}

func (i *Inspector) checkStringType(node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(model.DDL_CHECK_TYPE_CHAR_LENGTH)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(model.DDL_CHECK_TYPE_CHAR_LENGTH)
				}
			}
		}
	default:
		return nil
	}
	return nil
}
