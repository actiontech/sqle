package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"sqle/model"
	"strings"
)

// inspector rule code
const (
	SCHEMA_NOT_EXIST                     = "schema_not_exist"
	SCHEMA_EXIST                         = "schema_exist"
	TABLE_NOT_EXIST                      = "table_not_exist"
	TABLE_EXIST                          = "table_exist"
	DDL_CREATE_TABLE_NOT_EXIST           = "ddl_create_table_not_exist"
	DDL_CHECK_OBJECT_NAME_LENGTH         = "ddl_check_object_name_length"
	DDL_CHECK_PRIMARY_KEY_EXIST          = "ddl_check_primary_key_exist"
	DDL_CHECK_PRIMARY_KEY_TYPE           = "ddl_check_primary_key_type"
	DDL_DISABLE_VARCHAR_MAX              = "ddl_disable_varchar_max"
	DDL_CHECK_TYPE_CHAR_LENGTH           = "ddl_check_type_char_length"
	DDL_DISABLE_FOREIGN_KEY              = "ddl_disable_foreign_key"
	DDL_CHECK_INDEX_COUNT                = "ddl_check_index_count"
	DDL_CHECK_COMPOSITE_INDEX_MAX        = "ddl_check_composite_index_max"
	DDL_DISABLE_USING_KEYWORD            = "ddl_disable_using_keyword"
	DDL_TABLE_USING_INNODB_UTF8MB4       = "ddl_create_table_using_innodb"
	DDL_DISABLE_INDEX_DATA_TYPE_BLOB     = "ddl_disable_index_column_blob"
	DDL_CHECK_ALTER_TABLE_NEED_MERGE     = "ddl_check_alter_table_need_merge"
	DDL_DISABLE_DROP_STATEMENT           = "ddl_disable_drop_statement"
	DML_CHECK_INVALID_WHERE_CONDITION    = "ddl_check_invalid_where_condition"
	DML_DISABE_SELECT_ALL_COLUMN         = "dml_disable_select_all_column"
	DML_MYCAT_MUST_USING_SHARDING_CLOUNM = "dml_mycat_must_using_sharding_column"
)

type RuleHandler struct {
	Rule    model.Rule
	Message string
	Func    func(*Inspect, ast.StmtNode) error
}

var (
	RuleHandlerMap = map[string]RuleHandler{}
	DefaultRules   = []model.Rule{}
)

var RuleHandlers = []RuleHandler{
	RuleHandler{
		Rule: model.Rule{
			Name:  SCHEMA_NOT_EXIST,
			Desc:  "操作数据库时，数据库必须存在",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "schema %s 不存在",
		Func:    checkObjectNotExist,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  SCHEMA_EXIST,
			Desc:  "创建数据库时，数据库不能存在",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "schema %s 已存在",
		Func:    checkObjectExist,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  TABLE_NOT_EXIST,
			Desc:  "操作表时，表必须存在",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "表 %s 不存在",
		Func:    checkObjectNotExist,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  TABLE_EXIST,
			Desc:  "创建表时，表不能存在",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "表 %s 已存在",
		Func:    checkObjectExist,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CREATE_TABLE_NOT_EXIST,
			Desc:  "新建表必须加入if not exists create，保证重复执行不报错",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "新建表必须加入if not exists create，保证重复执行不报错",
		Func:    checkIfNotExist,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_LENGTH,
			Desc:  "表名、列名、索引名的长度不能大于64字节",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "表名、列名、索引名的长度不能大于64字节",
		Func:    checkNewObjectName,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PRIMARY_KEY_EXIST,
			Desc:  "表必须有主键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "表必须有主键",
		Func:    checkPrimaryKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PRIMARY_KEY_TYPE,
			Desc:  "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
		Func:    checkPrimaryKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_VARCHAR_MAX,
			Desc:  "禁止使用 varchar(max)",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止使用 varchar(max)",
		Func:    nil,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TYPE_CHAR_LENGTH,
			Desc:  "char长度大于20时，必须使用varchar类型",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "char长度大于20时，必须使用varchar类型",
		Func:    checkStringType,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_FOREIGN_KEY,
			Desc:  "禁止使用外键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止使用外键",
		Func:    checkForeignKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COUNT,
			Desc:  "索引个数建议不超过5个",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "索引个数建议不超过5个",
		Func:    checkIndex,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COMPOSITE_INDEX_MAX,
			Desc:  "复合索引的列数量不建议超过5个",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "复合索引的列数量不建议超过5个",
		Func:    checkIndex,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_USING_KEYWORD,
			Desc:  "数据库对象命名禁止使用关键字",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "数据库对象命名禁止使用关键字 %s",
		Func:    checkNewObjectName,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_TABLE_USING_INNODB_UTF8MB4,
			Desc:  "建议使用Innodb引擎,utf8mb4字符集",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议使用Innodb引擎,utf8mb4字符集",
		Func:    checkEngineAndCharacterSet,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_INVALID_WHERE_CONDITION,
			Desc:  "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
		Func:    checkSelectWhere,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_ALTER_TABLE_NEED_MERGE,
			Desc:  "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "已存在对该表的修改语句，建议合并成一个ALTER语句",
		Func:    checkMergeAlterTable,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_DISABE_SELECT_ALL_COLUMN,
			Desc:  "不建议使用select *",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议使用select *",
		Func:    checkSelectAll,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_DROP_STATEMENT,
			Desc:  "禁止除索引外的drop操作",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止除索引外的drop操作",
		Func:    disableDropStmt,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_MYCAT_MUST_USING_SHARDING_CLOUNM,
			Desc:  "mycat dml 必须使用分片字段",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "mycat dml 必须使用分片字段",
		Func:    checkMycatShardingColumn,
	},
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
		DefaultRules = append(DefaultRules, rh.Rule)
	}
}

func checkSelectAll(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					i.addResult(DML_DISABE_SELECT_ALL_COLUMN)
				}
			}
		}
	}
	return nil
}

func checkSelectWhere(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// where condition
		if stmt.Where == nil || !whereStmtHasOneColumn(stmt.Where) {
			i.addResult(DML_CHECK_INVALID_WHERE_CONDITION)
		}
	}
	return nil
}

func checkPrimaryKey(i *Inspect, node ast.StmtNode) error {
	var hasPk = false
	var pkIsAutoIncrementBigIntUnsigned = false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check primary key
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
	default:
		return nil
	}

	if !hasPk {
		i.addResult(DDL_CHECK_PRIMARY_KEY_EXIST)
	}
	if hasPk && !pkIsAutoIncrementBigIntUnsigned {
		i.addResult(DDL_CHECK_PRIMARY_KEY_TYPE)
	}
	return nil
}

func checkMergeAlterTable(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		tableName := i.getTableName(stmt.Table)
		_, ok := i.alterTableStmts[tableName]
		if ok {
			i.addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE)
			i.alterTableStmts[tableName] = append(i.alterTableStmts[tableName], stmt)
		} else {
			i.alterTableStmts[tableName] = []*ast.AlterTableStmt{stmt}
		}
	}

	return nil
}

func checkEngineAndCharacterSet(i *Inspect, node ast.StmtNode) error {
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
	i.addResult(DDL_TABLE_USING_INNODB_UTF8MB4)
	return nil
}

func disableAddIndexForColumnsTypeBlob(i *Inspect, node ast.StmtNode) error {
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
		i.addResult(DDL_DISABLE_INDEX_DATA_TYPE_BLOB)
	}
	return nil
}

func checkNewObjectName(i *Inspect, node ast.StmtNode) error {
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
			i.addResult(DDL_CHECK_OBJECT_NAME_LENGTH)
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
		i.addResult(DDL_DISABLE_USING_KEYWORD, strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func checkForeignKey(i *Inspect, node ast.StmtNode) error {
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
		i.addResult(DDL_DISABLE_FOREIGN_KEY)
	}
	return nil
}

func checkIndex(i *Inspect, node ast.StmtNode) error {
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
		i.addResult(DDL_CHECK_INDEX_COUNT)
	}
	if compositeIndexMax > 5 {
		i.addResult(DDL_CHECK_COMPOSITE_INDEX_MAX)
	}
	return nil
}

func checkStringType(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(DDL_CHECK_TYPE_CHAR_LENGTH)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(DDL_CHECK_TYPE_CHAR_LENGTH)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkObjectExist(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check schema
		schema := i.getSchemaName(stmt.Table)
		tableName := i.getTableName(stmt.Table)
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			// if schema not exist, table must not exist
			return nil

		} else {
			// check table if schema exist
			exist, err = i.isTableExist(tableName)
			if err != nil {
				return err
			}
			if exist {
				i.addResult(TABLE_EXIST, tableName)
			}
		}
	case *ast.CreateDatabaseStmt:
		schemaName := stmt.Name
		exist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if exist {
			i.addResult(SCHEMA_EXIST, schemaName)
		}
	}
	return nil
}

func checkObjectNotExist(i *Inspect, node ast.StmtNode) error {
	var tablesName = []string{}
	var schemasName = []string{}

	switch stmt := node.(type) {
	case *ast.UseStmt:
		schemasName = append(schemasName, stmt.DBName)

	case *ast.CreateTableStmt:
		schemasName = append(schemasName, i.getSchemaName(stmt.Table))

	case *ast.AlterTableStmt:
		schemasName = append(schemasName, i.getSchemaName(stmt.Table))
		tablesName = append(tablesName, i.getTableName(stmt.Table))

	case *ast.SelectStmt:
		for _, table := range getTables(stmt.From.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}
	case *ast.InsertStmt:
		for _, table := range getTables(stmt.Table.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}

	case *ast.DeleteStmt:
		if stmt.Tables != nil && stmt.Tables.Tables != nil {
			for _, table := range stmt.Tables.Tables {
				schemasName = append(schemasName, i.getSchemaName(table))
				tablesName = append(tablesName, i.getTableName(table))
			}
		}
		for _, table := range getTables(stmt.TableRefs.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}

	case *ast.UpdateStmt:
		for _, table := range getTables(stmt.TableRefs.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}
	}

	notExistSchemas := []string{}
	for _, schema := range schemasName {
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			notExistSchemas = append(notExistSchemas, schema)
		}
	}
	if len(notExistSchemas) > 0 {
		i.addResult(SCHEMA_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistSchemas), ", "))
	}

	notExistTables := []string{}
	for _, table := range tablesName {
		exist, err := i.isTableExist(table)
		if err != nil {
			return err
		}
		if !exist {
			notExistTables = append(notExistTables, table)
		}
	}
	if len(notExistTables) > 0 {
		i.addResult(TABLE_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistTables), ", "))
	}
	return nil
}

func checkIfNotExist(i *Inspect, node ast.StmtNode) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDL_CREATE_TABLE_NOT_EXIST)
		}
	}
	return nil
}

func disableDropStmt(i *Inspect, node ast.StmtNode) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	}
	return nil
}

func checkMycatShardingColumn(i *Inspect, node ast.StmtNode) error {
	if i.Task.Instance.DbType != model.DB_TYPE_MYCAT {
		return nil
	}
	config := i.Task.Instance.MycatConfig
	hasShardingColumn := false
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		tables := getTables(stmt.Table.TableRefs)
		// tables must be one on InsertIntoStmt in parser.go
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		if stmt.Columns != nil {
			for _, column := range stmt.Columns {
				if column.Name.L == strings.ToLower(shardingCoulmn) {
					hasShardingColumn = true
				}
			}
		}
		if stmt.Setlist != nil {
			for _, set := range stmt.Setlist {
				if set.Column.Name.L == strings.ToLower(shardingCoulmn) {
					hasShardingColumn = true
				}
			}
		}
	case *ast.UpdateStmt:
		tables := getTables(stmt.TableRefs.TableRefs)
		// multi table related update not supported on mycat
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		hasShardingColumn = whereStmtHasSpecificColumn(stmt.Where, shardingCoulmn)
	case *ast.DeleteStmt:
		// not support multi table related delete
		if stmt.IsMultiTable {
			return nil
		}
		tables := getTables(stmt.TableRefs.TableRefs)
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		hasShardingColumn = whereStmtHasSpecificColumn(stmt.Where, shardingCoulmn)
	}
	if !hasShardingColumn {
		i.addResult(DML_MYCAT_MUST_USING_SHARDING_CLOUNM)
	}
	return nil
}
