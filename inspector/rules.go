package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"sqle/model"
	"strings"
)

// inspector rule code
const (
	DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST       = "ddl_check_table_without_if_not_exists"
	DDL_CHECK_OBJECT_NAME_LENGTH               = "ddl_check_object_name_length"
	DDL_CHECK_OBJECT_NAME_USING_KEYWORD        = "ddl_check_object_name_using_keyword"
	DDL_CHECK_PK_NOT_EXIST                     = "ddl_check_pk_not_exist"
	DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED       = "ddl_check_pk_without_bigint_unsigned"
	DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT        = "ddl_check_pk_without_auto_increment"
	DDL_CHECK_COLUMN_VARCHAR_MAX               = "ddl_check_column_varchar_max"
	DDL_CHECK_COLUMN_CHAR_LENGTH               = "ddl_check_column_char_length"
	DDL_DISABLE_FK                             = "ddl_disable_fk"
	DDL_CHECK_INDEX_COUNT                      = "ddl_check_index_count"
	DDL_CHECK_COMPOSITE_INDEX_MAX              = "ddl_check_composite_index_max"
	DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4     = "ddl_check_table_without_innodb_utf8mb4"
	DDL_CHECK_INDEX_COLUMN_WITH_BLOB           = "ddl_check_index_column_with_blob"
	DDL_CHECK_ALTER_TABLE_NEED_MERGE           = "ddl_check_alter_table_need_merge"
	DDL_DISABLE_DROP_STATEMENT                 = "ddl_disable_drop_statement"
	DML_CHECK_WHERE_IS_INVALID                 = "all_check_where_is_invalid"
	DML_DISABE_SELECT_ALL_COLUMN               = "dml_disable_select_all_column"
	DML_CHECK_MYCAT_WITHOUT_SHARDING_CLOUNM    = "dml_check_mycat_without_sharding_column"
	DDL_CHECK_TABLE_WITHOUT_COMMENT            = "ddl_check_table_without_comment"
	DDL_CHECK_COLUMN_WITHOUT_COMMENT           = "ddl_check_column_without_comment"
	DDL_CHECK_INDEX_PREFIX                     = "ddl_check_index_prefix"
	DDL_CHECK_UNIQUE_INDEX_PRIFIX              = "ddl_check_unique_index_prefix"
	DDL_CHECK_COLUMN_WITHOUT_DEFAULT           = "ddl_check_column_without_default"
	DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT = "ddl_check_column_timestamp_without_default"
	DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL        = "ddl_check_column_blob_with_not_null"
	DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL  = "ddl_check_column_blob_default_is_not_null"
	DML_CHECK_WITH_LIMIT                       = "dml_check_with_limit"
	DML_CHECK_WITH_ORDER_BY                    = "dml_check_with_order_by"
)

// inspector config code
const (
	CONFIG_DML_ROLLBACK_MAX_ROWS = "dml_rollback_max_rows"
	CONFIG_DDL_OSC_MIN_SIZE      = "ddl_osc_min_size"
)

type RuleHandler struct {
	Rule    model.Rule
	Message string
	Func    func(*Inspect, ast.Node) error
}

var (
	RuleHandlerMap = map[string]RuleHandler{}
	DefaultRules   = []model.Rule{}
)

var RuleHandlers = []RuleHandler{
	// config
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DML_ROLLBACK_MAX_ROWS,
			Desc:  "在 DML 语句中预计影响行数超过指定值则不回滚",
			Value: "1000",
		},
		Func: nil,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DDL_OSC_MIN_SIZE,
			Desc:  "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
			Value: "16",
		},
		Func: nil,
	},

	// rule
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST,
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
			Name:  DDL_CHECK_PK_NOT_EXIST,
			Desc:  "表必须有主键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "表必须有主键",
		Func:    checkPrimaryKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
			Desc:  "主键建议使用自增",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "主键建议使用自增",
		Func:    checkPrimaryKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED,
			Desc:  "主键建议使用 bigint 无符号类型，即 bigint unsigned",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "主键建议使用 bigint 无符号类型，即 bigint unsigned",
		Func:    checkPrimaryKey,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_VARCHAR_MAX,
			Desc:  "禁止使用 varchar(max)",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止使用 varchar(max)",
		Func:    nil,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_CHAR_LENGTH,
			Desc:  "char长度大于20时，必须使用varchar类型",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "char长度大于20时，必须使用varchar类型",
		Func:    checkStringType,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_FK,
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
			Name:  DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
			Desc:  "数据库对象命名禁止使用关键字",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "数据库对象命名禁止使用关键字 %s",
		Func:    checkNewObjectName,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4,
			Desc:  "建议使用Innodb引擎,utf8mb4字符集",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议使用Innodb引擎,utf8mb4字符集",
		Func:    checkEngineAndCharacterSet,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COLUMN_WITH_BLOB,
			Desc:  "禁止将blob类型的列加入索引",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止将blob类型的列加入索引",
		Func:    disableAddIndexForColumnsTypeBlob,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_IS_INVALID,
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
			Name:  DML_CHECK_MYCAT_WITHOUT_SHARDING_CLOUNM,
			Desc:  "mycat dml 必须使用分片字段",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "mycat dml 必须使用分片字段",
		Func:    checkMycatShardingColumn,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_COMMENT,
			Desc:  "表建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "表建议添加注释",
		Func:    checkTableWithoutComment,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_COMMENT,
			Desc:  "列建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "列建议添加注释",
		Func:    checkColumnWithoutComment,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_PREFIX,
			Desc:  "普通索引必须要以\"idx_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "普通索引必须要以\"idx_\"为前缀",
		Func:    checkIndexPrefix,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_UNIQUE_INDEX_PRIFIX,
			Desc:  "unique索引必须要以\"uniq_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "unique索引必须要以\"uniq_\"为前缀",
		Func:    checkUniqIndexPrefix,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
			Desc:  "除了自增列及大字段列之外，每个列都必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "除了自增列及大字段列之外，每个列都必须添加默认值",
		Func:    checkColumnWithoutDefault,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
			Desc:  "timestamp 类型的列必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "timestamp 类型的列必须添加默认值",
		Func:    checkColumnTimestampWithoutDefault,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		Func:    checkColumnBlobNotNull,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
		Func:    checkColumnBlobDefaultNull,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_LIMIT,
			Desc:  "delete/update 语句不能有limit条件",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "delete/update 语句不能有limit条件",
		Func:    checkDMLWithLimit,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_ORDER_BY,
			Desc:  "delete/update 语句不能有order by",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "delete/update 语句不能有order by",
		Func:    checkDMLWithOrderBy,
	},
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
		DefaultRules = append(DefaultRules, rh.Rule)
	}
}

func checkSelectAll(i *Inspect, node ast.Node) error {
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

func checkSelectWhere(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// where condition
		if stmt.Where == nil || !whereStmtHasOneColumn(stmt.Where) {
			i.addResult(DML_CHECK_WHERE_IS_INVALID)
		}
	}
	return nil
}

func checkPrimaryKey(i *Inspect, node ast.Node) error {
	var hasPk = false
	var pkIsAutoIncrement = false
	var pkIsBigIntUnsigned = false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}
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
				if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) {
					pkIsBigIntUnsigned = true
				}
				if IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					pkIsAutoIncrement = true
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
							if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) {
								pkIsBigIntUnsigned = true
							}
							if IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrement = true
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
		i.addResult(DDL_CHECK_PK_NOT_EXIST)
	}
	if hasPk && !pkIsAutoIncrement {
		i.addResult(DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT)
	}
	if hasPk && !pkIsBigIntUnsigned {
		i.addResult(DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED)
	}
	return nil
}

func checkMergeAlterTable(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := i.getTableInfo(stmt.Table)
		if exist {
			if info.alterTableStmts != nil && len(info.alterTableStmts) > 0 {
				i.addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE)
			}
		}
	}

	return nil
}

func checkEngineAndCharacterSet(i *Inspect, node ast.Node) error {
	var engine string
	var characterSet string
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.ReferTable != nil {
			return nil
		}
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
	i.addResult(DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4)
	return nil
}

func disableAddIndexForColumnsTypeBlob(i *Inspect, node ast.Node) error {
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
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
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
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
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
		i.addResult(DDL_CHECK_INDEX_COLUMN_WITH_BLOB)
	}
	return nil
}

func checkNewObjectName(i *Inspect, node ast.Node) error {
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
		i.addResult(DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
			strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func checkForeignKey(i *Inspect, node ast.Node) error {
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
		i.addResult(DDL_DISABLE_FK)
	}
	return nil
}

func checkIndex(i *Inspect, node ast.Node) error {
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
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
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
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
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

func checkStringType(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(DDL_CHECK_COLUMN_CHAR_LENGTH)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(DDL_CHECK_COLUMN_CHAR_LENGTH)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkIfNotExist(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST)
		}
	}
	return nil
}

func disableDropStmt(i *Inspect, node ast.Node) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	}
	return nil
}

func checkMycatShardingColumn(i *Inspect, node ast.Node) error {
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
	default:
		return nil
	}
	if !hasShardingColumn {
		i.addResult(DML_CHECK_MYCAT_WITHOUT_SHARDING_CLOUNM)
	}
	return nil
}

func checkTableWithoutComment(i *Inspect, node ast.Node) error {
	var tableHasComment bool
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if has refer table, sql is create table ... like ...
		if stmt.ReferTable != nil {
			return nil
		}
		if stmt.Options != nil {
			for _, option := range stmt.Options {
				if option.Tp == ast.TableOptionComment {
					tableHasComment = true
					break
				}
			}
		}
		if !tableHasComment {
			i.addResult(DDL_CHECK_TABLE_WITHOUT_COMMENT)
		}
	}
	return nil
}

func checkColumnWithoutComment(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			columnHasComment := false
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionComment {
					columnHasComment = true
				}
			}
			if !columnHasComment {
				i.addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasComment := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionComment {
						columnHasComment = true
					}
				}
				if !columnHasComment {
					i.addResult(DDL_CHECK_COLUMN_WITHOUT_COMMENT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkIndexPrefix(i *Inspect, node ast.Node) error {
	indexesName := []string{}
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex:
				indexesName = append(indexesName, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex:
				indexesName = append(indexesName, spec.Constraint.Name)
			}
		}
	default:
		return nil
	}
	for _, name := range indexesName {
		if !strings.HasPrefix(name, "idx_") {
			i.addResult(DDL_CHECK_INDEX_PREFIX)
			return nil
		}
	}
	return nil
}

func checkUniqIndexPrefix(i *Inspect, node ast.Node) error {
	uniqueIndexesName := []string{}
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniq:
				uniqueIndexesName = append(uniqueIndexesName, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintUniq:
				uniqueIndexesName = append(uniqueIndexesName, spec.Constraint.Name)
			}
		}
	default:
		return nil
	}
	for _, name := range uniqueIndexesName {
		if !strings.HasPrefix(name, "uniq_") {
			i.addResult(DDL_CHECK_UNIQUE_INDEX_PRIFIX)
			return nil
		}
	}
	return nil
}

func checkColumnWithoutDefault(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col == nil {
				continue
			}
			isAutoIncrementColumn := false
			isBlobColumn := false
			columnHasDefault := false
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionAutoIncrement {
					isAutoIncrementColumn = true
				}

				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob:
					isBlobColumn = true
				}

				if option.Tp == ast.ColumnOptionDefaultValue {
					columnHasDefault = true
				}
			}
			if isAutoIncrementColumn || isBlobColumn {
				continue
			}
			if !columnHasDefault {
				i.addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasDefault := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionDefaultValue {
						columnHasDefault = true
					}
				}
				if !columnHasDefault {
					i.addResult(DDL_CHECK_COLUMN_WITHOUT_DEFAULT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnTimestampWithoutDefault(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			columnHasDefault := false
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionDefaultValue {
					columnHasDefault = true
				}
			}
			if !columnHasDefault && (col.Tp.Tp == mysql.TypeTimestamp || col.Tp.Tp == mysql.TypeDatetime) {
				i.addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasDefault := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionDefaultValue {
						columnHasDefault = true
					}
				}
				if !columnHasDefault && (col.Tp.Tp == mysql.TypeTimestamp || col.Tp.Tp == mysql.TypeDatetime) {
					i.addResult(DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnBlobNotNull(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			switch col.Tp.Tp {
			case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionNotNull {
						i.addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionNotNull {
							i.addResult(DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkColumnBlobDefaultNull(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			switch col.Tp.Tp {
			case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
				for _, opt := range col.Options {
					if opt.Tp == ast.ColumnOptionDefaultValue && opt.Expr.GetType().Tp != mysql.TypeNull {
						i.addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableModifyColumn, ast.AlterTableAlterColumn,
			ast.AlterTableChangeColumn, ast.AlterTableAddColumns) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionDefaultValue && opt.Expr.GetType().Tp != mysql.TypeNull {
							i.addResult(DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkDMLWithLimit(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			i.addResult(DML_CHECK_WITH_LIMIT)
		}
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			i.addResult(DML_CHECK_WITH_LIMIT)
		}
	}
	return nil
}

func checkDMLWithOrderBy(i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			i.addResult(DML_CHECK_WITH_ORDER_BY)
		}
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			i.addResult(DML_CHECK_WITH_ORDER_BY)
		}
	}
	return nil
}
