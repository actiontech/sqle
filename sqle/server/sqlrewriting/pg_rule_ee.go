//go:build enterprise
// +build enterprise

package sqlrewriting

// 现有规则名 到 重写功能使用的 新规则 ID 的映射
// 只映射了部分规则，用于demo版本
// TODO: 当SQLE使用 新规则ID后，可以去除这个映射
/*
SQLE现有PostgreSQL规则: 71
PostgreSQL映射规则: 44 (以下列出)
PostgreSQL未映射规则: 27
*/
var PostgreSQLRuleIdConvert = []ruleIdConvert{
	{
		CH:     "执行计划存在笛卡尔积，建议检查表关联条件",
		RuleId: "00091",
		Name:   "pg_028",
	},
	{
		RuleId: "00115",
		Name:   "pg_033",
		CH:     "不建议使用标量子查询",
	},
	{
		RuleId: "00192",
		Name:   "pg_051",
		CH:     "外键字段上必须创建索引",
	},
	{
		Name:   "pg_053",
		CH:     "索引最左侧的字段必须出现在查询条件内",
		RuleId: "00218",
	},
	{
		RuleId: "00033",
		Name:   "pg_049",
		CH:     "建表DDL必须包含更新时间字段，类型为TIMESTAMP",
	},
	{
		Name:   "pg_054",
		CH:     "数据库对象名必须只包含小写字母，下划线，数字。且不能以数字开头",
		RuleId: "00048",
	},
	{
		Name:   "pg_055",
		CH:     "绑定的变量个数不建议超过阈值",
		RuleId: "00002",
	},
	{
		Name:   "pg_001",
		CH:     "不建议使用select *",
		RuleId: "00053",
	},
	{
		Name:   "pg_023",
		CH:     "建议给列添加注释",
		RuleId: "00027",
	},
	{
		Name:   "pg_038",
		CH:     "查询语句不建议使用UNION",
		RuleId: "00090",
	},
	{
		Name:   "pg_008",
		CH:     "表不建议使用外键",
		RuleId: "00067",
	},
	{
		CH:     "不建议修改列的数据类型",
		RuleId: "00059",
		Name:   "pg_026",
	},
	{
		Name:   "pg_037",
		CH:     "不建议删除列",
		RuleId: "00071",
	},
	{
		Name:   "pg_056",
		CH:     "临时表必须使用固定前缀",
		RuleId: "00042",
	},
	{
		Name:   "pg_032",
		CH:     "禁止使用没有WHERE条件或者WHERE条件恒为TRUE的SQL",
		RuleId: "00001",
	},
	{
		Name:   "pg_022",
		CH:     "建议给表添加注释",
		RuleId: "00060",
	},
	{
		Name:   "pg_034",
		CH:     "不建议对多个字段使用OR条件查询",
		RuleId: "00143",
	},
	{
		CH:     "每个对象名称前面必须带有属主",
		RuleId: "00140",
		Name:   "pg_035",
	},
	{
		Name:   "pg_044",
		CH:     "多表关联时，不建议在WHERE条件中对不同表的字段使用OR条件",
		RuleId: "00143",
	},
	{
		Name:   "pg_016",
		CH:     "禁止使用左模糊查询",
		RuleId: "00086",
	},
	{
		RuleId: "00071",
		Name:   "pg_036",
		CH:     "不建议删除对象",
	},
	{
		Name:   "pg_052",
		CH:     "建议索引字段的区分度大于阈值",
		RuleId: "00039",
	},
	{
		Name:   "pg_040",
		CH:     "主键建议使用自增",
		RuleId: "00052",
	},
	{
		Name:   "pg_002",
		CH:     "delete 和 update 语句，必须带where条件",
		RuleId: "00001",
	},
	{
		Name:   "pg_003",
		CH:     "避免使用 having 子句",
		RuleId: "00128",
	},
	{
		CH:     "禁止使用触发器",
		RuleId: "00030",
		Name:   "pg_006",
	},
	{
		Name:   "pg_011",
		CH:     "复合索引的列数量不建议超过阈值",
		RuleId: "00005",
	},
	{
		Name:   "pg_031",
		CH:     "不建议创建冗余索引",
		RuleId: "00055",
	},
	{
		Name:   "pg_048",
		CH:     "建表DDL必须包含创建时间字段，类型为TIMESTAMP",
		RuleId: "00219",
	},
	{
		Name:   "pg_010",
		CH:     "单条SQL不建议过长",
		RuleId: "00107",
	},
	{
		Name:   "pg_042",
		CH:     "在WHERE条件中禁止对索引列使用函数或表达式",
		RuleId: "00111",
	},
	{
		Name:   "pg_050",
		CH:     "对象名称字符个数不建议超过阈值",
		RuleId: "00047",
	},
	{
		Name:   "pg_007",
		CH:     "表建议使用主键",
		RuleId: "00008",
	},
	{
		Name:   "pg_009",
		CH:     "表字段不建议过多",
		RuleId: "00020",
	},
	{
		Name:   "pg_012",
		CH:     "普通索引必须使用固定前缀",
		RuleId: "00040",
	},
	{
		RuleId: "00096",
		Name:   "pg_041",
		CH:     "表连接数不建议超过阈值",
	},
	{
		Name:   "pg_045",
		CH:     "不建议使用全表扫描",
		RuleId: "00139",
	},
	{
		Name:   "pg_059",
		CH:     "JOIN表时，关联字段的类型必须保持一致",
		RuleId: "00006",
	},
	{
		Name:   "pg_004",
		CH:     "禁止除索引外的 drop 操作",
		RuleId: "00066",
	},
	{
		Name:   "pg_015",
		CH:     "数据库对象命名禁止使用关键字",
		RuleId: "00049",
	},
	{
		Name:   "pg_021",
		CH:     "建议避免使用select for update",
		RuleId: "00099",
	},
	{
		Name:   "pg_027",
		CH:     "insert语句未指定列信息",
		RuleId: "00088",
	},
	{
		Name:   "pg_039",
		CH:     "同一张表的索引字段不建议超过阈值",
		RuleId: "00005",
	},
	{
		CH:     "禁止使用视图",
		RuleId: "00031",
		Name:   "pg_005",
	},
}
