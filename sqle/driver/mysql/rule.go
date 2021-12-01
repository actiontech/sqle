package mysql

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/ungerik/go-dry"
)

// rule type
const (
	RuleTypeGlobalConfig       = "全局配置"
	RuleTypeNamingConvention   = "命名规范"
	RuleTypeIndexingConvention = "索引规范"
	RuleTypeDDLConvention      = "DDL规范"
	RuleTypeDMLConvention      = "DML规范"
	RuleTypeUsageSuggestion    = "使用建议"
)

// inspector DDL rules
const (
	DDLCheckPKWithoutIfNotExists                = "ddl_check_table_without_if_not_exists"
	DDLCheckObjectNameLength                    = "ddl_check_object_name_length"
	DDLCheckObjectNameUsingKeyword              = "ddl_check_object_name_using_keyword"
	DDLCheckPKNotExist                          = "ddl_check_pk_not_exist"
	DDLCheckPKWithoutBigintUnsigned             = "ddl_check_pk_without_bigint_unsigned"
	DDLCheckPKWithoutAutoIncrement              = "ddl_check_pk_without_auto_increment"
	DDLCheckPKProhibitAutoIncrement             = "ddl_check_pk_prohibit_auto_increment"
	DDLCheckColumnCharLength                    = "ddl_check_column_char_length"
	DDLDisableFK                                = "ddl_disable_fk"
	DDLCheckIndexCount                          = "ddl_check_index_count"
	DDLCheckCompositeIndexMax                   = "ddl_check_composite_index_max"
	DDLCheckTableWithoutInnoDBUTF8MB4           = "ddl_check_table_without_innodb_utf8mb4"
	DDLCheckIndexedColumnWithBolb               = "ddl_check_index_column_with_blob"
	DDLCheckAlterTableNeedMerge                 = "ddl_check_alter_table_need_merge"
	DDLDisableDropStatement                     = "ddl_disable_drop_statement"
	DDLCheckTableWithoutComment                 = "ddl_check_table_without_comment"
	DDLCheckColumnWithoutComment                = "ddl_check_column_without_comment"
	DDLCheckIndexPrefix                         = "ddl_check_index_prefix"
	DDLCheckUniqueIndexPrefix                   = "ddl_check_unique_index_prefix"
	DDLCheckUniqueIndex                         = "ddl_check_unique_index"
	DDLCheckColumnWithoutDefault                = "ddl_check_column_without_default"
	DDLCheckColumnTimestampWitoutDefault        = "ddl_check_column_timestamp_without_default"
	DDLCheckColumnBlobWithNotNull               = "ddl_check_column_blob_with_not_null"
	DDLCheckColumnBlobDefaultIsNotNull          = "ddl_check_column_blob_default_is_not_null"
	DDLCheckColumnEnumNotice                    = "ddl_check_column_enum_notice"
	DDLCheckColumnSetNitice                     = "ddl_check_column_set_notice"
	DDLCheckColumnBlobNotice                    = "ddl_check_column_blob_notice"
	DDLCheckIndexesExistBeforeCreateConstraints = "ddl_check_indexes_exist_before_creat_constraints"
	DDLCheckDatabaseCollation                   = "ddl_check_collation_database"
	DDLCheckDecimalTypeColumn                   = "ddl_check_decimal_type_column"
	DDLCheckDatabaseSuffix                      = "ddl_check_database_suffix"
	DDLCheckPKName                              = "ddl_check_pk_name"
	DDLCheckTransactionIsolationLevel           = "ddl_check_transaction_isolation_level"
	DDLCheckTablePartition                      = "ddl_check_table_partition"
	DDLCheckIsExistLimitOffset                  = "ddl_check_is_exist_limit_offset"
	DDLCheckIndexOption                         = "ddl_check_index_option"
	DDLCheckOBjectNameUseCN                     = "ddl_check_object_name_using_cn"
	DDLCheckCreateView                          = "ddl_check_create_view"
	DDLCheckCreateTrigger                       = "ddl_check_create_trigger"
	DDLCheckCreateFunction                      = "ddl_check_create_function"
	DDLCheckCreateProcedure                     = "ddl_check_create_procedure"
)

// inspector DML rules
const (
	DMLCheckWithLimit                    = "dml_check_with_limit"
	DMLCheckWithOrderBy                  = "dml_check_with_order_by"
	DMLCheckWhereIsInvalid               = "all_check_where_is_invalid"
	DMLDisableSelectAllColumn            = "dml_disable_select_all_column"
	DMLCheckInsertColumnsExist           = "dml_check_insert_columns_exist"
	DMLCheckBatchInsertListsMax          = "dml_check_batch_insert_lists_max"
	DMLCheckWhereExistFunc               = "dml_check_where_exist_func"
	DMLCheckWhereExistNot                = "dml_check_where_exist_not"
	DMLCheckWhereExistImplicitConversion = "dml_check_where_exist_implicit_conversion"
	DMLCheckLimitMustExist               = "dml_check_limit_must_exist"
	DMLCheckWhereExistScalarSubquery     = "dml_check_where_exist_scalar_sub_queries"
	DMLWhereExistNull                    = "dml_check_where_exist_null"
	DMLCheckSelectForUpdate              = "dml_check_select_for_update"
	DMLCheckNeedlessFunc                 = "dml_check_needless_func"
	DMLCheckFuzzySearch                  = "dml_check_fuzzy_search"
	DMLCheckNumberOfJoinTables           = "dml_check_number_of_join_tables"
	DMLCheckIfAfterUnionDistinct         = "dml_check_is_after_union_distinct"
	DMLCheckExplainAccessTypeAll         = "dml_check_explain_access_type_all"
	DMLCheckExplainExtraUsingFilesort    = "dml_check_explain_extra_using_filesort"
	DMLCheckExplainExtraUsingTemporary   = "dml_check_explain_extra_using_temporary"
)

// inspector config code
const (
	ConfigDMLRollbackMaxRows = "dml_rollback_max_rows"
	ConfigDDLOSCMinSize      = "ddl_osc_min_size"
	ConfigDDLGhostMinSize    = "ddl_ghost_min_size"
)

type RuleHandler struct {
	Rule                 driver.Rule
	Message              string
	Func                 func(driver.Rule, *Inspect, ast.Node) error
	AllowOffline         bool
	NotAllowOfflineStmts []ast.Node
}

func (rh *RuleHandler) IsAllowOfflineRule(node ast.Node) bool {
	if !rh.AllowOffline {
		return false
	}
	for _, stmt := range rh.NotAllowOfflineStmts {
		if reflect.TypeOf(stmt) == reflect.TypeOf(node) {
			return false
		}
	}
	return true
}

var (
	RuleHandlerMap = map[string]RuleHandler{}

	// DefaultTemplateRules only use for unit test now. It should be removed later,
	// because Driver layer should not care about Rule template. TODO(@wy)
	DefaultTemplateRules = []driver.Rule{}
	InitRules            = []driver.Rule{}
)

var RuleHandlers = []RuleHandler{
	// config
	{
		Rule: driver.Rule{
			Name:     ConfigDMLRollbackMaxRows,
			Desc:     "在 DML 语句中预计影响行数超过指定值则不回滚",
			Value:    "1000",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeGlobalConfig,
		},
		Func: nil,
	},
	{
		Rule: driver.Rule{
			Name:     ConfigDDLOSCMinSize,
			Desc:     "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
			Value:    "16",
			Level:    driver.RuleLevelNormal,
			Category: RuleTypeGlobalConfig,
		},
		Func: nil,
	},

	{
		Rule: driver.Rule{
			Name:     ConfigDDLGhostMinSize,
			Desc:     "改表时，表空间超过指定大小(MB)时使用gh-ost上线",
			Value:    "16",
			Level:    driver.RuleLevelNormal,
			Category: RuleTypeGlobalConfig,
		},
		Func: nil,
	},

	// rule
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKWithoutIfNotExists,
			Desc:     "新建表必须加入if not exists create，保证重复执行不报错",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "新建表必须加入if not exists create，保证重复执行不报错",
		AllowOffline: true,
		Func:         checkIfNotExist,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckObjectNameLength,
			Desc:     "表名、列名、索引名的长度不能大于64字节",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "表名、列名、索引名的长度不能大于64字节",
		AllowOffline: true,
		Func:         checkNewObjectName,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKNotExist,
			Desc:     "表必须有主键",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexingConvention,
		},
		Message:              "表必须有主键",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                 checkPrimaryKey,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKWithoutAutoIncrement,
			Desc:     "主键建议使用自增",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexingConvention,
		},
		Message:              "主键建议使用自增",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                 checkPrimaryKey,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKWithoutBigintUnsigned,
			Desc:     "主键建议使用 bigint 无符号类型，即 bigint unsigned",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexingConvention,
		},
		Message:              "主键建议使用 bigint 无符号类型，即 bigint unsigned",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                 checkPrimaryKey,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnCharLength,
			Desc:     "char长度大于20时，必须使用varchar类型",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "char长度大于20时，必须使用varchar类型",
		AllowOffline: true,
		Func:         checkStringType,
	},
	{
		Rule: driver.Rule{
			Name:     DDLDisableFK,
			Desc:     "禁止使用外键",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexingConvention,
		},
		Message:      "禁止使用外键",
		AllowOffline: true,
		Func:         checkForeignKey,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexCount,
			Desc:     "索引个数建议不超过阈值",
			Level:    driver.RuleLevelNotice,
			Value:    "5",
			Category: RuleTypeIndexingConvention,
		},
		Message:              "索引个数建议不超过%v个",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                 checkIndex,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCompositeIndexMax,
			Desc:     "复合索引的列数量不建议超过阈值",
			Level:    driver.RuleLevelNotice,
			Value:    "3",
			Category: RuleTypeIndexingConvention,
		},
		Message:              "复合索引的列数量不建议超过%v个",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                 checkIndex,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckObjectNameUsingKeyword,
			Desc:     "数据库对象命名禁止使用保留字",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "数据库对象命名禁止使用保留字 %s",
		AllowOffline: true,
		Func:         checkNewObjectName,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckOBjectNameUseCN,
			Desc:     "数据库对象命名只能使用英文、下划线或数字，首字母必须是英文",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "数据库对象命名只能使用英文、下划线或数字，首字母必须是英文",
		AllowOffline: true,
		Func:         checkNewObjectName,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTableWithoutInnoDBUTF8MB4,
			Desc:     "建议使用Innodb引擎,utf8mb4字符集",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "建议使用Innodb引擎,utf8mb4字符集",
		AllowOffline: false,
		Func:         checkEngineAndCharacterSet,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexedColumnWithBolb,
			Desc:     "禁止将blob类型的列加入索引",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexingConvention,
		},
		Message:              "禁止将blob类型的列加入索引",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                 disableAddIndexForColumnsTypeBlob,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWhereIsInvalid,
			Desc:     "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message:      "禁止使用没有where条件的sql语句或者使用where 1=1等变相没有条件的sql",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckAlterTableNeedMerge,
			Desc:     "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "已存在对该表的修改语句，建议合并成一个ALTER语句",
		AllowOffline: false,
		Func:         checkMergeAlterTable,
	},
	{
		Rule: driver.Rule{
			Name:     DMLDisableSelectAllColumn,
			Desc:     "不建议使用select *",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "不建议使用select *",
		AllowOffline: true,
		Func:         checkSelectAll,
	},
	{
		Rule: driver.Rule{
			Name:     DDLDisableDropStatement,
			Desc:     "禁止除索引外的drop操作",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "禁止除索引外的drop操作",
		AllowOffline: true,
		Func:         disableDropStmt,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTableWithoutComment,
			Desc:     "表建议添加注释",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "表建议添加注释",
		AllowOffline: true,
		Func:         checkTableWithoutComment,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnWithoutComment,
			Desc:     "列建议添加注释",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "列建议添加注释",
		AllowOffline: true,
		Func:         checkColumnWithoutComment,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexPrefix,
			Desc:     "普通索引必须要以\"idx_\"为前缀",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "普通索引必须要以\"idx_\"为前缀",
		AllowOffline: true,
		Func:         checkIndexPrefix,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckUniqueIndexPrefix,
			Desc:     "unique索引必须要以\"uniq_\"为前缀",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "unique索引必须要以\"uniq_\"为前缀",
		AllowOffline: true,
		Func:         checkUniqIndexPrefix,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckUniqueIndex,
			Desc:     "unique索引名必须使用 IDX_UK_表名_字段名",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
		},
		Message:      "unique索引名必须使用 IDX_UK_表名_字段名",
		AllowOffline: true,
		Func:         checkUniqIndex,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnWithoutDefault,
			Desc:     "除了自增列及大字段列之外，每个列都必须添加默认值",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "除了自增列及大字段列之外，每个列都必须添加默认值",
		AllowOffline: true,
		Func:         checkColumnWithoutDefault,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnTimestampWitoutDefault,
			Desc:     "timestamp 类型的列必须添加默认值",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "timestamp 类型的列必须添加默认值",
		AllowOffline: true,
		Func:         checkColumnTimestampWithoutDefault,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnBlobWithNotNull,
			Desc:     "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		AllowOffline: true,
		Func:         checkColumnBlobNotNull,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnBlobDefaultIsNotNull,
			Desc:     "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "BLOB 和 TEXT 类型的字段不可指定非 NULL 的默认值",
		AllowOffline: true,
		Func:         checkColumnBlobDefaultNull,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWithLimit,
			Desc:     "delete/update 语句不能有limit条件",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message:      "delete/update 语句不能有limit条件",
		AllowOffline: true,
		Func:         checkDMLWithLimit,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWithOrderBy,
			Desc:     "delete/update 语句不能有order by",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message:      "delete/update 语句不能有order by",
		AllowOffline: true,
		Func:         checkDMLWithOrderBy,
	},
	{
		// TODO: 修改level以适配默认模板
		Rule: driver.Rule{
			Name:     DMLCheckInsertColumnsExist,
			Desc:     "insert 语句必须指定column",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "insert 语句必须指定column",
		AllowOffline: true,
		Func:         checkDMLWithInsertColumnExist,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckBatchInsertListsMax,
			Desc:     "单条insert语句，建议批量插入不超过阈值",
			Level:    driver.RuleLevelNotice,
			Value:    "5000",
			Category: RuleTypeDMLConvention,
		},
		Message:      "单条insert语句，建议批量插入不超过%v条",
		AllowOffline: true,
		Func:         checkDMLWithBatchInsertMaxLimits,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKProhibitAutoIncrement,
			Desc:     "主键禁止使用自增",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeIndexingConvention,
		},
		Message:              "主键禁止使用自增",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                 checkPrimaryKey,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWhereExistFunc,
			Desc:     "避免对条件字段使用函数操作",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "避免对条件字段使用函数操作",
		AllowOffline: false,
		Func:         checkWhereExistFunc,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWhereExistNot,
			Desc:     "不建议对条件字段使用负向查询",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "不建议对条件字段使用负向查询",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driver.Rule{
			Name:     DMLWhereExistNull,
			Desc:     "不建议对条件字段使用 NULL 值判断",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "不建议对条件字段使用 NULL 值判断",
		Func:         checkWhereExistNull,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWhereExistImplicitConversion,
			Desc:     "条件字段存在数值和字符的隐式转换",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message: "条件字段存在数值和字符的隐式转换",
		Func:    checkWhereColumnImplicitConversion,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckLimitMustExist,
			Desc:     "delete/update 语句必须有limit条件",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message:      "delete/update 语句必须有limit条件",
		Func:         checkDMLLimitExist,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckWhereExistScalarSubquery,
			Desc:     "避免使用标量子查询",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "避免使用标量子查询",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexesExistBeforeCreateConstraints,
			Desc:     "建议创建约束前,先行创建索引",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeIndexingConvention,
		},
		Message: "建议创建约束前,先行创建索引",
		Func:    checkIndexesExistBeforeCreatConstraints,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckSelectForUpdate,
			Desc:     "建议避免使用select for update",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "建议避免使用select for update",
		Func:         checkDMLSelectForUpdate,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckDatabaseCollation,
			Desc:     "建议使用规定的数据库排序规则",
			Level:    driver.RuleLevelNotice,
			Value:    "utf8mb4_0900_ai_ci",
			Category: RuleTypeDDLConvention,
		},
		Message: "建议使用规定的数据库排序规则为%s",
		Func:    checkCollationDatabase,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckDecimalTypeColumn,
			Desc:     "精确浮点数建议使用DECIMAL",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "精确浮点数建议使用DECIMAL",
		Func:         checkDecimalTypeColumn,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckNeedlessFunc,
			Desc:     "避免使用不必要的内置函数",
			Level:    driver.RuleLevelNotice,
			Value:    "sha(),sqrt(),md5()",
			Category: RuleTypeDMLConvention,
		},
		Message:      "避免使用不必要的内置函数[%v]",
		Func:         checkNeedlessFunc,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckDatabaseSuffix,
			Desc:     "数据库名称建议以\"_DB\"结尾",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeNamingConvention,
		},
		Message:      "数据库名称建议以\"_DB\"结尾",
		Func:         checkDatabaseSuffix,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckPKName,
			Desc:     "建议主键命名为\"PK_表名\"",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeNamingConvention,
		},
		Message:      "建议主键命名为\"PK_表名\"",
		Func:         checkPKIndexName,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTransactionIsolationLevel,
			Desc:     "事物隔离级别建议设置成RC",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "事物隔离级别建议设置成RC",
		Func:         checkTransactionIsolationLevel,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckFuzzySearch,
			Desc:     "禁止使用全模糊搜索或左模糊搜索",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDMLConvention,
		},
		Message:      "禁止使用全模糊搜索或左模糊搜索",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTablePartition,
			Desc:     "不建议使用分区表相关功能",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "不建议使用分区表相关功能",
		AllowOffline: true,
		Func:         checkTablePartition,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckNumberOfJoinTables,
			Desc:     "使用JOIN连接表查询建议不超过阈值",
			Level:    driver.RuleLevelNotice,
			Value:    "3",
			Category: RuleTypeDMLConvention,
		},
		Message:      "使用JOIN连接表查询建议不超过%v张",
		AllowOffline: true,
		Func:         checkNumberOfJoinTables,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckIfAfterUnionDistinct,
			Desc:     "建议使用UNION ALL,替代UNION",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "建议使用UNION ALL,替代UNION",
		AllowOffline: true,
		Func:         checkIsAfterUnionDistinct,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIsExistLimitOffset,
			Desc:     "使用LIMIT分页时,避免使用LIMIT M,N",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDMLConvention,
		},
		Message:      "使用LIMIT分页时,避免使用LIMIT M,N",
		AllowOffline: true,
		Func:         checkIsExistLimitOffset,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexOption,
			Desc:     "建议选择可选性超过阈值字段作为索引",
			Level:    driver.RuleLevelNotice,
			Value:    "0.7",
			Category: RuleTypeDMLConvention,
		},
		Message:      "创建索引的字段可选性未超过阈值:%v",
		AllowOffline: false,
		Func:         checkIndexOption,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnEnumNotice,
			Desc:     "不建议使用 ENUM 类型",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "不建议使用 ENUM 类型",
		AllowOffline: true,
		Func:         checkColumnEnumNotice,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnSetNitice,
			Desc:     "不建议使用 SET 类型",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "不建议使用 SET 类型",
		AllowOffline: true,
		Func:         checkColumnSetNotice,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckColumnBlobNotice,
			Desc:     "不建议使用 BLOB 或 TEXT 类型",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "不建议使用 BLOB 或 TEXT 类型",
		AllowOffline: true,
		Func:         checkColumnBlobNotice,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckExplainAccessTypeAll,
			Value:    "10000",
			Desc:     "查询的扫描不建议超过指定行数（默认值：10000）",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message:      "该查询的扫描行数为%v",
		AllowOffline: false,
		Func:         checkExplain,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckExplainExtraUsingFilesort,
			Desc:     "该查询使用了文件排序",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message:      "该查询使用了文件排序",
		AllowOffline: false,
		Func:         checkExplain,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckExplainExtraUsingTemporary,
			Desc:     "该查询使用了临时表",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
		},
		Message:      "该查询使用了临时表",
		AllowOffline: false,
		Func:         checkExplain,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCreateView,
			Desc:     "禁止使用视图",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "禁止使用视图",
		AllowOffline: true,
		Func:         checkCreateView,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCreateTrigger,
			Desc:     "禁止使用触发器",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "禁止使用触发器",
		AllowOffline: true,
		Func:         checkCreateTrigger,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCreateFunction,
			Desc:     "禁止使用自定义函数",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "禁止使用自定义函数",
		AllowOffline: true,
		Func:         checkCreateFunction,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCreateProcedure,
			Desc:     "禁止使用存储过程",
			Level:    driver.RuleLevelError,
			Category: RuleTypeUsageSuggestion,
		},
		Message:      "禁止使用存储过程",
		AllowOffline: true,
		Func:         checkCreateProcedure,
	},
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
		InitRules = append(InitRules, rh.Rule)
		if rh.Rule.Level == driver.RuleLevelError {
			DefaultTemplateRules = append(DefaultTemplateRules, rh.Rule)
		}
	}
}

func checkSelectAll(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					i.addResult(DMLDisableSelectAllColumn)
				}
			}
		}
	}
	return nil
}

func checkSelectWhere(rule driver.Rule, i *Inspect, node ast.Node) error {

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
		i.addResult(DMLCheckWhereIsInvalid)
		isAddResult = true
	}
	if where != nil && whereStmtExistNot(where) {
		i.addResult(DMLCheckWhereExistNot)
		isAddResult = true
	}
	if where != nil && whereStmtExistScalarSubQueries(where) {
		i.addResult(DMLCheckWhereExistScalarSubquery)
		isAddResult = true
	}
	if where != nil && checkWhereFuzzySearch(where) {
		i.addResult(DMLCheckFuzzySearch)
		isAddResult = true
	}
	return isAddResult
}
func checkWhereExistNull(rule driver.Rule, i *Inspect, node ast.Node) error {
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

func checkIndexesExistBeforeCreatConstraints(rule driver.Rule, i *Inspect, node ast.Node) error {
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
				i.addResult(DDLCheckIndexesExistBeforeCreateConstraints)
				return nil
			}
		}
	}
	return nil
}

func checkPrimaryKey(rule driver.Rule, i *Inspect, node ast.Node) error {
	var pkIsAutoIncrement = false
	var pkIsBigIntUnsigned = false
	inspectCol := func(col *ast.ColumnDef) {
		if IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
			pkIsAutoIncrement = true
		}
		if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) {
			pkIsBigIntUnsigned = true
		}
	}

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		var hasPk = false
		var pkColumnExist = false

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
				inspectCol(col)
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
							inspectCol(col)
						}
					}
				}
			}
		}
		if !hasPk {
			i.addResult(DDLCheckPKNotExist)
		}
		if hasPk && pkColumnExist && !pkIsAutoIncrement {
			i.addResult(DDLCheckPKWithoutAutoIncrement)
		}
		if hasPk && pkColumnExist && pkIsAutoIncrement {
			i.addResult(DDLCheckPKProhibitAutoIncrement)
		}
		if hasPk && pkColumnExist && !pkIsBigIntUnsigned {
			i.addResult(DDLCheckPKWithoutBigintUnsigned)
		}
	case *ast.AlterTableStmt:
		var alterPK bool
		if originTable, exist, err := i.getCreateTableStmt(stmt.Table); err == nil && exist {
			for _, spec := range stmt.Specs {
				switch spec.Tp {
				case ast.AlterTableAddColumns:
					for _, newColumn := range spec.NewColumns {
						if IsAllInOptions(newColumn.Options, ast.ColumnOptionPrimaryKey) {
							alterPK = true
							inspectCol(newColumn)
						}
					}
				case ast.AlterTableAddConstraint:
					if spec.Constraint.Tp == ast.ConstraintPrimaryKey {
						if len(spec.Constraint.Keys) == 1 {
							for _, col := range originTable.Cols {
								if col.Name.Name.L == spec.Constraint.Keys[0].Column.Name.L {
									alterPK = true
									inspectCol(col)
								}
							}
						}
					}
				}
			}

			if originPK, exist := getPrimaryKey(originTable); exist {
				for _, spec := range stmt.Specs {
					switch spec.Tp {
					case ast.AlterTableModifyColumn:
						for _, newColumn := range spec.NewColumns {
							if _, exist := originPK[newColumn.Name.Name.L]; exist {
								alterPK = true
								inspectCol(newColumn)
							}
						}
					case ast.AlterTableChangeColumn:
						if _, exist = originPK[spec.OldColumnName.Name.L]; exist {
							for _, newColumn := range spec.NewColumns {
								alterPK = true
								inspectCol(newColumn)
							}
						}
					}
				}
			}
		}
		if alterPK && !pkIsAutoIncrement {
			i.addResult(DDLCheckPKWithoutAutoIncrement)
		}
		if alterPK && pkIsAutoIncrement {
			i.addResult(DDLCheckPKProhibitAutoIncrement)
		}
		if alterPK && !pkIsBigIntUnsigned {
			i.addResult(DDLCheckPKWithoutBigintUnsigned)
		}
	default:
		return nil
	}
	return nil
}

func checkMergeAlterTable(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := i.getTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables) > 0 {
				i.addResult(DDLCheckAlterTableNeedMerge)
			}
		}
	}
	return nil
}

func checkEngineAndCharacterSet(rule driver.Rule, i *Inspect, node ast.Node) error {
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
	i.addResult(DDLCheckTableWithoutInnoDBUTF8MB4)
	return nil
}

func disableAddIndexForColumnsTypeBlob(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DDLCheckIndexedColumnWithBolb)
	}
	return nil
}

func checkNewObjectName(rule driver.Rule, i *Inspect, node ast.Node) error {
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
			i.addResult(DDLCheckObjectNameLength)
			break
		}
	}

	// check exist non-latin and underscore
	for _, name := range names {
		// CASE:
		// 	CREATE TABLE t1(id int, INDEX (id)); // when index name is anonymous, skip inspect it
		if name == "" {
			continue
		}
		if !unicode.Is(unicode.Latin, rune(name[0])) ||
			bytes.IndexFunc([]byte(name), func(r rune) bool {
				return !(unicode.Is(unicode.Latin, r) || string(r) == "_" || unicode.IsDigit(r))
			}) != -1 {

			i.addResult(DDLCheckOBjectNameUseCN)
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
		i.addResult(DDLCheckObjectNameUsingKeyword,
			strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func checkForeignKey(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DDLDisableFK)
	}
	return nil
}

func checkIndex(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DDLCheckIndexCount, value)
	}
	if compositeIndexMax > value {
		i.addResult(DDLCheckCompositeIndexMax, value)
	}
	return nil
}

func checkStringType(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(DDLCheckColumnCharLength)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(DDLCheckColumnCharLength)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkIfNotExist(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDLCheckPKWithoutIfNotExists)
		}
	}
	return nil
}

func disableDropStmt(rule driver.Rule, i *Inspect, node ast.Node) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDLDisableDropStatement)
	case *ast.DropTableStmt:
		i.addResult(DDLDisableDropStatement)
	}
	return nil
}

func checkTableWithoutComment(rule driver.Rule, i *Inspect, node ast.Node) error {
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
			i.addResult(DDLCheckTableWithoutComment)
		}
	}
	return nil
}

func checkColumnWithoutComment(rule driver.Rule, i *Inspect, node ast.Node) error {
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
				i.addResult(DDLCheckColumnWithoutComment)
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
					i.addResult(DDLCheckColumnWithoutComment)
					return nil
				}
			}
		}
	}
	return nil
}

func checkIndexPrefix(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		if !utils.HasPrefix(name, "idx_", false) {
			i.addResult(DDLCheckIndexPrefix)
			return nil
		}
	}
	return nil
}

func checkUniqIndexPrefix(rule driver.Rule, i *Inspect, node ast.Node) error {
	return checkIfUniqIndexSatisfy(rule, i, node, func(uniqIndexName, tableName string, indexedColNames []string) bool {
		return utils.HasPrefix(uniqIndexName, "uniq_", false)
	})
}

func checkUniqIndex(rule driver.Rule, i *Inspect, node ast.Node) error {
	return checkIfUniqIndexSatisfy(rule, i, node, func(uniqIndexName, tableName string, indexedColNames []string) bool {
		return strings.EqualFold(uniqIndexName, fmt.Sprintf("IDX_UK_%v_%v", tableName, strings.Join(indexedColNames, "_")))
	})
}

func checkIfUniqIndexSatisfy(
	rule driver.Rule,
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

func checkColumnWithoutDefault(rule driver.Rule, i *Inspect, node ast.Node) error {
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
				i.addResult(DDLCheckColumnWithoutDefault)
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
					i.addResult(DDLCheckColumnWithoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnTimestampWithoutDefault(rule driver.Rule, i *Inspect, node ast.Node) error {
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
				i.addResult(DDLCheckColumnTimestampWitoutDefault)
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
					i.addResult(DDLCheckColumnTimestampWitoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnBlobNotNull(rule driver.Rule, i *Inspect, node ast.Node) error {
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
						i.addResult(DDLCheckColumnBlobWithNotNull)
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
							i.addResult(DDLCheckColumnBlobWithNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkColumnEnumNotice(rule driver.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeEnum)
}

func checkColumnSetNotice(rule driver.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeSet)
}

func checkColumnBlobNotice(rule driver.Rule, i *Inspect, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, i, node, mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob)
}

func checkColumnShouldNotBeType(rule driver.Rule, i *Inspect, node ast.Node, colTypes ...byte) error {
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

func checkColumnBlobDefaultNull(rule driver.Rule, i *Inspect, node ast.Node) error {
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
						i.addResult(DDLCheckColumnBlobDefaultIsNotNull)
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
							i.addResult(DDLCheckColumnBlobDefaultIsNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkDMLWithLimit(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			i.addResult(DMLCheckWithLimit)
		}
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			i.addResult(DMLCheckWithLimit)
		}
	}
	return nil
}
func checkDMLLimitExist(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit == nil {
			i.addResult(DMLCheckLimitMustExist)
		}
	case *ast.DeleteStmt:
		if stmt.Limit == nil {
			i.addResult(DMLCheckLimitMustExist)
		}
	}
	return nil
}

func checkDMLWithOrderBy(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			i.addResult(DMLCheckWithOrderBy)
		}
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			i.addResult(DMLCheckWithOrderBy)
		}
	}
	return nil
}

func checkDMLWithInsertColumnExist(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Columns) == 0 {
			i.addResult(DMLCheckInsertColumnsExist)
		}
	}
	return nil
}

func checkDMLWithBatchInsertMaxLimits(rule driver.Rule, i *Inspect, node ast.Node) error {
	value, err := strconv.Atoi(rule.Value)
	if err != nil {
		return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	}
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Lists) > value {
			i.addResult(DMLCheckBatchInsertListsMax, value)
		}
	}
	return nil
}

func checkWhereExistFunc(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DMLCheckWhereExistFunc)
		return true
	}
	return false
}

func checkWhereColumnImplicitConversion(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DMLCheckWhereExistImplicitConversion)
		return true
	}
	return false
}

func checkDMLSelectForUpdate(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.LockTp == ast.SelectLockForUpdate {
			i.addResult(DMLCheckSelectForUpdate)
		}
	}
	return nil
}

func checkCollationDatabase(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DDLCheckDatabaseCollation, rule.Value)
	}
	return nil
}
func checkDecimalTypeColumn(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
				i.addResult(DDLCheckDecimalTypeColumn)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
					i.addResult(DDLCheckDecimalTypeColumn)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkNeedlessFunc(rule driver.Rule, i *Inspect, node ast.Node) error {
	needlessFuncArr := strings.Split(rule.Value, ",")
	sql := strings.ToLower(node.Text())
	for _, needlessFunc := range needlessFuncArr {
		needlessFunc = strings.ToLower(strings.TrimRight(needlessFunc, ")"))
		if strings.Contains(sql, needlessFunc) {
			i.addResult(DMLCheckNeedlessFunc, rule.Value)
			return nil
		}
	}
	return nil
}

func checkDatabaseSuffix(rule driver.Rule, i *Inspect, node ast.Node) error {
	databaseName := ""
	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		databaseName = stmt.Name
	case *ast.AlterDatabaseStmt:
		databaseName = stmt.Name
	default:
		return nil
	}
	if databaseName != "" && !utils.HasSuffix(databaseName, "_DB", false) {
		i.addResult(DDLCheckDatabaseSuffix)
		return nil
	}
	return nil
}

func checkPKIndexName(rule driver.Rule, i *Inspect, node ast.Node) error {
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
		i.addResult(DDLCheckPKName)
		return nil
	}
	return nil
}

func checkTransactionIsolationLevel(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SetStmt:
		for _, variable := range stmt.Variables {
			if dry.StringListContains([]string{"tx_isolation", "tx_isolation_one_shot"}, variable.Name) {
				switch node := variable.Value.(type) {
				case *parserdriver.ValueExpr:
					if node.Datum.GetString() != ast.ReadCommitted {
						i.addResult(DDLCheckTransactionIsolationLevel)
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

func checkTablePartition(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.PartitionNames != nil || spec.PartDefinitions != nil || spec.Partition != nil {
				i.addResult(DDLCheckTablePartition)
				return nil
			}
		}
	case *ast.CreateTableStmt:
		if stmt.Partition != nil {
			i.addResult(DDLCheckTablePartition)
			return nil
		}
	default:
		return nil
	}
	return nil
}
func checkNumberOfJoinTables(rule driver.Rule, i *Inspect, node ast.Node) error {
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
			i.addResult(DMLCheckNumberOfJoinTables, rule.Value)
		}
	default:
		return nil
	}
	return nil
}

func checkIsAfterUnionDistinct(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if ss.IsAfterUnionDistinct {
				i.addResult(DMLCheckIfAfterUnionDistinct)
				return nil
			}
		}
	default:
		return nil
	}

	return nil
}

func checkIsExistLimitOffset(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.Limit.Offset != nil {
			i.addResult(DDLCheckIsExistLimitOffset)
		}
	default:
		return nil
	}
	return nil
}

func checkIndexOption(rule driver.Rule, i *Inspect, node ast.Node) error {

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

func checkExplain(rule driver.Rule, i *Inspect, node ast.Node) error {
	// sql from MyBatis XML file is not the executable sql. so can't do explain for it.
	// TODO(@wy) ignore explain when audit Mybatis file
	//if i.Task.SQLSource == driver.TaskSQLSourceFromMyBatisXMLFile {
	//	return nil
	//}
	switch node.(type) {
	case *ast.SelectStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
	default:
		return nil
	}

	epRecords, err := i.getExecutionPlan(node.Text())
	if err != nil {
		// TODO: check dml related table or database is created, if not exist, explain will executed failure.
		i.Logger().Errorf("do explain error: %v, sql: %s", err, node.Text())
		return nil
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

func checkCreateView(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch node.(type) {
	case *ast.CreateViewStmt:
		i.addResult(rule.Name)
	}
	return nil
}

var createTriggerReg1 = regexp.MustCompile(`(?i)create[\s]+trigger[\s]+[\S\s]+before|after`)
var createTriggerReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+trigger[\s]+[\S\s]+before|after`)

// CREATE
//    [DEFINER = user]
//    TRIGGER trigger_name
//    trigger_time trigger_event
//    ON tbl_name FOR EACH ROW
//    [trigger_order]
//    trigger_body
//
// ref:https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
//
// For now, we do character matching for CREATE TRIGGER Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateTrigger(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createTriggerReg1.MatchString(node.Text()) ||
			createTriggerReg2.MatchString(node.Text()) {
			i.addResult(rule.Name)
		}
	}
	return nil
}

var createFunctionReg1 = regexp.MustCompile(`(?i)create[\s]+function[\s]+[\S\s]+returns`)
var createFunctionReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+function[\s]+[\S\s]+returns`)

// CREATE
//    [DEFINER = user]
//    FUNCTION sp_name ([func_parameter[,...]])
//    RETURNS type
//    [characteristic ...] routine_body
//
// ref: https://dev.mysql.com/doc/refman/5.7/en/create-procedure.html
// For now, we do character matching for CREATE FUNCTION Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateFunction(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createFunctionReg1.MatchString(node.Text()) ||
			createFunctionReg2.MatchString(node.Text()) {
			i.addResult(rule.Name)
		}
	}
	return nil
}

var createProcedureReg1 = regexp.MustCompile(`(?i)create[\s]+procedure[\s]+[\S\s]+`)
var createProcedureReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+procedure[\s]+[\S\s]+`)

// CREATE
//    [DEFINER = user]
//    PROCEDURE sp_name ([proc_parameter[,...]])
//    [characteristic ...] routine_body
//
// ref: https://dev.mysql.com/doc/refman/8.0/en/create-procedure.html
// For now, we do character matching for CREATE PROCEDURE Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateProcedure(rule driver.Rule, i *Inspect, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createProcedureReg1.MatchString(node.Text()) ||
			createProcedureReg2.MatchString(node.Text()) {
			i.addResult(rule.Name)
		}
	}
	return nil
}
