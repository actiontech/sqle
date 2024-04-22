package optimization

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var RuleHandler = []rulepkg.RuleHandler{
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLHintGroupByRequiresConditions,
			Desc:       "为GROUP BY显示添加 ORDER BY 条件(<MYSQL 5.7)",
			Annotation: "在早期版本的MySQL中，GROUP BY 默认进行排序，可通过添加 ORDER BY NULL 来取消此排序，提高查询效率。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "为GROUP BY显示添加 ORDER BY 条件(<MYSQL 5.7)",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLCheckWhereExistScalarSubquery,
			Desc:       "COUNT标量子查询重写",
			Annotation: "对于使用COUNT标量子查询来进行判断是否存在，可以重写为EXISTS子查询，从而避免一次聚集运算。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message:      "COUNT标量子查询重写",
		AllowOffline: true,
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLHintUseTruncateInsteadOfDelete,
			Desc:       "无条件的DELETE建议重写为Truncate",
			Annotation: "TRUNCATE TABLE 比 DELETE 速度快，且使用的系统和事务日志资源少，同时TRUNCATE后表所占用的空间会被释放，而DELETE后需要手工执行OPTIMIZE才能释放表空间",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "无条件的DELETE建议重写为Truncate",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLCheckWhereExistImplicitConversion,
			Desc:       "隐式类型转换导致索引失效",
			Annotation: "WHERE条件中使用与过滤字段不一致的数据类型会引发隐式数据类型转换，导致查询有无法命中索引的风险，在高并发、大数据量的情况下，不走索引会使得数据库的查询性能严重下降",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "隐式类型转换导致索引失效",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DDLCheckDatabaseCollation,
			Desc:       "排序字段方向不同导致索引失效",
			Annotation: "ORDER BY 子句中的所有表达式需要按统一的 ASC 或 DESC 方向排序，才能利用索引来避免排序；如果ORDER BY 语句对多个不同条件使用不同方向的排序无法使用索引",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "排序字段方向不同导致索引失效",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLCheckMathComputationOrFuncOnIndex,
			Desc:       "索引列上的运算导致索引失效",
			Annotation: "在索引列上的运算将导致索引失效，容易造成全表扫描，产生严重的性能问题。所以需要尽量将索引列上的运算转换到常量端进行。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeIndexInvalidation,
		},
		AllowOffline: false,
		Message:      "索引列上的运算导致索引失效",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLNotRecommendHaving,
			Desc:       "HAVING条件下推",
			Annotation: "从逻辑上，HAVING条件是在分组之后执行的，而WHERE子句上的条件可以在表访问的时候（索引访问）,或是表访问之后、分组之前执行，这两种条件都比在分组之后执行代价要小。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "HAVING条件下推",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLWhereExistNull,
			Desc:       "禁止使用=NULL判断空值",
			Annotation: "= null并不能判断表达式为空,= null总是被判断为假。判断表达式为空应该使用is null。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message:      "禁止使用=NULL判断空值",
		Func:         nil,
		AllowOffline: true,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLNotRecommendIn,
			Desc:       "IN子查询优化",
			Annotation: "IN子查询是指符合下面形式的子查询，IN子查询可以改写成等价的相关EXISTS子查询或是内连接，从而可以产生一个新的过滤条件。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "IN子查询优化",
		Func:    nil,
	},

	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLHintInNullOnlyFalse,
			Desc:       "IN可空子查询可能导致结果集不符合预期",
			Annotation: "查询条件永远非真，这将导致查询无匹配到的结果",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "IN可空子查询可能导致结果集不符合预期",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLNotRecommendNotWildcardLike,
			Desc:       "避免使用没有通配符的 LIKE 查询",
			Annotation: "不包含通配符的LIKE 查询逻辑上与等值查询相同，建议使用等值查询替代。而且不包含通配符的LIKE 查询逻辑通常是由于开发者错误导致的，可能不符合其期望的业务逻辑实现",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "避免使用没有通配符的 LIKE 查询",
		Func:    nil,
	},

	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLCheckNotEqualSymbol,
			Desc:       "建议使用'<>'代替'!='",
			Annotation: "'!=' 是非标准的运算符，'<>' 才是SQL中标准的不等于运算符",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message: "建议使用'<>'代替'!='",
		Func:    nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       rulepkg.DMLCheckLimitOffsetNum,
			Desc:       "OFFSET的值超过阈值",
			Annotation: "使用LIMIT和OFFSET子句可以分别控制查询结果的数量和指定从哪一行开始返回数据。但是，当OFFSET值较大时，查询效率会降低，因为系统必须扫描更多数据才能找到起始行，这在大数据集中尤其会导致性能问题和资源消耗。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		Message:      "OFFSET的值超过阈值",
		AllowOffline: true,
		Func:         nil,
	},

	{
		Rule: driverV2.Rule{
			Name:       DMLRuleDistinctEliminationRewrite,
			Desc:       "子查询中的DISTINCT消除",
			Annotation: "对于仅进行存在性测试的子查询,如果子查询包含DISTINCT通常可以删除,以避免一次去重操作。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "子查询中的DISTINCT消除",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleExists2JoinRewrite,
			Desc:       "EXISTS查询转换为表连接",
			Annotation: "EXISTS子查询可以在适当情况下转换为JOIN来优化查询，提高数据库处理效率和性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "EXISTS查询转换为表连接",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleFilterPredicatePushDownRewrite,
			Desc:       "过滤谓词下推",
			Annotation: "滤条件下推（FPPD）是一种通过将过滤条件提前应用于内部查询块，以减少数据处理量并提升SQL执行效率。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "过滤谓词下推",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleGroupingFromDiffTablesRewrite,
			Desc:       "GROUPBY字段来自不同表",
			Annotation: "如果分组字段来自不同的表，数据库优化器将没有办法利用索引的有序性来避免一次排序，如果存在等值条件，可以替换这些字段为来自同一张表的字段，以利用索引优化排序和提高查询效率。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "GROUPBY字段来自不同表",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleJoinEliminationRewrite,
			Desc:       "表连接消除",
			Annotation: "在不影响结果的情况下通过删除不必要的表连接来简化查询并提升性能，适用于查询仅涉及到主表主键列的场景。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "表连接消除",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleLimitClausePushDownRewrite,
			Desc:       "LIMIT下推至UNION分支",
			Annotation: "Limit子句下推优化通过尽可能的 “下压” Limit子句，提前过滤掉部分数据, 减少中间结果集的大小，减少后续计算需要处理的数据量, 以提高查询性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "LIMIT下推至UNION分支",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleMaxMinAggRewrite,
			Desc:       "MAX/MIN子查询重写",
			Annotation: "对于使用MAX/MIN的子查询，可以通过重写从而利用索引的有序来避免一次聚集运算。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "MAX/MIN子查询重写",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleMoveOrder2LeadingRewrite,
			Desc:       "ORDER子句重排序优化",
			Annotation: "如果一个查询中既包含来自同一个表的排序字段也包含分组字段，但字段顺序不同，可以通过调整分组字段顺序，使其和排序字段顺序一致，这样数据库可以避免一次排序操作。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "ORDER子句重排序优化",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleOrCond4SelectRewrite,
			Desc:       "OR条件的SELECT重写",
			Annotation: "如果使用OR条件的查询语句，数据库优化器有可能无法使用索引来完成查询，可以把查询语句重写为UNION或UNION ALL查询，以便使用索引提升查询性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "OR条件的SELECT重写",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleOrCond4UpDeleteRewrite,
			Desc:       "OR条件的UPDELETE重写",
			Annotation: "如果有使用OR条件的UPDATE或DELETE语句，数据库优化器有可能无法使用索引来完成操作，可以把它重写为多个DELETE语句，利用索引提升查询性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "OR条件的UPDELETE重写",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleOrderEliminationInSubqueryRewrite,
			Desc:       "IN子查询中没有LIMIT的排序消除",
			Annotation: "如果子查询没有LIMIT子句，那么子查询的排序操作就没有意义，可以将其删除而不影响最终的结果。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "IN子查询中没有LIMIT的排序消除",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleOrderingFromDiffTablesRewrite,
			Desc:       "避免ORDERBY字段来自不同表",
			Annotation: "当排序字段来自不同表时，若存在等值条件，可替换这些字段为来自同一张表的字段，利用索引避免额外排序，提升效率。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "避免ORDERBY字段来自不同表",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleOuter2InnerConversionRewrite,
			Desc:       "外连接优化",
			Annotation: "外连接优化指的是满足一定条件（外表具有NULL拒绝条件）的外连接可以转化为内连接，从而可以让数据库优化器可以选择更优的执行计划，提升SQL查询的性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "外连接优化",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleProjectionPushdownRewrite,
			Desc:       "投影下推(PROJECTION PUSHDOWN)",
			Annotation: "投影下推指的通过删除DT子查询中无意义的列（在外查询中没有使用），来减少IO和网络的代价，同时提升优化器在进行表访问的规划时，采用无需回表的优化选项的几率。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "投影下推(PROJECTION PUSHDOWN)",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleQualifierSubQueryRewrite,
			Desc:       "修饰子查询重写优化",
			Annotation: "ANY/SOME/ALL修饰的子查询用于比较值关系，但效率低下因为它们逐行处理比较。通过查询重写可以提升这类子查询的执行效率。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "修饰子查询重写优化",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleQueryFoldingRewrite,
			Desc:       "查询折叠(QUERY FOLDING)",
			Annotation: "查询折叠指的是把视图、CTE或是DT子查询展开，并与引用它的查询语句合并，来减少序列化中间结果集，或是触发更优的关于表连接规划的优化技术。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "查询折叠(QUERY FOLDING)",
		Func:         nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLRuleSATTCRewrite,
			Desc:       "SATTC重写优化",
			Annotation: "SAT-TC重写优化通过分析和处理查询条件的逻辑关系，以发现矛盾、简化条件或推断新条件，从而帮助数据库优化器制定更高效的执行计划，提升SQL性能。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
		},
		AllowOffline: false,
		Message:      "SATTC重写优化",
		Func:         nil,
	},
}
