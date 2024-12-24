//go:build enterprise
// +build enterprise

package sqlrewriting

// 现有规则名 到 重写功能使用的 新规则 ID 的映射
// 只映射了部分规则，用于demo版本
// TODO: 当SQLE使用 新规则ID后，可以去除这个映射
/*
SQLE现有MySQL规则: 154
MySQL映射规则: 118 (以下列出)
MySQL未映射规则: 36
*/
var MySQLRuleIdConvert = []ruleIdConvert{
	{
		Name:   "dml_check_in_query_limit",
		CH:     "WHERE条件内IN语句中的参数个数不能超过阈值",
		RuleId: "00087",
	},
	{
		Name:   "ddl_check_column_char_length",
		CH:     "CHAR长度大于20时，必须使用VARCHAR类型",
		RuleId: "00018",
	},
	{
		RuleId: "00063",
		Name:   "ddl_check_unique_index",
		CH:     "建议UNIQUE索引名使用 IDX_UK_表名_字段名",
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
		Name:   "dml_check_where_exist_implicit_conversion",
		CH:     "不建议在WHERE条件中使用与过滤字段不一致的数据类型",
		RuleId: "00112",
	},
	{
		CH:     "单字段上的索引数量不建议超过阈值",
		RuleId: "00043",
		Name:   "ddl_check_index_too_many",
	},
	{
		Name:   "ddl_check_pk_without_auto_increment",
		CH:     "主键建议使用自增",
		RuleId: "00052",
	},
	{
		RuleId: "00003",
		Name:   "ddl_check_all_index_not_null_constraint",
		CH:     "建议为至少一个索引添加非空约束",
	},
	{
		Name:   "dml_hint_sum_func_tips",
		CH:     "避免使用 SUM(COL)",
		RuleId: "00122",
	},
	{
		Name:   "ddl_check_create_time_column",
		CH:     "建议建表DDL包含创建时间字段且默认值为CURRENT_TIMESTAMP",
		RuleId: "00219",
	},
	{
		Name:   "dml_check_sort_direction",
		CH:     "不建议在 ORDER BY 语句中对多个不同条件使用不同方向的排序",
		RuleId: "00104",
	},
	{
		Name:   "dml_check_where_exist_func",
		CH:     "避免对条件字段使用函数操作",
		RuleId: "00111",
	},
	{
		Name:   "all_check_prepare_statement_placeholders",
		CH:     "绑定的变量个数不建议超过阈值",
		RuleId: "00002",
	},
	{
		Name:   "dml_hint_in_null_only_false",
		CH:     "避免使用 IN (NULL) 或者 NOT IN (NULL)",
		RuleId: "00120",
	},
	{
		Name:   "dml_check_insert_select",
		CH:     "不建议使用INSERT ... SELECT",
		RuleId: "00089",
	},
	{
		Name:   "ddl_check_index_count",
		CH:     "索引个数建议不超过阈值",
		RuleId: "00037",
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
		Name:   "dml_check_affected_rows",
		CH:     "UPDATE/DELETE操作影响行数不建议超过阈值",
		RuleId: "00076",
	},
	{
		RuleId: "00108",
		Name:   "dml_check_sub_query_depth",
		CH:     "子查询嵌套层数不建议超过阈值",
	},
	{
		Name:   "dml_check_explain_extra_using_temporary",
		CH:     "不建议使用临时表",
		RuleId: "00084",
	},
	{
		Name:   "ddl_check_create_trigger",
		CH:     "禁止使用触发器",
		RuleId: "00030",
	},
	{
		Name:   "dml_check_where_exist_not",
		CH:     "不建议对条件字段使用负向查询",
		RuleId: "00113",
	},
	{
		Name:   "dml_check_alias",
		CH:     "别名不建议与表或列的名字相同",
		RuleId: "00079",
	},
	{
		Name:   "ddl_check_object_name_is_upper_and_lower_letter_mixed",
		CH:     "数据库对象命名不建议大小写字母混合",
		RuleId: "00046",
	},
	{
		Name:   "ddl_check_pk_without_bigint_unsigned",
		CH:     "主键建议使用 BIGINT 无符号类型，即 BIGINT UNSIGNED",
		RuleId: "00054",
	},
	{
		RuleId: "00001",
		Name:   "all_check_where_is_invalid",
		CH:     "禁止使用没有WHERE条件或者WHERE条件恒为TRUE的SQL",
	},
	{
		Name:   "ddl_check_table_without_comment",
		CH:     "表建议添加注释",
		RuleId: "00060",
	},
	{
		Name:   "ddl_hint_update_table_charset_will_not_update_field_charset",
		CH:     "不建议修改表的默认字符集",
		RuleId: "00073",
	},
	{
		CH:     "建议使用'<>'代替'!='",
		RuleId: "00095",
		Name:   "dml_check_not_equal_symbol",
	},
	{
		Name:   "ddl_check_pk_not_exist",
		CH:     "表必须有主键",
		RuleId: "00008",
	},
	{
		Name:   "ddl_check_composite_index_max",
		CH:     "复合索引的列数量不建议超过阈值",
		RuleId: "00005",
	},
	{
		Name:   "ddl_check_table_character_set",
		CH:     "建议使用指定数据库字符集",
		RuleId: "00015",
	},
	{
		Name:   "dml_check_update_or_delete_has_where",
		CH:     "建议UPDATE/DELETE操作使用WHERE条件",
		RuleId: "00001",
	},
	{
		Name:   "ddl_check_index_prefix",
		CH:     "建议普通索引使用固定前缀",
		RuleId: "00040",
	},
	{
		Name:   "dml_check_select_for_update",
		CH:     "不建议使用SELECT FOR UPDATE",
		RuleId: "00099",
	},
	{
		Name:   "dml_check_needless_func",
		CH:     "避免使用不必要的内置函数",
		RuleId: "00094",
	},
	{
		CH:     "不建议使用TRUNCATE操作",
		RuleId: "00123",
		Name:   "dml_hint_truncate_tips",
	},
	{
		Name:   "dml_disable_select_all_column",
		CH:     "不建议使用SELECT *",
		RuleId: "00053",
	},
	{
		CH:     "不推荐使用子查询",
		RuleId: "00132",
		Name:   "dml_not_recommend_subquery",
	},
	{
		Name:   "ddl_check_database_suffix",
		CH:     "建议数据库名称使用固定后缀结尾",
		RuleId: "00032",
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
		Name:   "dml_check_limit_must_exist",
		CH:     "建议DELETE/UPDATE 语句带有LIMIT条件",
		RuleId: "00092",
	},
	{
		Name:   "ddl_check_object_name_length",
		CH:     "表名、列名、索引名的长度不建议超过阈值",
		RuleId: "00047",
	},
	{
		RuleId: "00006",
		Name:   "dml_check_join_field_type",
		CH:     "建议JOIN字段类型保持一致",
	},
	{
		Name:   "ddl_check_pk_prohibit_auto_increment",
		CH:     "不建议主键使用自增",
		RuleId: "00051",
	},
	{
		CH:     "ALTER表字段禁止使用FIRST,AFTER",
		RuleId: "00065",
		Name:   "ddl_disable_alter_field_use_first_and_after",
	},
	{
		Name:   "dml_check_with_order_by",
		CH:     "DELETE/UPDATE 语句不能有ORDER BY",
		RuleId: "00102",
	},
	{
		Name:   "ddl_hint_drop_primary_key",
		CH:     "禁止进行删除主键的操作",
		RuleId: "00010",
	},
	{
		Name:   "dml_check_insert_columns_exist",
		CH:     "INSERT 语句需要指定COLUMN",
		RuleId: "00088",
	},
	{
		Name:   "dml_not_recommend_func_in_where",
		CH:     "应避免在 WHERE 条件中使用函数或其他运算符",
		RuleId: "00111",
	},
	{
		Name:   "ddl_check_auto_increment",
		CH:     "表的初始AUTO_INCREMENT值建议为0",
		RuleId: "00004",
	},
	{
		RuleId: "00062",
		Name:   "ddl_check_transaction_isolation_level",
		CH:     "事物隔离级别建议设置成RC",
	},
	{
		Name:   "ddl_check_column_blob_notice",
		CH:     "不建议使用 BLOB 或 TEXT 类型",
		RuleId: "00017",
	},
	{
		RuleId: "00124",
		Name:   "dml_hint_use_truncate_instead_of_delete",
		CH:     "删除全表时建议使用 TRUNCATE 替代 DELETE",
	},
	{
		Name:   "ddl_check_column_quantity",
		CH:     "表的列数不建议超过阈值",
		RuleId: "00022",
	},
	{
		Name:   "ddl_check_column_type_integer",
		CH:     "整型定义建议采用 INT(10) 或 BIGINT(20)",
		RuleId: "00026",
	},
	{
		Name:   "ddl_check_table_size",
		CH:     "不建议对数据量过大的表执行DDL操作",
		RuleId: "00059",
	},
	{
		Name:   "ddl_check_table_db_engine",
		CH:     "建议使用指定数据库引擎",
		RuleId: "00057",
	},
	{
		CH:     "建表时，自增字段只能设置一个",
		RuleId: "00007",
		Name:   "ddl_check_auto_increment_field_num",
	},
	{
		CH:     "精确浮点数建议使用DECIMAL",
		RuleId: "00013",
		Name:   "ddl_check_decimal_type_column",
	},
	{
		Name:   "ddl_check_table_partition",
		CH:     "不建议使用分区表相关功能",
		RuleId: "00058",
	},
	{
		Name:   "ddl_check_is_exist_limit_offset",
		CH:     "使用分页查询时，避免使用偏移量",
		RuleId: "00045",
	},
	{
		Name:   "ddl_check_column_without_comment",
		CH:     "列建议添加注释",
		RuleId: "00027",
	},
	{
		Name:   "dml_check_where_exist_null",
		CH:     "不建议对条件字段使用 NULL 值判断",
		RuleId: "00001",
	},
	{
		Name:   "dml_check_where_exist_scalar_sub_queries",
		CH:     "不建议使用标量子查询",
		RuleId: "00115",
	},
	{
		Name:   "ddl_table_column_charset_same",
		CH:     "建议列与表使用同一个字符集",
		RuleId: "00075",
	},
	{
		Name:   "dml_hint_group_by_requires_conditions",
		CH:     "建议为GROUP BY语句添加ORDER BY条件",
		RuleId: "00119",
	},
	{
		Name:   "ddl_check_varchar_size",
		CH:     "定义VARCHAR 长度时不建议大于阈值",
		RuleId: "00064",
	},
	{
		Name:   "ddl_check_create_function",
		CH:     "禁止使用自定义函数",
		RuleId: "00014",
	},
	{
		RuleId: "00126",
		Name:   "dml_not_recommend_group_by_constant",
		CH:     "不建议对常量进行 GROUP BY",
	},
	{
		Name:   "dml_not_recommend_update_pk",
		CH:     "不建议UPDATE主键",
		RuleId: "00134",
	},
	{
		Name:   "dml_check_using_index",
		CH:     "SQL查询条件需要走索引",
		RuleId: "00110",
	},
	{
		Name:   "ddl_disable_fk",
		CH:     "禁止使用外键",
		RuleId: "00067",
	},
	{
		Name:   "ddl_check_object_name_using_cn",
		CH:     "数据库对象命名只能使用英文、下划线或数字，首字母必须是英文",
		RuleId: "00048",
	},
	{
		RuleId: "00012",
		Name:   "ddl_check_bigint_instead_of_decimal",
		CH:     "建议用BIGINT类型代替DECIMAL",
	},
	{
		CH:     "不建议ORDER BY 的条件为表达式",
		RuleId: "00127",
		Name:   "dml_not_recommend_group_by_expression",
	},
	{
		Name:   "dml_check_sql_length",
		CH:     "建议将过长的SQL分解成几个简单的SQL",
		RuleId: "00107",
	},
	{
		Name:   "dml_check_order_by_field_length",
		CH:     "禁止对长字段排序",
		RuleId: "00097",
	},
	{
		Name:   "dml_check_same_table_joined_multiple_times",
		CH:     "不建议对同一张表连接多次",
		RuleId: "00098",
	},
	{
		Name:   "ddl_check_alter_table_need_merge",
		CH:     "存在多条对同一个表的修改语句，建议合并成一个ALTER语句",
		RuleId: "00011",
	},
	{
		Name:   "dml_check_explain_extra_using_filesort",
		CH:     "不建议使用文件排序",
		RuleId: "00082",
	},
	{
		RuleId: "00049",
		Name:   "ddl_check_object_name_using_keyword",
		CH:     "数据库对象命名禁止使用保留字",
	},
	{
		RuleId: "00086",
		Name:   "dml_check_fuzzy_search",
		CH:     "禁止使用全模糊搜索或左模糊搜索",
	},
	{
		Name:   "ddl_hint_drop_foreign_key",
		CH:     "禁止进行删除外键的操作",
		RuleId: "00072",
	},
	{
		RuleId: "00023",
		Name:   "ddl_check_column_quantity_in_pk",
		CH:     "主键包含的列数不建议超过阈值",
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
		CH:     "不建议对表进行索引跳跃扫描",
		RuleId: "00083",
		Name:   "dml_check_explain_extra_using_index_for_skip_scan",
	},
	{
		Name:   "ddl_check_table_without_if_not_exists",
		CH:     "新建表建议加入 IF NOT EXISTS，保证重复执行不报错",
		RuleId: "00061",
	},
	{
		Name:   "dml_check_join_has_on",
		CH:     "建议连接操作指定连接条件",
		RuleId: "00091",
	},
	{
		Name:   "ddl_check_column_blob_default_is_not_null",
		CH:     "BLOB 和 TEXT 类型的字段默认值只能为NULL",
		RuleId: "00016",
	},
	{
		Name:   "dml_check_batch_insert_lists_max",
		CH:     "单条INSERT语句，建议批量插入不超过阈值",
		RuleId: "00080",
	},
	{
		Name:   "ddl_check_create_procedure",
		CH:     "禁止使用存储过程",
		RuleId: "00029",
	},
	{
		CH:     "建议使用规定的数据库排序规则",
		RuleId: "00015",
		Name:   "ddl_check_collation_database",
	},
	{
		RuleId: "00019",
		Name:   "ddl_check_column_enum_notice",
		CH:     "不建议使用 ENUM 类型",
	},
	{
		CH:     "禁止使用视图",
		RuleId: "00031",
		Name:   "ddl_check_create_view",
	},
	{
		Name:   "dml_not_recommend_having",
		CH:     "不建议使用 HAVING 子句",
		RuleId: "00128",
	},
	{
		Name:   "ddl_check_redundant_index",
		CH:     "不建议创建冗余索引",
		RuleId: "00055",
	},
	{
		Name:   "ddl_check_field_not_null_must_contain_default_value",
		CH:     "建议字段约束为NOT NULL时带默认值",
		RuleId: "00034",
	},
	{
		Name:   "ddl_check_unique_index_prefix",
		CH:     "建议UNIQUE索引使用固定前缀",
		RuleId: "00041",
	},
	{
		CH:     "建议连库查询时，确保SQL执行计划中使用的索引区分度大于阈值",
		RuleId: "00039",
		Name:   "dml_check_index_selectivity",
	},
	{
		Name:   "dml_check_limit_offset_num",
		CH:     "不建议LIMIT的偏移OFFSET大于阈值",
		RuleId: "00045",
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
		Name:   "ddl_check_update_time_column",
		CH:     "建表DDL需要包含更新时间字段且默认值为CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
		RuleId: "00033",
	},
	{
		Name:   "ddl_check_column_timestamp_without_default",
		CH:     "TIMESTAMP 类型的列必须添加默认值",
		RuleId: "00025",
	},
	{
		CH:     "不建议在子查询中使用LIMIT",
		RuleId: "00109",
		Name:   "dml_check_subquery_limit",
	},
	{
		Name:   "ddl_check_column_blob_with_not_null",
		CH:     "BLOB 和 TEXT 类型的字段不建议设置为 NOT NULL",
		RuleId: "00016",
	},
	{
		RuleId: "00100",
		Name:   "dml_check_select_limit",
		CH:     "SELECT 语句需要带LIMIT",
	},
	{
		Name:   "ddl_disable_type_timestamp",
		CH:     "不建议使用TIMESTAMP字段",
		RuleId: "00068",
	},
	{
		Name:   "ddl_hint_drop_column",
		CH:     "禁止进行删除列的操作",
		RuleId: "00071",
	},
	{
		Name:   "ddl_disable_drop_statement",
		CH:     "禁止除索引外的DROP操作",
		RuleId: "00066",
	},
	{
		Name:   "dml_hint_limit_must_be_combined_with_order_by",
		CH:     "LIMIT 查询建议使用ORDER BY",
		RuleId: "00121",
	},
	{
		CH:     "SELECT 语句不能有ORDER BY",
		RuleId: "00101",
		Name:   "dml_check_select_with_order_by",
	},
	{
		Name:   "ddl_check_index_option",
		CH:     "建议索引字段对区分度大于阈值",
		RuleId: "00039",
	},
	{
		Name:   "dml_check_math_computation_or_func_on_index",
		CH:     "禁止对索引列进行数学运算和使用函数",
		RuleId: "00009",
	},
}
