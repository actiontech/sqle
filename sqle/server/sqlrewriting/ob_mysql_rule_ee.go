//go:build enterprise
// +build enterprise

package sqlrewriting

// 现有规则名 到 重写功能使用的 新规则 ID 的映射
// 只映射了部分规则，用于demo版本
// TODO: 当SQLE使用 新规则ID后，可以去除这个映射
/*
SQLE现有OBMySQL规则(基于sqle-ob-mysql-plugin仓库 99972e712b2a165a43155ecade6b3d0fa33ab3b0统计): 150
OB MySQL映射规则: 114
OB MySQL未映射规则: 36
// 以下是未映射规则
{
	Name:   "dml_enable_explain_pre_check",
	CH:     "使用EXPLAIN加强预检查能力",
},
{
	Name:   "sql_is_executed",
	CH:     "停用上线审核模式",
},
{
	Name:   "dml_rollback_max_rows",
	CH:     "在 DML 语句中预计影响行数超过指定值则不回滚",
},
{
	Name:   "ddl_check_index_not_null_constraint",
	CH:     "索引字段需要有非空约束",
},
{
	Name:   "ddl_check_index_column_with_blob",
	CH:     "禁止将BLOB类型的列加入索引",
},
{
	Name:   "ddl_check_column_without_default",
	CH:     "除了自增列及大字段列之外，每个列都必须添加默认值",
},
{
	Name:   "dml_check_with_limit",
	CH:     "DELETE/UPDATE 语句不能有LIMIT条件",
},
{
	Name:   "ddl_check_indexes_exist_before_creat_constraints",
	CH:     "对字段创建约束前，建议先创建索引",
},
{
	Name:   "ddl_check_pk_name",
	CH:     "建议主键命名为\"PK_表名\"",
},
{
	Name:   "ddl_check_column_set_notice",
	CH:     "不建议使用 SET 类型",
},
{
	Name:   "dml_not_recommend_not_wildcard_like",
	CH:     "不建议使用没有通配符的 LIKE 查询",
},
{
	Name:   "dml_not_recommend_in",
	CH:     "不建议使用IN",
},
{
	Name:   "dml_check_spaces_around_the_string",
	CH:     "引号中的字符串开头或结尾不建议包含空格",
},
{
	Name:   "dml_not_recommend_sysdate",
	CH:     "不建议使用 SYSDATE() 函数",
},
{
	Name:   "dml_hint_count_func_with_col",
	CH:     "避免使用 COUNT(COL)",
},
{
	Name:   "ddl_check_table_rows",
	CH:     "表行数超过阈值，建议对表进行拆分",
},
{
	Name:   "ddl_check_composite_index_distinction",
	CH:     "建议在组合索引中将区分度高的字段靠前放",
},
{
	Name:   "ddl_avoid_text",
	CH:     "使用TEXT 类型的字段建议和原表进行分拆，与原表主键单独组成另外一个表进行存放",
},
{
	Name:   "dml_check_select_rows",
	CH:     "查询数据量超过阈值，筛选条件必须带上主键或者索引",
},
{
	Name:   "dml_check_scan_rows",
	CH:     "扫描行数超过阈值，筛选条件必须带上主键或者索引",
},
{
	Name:   "dml_must_use_left_most_prefix",
	CH:     "使用联合索引时，必须使用联合索引的首字段",
},
{
	Name:   "dml_must_match_left_most_prefix",
	CH:     "禁止对联合索引左侧字段进行IN 、OR等非等值查询",
},
{
	Name:   "dml_check_join_field_use_index",
	CH:     "JOIN字段必须包含索引",
},
{
	Name:   "dml_check_join_field_character_set_Collation",
	CH:     "连接表字段的字符集和排序规则必须一致",
},
{
	Name:   "dml_sql_explain_lowest_level",
	CH:     "SQL执行计划中type字段建议满足规定的级别",
},
{
	Name:   "ddl_avoid_full_text",
	CH:     "禁止使用全文索引",
},
{
	Name:   "ddl_avoid_geometry",
	CH:     "禁止使用空间字段和空间索引",
},
{
	Name:   "dml_avoid_where_equal_null",
	CH:     "WHERE子句中禁止将NULL值与其他字段或值进行比较运算",
},
{
	Name:   "ddl_avoid_event",
	CH:     "禁止使用event",
},
{
	Name:   "ddl_check_char_length",
	CH:     "禁止char, varchar类型字段字符长度总和超过阈值",
},
{
	Name:   "dml_avoid_count_column_name",
	CH:     "不推荐使用 count(列名) 来替代 count(*)",
},
{
	Name:   "ddl_check_database_name_length",
	CH:     "库名长度不能超过指定字节",
},
{
	Name:   "ddl_should_not_add_primary_key_after_creating_table",
	CH:     "建表后不允许添加主键",
},
{
	Name:   "ddl_should_not_use_json_type",
	CH:     "不建议使用JSON类型",
},
{
	Name:   "ddl_should_not_create_database_naming_with_test",
	CH:     "禁止创建名为 test 或以 test 开头的库",
},
{
	Name:   "should_not_use_database_oceanbase_and_test",
	CH:     "禁止使用或操作 oceanbase 库和 test 库",
},
*/
// 以下是映射规则
var OBMySQLRuleIdConvert = []ruleIdConvert{
	{
		Name:   "ddl_check_index_too_many",
		CH:     "单字段上的索引数量不建议超过阈值",
		RuleId: "00043",
	},
	{
		Name:   "ddl_check_redundant_index",
		CH:     "不建议创建冗余索引",
		RuleId: "00055",
	},
	{
		Name:   "ddl_check_table_without_if_not_exists",
		CH:     "新建表建议加入 IF NOT EXISTS，保证重复执行不报错",
		RuleId: "00061",
	},
	{
		Name:   "ddl_check_object_name_length",
		CH:     "表名、列名、索引名的长度不建议超过阈值",
		RuleId: "00047",
	},
	{
		Name:   "ddl_check_object_name_is_upper_and_lower_letter_mixed",
		CH:     "数据库对象命名不建议大小写字母混合",
		RuleId: "00046",
	},
	{
		Name:   "ddl_check_pk_not_exist",
		CH:     "表必须有主键",
		RuleId: "00008",
	},
	{
		Name:   "ddl_check_pk_without_auto_increment",
		CH:     "主键建议使用自增",
		RuleId: "00052",
	},
	{
		Name:   "ddl_check_pk_without_bigint_unsigned",
		CH:     "主键建议使用 BIGINT 无符号类型，即 BIGINT UNSIGNED",
		RuleId: "00054",
	},
	{
		Name:   "dml_check_join_field_type",
		CH:     "建议JOIN字段类型保持一致",
		RuleId: "00006",
	},
	{
		Name:   "dml_check_join_has_on",
		CH:     "建议连接操作指定连接条件",
		RuleId: "00091",
	},
	{
		Name:   "ddl_check_column_char_length",
		CH:     "CHAR长度大于20时，必须使用VARCHAR类型",
		RuleId: "00018",
	},
	{
		Name:   "ddl_check_field_not_null_must_contain_default_value",
		CH:     "建议字段约束为NOT NULL时带默认值",
		RuleId: "00034",
	},
	{
		Name:   "ddl_disable_fk",
		CH:     "禁止使用外键",
		RuleId: "00067",
	},
	{
		Name:   "ddl_check_create_time_column",
		CH:     "建议建表DDL包含创建时间字段且默认值为CURRENT_TIMESTAMP",
		RuleId: "00219",
	},
	{
		Name:   "ddl_check_index_count",
		CH:     "索引个数建议不超过阈值",
		RuleId: "00037",
	},
	{
		Name:   "ddl_check_update_time_column",
		CH:     "建表DDL需要包含更新时间字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
		RuleId: "00033",
	},
	{
		Name:   "ddl_check_composite_index_max",
		CH:     "复合索引的列数量不建议超过阈值",
		RuleId: "00005",
	},
	{
		Name:   "ddl_check_object_name_using_keyword",
		CH:     "数据库对象命名禁止使用保留字",
		RuleId: "00049",
	},
	{
		Name:   "ddl_check_object_name_using_cn",
		CH:     "数据库对象命名只能使用英文、下划线或数字，首字母必须是英文",
		RuleId: "00048",
	},
	{
		Name:   "ddl_check_table_character_set",
		CH:     "建议使用指定数据库字符集",
		RuleId: "00015",
	},
	{
		Name:   "all_check_where_is_invalid",
		CH:     "禁止使用没有WHERE条件或者WHERE条件恒为TRUE的SQL",
		RuleId: "00001",
	},
	{
		Name:   "ddl_check_alter_table_need_merge",
		CH:     "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
		RuleId: "00011",
	},
	{
		Name:   "dml_disable_select_all_column",
		CH:     "不建议使用SELECT *",
		RuleId: "00053",
	},
	{
		Name:   "ddl_disable_drop_statement",
		CH:     "禁止除索引外的DROP操作",
		RuleId: "00066",
	},
	{
		Name:   "ddl_check_table_without_comment",
		CH:     "表建议添加注释",
		RuleId: "00060",
	},
	{
		Name:   "ddl_check_column_without_comment",
		CH:     "列建议添加注释",
		RuleId: "00027",
	},
	{
		Name:   "ddl_check_index_prefix",
		CH:     "建议普通索引使用固定前缀",
		RuleId: "00040",
	},
	{
		Name:   "ddl_check_unique_index_prefix",
		CH:     "建议UNIQUE索引使用固定前缀",
		RuleId: "00041",
	},
	{
		Name:   "ddl_check_unique_index",
		CH:     "建议UNIQUE索引名使用 IDX_UK_表名_字段名",
		RuleId: "00063",
	},
	{
		Name:   "ddl_check_column_timestamp_without_default",
		CH:     "TIMESTAMP 类型的列必须添加默认值",
		RuleId: "00025",
	},
	{
		Name:   "ddl_check_column_blob_with_not_null",
		CH:     "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		RuleId: "00016",
	},
	{
		Name:   "ddl_check_column_blob_default_is_not_null",
		CH:     "BLOB 和 TEXT 类型的字段默认值只能为NULL",
		RuleId: "00016",
	},
	{
		Name:   "ddl_check_auto_increment_field_num",
		CH:     "建表时，自增字段只能设置一个",
		RuleId: "00007",
	},
	{
		Name:   "ddl_check_all_index_not_null_constraint",
		CH:     "建议为至少一个索引添加非空约束",
		RuleId: "00003",
	},
	{
		Name:   "dml_check_select_limit",
		CH:     "SELECT 语句需要带LIMIT",
		RuleId: "00100",
	},
	{
		Name:   "dml_check_with_order_by",
		CH:     "DELETE/UPDATE 语句不能有ORDER BY",
		RuleId: "00102",
	},
	{
		Name:   "dml_check_select_with_order_by",
		CH:     "SELECT 语句不能有ORDER BY",
		RuleId: "00101",
	},
	{
		Name:   "dml_check_insert_columns_exist",
		CH:     "INSERT 语句需要指定COLUMN",
		RuleId: "00088",
	},
	{
		Name:   "dml_check_batch_insert_lists_max",
		CH:     "单条INSERT语句，建议批量插入不超过阈值",
		RuleId: "00080",
	},
	{
		Name:   "dml_check_in_query_limit",
		CH:     "WHERE条件内IN语句中的参数个数不能超过阈值",
		RuleId: "00087",
	},
	{
		Name:   "ddl_check_pk_prohibit_auto_increment",
		CH:     "不建议主键使用自增",
		RuleId: "00051",
	},
	{
		Name:   "dml_check_where_exist_func",
		CH:     "避免对条件字段使用函数操作",
		RuleId: "00111",
	},
	{
		Name:   "dml_check_where_exist_not",
		CH:     "不建议对条件字段使用负向查询",
		RuleId: "00113",
	},
	{
		Name:   "dml_check_where_exist_null",
		CH:     "不建议对条件字段使用 NULL 值判断",
		RuleId: "00001",
	},
	{
		Name:   "dml_check_where_exist_implicit_conversion",
		CH:     "不建议在WHERE条件中使用与过滤字段不一致的数据类型",
		RuleId: "00112",
	},
	{
		Name:   "dml_check_limit_must_exist",
		CH:     "建议DELETE/UPDATE 语句带有LIMIT条件",
		RuleId: "00092",
	},
	{
		Name:   "dml_check_where_exist_scalar_sub_queries",
		CH:     "不建议使用标量子查询",
		RuleId: "00115",
	},
	{
		Name:   "dml_check_select_for_update",
		CH:     "不建议使用SELECT FOR UPDATE",
		RuleId: "00099",
	},
	{
		Name:   "ddl_check_collation_database",
		CH:     "建议使用规定的数据库排序规则",
		RuleId: "00015",
	},
	{
		Name:   "ddl_check_decimal_type_column",
		CH:     "精确浮点数建议使用DECIMAL",
		RuleId: "00013",
	},
	{
		Name:   "ddl_check_bigint_instead_of_decimal",
		CH:     "建议用BIGINT类型代替DECIMAL",
		RuleId: "00012",
	},
	{
		Name:   "dml_check_sub_query_depth",
		CH:     "子查询嵌套层数不建议超过阈值",
		RuleId: "00108",
	},
	{
		Name:   "dml_check_needless_func",
		CH:     "避免使用不必要的内置函数",
		RuleId: "00094",
	},
	{
		Name:   "ddl_check_database_suffix",
		CH:     "建议数据库名称使用固定后缀结尾",
		RuleId: "00032",
	},
	{
		Name:   "ddl_check_transaction_isolation_level",
		CH:     "事物隔离级别建议设置成RC",
		RuleId: "00062",
	},
	{
		Name:   "dml_check_fuzzy_search",
		CH:     "禁止使用全模糊搜索或左模糊搜索",
		RuleId: "00086",
	},
	{
		Name:   "ddl_check_table_partition",
		CH:     "不建议使用分区表相关功能",
		RuleId: "00058",
	},
	{
		Name:   "dml_check_number_of_join_tables",
		CH:     "使用JOIN连接表查询建议不超过阈值",
		RuleId: "00096",
	},
	{
		Name:   "dml_check_is_after_union_distinct",
		CH:     "建议使用UNION ALL,替代UNION",
		RuleId: "00090",
	},
	{
		Name:   "ddl_check_is_exist_limit_offset",
		CH:     "使用分页查询时，避免使用偏移量",
		RuleId: "00045",
	},
	{
		Name:   "ddl_check_index_option",
		CH:     "建议索引字段对区分度大于阈值",
		RuleId: "00039",
	},
	{
		Name:   "ddl_check_column_enum_notice",
		CH:     "不建议使用 ENUM 类型",
		RuleId: "00019",
	},
	{
		Name:   "ddl_check_column_blob_notice",
		CH:     "不建议使用 BLOB 或 TEXT 类型",
		RuleId: "00017",
	},
	{
		Name:   "ddl_check_create_view",
		CH:     "禁止使用视图",
		RuleId: "00031",
	},
	{
		Name:   "ddl_check_create_trigger",
		CH:     "禁止使用触发器",
		RuleId: "00030",
	},
	{
		Name:   "ddl_check_create_function",
		CH:     "禁止使用自定义函数",
		RuleId: "00014",
	},
	{
		Name:   "ddl_check_create_procedure",
		CH:     "禁止使用存储过程",
		RuleId: "00029",
	},
	{
		Name:   "ddl_disable_type_timestamp",
		CH:     "不建议使用TIMESTAMP字段",
		RuleId: "00068",
	},
	{
		Name:   "dml_check_alias",
		CH:     "别名不建议与表或列的名字相同",
		RuleId: "00079",
	},
	{
		Name:   "ddl_hint_update_table_charset_will_not_update_field_charset",
		CH:     "不建议修改表的默认字符集",
		RuleId: "00073",
	},
	{
		Name:   "ddl_hint_drop_column",
		CH:     "禁止进行删除列的操作",
		RuleId: "00071",
	},
	{
		Name:   "ddl_hint_drop_primary_key",
		CH:     "禁止进行删除主键的操作",
		RuleId: "00010",
	},
	{
		Name:   "ddl_hint_drop_foreign_key",
		CH:     "禁止进行删除外键的操作",
		RuleId: "00072",
	},
	{
		Name:   "dml_hint_in_null_only_false",
		CH:     "避免使用 IN (NULL) 或者 NOT IN (NULL)",
		RuleId: "00120",
	},
	{
		Name:   "ddl_check_full_width_quotation_marks",
		CH:     "DDL语句中不建议使用中文全角引号",
		RuleId: "00035",
	},
	{
		Name:   "dml_not_recommend_order_by_rand",
		CH:     "不建议使用 ORDER BY RAND()",
		RuleId: "00131",
	},
	{
		Name:   "dml_not_recommend_group_by_constant",
		CH:     "不建议对常量进行 GROUP BY",
		RuleId: "00126",
	},
	{
		Name:   "dml_check_sort_direction",
		CH:     "不建议在 ORDER BY 语句中对多个不同条件使用不同方向的排序",
		RuleId: "00104",
	},
	{
		Name:   "dml_hint_group_by_requires_conditions",
		CH:     "建议为GROUP BY语句添加ORDER BY条件",
		RuleId: "00119",
	},
	{
		Name:   "dml_not_recommend_group_by_expression",
		CH:     "不建议ORDER BY 的条件为表达式",
		RuleId: "00127",
	},
	{
		Name:   "dml_check_sql_length",
		CH:     "建议将过长的SQL分解成几个简单的SQL",
		RuleId: "00107",
	},
	{
		Name:   "dml_not_recommend_having",
		CH:     "不建议使用 HAVING 子句",
		RuleId: "00128",
	},
	{
		Name:   "dml_hint_use_truncate_instead_of_delete",
		CH:     "删除全表时建议使用 TRUNCATE 替代 DELETE",
		RuleId: "00124",
	},
	{
		Name:   "dml_not_recommend_update_pk",
		CH:     "不建议UPDATE主键",
		RuleId: "00134",
	},
	{
		Name:   "ddl_check_column_quantity",
		CH:     "表的列数不建议超过阈值",
		RuleId: "00022",
	},
	{
		Name:   "ddl_table_column_charset_same",
		CH:     "建议列与表使用同一个字符集",
		RuleId: "00075",
	},
	{
		Name:   "ddl_check_column_type_integer",
		CH:     "整型定义建议采用 INT(10) 或 BIGINT(20)",
		RuleId: "00026",
	},
	{
		Name:   "ddl_check_varchar_size",
		CH:     "定义VARCHAR 长度时不建议大于阈值",
		RuleId: "00064",
	},
	{
		Name:   "dml_not_recommend_func_in_where",
		CH:     "应避免在 WHERE 条件中使用函数或其他运算符",
		RuleId: "00111",
	},
	{
		Name:   "dml_hint_sum_func_tips",
		CH:     "避免使用 SUM(COL)",
		RuleId: "00122",
	},
	{
		Name:   "ddl_check_column_quantity_in_pk",
		CH:     "主键包含的列数不建议超过阈值",
		RuleId: "00023",
	},
	{
		Name:   "dml_hint_limit_must_be_combined_with_order_by",
		CH:     "LIMIT 查询建议使用ORDER BY",
		RuleId: "00121",
	},
	{
		Name:   "dml_hint_truncate_tips",
		CH:     "不建议使用TRUNCATE操作",
		RuleId: "00123",
	},
	{
		Name:   "dml_hint_delete_tips",
		CH:     "建议在执行DELETE/DROP/TRUNCATE等操作前进行备份",
		RuleId: "00118",
	},
	{
		Name:   "dml_check_sql_injection_func",
		CH:     "不建议使用常见 SQL 注入函数",
		RuleId: "00106",
	},
	{
		Name:   "dml_check_not_equal_symbol",
		CH:     "建议使用'<>'代替'!='",
		RuleId: "00095",
	},
	{
		Name:   "dml_not_recommend_subquery",
		CH:     "不推荐使用子查询",
		RuleId: "00132",
	},
	{
		Name:   "dml_check_subquery_limit",
		CH:     "不建议在子查询中使用LIMIT",
		RuleId: "00109",
	},
	{
		Name:   "ddl_check_auto_increment",
		CH:     "表的初始AUTO_INCREMENT值建议为0",
		RuleId: "00004",
	},
	{
		Name:   "ddl_not_allow_renaming",
		CH:     "禁止使用RENAME或CHANGE对表名字段名进行修改",
		RuleId: "00074",
	},
	{
		Name:   "dml_check_explain_full_index_scan",
		CH:     "不建议对表进行全索引扫描",
		RuleId: "00085",
	},
	{
		Name:   "dml_check_limit_offset_num",
		CH:     "不建议LIMIT的偏移OFFSET大于阈值",
		RuleId: "00045",
	},
	{
		Name:   "dml_check_update_or_delete_has_where",
		CH:     "建议UPDATE/DELETE操作使用WHERE条件",
		RuleId: "00001",
	},
	{
		Name:   "dml_check_order_by_field_length",
		CH:     "禁止对长字段排序",
		RuleId: "00097",
	},
	{
		Name:   "all_check_prepare_statement_placeholders",
		CH:     "绑定的变量个数不建议超过阈值",
		RuleId: "00002",
	},
	{
		Name:   "dml_check_explain_extra_using_index_for_skip_scan",
		CH:     "不建议对表进行索引跳跃扫描",
		RuleId: "00083",
	},
	{
		Name:   "dml_check_affected_rows",
		CH:     "UPDATE/DELETE操作影响行数不建议超过阈值",
		RuleId: "00076",
	},
	{
		Name:   "dml_check_same_table_joined_multiple_times",
		CH:     "不建议对同一张表连接多次",
		RuleId: "00098",
	},
	{
		Name:   "dml_check_using_index",
		CH:     "SQL查询条件需要走索引",
		RuleId: "00110",
	},
	{
		Name:   "dml_check_insert_select",
		CH:     "不建议使用INSERT ... SELECT",
		RuleId: "00089",
	},
	{
		Name:   "dml_check_aggregate",
		CH:     "不建议使用聚合函数",
		RuleId: "00078",
	},
	{
		Name:   "ddl_check_column_not_null",
		CH:     "表字段建议有NOT NULL约束",
		RuleId: "00021",
	},
	{
		Name:   "dml_check_index_selectivity",
		CH:     "建议连库查询时，确保SQL执行计划中使用的索引区分度大于阈值",
		RuleId: "00039",
	},
	{
		Name:   "dml_check_math_computation_or_func_on_index",
		CH:     "禁止对索引列进行数学运算和使用函数",
		RuleId: "00009",
	},
}
