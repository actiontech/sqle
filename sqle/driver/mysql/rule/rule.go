package rule

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/keyword"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/parser/types"
	tidbTypes "github.com/pingcap/tidb/types"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	dry "github.com/ungerik/go-dry"
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
	RuleTypeIndexInvalidation  = "索引失效"
)

const (
	AllCheckPrepareStatementPlaceholders = "all_check_prepare_statement_placeholders"
)

// inspector DDL rules
const (
	DDLCheckPKWithoutIfNotExists                       = "ddl_check_table_without_if_not_exists"
	DDLCheckObjectNameLength                           = "ddl_check_object_name_length"
	DDLCheckObjectNameUsingKeyword                     = "ddl_check_object_name_using_keyword"
	DDLCheckPKNotExist                                 = "ddl_check_pk_not_exist"
	DDLCheckPKWithoutBigintUnsigned                    = "ddl_check_pk_without_bigint_unsigned"
	DDLCheckPKWithoutAutoIncrement                     = "ddl_check_pk_without_auto_increment"
	DDLCheckPKProhibitAutoIncrement                    = "ddl_check_pk_prohibit_auto_increment"
	DDLCheckColumnCharLength                           = "ddl_check_column_char_length"
	DDLDisableFK                                       = "ddl_disable_fk"
	DDLCheckIndexCount                                 = "ddl_check_index_count"
	DDLCheckIndexNotNullConstraint                     = "ddl_check_index_not_null_constraint"
	DDLCheckCompositeIndexMax                          = "ddl_check_composite_index_max"
	DDLCheckTableDBEngine                              = "ddl_check_table_db_engine"
	DDLCheckTableCharacterSet                          = "ddl_check_table_character_set"
	DDLCheckIndexedColumnWithBlob                      = "ddl_check_index_column_with_blob"
	DDLCheckAlterTableNeedMerge                        = "ddl_check_alter_table_need_merge"
	DDLDisableDropStatement                            = "ddl_disable_drop_statement"
	DDLCheckTableWithoutComment                        = "ddl_check_table_without_comment"
	DDLCheckColumnWithoutComment                       = "ddl_check_column_without_comment"
	DDLCheckIndexPrefix                                = "ddl_check_index_prefix"
	DDLCheckUniqueIndexPrefix                          = "ddl_check_unique_index_prefix"
	DDLCheckUniqueIndex                                = "ddl_check_unique_index"
	DDLCheckColumnWithoutDefault                       = "ddl_check_column_without_default"
	DDLCheckColumnTimestampWithoutDefault              = "ddl_check_column_timestamp_without_default"
	DDLCheckColumnBlobWithNotNull                      = "ddl_check_column_blob_with_not_null"
	DDLCheckColumnBlobDefaultIsNotNull                 = "ddl_check_column_blob_default_is_not_null"
	DDLCheckColumnEnumNotice                           = "ddl_check_column_enum_notice"
	DDLCheckColumnSetNotice                            = "ddl_check_column_set_notice"
	DDLCheckColumnBlobNotice                           = "ddl_check_column_blob_notice"
	DDLCheckIndexesExistBeforeCreateConstraints        = "ddl_check_indexes_exist_before_creat_constraints"
	DDLCheckDatabaseCollation                          = "ddl_check_collation_database"
	DDLCheckDecimalTypeColumn                          = "ddl_check_decimal_type_column"
	DDLCheckBigintInsteadOfDecimal                     = "ddl_check_bigint_instead_of_decimal"
	DDLCheckDatabaseSuffix                             = "ddl_check_database_suffix"
	DDLCheckPKName                                     = "ddl_check_pk_name"
	DDLCheckTransactionIsolationLevel                  = "ddl_check_transaction_isolation_level"
	DDLCheckTablePartition                             = "ddl_check_table_partition"
	DDLCheckIsExistLimitOffset                         = "ddl_check_is_exist_limit_offset"
	DDLCheckIndexOption                                = "ddl_check_index_option"
	DDLCheckObjectNameUseCN                            = "ddl_check_object_name_using_cn"
	DDLCheckCreateView                                 = "ddl_check_create_view"
	DDLCheckCreateTrigger                              = "ddl_check_create_trigger"
	DDLCheckCreateFunction                             = "ddl_check_create_function"
	DDLCheckCreateProcedure                            = "ddl_check_create_procedure"
	DDLCheckTableSize                                  = "ddl_check_table_size"
	DDLCheckIndexTooMany                               = "ddl_check_index_too_many"
	DDLCheckRedundantIndex                             = "ddl_check_redundant_index"
	DDLDisableTypeTimestamp                            = "ddl_disable_type_timestamp"
	DDLDisableAlterFieldUseFirstAndAfter               = "ddl_disable_alter_field_use_first_and_after"
	DDLCheckCreateTimeColumn                           = "ddl_check_create_time_column"
	DDLCheckUpdateTimeColumn                           = "ddl_check_update_time_column"
	DDLHintUpdateTableCharsetWillNotUpdateFieldCharset = "ddl_hint_update_table_charset_will_not_update_field_charset"
	DDLHintDropColumn                                  = "ddl_hint_drop_column"
	DDLHintDropPrimaryKey                              = "ddl_hint_drop_primary_key"
	DDLHintDropForeignKey                              = "ddl_hint_drop_foreign_key"
	DDLCheckFullWidthQuotationMarks                    = "ddl_check_full_width_quotation_marks"
	DDLCheckColumnQuantity                             = "ddl_check_column_quantity"
	DDLRecommendTableColumnCharsetSame                 = "ddl_table_column_charset_same"
	DDLCheckColumnTypeInteger                          = "ddl_check_column_type_integer"
	DDLCheckVarcharSize                                = "ddl_check_varchar_size"
	DDLCheckColumnQuantityInPK                         = "ddl_check_column_quantity_in_pk"
	DDLCheckAutoIncrement                              = "ddl_check_auto_increment"
	DDLNotAllowRenaming                                = "ddl_not_allow_renaming"
	DDLCheckObjectNameIsUpperAndLowerLetterMixed       = "ddl_check_object_name_is_upper_and_lower_letter_mixed"
	DDLCheckFieldNotNUllMustContainDefaultValue        = "ddl_check_field_not_null_must_contain_default_value"
	DDLCheckAutoIncrementFieldNum                      = "ddl_check_auto_increment_field_num"
	DDLCheckAllIndexNotNullConstraint                  = "ddl_check_all_index_not_null_constraint"
	DDLCheckColumnNotNULL                              = "ddl_check_column_not_null"
	DDLCheckTableRows                                  = "ddl_check_table_rows"
	DDLCheckCompositeIndexDistinction                  = "ddl_check_composite_index_distinction"
	DDLAvoidText                                       = "ddl_avoid_text"
	DDLAvoidFullText                                   = "ddl_avoid_full_text"
	DDLAvoidGeometry                                   = "ddl_avoid_geometry"
	DDLAvoidEvent                                      = "ddl_avoid_event"
	DDLCheckCharLength                                 = "ddl_check_char_length"
)

// inspector DML rules
const (
	DMLCheckWithLimit                         = "dml_check_with_limit"
	DMLCheckSelectLimit                       = "dml_check_select_limit"
	DMLCheckWithOrderBy                       = "dml_check_with_order_by"
	DMLCheckSelectWithOrderBy                 = "dml_check_select_with_order_by"
	DMLCheckWhereIsInvalid                    = "all_check_where_is_invalid"
	DMLDisableSelectAllColumn                 = "dml_disable_select_all_column"
	DMLCheckInsertColumnsExist                = "dml_check_insert_columns_exist"
	DMLCheckBatchInsertListsMax               = "dml_check_batch_insert_lists_max"
	DMLCheckInQueryNumber                     = "dml_check_in_query_limit"
	DMLCheckWhereExistFunc                    = "dml_check_where_exist_func"
	DMLCheckWhereExistNot                     = "dml_check_where_exist_not"
	DMLCheckWhereExistImplicitConversion      = "dml_check_where_exist_implicit_conversion"
	DMLCheckLimitMustExist                    = "dml_check_limit_must_exist"
	DMLCheckWhereExistScalarSubquery          = "dml_check_where_exist_scalar_sub_queries"
	DMLWhereExistNull                         = "dml_check_where_exist_null"
	DMLCheckSelectForUpdate                   = "dml_check_select_for_update"
	DMLCheckNeedlessFunc                      = "dml_check_needless_func"
	DMLCheckFuzzySearch                       = "dml_check_fuzzy_search"
	DMLCheckNumberOfJoinTables                = "dml_check_number_of_join_tables"
	DMLCheckIfAfterUnionDistinct              = "dml_check_is_after_union_distinct"
	DMLCheckExplainAccessTypeAll              = "dml_check_explain_access_type_all"
	DMLCheckExplainExtraUsingFilesort         = "dml_check_explain_extra_using_filesort"
	DMLCheckExplainExtraUsingTemporary        = "dml_check_explain_extra_using_temporary"
	DMLCheckTableSize                         = "dml_check_table_size"
	DMLCheckJoinFieldType                     = "dml_check_join_field_type"
	DMLCheckHasJoinCondition                  = "dml_check_join_has_on"
	DMLCheckAlias                             = "dml_check_alias"
	DMLNotRecommendNotWildcardLike            = "dml_not_recommend_not_wildcard_like"
	DMLHintInNullOnlyFalse                    = "dml_hint_in_null_only_false"
	DMLNotRecommendIn                         = "dml_not_recommend_in"
	DMLCheckSpacesAroundTheString             = "dml_check_spaces_around_the_string"
	DMLNotRecommendOrderByRand                = "dml_not_recommend_order_by_rand"
	DMLNotRecommendGroupByConstant            = "dml_not_recommend_group_by_constant"
	DMLCheckSortDirection                     = "dml_check_sort_direction"
	DMLHintGroupByRequiresConditions          = "dml_hint_group_by_requires_conditions"
	DMLNotRecommendGroupByExpression          = "dml_not_recommend_group_by_expression"
	DMLCheckSQLLength                         = "dml_check_sql_length"
	DMLNotRecommendHaving                     = "dml_not_recommend_having"
	DMLHintUseTruncateInsteadOfDelete         = "dml_hint_use_truncate_instead_of_delete"
	DMLNotRecommendUpdatePK                   = "dml_not_recommend_update_pk"
	DMLNotRecommendFuncInWhere                = "dml_not_recommend_func_in_where"
	DMLNotRecommendSysdate                    = "dml_not_recommend_sysdate"
	DMLHintSumFuncTips                        = "dml_hint_sum_func_tips"
	DMLHintCountFuncWithCol                   = "dml_hint_count_func_with_col"
	DMLHintLimitMustBeCombinedWithOrderBy     = "dml_hint_limit_must_be_combined_with_order_by"
	DMLHintTruncateTips                       = "dml_hint_truncate_tips"
	DMLHintDeleteTips                         = "dml_hint_delete_tips"
	DMLCheckSQLInjectionFunc                  = "dml_check_sql_injection_func"
	DMLCheckNotEqualSymbol                    = "dml_check_not_equal_symbol"
	DMLNotRecommendSubquery                   = "dml_not_recommend_subquery"
	DMLCheckSubqueryLimit                     = "dml_check_subquery_limit"
	DMLCheckSubQueryNestNum                   = "dml_check_sub_query_depth"
	DMLCheckExplainFullIndexScan              = "dml_check_explain_full_index_scan"
	DMLCheckExplainExtraUsingIndexForSkipScan = "dml_check_explain_extra_using_index_for_skip_scan"
	DMLCheckAffectedRows                      = "dml_check_affected_rows"
	DMLCheckLimitOffsetNum                    = "dml_check_limit_offset_num"
	DMLCheckUpdateOrDeleteHasWhere            = "dml_check_update_or_delete_has_where"
	DMLCheckSortColumnLength                  = "dml_check_order_by_field_length"
	DMLCheckSameTableJoinedMultipleTimes      = "dml_check_same_table_joined_multiple_times"
	DMLCheckInsertSelect                      = "dml_check_insert_select"
	DMLCheckAggregate                         = "dml_check_aggregate"
	DMLCheckExplainUsingIndex                 = "dml_check_using_index"
	DMLCheckIndexSelectivity                  = "dml_check_index_selectivity"
	DMLCheckSelectRows                        = "dml_check_select_rows"
	DMLCheckScanRows                          = "dml_check_scan_rows"
	DMLMustMatchLeftMostPrefix                = "dml_must_match_left_most_prefix"
	DMLMustUseLeftMostPrefix                  = "dml_must_use_left_most_prefix"
	DMLCheckMathComputationOrFuncOnIndex      = "dml_check_math_computation_or_func_on_index"
	DMLCheckJoinFieldUseIndex                 = "dml_check_join_field_use_index"
	DMLCheckJoinFieldCharacterSetAndCollation = "dml_check_join_field_character_set_Collation"
	DMLSQLExplainLowestLevel                  = "dml_sql_explain_lowest_level"
	DMLAvoidWhereEqualNull                    = "dml_avoid_where_equal_null"
)

// inspector config code
const (
	ConfigDMLRollbackMaxRows       = "dml_rollback_max_rows"
	ConfigDDLOSCMinSize            = "ddl_osc_min_size"
	ConfigDDLGhostMinSize          = "ddl_ghost_min_size"
	ConfigOptimizeIndexEnabled     = "optimize_index_enabled"
	ConfigDMLExplainPreCheckEnable = "dml_enable_explain_pre_check"
	ConfigSQLIsExecuted            = "sql_is_executed"
)

// 计算单位
const (
	TenThousand = 10000
)

type RuleHandlerInput struct {
	Ctx  *session.Context
	Rule driverV2.Rule
	Res  *driverV2.AuditResults
	Node ast.Node
}

type RuleHandlerFunc func(input *RuleHandlerInput) error

type RuleHandler struct {
	Rule                 driverV2.Rule
	Message              string
	Func                 RuleHandlerFunc
	AllowOffline         bool
	NotAllowOfflineStmts []ast.Node
	// 开始事后审核时将会跳过这个值为ture的规则
	OnlyAuditNotExecutedSQL bool
	// 事后审核时将会跳过下方列表中的类型
	NotSupportExecutedSQLAuditStmts []ast.Node
}

func init() {
	defaultRulesKnowledge, err := getDefaultRulesKnowledge()
	if err != nil {
		panic(fmt.Errorf("get default rules knowledge failed: %v", err))
	}
	for i, rh := range RuleHandlers {
		if knowledge, ok := defaultRulesKnowledge[rh.Rule.Name]; ok {
			rh.Rule.Knowledge = driverV2.RuleKnowledge{Content: knowledge}
			RuleHandlers[i] = rh
		}
		RuleHandlerMap[rh.Rule.Name] = rh
	}
}

// In order to reuse some code, some rules use the same rule handler.
// Then following code is the side effect of the purpose.
//
// It's not a good idea to use the same rule handler for different rules.
// FIXME: once we map one rule to one rule handler, we should remove the side effect.
func addResult(result *driverV2.AuditResults, currentRule driverV2.Rule, ruleName string, args ...interface{}) {
	// if rule is not current rule, ignore save the message.
	if ruleName != currentRule.Name {
		return
	}
	level := currentRule.Level
	message := RuleHandlerMap[ruleName].Message
	result.Add(level, ruleName, message, args...)
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

func (rh *RuleHandler) IsDisableExecutedSQLRule(node ast.Node) bool {
	for _, stmt := range rh.NotSupportExecutedSQLAuditStmts {
		if reflect.TypeOf(stmt) == reflect.TypeOf(node) {
			return true
		}
	}
	return false
}

var (
	RuleHandlerMap = map[string]RuleHandler{}
)

const DefaultSingleParamKeyName = "first_key" // For most of the rules, it is just has one param, this is first params.

const (
	DefaultMultiParamsFirstKeyName  = "multi_params_first_key"
	DefaultMultiParamsSecondKeyName = "multi_params_second_key"
)

func checkMathComputationOrFuncOnIndex(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		selectStmtExtractor := util.SelectStmtExtractor{}
		stmt.Accept(&selectStmtExtractor)

		for _, selectStmt := range selectStmtExtractor.SelectStmts {
			if ExistMathComputationOrFuncOnIndex(input, selectStmt, selectStmt.Where) {
				addResult(input.Res, input.Rule, DMLCheckMathComputationOrFuncOnIndex)
			}
		}
	case *ast.UpdateStmt:
		if ExistMathComputationOrFuncOnIndex(input, stmt, stmt.Where) {
			addResult(input.Res, input.Rule, DMLCheckMathComputationOrFuncOnIndex)
		}
	case *ast.DeleteStmt:
		if ExistMathComputationOrFuncOnIndex(input, stmt, stmt.Where) {
			addResult(input.Res, input.Rule, DMLCheckMathComputationOrFuncOnIndex)
		}
	default:
		return nil
	}

	return nil
}

func ExistMathComputationOrFuncOnIndex(input *RuleHandlerInput, node ast.Node, whereClause ast.ExprNode) bool {
	if whereClause == nil {
		return false
	}

	tableNameExtractor := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
	node.Accept(&tableNameExtractor)

	indexNameMap := make(map[string]struct{})
	for _, tableName := range tableNameExtractor.TableNames {
		schemaName := input.Ctx.GetSchemaName(tableName)
		indexesInfo, err := input.Ctx.GetTableIndexesInfo(schemaName, tableName.Name.String())
		if err != nil {
			continue
		}

		for _, indexInfo := range indexesInfo {
			indexNameMap[indexInfo.ColumnName] = struct{}{}
		}
	}

	computationOrFuncExtractor := mathComputationOrFuncExtractor{columnList: make([]*ast.ColumnName, 0)}
	whereClause.Accept(&computationOrFuncExtractor)

	for _, column := range computationOrFuncExtractor.columnList {
		if _, ok := indexNameMap[column.Name.O]; ok {
			return true
		}
	}

	return false
}

type mathComputationOrFuncExtractor struct {
	columnList []*ast.ColumnName
}

func (mc *mathComputationOrFuncExtractor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch stmt := in.(type) {
	case *ast.FuncCallExpr:
		for _, columnNameExpr := range stmt.Args {
			col, ok := columnNameExpr.(*ast.ColumnNameExpr)
			if !ok {
				continue
			}
			mc.columnList = append(mc.columnList, col.Name)
		}
	case *ast.BinaryOperationExpr:
		// https://dev.mysql.com/doc/refman/8.0/en/arithmetic-functions.html
		if !IsMathComputation(stmt) {
			return stmt, false
		}

		if col, ok := stmt.L.(*ast.ColumnNameExpr); ok {
			mc.columnList = append(mc.columnList, col.Name)
		}

		if col, ok := stmt.R.(*ast.ColumnNameExpr); ok {
			mc.columnList = append(mc.columnList, col.Name)
		}
	case *ast.UnaryOperationExpr:
		if stmt.Op == opcode.Minus {
			col, ok := stmt.V.(*ast.ColumnNameExpr)
			if !ok {
				return stmt, false
			}
			mc.columnList = append(mc.columnList, col.Name)
		}
	}

	return in, false
}

func IsMathComputation(stmt *ast.BinaryOperationExpr) bool {
	return stmt.Op == opcode.Plus || stmt.Op == opcode.Minus || stmt.Op == opcode.Mul || stmt.Op == opcode.Div || stmt.Op == opcode.IntDiv || stmt.Op == opcode.Mod
}

func (mc *mathComputationOrFuncExtractor) Leave(in ast.Node) (node ast.Node, ok bool) {
	return in, true
}

func checkFieldNotNUllMustContainDefaultValue(input *RuleHandlerInput) error {
	names := make([]string, 0)

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 获取主键的列名
		// 联合主键的情况，只需要取第一个字段的列名，因为自增字段必须是联合主键的第一个字段，否则建表会报错
		var primaryKeyColName string
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				primaryKeyColName = constraint.Keys[0].Column.Name.O
				break
			}
		}

		for _, col := range stmt.Cols {
			if col.Options == nil {
				continue
			}

			// 跳过主键自增的列，因为主键自增的列不需要设置默认值
			if (isFieldContainColumnOptionType(col, ast.ColumnOptionPrimaryKey) || primaryKeyColName == col.Name.Name.O) &&
				isFieldContainColumnOptionType(col, ast.ColumnOptionAutoIncrement) {
				continue
			}

			if isFieldContainColumnOptionType(col, ast.ColumnOptionNotNull) && !isFieldContainColumnOptionType(col, ast.ColumnOptionDefaultValue) {
				names = append(names, col.Name.Name.String())
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Options == nil {
					continue
				}

				if isFieldContainColumnOptionType(col, ast.ColumnOptionPrimaryKey) && isFieldContainColumnOptionType(col, ast.ColumnOptionAutoIncrement) {
					continue
				}

				if isFieldContainColumnOptionType(col, ast.ColumnOptionNotNull) && !isFieldContainColumnOptionType(col, ast.ColumnOptionDefaultValue) {
					names = append(names, col.Name.Name.String())
				}
			}
		}
	default:
		return nil
	}

	if len(names) > 0 {
		addResult(input.Res, input.Rule, DDLCheckFieldNotNUllMustContainDefaultValue, strings.Join(names, ","))
	}

	return nil
}

func isFieldContainColumnOptionType(field *ast.ColumnDef, optionType ast.ColumnOptionType) bool {
	for _, option := range field.Options {
		if option.Tp == optionType {
			return true
		}
	}
	return false
}

func checkSubQueryNestNum(in *RuleHandlerInput) error {
	if _, ok := in.Node.(ast.DMLNode); ok {
		var maxNestNum int
		subQueryNestNumExtract := util.SubQueryMaxNestNumExtractor{MaxNestNum: &maxNestNum, CurrentNestNum: 1}
		in.Node.Accept(&subQueryNestNumExtract)
		expectNestNum := in.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
		if *subQueryNestNumExtract.MaxNestNum > expectNestNum {
			addResult(in.Res, in.Rule, DMLCheckSubQueryNestNum, expectNestNum)
		}
	}
	return nil
}

func getCreateTableAndOnCondition(input *RuleHandlerInput) (map[string]*ast.CreateTableStmt, []*ast.OnCondition) {
	//nolint:staticcheck
	tableNameCreateTableStmtMap := make(map[string]*ast.CreateTableStmt)
	//nolint:staticcheck
	onConditions := make([]*ast.OnCondition, 0)

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			return nil, nil
		}
		tableNameCreateTableStmtMap = input.Ctx.GetTableNameCreateTableStmtMap(stmt.From.TableRefs)
		onConditions = util.GetTableFromOnCondition(stmt.From.TableRefs)
	case *ast.UpdateStmt:
		if stmt.TableRefs == nil {
			return nil, nil
		}
		tableNameCreateTableStmtMap = input.Ctx.GetTableNameCreateTableStmtMap(stmt.TableRefs.TableRefs)
		onConditions = util.GetTableFromOnCondition(stmt.TableRefs.TableRefs)
	case *ast.DeleteStmt:
		if stmt.TableRefs == nil {
			return nil, nil
		}
		tableNameCreateTableStmtMap = input.Ctx.GetTableNameCreateTableStmtMap(stmt.TableRefs.TableRefs)
		onConditions = util.GetTableFromOnCondition(stmt.TableRefs.TableRefs)
	default:
		return nil, nil
	}
	return tableNameCreateTableStmtMap, onConditions
}

func getCreateTableAndOnConditionForJoinType(input *RuleHandlerInput) (map[string]*ast.CreateTableStmt, []*ast.OnCondition) {
	var ctx *session.Context = input.Ctx
	var joinStmt *ast.Join
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			return nil, nil
		}
		joinStmt = stmt.From.TableRefs
	case *ast.UpdateStmt:
		if stmt.TableRefs == nil {
			return nil, nil
		}
		joinStmt = stmt.TableRefs.TableRefs
	case *ast.DeleteStmt:
		if stmt.TableRefs == nil {
			return nil, nil
		}
		joinStmt = stmt.TableRefs.TableRefs
	default:
		return nil, nil
	}
	tableNameCreateTableStmtMap := getTableNameCreateTableStmtMapForJoinType(ctx, joinStmt)
	onConditions := util.GetTableFromOnCondition(joinStmt)
	return tableNameCreateTableStmtMap, onConditions
}

func checkJoinFieldType(input *RuleHandlerInput) error {
	tableNameCreateTableStmtMap, onConditions := getCreateTableAndOnConditionForJoinType(input)
	if tableNameCreateTableStmtMap == nil && onConditions == nil {
		return nil
	}

	for _, onCondition := range onConditions {
		leftType, rightType := getOnConditionLeftAndRightType(onCondition, tableNameCreateTableStmtMap)
		// 没有类型的情况下不检查
		if leftType == 0 || rightType == 0 {
			continue
		}
		if leftType != rightType {
			addResult(input.Res, input.Rule, DMLCheckJoinFieldType)
		}
	}

	return nil
}

func checkHasJoinCondition(input *RuleHandlerInput) error {
	joinNode := getJoinNodeFromNode(input.Node)
	whereStmt := getWhereExpr(input.Node)
	if joinNode == nil {
		return nil
	}
	joinTables, hasCondition := checkJoinConditionInJoinNode(input.Ctx, whereStmt, joinNode)
	if joinTables && !hasCondition {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func doesNotJoinTables(tableRefs *ast.Join) bool {
	return tableRefs == nil || tableRefs.Left == nil || tableRefs.Right == nil
}

func checkJoinConditionInJoinNode(ctx *session.Context, whereStmt ast.ExprNode, joinNode *ast.Join) (joinTables, hasCondition bool) {
	if joinNode == nil {
		return false, false
	}
	if doesNotJoinTables(joinNode) {
		// 非JOIN两表的JOIN节点 一般是叶子节点 不检查
		return false, false
	}

	// 深度遍历左子树类型为ast.Join的节点 一旦有节点是JOIN两表的节点，并且没有连接条件，则返回
	if l, ok := joinNode.Left.(*ast.Join); ok {
		joinTables, hasCondition = checkJoinConditionInJoinNode(ctx, whereStmt, l)
		if joinTables && !hasCondition {
			return joinTables, hasCondition
		}
	}

	// 判断该节点是否有显式声明连接条件
	if isJoinConditionInOnClause(joinNode) {
		return true, true
	}
	if isJoinConditionInUsingClause(joinNode) {
		return true, true
	}
	if isJoinConditionInWhereStmt(ctx, whereStmt, joinNode) {
		return true, true
	}
	return true, false
}

func isJoinConditionInOnClause(joinNode *ast.Join) bool {
	return joinNode.On != nil
}

func isJoinConditionInUsingClause(joinNode *ast.Join) bool {
	return len(joinNode.Using) > 0
}

func isJoinConditionInWhereStmt(ctx *session.Context, stmt ast.ExprNode, node *ast.Join) bool {
	if stmt == nil {
		return false
	}

	equalConditionVisitor := util.EqualConditionVisitor{}
	stmt.Accept(&equalConditionVisitor)

	for _, column := range equalConditionVisitor.ConditionList {
		/*
			当一个Join节点没有ON或者Using的连接条件时，需要检查Where语句中是否包含连接条件
			Where语句中包含连接条件的判断依据是：
			Where语句中等值条件两侧的不同表的两列，其中一列属于Join右子节点对应的表，另一列属于Join左子树中任意一张表
		*/
		if columnInTable(ctx, node.Right, column.Left) && columnInNode(ctx, node.Left, column.Right) {
			return true
		}
		if columnInTable(ctx, node.Right, column.Right) && columnInNode(ctx, node.Left, column.Left) {
			return true
		}
	}
	return false
}

func columnInTable(ctx *session.Context, node ast.ResultSetNode, columnName *ast.ColumnName) bool {
	if node == nil {
		return false
	}
	switch t := node.(type) {
	case *ast.TableSource:
		return getTableSourceByColumnName(ctx, []*ast.TableSource{t}, columnName) != nil
	}
	return false
}

// 迭代检查表名称是否与JOIN节点中的tableSource的表名或表别名匹配
func columnInNode(ctx *session.Context, node ast.ResultSetNode, columnName *ast.ColumnName) bool {
	if node == nil {
		return false
	}
	switch t := node.(type) {
	case *ast.TableSource:
		return getTableSourceByColumnName(ctx, []*ast.TableSource{t}, columnName) != nil
	case *ast.Join:
		if columnInNode(ctx, t.Right, columnName) {
			return true
		}
		if columnInNode(ctx, t.Left, columnName) {
			return true
		}
	}
	return false
}

func getTableNameCreateTableStmtMapForJoinType(sessionContext *session.Context, joinStmt *ast.Join) map[string]*ast.CreateTableStmt {
	tableNameCreateTableStmtMap := make(map[string]*ast.CreateTableStmt)
	tableSources := util.GetTableSources(joinStmt)
	for _, tableSource := range tableSources {
		tableNameExtractor := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
		tableSource.Source.Accept(&tableNameExtractor)
		if len(tableNameExtractor.TableNames) > 1 {
			log.Logger().Warn("规则:建议JOIN字段类型保持一致,不支持JOIN的表由多表构成")
			continue
		}
		for tableName, tableNameStmt := range tableNameExtractor.TableNames {
			createTableStmt, exist, err := sessionContext.GetCreateTableStmt(tableNameStmt)
			if err != nil || !exist {
				continue
			}
			tableNameCreateTableStmtMap[tableName] = createTableStmt
			// !临时方案：只支持别名对应的临时表只含有一个表，不支持JOIN的表由多表构成
			// TODO AS语句中的别名作为表的别名时，表别名所对应的表可能是数据库的库表，也有可能是语句中构建的临时表。其中，临时表的可能性有很多种，例如：子查询的结果作为表，JOIN得到的表，其中还可能存在层层嵌套的关系。如果要获取到ON语句块中列的实际表名称，需要递归地构建别名:列名:表名(这个表名可能还是别名)的映射关系
			if tableSource.AsName.String() != "" {
				tableNameCreateTableStmtMap[tableSource.AsName.String()] = createTableStmt
			}
			// TODO: 跨库的 JOIN 无法区分
		}
	}
	return tableNameCreateTableStmtMap
}

func getOnConditionLeftAndRightType(onCondition *ast.OnCondition, createTableStmtMap map[string]*ast.CreateTableStmt) (byte, byte) {
	var leftType, rightType byte
	// onCondition在中的ColumnNameExpr.Refer为nil无法索引到原表名和表别名
	if binaryOperation, ok := onCondition.Expr.(*ast.BinaryOperationExpr); ok {
		switch node := binaryOperation.L.(type) {
		// 当使用类型转换时 列的类型被显式转化为对应类型 支持CAST和CONVERT函数
		case *ast.FuncCastExpr:
			leftType = node.Tp.Tp
		default:
			// 默认获取子树的所有列 对应等号一侧 一般连接键只会有一个 不支持多个列的组合
			lVisitor := util.ColumnNameVisitor{}
			binaryOperation.L.Accept(&lVisitor)
			if len(lVisitor.ColumnNameList) > 1 {
				log.Logger().Warn("规则:建议JOIN字段类型保持一致,连接键不支持多个列的组合")
			}
			if len(lVisitor.ColumnNameList) == 1 {
				leftType = getColumnType(lVisitor.ColumnNameList[0], createTableStmtMap)
			}
		}

		switch node := binaryOperation.R.(type) {
		case *ast.FuncCastExpr:
			rightType = node.Tp.Tp
		default:
			rVisitor := util.ColumnNameVisitor{}
			binaryOperation.R.Accept(&rVisitor)
			if len(rVisitor.ColumnNameList) > 1 {
				log.Logger().Warn("规则:建议JOIN字段类型保持一致,连接键不支持多个列的组合")
			}
			if len(rVisitor.ColumnNameList) > 0 {
				rightType = getColumnType(rVisitor.ColumnNameList[0], createTableStmtMap)
			}
		}

	}

	return leftType, rightType
}

func getColumnType(columnName *ast.ColumnNameExpr, createTableStmtMap map[string]*ast.CreateTableStmt) byte {
	var columnType byte
	if createTableStmt, ok := createTableStmtMap[columnName.Name.Table.L]; ok {
		for _, col := range createTableStmt.Cols {
			if col.Tp == nil {
				continue
			}

			if col.Name.Name.L == columnName.Name.Name.L {
				columnType = col.Tp.Tp
			}
		}
	}

	return columnType
}

func checkFieldCreateTime(input *RuleHandlerInput) error {
	var hasCreateTimeAndDefaultValue bool
	createTimeFieldName := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if strings.EqualFold(col.Name.Name.O, createTimeFieldName) && hasDefaultValueCurrentTimeStamp(col.Options) {
				hasCreateTimeAndDefaultValue = true
			}
		}
	default:
		return nil
	}

	if !hasCreateTimeAndDefaultValue {
		addResult(input.Res, input.Rule, DDLCheckCreateTimeColumn, createTimeFieldName)
	}

	return nil
}

func checkSelectWithOrderBy(input *RuleHandlerInput) error {
	var hasOrderBy bool
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.OrderBy != nil {
			hasOrderBy = true
			break
		}

		selectStmtExtractor := util.SelectStmtExtractor{}
		stmt.Accept(&selectStmtExtractor)

		for _, selectStmt := range selectStmtExtractor.SelectStmts {
			if selectStmt.OrderBy != nil {
				hasOrderBy = true
			}
		}
	}

	if hasOrderBy {
		addResult(input.Res, input.Rule, DMLCheckSelectWithOrderBy)
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

func checkInQueryLimit(input *RuleHandlerInput) error {
	dmlNode, ok := input.Node.(ast.DMLNode)
	if !ok {
		return nil
	}

	whereVisitor := &util.WhereVisitor{}
	dmlNode.Accept(whereVisitor)
	paramThresholdNumber := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()

	for _, whereExpr := range whereVisitor.WhereList {
		util.ScanWhereStmt(func(expr ast.ExprNode) bool {
			switch stmt := expr.(type) {
			case *ast.PatternInExpr:
				inQueryParamActualNumber := len(stmt.List)
				if inQueryParamActualNumber > paramThresholdNumber {
					addResult(input.Res, input.Rule, DMLCheckInQueryNumber, inQueryParamActualNumber, paramThresholdNumber)
				}
				return true
			}

			return false
		}, whereExpr)
	}

	return nil
}

func checkFieldUpdateTime(input *RuleHandlerInput) error {
	var hasUpdateTimeAndDefaultValue bool
	updateTimeFieldName := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if strings.EqualFold(col.Name.Name.O, updateTimeFieldName) && hasDefaultValueUpdateTimeStamp(col.Options) {
				hasUpdateTimeAndDefaultValue = true
			}
		}
	default:
		return nil
	}

	if !hasUpdateTimeAndDefaultValue {
		addResult(input.Res, input.Rule, DDLCheckUpdateTimeColumn, updateTimeFieldName)
	}

	return nil
}

func hasDefaultValueUpdateTimeStamp(options []*ast.ColumnOption) bool {
	var hasDefaultCurrentStamp, hasUpdateCurrentTimestamp bool
	for _, option := range options {
		if hasDefaultValueCurrentTimestamp(option) {
			hasDefaultCurrentStamp = true
		}
		if hasUpdateValueCurrentTimestamp(option) {
			hasUpdateCurrentTimestamp = true
		}
	}

	if hasDefaultCurrentStamp && hasUpdateCurrentTimestamp {
		return true
	}

	return false
}

func hasUpdateValueCurrentTimestamp(option *ast.ColumnOption) bool {
	if option.Tp == ast.ColumnOptionOnUpdate {
		funcCallExpr, ok := option.Expr.(*ast.FuncCallExpr)
		if !ok {
			return false
		}

		if funcCallExpr.FnName.L == "current_timestamp" {
			return true
		}
	}

	return false
}

func hasDefaultValueCurrentTimestamp(option *ast.ColumnOption) bool {
	if option.Tp == ast.ColumnOptionDefaultValue {
		funcCallExpr, ok := option.Expr.(*ast.FuncCallExpr)
		if !ok {
			return false
		}

		if funcCallExpr.FnName.L == "current_timestamp" {
			return true
		}
	}

	return false
}

func disableUseTypeTimestampField(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if stmt.Cols == nil {
			return nil
		}
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeTimestamp {
				addResult(input.Res, input.Rule, DDLDisableTypeTimestamp)
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
					addResult(input.Res, input.Rule, DDLDisableTypeTimestamp)
					return nil
				}
			}
		}
	}

	return nil
}

func checkBigintInsteadOfDecimal(input *RuleHandlerInput) error {
	var columnNames []string
	switch stmt := input.Node.(type) {
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
		addResult(input.Res, input.Rule, DDLCheckBigintInsteadOfDecimal, strings.Join(columnNames, ","))
	}

	return nil
}

func disableAlterUseFirstAndAfter(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		specs := util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns, ast.AlterTableModifyColumn,
			ast.AlterTableChangeColumn)

		for _, spec := range specs {
			if spec.Position == nil {
				continue
			}
			if spec.Position.Tp == ast.ColumnPositionFirst || spec.Position.Tp == ast.ColumnPositionAfter {
				addResult(input.Res, input.Rule, DDLDisableAlterFieldUseFirstAndAfter)
			}
		}
	}

	return nil
}

func checkSelectAll(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					addResult(input.Res, input.Rule, DMLDisableSelectAllColumn)
				}
			}
		}
	}
	return nil
}

func isSelectCount(selectStmt *ast.SelectStmt) bool {
	if len(selectStmt.Fields.Fields) == 1 {
		if fu, ok := selectStmt.Fields.Fields[0].Expr.(*ast.AggregateFuncExpr); ok && strings.ToLower(fu.F) == "count" {
			return true
		}
	}
	return false
}

func checkSelectWhere(input *RuleHandlerInput) error {

	visitor := util.WhereVisitor{}
	if input.Rule.Name == DMLCheckWhereIsInvalid {
		visitor.WhetherContainNil = true
	}
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			//If from is null skip check. EX: select 1;select version
			return nil
		}
		if input.Rule.Name == DMLCheckWhereIsInvalid && isSelectCount(stmt) {
			// 只做count()计数，不要求一定有有效的where条件
			return nil
		}
		stmt.Accept(&visitor)
	case *ast.UpdateStmt, *ast.DeleteStmt, *ast.UnionStmt:
		stmt.Accept(&visitor)
	default:
		return nil
	}
	checkWhere(input.Rule, input.Res, visitor.WhereList)

	return nil
}

func checkWhere(rule driverV2.Rule, res *driverV2.AuditResults, whereList []ast.ExprNode) {
	switch rule.Name {
	case DMLCheckWhereIsInvalid:
		if len(whereList) == 0 {
			addResult(res, rule, DMLCheckWhereIsInvalid)
		}
		for _, where := range whereList {
			if where == nil {
				addResult(res, rule, DMLCheckWhereIsInvalid)
				break
			}
			if !util.WhereStmtNotAlwaysTrue(where) {
				addResult(res, rule, DMLCheckWhereIsInvalid)
				break
			}
		}
	case DMLCheckWhereExistNot:
		for _, where := range whereList {
			if util.WhereStmtExistNot(where) {
				addResult(res, rule, DMLCheckWhereExistNot)
				break
			}
		}
	case DMLCheckWhereExistScalarSubquery:
		for _, where := range whereList {
			if util.WhereStmtExistScalarSubQueries(where) {
				addResult(res, rule, DMLCheckWhereExistScalarSubquery)
				break
			}
		}
	case DMLCheckFuzzySearch:
		for _, where := range whereList {
			if util.CheckWhereFuzzySearch(where) {
				addResult(res, rule, DMLCheckFuzzySearch)
				break
			}
		}
	}
}
func checkWhereExistNull(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		var existNull bool
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			if _, ok := expr.(*ast.IsNullExpr); ok {
				existNull = true
				return true
			}
			return false
		}, where)
		if existNull {
			addResult(input.Res, input.Rule, input.Rule.Name)
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

func checkIndexesExistBeforeCreatConstraints(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
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
				addResult(input.Res, input.Rule, DDLCheckIndexesExistBeforeCreateConstraints)
				return nil
			}
		}
	}
	return nil
}

func checkPrimaryKey(input *RuleHandlerInput) error {
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

	switch stmt := input.Node.(type) {
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
			addResult(input.Res, input.Rule, DDLCheckPKNotExist)
		}
		if hasPk && pkColumnExist && !pkIsAutoIncrement {
			addResult(input.Res, input.Rule, DDLCheckPKWithoutAutoIncrement)
		}
		if hasPk && pkColumnExist && pkIsAutoIncrement {
			addResult(input.Res, input.Rule, DDLCheckPKProhibitAutoIncrement)
		}
		if hasPk && pkColumnExist && !pkIsBigIntUnsigned {
			addResult(input.Res, input.Rule, DDLCheckPKWithoutBigintUnsigned)
		}
	case *ast.AlterTableStmt:
		var alterPK bool
		if originTable, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table); err == nil && exist {
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
			addResult(input.Res, input.Rule, DDLCheckPKWithoutAutoIncrement)
		}
		if alterPK && pkIsAutoIncrement {
			addResult(input.Res, input.Rule, DDLCheckPKProhibitAutoIncrement)
		}
		if alterPK && !pkIsBigIntUnsigned {
			addResult(input.Res, input.Rule, DDLCheckPKWithoutBigintUnsigned)
		}
	default:
		return nil
	}
	return nil
}

func checkMergeAlterTable(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		info, exist := input.Ctx.GetTableInfo(stmt.Table)
		if exist {
			if info.AlterTables != nil && len(info.AlterTables) > 0 {
				addResult(input.Res, input.Rule, DDLCheckAlterTableNeedMerge)
			}
		}
	}
	return nil
}

func checkEngine(input *RuleHandlerInput) error {
	var tableName *ast.TableName
	var engine string
	var err error
	schemaName := ""
	switch stmt := input.Node.(type) {
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
		engine, err = input.Ctx.GetSchemaEngine(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	expectEngine := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if !strings.EqualFold(engine, expectEngine) {
		addResult(input.Res, input.Rule, DDLCheckTableDBEngine, expectEngine)
		return nil
	}
	return nil
}

func getSingleColumnCSFromColumnsDef(column *ast.ColumnDef) (string, bool) {
	// Just string data type and not binary can be set "character set".
	if column.Tp == nil || column.Tp.EvalType() != types.ETString || mysql.HasBinaryFlag(column.Tp.Flag) {
		return "", false
	}
	if column.Tp.Charset == "" {
		return "", true
	}
	return column.Tp.Charset, true
}

func getColumnCSFromColumnsDef(columns []*ast.ColumnDef) []string {
	columnCharacterSets := []string{}
	for _, column := range columns {
		charset, hasCharSet := getSingleColumnCSFromColumnsDef(column)
		if !hasCharSet {
			continue
		}
		if charset == "" {
			continue
		}
		columnCharacterSets = append(columnCharacterSets, charset)
	}
	return columnCharacterSets
}

func checkCharacterSet(input *RuleHandlerInput) error {
	var tableName *ast.TableName
	var characterSet string
	var columnCharacterSets []string

	var err error
	schemaName := ""
	switch stmt := input.Node.(type) {
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
		characterSet, err = input.Ctx.GetSchemaCharacter(tableName, schemaName)
		if err != nil {
			return err
		}
	}
	expectCS := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if !strings.EqualFold(characterSet, expectCS) {
		addResult(input.Res, input.Rule, DDLCheckTableCharacterSet, expectCS)
		return nil
	}
	for _, cs := range columnCharacterSets {
		if !strings.EqualFold(cs, expectCS) {
			addResult(input.Res, input.Rule, DDLCheckTableCharacterSet, expectCS)
			return nil
		}
	}
	return nil
}

func disableAddIndexForColumnsTypeBlob(input *RuleHandlerInput) error {
	isTypeBlobCols := map[string]bool{}
	indexDataTypeIsBlob := false
	switch stmt := input.Node.(type) {
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
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
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
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
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
		addResult(input.Res, input.Rule, DDLCheckIndexedColumnWithBlob)
	}
	return nil
}

func checkIsObjectNameUpperAndLowerLetterMixed(input *RuleHandlerInput) error {
	names := getObjectNames(input.Node)

	invalidNames := make([]string, 0)
	for _, name := range names {
		if !utils.IsUpperAndLowerLetterMixed(name) {
			continue
		}
		invalidNames = append(invalidNames, name)
	}

	if len(invalidNames) > 0 {
		addResult(input.Res, input.Rule, DDLCheckObjectNameIsUpperAndLowerLetterMixed, strings.Join(invalidNames, ","))
	}

	return nil
}

func checkNewObjectName(input *RuleHandlerInput) error {
	names := getObjectNames(input.Node)

	// check length
	if input.Rule.Name == DDLCheckObjectNameLength {
		length := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
		//length, err := strconv.Atoi(input.Rule.Value)
		//if err != nil {
		//	return fmt.Errorf("parsing input.Rule[%v] value error: %v", input.Rule.Name, err)
		//}
		for _, name := range names {
			if len(name) > length {
				addResult(input.Res, input.Rule, DDLCheckObjectNameLength, length)
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

			addResult(input.Res, input.Rule, DDLCheckObjectNameUseCN)
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
		addResult(input.Res, input.Rule, DDLCheckObjectNameUsingKeyword,
			strings.Join(util.RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func getObjectNames(node ast.Node) []string {
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

	return names
}

func checkForeignKey(input *RuleHandlerInput) error {
	hasFk := false

	switch stmt := input.Node.(type) {
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
		addResult(input.Res, input.Rule, DDLDisableFK)
	}
	return nil
}

func checkIndex(input *RuleHandlerInput) error {
	indexCounter := 0
	compositeIndexMax := 0
	singleIndexCounter := map[string] /*index*/ int /*count*/ {}
	tableIndexs, newIndexs := []index{}, []index{}
	switch stmt := input.Node.(type) {
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
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
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
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
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
	//value, err := strconv.Atoi(input.Rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing input.Rule[%v] value error: %v", input.Rule.Name, err)
	//}
	expectCounter := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	if input.Rule.Name == DDLCheckIndexCount && indexCounter > expectCounter {
		addResult(input.Res, input.Rule, DDLCheckIndexCount, expectCounter)
	}
	if input.Rule.Name == DDLCheckCompositeIndexMax && compositeIndexMax > expectCounter {
		addResult(input.Res, input.Rule, DDLCheckCompositeIndexMax, expectCounter)
	}
	if input.Rule.Name == DDLCheckIndexTooMany {
		manyKeys := []string{}
		for s, i := range singleIndexCounter {
			if i > expectCounter {
				manyKeys = append(manyKeys, s)
			}
		}
		if len(manyKeys) > 0 {
			addResult(input.Res, input.Rule, DDLCheckIndexTooMany, strings.Join(manyKeys, " , "), expectCounter)
		}
	}
	if input.Rule.Name == DDLCheckRedundantIndex {
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
			addResult(input.Res, input.Rule, DDLCheckRedundantIndex, errStr)
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

func getIndexAndNotNullCols(input *RuleHandlerInput) ([]string, map[string]struct{}, error) {
	indexCols := []string{}
	colsWithNotNullConstraint := make(map[string] /*column name*/ struct{})

	checkNewColumns := func(newColumns []*ast.ColumnDef) {
		for _, column := range newColumns {
			hasNotNull, hasIndex := false, false
			for _, option := range column.Options {
				switch option.Tp {
				case ast.ColumnOptionUniqKey, ast.ColumnOptionPrimaryKey:
					hasIndex = true
				case ast.ColumnOptionNotNull:
					hasNotNull = true
				}
			}
			if hasIndex && !hasNotNull {
				indexCols = append(indexCols, column.Name.Name.L)
			}
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			for _, option := range col.Options {
				switch option.Tp {
				case ast.ColumnOptionNotNull:
					colsWithNotNullConstraint[col.Name.Name.L] = struct{}{}
				case ast.ColumnOptionPrimaryKey, ast.ColumnOptionUniqKey:
					indexCols = append(indexCols, col.Name.Name.L)
				}
			}
		}

		// check index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintUniq, ast.ConstraintKey, ast.ConstraintUniqKey, ast.ConstraintPrimaryKey:
				for _, k := range constraint.Keys {
					indexCols = append(indexCols, k.Column.Name.L)
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil {
				switch spec.Constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					for _, key := range spec.Constraint.Keys {
						indexCols = append(indexCols, key.Column.Name.L)
					}
				}
			}

			switch spec.Tp {
			case ast.AlterTableAddConstraint:
				if spec.Constraint == nil {
					continue
				}
				for _, key := range spec.Constraint.Keys {
					indexCols = append(indexCols, key.Column.Name.L)
				}
			case ast.AlterTableAddColumns, ast.AlterTableModifyColumn:
				checkNewColumns(spec.NewColumns)
			}
		}
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return indexCols, colsWithNotNullConstraint, err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				for _, option := range col.Options {
					switch option.Tp {
					case ast.ColumnOptionNotNull:
						colsWithNotNullConstraint[col.Name.Name.L] = struct{}{}
					}
				}
			}
		}
	case *ast.CreateIndexStmt:
		createTableStmt, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return indexCols, colsWithNotNullConstraint, err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				for _, option := range col.Options {
					switch option.Tp {
					case ast.ColumnOptionNotNull:
						colsWithNotNullConstraint[col.Name.Name.L] = struct{}{}
					}
				}
			}
		}
		for _, specification := range stmt.IndexPartSpecifications {
			indexCols = append(indexCols, specification.Column.Name.L)
		}
	default:
		return indexCols, colsWithNotNullConstraint, nil
	}
	return indexCols, colsWithNotNullConstraint, nil
}

func checkIndexNotNullConstraint(input *RuleHandlerInput) error {
	indexCols, colsWithNotNullConstraint, err := getIndexAndNotNullCols(input)
	if err != nil {
		return err
	}

	idxColsWithoutNotNull := []string{}
	indexCols = utils.RemoveDuplicate(indexCols)
	for _, k := range indexCols {
		if _, ok := colsWithNotNullConstraint[k]; !ok {
			idxColsWithoutNotNull = append(idxColsWithoutNotNull, k)
		}
	}
	if len(idxColsWithoutNotNull) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(idxColsWithoutNotNull, ","))
	}
	return nil
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
		} else if isExistRedundancyIndex(lastIndex, ind) {
			redundancy[ind.String()] = lastNormalIndex.String()
		} else {
			lastNormalIndex = ind
		}
		lastIndex = ind
	}

	return
}

func isExistRedundancyIndex(lastIndex index, ind index) bool {
	return utils.IsPrefixSubStrArray(lastIndex.Column, ind.Column)
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

func checkStringType(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				addResult(input.Res, input.Rule, DDLCheckColumnCharLength)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					addResult(input.Res, input.Rule, DDLCheckColumnCharLength)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkIfNotExist(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			addResult(input.Res, input.Rule, DDLCheckPKWithoutIfNotExists)
		}
	}
	return nil
}

func checkDDLTableSize(input *RuleHandlerInput) error {
	min := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	tables := []*ast.TableName{}
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		tables = append(tables, stmt.Table)
	case *ast.DropTableStmt:
		tables = append(tables, stmt.Tables...)
	default:
		return nil
	}

	beyond := []string{}
	for _, table := range tables {
		size, err := input.Ctx.GetTableSize(table)
		if err != nil {
			return err
		}
		if float64(min) < size {
			beyond = append(beyond, table.Name.String())
		}
	}

	if len(beyond) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(beyond, " , "), min)
	}
	return nil
}

func checkDMLTableSize(input *RuleHandlerInput) error {
	min := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	tables := []*ast.TableName{}
	switch stmt := input.Node.(type) {
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
		size, err := input.Ctx.GetTableSize(table)
		if err != nil {
			return err
		}
		if float64(min) < size {
			beyond = append(beyond, table.Name.String())
		}
	}

	if len(beyond) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(beyond, " , "), min)
	}
	return nil
}

func disableDropStmt(input *RuleHandlerInput) error {
	// specific check
	switch input.Node.(type) {
	case *ast.DropDatabaseStmt:
		addResult(input.Res, input.Rule, DDLDisableDropStatement)
	case *ast.DropTableStmt:
		addResult(input.Res, input.Rule, DDLDisableDropStatement)
	}
	return nil
}

func checkTableWithoutComment(input *RuleHandlerInput) error {
	var tableHasComment bool
	switch stmt := input.Node.(type) {
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
			addResult(input.Res, input.Rule, DDLCheckTableWithoutComment)
		}
	}
	return nil
}

func checkColumnWithoutComment(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
				addResult(input.Res, input.Rule, DDLCheckColumnWithoutComment)
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
					addResult(input.Res, input.Rule, DDLCheckColumnWithoutComment)
					return nil
				}
			}
		}
	}
	return nil
}

func checkIndexPrefix(input *RuleHandlerInput) error {
	indexesName := []string{}
	switch stmt := input.Node.(type) {
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
	prefix := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	for _, name := range indexesName {
		if !utils.HasPrefix(name, prefix, false) {
			addResult(input.Res, input.Rule, DDLCheckIndexPrefix, prefix)
			return nil
		}
	}
	return nil
}

func checkUniqIndexPrefix(input *RuleHandlerInput) error {
	_, indexes := getTableUniqIndex(input.Node)
	prefix := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	for index := range indexes {
		if !utils.HasPrefix(index, prefix, false) {
			addResult(input.Res, input.Rule, DDLCheckUniqueIndexPrefix, prefix)
			return nil
		}
	}
	return nil
}

func checkUniqIndex(input *RuleHandlerInput) error {
	tableName, indexes := getTableUniqIndex(input.Node)
	for index, indexedCols := range indexes {
		if !strings.EqualFold(index, fmt.Sprintf("IDX_UK_%v_%v", tableName, strings.Join(indexedCols, "_"))) {
			addResult(input.Res, input.Rule, DDLCheckUniqueIndex)
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

func checkColumnWithoutDefault(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
				addResult(input.Res, input.Rule, DDLCheckColumnWithoutDefault)
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
					addResult(input.Res, input.Rule, DDLCheckColumnWithoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnTimestampWithoutDefault(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
				addResult(input.Res, input.Rule, DDLCheckColumnTimestampWithoutDefault)
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
					addResult(input.Res, input.Rule, DDLCheckColumnTimestampWithoutDefault)
					return nil
				}
			}
		}
	}
	return nil
}

func checkColumnBlobNotNull(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
						addResult(input.Res, input.Rule, DDLCheckColumnBlobWithNotNull)
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
							addResult(input.Res, input.Rule, DDLCheckColumnBlobWithNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkColumnEnumNotice(input *RuleHandlerInput) error {
	return checkColumnShouldNotBeType(input.Rule, input.Res, input.Node, mysql.TypeEnum)
}

func checkColumnSetNotice(input *RuleHandlerInput) error {
	return checkColumnShouldNotBeType(input.Rule, input.Res, input.Node, mysql.TypeSet)
}

func checkColumnBlobNotice(input *RuleHandlerInput) error {
	return checkColumnShouldNotBeType(input.Rule, input.Res, input.Node, mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob)
}

func checkColumnShouldNotBeType(rule driverV2.Rule, res *driverV2.AuditResults, node ast.Node, colTypes ...byte) error {
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

func checkColumnBlobDefaultNull(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
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
						addResult(input.Res, input.Rule, DDLCheckColumnBlobDefaultIsNotNull)
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
							addResult(input.Res, input.Rule, DDLCheckColumnBlobDefaultIsNotNull)
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func checkDMLWithLimit(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit != nil {
			addResult(input.Res, input.Rule, DMLCheckWithLimit)
		}
	case *ast.DeleteStmt:
		if stmt.Limit != nil {
			addResult(input.Res, input.Rule, DMLCheckWithLimit)
		}
	}
	return nil
}

func checkSelectLimit(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:

		// 类似 select 1 和 select sleep(1) 这种不是真正查询的SQL, 没有检查limit的必要
		// select count() 没有limit的必要
		if stmt.From == nil || isSelectCount(stmt) {
			return nil
		}

		max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()

		if stmt.Limit == nil {
			addResult(input.Res, input.Rule, DMLCheckSelectLimit, max)
			return nil
		}

		value, ok := stmt.Limit.Count.(ast.ValueExpr)
		if !ok {
			return nil
		}
		limit, err := strconv.Atoi(fmt.Sprintf("%v", value.GetValue()))
		if err != nil {
			// 当limit的值为 ? 时此处会报错, 此时应当跳过检查
			//nolint:nilerr
			return nil
		}
		if limit > max {
			addResult(input.Res, input.Rule, DMLCheckSelectLimit, max)
			return nil
		}
		return nil
	default:
		return nil
	}
}

func checkDMLLimitExist(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		if stmt.Limit == nil {
			addResult(input.Res, input.Rule, DMLCheckLimitMustExist)
		}
	case *ast.DeleteStmt:
		if stmt.Limit == nil {
			addResult(input.Res, input.Rule, DMLCheckLimitMustExist)
		}
	}
	return nil
}

func checkDMLWithOrderBy(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		if stmt.Order != nil {
			addResult(input.Res, input.Rule, DMLCheckWithOrderBy)
		}
	case *ast.DeleteStmt:
		if stmt.Order != nil {
			addResult(input.Res, input.Rule, DMLCheckWithOrderBy)
		}
	}
	return nil
}

func checkDMLWithInsertColumnExist(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Columns) == 0 {
			addResult(input.Res, input.Rule, DMLCheckInsertColumnsExist)
		}
	}
	return nil
}

func checkDMLWithBatchInsertMaxLimits(input *RuleHandlerInput) error {
	max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	//value, err := strconv.Atoi(input.Rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing input.Rule[%v] value error: %v", input.Rule.Name, err)
	//}
	switch stmt := input.Node.(type) {
	case *ast.InsertStmt:
		if len(stmt.Lists) > max {
			addResult(input.Res, input.Rule, DMLCheckBatchInsertListsMax, max)
		}
	}
	return nil
}

func checkWhereExistFunc(input *RuleHandlerInput) error {
	tables := []*ast.TableName{}
	switch stmt := input.Node.(type) {
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
			checkExistFunc(input.Ctx, input.Rule, input.Res, tables, stmt.Where)
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
			checkExistFunc(input.Ctx, input.Rule, input.Res, tables, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			checkExistFunc(input.Ctx, input.Rule, input.Res, util.GetTables(stmt.TableRefs.TableRefs), stmt.Where)
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
			if checkExistFunc(input.Ctx, input.Rule, input.Res, tables, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}

func checkExistFunc(ctx *session.Context, rule driverV2.Rule, res *driverV2.AuditResults, tables []*ast.TableName, where ast.ExprNode) bool {
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

func checkWhereColumnImplicitConversion(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.From.TableRefs)
			// not select from table statement
			if len(tableSources) < 1 {
				break
			}
			checkWhereColumnImplicitConversionFunc(input.Ctx, input.Rule, input.Res, tableSources, stmt.Where)
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
			checkWhereColumnImplicitConversionFunc(input.Ctx, input.Rule, input.Res, tableSources, stmt.Where)
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
			checkWhereColumnImplicitConversionFunc(input.Ctx, input.Rule, input.Res, tableSources, stmt.Where)
		}
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			tableSources := util.GetTableSources(ss.From.TableRefs)
			if len(tableSources) < 1 {
				continue
			}
			if checkWhereColumnImplicitConversionFunc(input.Ctx, input.Rule, input.Res, tableSources, ss.Where) {
				break
			}
		}
	default:
		return nil
	}
	return nil
}

func checkWhereColumnImplicitConversionFunc(ctx *session.Context, rule driverV2.Rule, res *driverV2.AuditResults, tableSources []*ast.TableSource, where ast.ExprNode) bool {
	var hasImplicitConversionColumn bool
	if where == nil {
		return hasImplicitConversionColumn
	}

	util.ScanColumnValueFromExpr(where, func(cn *ast.ColumnName, values []*parserdriver.ValueExpr) bool {
		ts := getTableSourceByColumnName(ctx, tableSources, cn)
		if ts == nil {
			return false
		}
		// 暂时不处理子查询的情况
		switch source := ts.Source.(type) {
		case *ast.TableName:
			createTableStmt, exist, err := ctx.GetCreateTableStmt(source)
			if err != nil {
				return false
			}
			if !exist {
				return false
			}

			for _, col := range createTableStmt.Cols {
				if col.Name.Name.L != cn.Name.L {
					continue
				}

				// datetime, date, timestamp, time, year 类型的列不做检查
				// 因为这些类型的列不会发生隐式转换,mysql可以自动识别各种日期格式
				switch col.Tp.Tp {
				case mysql.TypeDatetime, mysql.TypeDate, mysql.TypeTimestamp, mysql.TypeDuration, mysql.TypeYear:
					continue
				}

				for _, v := range values {
					if !checkColumnTypeIsMatch(v, col.Tp.Tp) {
						addResult(res, rule, DMLCheckWhereExistImplicitConversion)
						hasImplicitConversionColumn = true
						return true
					}
				}
			}
		}
		return false
	})
	return hasImplicitConversionColumn
}

func getTableSourceByColumnName(ctx *session.Context, tableSources []*ast.TableSource, columnName *ast.ColumnName) *ast.TableSource {
	for _, ts := range tableSources {
		switch source := ts.Source.(type) {
		case *ast.TableName:
			if ts.AsName.L == columnName.Table.L {
				return ts
			}
			columnTableName := &ast.TableName{
				Schema: columnName.Schema,
				Name:   columnName.Table,
			}
			if ctx.GetSchemaName(source) == ctx.GetSchemaName(columnTableName) && source.Name.L == columnTableName.Name.L {
				return ts
			}
		}
	}
	return nil
}

func checkColumnTypeIsMatch(value *parserdriver.ValueExpr, ColumnType byte) bool {
	var columnTypeStr string
	switch ColumnType {
	case mysql.TypeVarchar, mysql.TypeString:
		columnTypeStr = "string"
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong, mysql.TypeDouble, mysql.TypeFloat, mysql.TypeNewDecimal:
		columnTypeStr = "int"
	default:
		columnTypeStr = "unknown"
	}

	var valueTypeStr string
	switch value.Datum.GetValue().(type) {
	case string:
		valueTypeStr = "string"
	case int, int8, int16, int32, int64, *tidbTypes.MyDecimal:
		valueTypeStr = "int"
	default:
		valueTypeStr = "unknown"
	}
	return valueTypeStr == columnTypeStr
}

func checkDMLSelectForUpdate(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.LockTp == ast.SelectLockForUpdate {
			addResult(input.Res, input.Rule, DMLCheckSelectForUpdate)
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

func checkCollationDatabase(input *RuleHandlerInput) error {
	var collationDatabase string
	var columnCollations []string
	var err error

	switch stmt := input.Node.(type) {
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
			collationDatabase, err = input.Ctx.GetCollationDatabase(tableName, "")
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
			collationDatabase, err = input.Ctx.GetCollationDatabase(nil, schemaName)
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
	expectCollation := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()

	// if collationDatabase empty, it means that we are not "create object"
	// and collation not change in "update object", so don't to check it.
	if collationDatabase != "" && !strings.EqualFold(collationDatabase, expectCollation) {
		addResult(input.Res, input.Rule, DDLCheckDatabaseCollation, expectCollation)
	}

	for _, cs := range columnCollations {
		if !strings.EqualFold(cs, expectCollation) {
			addResult(input.Res, input.Rule, DDLCheckDatabaseCollation, expectCollation)
			return nil
		}
	}
	return nil
}
func checkDecimalTypeColumn(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
				addResult(input.Res, input.Rule, DDLCheckDecimalTypeColumn)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp != nil && (col.Tp.Tp == mysql.TypeFloat || col.Tp.Tp == mysql.TypeDouble) {
					addResult(input.Res, input.Rule, DDLCheckDecimalTypeColumn)
					return nil
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func checkNeedlessFunc(input *RuleHandlerInput) error {
	funcArrStr := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	needlessFuncArr := strings.Split(funcArrStr, ",")
	funcExtractor := &functionVisitor{}
	input.Node.Accept(funcExtractor)
	functions := funcExtractor.functions
	for _, needlessFunc := range needlessFuncArr {
		for _, sqlFunc := range functions {
			needlessFunc = strings.ToLower(strings.TrimRight(needlessFunc, "()"))
			sqlFunc = strings.ToLower(sqlFunc)
			if needlessFunc == sqlFunc {
				addResult(input.Res, input.Rule, DMLCheckNeedlessFunc, funcArrStr)
				return nil
			}
		}
	}
	return nil
}

type functionVisitor struct {
	functions []string
}

func (v *functionVisitor) Enter(in ast.Node) (node ast.Node, skipChildren bool) {
	switch node := in.(type) {
	case *ast.FuncCallExpr:
		v.functions = append(v.functions, node.FnName.O)
	}
	return in, false
}

func (v *functionVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func checkDatabaseSuffix(input *RuleHandlerInput) error {
	databaseName := ""
	switch stmt := input.Node.(type) {
	case *ast.CreateDatabaseStmt:
		databaseName = stmt.Name
	case *ast.AlterDatabaseStmt:
		databaseName = stmt.Name
	default:
		return nil
	}
	suffix := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	if databaseName != "" && !utils.HasSuffix(databaseName, suffix, false) {
		addResult(input.Res, input.Rule, DDLCheckDatabaseSuffix, suffix)
		return nil
	}
	return nil
}

/*
规则：建议主键命名为"PK_表名"

	触发条件:
	1 是创建或变更表的DDL语句
	2 DDL语句中创建或修改了主键
	3 未对该主键命名或主键命名不遵循PK_表名的规范
*/
func checkPKIndexName(input *RuleHandlerInput) error {
	indexesName := ""
	tableName := ""
	var hasPrimaryKey bool
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				indexesName = constraint.Name
				tableName = stmt.Table.Name.String()
				hasPrimaryKey = true
				break
			}
		}
	case *ast.AlterTableStmt:
		tableName = strings.ToUpper(stmt.Table.Name.String())
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintPrimaryKey {
				indexesName = spec.Constraint.Name
				tableName = stmt.Table.Name.String()
				hasPrimaryKey = true
				break
			}
		}
	default:
		return nil
	}
	if hasPrimaryKey && !strings.EqualFold(indexesName, "PK_"+tableName) {
		addResult(input.Res, input.Rule, DDLCheckPKName)
		return nil
	}
	return nil
}

func checkTransactionIsolationLevel(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SetStmt:
		for _, variable := range stmt.Variables {
			if dry.StringListContains([]string{"tx_isolation", "tx_isolation_one_shot"}, variable.Name) {
				switch node := variable.Value.(type) {
				case *parserdriver.ValueExpr:
					if node.Datum.GetString() != ast.ReadCommitted {
						addResult(input.Res, input.Rule, DDLCheckTransactionIsolationLevel)
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

func checkTablePartition(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.PartitionNames != nil || spec.PartDefinitions != nil || spec.Partition != nil {
				addResult(input.Res, input.Rule, DDLCheckTablePartition)
				return nil
			}
		}
	case *ast.CreateTableStmt:
		if stmt.Partition != nil {
			addResult(input.Res, input.Rule, DDLCheckTablePartition)
			return nil
		}
	default:
		return nil
	}
	return nil
}
func checkNumberOfJoinTables(input *RuleHandlerInput) error {
	nums := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	//nums, err := strconv.Atoi(input.Rule.Value)
	//if err != nil {
	//	return fmt.Errorf("parsing input.Rule[%v] value error: %v", input.Rule.Name, err)
	//}
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil { //If from is null skip check. EX: select 1;select version
			return nil
		}
		if nums < util.GetNumberOfJoinTables(stmt.From.TableRefs) {
			addResult(input.Res, input.Rule, DMLCheckNumberOfJoinTables, nums)
		}
	default:
		return nil
	}
	return nil
}

func checkIsAfterUnionDistinct(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UnionStmt:
		for _, ss := range stmt.SelectList.Selects {
			if ss.IsAfterUnionDistinct {
				addResult(input.Res, input.Rule, DMLCheckIfAfterUnionDistinct)
				return nil
			}
		}
	default:
		return nil
	}

	return nil
}

func checkIsExistLimitOffset(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.Limit.Offset != nil {
			addResult(input.Res, input.Rule, DDLCheckIsExistLimitOffset)
		}
	default:
		return nil
	}
	return nil
}

func checkIndexOption(input *RuleHandlerInput) error {

	var tableName *ast.TableName
	indexColumns := make([]string, 0)
	switch stmt := input.Node.(type) {
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

	columnSelectivityMap, err := input.Ctx.GetSelectivityOfColumns(tableName, indexColumns)
	if err != nil {
		log.NewEntry().Errorf("get selectivity of columns failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}

	max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	var maxSelectivity float64 = -1
	for _, selectivity := range columnSelectivityMap {
		if selectivity > maxSelectivity {
			maxSelectivity = selectivity
		}
	}
	if maxSelectivity > 0 && maxSelectivity < float64(max) {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(indexColumns, ", "), max)
	}
	return nil
}

func isColumnUsingIndex(column string, constraints []*ast.Constraint) bool {
	for _, constraint := range constraints {
		for _, key := range constraint.Keys {
			if key.Column.Name.L == column {
				return true
			}
		}
	}
	return false
}

func checkWhereConditionUseIndex(ctx *session.Context, whereVisitor *util.WhereWithTableVisitor) bool {
	for _, whereExpr := range whereVisitor.WhereStmts {
		if whereExpr.WhereStmt == nil {
			return false
		}

		isUsingIndex := false

		if whereExpr.TableRef == nil {
			continue
		}

		tableNameCreateTableStmtMap := ctx.GetTableNameCreateTableStmtMap(whereExpr.TableRef)
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.ColumnNameExpr:
				tableName := x.Name.Table.L
				columnName := x.Name.Name.L
				// 代表单表查询并且没有使用表别名
				if tableName == "" {
					for _, createTableStmt := range tableNameCreateTableStmtMap {
						if isColumnUsingIndex(columnName, createTableStmt.Constraints) {
							isUsingIndex = true
						}
					}
				} else {
					createStmt, ok := tableNameCreateTableStmtMap[tableName]
					if ok && isColumnUsingIndex(columnName, createStmt.Constraints) {
						isUsingIndex = true
					}
				}
			}
			return false
		}, *whereExpr.WhereStmt)
		if !isUsingIndex {
			return false
		}
	}
	return true
}

func checkExplain(input *RuleHandlerInput) error {
	// sql from MyBatis XML file is not the executable sql. so can't do explain for it.
	// TODO(@wy) ignore explain when audit Mybatis file
	//if i.Task.SQLSource == driverV2.TaskSQLSourceFromMyBatisXMLFile {
	//	return nil
	//}
	switch input.Node.(type) {
	case *ast.SelectStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
	default:
		return nil
	}

	epRecords, err := input.Ctx.GetExecutionPlan(input.Node.Text())
	if err != nil {
		// TODO: check dml related table or database is created, if not exist, explain will executed failure.
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)

		// xml解析出来的sql获取执行计划会失败
		// 需要根据查询条件中的字段判断是否使用了索引
		if input.Rule.Name != DMLCheckExplainUsingIndex {
			return nil
		}
		// 验证where条件是否使用了索引字段
		wv := &util.WhereWithTableVisitor{}
		input.Node.Accept(wv)
		if !checkWhereConditionUseIndex(input.Ctx, wv) {
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}
		// 验证连表查询中连接字段是否使用索引
		isUsingIndex, err := judgeJoinFieldUseIndex(input)
		if err == nil && !isUsingIndex {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}

		return nil
	}
	for _, record := range epRecords {
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingFilesort) {
			addResult(input.Res, input.Rule, DMLCheckExplainExtraUsingFilesort)
		}
		if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingTemporary) {
			addResult(input.Res, input.Rule, DMLCheckExplainExtraUsingTemporary)
		}

		//defaultRule := RuleHandlerMap[DMLCheckExplainAccessTypeAll].Rule
		max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
		if record.Type == executor.ExplainRecordAccessTypeAll && record.Rows > int64(max) {
			addResult(input.Res, input.Rule, DMLCheckExplainAccessTypeAll, record.Rows)
		}

		if input.Rule.Name == DMLCheckExplainFullIndexScan &&
			record.Type == executor.ExplainRecordAccessTypeIndex {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}

		if input.Rule.Name == DMLCheckExplainExtraUsingIndexForSkipScan &&
			strings.Contains(record.Extra, executor.ExplainRecordExtraUsingIndexForSkipScan) {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		if input.Rule.Name == DMLCheckExplainUsingIndex && record.Key == "" {
			if strings.Contains(record.Extra, executor.ExplainRecordExtraUsingWhere) {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		}

	}
	return nil
}

func checkCreateView(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.CreateViewStmt:
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

var createTriggerReg1 = regexp.MustCompile(`(?i)create[\s]+trigger[\s]+[\S\s]+(before|after)+`)
var createTriggerReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+trigger[\s]+[\S\s]+(before|after)+`)

// CREATE
//
//	[DEFINER = user]
//	TRIGGER trigger_name
//	trigger_time trigger_event
//	ON tbl_name FOR EACH ROW
//	[trigger_order]
//	trigger_body
//
// ref:https://dev.mysql.com/doc/refman/8.0/en/create-trigger.html
//
// For now, we do character matching for CREATE TRIGGER Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateTrigger(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createTriggerReg1.MatchString(input.Node.Text()) ||
			createTriggerReg2.MatchString(input.Node.Text()) {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

var createFunctionReg1 = regexp.MustCompile(`(?i)create[\s]+function[\s]+[\S\s]+returns`)
var createFunctionReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+function[\s]+[\S\s]+returns`)

// CREATE
//
//	[DEFINER = user]
//	FUNCTION sp_name ([func_parameter[,...]])
//	RETURNS type
//	[characteristic ...] routine_body
//
// ref: https://dev.mysql.com/doc/refman/5.7/en/create-procedure.html
// For now, we do character matching for CREATE FUNCTION Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateFunction(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createFunctionReg1.MatchString(input.Node.Text()) ||
			createFunctionReg2.MatchString(input.Node.Text()) {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

var createProcedureReg1 = regexp.MustCompile(`(?i)create[\s]+procedure[\s]+[\S\s]+`)
var createProcedureReg2 = regexp.MustCompile(`(?i)create[\s]+[\s\S]+[\s]+procedure[\s]+[\S\s]+`)

// CREATE
//
//	[DEFINER = user]
//	PROCEDURE sp_name ([proc_parameter[,...]])
//	[characteristic ...] routine_body
//
// ref: https://dev.mysql.com/doc/refman/8.0/en/create-procedure.html
// For now, we do character matching for CREATE PROCEDURE Statement. Maybe we need
// more accurate match by adding such syntax support to parser.
func checkCreateProcedure(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.UnparsedStmt:
		if createProcedureReg1.MatchString(input.Node.Text()) ||
			createProcedureReg2.MatchString(input.Node.Text()) {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func checkAlias(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		repeats := []string{}
		fields := map[string]struct{}{}
		if stmt.From != nil {
			if source, ok := stmt.From.TableRefs.Left.(*ast.TableSource); ok {
				if tableName, ok := source.Source.(*ast.TableName); ok {
					fields[tableName.Name.L] = struct{}{}
				}

			}
		}
		for _, field := range stmt.Fields.Fields {
			if selectColumn, ok := field.Expr.(*ast.ColumnNameExpr); ok && selectColumn.Name.Name.L != "" {
				fields[selectColumn.Name.Name.L] = struct{}{}
			}
		}
		for _, field := range stmt.Fields.Fields {
			if _, ok := fields[field.AsName.L]; ok {
				repeats = append(repeats, field.AsName.String())
			}
		}
		if len(repeats) > 0 {
			addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(repeats, ","))
		}
		return nil
	default:
		return nil
	}
}

func hintUpdateTableCharsetWillNotUpdateFieldCharset(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, option := range spec.Options {
				if option.Tp == ast.TableOptionCharset {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func hintDropColumn(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropColumn {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func hintDropPrimaryKey(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DropIndexStmt:
		if strings.ToLower(stmt.IndexName) == "primary" {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropPrimaryKey {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func hintDropForeignKey(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		if len(stmt.Specs) > 0 {
			for _, spec := range stmt.Specs {
				if spec.Tp == ast.AlterTableDropForeignKey {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func notRecommendNotWildcardLike(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternLikeExpr:
				switch pattern := x.Pattern.(type) {
				case *parserdriver.ValueExpr:
					datum := pattern.Datum.GetString()
					if !strings.HasPrefix(datum, "%") && !strings.HasSuffix(datum, "%") {
						trigger = true
						return true
					}
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func hintInNullOnlyFalse(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch x := expr.(type) {
			case *ast.PatternInExpr:
				for _, exprNode := range x.List {
					switch pattern := exprNode.(type) {
					case *parserdriver.ValueExpr:
						if pattern.Datum.Kind() == tidbTypes.KindNull {
							trigger = true
							return true
						}
					}
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func notRecommendIn(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.PatternInExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func checkSpacesAroundTheString(input *RuleHandlerInput) error {
	visitor := &checkSpacesAroundTheStringVisitor{}
	input.Node.Accept(visitor)
	if visitor.HasPrefixOrSuffixSpace {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

type checkSpacesAroundTheStringVisitor struct {
	HasPrefixOrSuffixSpace bool
}

func (g *checkSpacesAroundTheStringVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if g.HasPrefixOrSuffixSpace {
		return n, false
	}

	if stmt, ok := n.(*parserdriver.ValueExpr); ok && stmt.Datum.Kind() == tidbTypes.KindString {
		if strings.HasPrefix(stmt.GetDatumString(), " ") || strings.HasSuffix(stmt.GetDatumString(), " ") {
			g.HasPrefixOrSuffixSpace = true
		}
	}

	return n, false
}

func (g *checkSpacesAroundTheStringVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func checkFullWidthQuotationMarks(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case ast.DDLNode:
		if strings.Contains(input.Node.Text(), "“") {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func notRecommendOrderByRand(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			if expr, ok := orderBy.Items[0].Expr.(*ast.FuncCallExpr); ok && expr.FnName.L == "rand" {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

func notRecommendGroupByConstant(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		groupBy := stmt.GroupBy
		if groupBy != nil {
			if _, ok := groupBy.Items[0].Expr.(*ast.PositionExpr); ok {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		}
		return nil
	default:
		return nil
	}
}

func checkSortDirection(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			isDesc := false
			for i, item := range orderBy.Items {
				if i == 0 {
					isDesc = item.Desc
				}
				if item.Desc != isDesc {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func hintGroupByRequiresConditions(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.GroupBy != nil && stmt.OrderBy == nil {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func notRecommendGroupByExpression(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		orderBy := stmt.OrderBy
		if orderBy != nil {
			for _, item := range orderBy.Items {
				if _, ok := item.Expr.(*ast.BinaryOperationExpr); ok {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func checkSQLLength(input *RuleHandlerInput) error {
	if len(input.Node.Text()) > input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func notRecommendHaving(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.Having != nil {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func hintUseTruncateInsteadOfDelete(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.DeleteStmt:
		if stmt.Where == nil {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func notRecommendUpdatePK(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		source, ok := stmt.TableRefs.TableRefs.Left.(*ast.TableSource)
		if !ok {
			return nil
		}
		t, ok := source.Source.(*ast.TableName)
		if !ok {
			return nil
		}
		createTable, exist, err := input.Ctx.GetCreateTableStmt(t)
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		primary := map[string]struct{}{}
		for _, col := range createTable.Constraints {
			if col.Tp == ast.ConstraintPrimaryKey {
				for _, key := range col.Keys {
					primary[key.Column.Name.L] = struct{}{}
				}
				break
			}
		}
		for _, assignment := range stmt.List {
			if _, ok := primary[assignment.Column.Name.L]; ok {
				addResult(input.Res, input.Rule, input.Rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func checkColumnQuantity(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		if len(stmt.Cols) > input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	default:
		return nil
	}
}

/*
recommendTableColumnCharsetSame

	触发条件：
	1 若DDL语句是建表语句：
		1.1 若列声明了字符集或排序规则。列的字符集或排序规则对应的字符集与表字符集不同，触发规则
	2 若DDL语句是修改表语句：
		2.1 若只修改列字符集。字符集与原表不同，触发规则。
		2.2 若修改了表的字符集，并且使用CONVERT，对比CONVERT的目标字符集，以及CONVERT后续修改列的语句
		2.3 !若修改了表的字符集，但没有使用CONVERT，不触发规则，暂不支持这种情形

	注意：若建表语句的字符集缺失，会按照mysql选择使用哪种字符集的逻辑获取字符集

	建表语句或列选择使用哪种字符集以及排序规则的逻辑如下：
	1 若符集以及排序都指定，则根据指定的字符集以及排序规则设定
	2 若未指定字符集，但指定了排序规则，字符集设定为排序规则关联的字符集
	3 若未指定排序规则，但指定了字符集，排序规则设定为字符集关联的排序规则
	4 若表的字符集和排序规则都不指定，二者将被设定为数据库的字符集及排序的默认值
	5 若列的字符集和排序规则都不知道，二者将被设定为数据表的字符集及排序的默认值
	6 若ALTER语句中使用CONVERT TO修改表的字符集，表中所有列的字符集都会被修改为目标字符集

	不支持：
	1 一条ALTER语句中反复修改同一列的字符集
	2 检查修改的列是否在表中
	3 修改表不使用CONVERT

参考文档1：
https://dev.mysql.com/doc/refman/8.0/en/charset-database.html

参考文档2：
https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
*/
func recommendTableColumnCharsetSame(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 获取建表语句中指定了字符集或排序的列
		columnWithCharset := getColumnWithCharset(stmt, input)
		if len(columnWithCharset) == 0 {
			return nil
		}
		// 获取建表语句中的字符集
		charset := getCharsetFromCreateTableStmt(input.Ctx, stmt)
		if charset.StrValue == "" {
			log.Logger().Warnf("skip rule:%s. reason: for sql %s, rule failed to obtain character set for comparison", input.Rule.Name, input.Node.Text())
			// 未能获取字符集 无法比较 返回
			return nil
		}
		for _, column := range columnWithCharset {
			if column.Tp.Charset != charset.StrValue {
				addResult(input.Res, input.Rule, input.Rule.Name)
				break
			}
		}
	case *ast.AlterTableStmt:
		var columnWithCharset []*ast.ColumnDef
		var newCharset *ast.TableOption
		var useConvert bool
		for _, spec := range stmt.Specs {
			// 修改的列
			for _, col := range spec.NewColumns {
				if col.Tp.Charset != "" {
					columnWithCharset = append(columnWithCharset, col)
				}
			}
			// 获取更改后的字符集以及更改的列
			charset, _ := getCharsetAndCollation(spec.Options)
			if charset.StrValue != "" {
				if charset.UintValue == ast.TableOptionCharsetWithConvertTo {
					// 使用CONVERT TO 则表的字符集会统一为该字符集 清空指定charset的列
					columnWithCharset = make([]*ast.ColumnDef, 0)
					useConvert = true
				}
				newCharset = charset
			}
		}

		if newCharset != nil {
			if useConvert {
				for _, column := range columnWithCharset {
					if column.Tp.Charset != newCharset.StrValue {
						addResult(input.Res, input.Rule, input.Rule.Name)
						break
					}
				}
				return nil
			}
			/*
				暂不支持修改表字符集但不使用CONVERT的情况
				若不使用CONVERT，仅修改表的字符集，不修改表中列的字符集
				需要判断最终的表字符集和列字符集是否一致，
				1. 获取原表各列字符集:
					SELECT column_name, character_set_name
					FROM information_schema.columns
					WHERE table_name = 'your_table_name';
				2. 根据SQL语句修改列的字符集到目标字符集
				3. 判断最终表字符集和最终列字符集是否一致
			*/
			log.Logger().Warnf("skip rule:%s. reason: for sql %s,alter the table character but not using CONVERT TO is currently not supported.", input.Rule.Name, input.Node.Text())
			return nil
		}
		if newCharset == nil {
			if len(columnWithCharset) == 0 {
				// 没有指定字符集的列时，列的字符集被设定为表字符集
				return nil
			}
			// 若未更改表的字符集，则获取原表字符集作为表字符集
			originTable, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
			if err != nil {
				log.Logger().Errorf("skip rule:%s. reason: for sql %s, an error occur when rule try to obtain the corresponding table,err:%v", input.Rule.Name, input.Node.Text(), err)
				return nil
			}
			if !exist {
				log.Logger().Warnf("skip rule:%s. reason: for sql %s,the corresponding table is not exist", input.Rule.Name, input.Node.Text())
				return nil
			}
			// 若没有修改表字符集
			charset := getCharsetFromCreateTableStmt(input.Ctx, originTable)
			if charset.StrValue == "" {
				// 未能获取字符集 无法比较 返回
				log.Logger().Warnf("skip rule:%s. reason: for sql %s, rule failed to obtain character set for comparison", input.Rule.Name, input.Node.Text())
				return nil
			}
			for _, column := range columnWithCharset {
				if column.Tp.Charset != charset.StrValue {
					addResult(input.Res, input.Rule, input.Rule.Name)
					break
				}
			}
		}

	default:
		return nil
	}
	return nil
}

func getColumnWithCharset(stmt *ast.CreateTableStmt, input *RuleHandlerInput) []*ast.ColumnDef {
	var columnWithCharset []*ast.ColumnDef
	for _, col := range stmt.Cols {

		if col.Tp.Charset != "" {
			columnWithCharset = append(columnWithCharset, col)
		} else if col.Tp.Collate != "" {
			col.Tp.Charset, _ = input.Ctx.GetSchemaCharacterByCollation(col.Tp.Collate)
			columnWithCharset = append(columnWithCharset, col)
		} else if len(col.Options) > 0 {
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionCollate {
					col.Tp.Charset, _ = input.Ctx.GetSchemaCharacterByCollation(option.StrValue)
					columnWithCharset = append(columnWithCharset, col)
				}
			}
		}
	}
	return columnWithCharset
}

func getCharsetAndCollation(options []*ast.TableOption) (*ast.TableOption, *ast.TableOption) {
	charset := &ast.TableOption{}
	collation := &ast.TableOption{}
	for _, option := range options {
		if option.Tp == ast.TableOptionCharset {
			charset = option
		}
		if option.Tp == ast.TableOptionCollate {
			collation = option
		}
	}
	return charset, collation
}

// 获取建表语句中的字符集
func getCharsetFromCreateTableStmt(ctx *session.Context, stmt *ast.CreateTableStmt) *ast.TableOption {
	charset, collation := getCharsetAndCollation(stmt.Options)
	if charset.StrValue == "" && collation.StrValue == "" {
		// 没有指定表字符集以及排序时，表的字符集和排序被设定为数据库默认字符集以及排序
		charset.StrValue, _ = ctx.GetSchemaCharacter(stmt.Table, "")
	}
	if charset.StrValue == "" && collation.StrValue != "" {
		// 指定了表排序但未指定表字符集时，表字符集被设定为排序对应的字符集
		charset.StrValue, _ = ctx.GetSchemaCharacterByCollation(collation.StrValue)
	}
	return charset
}

func checkColumnTypeInteger(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if (col.Tp.Tp == mysql.TypeLong && col.Tp.Flen != 10) || (col.Tp.Tp == mysql.TypeLonglong && col.Tp.Flen != 20) {
				addResult(input.Res, input.Rule, input.Rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func checkVarcharSize(input *RuleHandlerInput) error {
	maxVarcharLen := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if col.Tp == nil {
				continue
			}
			if col.Tp.Tp == mysql.TypeVarchar && col.Tp.Flen > maxVarcharLen {
				addResult(input.Res, input.Rule, input.Rule.Name, maxVarcharLen)
				break
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, column := range spec.NewColumns {
				if column.Tp == nil {
					continue
				}
				if column.Tp.Tp == mysql.TypeVarchar && column.Tp.Flen > maxVarcharLen {
					addResult(input.Res, input.Rule, input.Rule.Name, maxVarcharLen)
					break
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func containsOp(ops []opcode.Op, op opcode.Op) bool {
	for i := range ops {
		if op == ops[i] {
			return true
		}
	}
	return false
}

/*
应避免在 WHERE子句中使用函数或其他运算符

	触发条件：
	1 在WHERE子句中使用了函数，并且函数作用在至少一列上
	2 在WHERE二元操作符中使用运算符，包括：位运算符和算数运算符
	3 如果WHERE子句中使用了像sysdate()、now()这样的函数，不作用在任何列上，则不触发规则
*/
func notRecommendFuncInWhere(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch stmt := expr.(type) {
			case *ast.FuncCallExpr:
				visitor := util.ColumnNameVisitor{}
				stmt.Accept(&visitor)
				if len(visitor.ColumnNameList) > 0 {
					trigger = true
					return true
				}
			case *ast.BinaryOperationExpr:
				ops := []opcode.Op{
					opcode.LeftShift, opcode.RightShift, opcode.And, opcode.Or, opcode.BitNeg, opcode.Xor, // 位运算符
					opcode.Plus, opcode.Minus, opcode.Mul, opcode.Div, opcode.Mod, // 算术运算符
				}
				if containsOp(ops, stmt.Op) {
					trigger = true
					return true
				}
				return false
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func notRecommendSysdate(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.FuncCallExpr); ok && fu.FnName.L == "sysdate" {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.FuncCallExpr:
				if pattern.FnName.L == "sysdate" {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func hintSumFuncTips(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.AggregateFuncExpr); ok && strings.ToLower(fu.F) == "sum" {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.AggregateFuncExpr:
				if strings.ToLower(pattern.F) == "sum" {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func hintCountFuncWithCol(input *RuleHandlerInput) error {
	extractor := util.SelectStmtExtractor{}
	input.Node.Accept(&extractor)
	for _, stmt := range extractor.SelectStmts {
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.AggregateFuncExpr); ok && strings.ToLower(fu.F) == "count" {
				if fu.Distinct {
					continue
				}
				for _, arg := range fu.Args {
					if _, ok := arg.(*ast.ColumnNameExpr); ok {
						addResult(input.Res, input.Rule, input.Rule.Name)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func checkColumnQuantityInPK(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey && len(constraint.Keys) > input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int() {
				addResult(input.Res, input.Rule, input.Rule.Name)
				break
			}
		}
		return nil
	default:
		return nil
	}
}

func hintLimitMustBeCombinedWithOrderBy(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.OrderBy == nil {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
		return nil
	default:
		return nil
	}
}

func hintTruncateTips(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.TruncateTableStmt:
		addResult(input.Res, input.Rule, input.Rule.Name)
		return nil
	default:
		return nil
	}
}

func hintDeleteTips(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.TruncateTableStmt, *ast.DeleteStmt, *ast.DropTableStmt:
		addResult(input.Res, input.Rule, input.Rule.Name)
		return nil
	default:
		return nil
	}
}

func checkSQLInjectionFunc(input *RuleHandlerInput) error {
	funcs := []string{"sleep", "benchmark", "get_lock", "release_lock"}
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		for _, f := range stmt.Fields.Fields {
			if fu, ok := f.Expr.(*ast.FuncCallExpr); ok && inSlice(funcs, fu.FnName.L) {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	}
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.FuncCallExpr:
				if inSlice(funcs, pattern.FnName.L) {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func inSlice(ss []string, s string) bool {
	for _, s2 := range ss {
		if s2 == s {
			return true
		}
	}
	return false
}

func notRecommendSubquery(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch expr.(type) {
			case *ast.SubqueryExpr:
				trigger = true
				return true
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func checkNotEqualSymbol(input *RuleHandlerInput) error {
	if strings.Contains(input.Node.Text(), "!=") {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func checkSubqueryLimit(input *RuleHandlerInput) error {
	if where := getWhereExpr(input.Node); where != nil {
		trigger := false
		util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
			switch pattern := expr.(type) {
			case *ast.SubqueryExpr:
				if stmt, ok := pattern.Query.(*ast.SelectStmt); ok && stmt.Limit != nil && pattern.Query != nil {
					trigger = true
					return true
				}
			}
			return false
		}, where)
		if trigger {
			addResult(input.Res, input.Rule, input.Rule.Name)
		}
	}
	return nil
}

func checkAutoIncrement(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	default:
		return nil
	case *ast.CreateTableStmt:
		for _, option := range stmt.Options {
			if option.Tp == ast.TableOptionAutoIncrement && option.UintValue != 0 {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		}
		return nil
	}
}

func ddlNotAllowRenaming(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.RenameTableStmt:
		addResult(input.Res, input.Rule, input.Rule.Name)
		return nil
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Tp == ast.AlterTableChangeColumn ||
				spec.Tp == ast.AlterTableRenameTable ||
				spec.Tp == ast.AlterTableRenameColumn {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	}
	return nil
}

func checkLimitOffsetNum(input *RuleHandlerInput) error {
	maxOffset := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.Limit != nil && stmt.Limit.Offset != nil {
			offsetVal, ok := stmt.Limit.Offset.(*parserdriver.ValueExpr)
			if !ok {
				return nil
			}
			offset := offsetVal.Datum.GetInt64()
			if offset > int64(maxOffset) {
				addResult(input.Res, input.Rule, DMLCheckLimitOffsetNum, offset, maxOffset)
			}
		}
	default:
		return nil
	}
	return nil
}

func checkUpdateOrDeleteHasWhere(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		if stmt.Where == nil {
			addResult(input.Res, input.Rule, DMLCheckUpdateOrDeleteHasWhere)
		}
	case *ast.DeleteStmt:
		if stmt.Where == nil {
			addResult(input.Res, input.Rule, DMLCheckUpdateOrDeleteHasWhere)
		}
	default:
		return nil
	}
	return nil
}

func checkSortColumnLength(input *RuleHandlerInput) error {
	maxLength := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()

	type col struct {
		Table   *ast.TableName
		ColName string
	}
	checkColumns := []col{}

	buildCheckColumns := func(colName *ast.ColumnNameExpr, singleTableSource *ast.TableName) {
		var table *ast.TableName
		if singleTableSource == nil { // 这种情况是查询多表
			if colName.Name.Table.O == "" { // 查询多表的情况下order by的字段没有指定表名，简单处理，暂不对这个字段做校验。  todo 需要校验这种情况
				return
			}
			table = &ast.TableName{
				Schema: colName.Name.Schema,
				Name:   colName.Name.Table,
			}
		} else {
			table = &ast.TableName{
				Schema: singleTableSource.Schema,
				Name:   singleTableSource.Name,
			}
		}
		checkColumns = append(checkColumns, col{
			Table:   table,
			ColName: colName.Name.Name.L,
		})
	}

	gatherColFromOrderByClause := func(orderBy *ast.OrderByClause, singleTableSource *ast.TableName) {
		if orderBy != nil {
			for _, item := range orderBy.Items {
				if item == nil {
					continue
				}
				colName, ok := item.Expr.(*ast.ColumnNameExpr)
				if !ok {
					continue
				}
				buildCheckColumns(colName, singleTableSource)
			}
		}
	}

	gatherColFromSelectStmt := func(stmt *ast.SelectStmt, singleTableSource *ast.TableName) {
		gatherColFromOrderByClause(stmt.OrderBy, singleTableSource)
		if stmt.GroupBy != nil {
			for _, item := range stmt.GroupBy.Items {
				if item == nil {
					continue
				}
				colName, ok := item.Expr.(*ast.ColumnNameExpr)
				if !ok {
					continue
				}
				buildCheckColumns(colName, singleTableSource)
			}
		}
		if stmt.Distinct {
			if stmt.Fields != nil {
				for _, field := range stmt.Fields.Fields {
					if field == nil {
						continue
					}
					colName, ok := field.Expr.(*ast.ColumnNameExpr)
					if !ok {
						continue
					}
					buildCheckColumns(colName, singleTableSource)
				}
			}
		}
	}

	invalidCols := []string{}
	checkColLen := func(column col) error {
		table, exist, err := input.Ctx.GetCreateTableStmt(column.Table)
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		for _, def := range table.Cols {
			if def.Name.Name.L != column.ColName || def.Tp.Flen <= maxLength {
				continue
			}
			invalidCols = append(invalidCols, fmt.Sprintf("%v.%v", column.Table.Name.L, column.ColName))
		}
		return nil
	}

	var singleTable *ast.TableName
	// 简单处理表名：
	// 只在单表查询时通过from获取表名；
	// 多表查询时如果order by某个列没有指定表名，则不会检查这个列（这种情况应该不常见，暂时这样处理）
	// e.g. SELECT tb1.a,tb6.b FROM tb1,tb6 ORDER BY tb1.a,b  ->  字段b将不会被校验   todo 需要校验这种情况
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil || stmt.From.TableRefs == nil {
			return nil
		}
		// join子查询里的order by不做处理
		t, ok := stmt.From.TableRefs.Left.(*ast.TableSource)
		if ok && t != nil && stmt.From.TableRefs.Right == nil {
			temp, ok := t.Source.(*ast.TableName)
			if ok {
				singleTable = temp
			}
		}
		gatherColFromSelectStmt(stmt, singleTable)
	case *ast.UnionStmt:
		// join子查询里的order by不做处理
		// 因为union会对字段进行隐式排序，而order by的字段一定是union的字段，所以不需要额外对union语句的order by等函数的字段进行检查
		if stmt.SelectList == nil {
			return nil
		}
		for _, s := range stmt.SelectList.Selects {
			if s.From == nil || s.From.TableRefs == nil {
				continue
			}
			t, ok := s.From.TableRefs.Left.(*ast.TableSource)
			if ok && t != nil && s.From.TableRefs.Right == nil {
				temp, ok := t.Source.(*ast.TableName)
				if ok {
					singleTable = temp

					// 收集select的普通目标列
					if s.Fields != nil {
						for _, field := range s.Fields.Fields {
							if c, ok := field.Expr.(*ast.ColumnNameExpr); ok && c.Name != nil {
								checkColumns = append(checkColumns, col{
									Table:   temp,
									ColName: c.Name.Name.L,
								})
							}
						}
					}
				}
			}
			gatherColFromSelectStmt(s, singleTable) // 收集group by、distinct里的列
			gatherColFromOrderByClause(s.OrderBy, singleTable)
		}
	case *ast.DeleteStmt:
		if stmt.TableRefs == nil || stmt.TableRefs.TableRefs == nil {
			return nil
		}
		t, ok := stmt.TableRefs.TableRefs.Left.(*ast.TableSource)
		if ok && t != nil && stmt.TableRefs.TableRefs.Right == nil {
			temp, ok := t.Source.(*ast.TableName)
			if ok {
				singleTable = temp
			}
		}
		gatherColFromOrderByClause(stmt.Order, singleTable)
	case *ast.UpdateStmt:
		if stmt.TableRefs == nil || stmt.TableRefs.TableRefs == nil {
			return nil
		}
		t, ok := stmt.TableRefs.TableRefs.Left.(*ast.TableSource)
		if ok && t != nil && stmt.TableRefs.TableRefs.Right == nil {
			temp, ok := t.Source.(*ast.TableName)
			if ok {
				singleTable = temp
			}
		}
		gatherColFromOrderByClause(stmt.Order, singleTable)
	default:
		return nil
	}

	for _, column := range checkColumns {
		if err := checkColLen(column); err != nil {
			return err
		}
	}

	if len(invalidCols) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(invalidCols, ","))
	}

	return nil
}

func checkAffectedRows(input *RuleHandlerInput) error {

	switch input.Node.(type) {
	case *ast.UpdateStmt, *ast.DeleteStmt:
	default:
		return nil
	}

	affectCount, err := util.GetAffectedRowNum(
		context.TODO(), input.Node.Text(), input.Ctx.GetExecutor())
	if err != nil {
		log.NewEntry().Errorf("rule: %v; SQL: %v; get affected row number failed: %v", input.Rule.Name, input.Node.Text(), err)
		return nil
	}

	affectCountLimit := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	if affectCount > int64(affectCountLimit) {
		addResult(input.Res, input.Rule, input.Rule.Name, affectCount, affectCountLimit)
	}

	return nil
}

// NOTE: ParamMarkerExpr is actually "?".
// ref: https://docs.pingcap.com/zh/tidb/dev/expression-syntax#%E8%A1%A8%E8%BE%BE%E5%BC%8F%E8%AF%AD%E6%B3%95-expression-syntax
// ref: https://github.com/pingcap/tidb/blob/master/types/parser_driver/value_expr.go#L247
func checkPrepareStatementPlaceholders(input *RuleHandlerInput) error {

	placeholdersCount := 0
	placeholdersLimit := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:

		if whereStmt, ok := stmt.Where.(*ast.PatternInExpr); ok && stmt.Where != nil {
			for i := range whereStmt.List {
				item := whereStmt.List[i]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}

		if stmt.Fields != nil {
			for i := range stmt.Fields.Fields {
				item := stmt.Fields.Fields[i]
				if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
					placeholdersCount++
				}
			}
		}

		if stmt.GroupBy != nil {
			for i := range stmt.GroupBy.Items {
				item := stmt.GroupBy.Items[i]
				if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
					placeholdersCount++
				}
			}
		}

		if stmt.Having != nil && stmt.Having.Expr != nil {
			item := stmt.Having.Expr
			if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
				placeholdersCount++
			}
		}

		if stmt.OrderBy != nil {
			for i := range stmt.OrderBy.Items {
				item := stmt.OrderBy.Items[i]
				if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
					placeholdersCount++
				}
			}
		}

	case *ast.InsertStmt:
		for i := range stmt.Lists {
			for j := range stmt.Lists[i] {
				item := stmt.Lists[i][j]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}
		for i := range stmt.Setlist {
			if _, ok := stmt.Setlist[i].Expr.(*parserdriver.ParamMarkerExpr); ok && stmt.Setlist[i].Expr != nil {
				placeholdersCount++
			}
		}
		for i := range stmt.OnDuplicate {
			if _, ok := stmt.OnDuplicate[i].Expr.(*parserdriver.ParamMarkerExpr); ok && stmt.OnDuplicate[i].Expr != nil {
				placeholdersCount++
			}
		}

	case *ast.UpdateStmt:
		for i := range stmt.List {
			item := stmt.List[i]
			if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
				placeholdersCount++
			}
		}
		if whereStmt, ok := stmt.Where.(*ast.PatternInExpr); ok && stmt.Where != nil {
			for i := range whereStmt.List {
				item := whereStmt.List[i]
				if _, ok := item.(*parserdriver.ParamMarkerExpr); ok {
					placeholdersCount++
				}
			}
		}
		if stmt.Order != nil {
			for i := range stmt.Order.Items {
				item := stmt.Order.Items[i]
				if _, ok := item.Expr.(*parserdriver.ParamMarkerExpr); ok && item.Expr != nil {
					placeholdersCount++
				}
			}
		}

	}

	if placeholdersCount > placeholdersLimit {
		addResult(input.Res, input.Rule, input.Rule.Name, placeholdersCount, placeholdersLimit)
	}
	return nil
}

func checkAutoIncrementFieldNum(input *RuleHandlerInput) error {
	autoIncrementFieldNums := 0
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			for _, option := range col.Options {
				if option.Tp == ast.ColumnOptionAutoIncrement {
					autoIncrementFieldNums += 1
					break
				}
			}
		}
	default:
		return nil
	}

	if autoIncrementFieldNums > 1 {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}

	return nil
}

func getTableNameWithSchema(stmt *ast.TableName, c *session.Context) string {
	var tableWithSchema string

	if stmt.Schema.String() == "" {
		currentSchema := c.CurrentSchema()
		if currentSchema != "" {
			tableWithSchema = fmt.Sprintf("`%s`.`%s`", currentSchema, stmt.Name)
		} else {
			tableWithSchema = fmt.Sprintf("`%s`", stmt.Name)
		}
	} else {
		tableWithSchema = fmt.Sprintf("`%s`.`%s`", stmt.Schema, stmt.Name)
	}

	if c.IsLowerCaseTableName() {
		tableWithSchema = strings.ToLower(tableWithSchema)
	}

	return tableWithSchema
}

func checkSameTableJoinedMultipleTimes(input *RuleHandlerInput) error {
	var repeatTables []string

	if _, ok := input.Node.(ast.DMLNode); ok {
		selectVisitor := &util.SelectVisitor{}
		input.Node.Accept(selectVisitor)

		for _, selectNode := range selectVisitor.SelectList {
			tableJoinedNums := make(map[string]int)

			if selectNode.From != nil {
				tableSources := util.GetTableSources(selectNode.From.TableRefs)
				for _, tableSource := range tableSources {
					switch source := tableSource.Source.(type) {
					case *ast.TableName:
						tableName := getTableNameWithSchema(source, input.Ctx)
						tableJoinedNums[tableName] += 1
					}
				}

				for tableName, joinedNums := range tableJoinedNums {
					if joinedNums > 1 {
						repeatTables = append(repeatTables, tableName)
					}
				}
			}
		}
	}

	repeatTables = utils.RemoveDuplicate(repeatTables)
	if len(repeatTables) > 0 {
		tablesString := strings.Join(repeatTables, ",")
		addResult(input.Res, input.Rule, input.Rule.Name, tablesString)
	}

	return nil
}

func checkAllIndexNotNullConstraint(input *RuleHandlerInput) error {
	indexCols, colsWithNotNullConstraint, err := getIndexAndNotNullCols(input)
	if err != nil {
		return err
	}

	idxColsWithoutNotNull := []string{}
	indexCols = utils.RemoveDuplicate(indexCols)
	for _, k := range indexCols {
		if _, ok := colsWithNotNullConstraint[k]; !ok {
			idxColsWithoutNotNull = append(idxColsWithoutNotNull, k)
		}
	}
	if len(idxColsWithoutNotNull) == len(indexCols) && len(indexCols) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func checkInsertSelect(input *RuleHandlerInput) error {
	if stmt, ok := input.Node.(*ast.InsertStmt); ok {
		if stmt.Select != nil {
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}
	}
	return nil
}

func checkAggregateFunc(input *RuleHandlerInput) error {
	if _, ok := input.Node.(ast.DMLNode); !ok {
		return nil
	}
	selectVisitor := &util.SelectVisitor{}
	input.Node.Accept(selectVisitor)
	for _, selectNode := range selectVisitor.SelectList {
		if selectNode.Having != nil {
			isHavingUseFunc := false
			util.ScanWhereStmt(func(expr ast.ExprNode) bool {
				switch expr.(type) {
				case *ast.AggregateFuncExpr:
					isHavingUseFunc = true
					return true
				}
				return false
			}, selectNode.Having.Expr)

			if isHavingUseFunc {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
		for _, field := range selectNode.Fields.Fields {
			if _, ok := field.Expr.(*ast.AggregateFuncExpr); ok {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	}
	return nil
}

func checkColumnNotNull(input *RuleHandlerInput) error {
	notNullColumns := []string{}
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, newColumn := range spec.NewColumns {
				ok := util.IsAllInOptions(newColumn.Options, ast.ColumnOptionNotNull)
				if !ok {
					notNullColumns = append(notNullColumns, newColumn.Name.OrigColName())
				}
			}
		}
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			ok := util.IsAllInOptions(col.Options, ast.ColumnOptionNotNull)
			if !ok {
				notNullColumns = append(notNullColumns, col.Name.OrigColName())
			}
		}
	}
	if len(notNullColumns) > 0 {
		notNullColString := strings.Join(notNullColumns, ",")
		addResult(input.Res, input.Rule, input.Rule.Name, notNullColString)
	}
	return nil
}

func checkIndexSelectivity(input *RuleHandlerInput) error {
	if _, ok := input.Node.(*ast.SelectStmt); !ok {
		return nil
	}
	selectVisitor := &util.SelectVisitor{}
	input.Node.Accept(selectVisitor)
	explainRecords, err := input.Ctx.GetExecutionPlan(input.Node.Text())
	if err != nil {
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}
	for _, record := range explainRecords {
		indexes := strings.Split(record.Key, ",")
		recordTable := record.Table
		if len(indexes) == 0 || recordTable == "" {
			// 若执行计划没有使用索引 则跳过
			continue
		}
		for _, selectNode := range selectVisitor.SelectList {
			if selectNode.From == nil || selectNode.From.TableRefs == nil {
				continue
			}
			tables := util.GetTables(selectNode.From.TableRefs)
			for _, tableName := range tables {
				if tableName.Name.L != recordTable {
					// 只检查 使用索引对应的表
					continue
				}
				indexSelectivityMap, err := input.Ctx.GetSelectivityOfIndex(tableName, indexes)
				if err != nil {
					log.NewEntry().Errorf("get selectivity of index failed, sqle: %v, error: %v", input.Node.Text(), err)
					continue
				}
				max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
				for indexName, selectivity := range indexSelectivityMap {
					if selectivity > 0 && selectivity < float64(max) {
						addResult(input.Res, input.Rule, input.Rule.Name, indexName, max)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func checkTableRows(input *RuleHandlerInput) error {
	limitRowsString := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	limitRowsInt, err := strconv.Atoi(limitRowsString)
	if err != nil {
		return err
	}

	stmt, ok := input.Node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}

	exist, err := input.Ctx.IsTableExist(stmt.Table)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	rowsCount, err := input.Ctx.GetTableRowCount(stmt.Table)
	if err != nil {
		return err
	}

	if rowsCount > limitRowsInt*TenThousand {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func checkCompositeIndexSelectivity(input *RuleHandlerInput) error {
	indexSlices := [][]string{}
	table := &ast.TableName{}
	switch stmt := input.Node.(type) {
	case *ast.CreateIndexStmt:
		singleIndexSlice := []string{}
		if len(stmt.IndexPartSpecifications) == 1 {
			return nil
		}
		for _, indexPart := range stmt.IndexPartSpecifications {
			singleIndexSlice = append(singleIndexSlice, indexPart.Column.Name.O)
		}
		indexSlices = append(indexSlices, singleIndexSlice)
		table = stmt.Table
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint == nil {
				continue
			}
			if spec.Constraint.Tp != ast.ConstraintIndex && spec.Constraint.Tp != ast.ConstraintUniq {
				continue
			}
			singleIndexSlice := []string{}
			if len(spec.Constraint.Keys) == 1 {
				continue
			}
			for _, key := range spec.Constraint.Keys {
				singleIndexSlice = append(singleIndexSlice, key.Column.Name.O)
			}
			indexSlices = append(indexSlices, singleIndexSlice)
		}
		table = stmt.Table
	case *ast.CreateTableStmt:
		if stmt.Constraints == nil {
			return nil
		}
		for _, con := range stmt.Constraints {
			if con.Tp != ast.ConstraintIndex && con.Tp != ast.ConstraintUniq {
				continue
			}
			singleIndexSlice := []string{}
			if len(con.Keys) == 1 {
				continue
			}
			for _, key := range con.Keys {
				singleIndexSlice = append(singleIndexSlice, key.Column.Name.O)
			}
			indexSlices = append(indexSlices, singleIndexSlice)
		}
		table = stmt.Table
	}
	exist, err := input.Ctx.IsTableExist(table)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	l := log.NewEntry()
	noticeInfos := []string{}
	for _, singleIndexSlice := range indexSlices {
		var indexSelectValueSlice []struct {
			Index string
			Value float64
		}
		sortIndexes := make([]string, len(singleIndexSlice))
		colSelectivityMap, err := input.Ctx.GetSelectivityOfColumns(table, singleIndexSlice)
		if err != nil {
			l.Errorf("get columns selectivity error: %v", err)
			return nil
		}
		for _, indexColumn := range singleIndexSlice {
			selectivityValue, ok := colSelectivityMap[indexColumn]
			if !ok {
				l.Errorf("do not get column selectivity, column: %v", indexColumn)
				return nil
			}
			indexSelectValueSlice = append(indexSelectValueSlice, struct {
				Index string
				Value float64
			}{indexColumn, selectivityValue})
		}
		sort.Slice(indexSelectValueSlice, func(i, j int) bool {
			return indexSelectValueSlice[i].Value > indexSelectValueSlice[j].Value
		})
		for i, kv := range indexSelectValueSlice {
			sortIndexes[i] = kv.Index
		}
		for ind, indexColumn := range singleIndexSlice {
			if indexColumn != indexSelectValueSlice[ind].Index {
				noticeInfos = append(noticeInfos, fmt.Sprintf("(%s)可调整为(%s)", strings.Join(singleIndexSlice, "，"), strings.Join(sortIndexes, "，")))
				break
			}
		}
	}
	if len(noticeInfos) > 0 {
		addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(noticeInfos, "，"))
	}
	return nil
}

func judgeTextField(col *ast.ColumnDef) bool {
	if col.Tp.Tp == mysql.TypeBlob || col.Tp.Tp == mysql.TypeTinyBlob || col.Tp.Tp == mysql.TypeMediumBlob || col.Tp.Tp == mysql.TypeLongBlob {
		// mysql blob字段为二进制对象
		// https://dev.mysql.com/doc/refman/8.0/en/blob.html
		if col.Tp.Flag != mysql.BinaryFlag {
			return true
		}
	}
	return false
}

func checkText(input *RuleHandlerInput) error {
	textColumns := []string{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		var hasPk bool
		columnsWithoutPkAndText := make(map[string]struct{})
		for _, col := range stmt.Cols {
			isText := judgeTextField(col)
			if isText {
				textColumns = append(textColumns, col.Name.Name.O)
				continue
			}
			if util.IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
				hasPk = true
				continue
			}
			columnsWithoutPkAndText[col.Name.Name.O] = struct{}{}
		}
		for _, constraint := range stmt.Constraints {
			if constraint.Tp != ast.ConstraintPrimaryKey {
				continue
			}
			hasPk = true
			// 移除columnsWithoutPkAndText中主键的字段
			for _, key := range constraint.Keys {
				columnName := key.Column.Name.O
				delete(columnsWithoutPkAndText, columnName)
			}
		}
		if hasPk && len(textColumns) > 0 && len(columnsWithoutPkAndText) > 0 {
			addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(textColumns, "，"))
		}
	case *ast.AlterTableStmt:
		for _, col := range stmt.Specs {
			if col.Tp != ast.AlterTableAddColumns {
				continue
			}
			for _, newColumn := range col.NewColumns {
				isText := judgeTextField(newColumn)
				if isText {
					textColumns = append(textColumns, newColumn.Name.Name.O)
				}
			}
		}
		if len(textColumns) == 0 {
			return nil
		}
		originTable, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table)
		if err != nil {
			return err
		}
		if !exist {
			return nil
		}
		originPK, hasPk := util.GetPrimaryKey(originTable)
		if !hasPk {
			return nil
		}
		originTableAllColumns := []string{}
		for _, col := range originTable.Cols {
			originTableAllColumns = append(originTableAllColumns, col.Name.Name.L)
		}
		// 判断原表是否只存在主键
		if len(originPK) != len(originTableAllColumns) {
			addResult(input.Res, input.Rule, input.Rule.Name, strings.Join(textColumns, "，"))
		}
	}
	return nil
}

func checkSelectRows(input *RuleHandlerInput) error {
	if _, ok := input.Node.(*ast.SelectStmt); !ok {
		return nil
	}
	epRecords, err := input.Ctx.GetExecutionPlan(input.Node.Text())
	if err != nil {
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}

	var notUseIndex bool
	for _, record := range epRecords {
		if record.Type == executor.ExplainRecordAccessTypeIndex || record.Type == executor.ExplainRecordAccessTypeAll {
			notUseIndex = true
			break
		}
	}

	if !notUseIndex {
		return nil
	}
	affectCount, err := util.GetAffectedRowNum(context.TODO(), input.Node.Text(), input.Ctx.GetExecutor())
	if err != nil {
		return err
	}
	max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	if affectCount > int64(max)*int64(TenThousand) {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}

	return nil
}

func checkScanRows(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.SelectStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UpdateStmt:
	default:
		return nil
	}

	max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	epRecords, err := input.Ctx.GetExecutionPlan(input.Node.Text())
	if err != nil {
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}
	for _, record := range epRecords {
		if record.Rows > int64(max)*int64(TenThousand) {
			if record.Type == executor.ExplainRecordAccessTypeIndex || record.Type == executor.ExplainRecordAccessTypeAll {
				addResult(input.Res, input.Rule, input.Rule.Name)
				break
			}
		}
	}
	return nil
}

func mustMatchLeftMostPrefix(input *RuleHandlerInput) error {
	tables := []*ast.TableSource{}
	type colSets struct {
		AllCols    []string
		ColsWithEq []string
		ColsWithOr []string
	}
	tablesFromCondition := make(map[string]colSets)
	defaultTable := ""
	getTableName := func(tableName string) string {
		if tableName == "" {
			return defaultTable
		}
		return tableName
	}

	gatherColsWithOr := func(expr ast.ExprNode) {
		var col *ast.ColumnNameExpr
		var ok bool
		switch x := expr.(type) {
		case *ast.BinaryOperationExpr:
			col, ok = x.L.(*ast.ColumnNameExpr)
			if !ok {
				return
			}
		case *ast.PatternInExpr:
			col, ok = x.Expr.(*ast.ColumnNameExpr)
			if !ok {
				return
			}
		case *ast.PatternLikeExpr:
			col, ok = x.Expr.(*ast.ColumnNameExpr)
			if !ok {
				return
			}
		}

		if col == nil {
			return
		}
		tableName := getTableName(col.Name.Table.L)
		sets := tablesFromCondition[tableName]
		sets.AllCols = append(tablesFromCondition[tableName].AllCols, col.Name.Name.L)
		sets.ColsWithOr = append(tablesFromCondition[tableName].ColsWithOr, col.Name.Name.L)
		tablesFromCondition[tableName] = sets
	}
	var gatherColFromConditions func(expr ast.ExprNode) (skip bool)
	gatherColFromConditions = func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.BinaryOperationExpr:
			if x.Op == opcode.LogicOr || x.Op == opcode.LogicXor {
				gatherColsWithOr(x.L)
				gatherColsWithOr(x.R)
				return false
			}
			// select * from tb1 where a > 1
			// select * from tb1 where a = 1
			col, ok := x.L.(*ast.ColumnNameExpr)
			if !ok {
				return false
			}
			tableName := getTableName(col.Name.Table.L)
			if x.Op == opcode.EQ {
				sets := tablesFromCondition[tableName]
				sets.AllCols = append(tablesFromCondition[tableName].AllCols, col.Name.Name.L)
				sets.ColsWithEq = append(tablesFromCondition[tableName].ColsWithEq, col.Name.Name.L)
				tablesFromCondition[tableName] = sets
			} else {
				sets := tablesFromCondition[tableName]
				sets.AllCols = append(tablesFromCondition[tableName].AllCols, col.Name.Name.L)
				tablesFromCondition[tableName] = sets
			}
		case *ast.SubqueryExpr:
			if selectStmt, ok := x.Query.(*ast.SelectStmt); ok && selectStmt != nil {
				util.ScanWhereStmt(gatherColFromConditions, selectStmt.Where)
			}
		case *ast.PatternInExpr:
			//select * from tb1 where a IN(1,2)
			col, ok := x.Expr.(*ast.ColumnNameExpr)
			if !ok {
				return false
			}
			tableName := getTableName(col.Name.Table.L)
			sets := tablesFromCondition[tableName]
			sets.AllCols = append(tablesFromCondition[tableName].AllCols, col.Name.Name.L)
			tablesFromCondition[tableName] = sets
		case *ast.PatternLikeExpr:
			// select * from tb1 where a LIKE '%abc'
			// todo issue1783
			//col, ok := x.Expr.(*ast.ColumnNameExpr)
			//if !ok {
			//	return false
			//}
			//tableName := getTableName(col.Name.Table.L)
			//tablesFromCondition[tableName] = colSets{
			//	AllCols:    append(tablesFromCondition[tableName].AllCols, col.Name.Name.L),
			//	ColsWithEq: tablesFromCondition[tableName].ColsWithEq,
			//}
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil || stmt.Where == nil {
			return nil
		}
		if t, ok := stmt.From.TableRefs.Left.(*ast.TableSource); ok && t != nil {
			if name, ok := t.Source.(*ast.TableName); ok && name != nil {
				defaultTable = name.Name.L
			}
		}
		tables = util.GetTableSources(stmt.From.TableRefs)
		util.ScanWhereStmt(gatherColFromConditions, stmt.Where)
	case *ast.UpdateStmt:
		if stmt.Where == nil {
			return nil
		}
		if t, ok := stmt.TableRefs.TableRefs.Left.(*ast.TableSource); ok && t != nil {
			if name, ok := t.Source.(*ast.TableName); ok && name != nil {
				defaultTable = name.Name.L
			}
		}
		tables = util.GetTableSources(stmt.TableRefs.TableRefs)
		util.ScanWhereStmt(gatherColFromConditions, stmt.Where)
	case *ast.DeleteStmt:
		if stmt.Where == nil {
			return nil
		}
		if t, ok := stmt.TableRefs.TableRefs.Left.(*ast.TableSource); ok && t != nil {
			if name, ok := t.Source.(*ast.TableName); ok && name != nil {
				defaultTable = name.Name.L
			}
		}
		tables = util.GetTableSources(stmt.TableRefs.TableRefs)
		util.ScanWhereStmt(gatherColFromConditions, stmt.Where)
	case *ast.UnionStmt:
		for _, selectStmt := range stmt.SelectList.Selects {
			if selectStmt.Where == nil {
				continue
			}
			tables = util.GetTableSources(selectStmt.From.TableRefs)
			if t, ok := selectStmt.From.TableRefs.Left.(*ast.TableSource); ok && t != nil {
				if name, ok := t.Source.(*ast.TableName); ok && name != nil {
					defaultTable = name.Name.L
				}
			}
			util.ScanWhereStmt(gatherColFromConditions, selectStmt.Where)
		}
	default:
		return nil
	}

	isAllSubquery := true
	for _, table := range tables {
		if _, ok := table.Source.(*ast.TableName); ok {
			isAllSubquery = false
			break
		}
	}
	if isAllSubquery {
		return nil
	}

	for alias, cols := range tablesFromCondition {
		table, err := util.ConvertAliasToTable(alias, tables)
		if err != nil {
			log.NewEntry().Errorf("convert table alias failed, sqle: %v, error: %v", input.Node.Text(), err)
			return nil
		}
		createTable, exist, err := input.Ctx.GetCreateTableStmt(table)
		if err != nil {
			log.NewEntry().Errorf("get create table statement failed: %v", err)
			return nil
		}
		if !exist {
			log.NewEntry().Errorf("table [%v] doesn't exist", table.Name.O)
			return nil
		}

		if input.Rule.Name == DMLMustMatchLeftMostPrefix {
			if !isColumnMatchedALeftMostPrefix(cols.AllCols, cols.ColsWithEq, cols.ColsWithOr, createTable.Constraints) {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		} else if input.Rule.Name == DMLMustUseLeftMostPrefix {
			if !isColumnUseLeftMostPrefix(cols.AllCols, createTable.Constraints) {
				addResult(input.Res, input.Rule, input.Rule.Name)
			}
		}
	}
	return nil
}

func isColumnMatchedALeftMostPrefix(allCols []string, colsWithEQ, colsWithOr []string, constraints []*ast.Constraint) bool {
	multiConstraints := make([]*ast.Constraint, 0)
	for _, constraint := range constraints {
		if len(constraint.Keys) == 1 {
			continue
		}
		multiConstraints = append(multiConstraints, constraint)
	}
	walkConstraint := func(constraint *ast.Constraint) bool {
		for i, key := range constraint.Keys {
			for _, col := range allCols {
				if col != key.Column.Name.L {
					// 不是这个索引字段，跳过
					continue
				}
				if i == 0 && (!utils.StringsContains(colsWithEQ, col) || utils.StringsContains(colsWithOr, col)) {
					// 1.是最左字段 2.该字段没有使用等值查询或使用了or
					return false
				}
			}
		}
		return true
	}

	for _, constraint := range multiConstraints {
		if !walkConstraint(constraint) {
			return false
		}
	}
	return true
}

func checkSingleIndex(allCols []string, constraint *ast.Constraint) bool {
	singleIndexColumn := constraint.Keys[0].Column.Name.L
	for _, col := range allCols {
		if col == singleIndexColumn {
			return true
		}
	}
	return false
}

func isColumnUseLeftMostPrefix(allCols []string, constraints []*ast.Constraint) bool {
	multiConstraints := make([]*ast.Constraint, 0)
	for _, constraint := range constraints {
		if len(constraint.Keys) == 1 {
			if checkSingleIndex(allCols, constraint) {
				return true
			}
			continue
		}
		multiConstraints = append(multiConstraints, constraint)
	}
	walkConstraint := func(constraint *ast.Constraint) bool {
		isCurrentIndexUsed := false
		for i, key := range constraint.Keys {
			for _, col := range allCols {
				if col != key.Column.Name.L {
					// 不是这个索引的字段，跳过
					continue
				}
				isCurrentIndexUsed = true
				if i == 0 {
					return true
				}
			}
		}
		return !isCurrentIndexUsed
	}

	for _, constraint := range multiConstraints {
		if !walkConstraint(constraint) {
			return false
		}
	}
	return true
}

func checkJoinFieldUseIndex(input *RuleHandlerInput) error {
	isUsingIndex, err := judgeJoinFieldUseIndex(input)
	if err == nil && !isUsingIndex {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

/*
judgeJoinFieldUseIndex 判断Join语句中被驱动表中作为连接条件的列是否属于索引

	触发条件：
		A. CrossJoin，RightJoin和LeftJoin
			1. 分别判断ON USING WHERE中的连接条件是否有索引
	连接条件：等值条件两侧为不同表的列
	支持情况：
		支持：
			1. 多表JOIN
			2. 判断单列索引和复合索引
		不支持：
			1. 子查询中JOIN多表的判断
*/
func judgeJoinFieldUseIndex(input *RuleHandlerInput) (bool, error) {
	joinNode := getJoinNodeFromNode(input.Node)
	if doesNotJoinTables(joinNode) {
		// 如果SQL没有JOIN多表，则不需要审核
		return true, fmt.Errorf("sql have not join node")
	}
	tableNameCreateTableStmtMap := input.Ctx.GetTableNameCreateTableStmtMap(joinNode)
	tableIndexes := make(map[string][]*ast.Constraint, len(tableNameCreateTableStmtMap))
	for tableName, createTableStmt := range tableNameCreateTableStmtMap {
		tableIndexes[tableName] = createTableStmt.Constraints
	}

	if joinNodes, hasIndex := joinConditionInJoinNodeHasIndex(input.Ctx, joinNode, tableIndexes); joinNodes && !hasIndex {
		return false, nil
	}

	whereStmt := getWhereStmtFromNode(input.Node)
	if joinNodes, hasIndex := joinConditionInWhereStmtHasIndex(input.Ctx, joinNode, whereStmt, tableIndexes); joinNodes && !hasIndex {
		return false, nil
	}
	return true, nil
}

func joinConditionInWhereStmtHasIndex(ctx *session.Context, joinNode *ast.Join, whereStmt ast.ExprNode, tableIndex map[string][]*ast.Constraint) (joinTables, hasIndex bool) {
	if doesNotJoinTables(joinNode) {
		return false, false
	}
	if whereStmt == nil {
		return false, false
	}

	visitor := util.EqualConditionVisitor{}
	whereStmt.Accept(&visitor)
	tableColumnMap := make(tableColumnMap)
	for _, condition := range visitor.ConditionList {
		tableColumnMap.add(condition.Left.Table.L, condition.Left.Name.L)
		tableColumnMap.add(condition.Right.Table.L, condition.Right.Name.L)
	}
	for tableName, columnMap := range tableColumnMap {
		if constraints, ok := tableIndex[tableName]; ok {
			if !util.IsIndex(columnMap, constraints) {
				return true, false
			}
		}
	}

	return true, true
}

func joinConditionInJoinNodeHasIndex(ctx *session.Context, joinNode *ast.Join, tableIndex map[string] /*table name or alias name*/ []*ast.Constraint) (joinTables, hasIndex bool) {
	if doesNotJoinTables(joinNode) {
		return false, false
	}
	// 深度遍历左子树类型为ast.Join的节点
	if l, ok := joinNode.Left.(*ast.Join); ok {
		joinTables, hasIndex = joinConditionInJoinNodeHasIndex(ctx, l, tableIndex)
		if joinTables && !hasIndex {
			return joinTables, hasIndex
		}
	}

	tableColumnMap := make(tableColumnMap)

	if isJoinConditionInOnClause(joinNode) {
		visitor := util.EqualConditionVisitor{}
		joinNode.On.Accept(&visitor)
		for _, condition := range visitor.ConditionList {
			tableColumnMap.add(condition.Left.Table.L, condition.Left.Name.L)
			tableColumnMap.add(condition.Right.Table.L, condition.Right.Name.L)
		}
	}

	if isJoinConditionInUsingClause(joinNode) {
		leftTableSource, rightTableSource := getTableSourcesBesideJoin(joinNode)
		if leftTableSource != nil {
			tableName := getTableName(leftTableSource)
			for _, columnName := range joinNode.Using {
				tableColumnMap.add(tableName, columnName.Name.L)
			}
		}
		if rightTableSource != nil {
			tableName := getTableName(rightTableSource)
			for _, columnName := range joinNode.Using {
				tableColumnMap.add(tableName, columnName.Name.L)
			}
		}
	}

	for tableName, columnMap := range tableColumnMap {
		if constraints, ok := tableIndex[tableName]; ok {
			if !util.IsIndex(columnMap, constraints) {
				return true, false
			}
		}
	}
	return true, true
}

type tableColumnMap map[string] /*table name or alias name*/ map[string] /*column name*/ struct{}

func (m tableColumnMap) add(tableName, columnName string) {
	if m[tableName] == nil {
		m[tableName] = make(map[string]struct{})
	}
	m[tableName][columnName] = struct{}{}
}

/*
示例SQL:

	select * from a join b join c join d;
						   ↑
	如果*ast.Join节点是↑所指的join，拿到的是b和c两张表的tableSource

	select * from a join b join c join d;
								  ↑
	如果*ast.Join节点是↑所指的join，拿到的是c和d两张表的tableSource

	SQL: select * from a join(1) b join(2) c join(3) d;中join的抽象语法树如下:

		 join(3)
		 /     \
	   join(2)  d
	   /     \
	 join(1)  c
	 /    \
	a      b
*/
func getTableSourcesBesideJoin(joinNode *ast.Join) (left *ast.TableSource, right *ast.TableSource) {
	if tableSource, ok := joinNode.Right.(*ast.TableSource); ok {
		right = tableSource
	}
	if tableSource, ok := joinNode.Left.(*ast.TableSource); ok {
		left = tableSource
	}
	if join, ok := joinNode.Left.(*ast.Join); ok {
		if tableSource, ok := join.Right.(*ast.TableSource); ok {
			left = tableSource
		}
	}
	return
}

func getTableName(tableSource *ast.TableSource) string {
	if tableNameStmt, ok := tableSource.Source.(*ast.TableName); ok {
		if tableSource.AsName.L != "" {
			// 如果使用了别名，就应该用别名来引用列
			return tableSource.AsName.L
		}
		return tableNameStmt.Name.L
	}
	return ""
}

func getJoinNodeFromNode(node ast.Node) *ast.Join {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.From == nil {
			return nil
		}
		return stmt.From.TableRefs
	case *ast.UpdateStmt:
		if stmt.TableRefs == nil {
			return nil
		}
		return stmt.TableRefs.TableRefs
	case *ast.DeleteStmt:
		if stmt.TableRefs == nil {
			return nil
		}
		return stmt.TableRefs.TableRefs
	default:
		return nil
	}
}

func getWhereStmtFromNode(node ast.Node) ast.ExprNode {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if stmt.Where != nil {
			return stmt.Where
		}
	case *ast.UpdateStmt:
		if stmt.Where != nil {
			return stmt.Where
		}
	case *ast.DeleteStmt:
		if stmt.Where != nil {
			return stmt.Where
		}
	default:
		// 不检查JOIN不会存在的语句
		return nil
	}
	return nil
}

func getTableDefaultCharset(createTableStmt *ast.CreateTableStmt) string {
	characterSet := ""
	for _, op := range createTableStmt.Options {
		switch op.Tp {
		case ast.TableOptionCharset:
			characterSet = op.StrValue
		}
	}
	return characterSet
}

func getTableDefaultCollation(createTableStmt *ast.CreateTableStmt) string {
	collation := ""
	for _, op := range createTableStmt.Options {
		switch op.Tp {
		case ast.TableOptionCollate:
			collation = op.StrValue
		}
	}
	return collation
}

func getColumnCSCollation(columnName *ast.ColumnNameExpr, tableColumnCSCT map[string]map[string]columnCSCollation) columnCSCollation {
	var cSCollation columnCSCollation
	if columnCSCT, ok := tableColumnCSCT[columnName.Name.Table.L]; ok {
		cSCollation = columnCSCT[columnName.Name.Name.L]
	}

	return cSCollation
}

func getOnConditionLeftAndRightCSCollation(onCondition *ast.OnCondition, tableColumnCSCT map[string]map[string]columnCSCollation) (columnCSCollation, columnCSCollation) {
	var leftCSCollation, rightCSCollation columnCSCollation

	if binaryOperation, ok := onCondition.Expr.(*ast.BinaryOperationExpr); ok {
		if columnName, ok := binaryOperation.L.(*ast.ColumnNameExpr); ok {
			leftCSCollation = getColumnCSCollation(columnName, tableColumnCSCT)
		}

		if columnName, ok := binaryOperation.R.(*ast.ColumnNameExpr); ok {
			rightCSCollation = getColumnCSCollation(columnName, tableColumnCSCT)
		}
	}

	return leftCSCollation, rightCSCollation
}

type columnCSCollation struct {
	Charset   string
	Collation string
}

func checkJoinFieldCharacterSetAndCollation(input *RuleHandlerInput) error {
	tableNameCreateTableStmtMap, onConditions := getCreateTableAndOnCondition(input)
	if tableNameCreateTableStmtMap == nil || onConditions == nil {
		return nil
	}

	// 存储表名和列名的字符集排序规则映射关系, {tableName: {columnName: columnCSCollation}}
	tableColumnCSCT := make(map[string]map[string]columnCSCollation)
	for tableName, createTableStmt := range tableNameCreateTableStmtMap {
		tableDefaultCS := getTableDefaultCharset(createTableStmt)
		tableDefaultCollation := getTableDefaultCollation(createTableStmt)
		for _, col := range createTableStmt.Cols {
			cSCollation := columnCSCollation{}
			charset, hasCharset := getSingleColumnCSFromColumnsDef(col)
			if !hasCharset {
				continue
			}
			if charset == "" {
				charset = tableDefaultCS
			}
			cSCollation.Charset = charset
			for _, op := range col.Options {
				if op.Tp == ast.ColumnOptionCollate {
					cSCollation.Collation = op.StrValue
					break
				}
			}
			if cSCollation.Collation == "" {
				cSCollation.Collation = tableDefaultCollation
			}
			if tableColumnCSCT[tableName] == nil {
				tableColumnCSCT[tableName] = make(map[string]columnCSCollation)
			}
			tableColumnCSCT[tableName][col.Name.Name.L] = cSCollation
		}
	}

	for _, onCondition := range onConditions {
		leftCSCollation, righCSCollation := getOnConditionLeftAndRightCSCollation(onCondition, tableColumnCSCT)
		if leftCSCollation.Charset == "" || righCSCollation.Charset == "" {
			continue
		}
		if leftCSCollation.Charset != righCSCollation.Charset {
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}
		if leftCSCollation.Collation != righCSCollation.Collation {
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}

	}

	return nil
}

func checkSQLExplainLowestLevel(input *RuleHandlerInput) error {
	switch input.Node.(type) {
	case *ast.SelectStmt, *ast.DeleteStmt, *ast.UpdateStmt:
	default:
		return nil
	}

	levelStr := input.Rule.Params.GetParam(DefaultSingleParamKeyName).String()
	splitStr := strings.Split(levelStr, ",")
	levelMap := make(map[string]struct{})
	for _, s := range splitStr {
		s = strings.ToLower(strings.TrimSpace(s))
		levelMap[s] = struct{}{}
	}

	epRecords, err := input.Ctx.GetExecutionPlan(input.Node.Text())
	if err != nil {
		log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}
	for _, record := range epRecords {
		explainType := strings.ToLower(record.Type)
		if _, ok := levelMap[explainType]; !ok {
			addResult(input.Res, input.Rule, DMLSQLExplainLowestLevel, levelStr)
			return nil
		}
	}
	return nil
}

func avoidFullText(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintFulltext:
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
			switch spec.Constraint.Tp {
			case ast.ConstraintFulltext:
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	case *ast.CreateIndexStmt:
		switch stmt.KeyType {
		case ast.IndexKeyTypeFullText:
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}
	}
	return nil
}

func avoidGeometry(input *RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsGeometryColumn(col) {
				addResult(input.Res, input.Rule, input.Rule.Name)
				return nil
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint, ast.AlterTableAddColumns) {
			if spec.Constraint == nil {
				for _, newColumn := range spec.NewColumns {
					if util.IsGeometryColumn(newColumn) {
						addResult(input.Res, input.Rule, input.Rule.Name)
						return nil
					}
				}
			} else {
				switch spec.Constraint.Tp {
				case ast.ConstraintSpatial:
					addResult(input.Res, input.Rule, input.Rule.Name)
					return nil
				}
			}
		}
	case *ast.CreateIndexStmt:
		switch stmt.KeyType {
		case ast.IndexKeyTypeSpatial:
			addResult(input.Res, input.Rule, input.Rule.Name)
			return nil
		}
	}
	return nil
}

func avoidWhereEqualNull(input *RuleHandlerInput) error {
	dmlNode, ok := input.Node.(ast.DMLNode)
	if !ok {
		return nil
	}

	whereVisitor := &util.WhereVisitor{}
	dmlNode.Accept(whereVisitor)
	for _, whereExpr := range whereVisitor.WhereList {
		util.ScanWhereStmt(func(expr ast.ExprNode) bool {
			switch stmt := expr.(type) {
			case *ast.BinaryOperationExpr:
				for _, binExpr := range []ast.ExprNode{stmt.L, stmt.R} {
					value, ok := binExpr.(*parserdriver.ValueExpr)
					if ok {
						if value.Type.Tp == mysql.TypeNull {
							switch stmt.Op {
							case opcode.GE, opcode.LE, opcode.EQ, opcode.NE, opcode.LT, opcode.GT:
								addResult(input.Res, input.Rule, input.Rule.Name)
								return true
							}
						}
					}
				}
			}
			return false
		}, whereExpr)
	}
	return nil
}

func avoidEvent(input *RuleHandlerInput) error {
	if util.IsEventSQL(input.Node.Text()) {
		addResult(input.Res, input.Rule, input.Rule.Name)
	}
	return nil
}

func getCharLengthFromColumn(col *ast.ColumnDef) int {
	charLength := 0
	switch col.Tp.Tp {
	case mysql.TypeString:
		// charset为binary，字段类型为binary
		if col.Tp.Charset != "binary" {
			charLength = col.Tp.Flen
		}
	case mysql.TypeVarchar:
		// charset为binary，字段类型为varbinary
		if col.Tp.Charset != "binary" {
			charLength = col.Tp.Flen
		}
	}
	return charLength
}

func checkCharLength(input *RuleHandlerInput) error {
	max := input.Rule.Params.GetParam(DefaultSingleParamKeyName).Int()
	charLength := 0
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			charLength += getCharLengthFromColumn(col)
		}
	case *ast.AlterTableStmt:
		if originTable, exist, err := input.Ctx.GetCreateTableStmt(stmt.Table); err == nil && exist {
			modifyColMap := make(map[string]struct{})
			for _, col := range stmt.Specs {
				if col.Tp == ast.AlterTableAddColumns {
					for _, newColumn := range col.NewColumns {
						charLength += getCharLengthFromColumn(newColumn)
					}
				} else if col.Tp == ast.AlterTableModifyColumn {
					for _, modifyCol := range col.NewColumns {
						modifyColMap[modifyCol.Name.Name.O] = struct{}{}
						charLength += getCharLengthFromColumn(modifyCol)
					}
				}
			}
			if charLength > 0 {
				for _, col := range originTable.Cols {
					// 获取建表语句char总和时，排除modify字段
					if _, ok := modifyColMap[col.Name.Name.O]; ok {
						continue
					}
					charLength += getCharLengthFromColumn(col)
				}
			}
		}
	}
	if charLength > max {
		addResult(input.Res, input.Rule, input.Rule.Name, max)
	}
	return nil
}
