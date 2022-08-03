package rule

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/keyword"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"
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
	RuleTypeIndexOptimization  = "索引优化"
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
	DDLCheckTableDBEngine                       = "ddl_check_table_db_engine"
	DDLCheckTableCharacterSet                   = "ddl_check_table_character_set"
	DDLCheckIndexedColumnWithBlob               = "ddl_check_index_column_with_blob"
	DDLCheckAlterTableNeedMerge                 = "ddl_check_alter_table_need_merge"
	DDLDisableDropStatement                     = "ddl_disable_drop_statement"
	DDLCheckTableWithoutComment                 = "ddl_check_table_without_comment"
	DDLCheckColumnWithoutComment                = "ddl_check_column_without_comment"
	DDLCheckIndexPrefix                         = "ddl_check_index_prefix"
	DDLCheckUniqueIndexPrefix                   = "ddl_check_unique_index_prefix"
	DDLCheckUniqueIndex                         = "ddl_check_unique_index"
	DDLCheckColumnWithoutDefault                = "ddl_check_column_without_default"
	DDLCheckColumnTimestampWithoutDefault       = "ddl_check_column_timestamp_without_default"
	DDLCheckColumnBlobWithNotNull               = "ddl_check_column_blob_with_not_null"
	DDLCheckColumnBlobDefaultIsNotNull          = "ddl_check_column_blob_default_is_not_null"
	DDLCheckColumnEnumNotice                    = "ddl_check_column_enum_notice"
	DDLCheckColumnSetNotice                     = "ddl_check_column_set_notice"
	DDLCheckColumnBlobNotice                    = "ddl_check_column_blob_notice"
	DDLCheckIndexesExistBeforeCreateConstraints = "ddl_check_indexes_exist_before_creat_constraints"
	DDLCheckDatabaseCollation                   = "ddl_check_collation_database"
	DDLCheckDecimalTypeColumn                   = "ddl_check_decimal_type_column"
	DDLCheckBigintInsteadOfDecimal              = "ddl_check_bigint_instead_of_decimal"
	DDLCheckDatabaseSuffix                      = "ddl_check_database_suffix"
	DDLCheckPKName                              = "ddl_check_pk_name"
	DDLCheckTransactionIsolationLevel           = "ddl_check_transaction_isolation_level"
	DDLCheckTablePartition                      = "ddl_check_table_partition"
	DDLCheckIsExistLimitOffset                  = "ddl_check_is_exist_limit_offset"
	DDLCheckIndexOption                         = "ddl_check_index_option"
	DDLCheckObjectNameUseCN                     = "ddl_check_object_name_using_cn"
	DDLCheckCreateView                          = "ddl_check_create_view"
	DDLCheckCreateTrigger                       = "ddl_check_create_trigger"
	DDLCheckCreateFunction                      = "ddl_check_create_function"
	DDLCheckCreateProcedure                     = "ddl_check_create_procedure"
	DDLCheckTableSize                           = "ddl_check_table_size"
	DDLCheckIndexTooMany                        = "ddl_check_index_too_many"
	DDLCheckRedundantIndex                      = "ddl_check_redundant_index"
	DDLDisableTypeTimestamp                     = "ddl_disable_type_timestamp"
	DDLDisableAlterFieldUseFirstAndAfter        = "ddl_disable_alter_field_use_first_and_after"
	DDLCheckCreateTimeColumn                    = "ddl_check_create_time_column"
)

// inspector DML rules
const (
	DMLCheckWithLimit                    = "dml_check_with_limit"
	DMLCheckWithOrderBy                  = "dml_check_with_order_by"
	DMLCheckWhereIsInvalid               = "all_check_where_is_invalid"
	DMLDisableSelectAllColumn            = "dml_disable_select_all_column"
	DMLCheckInsertColumnsExist           = "dml_check_insert_columns_exist"
	DMLCheckBatchInsertListsMax          = "dml_check_batch_insert_lists_max"
	DMLCheckInQueryNumber                = "dml_check_in_query_limit"
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
	DMLCheckTableSize                    = "dml_check_table_size"
)

// inspector config code
const (
	ConfigDMLRollbackMaxRows       = "dml_rollback_max_rows"
	ConfigDDLOSCMinSize            = "ddl_osc_min_size"
	ConfigDDLGhostMinSize          = "ddl_ghost_min_size"
	ConfigOptimizeIndexEnabled     = "optimize_index_enabled"
	ConfigDMLExplainPreCheckEnable = "dml_enable_explain_pre_check"
)

type RuleHandlerInput struct {
	Ctx  *session.Context
	Rule driver.Rule
	Res  *driver.AuditResult
	Node ast.Node

	require AdditionalRequire
}

type AdditionalRequire struct{}

type RuleHandlerFunc func(input *RuleHandlerInput) error

type RuleHandler struct {
	Rule                 driver.Rule
	Message              string
	Func                 RuleHandlerFunc
	AllowOffline         bool
	NotAllowOfflineStmts []ast.Node
}

// In order to reuse some code, some rules use the same rule handler.
// Then following code is the side effect of the purpose.
//
// It's not a good idea to use the same rule handler for different rules.
// FIXME: once we map one rule to one rule handler, we should remove the side effect.
func addResult(result *driver.AuditResult, currentRule driver.Rule, ruleName string, args ...interface{}) {
	// if rule is not current rule, ignore save the message.
	if ruleName != currentRule.Name {
		return
	}
	level := currentRule.Level
	message := RuleHandlerMap[ruleName].Message
	result.Add(level, message, args...)
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
)

const DefaultSingleParamKeyName = "first_key" // For most of the rules, it is just has one param, this is first params.

const (
	DefaultMultiParamsFirstKeyName  = "multi_params_first_key"
	DefaultMultiParamsSecondKeyName = "multi_params_second_key"
)

var RuleHandlers = []RuleHandler{
	// config
	{
		Rule: driver.Rule{
			Name: ConfigDMLRollbackMaxRows,
			Desc: "在 DML 语句中预计影响行数超过指定值则不回滚",
			//Value:    "1000",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeGlobalConfig,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "1000",
					Desc:  "最大影响行数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Func: nil,
	},
	{
		Rule: driver.Rule{
			Name: ConfigDDLOSCMinSize,
			Desc: "改表时，表空间超过指定大小(MB)审核时输出osc改写建议",
			//Value:    "16",
			Level:    driver.RuleLevelNormal,
			Category: RuleTypeGlobalConfig,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "16",
					Desc:  "表空间大小（MB）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Func: nil,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTableSize,
			Desc:     "检查DDL操作的表是否超过指定数据量",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "16",
					Desc:  "表空间大小（MB）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "执行DDL的表 %v 空间超过 %vMB",
		Func:    checkDDLTableSize,
	}, {
		Rule: driver.Rule{
			Name:     DDLCheckIndexTooMany,
			Desc:     "检查DDL创建的新索引对应字段是否已存在过多索引",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeIndexingConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "2",
					Desc:  "单字段的索引数最大值",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "字段 %v 上的索引数量超过%v个",
		Func:    checkIndex,
	},
	{
		Rule: driver.Rule{
			Name:     ConfigDMLExplainPreCheckEnable,
			Desc:     "使用explain加强预检查能力",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeGlobalConfig,
		},
		Func: nil,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckRedundantIndex,
			Desc:     "检查DDL是否创建冗余的索引",
			Level:    driver.RuleLevelError,
			Category: RuleTypeIndexOptimization,
		},
		Message:      "%v",
		AllowOffline: true,
		Func:         checkIndex,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckTableSize,
			Desc:     "检查DML操作的表是否超过指定数据量",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "16",
					Desc:  "表空间大小（MB）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "执行DML的表 %v 空间超过 %vMB",
		Func:    checkDMLTableSize,
	},

	{
		Rule: driver.Rule{
			Name:     ConfigOptimizeIndexEnabled,
			Desc:     "索引创建建议",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeIndexOptimization,
			Params: params.Params{
				&params.Param{
					Key:   DefaultMultiParamsFirstKeyName,
					Value: "1000000",
					Desc:  "计算列基数阈值",
					Type:  params.ParamTypeInt,
				},
				&params.Param{
					Key:   DefaultMultiParamsSecondKeyName,
					Value: "3",
					Desc:  "联合索引最大列数",
					Type:  params.ParamTypeInt,
				},
			},
		},
	},

	{
		Rule: driver.Rule{
			Name: ConfigDDLGhostMinSize,
			Desc: "改表时，表空间超过指定大小(MB)时使用gh-ost上线",
			//Value:    "16",
			Level:    driver.RuleLevelNormal,
			Category: RuleTypeGlobalConfig,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "16",
					Desc:  "表空间大小（MB）",
					Type:  params.ParamTypeInt,
				},
			},
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
			Desc:     "表名、列名、索引名的长度不能大于指定字节",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
			//Value:    "64",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "64",
					Desc:  "最大长度（字节）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "表名、列名、索引名的长度不能大于%v字节",
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
			Name:     DDLDisableAlterFieldUseFirstAndAfter,
			Desc:     "alter表字段禁止使用first,after",
			Level:    driver.RuleLevelError,
			Category: RuleTypeDDLConvention,
		},
		Message:      "alter表字段禁止使用first,after",
		AllowOffline: true,
		Func:         disableAlterUseFirstAndAfter,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckCreateTimeColumn,
			Desc:     "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
		},
		Message:      "建表DDL必须包含CREATE_TIME字段且默认值为CURRENT_TIMESTAMP",
		AllowOffline: true,
		Func:         checkFieldCreateTime,
	},
	{
		Rule: driver.Rule{
			Name:  DDLCheckIndexCount,
			Desc:  "索引个数建议不超过阈值",
			Level: driver.RuleLevelNotice,
			//Value:    "5",
			Category: RuleTypeIndexingConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "最大索引个数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:              "索引个数建议不超过%v个",
		AllowOffline:         true,
		NotAllowOfflineStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                 checkIndex,
	},
	{
		Rule: driver.Rule{
			Name:  DDLCheckCompositeIndexMax,
			Desc:  "复合索引的列数量不建议超过阈值",
			Level: driver.RuleLevelNotice,
			//Value:    "3",
			Category: RuleTypeIndexingConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "3",
					Desc:  "最大索引列数量",
					Type:  params.ParamTypeInt,
				},
			},
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
			Name:     DDLCheckObjectNameUseCN,
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
			Name:     DDLCheckTableDBEngine,
			Desc:     "必须使用指定数据库引擎",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
			//Value:    "Innodb",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "Innodb",
					Desc:  "数据库引擎",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "必须使用%v数据库引擎",
		AllowOffline: false,
		Func:         checkEngine,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckTableCharacterSet,
			Desc:     "必须使用指定数据库字符集",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
			//Value:    "utf8mb4",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "utf8mb4",
					Desc:  "数据库字符集",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "必须使用%v数据库字符集",
		AllowOffline: false,
		Func:         checkCharacterSet,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckIndexedColumnWithBlob,
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
			Desc:     "普通索引必须使用固定前缀",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
			//Value:    "idx_",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "idx_",
					Desc:  "索引前缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "普通索引必须要以\"%v\"为前缀",
		AllowOffline: true,
		Func:         checkIndexPrefix,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckUniqueIndexPrefix,
			Desc:     "unique索引必须使用固定前缀",
			Level:    driver.RuleLevelError,
			Category: RuleTypeNamingConvention,
			//Value:    "uniq_",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "uniq_",
					Desc:  "索引前缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "unique索引必须要以\"%v\"为前缀",
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
			Name:     DDLCheckColumnTimestampWithoutDefault,
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
			Name:  DMLCheckBatchInsertListsMax,
			Desc:  "单条insert语句，建议批量插入不超过阈值",
			Level: driver.RuleLevelNotice,
			//Value:    "5000",
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "5000",
					Desc:  "最大插入行数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "单条insert语句，建议批量插入不超过%v条",
		AllowOffline: true,
		Func:         checkDMLWithBatchInsertMaxLimits,
	},
	{
		Rule: driver.Rule{
			Name:     DMLCheckInQueryNumber,
			Desc:     "where条件内in语句中的参数个数不能超过阈值",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "1000",
					Desc:  "in语句参数最大个数",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "in语句中的参数已有%v个，超过阙值%v",
		AllowOffline: true,
		Func:         checkInQueryLimit,
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
			Name:  DDLCheckDatabaseCollation,
			Desc:  "建议使用规定的数据库排序规则",
			Level: driver.RuleLevelNotice,
			//Value:    "utf8mb4_0900_ai_ci",
			Category: RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "utf8mb4_0900_ai_ci",
					Desc:  "数据库排序规则",
					Type:  params.ParamTypeString,
				},
			},
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
			Name:     DDLCheckBigintInsteadOfDecimal,
			Desc:     "建议用BIGINT类型代替DECIMAL",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeDDLConvention,
		},
		Message:      "建议列%s用BIGINT类型代替DECIMAL",
		Func:         checkBigintInsteadOfDecimal,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:  DMLCheckNeedlessFunc,
			Desc:  "避免使用不必要的内置函数",
			Level: driver.RuleLevelNotice,
			//Value:    "sha(),sqrt(),md5()",
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "sha(),sqrt(),md5()",
					Desc:  "指定的函数集合（逗号分割）",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "避免使用不必要的内置函数%v",
		Func:         checkNeedlessFunc,
		AllowOffline: true,
	},
	{
		Rule: driver.Rule{
			Name:     DDLCheckDatabaseSuffix,
			Desc:     "数据库名称必须使用固定后缀结尾",
			Level:    driver.RuleLevelNotice,
			Category: RuleTypeNamingConvention,
			//Value:    "_DB",
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "_DB",
					Desc:  "数据库名称后缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "数据库名称必须以\"%v\"结尾",
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
			Name:  DMLCheckNumberOfJoinTables,
			Desc:  "使用JOIN连接表查询建议不超过阈值",
			Level: driver.RuleLevelNotice,
			//Value:    "3",
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "3",
					Desc:  "最大连接表个数",
					Type:  params.ParamTypeInt,
				},
			},
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
			Name:  DDLCheckIndexOption,
			Desc:  "建议选择可选性超过阈值字段作为索引",
			Level: driver.RuleLevelNotice,
			//Value:    "0.7",
			Category: RuleTypeIndexOptimization,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "70",
					Desc:  "可选择性（百分比）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "索引 %v 未超过可选性阈值 百分之%v, 不建议选为索引",
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
			Name:     DDLCheckColumnSetNotice,
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
			Name: DMLCheckExplainAccessTypeAll,
			//Value:    "10000",
			Desc:     "查询的扫描不建议超过指定行数（默认值：10000）",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "10000",
					Desc:  "最大扫描行数",
					Type:  params.ParamTypeInt,
				},
			},
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
	{
		Rule: driver.Rule{
			Name:     DDLDisableTypeTimestamp,
			Desc:     "禁止使用TIMESTAMP字段",
			Level:    driver.RuleLevelWarn,
			Category: RuleTypeDDLConvention,
		},
		Message:      "禁止使用TIMESTAMP字段",
		AllowOffline: true,
		Func:         disableUseTypeTimestampField,
	},
}

func checkFieldCreateTime(_ *session.Context, rule driver.Rule, result *driver.AuditResult, node ast.Node) error {
	var hasCreateTimeAndDefaultValue bool
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Name.Name.L == "create_time" && hasDefaultValueCurrentTimeStamp(col.Options) {
				hasCreateTimeAndDefaultValue = true
			}
		}
	default:
		return nil
	}

	if !hasCreateTimeAndDefaultValue {
		addResult(result, rule, DDLCheckCreateTimeColumn)
	}

	return nil
}

func hasDefaultValueCurrentTimeStamp(options []*ast.ColumnOption) bool {
	for _, option := range options {
		if option.Tp == ast.ColumnOptionDefaultValue {
			funcCallExpr, ok := option.Expr.(*ast.FuncCallExpr)
			if !ok {
				return false
			}
			if funcCallExpr.FnName.L == "current_timestamp" {
				return true
			}
		}
	}

	return false
}

func checkInQueryLimit(_ *session.Context, rule driver.Rule, result *driver.AuditResult, node ast.Node) error {
	where := getWhereExpr(node)
	if where == nil {
		return nil
	}

	paramThresholdNumber := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	util.ScanWhereStmt(func(expr ast.ExprNode) bool {
		switch stmt := expr.(type) {
		case *ast.PatternInExpr:
			inQueryParamActualNumber := len(stmt.List)
			if inQueryParamActualNumber > paramThresholdNumber {
				addResult(result, rule, DMLCheckInQueryNumber, inQueryParamActualNumber, paramThresholdNumber)
			}
			return true
		}

		return false
	}, where)

	return nil
}

func disableUseTypeTimestampField(_ *session.Context, rule driver.Rule, result *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeTimestamp {
				addResult(result, rule, DDLDisableTypeTimestamp)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn)
		for _, spec := range specs {
			for _, newColumn := range spec.NewColumns {
				if newColumn.Tp.Tp == mysql.TypeTimestamp {
					addResult(result, rule, DDLDisableTypeTimestamp)
					return nil
				}
			}
		}
	}

	return nil
}

func checkBigintInsteadOfDecimal(_ *session.Context, rule driver.Rule, result *driver.AuditResult, node ast.Node) error {
	var columnNames []string
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			if col.Tp.Tp == mysql.TypeNewDecimal {
				columnNames = append(columnNames, col.Name.Name.O)
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn)

		for _, spec := range specs {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				if col.Tp.Tp == mysql.TypeNewDecimal {
					columnNames = append(columnNames, col.Name.Name.O)
				}
			}
		}
	default:
		return nil
	}

	if len(columnNames) > 0 {
		addResult(result, rule, DDLCheckBigintInsteadOfDecimal, strings.Join(columnNames, ","))
	}

	return nil
}

func disableAlterUseFirstAndAfter(_ *session.Context, rule driver.Rule, result *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn)

		for _, spec := range specs {
			if spec.Position == nil {
				continue
			}
			if spec.Position.Tp == ast.ColumnPositionFirst || spec.Position.Tp == ast.ColumnPositionAfter {
				addResult(result, rule, DDLDisableAlterFieldUseFirstAndAfter)
			}
		}
	}

	return nil
}

func init() {
	for _, rh := range RuleHandlers {
		RuleHandlerMap[rh.Rule.Name] = rh
	}
}

func checkSelectAll(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					addResult(res, rule, DMLDisableSelectAllColumn)
				}
			}
		}
	}
	return nil
}

func checkSelectWhere(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {

	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		checkWhere(rule, res, stmt.Where)

	case *ast.UpdateStmt:
		checkWhere(rule, res, stmt.Where)
	case *ast.DeleteStmt:
		checkWhere(rule, res, stmt.Where)
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if checkWhere(rule, res, ss.Where) {
				break
			}
		}
	default:
		return nil
	}

	return nil
}

func checkWhere(rule driver.Rule, res *driver.AuditResult, where ast.ExprNode) bool {
	isAddResult := false

	if where == nil || !util.WhereStmtHasOneColumn(where) {
		addResult(res, rule, DMLCheckWhereIsInvalid)
		isAddResult = true
	}
	if where != nil && util.WhereStmtExistNot(where) {
		addResult(res, rule, DMLCheckWhereExistNot)
		isAddResult = true
	}
	if where != nil && util.WhereStmtExistScalarSubQueries(where) {
		addResult(res, rule, DMLCheckWhereExistScalarSubquery)
		isAddResult = true
	}
	if where != nil && util.CheckWhereFuzzySearch(where) {
		addResult(res, rule, DMLCheckFuzzySearch)
		isAddResult = true
	}
	return isAddResult
}
func checkWhereExistNull(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	if where := getWhereExpr(node); where != nil {
		var existNull bool
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			if _, ok := expr.(*ast.IsNullExpr); ok {
				existNull = true
				return true
			}
			return false
		}, where)
		if existNull {
			addResult(res, rule, rule.Name)
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

func checkIndexesExistBeforeCreatConstraints(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		constraintMap := make(map[string]struct{})
		cols := []string{}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			if spec.Constraint != nil && (spec.Constraint.Tp == ast.ConstraintPrimaryKey ||
				spec.Constraint.Tp == ast.ConstraintUniq || spec.Constraint.Tp == ast.ConstraintUniqKey) {
				for _, key := range spec.Constraint.Keys {
					cols = append(cols, key.Column.Name.String())
				}
			}
		}
		createTableStmt, exist, err := ctx.GetCreateTableStmt(stmt.Table)
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
				addResult(res, rule, DDLCheckIndexesExistBeforeCreateConstraints)
				return nil
			}
		}
	}
	return nil
}

func checkPrimaryKey(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	var pkIsAutoIncrement = false
	var pkIsBigIntUnsigned = false
	inspectCol := func(col *ast.ColumnDef) {
		if util.IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
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
			if util.IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
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
			addResult(res, rule, DDLCheckPKNotExist)
		}
		if hasPk && pkColumnExist && !pkIsAutoIncrement {
			addResult(res, rule, DDLCheckPKWithoutAutoIncrement)
		}
		if hasPk && pkColumnExist && pkIsAutoIncrement {
			addResult(res, rule, DDLCheckPKProhibitAutoIncrement)
		}
		if hasPk && pkColumnExist && !pkIsBigIntUnsigned {
			addResult(res, rule, DDLCheckPKWithoutBigintUnsigned)
		}
	case *ast.AlterTableStmt:
		var alterPK bool
		if originTable, exist, err := ctx.GetCreateTableStmt(stmt.Table); err == nil && exist {
			for _, spec := range stmt.Specs {
				switch spec.Tp {
				case ast.AlterTableAddColumns:
					for _, newColumn := range spec.NewColumns {
						if util.IsAllInOptions(newColumn.Options, ast.ColumnOptionPrimaryKey) {
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

			if originPK, exist := util.GetPrimaryKey(originTable); exist {
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
			addResult(res, rule, DDLCheckPKWithoutAutoIncrement)
		}
		if alterPK && pkIsAutoIncrement {
			addResult(res, rule, DDLCheckPKProhibitAutoIncrement)
		}
		if alterPK && !pkIsBigIntUnsigned {
			addResult(res, rule, DDLCheckPKWithoutBigintUnsigned)
		}
	default:
		return nil
	}
	return nil
}

func checkMergeAlterTable(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := ctx.GetTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables) > 0 {
				addResult(res, rule, DDLCheckAlterTableNeedMerge)
			}
		}
	}
	return nil
}

func checkEngine(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	var tableName *ast.TableName
	var engine string
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
			}
		}
	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, ss := range stmt.Specs {
			for _, op := range ss.Options {
				switch op.Tp {
				case ast.TableOptionEngine:
					engine = op.StrValue
				}
			}
		}
	default:
		return nil
	}
	if engine == "" {
		engine, err = ctx.GetSchemaEngine(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	expectEngine := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if !strings.EqualFold(engine, expectEngine) {
		addResult(res, rule, DDLCheckTableDBEngine, expectEngine)
		return nil
	}
	return nil
}

func getColumnCSFromColumnsDef(columns []*ast.ColumnDef) []string {
	columnCharacterSets := []string{}
	for _, column := range columns {
		// Just string data type and not binary can be set "character set".
		if column.Tp.EvalType() != types.ETString || mysql.HasBinaryFlag(column.Tp.Flag) {
			continue
		}
		if column.Tp.Charset == "" {
			continue
		}
		columnCharacterSets = append(columnCharacterSets, column.Tp.Charset)
	}
	return columnCharacterSets
}

func checkCharacterSet(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	var tableName *ast.TableName
	var characterSet string
	var columnCharacterSets []string

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
			case ast.TableOptionCharset:
				characterSet = op.StrValue
			}
		}
		// https://github.com/actiontech/sqle/issues/389
		// character set can ben defined in columns, like:
		// create table t1 (
		//    id varchar(255) character set utf8
		// )
		columnCharacterSets = getColumnCSFromColumnsDef(stmt.Cols)

	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, ss := range stmt.Specs {
			for _, op := range ss.Options {
				switch op.Tp {
				case ast.TableOptionCharset:
					characterSet = op.StrValue
				}
			}
			// https://github.com/actiontech/sqle/issues/389
			columnCharacterSets = append(columnCharacterSets, getColumnCSFromColumnsDef(ss.NewColumns)...)
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

	if characterSet == "" {
		characterSet, err = ctx.GetSchemaCharacter(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	expectCS := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if !strings.EqualFold(characterSet, expectCS) {
		addResult(res, rule, DDLCheckTableCharacterSet, expectCS)
		return nil
	}
	for _, cs := range columnCharacterSets {
		if !strings.EqualFold(cs, expectCS) {
			addResult(res, rule, DDLCheckTableCharacterSet, expectCS)
			return nil
		}
	}
	return nil
}

func disableAddIndexForColumnsTypeBlob(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	isTypeBlobCols := map[string]bool{}
	indexDataTypeIsBlob := false
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
				if util.HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
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
		createTableStmt, exist, err := ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		// collect columns type from alter table
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn) {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
					if util.HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
						indexDataTypeIsBlob = true
						break
					}
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
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
		createTableStmt, exist, err := ctx.GetCreateTableStmt(stmt.Table)
		if err != nil || !exist {
			return err
		}
		for _, col := range createTableStmt.Cols {
			if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, indexColumns := range stmt.IndexPartSpecifications {
			if isTypeBlobCols[indexColumns.Column.Name.String()] {
				indexDataTypeIsBlob = true
				break
			}
		}
	default:
		return nil
	}
	if indexDataTypeIsBlob {
		addResult(res, rule, DDLCheckIndexedColumnWithBlob)
	}
	return nil
}

func checkNewObjectName(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
	if rule.Name == DDLCheckObjectNameLength {
		length := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
		//length, err := strconv.Atoi(rule.Value)
		//if err != nil {
		//	return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
		//}
		for _, name := range names {
			if len(name) > length {
				addResult(res, rule, DDLCheckObjectNameLength, length)
				break
			}
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

			addResult(res, rule, DDLCheckObjectNameUseCN)
			break
		}
	}

	// check keyword
	invalidNames := []string{}
	for _, name := range names {
		if keyword.IsMysqlReservedKeyword(name) {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		addResult(res, rule, DDLCheckObjectNameUsingKeyword,
			strings.Join(util.RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func checkForeignKey(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
		addResult(res, rule, DDLDisableFK)
	}
	return nil
}

func checkIndex(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	indexCounter := 0
	compositeIndexMax := 0
	singleIndexCounter := map[string] /*index*/ int /*count*/ {}
	tableIndexs, newIndexs := []index{}, []index{}
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
			singleConstraint := index{Name: constraint.Name, Column: []string{}}
			for _, key := range constraint.Keys {
				singleConstraint.Column = append(singleConstraint.Column, key.Column.Name.L)
				singleIndexCounter[key.Column.Name.L]++
			}
			newIndexs = append(newIndexs, singleConstraint)
		}
	case *ast.AlterTableStmt:
		hasAddConstraint := false
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
			switch spec.Tp {
			case ast.AlterTableAddConstraint:
				hasAddConstraint = true
				singleConstraint := index{Name: spec.Constraint.Name, Column: []string{}}
				for _, key := range spec.Constraint.Keys {
					singleConstraint.Column = append(singleConstraint.Column, key.Column.Name.L)
					singleIndexCounter[key.Column.Name.L]++
				}
				newIndexs = append(newIndexs, singleConstraint)
			}
		}
		createTableStmt, exist, err := ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
				singleConstraint := index{Name: constraint.Name, Column: []string{}}
				for _, key := range constraint.Keys {
					singleConstraint.Column = append(singleConstraint.Column, key.Column.Name.L)
					if hasAddConstraint {
						singleIndexCounter[key.Column.Name.L]++
					}
				}
				tableIndexs = append(tableIndexs, singleConstraint)
			}
		}

	case *ast.CreateIndexStmt:
		indexCounter++
		if compositeIndexMax < len(stmt.IndexPartSpecifications) {
			compositeIndexMax = len(stmt.IndexPartSpecifications)
		}
		singleConstraint := index{Name: stmt.IndexName, Column: []string{}}
		for _, key := range stmt.IndexPartSpecifications {
			singleConstraint.Column = append(singleConstraint.Column, key.Column.Name.L)
			singleIndexCounter[key.Column.Name.L]++
		}
		newIndexs = append(newIndexs, singleConstraint)
		createTableStmt, exist, err := ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
				singleConstraint := index{Name: constraint.Name, Column: []string{}}
				for _, key := range constraint.Keys {
					singleConstraint.Column = append(singleConstraint.Column, key.Column.Name.L)
					singleIndexCounter[key.Column.Name.L]++
				}
				tableIndexs = append(tableIndexs, singleConstraint)
			}
		}
	default:
		return nil
	}
	//value, err := strconv.Atoi(rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	//}
	expectCounter := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	if rule.Name == DDLCheckIndexCount && indexCounter > expectCounter {
		addResult(res, rule, DDLCheckIndexCount, expectCounter)
	}
	if rule.Name == DDLCheckCompositeIndexMax && compositeIndexMax > expectCounter {
		addResult(res, rule, DDLCheckCompositeIndexMax, expectCounter)
	}
	if rule.Name == DDLCheckIndexTooMany {
		manyKeys := []string{}
		for s, i := range singleIndexCounter {
			if i > expectCounter {
				manyKeys = append(manyKeys, s)
			}
		}
		if len(manyKeys) > 0 {
			addResult(res, rule, DDLCheckIndexTooMany, strings.Join(manyKeys, " , "), expectCounter)
		}
	}
	if rule.Name == DDLCheckRedundantIndex {
		// here's a false positive
		//nolint:staticcheck
		repeat, redundancy := []string{}, map[string]string{}
		if len(tableIndexs) == 0 {
			repeat, redundancy = checkRedundantIndex(newIndexs)
		} else {
			repeat, redundancy = checkAlterTableRedundantIndex(newIndexs, tableIndexs)
		}

		errStr := ""
		if len(repeat) > 0 {
			errStr = fmt.Sprintf("存在重复索引:%v; ", strings.Join(repeat, " , "))
		}
		for red, source := range redundancy {
			errStr += fmt.Sprintf("已存在索引 %v , 索引 %v 为冗余索引; ", source, red)
		}
		if errStr != "" {
			addResult(res, rule, DDLCheckRedundantIndex, errStr)
		}
	}
	return nil
}

// MySQL column index
type index struct {
	Name   string
	Column []string
}

func (i index) ColumnString() string {
	return strings.Join(i.Column, ",")
}

func (i index) String() string {
	return fmt.Sprintf("%v(%v)", i.Name, i.ColumnString())
}

func checkRedundantIndex(indexs []index) (repeat []string /*column name*/, redundancy map[string] /* redundancy index's column name or index name*/ string /*source column name or index name*/) {
	redundancy = map[string]string{}
	repeat = []string{}
	if len(indexs) == 0 {
		return
	}
	sort.SliceStable(indexs, func(i, j int) bool {
		return indexs[i].ColumnString() < indexs[j].ColumnString()
	})
	lastIndex, lastNormalIndex := indexs[len(indexs)-1], indexs[len(indexs)-1]

	for i := len(indexs) - 2; i >= 0; i-- {
		ind := indexs[i]
		if ind.ColumnString() == lastIndex.ColumnString() &&
			(len(repeat) == 0 || repeat[len(repeat)-1] != ind.String()) {
			repeat = append(repeat, ind.String())
		} else if strings.HasPrefix(lastNormalIndex.ColumnString(), ind.ColumnString()) {
			redundancy[ind.String()] = lastNormalIndex.String()
		} else {
			lastNormalIndex = ind
		}
		lastIndex = ind
	}

	return
}

func checkAlterTableRedundantIndex(newIndexs, tableIndexs []index) (repeat []string /*column name*/, redundancy map[string] /* redundancy index's column name or index name*/ string /*source column name or index name*/) {
	repeat, redundancy = checkRedundantIndex(append(newIndexs, tableIndexs...))

	for i := len(repeat) - 1; i >= 0; i-- {
		hasIndex := false
		for _, newIndex := range newIndexs {
			if newIndex.String() == repeat[i] {
				hasIndex = true
				break
			}
		}
		if !hasIndex {
			repeat = append(repeat[:i], repeat[i+1:]...)
		}
	}

	for r, s := range redundancy {
		hasIndex := false
		for _, newIndex := range newIndexs {
			if r == newIndex.String() || s == newIndex.String() {
				hasIndex = true
				break
			}
		}
		if !hasIndex {
			delete(redundancy, r)
		}
	}

	return
}

func checkStringType(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				addResult(res, rule, DDLCheckColumnCharLength)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					addResult(res, rule, DDLCheckColumnCharLength)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkIfNotExist(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			addResult(res, rule, DDLCheckPKWithoutIfNotExists)
		}
	}
	return nil
}

func checkDDLTableSize(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	min := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		tables = append(tables, stmt.Table)
	case *ast.DropTableStmt:
		tables = append(tables, stmt.Tables...)
	default:
		return nil
	}

	beyond := []string{}
	for _, table := range tables {
		size, err := ctx.GetTableSize(table)
		if err != nil {
			return err
		}
		if float64(min) < size {
			beyond = append(beyond, table.Name.String())
		}
	}

	if len(beyond) > 0 {
		addResult(res, rule, rule.Name, strings.Join(beyond, " , "), min)
	}
	return nil
}

func checkDMLTableSize(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	min := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			return nil
		}
		tables = append(tables, util.GetTables(stmt.From.TableRefs)...)
	case *ast.InsertStmt:
		tables = append(tables, util.GetTables(stmt.Table.TableRefs)...)
	case *ast.UpdateStmt:
		tables = append(tables, util.GetTables(stmt.TableRefs.TableRefs)...)
	case *ast.DeleteStmt:
		tables = append(tables, util.GetTables(stmt.TableRefs.TableRefs)...)
		if stmt.Tables != nil {
			tables = append(tables, stmt.Tables.Tables...)
		}
	case *ast.LockTablesStmt:
		for _, lock := range stmt.TableLocks {
			tables = append(tables, lock.Table)
		}
	default:
		return nil
	}

	beyond := []string{}
	for _, table := range tables {
		size, err := ctx.GetTableSize(table)
		if err != nil {
			return err
		}
		if float64(min) < size {
			beyond = append(beyond, table.Name.String())
		}
	}

	if len(beyond) > 0 {
		addResult(res, rule, rule.Name, strings.Join(beyond, " , "), min)
	}
	return nil
}

func disableDropStmt(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		addResult(res, rule, DDLDisableDropStatement)
	case *ast.DropTableStmt:
		addResult(res, rule, DDLDisableDropStatement)
	}
	return nil
}

func checkTableWithoutComment(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
			addResult(res, rule, DDLCheckTableWithoutComment)
		}
	}
	return nil
}

func checkColumnWithoutComment(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
				addResult(res, rule, DDLCheckColumnWithoutComment)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasComment := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionComment {
						columnHasComment = true
					}
				}
				if !columnHasComment {
					addResult(res, rule, DDLCheckColumnWithoutComment)
					return nil
				}
			}
		}
	}
	return nil
}

func checkIndexPrefix(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex:
				indexesName = append(indexesName, spec.Constraint.Name)
			}
		}
	case *ast.CreateIndexStmt:
		if stmt.KeyType == ast.IndexKeyTypeNone {
			indexesName = append(indexesName, stmt.IndexName)
		}
	default:
		return nil
	}
	prefix := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	for _, name := range indexesName {
		if !utils.HasPrefix(name, prefix, false) {
			addResult(res, rule, DDLCheckIndexPrefix, prefix)
			return nil
		}
	}
	return nil
}

func checkUniqIndexPrefix(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	_, indexes := getTableUniqIndex(node)
	prefix := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	for index := range indexes {
		if !utils.HasPrefix(index, prefix, false) {
			addResult(res, rule, DDLCheckUniqueIndexPrefix, prefix)
			return nil
		}
	}
	return nil
}

func checkUniqIndex(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	tableName, indexes := getTableUniqIndex(node)
	for index, indexedCols := range indexes {
		if !strings.EqualFold(index, fmt.Sprintf("IDX_UK_%v_%v", tableName, strings.Join(indexedCols, "_"))) {
			addResult(res, rule, DDLCheckUniqueIndex)
			return nil
		}
	}
	return nil
}

func getTableUniqIndex(node ast.Node) (string, map[string][]string) {
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
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintUniq:
				for _, key := range spec.Constraint.Keys {
					indexes[spec.Constraint.Name] = append(indexes[spec.Constraint.Name], key.Column.Name.String())
				}
			}
		}
	case *ast.CreateIndexStmt:
		tableName = stmt.Table.Name.String()
		if stmt.KeyType == ast.IndexKeyTypeUnique {
			for _, indexCol := range stmt.IndexPartSpecifications {
				indexes[stmt.IndexName] = append(indexes[stmt.IndexName], indexCol.Column.Name.String())
			}
		}
	default:
	}
	return tableName, indexes
}

func checkColumnWithoutDefault(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
			if util.HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
				isAutoIncrementColumn = true
			}
			if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
				isBlobColumn = true
			}
			if util.HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
				columnHasDefault = true
			}
			if isAutoIncrementColumn || isBlobColumn {
				continue
			}
			if !columnHasDefault {
				addResult(res, rule, DDLCheckColumnWithoutDefault)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				isAutoIncrementColumn := false
				isBlobColumn := false
				columnHasDefault := false

				if util.HasOneInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					isAutoIncrementColumn = true
				}
				if util.MysqlDataTypeIsBlob(col.Tp.Tp) {
					isBlobColumn = true
				}
				if util.HasOneInOptions(col.Options, ast.ColumnOptionDefaultValue) {
					columnHasDefault = true
				}

				if isAutoIncrementColumn || isBlobColumn {
					continue
				}
				if !columnHasDefault {
					addResult(res, rule, DDLCheckColumnWithoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnTimestampWithoutDefault(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
				addResult(res, rule, DDLCheckColumnTimestampWithoutDefault)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn) {
			for _, col := range spec.NewColumns {
				columnHasDefault := false
				for _, op := range col.Options {
					if op.Tp == ast.ColumnOptionDefaultValue {
						columnHasDefault = true
					}
				}
				if !columnHasDefault && (col.Tp.Tp == mysql.TypeTimestamp || col.Tp.Tp == mysql.TypeDatetime) {
					addResult(res, rule, DDLCheckColumnTimestampWithoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnBlobNotNull(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
						addResult(res, rule, DDLCheckColumnBlobWithNotNull)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableChangeColumn,
			ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionNotNull {
							addResult(res, rule, DDLCheckColumnBlobWithNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkColumnEnumNotice(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, res, node, mysql.TypeEnum)
}

func checkColumnSetNotice(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, res, node, mysql.TypeSet)
}

func checkColumnBlobNotice(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	return checkColumnShouldNotBeType(rule, res, node, mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob)
}

func checkColumnShouldNotBeType(rule driver.Rule, res *driver.AuditResult, node ast.Node, colTypes ...byte) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col == nil {
				continue
			}
			if bytes.Contains(colTypes, []byte{col.Tp.Tp}) {
				addResult(res, rule, rule.Name)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(
			stmt.Specs,
			ast.AlterTableAddColumns,
			ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn) {

			for _, newCol := range spec.NewColumns {
				if newCol.Tp == nil {
					continue
				}

				if bytes.Contains(colTypes, []byte{newCol.Tp.Tp}) {
					addResult(res, rule, rule.Name)
					return nil
				}
			}
		}
	}

	return nil
}

func checkColumnBlobDefaultNull(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
						addResult(res, rule, DDLCheckColumnBlobDefaultIsNotNull)
						return nil
					}
				}
			}
		}
	case *ast.AlterTableStmt:
		if stmt.Specs == nil {
			return nil
		}
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableModifyColumn, ast.AlterTableAlterColumn,
			ast.AlterTableChangeColumn, ast.AlterTableAddColumns) {
			for _, col := range spec.NewColumns {
				if col.Tp == nil {
					continue
				}
				switch col.Tp.Tp {
				case mysql.TypeBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeLongBlob:
					for _, opt := range col.Options {
						if opt.Tp == ast.ColumnOptionDefaultValue && opt.Expr.GetType().Tp != mysql.TypeNull {
							addResult(res, rule, DDLCheckColumnBlobDefaultIsNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkDMLWithLimit(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			addResult(res, rule, DMLCheckWithLimit)
		}
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			addResult(res, rule, DMLCheckWithLimit)
		}
	}
	return nil
}
func checkDMLLimitExist(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit == nil {
			addResult(res, rule, DMLCheckLimitMustExist)
		}
	case *ast.DeleteStmt:
		if stmt.Limit == nil {
			addResult(res, rule, DMLCheckLimitMustExist)
		}
	}
	return nil
}

func checkDMLWithOrderBy(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			addResult(res, rule, DMLCheckWithOrderBy)
		}
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			addResult(res, rule, DMLCheckWithOrderBy)
		}
	}
	return nil
}

func checkDMLWithInsertColumnExist(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Columns) == 0 {
			addResult(res, rule, DMLCheckInsertColumnsExist)
		}
	}
	return nil
}

func checkDMLWithBatchInsertMaxLimits(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	max := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	//value, err := strconv.Atoi(rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	//}
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Lists) > max {
			addResult(res, rule, DMLCheckBatchInsertListsMax, max)
		}
	}
	return nil
}

func checkWhereExistFunc(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.From.TableRefs)
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
			checkExistFunc(ctx, rule, res, tables, stmt.Where)
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkExistFunc(ctx, rule, res, tables, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			checkExistFunc(ctx, rule, res, util.GetTables(stmt.TableRefs.TableRefs), stmt.Where)
		}
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			tableSources := util.GetTableSources(ss.From.TableRefs)
			if len(tableSources) < 1 {
				continue
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			if checkExistFunc(ctx, rule, res, tables, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}

func checkExistFunc(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, tables []*ast.TableName, where ast.ExprNode) bool {
	if where == nil {
		return false
	}
	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := ctx.GetCreateTableStmt(tableName)
		if exist && err == nil {
			cols = append(cols, createTableStmt.Cols...)
		}
	}
	colMap := make(map[string]struct{})
	for _, col := range cols {
		colMap[col.Name.String()] = struct{}{}
	}
	if util.IsFuncUsedOnColumnInWhereStmt(colMap, where) {
		addResult(res, rule, DMLCheckWhereExistFunc)
		return true
	}
	return false
}

func checkWhereColumnImplicitConversion(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	tables := []*ast.TableName{}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.From.TableRefs)
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
			checkWhereColumnImplicitConversionFunc(ctx, rule, res, tables, stmt.Where)
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			checkWhereColumnImplicitConversionFunc(ctx, rule, res, tables, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			checkWhereColumnImplicitConversionFunc(ctx, rule, res, util.GetTables(stmt.TableRefs.TableRefs), stmt.Where)
		}
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			tableSources := util.GetTableSources(ss.From.TableRefs)
			if len(tableSources) < 1 {
				continue
			}
			for _, tableSource := range tableSources {
				switch source := tableSource.Source.(type) {
				case *ast.TableName:
					tables = append(tables, source)
				}
			}
			if checkWhereColumnImplicitConversionFunc(ctx, rule, res, tables, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}
func checkWhereColumnImplicitConversionFunc(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, tables []*ast.TableName, where ast.ExprNode) bool {
	if where == nil {
		return false
	}
	var cols []*ast.ColumnDef
	for _, tableName := range tables {
		createTableStmt, exist, err := ctx.GetCreateTableStmt(tableName)
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
	if util.IsColumnImplicitConversionInWhereStmt(colMap, where) {
		addResult(res, rule, DMLCheckWhereExistImplicitConversion)
		return true
	}
	return false
}

func checkDMLSelectForUpdate(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.LockTp == ast.SelectLockForUpdate {
			addResult(res, rule, DMLCheckSelectForUpdate)
		}
	}
	return nil
}

func getColumnCollationsFromColumnsDef(columns []*ast.ColumnDef) []string {
	columnCollations := []string{}
	for _, column := range columns {
		for _, op := range column.Options {
			if op.Tp == ast.ColumnOptionCollate {
				columnCollations = append(columnCollations, op.StrValue)
				break
			}
		}
	}
	return columnCollations
}

func checkCollationDatabase(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	var collationDatabase string
	var columnCollations []string
	var err error

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		tableName := stmt.Table
		if stmt.ReferTable != nil {
			return nil
		}
		for _, op := range stmt.Options {
			if op.Tp == ast.TableOptionCollate {
				collationDatabase = op.StrValue
				break
			}
		}
		// if create table not define collation, using default.
		if collationDatabase == "" {
			collationDatabase, err = ctx.GetCollationDatabase(tableName, "")
			if err != nil {
				return err
			}
		}

		// https://github.com/actiontech/sqle/issues/443
		// character set can ben defined in columns, like:
		// create table t1 (
		//    id varchar(255) collate utf8mb4_bin
		// )
		columnCollations = getColumnCollationsFromColumnsDef(stmt.Cols)

	case *ast.AlterTableStmt:
		for _, ss := range stmt.Specs {
			for _, op := range ss.Options {
				if op.Tp == ast.TableOptionCollate {
					collationDatabase = op.StrValue
					break
				}
			}
			// https://github.com/actiontech/sqle/issues/443
			columnCollations = append(columnCollations, getColumnCollationsFromColumnsDef(ss.NewColumns)...)
		}
	case *ast.CreateDatabaseStmt:
		schemaName := stmt.Name
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCollate {
				collationDatabase = ss.Value
				break
			}
		}
		// if create schema not define collation, using default.
		if collationDatabase == "" {
			collationDatabase, err = ctx.GetCollationDatabase(nil, schemaName)
			if err != nil {
				return err
			}
		}
	case *ast.AlterDatabaseStmt:
		for _, ss := range stmt.Options {
			if ss.Tp == ast.DatabaseOptionCollate {
				collationDatabase = ss.Value
				break
			}
		}
	default:
		return nil
	}
	expectCollation := rule.Params.GetParam(DefaultSingleParamKeyName).String()

	// if collationDatabase empty, it means that we are not "create object"
	// and collation not change in "update object", so don't to check it.
	if collationDatabase != "" && !strings.EqualFold(collationDatabase, expectCollation) {
		addResult(res, rule, DDLCheckDatabaseCollation, expectCollation)
	}

	for _, cs := range columnCollations {
		if !strings.EqualFold(cs, expectCollation) {
			addResult(res, rule, DDLCheckDatabaseCollation, expectCollation)
			return nil
		}
	}
	return nil
}
func checkDecimalTypeColumn(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
				addResult(res, rule, DDLCheckDecimalTypeColumn)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
					addResult(res, rule, DDLCheckDecimalTypeColumn)
					return nil
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkNeedlessFunc(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	funcArrStr := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	needlessFuncArr := strings.Split(funcArrStr, ",")
	sql := strings.ToLower(node.Text())
	for _, needlessFunc := range needlessFuncArr {
		needlessFunc = strings.ToLower(strings.TrimRight(needlessFunc, ")"))
		if strings.Contains(sql, needlessFunc) {
			addResult(res, rule, DMLCheckNeedlessFunc, funcArrStr)
			return nil
		}
	}
	return nil
}

func checkDatabaseSuffix(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	databaseName := ""
	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		databaseName = stmt.Name
	case *ast.AlterDatabaseStmt:
		databaseName = stmt.Name
	default:
		return nil
	}
	suffix := rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if databaseName != "" && !utils.HasSuffix(databaseName, suffix, false) {
		addResult(res, rule, DDLCheckDatabaseSuffix, suffix)
		return nil
	}
	return nil
}

func checkPKIndexName(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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
		addResult(res, rule, DDLCheckPKName)
		return nil
	}
	return nil
}

func checkTransactionIsolationLevel(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SetStmt:
		for _, variable := range stmt.Variables {
			if dry.StringListContains([]string{"tx_isolation", "tx_isolation_one_shot"}, variable.Name) {
				switch node := variable.Value.(type) {
				case *parserdriver.ValueExpr:
					if node.Datum.GetString() != ast.ReadCommitted {
						addResult(res, rule, DDLCheckTransactionIsolationLevel)
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

func checkTablePartition(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.PartitionNames != nil || spec.PartDefinitions != nil || spec.Partition != nil {
				addResult(res, rule, DDLCheckTablePartition)
				return nil
			}
		}
	case *ast.CreateTableStmt:
		if stmt.Partition != nil {
			addResult(res, rule, DDLCheckTablePartition)
			return nil
		}
	default:
		return nil
	}
	return nil
}
func checkNumberOfJoinTables(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	nums := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	//nums, err := strconv.Atoi(rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing rule[%v] value error: %v", rule.Name, err)
	//}
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		if nums < util.GetNumberOfJoinTables(stmt.From.TableRefs) {
			addResult(res, rule, DMLCheckNumberOfJoinTables, nums)
		}
	default:
		return nil
	}
	return nil
}

func checkIsAfterUnionDistinct(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if ss.IsAfterUnionDistinct {
				addResult(res, rule, DMLCheckIfAfterUnionDistinct)
				return nil
			}
		}
	default:
		return nil
	}

	return nil
}

func checkIsExistLimitOffset(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.Limit.Offset != nil {
			addResult(res, rule, DDLCheckIsExistLimitOffset)
		}
	default:
		return nil
	}
	return nil
}

func checkIndexOption(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {

	var tableName *ast.TableName
	indexColumns := make([]string, 0)
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		tableName = stmt.Table
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			if spec.Constraint == nil {
				continue
			}
			for _, key := range spec.Constraint.Keys {
				indexColumns = append(indexColumns, key.Column.Name.String())
			}
		}
	case *ast.CreateIndexStmt:
		tableName = stmt.Table
		for _, indexCol := range stmt.IndexPartSpecifications {
			indexColumns = append(indexColumns, indexCol.Column.Name.String())
		}
	default:
		return nil
	}
	if len(indexColumns) == 0 {
		return nil
	}
	maxIndexOption, err := ctx.GetMaxIndexOptionForTable(tableName, indexColumns)
	if err != nil {
		return err
	}
	// todo: using number compare, don't use string compare
	max := rule.Params.GetParam(DefaultSingleParamKeyName).Int()

	if maxIndexOption > 0 && float64(max) > maxIndexOption {
		addResult(res, rule, rule.Name, strings.Join(indexColumns, ", "), max)
	}
	return nil
}

func checkExplain(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
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

	epRecords, err := ctx.GetExecutionPlan(node.Text())
	if err != nil {
		// TODO: check dml related table or database is created, if not exist, explain will executed failure.
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", node.Text(), err)
		return nil
	}
	for _, record := range epRecords {
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingFilesort) {
			addResult(res, rule, DMLCheckExplainExtraUsingFilesort)
		}
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingTemporary) {
			addResult(res, rule, DMLCheckExplainExtraUsingTemporary)
		}

		//defaultRule := RuleHandlerMap[DMLCheckExplainAccessTypeAll].Rule
		max := rule.Params.GetParam(DefaultSingleParamKeyName).Int()
		if record.Type == executor.ExplainRecordAccessTypeAll && record.Rows > int64(max) {
			addResult(res, rule, DMLCheckExplainAccessTypeAll, record.Rows)
		}
	}
	return nil
}

func checkCreateView(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case *ast.CreateViewStmt:
		addResult(res, rule, rule.Name)
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
func checkCreateTrigger(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createTriggerReg1.MatchString(node.Text()) ||
			createTriggerReg2.MatchString(node.Text()) {
			addResult(res, rule, rule.Name)
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
func checkCreateFunction(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createFunctionReg1.MatchString(node.Text()) ||
			createFunctionReg2.MatchString(node.Text()) {
			addResult(res, rule, rule.Name)
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
func checkCreateProcedure(ctx *session.Context, rule driver.Rule, res *driver.AuditResult, node ast.Node) error {
	switch node.(type) {
	case *ast.UnparsedStmt:
		if createProcedureReg1.MatchString(node.Text()) ||
			createProcedureReg2.MatchString(node.Text()) {
			addResult(res, rule, rule.Name)
		}
	}
	return nil
}
