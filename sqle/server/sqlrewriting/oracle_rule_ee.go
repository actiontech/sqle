//go:build enterprise
// +build enterprise

package sqlrewriting

// 现有规则名 到 重写功能使用的 新规则 ID 的映射
// 只映射了部分规则，用于demo版本
// TODO: 当SQLE使用 新规则ID后，可以去除这个映射
/*
SQLE现有Oracle规则: 85
Oracle映射规则: 74 (以下列出)
Oracle未映射规则: 11
*/
var OracleRuleIdConvert = []ruleIdConvert{
	{
		CH:     "序列名称不规范",
		RuleId: "00160",
		Name:   "Oracle_035",
	},
	{
		Name:   "Oracle_002",
		CH:     "delete 和 update 语句，必须带where条件",
		RuleId: "00001",
	},
	{
		Name:   "Oracle_024",
		CH:     "数据库对象命名禁止使用关键字",
		RuleId: "00049",
	},
	{
		Name:   "Oracle_027",
		CH:     "序列步长大于1",
		RuleId: "00161",
	},
	{
		Name:   "Oracle_078",
		CH:     "UPDATE和DELETE语句评估影响行数过大",
		RuleId: "00076",
	},
	{
		CH:     "索引列区分度低",
		RuleId: "00039",
		Name:   "Oracle_080",
	},
	{
		Name:   "Oracle_010",
		CH:     "单条SQL不建议过长",
		RuleId: "00107",
	},
	{
		Name:   "Oracle_031",
		CH:     "CREATE TABLE/INDEX使用禁止的表空间",
		RuleId: "00151",
	},
	{
		CH:     "单条insert语句，建议批量插入不超过阈值",
		RuleId: "00080",
		Name:   "Oracle_065",
	},
	{
		Name:   "Oracle_048",
		CH:     "建表DDL必须包含创建时间字段且默认值为SYSDATE，创建时间字段名: CREATE_TIME",
		RuleId: "00219",
	},
	{
		Name:   "Oracle_049",
		CH:     "建表DDL必须包含更新时间字段且默认值为SYSDATE，更新时间字段名: UPDATE_TIME",
		RuleId: "00033",
	},
	{
		CH:     "序列cache值设置不合理",
		RuleId: "00168",
		Name:   "Oracle_042",
	},
	{
		CH:     "禁止使用位图索引",
		RuleId: "00159",
		Name:   "Oracle_034",
	},
	{
		Name:   "Oracle_066",
		CH:     "DML语句中使用了order by",
		RuleId: "00101",
	},
	{
		Name:   "Oracle_062",
		CH:     "不建议SQL中包含hint指令",
		RuleId: "00176",
	},
	{
		Name:   "Oracle_051",
		CH:     "不建议使用 BLOB 类型",
		RuleId: "00017",
	},
	{
		Name:   "Oracle_025",
		CH:     "表建议添加注释",
		RuleId: "00060",
	},
	{
		Name:   "Oracle_023",
		CH:     "普通索引必须使用固定前缀",
		RuleId: "00040",
	},
	{
		Name:   "Oracle_041",
		CH:     "序列不推荐设置不循环",
		RuleId: "00167",
	},
	{
		Name:   "Oracle_013",
		CH:     "unique索引必须使用固定前缀",
		RuleId: "00041",
	},
	{
		Name:   "Oracle_075",
		CH:     "SQL语句使用多层次嵌套",
		RuleId: "00108",
	},
	{
		Name:   "Oracle_069",
		CH:     "SQL语句存在全表排序操作",
		RuleId: "00178",
	},
	{
		Name:   "Oracle_037",
		CH:     "创建索引或重建索引时缺少并行",
		RuleId: "00163",
	},
	{
		Name:   "Oracle_036",
		CH:     "建表索引时未指定索引表空间",
		RuleId: "00162",
	},
	{
		Name:   "Oracle_055",
		CH:     "建议避免使用select for update",
		RuleId: "00099",
	},
	{
		Name:   "Oracle_059",
		CH:     "禁止使用没有where条件的sql语句",
		RuleId: "00001",
	},
	{
		CH:     "创建或重建索引缺少online",
		RuleId: "00165",
		Name:   "Oracle_039",
	},
	{
		RuleId: "00001",
		Name:   "Oracle_061",
		CH:     "禁止where条件中出现1=1",
	},
	{
		Name:   "Oracle_067",
		CH:     "INSERT语句未指定字段",
		RuleId: "00088",
	},
	{
		Name:   "Oracle_050",
		CH:     "表中包含有太多的列",
		RuleId: "00020",
	},
	{
		RuleId: "00079",
		Name:   "Oracle_063",
		CH:     "别名不要与表或列的名字相同",
	},
	{
		CH:     "条件字段做函数操作",
		RuleId: "00111",
		Name:   "Oracle_074",
	},
	{
		CH:     "char长度大于20时，必须使用varchar2类型",
		RuleId: "00018",
		Name:   "Oracle_052",
	},
	{
		Name:   "Oracle_073",
		CH:     "使用标量子查询",
		RuleId: "00115",
	},
	{
		RuleId: "00166",
		Name:   "Oracle_040",
		CH:     "序列建议设置cache值",
	},
	{
		Name:   "Oracle_056",
		CH:     "禁止使用并行属性",
		RuleId: "00164",
	},
	{
		CH:     "参与连接操作的表数量太多",
		RuleId: "00096",
		Name:   "Oracle_076",
	},
	{
		Name:   "Oracle_030",
		CH:     "建议使用DATE替代TIMESTAMP类型",
		RuleId: "00157",
	},
	{
		RuleId: "00086",
		Name:   "Oracle_064",
		CH:     "使用了全模糊查询或左模糊查询",
	},
	{
		RuleId: "00153",
		Name:   "Oracle_045",
		CH:     "创建表建议添加索引",
	},
	{
		Name:   "Oracle_047",
		CH:     "禁止删除列",
		RuleId: "00071",
	},
	{
		Name:   "Oracle_068",
		CH:     "ORDER BY字段数过多",
		RuleId: "00177",
	},
	{
		Name:   "Oracle_046",
		CH:     "相同字段类型不能修改字段长度(变小)",
		RuleId: "00170",
	},
	{
		Name:   "Oracle_070",
		CH:     "SQL中存在嵌套子查询",
		RuleId: "00108",
	},
	{
		Name:   "Oracle_054",
		CH:     "视图定义中禁止带rownum",
		RuleId: "00173",
	},
	{
		Name:   "Oracle_026",
		CH:     "列建议添加注释",
		RuleId: "00027",
	},
	{
		Name:   "Oracle_020",
		CH:     "复合索引的列数量不建议超过阈值",
		RuleId: "00005",
	},
	{
		Name:   "Oracle_072",
		CH:     "对条件字段使用负向查询",
		RuleId: "00113",
	},
	{
		CH:     "数据库对象命名只能使用英文、下划线或数字，首字母必须是英文",
		RuleId: "00048",
		Name:   "Oracle_029",
	},
	{
		Name:   "Oracle_007",
		CH:     "表建议使用主键",
		RuleId: "00008",
	},
	{
		Name:   "Oracle_003",
		CH:     "避免使用 having 子句",
		RuleId: "00128",
	},
	{
		Name:   "Oracle_006",
		CH:     "禁止创建触发器",
		RuleId: "00030",
	},
	{
		CH:     "禁止创建视图",
		RuleId: "00031",
		Name:   "Oracle_005",
	},
	{
		Name:   "Oracle_022",
		CH:     "表名、列名、索引名的长度不能大于指定字节",
		RuleId: "00047",
	},
	{
		Name:   "Oracle_004",
		CH:     "禁止除索引外的 drop 操作",
		RuleId: "00066",
	},
	{
		Name:   "Oracle_060",
		CH:     "建议使用UNION ALL替代UNION",
		RuleId: "00090",
	},
	{
		Name:   "Oracle_028",
		CH:     "创建表时未指定表空间",
		RuleId: "00148",
	},
	{
		RuleId: "00132",
		Name:   "Oracle_009",
		CH:     "不推荐使用子查询",
	},
	{
		Name:   "Oracle_001",
		CH:     "禁止使用 select *",
		RuleId: "00053",
	},
	{
		Name:   "Oracle_011",
		CH:     "索引个数建议不超过阈值",
		RuleId: "00037",
	},
	{
		RuleId: "00020",
		Name:   "Oracle_012",
		CH:     "表字段过多",
	},
	{
		Name:   "Oracle_033",
		CH:     "建议表名和关键字大写",
		RuleId: "00158",
	},
	{
		Name:   "Oracle_032",
		CH:     "where条件内in语句中的参数个数不能超过阈值",
		RuleId: "00087",
	},
	{
		Name:   "Oracle_016",
		CH:     "表关联个数过多",
		RuleId: "00096",
	},
	{
		Name:   "Oracle_008",
		CH:     "表不建议使用外键",
		RuleId: "00067",
	},
	{
		Name:   "Oracle_014",
		CH:     "表关联存在笛卡尔积",
		RuleId: "00091",
	},
	{
		CH:     "执行计划中存在索引快速全扫",
		RuleId: "00085",
		Name:   "Oracle_082",
	},
	{
		Name:   "Oracle_058",
		CH:     "执行计划中存在bitmap-and步骤",
		RuleId: "00175",
	},
	{
		Name:   "Oracle_017",
		CH:     "对大表进行全表扫描",
		RuleId: "00139",
	},
	{
		Name:   "Oracle_019",
		CH:     "对大索引执行全扫描",
		RuleId: "00156",
	},
	{
		CH:     "执行计划中存在索引全扫",
		RuleId: "00182",
		Name:   "Oracle_083",
	},
	{
		Name:   "Oracle_081",
		CH:     "执行计划中存在filter步骤",
		RuleId: "00180",
	},
	{
		Name:   "Oracle_079",
		CH:     "建议避免数据类型隐式转换",
		RuleId: "00179",
	},
	{
		Name:   "Oracle_018",
		CH:     "存在索引执行跳跃扫描",
		RuleId: "00083",
	},
}
