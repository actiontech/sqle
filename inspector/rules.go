package inspector

import (
	"github.com/pingcap/tidb/ast"
	"sqle/model"
)

const (
	RULE_LEVEL_NORMAL = "normal"
	RULE_LEVEL_NOTICE = "notice"
	RULE_LEVEL_WARN   = "warn"
	RULE_LEVEL_ERROR  = "error"
)

var RuleLevelMap = map[string]int{
	RULE_LEVEL_NORMAL: 0,
	RULE_LEVEL_NOTICE: 1,
	RULE_LEVEL_WARN:   2,
	RULE_LEVEL_ERROR:  3,
}

// inspector rule code
const (
	SCHEMA_NOT_EXIST                  = "schema_not_exist"
	SCHEMA_EXIST                      = "schema_exist"
	TABLE_NOT_EXIST                   = "table_not_exist"
	TABLE_EXIST                       = "table_exist"
	DDL_CREATE_TABLE_NOT_EXIST        = "ddl_create_table_not_exist"
	DDL_CHECK_OBJECT_NAME_LENGTH      = "ddl_check_object_name_length"
	DDL_CHECK_PRIMARY_KEY_EXIST       = "ddl_check_primary_key_exist"
	DDL_CHECK_PRIMARY_KEY_TYPE        = "ddl_check_primary_key_type"
	DDL_DISABLE_VARCHAR_MAX           = "ddl_disable_varchar_max"
	DDL_CHECK_TYPE_CHAR_LENGTH        = "ddl_check_type_char_length"
	DDL_DISABLE_FOREIGN_KEY           = "ddl_disable_foreign_key"
	DDL_CHECK_INDEX_COUNT             = "ddl_check_index_count"
	DDL_CHECK_COMPOSITE_INDEX_MAX     = "ddl_check_composite_index_max"
	DDL_DISABLE_USING_KEYWORD         = "ddl_disable_using_keyword"
	DDL_TABLE_USING_INNODB_UTF8MB4    = "ddl_create_table_using_innodb"
	DDL_DISABLE_INDEX_DATA_TYPE_BLOB  = "ddl_disable_index_column_blob"
	DDL_CHECK_ALTER_TABLE_NEED_MERGE  = "ddl_check_alter_table_need_merge"
	DDL_DISABLE_DROP_STATEMENT        = "ddl_disable_drop_statement"
	DML_CHECK_INVALID_WHERE_CONDITION = "ddl_check_invalid_where_condition"
	DML_DISABE_SELECT_ALL_COLUMN      = "dml_disable_select_all_column"
)

var DefaultRules = []model.Rule{
	model.Rule{
		Name:    SCHEMA_NOT_EXIST,
		Desc:    "操作数据库时，数据库必须存在",
		Message: "schema %s 不存在",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    SCHEMA_EXIST,
		Desc:    "创建数据库时，数据库不能存在",
		Message: "schema %s 已存在",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    TABLE_NOT_EXIST,
		Desc:    "操作表时，表必须存在",
		Message: "表 %s 不存在",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    TABLE_EXIST,
		Desc:    "创建表时，表不能存在",
		Message: "表 %s 已存在",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CREATE_TABLE_NOT_EXIST,
		Desc:    "新建表必须加入if not exists create，保证重复执行不报错",
		Message: "新建表必须加入if not exists create，保证重复执行不报错",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_OBJECT_NAME_LENGTH,
		Desc:    "表名、列名、索引名的长度不能大于64字节",
		Message: "表名、列名、索引名的长度不能大于64字节",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_PRIMARY_KEY_EXIST,
		Desc:    "表必须有主键",
		Message: "表必须有主键",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_PRIMARY_KEY_TYPE,
		Desc:    "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
		Message: "主键建议使用自增，且为bigint无符号类型，即bigint unsigned",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_DISABLE_VARCHAR_MAX,
		Desc:    "禁止使用 varchar(max)",
		Message: "禁止使用 varchar(max)",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_TYPE_CHAR_LENGTH,
		Desc:    "char长度大于20时，必须使用varchar类型",
		Message: "char长度大于20时，必须使用varchar类型",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_DISABLE_FOREIGN_KEY,
		Desc:    "禁止使用外键",
		Message: "禁止使用外键",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_INDEX_COUNT,
		Desc:    "索引个数建议不超过5个",
		Message: "索引个数建议不超过5个",
		Level:   RULE_LEVEL_NOTICE,
	},
	model.Rule{
		Name:    DDL_CHECK_COMPOSITE_INDEX_MAX,
		Desc:    "复合索引的列数量不建议超过5个",
		Message: "复合索引的列数量不建议超过5个",
		Level:   RULE_LEVEL_NOTICE,
	},
	model.Rule{
		Name:    DDL_DISABLE_USING_KEYWORD,
		Desc:    "数据库对象命名禁止使用关键字",
		Message: "数据库对象命名禁止使用关键字 %s",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_TABLE_USING_INNODB_UTF8MB4,
		Desc:    "建议使用Innodb引擎,utf8mb4字符集",
		Message: "建议使用Innodb引擎,utf8mb4字符集",
		Level:   RULE_LEVEL_NOTICE,
	},
	model.Rule{
		Name:    DDL_DISABLE_INDEX_DATA_TYPE_BLOB,
		Desc:    "禁止将blob类型的列加入索引",
		Message: "禁止将blob类型的列加入索引",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DML_CHECK_INVALID_WHERE_CONDITION,
		Desc:    "必须使用有效的 where 条件查询",
		Message: "schema %s 已存在",
		Level:   RULE_LEVEL_ERROR,
	},
	model.Rule{
		Name:    DDL_CHECK_ALTER_TABLE_NEED_MERGE,
		Desc:    "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
		Message: "已存在对该表的修改语句，建议合并成一个ALTER语句",
		Level:   RULE_LEVEL_NOTICE,
	},
	model.Rule{
		Name:    DML_DISABE_SELECT_ALL_COLUMN,
		Desc:    "不建议使用select *",
		Message: "不建议使用select *",
		Level:   RULE_LEVEL_NOTICE,
	},
	model.Rule{
		Name:    DDL_DISABLE_DROP_STATEMENT,
		Desc:    "禁止除索引外的drop操作",
		Message: "禁止除索引外的drop操作",
		Level:   RULE_LEVEL_ERROR,
	},
}

func (i *Inspector) initRulesFunc() {
	i.RulesFunc = map[string]func(stmt ast.StmtNode, rule string) error{
		SCHEMA_NOT_EXIST:                  i.checkObjectNotExist,
		TABLE_NOT_EXIST:                   i.checkObjectNotExist,
		SCHEMA_EXIST:                      i.checkObjectExist,
		TABLE_EXIST:                       i.checkObjectExist,
		DDL_CREATE_TABLE_NOT_EXIST:        i.checkIfNotExist,
		DDL_CHECK_OBJECT_NAME_LENGTH:      i.checkNewObjectName,
		DDL_CHECK_PRIMARY_KEY_EXIST:       i.checkPrimaryKey,
		DDL_CHECK_PRIMARY_KEY_TYPE:        i.checkPrimaryKey,
		DDL_DISABLE_VARCHAR_MAX:           nil,
		DDL_CHECK_TYPE_CHAR_LENGTH:        i.checkStringType,
		DDL_DISABLE_FOREIGN_KEY:           i.checkForeignKey,
		DDL_CHECK_INDEX_COUNT:             i.checkIndex,
		DDL_CHECK_COMPOSITE_INDEX_MAX:     i.checkIndex,
		DDL_DISABLE_USING_KEYWORD:         i.checkNewObjectName,
		DDL_TABLE_USING_INNODB_UTF8MB4:    i.checkEngineAndCharacterSet,
		DDL_DISABLE_INDEX_DATA_TYPE_BLOB:  i.disableAddIndexForColumnsTypeBlob,
		DDL_CHECK_ALTER_TABLE_NEED_MERGE:  i.checkMergeAlterTable,
		DDL_DISABLE_DROP_STATEMENT:        i.disableDropStmt,
		DML_CHECK_INVALID_WHERE_CONDITION: i.checkSelectWhere,
		DML_DISABE_SELECT_ALL_COLUMN:      i.checkSelectAll,
	}
}
