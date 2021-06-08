package inspector

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"actiontech.cloud/sqle/sqle/sqle/executor"
	"actiontech.cloud/universe/ucommon/v4/util"

	driver "github.com/pingcap/tidb/types/parser_driver"

	"actiontech.cloud/sqle/sqle/sqle/model"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

// inspector DDL rules
const (
	DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST             = "ddl_check_table_without_if_not_exists"
	DDL_CHECK_OBJECT_NAME_LENGTH                     = "ddl_check_object_name_length"
	DDL_CHECK_OBJECT_NAME_USING_KEYWORD              = "ddl_check_object_name_using_keyword"
	DDL_CHECK_PK_NOT_EXIST                           = "ddl_check_pk_not_exist"
	DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED             = "ddl_check_pk_without_bigint_unsigned"
	DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT              = "ddl_check_pk_without_auto_increment"
	DDL_CHECK_COLUMN_VARCHAR_MAX                     = "ddl_check_column_varchar_max"
	DDL_CHECK_COLUMN_CHAR_LENGTH                     = "ddl_check_column_char_length"
	DDL_DISABLE_FK                                   = "ddl_disable_fk"
	DDL_CHECK_INDEX_COUNT                            = "ddl_check_index_count"
	DDL_CHECK_COMPOSITE_INDEX_MAX                    = "ddl_check_composite_index_max"
	DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4           = "ddl_check_table_without_innodb_utf8mb4"
	DDL_CHECK_INDEX_COLUMN_WITH_BLOB                 = "ddl_check_index_column_with_blob"
	DDL_CHECK_ALTER_TABLE_NEED_MERGE                 = "ddl_check_alter_table_need_merge"
	DDL_DISABLE_DROP_STATEMENT                       = "ddl_disable_drop_statement"
	DDL_CHECK_TABLE_WITHOUT_COMMENT                  = "ddl_check_table_without_comment"
	DDL_CHECK_COLUMN_WITHOUT_COMMENT                 = "ddl_check_column_without_comment"
	DDL_CHECK_INDEX_PREFIX                           = "ddl_check_index_prefix"
	DDL_CHECK_UNIQUE_INDEX_PRIFIX                    = "ddl_check_unique_index_prefix"
	DDL_CHECK_UNIQUE_INDEX                           = "ddl_check_unique_index"
	DDL_CHECK_COLUMN_WITHOUT_DEFAULT                 = "ddl_check_column_without_default"
	DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT       = "ddl_check_column_timestamp_without_default"
	DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL              = "ddl_check_column_blob_with_not_null"
	DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL        = "ddl_check_column_blob_default_is_not_null"
	DDL_CHECK_COLUMN_ENUM_NOTICE                     = "ddl_check_column_enum_notice"
	DDL_CHECK_COLUMN_SET_NOTICE                      = "ddl_check_column_set_notice"
	DDL_CHECK_COLUMN_BLOB_NOTICE                     = "ddl_check_column_blob_notice"
	DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT             = "ddl_check_pk_prohibit_auto_increment"
	DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS = "ddl_check_indexes_exist_before_creat_constraints"
	DDL_CHECK_COLLATION_DATABASE                     = "ddl_check_collation_database"
	DDL_CHECK_DECIMAL_TYPE_COLUMN                    = "ddl_check_decimal_type_column"
	DDL_CHECK_DATABASE_SUFFIX                        = "ddl_check_database_suffix"
	DDL_CHECK_PK_NAME                                = "ddl_check_pk_name"
	DDL_CHECK_TRANSACTION_ISOLATION_LEVEL            = "ddl_check_transaction_isolation_level"
	DDL_CHECK_TABLE_PARTITION                        = "ddl_check_table_partition"
	DDL_CHECK_IS_EXIST_LIMIT_OFFSET                  = "ddl_check_is_exist_limit_offset"
	DDL_CHECK_INDEX_OPTION                           = "ddl_check_index_option"
	DDL_CHECK_OBJECT_NAME_USING_CN                   = "ddl_check_object_name_using_cn"
)

// inspector DML rules
const (
	DML_CHECK_WITH_LIMIT                      = "dml_check_with_limit"
	DML_CHECK_WITH_ORDER_BY                   = "dml_check_with_order_by"
	DML_CHECK_WHERE_IS_INVALID                = "all_check_where_is_invalid"
	DML_DISABE_SELECT_ALL_COLUMN              = "dml_disable_select_all_column"
	DML_CHECK_INSERT_COLUMNS_EXIST            = "dml_check_insert_columns_exist"
	DML_CHECK_BATCH_INSERT_LISTS_MAX          = "dml_check_batch_insert_lists_max"
	DML_CHECK_WHERE_EXIST_FUNC                = "dml_check_where_exist_func"
	DML_CHECK_WHERE_EXIST_NOT                 = "dml_check_where_exist_not"
	DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION = "dml_check_where_exist_implicit_conversion"
	DML_CHECK_LIMIT_MUST_EXIST                = "dml_check_limit_must_exist"
	DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES  = "dml_check_where_exist_scalar_sub_queries"
	DML_CHECK_WHERE_EXIST_NULL                = "dml_check_where_exist_null"
	DML_CHECK_SELECT_FOR_UPDATE               = "dml_check_select_for_update"
	DML_CHECK_NEEDLESS_FUNC                   = "dml_check_needless_func"
	DML_CHECK_FUZZY_SEARCH                    = "dml_check_fuzzy_search"
	DML_CHECK_NUMBER_OF_JOIN_TABLES           = "dml_check_number_of_join_tables"
	DML_CHECK_IS_AFTER_UNION_DISTINCT         = "dml_check_is_after_union_distinct"
	DMLCheckExplainAccessTypeAll              = "dml_check_explain_access_type_all"
	DMLCheckExplainExtraUsingFilesort         = "dml_check_explain_extra_using_filesort"
	DMLCheckExplainExtraUsingTemporary        = "dml_check_explain_extra_using_temporary"
)

// inspector config code
const (
	CONFIG_DML_ROLLBACK_MAX_ROWS = "dml_rollback_max_rows"
	CONFIG_DDL_OSC_MIN_SIZE      = "ddl_osc_min_size"
)

type RuleHandler struct {
	Rule          model.Rule
	Message       string
	Func          func(model.Rule, *Inspect, ast.Node) error
	IsDefaultRule bool
}

var (
	RuleHandlerMap       = map[string]RuleHandler{}
	DefaultTemplateRules = []model.Rule{}
	InitRules            = []model.Rule{}
)

var RuleHandlers = []RuleHandler{
	// config
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DML_ROLLBACK_MAX_ROWS,
			Desc:  "在 DML 语句中预计影响行数超过指定值则不回滚",
			Value: "1000",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Func:          nil,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  CONFIG_DDL_OSC_MIN_SIZE,
			Desc:  "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
			Value: "16",
			Level: model.RULE_LEVEL_NORMAL,
		},
		Func:          nil,
		IsDefaultRule: true,
	},

	// rule
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST,
			Desc:  "新建表必须加入if not exists create，保证重复执行不报错",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "新建表必须加入if not exists create，保证重复执行不报错",
		Func:          checkIfNotExist,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_LENGTH,
			Desc:  "表名、列名、索引名的长度不能大于64字节",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "表名、列名、索引名的长度不能大于64字节",
		Func:          checkNewObjectName,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_NOT_EXIST,
			Desc:  "表必须有主键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "表必须有主键",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT,
			Desc:  "主键建议使用自增",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "主键建议使用自增",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED,
			Desc:  "主键建议使用 bigint 无符号类型，即 bigint unsigned",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "主键建议使用 bigint 无符号类型，即 bigint unsigned",
		Func:          checkPrimaryKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_VARCHAR_MAX,
			Desc:  "禁止使用 varchar(max)",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用 varchar(max)",
		Func:          nil,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_CHAR_LENGTH,
			Desc:  "char长度大于20时，必须使用varchar类型",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "char长度大于20时，必须使用varchar类型",
		Func:          checkStringType,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_FK,
			Desc:  "禁止使用外键",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用外键",
		Func:          checkForeignKey,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COUNT,
			Desc:  "索引个数建议不超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "5",
		},
		Message:       "索引个数建议不超过%v个",
		Func:          checkIndex,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COMPOSITE_INDEX_MAX,
			Desc:  "复合索引的列数量不建议超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "3",
		},
		Message:       "复合索引的列数量不建议超过%v个",
		Func:          checkIndex,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_USING_KEYWORD,
			Desc:  "数据库对象命名禁止使用关键字",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "数据库对象命名禁止使用关键字 %s",
		Func:          checkNewObjectName,
		IsDefaultRule: true,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_OBJECT_NAME_USING_CN,
			Desc:  "数据库对象命名不能使用英文、下划线、数字之外的字符",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "数据库对象命名不能使用英文、下划线、数字之外的字符",
		Func:          checkNewObjectName,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4,
			Desc:  "建议使用Innodb引擎,utf8mb4字符集",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "建议使用Innodb引擎,utf8mb4字符集",
		Func:          checkEngineAndCharacterSet,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_COLUMN_WITH_BLOB,
			Desc:  "禁止将blob类型的列加入索引",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止将blob类型的列加入索引",
		Func:          disableAddIndexForColumnsTypeBlob,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_IS_INVALID,
			Desc:  "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
		Func:          checkSelectWhere,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_ALTER_TABLE_NEED_MERGE,
			Desc:  "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "已存在对该表的修改语句，建议合并成一个ALTER语句",
		Func:          checkMergeAlterTable,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_DISABE_SELECT_ALL_COLUMN,
			Desc:  "不建议使用select *",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "不建议使用select *",
		Func:          checkSelectAll,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_DISABLE_DROP_STATEMENT,
			Desc:  "禁止除索引外的drop操作",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "禁止除索引外的drop操作",
		Func:          disableDropStmt,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_WITHOUT_COMMENT,
			Desc:  "表建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "表建议添加注释",
		Func:          checkTableWithoutComment,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_COMMENT,
			Desc:  "列建议添加注释",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message:       "列建议添加注释",
		Func:          checkColumnWithoutComment,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_PREFIX,
			Desc:  "普通索引必须要以\"idx_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "普通索引必须要以\"idx_\"为前缀",
		Func:          checkIndexPrefix,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_UNIQUE_INDEX_PRIFIX,
			Desc:  "unique索引必须要以\"uniq_\"为前缀",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "unique索引必须要以\"uniq_\"为前缀",
		Func:          checkUniqIndexPrefix,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_UNIQUE_INDEX,
			Desc:  "unique索引名必须使用 IDX_UK_表名_字段名",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "unique索引名必须使用 IDX_UK_表名_字段名",
		Func:    checkUniqIndex,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_WITHOUT_DEFAULT,
			Desc:  "除了自增列及大字段列之外，每个列都必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "除了自增列及大字段列之外，每个列都必须添加默认值",
		Func:          checkColumnWithoutDefault,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_TIMESTAMP_WITHOUT_DEFAULT,
			Desc:  "timestamp 类型的列必须添加默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "timestamp 类型的列必须添加默认值",
		Func:          checkColumnTimestampWithoutDefault,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_WITH_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		Func:          checkColumnBlobNotNull,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_DEFAULT_IS_NOT_NULL,
			Desc:  "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
		Func:          checkColumnBlobDefaultNull,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_LIMIT,
			Desc:  "delete/update 语句不能有limit条件",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "delete/update 语句不能有limit条件",
		Func:          checkDMLWithLimit,
		IsDefaultRule: true,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WITH_ORDER_BY,
			Desc:  "delete/update 语句不能有order by",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message:       "delete/update 语句不能有order by",
		Func:          checkDMLWithOrderBy,
		IsDefaultRule: true,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_INSERT_COLUMNS_EXIST,
			Desc:  "insert 语句必须指定column",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "insert 语句必须指定column",
		Func:    checkDMLWithInsertColumnExist,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_BATCH_INSERT_LISTS_MAX,
			Desc:  "单条insert语句，建议批量插入不超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "5000",
		},
		Message: "单条insert语句，建议批量插入不超过%v条",
		Func:    checkDMLWithBatchInsertMaxLimits,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT,
			Desc:  "主键禁止使用自增",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "主键禁止使用自增",
		Func:    checkPrimaryKey,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_FUNC,
			Desc:  "避免对条件字段使用函数操作",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "避免对条件字段使用函数操作",
		Func:    checkWhereExistFunc,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_NOT,
			Desc:  "不建议对条件字段使用负向查询",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议对条件字段使用负向查询",
		Func:    checkSelectWhere,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_NULL,
			Desc:  "不建议对条件字段使用 NULL 值判断",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议对条件字段使用 NULL 值判断",
		Func:    checkWhereExistNull,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION,
			Desc:  "条件字段存在数值和字符的隐式转换",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "条件字段存在数值和字符的隐式转换",
		Func:    checkWhereColumnImplicitConversion,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_LIMIT_MUST_EXIST,
			Desc:  "delete/update 语句必须有limit条件",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "delete/update 语句必须有limit条件",
		Func:    checkDMLLimitExist,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES,
			Desc:  "避免使用标量子查询",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "避免使用标量子查询",
		Func:    checkSelectWhere,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS,
			Desc:  "建议创建约束前,先行创建索引",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议创建约束前,先行创建索引",
		Func:    checkIndexesExistBeforeCreatConstraints,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_SELECT_FOR_UPDATE,
			Desc:  "建议避免使用select for update",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议避免使用select for update",
		Func:    checkDMLSelectForUpdate,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLLATION_DATABASE,
			Desc:  "建议使用规定的数据库排序规则",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "utf8mb4_0900_ai_ci",
		},
		Message: "建议使用规定的数据库排序规则为%s",
		Func:    checkCollationDatabase,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_DECIMAL_TYPE_COLUMN,
			Desc:  "精确浮点数建议使用DECIMAL",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "精确浮点数建议使用DECIMAL",
		Func:    checkDecimalTypeColumn,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_NEEDLESS_FUNC,
			Desc:  "避免使用不必要的内置函数",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "sha(),sqrt(),md5()",
		},
		Message: "避免使用不必要的内置函数[%v]",
		Func:    checkNeedlessFunc,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_DATABASE_SUFFIX,
			Desc:  "数据库名称建议以\"_DB\"结尾",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "数据库名称建议以\"_DB\"结尾",
		Func:    checkDatabaseSuffix,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_PK_NAME,
			Desc:  "建议主键命名为\"PK_表名\"",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议主键命名为\"PK_表名\"",
		Func:    checkPKIndexName,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TRANSACTION_ISOLATION_LEVEL,
			Desc:  "事物隔离级别建议设置成RC",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "事物隔离级别建议设置成RC",
		Func:    checkTransactionIsolationLevel,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_FUZZY_SEARCH,
			Desc:  "禁止使用全模糊搜索或左模糊搜索",
			Level: model.RULE_LEVEL_ERROR,
		},
		Message: "禁止使用全模糊搜索或左模糊搜索",
		Func:    checkSelectWhere,
	},
	RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_TABLE_PARTITION,
			Desc:  "不建议使用分区表相关功能",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议使用分区表相关功能",
		Func:    checkTablePartition,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_NUMBER_OF_JOIN_TABLES,
			Desc:  "使用JOIN连接表查询建议不超过阈值",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "3",
		},
		Message: "使用JOIN连接表查询建议不超过%v张",
		Func:    checkNumberOfJoinTables,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DML_CHECK_IS_AFTER_UNION_DISTINCT,
			Desc:  "建议使用UNION ALL,替代UNION",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "建议使用UNION ALL,替代UNION",
		Func:    checkIsAfterUnionDistinct,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_IS_EXIST_LIMIT_OFFSET,
			Desc:  "使用LIMIT分页时,避免使用LIMIT M,N",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "使用LIMIT分页时,避免使用LIMIT M,N",
		Func:    checkIsExistLimitOffset,
	}, RuleHandler{
		Rule: model.Rule{
			Name:  DDL_CHECK_INDEX_OPTION,
			Desc:  "建议选择可选性超过阈值字段作为索引",
			Level: model.RULE_LEVEL_NOTICE,
			Value: "0.7",
		},
		Message: "创建索引的字段可选性未超过阈值:%v",
		Func:    checkIndexOption,
	},
	{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_ENUM_NOTICE,
			Desc:  "不建议使用 ENUM 类型",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议使用 ENUM 类型",
		Func:    checkColumnEnumNotice,
	},
	{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_SET_NOTICE,
			Desc:  "不建议使用 SET 类型",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议使用 SET 类型",
		Func:    checkColumnSetNotice,
	},
	{
		Rule: model.Rule{
			Name:  DDL_CHECK_COLUMN_BLOB_NOTICE,
			Desc:  "不建议使用 BLOB 或 TEXT 类型",
			Level: model.RULE_LEVEL_NOTICE,
		},
		Message: "不建议使用 BLOB 或 TEXT 类型",
		Func:    checkColumnBlobNotice,
	},
	{
		Rule: model.Rule{
			Name:  DMLCheckExplainAccessTypeAll,
			Value: "10000",
			Desc:  "查询的扫描不建议超过指定行数（默认值：10000）",
			Level: model.RULE_LEVEL_WARN,
		},
		Message: "该查询的扫描行数为%v",
		Func:    checkExplain,
	},
	{
		Rule: model.Rule{
			Name:  DMLCheckExplainExtraUsingFilesort,
			Desc:  "该查询使用了文件排序",
			Level: model.RULE_LEVEL_WARN,
		},
		Message: "该查询使用了文件排序",
		Func:    checkExplain,
	},
	{
		Rule: model.Rule{
			Name:  DMLCheckExplainExtraUsingTemporary,
			Desc:  "该查询使用了临时表",
			Level: model.RULE_LEVEL_WARN,
		},
		Message: "该查询使用了临时表",
		Func:    checkExplain,
	},
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
		InitRules = append(InitRules, rh.Rule)
		if rh.IsDefaultRule {
			DefaultTemplateRules = append(DefaultTemplateRules, rh.Rule)
		}
	}
}

func checkSelectAll(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkSelectWhere(rule model.Rule, i *Inspect, node ast.Node) error {

	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		checkWhere(i, stmt.Where)

	case *ast.UpdateStmt:
		checkWhere(i, stmt.Where)
	case *ast.DeleteStmt:
		checkWhere(i, stmt.Where)
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if checkWhere(i, ss.Where) {
				break
			}
		}
	default:
		return nil
	}

	return nil
}

func checkWhere(i *Inspect, where ast.ExprNode) bool {
	isAddResult := false

	if where == nil || !whereStmtHasOneColumn(where) {
		i.addResult(DML_CHECK_WHERE_IS_INVALID)
		isAddResult = true
	}
	if where != nil && whereStmtExistNot(where) {
		i.addResult(DML_CHECK_WHERE_EXIST_NOT)
		isAddResult = true
	}
	if where != nil && whereStmtExistScalarSubQueries(where) {
		i.addResult(DML_CHECK_WHERE_EXIST_SCALAR_SUB_QUERIES)
		isAddResult = true
	}
	if where != nil && checkWhereFuzzySearch(where) {
		i.addResult(DML_CHECK_FUZZY_SEARCH)
		isAddResult = true
	}
	return isAddResult
}
func checkWhereExistNull(rule model.Rule, i *Inspect, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		var existNull bool
		scanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			if _, ok := expr.(*ast.IsNullExpr); ok {
				existNull = true
				return true
			}
			return false
		}, where)
		if existNull {
			i.addResult(rule.Name)
		}
	}
	return nil
}

func getWhereExpr(node ast.Node) (where ast.ExprNode) {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		where = stmt.Where
	case *ast.UpdateStmt:
		where = stmt.Where
	case *ast.DeleteStmt:
		where = stmt.Where

	}
	return
}

func checkIndexesExistBeforeCreatConstraints(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		constraintMap := make(map[string]struct{})
		cols := []string{}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			if spec.Constraint != nil && (spec.Constraint.Tp == ast.ConstraintPrimaryKey ||
				spec.Constraint.Tp == ast.ConstraintUniq || spec.Constraint.Tp == ast.ConstraintUniqKey) {
				for _, key := range spec.Constraint.Keys {
					cols = append(cols, key.Column.Name.String())
				}
			}
		}
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		for _, constraints := range createTableStmt.Constraints {
			for _, key := range constraints.Keys {
				constraintMap[key.Column.Name.String()] = struct{}{}
			}
		}
		for _, col := range cols {
			if _, ok := constraintMap[col]; !ok {
				i.addResult(DDL_CHECK_INDEXES_EXIST_BEFORE_CREAT_CONSTRAINTS)
				return nil
			}
		}
	}
	return nil
}

func checkPrimaryKey(rule model.Rule, i *Inspect, node ast.Node) error {
	var hasPk = false
	var pkColumnExist = false
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
				pkColumnExist = true
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
							pkColumnExist = true
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
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableAddColumns:
				for _, newColumn := range spec.NewColumns {
					if IsAllInOptions(newColumn.Options, ast.ColumnOptionPrimaryKey) {
						hasPk = true
						pkColumnExist = true
						if IsAllInOptions(newColumn.Options, ast.ColumnOptionAutoIncrement) {
							pkIsAutoIncrement = true
						}
					}
				}
			}
		}

		if originTable, exist, err := i.getCreateTableStmt(stmt.Table); err == nil && exist {
			if originPK, exist := getPrimaryKey(originTable); exist {

				hasPk = true
				pkColumnExist = true
				for _, spec := range stmt.Specs {
					switch spec.Tp {
					case ast.AlterTableModifyColumn:
						for _, newColumn := range spec.NewColumns {
							if _, exist := originPK[newColumn.Name.Name.L]; exist &&
								IsAllInOptions(newColumn.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrement = true
							}
						}
					case ast.AlterTableChangeColumn:
						if _, exist = originPK[spec.OldColumnName.Name.L]; exist {
							for _, newColumn := range spec.NewColumns {
								if IsAllInOptions(newColumn.Options, ast.ColumnOptionAutoIncrement) {
									pkIsAutoIncrement = true
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
		i.addResult(DDL_CHECK_PK_NOT_EXIST)
	}
	if hasPk && pkColumnExist && !pkIsAutoIncrement {
		i.addResult(DDL_CHECK_PK_WITHOUT_AUTO_INCREMENT)
	}
	if hasPk && pkColumnExist && pkIsAutoIncrement {
		i.addResult(DDL_CHECK_PK_PROHIBIT_AUTO_INCREMENT)
	}
	if hasPk && pkColumnExist && !pkIsBigIntUnsigned {
		i.addResult(DDL_CHECK_PK_WITHOUT_BIGINT_UNSIGNED)
	}

	return nil
}

func checkMergeAlterTable(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := i.getTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables) > 0 {
				i.addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE)
			}
		}
	}
	return nil
}

func checkEngineAndCharacterSet(rule model.Rule, i *Inspect, node ast.Node) error {
	var tableName *ast.TableName
	var engine string
	var characterSet string
	var err error
	schemaName := ""
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		tableName = stmt.Table
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
	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, ss := range stmt.Specs {
			for _, op := range ss.Options {
				switch op.Tp {
				case ast.TableOptionEngine:
					engine = op.StrValue
				case ast.TableOptionCharset:
					characterSet = op.StrValue
				}
			}
		}
	case *ast.CreateDatabaseStmt:
		schemaName = stmt.Name
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCharset {
				characterSet = ss.Value
				break
			}
		}
	case *ast.AlterDatabaseStmt:
		schemaName = stmt.Name
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCharset {
				characterSet = ss.Value
				break
			}
		}
	default:
		return nil
	}
	if engine == "" {
		engine, err = i.getSchemaEngine(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	if characterSet == "" {
		characterSet, err = i.getSchemaCharacter(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	if strings.ToLower(engine) == "innodb" && strings.ToLower(characterSet) == "utf8mb4" {
		return nil
	}
	i.addResult(DDL_CHECK_TABLE_WITHOUT_INNODB_UTF8MB4)
	return nil
}

func disableAddIndexForColumnsTypeBlob(rule model.Rule, i *Inspect, node ast.Node) error {
	isTypeBlobCols := map[string]bool{}
	indexDataTypeIsBlob := false
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
				if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
					indexDataTypeIsBlob = true
					break
				}
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				for _, col := range constraint.Keys {
					if isTypeBlobCols[col.Column.Name.String()] {
						indexDataTypeIsBlob = true
						break
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		// collect columns type from original table
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
		// collect columns type from alter table
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn) {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
						indexDataTypeIsBlob = true
						break
					}
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniq:
				for _, col := range spec.Constraint.Keys {
					if isTypeBlobCols[col.Column.Name.String()] {
						indexDataTypeIsBlob = true
						break
					}
				}
			}
		}
	case *ast.CreateIndexStmt:
		createTableStmt, exist, err := i.getCreateTableStmt(stmt.Table)
		if err != nil || !exist {
			return err
		}
		for _, col := range createTableStmt.Cols {
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
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

func checkNewObjectName(rule model.Rule, i *Inspect, node ast.Node) error {
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

	// check exist non-latin and underscore
	for _, name := range names {
		if bytes.IndexFunc([]byte(name), func(r rune) bool {
			return !(unicode.Is(unicode.Latin, r) || string(r) == "_" || unicode.IsDigit(r))
		}) != -1 {
			i.addResult(DDL_CHECK_OBJECT_NAME_USING_CN)
			break
		}

		if idx := bytes.IndexFunc([]byte(name), func(r rune) bool {
			return string(r) == "_"
		}); idx == -1 || idx == 0 || idx == len(name)-1 {
			i.addResult(DDL_CHECK_OBJECT_NAME_USING_CN)
			break
		}
	}

	// check keyword
	invalidNames := []string{}
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

func checkForeignKey(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkIndex(rule model.Rule, i *Inspect, node ast.Node) error {
	indexCounter := 0
	compositeIndexMax := 0
	value, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
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
	if indexCounter > value {
		i.addResult(DDL_CHECK_INDEX_COUNT, value)
	}
	if compositeIndexMax > value {
		i.addResult(DDL_CHECK_COMPOSITE_INDEX_MAX, value)
	}
	return nil
}

func checkStringType(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkIfNotExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDL_CHECK_TABLE_WITHOUT_IF_NOT_EXIST)
		}
	}
	return nil
}

func disableDropStmt(rule model.Rule, i *Inspect, node ast.Node) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	}
	return nil
}

func checkTableWithoutComment(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkColumnWithoutComment(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkIndexPrefix(rule model.Rule, i *Inspect, node ast.Node) error {
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
	case *ast.CreateIndexStmt:
		if !stmt.Unique {
			indexesName = append(indexesName, stmt.IndexName)
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

func checkUniqIndexPrefix(rule model.Rule, i *Inspect, node ast.Node) error {
	return checkIfUniqIndexSatisfy(rule, i, node, func(uniqIndexName, tableName string, indexedColNames []string) bool {
		return strings.HasPrefix(uniqIndexName, "uniq_")
	})
}

func checkUniqIndex(rule model.Rule, i *Inspect, node ast.Node) error {
	return checkIfUniqIndexSatisfy(rule, i, node, func(uniqIndexName, tableName string, indexedColNames []string) bool {
		return strings.EqualFold(uniqIndexName, fmt.Sprintf("IDX_UK_%v_%v", tableName, strings.Join(indexedColNames, "_")))
	})
}

func checkIfUniqIndexSatisfy(
	rule model.Rule,
	i *Inspect,
	node ast.Node,
	isSatisfy func(uniqIndexName, tableName string, indexedColNames []string) bool) error {

	var tableName string
	var indexes = make(map[string] /*unique index name*/ []string /*indexed columns*/)

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		tableName = stmt.Table.Name.String()
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniq:
				for _, key := range constraint.Keys {
					indexes[constraint.Name] = append(indexes[constraint.Name], key.Column.Name.String())
				}
			}
		}
	case *ast.AlterTableStmt:
		tableName = stmt.Table.Name.String()
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintUniq:
				for _, key := range spec.Constraint.Keys {
					indexes[spec.Constraint.Name] = append(indexes[spec.Constraint.Name], key.Column.Name.String())
				}
			}
		}
	case *ast.CreateIndexStmt:
		tableName = stmt.Table.Name.String()
		if stmt.Unique {
			for _, indexCol := range stmt.IndexColNames {
				indexes[stmt.IndexName] = append(indexes[stmt.IndexName], indexCol.Column.Name.String())
			}
		}
	default:
		return nil
	}

	for index, indexedCols := range indexes {
		if !isSatisfy(index, tableName, indexedCols) {
			i.addResult(rule.Name)
			return nil
		}
	}
	return nil
}

func checkColumnWithoutDefault(rule model.Rule, i *Inspect, node ast.Node) error {
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
			if HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
				isAutoIncrementColumn = true
			}
			if MysqlDataTypeIsBlob(col.Tp.Tp) {
				isBlobColumn = true
			}
			if HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
				columnHasDefault = true
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
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				isAutoIncrementColumn := false
				isBlobColumn := false
				columnHasDefault := false

				if HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					isAutoIncrementColumn = true
				}
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					isBlobColumn = true
				}
				if HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					columnHasDefault = true
				}

				if isAutoIncrementColumn || isBlobColumn {
					continue
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

func checkColumnTimestampWithoutDefault(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkColumnBlobNotNull(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkColumnEnumNotice(rule model.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeEnum)
}

func checkColumnSetNotice(rule model.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeSet)
}

func checkColumnBlobNotice(rule model.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob)
}

func checkColumnShouldNotBeType(rule model.Rule, i *Inspect, node ast.Node, colTypes ...byte) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col == nil {
				continue
			}
			if bytes.Contains(colTypes, []byte{col.Tp.Tp}) {
				i.addResult(rule.Name)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range getAlterTableSpecByTp(
			stmt.Specs,
			ast.AlterTableAddColumns,
			ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn) {

			for _, newCol := range spec.NewColumns {
				if newCol.Tp == nil {
					continue
				}

				if bytes.Contains(colTypes, []byte{newCol.Tp.Tp}) {
					i.addResult(rule.Name)
					return nil
				}
			}
		}
	}

	return nil
}

func checkColumnBlobDefaultNull(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkDMLWithLimit(rule model.Rule, i *Inspect, node ast.Node) error {
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
func checkDMLLimitExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit == nil {
			i.addResult(DML_CHECK_LIMIT_MUST_EXIST)
		}
	case *ast.DeleteStmt:
		if stmt.Limit == nil {
			i.addResult(DML_CHECK_LIMIT_MUST_EXIST)
		}
	}
	return nil
}

func checkDMLWithOrderBy(rule model.Rule, i *Inspect, node ast.Node) error {
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

func checkDMLWithInsertColumnExist(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Columns) == 0 {
			i.addResult(DML_CHECK_INSERT_COLUMNS_EXIST)
		}
	}
	return nil
}

func checkDMLWithBatchInsertMaxLimits(rule model.Rule, i *Inspect, node ast.Node) error {
	value, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Lists) > value {
			i.addResult(DML_CHECK_BATCH_INSERT_LISTS_MAX, value)
		}
	}
	return nil
}

func checkWhereExistFunc(rule model.Rule, i *Inspect, node ast.Node) error {
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			tableSources := getTableSources(stmt.From.TableRefs)
			// not select from table statement
			if len(tableSources) < 1 {
				break
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkExistFunc(i, tables, stmt.Where)
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			tableSources := getTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkExistFunc(i, tables, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			checkExistFunc(i, getTables(stmt.TableRefs.TableRefs), stmt.Where)
		}
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			tableSources := getTableSources(ss.From.TableRefs)
			if len(tableSources) < 1 {
				continue
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			if checkExistFunc(i, tables, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}
func checkExistFunc(i *Inspect, tables []*ast.TableName, where ast.ExprNode) bool {
	if where == nil {
		return false
	}
	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := i.getCreateTableStmt(tableName)
		if exist && err == nil {
			cols = append(cols, createTableStmt.Cols...)
		}
	}
	colMap := make(map[string]struct{})
	for _, col := range cols {
		colMap[col.Name.String()] = struct{}{}
	}
	if isFuncUsedOnColumnInWhereStmt(colMap, where) {
		i.addResult(DML_CHECK_WHERE_EXIST_FUNC)
		return true
	}
	return false
}

func checkWhereColumnImplicitConversion(rule model.Rule, i *Inspect, node ast.Node) error {
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			tableSources := getTableSources(stmt.From.TableRefs)
			// not select from table statement
			if len(tableSources) < 1 {
				break
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkWhereColumnImplicitConversionFunc(i, tables, stmt.Where)
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			tableSources := getTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkWhereColumnImplicitConversionFunc(i, tables, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			checkWhereColumnImplicitConversionFunc(i, getTables(stmt.TableRefs.TableRefs), stmt.Where)
		}
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			tableSources := getTableSources(ss.From.TableRefs)
			if len(tableSources) < 1 {
				continue
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			if checkWhereColumnImplicitConversionFunc(i, tables, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}
func checkWhereColumnImplicitConversionFunc(i *Inspect, tables []*ast.TableName, where ast.ExprNode) bool {
	if where == nil {
		return false
	}
	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := i.getCreateTableStmt(tableName)
		if exist && err == nil {
			cols = append(cols, createTableStmt.Cols...)
		}
	}
	colMap := make(map[string]string)
	for _, col := range cols {
		colType := ""
		if col.Tp == nil {
			continue
		}
		switch col.Tp.Tp {
		case mysql.TypeVarchar, mysql.TypeString:
			colType = "string"
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong, mysql.TypeDouble, mysql.TypeFloat, mysql.TypeNewDecimal:
			colType = "int"
		}
		if colType != "" {
			colMap[col.Name.String()] = colType
		}

	}
	if isColumnImplicitConversionInWhereStmt(colMap, where) {
		i.addResult(DML_CHECK_WHERE_EXIST_IMPLICIT_CONVERSION)
		return true
	}
	return false
}

func checkDMLSelectForUpdate(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.LockTp == ast.SelectLockForUpdate {
			i.addResult(DML_CHECK_SELECT_FOR_UPDATE)
		}
	}
	return nil
}

func checkCollationDatabase(rule model.Rule, i *Inspect, node ast.Node) error {
	var tableName *ast.TableName
	var collationDatabase string
	var err error
	schemaName := ""
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		tableName = stmt.Table
		if stmt.ReferTable != nil {
			return nil
		}
		for _, op := range stmt.Options {
			if op.Tp == ast.TableOptionCollate {
				collationDatabase = op.StrValue
				break
			}
		}
	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, ss := range stmt.Specs {
			for _, op := range ss.Options {
				if op.Tp == ast.TableOptionCollate {
					collationDatabase = op.StrValue
					break
				}
			}
		}
	case *ast.CreateDatabaseStmt:
		schemaName = stmt.Name
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCollate {
				collationDatabase = ss.Value
				break
			}
		}
	case *ast.AlterDatabaseStmt:
		schemaName = stmt.Name
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCollate {
				collationDatabase = ss.Value
				break
			}
		}
	default:
		return nil
	}
	if collationDatabase == "" && (tableName != nil || schemaName != "") {
		collationDatabase, err = i.getCollationDatabase(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	if !strings.EqualFold(collationDatabase, rule.Value) {
		i.addResult(DDL_CHECK_COLLATION_DATABASE, rule.Value)
	}
	return nil
}
func checkDecimalTypeColumn(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
				i.addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
					i.addResult(DDL_CHECK_DECIMAL_TYPE_COLUMN)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkNeedlessFunc(rule model.Rule, i *Inspect, node ast.Node) error {
	needlessFuncArr := strings.Split(rule.Value, ",")
	sql := strings.ToLower(node.Text())
	for _, needlessFunc := range needlessFuncArr {
		needlessFunc = strings.ToLower(strings.TrimRight(needlessFunc, ")"))
		if strings.Contains(sql, needlessFunc) {
			i.addResult(DML_CHECK_NEEDLESS_FUNC, rule.Value)
			return nil
		}
	}
	return nil
}

func checkDatabaseSuffix(rule model.Rule, i *Inspect, node ast.Node) error {
	databaseName := ""
	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		databaseName = stmt.Name
	case *ast.AlterDatabaseStmt:
		databaseName = stmt.Name
	default:
		return nil
	}
	if databaseName != "" && !strings.HasSuffix(strings.ToUpper(databaseName), "_DB") {
		i.addResult(DDL_CHECK_DATABASE_SUFFIX)
		return nil
	}
	return nil
}

func checkPKIndexName(rule model.Rule, i *Inspect, node ast.Node) error {
	indexesName := ""
	tableName := ""
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				indexesName = constraint.Name
				tableName = stmt.Table.Name.String()
				break
			}
		}
	case *ast.AlterTableStmt:
		tableName = strings.ToUpper(stmt.Table.Name.String())
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintPrimaryKey {
				indexesName = spec.Constraint.Name
				tableName = stmt.Table.Name.String()
				break
			}
		}
	default:
		return nil
	}
	if indexesName != "" && !strings.EqualFold(indexesName, "PK_"+tableName) {
		i.addResult(DDL_CHECK_PK_NAME)
		return nil
	}
	return nil
}

func checkTransactionIsolationLevel(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SetStmt:
		for _, variable := range stmt.Variables {
			if util.Contains([]string{"tx_isolation", "tx_isolation_one_shot"}, variable.Name) {
				switch node := variable.Value.(type) {
				case *driver.ValueExpr:
					if node.Datum.GetString() != ast.ReadCommitted {
						i.addResult(DDL_CHECK_TRANSACTION_ISOLATION_LEVEL)
						return nil
					}
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkTablePartition(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.PartitionNames != nil || spec.PartDefinitions != nil || spec.Partition != nil {
				i.addResult(DDL_CHECK_TABLE_PARTITION)
				return nil
			}
		}
	case *ast.CreateTableStmt:
		if stmt.Partition != nil {
			i.addResult(DDL_CHECK_TABLE_PARTITION)
			return nil
		}
	default:
		return nil
	}
	return nil
}
func checkNumberOfJoinTables(rule model.Rule, i *Inspect, node ast.Node) error {
	nums, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		if nums < getNumberOfJoinTables(stmt.From.TableRefs) {
			i.addResult(DML_CHECK_NUMBER_OF_JOIN_TABLES, rule.Value)
		}
	default:
		return nil
	}
	return nil
}

func checkIsAfterUnionDistinct(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if ss.IsAfterUnionDistinct {
				i.addResult(DML_CHECK_IS_AFTER_UNION_DISTINCT)
				return nil
			}
		}
	default:
		return nil
	}

	return nil
}

func checkIsExistLimitOffset(rule model.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.Limit.Offset != nil {
			i.addResult(DDL_CHECK_IS_EXIST_LIMIT_OFFSET)
		}
	default:
		return nil
	}
	return nil
}

func checkIndexOption(rule model.Rule, i *Inspect, node ast.Node) error {

	var tableName *ast.TableName
	indexColumns := make([]string, 0)
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			if spec.Constraint == nil {
				continue
			}
			for _, key := range spec.Constraint.Keys {
				indexColumns = append(indexColumns, key.Column.Name.String())
			}
		}
	case *ast.CreateIndexStmt:
		tableName = stmt.Table
		for _, indexCol := range stmt.IndexColNames {
			indexColumns = append(indexColumns, indexCol.Column.Name.String())
		}
	default:
		return nil
	}
	if len(indexColumns) == 0 {
		return nil
	}
	maxIndexOption, err := i.getMaxIndexOptionForTable(tableName, indexColumns)
	if err != nil {
		return err
	}
	if maxIndexOption != "" && strings.Compare(rule.Value, maxIndexOption) > 0 {
		i.addResult(rule.Name, rule.Value)
	}
	return nil
}

func checkExplain(rule model.Rule, i *Inspect, node ast.Node) error {
	switch node.(type) {
	case *ast.SelectStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
	default:
		return nil
	}

	epRecords, err := i.getExecutionPlan(node.Text())
	if err != nil {
		return err
	}
	for _, record := range epRecords {
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingFilesort) {
			i.addResult(DMLCheckExplainExtraUsingFilesort)
		}
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingTemporary) {
			i.addResult(DMLCheckExplainExtraUsingTemporary)
		}

		defaultRule := RuleHandlerMap[DMLCheckExplainAccessTypeAll].Rule
		if record.Type == executor.ExplainRecordAccessTypeAll && record.Rows > rule.GetValueInt(&defaultRule) {
			i.addResult(DMLCheckExplainAccessTypeAll, record.Rows)
		}
	}
	return nil
}
