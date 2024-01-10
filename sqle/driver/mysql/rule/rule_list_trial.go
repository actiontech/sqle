//go:build trial
// +build trial

package rule

import (
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

var RuleHandlers = []RuleHandler{
	{
		Rule: driverV2.Rule{
			Name:       DMLCheckFuzzySearch,
			Desc:       "禁止使用全模糊搜索或左模糊搜索",
			Annotation: "使用全模糊搜索或左模糊搜索将导致查询无法使用索引，导致全表扫描",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeDMLConvention,
		},
		Message:      "禁止使用全模糊搜索或左模糊搜索",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLCheckJoinFieldType,
			Desc:       "建议JOIN字段类型保持一致",
			Annotation: "JOIN字段类型不一致会导致类型不匹配发生隐式准换，建议开启此规则，避免索引失效",
			Level:      driverV2.RuleLevelWarn,
			Category:   RuleTypeDMLConvention,
		},
		Message:      "建议JOIN字段类型保持一致, 否则会导致隐式转换",
		AllowOffline: false,
		Func:         checkJoinFieldType,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLRecommendTableColumnCharsetSame,
			Desc:       "建议列与表使用同一个字符集",
			Annotation: "统一字符集可以避免由于字符集转换产生的乱码，不同的字符集进行比较前需要进行转换会造成索引失效",
			Level:      driverV2.RuleLevelWarn,
			Category:   RuleTypeDDLConvention,
		},
		Message: "建议列与表使用同一个字符集",
		Func:    recommendTableColumnCharsetSame,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckColumnTimestampWithoutDefault,
			Desc:       "TIMESTAMP 类型的列必须添加默认值",
			Annotation: "TIMESTAMP添加默认值，可避免出现全为0的日期格式与业务预期不符",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeDDLConvention,
		},
		Message:      "TIMESTAMP 类型的列必须添加默认值",
		AllowOffline: true,
		Func:         checkColumnTimestampWithoutDefault,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckIndexPrefix,
			Desc:       "建议普通索引使用固定前缀",
			Annotation: "通过配置该规则可以规范指定业务的索引命名规则，具体命名规范可以自定义设置，默认提示值：idx_",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeNamingConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "idx_",
					Desc:  "索引前缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "建议普通索引要以\"%v\"为前缀",
		AllowOffline: true,
		Func:         checkIndexPrefix,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckPKNotExist,
			Desc:       "表必须有主键",
			Annotation: "主键使数据达到全局唯一，可提高数据检索效率",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeIndexingConvention,
		},
		Message:                         "表必须有主键",
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                            checkPrimaryKey,
	},
	{
		Rule: driverV2.Rule{
			Name:       ConfigDMLExplainPreCheckEnable,
			Desc:       "使用EXPLAIN加强预检查能力",
			Annotation: "通过 EXPLAIN 的形式将待上线的DML进行SQL是否能正确执行的检查，提前发现语句的错误，提高上线成功率",
			Level:      driverV2.RuleLevelWarn,
			Category:   RuleTypeGlobalConfig,
		},
		Func: nil,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckIndexCount,
			Desc:       "索引个数建议不超过阈值",
			Annotation: "在表上建立的每个索引都会增加存储开销，索引对于插入、删除、更新操作也会增加处理上的开销，太多与不充分、不正确的索引对性能都毫无益处；具体规则阈值可以根据业务需求调整，默认值：5",
			Level:      driverV2.RuleLevelNotice,
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
		Message:                         "索引个数建议不超过%v个",
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                            checkIndex,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckPKWithoutAutoIncrement,
			Desc:       "主键建议使用自增",
			Annotation: "自增主键，数字型速度快，而且是增量增长，占用空间小，更快速的做数据插入操作，避免增加维护索引的开销",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeIndexingConvention,
		},
		Message:                         "主键建议使用自增",
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                            checkPrimaryKey,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckObjectNameUsingKeyword,
			Desc:       "数据库对象命名禁止使用保留字",
			Annotation: "通过配置该规则可以规范指定业务的数据对象命名规则，避免发生冲突，以及混淆",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeNamingConvention,
		},
		Message:      "数据库对象命名禁止使用保留字 %s",
		AllowOffline: true,
		Func:         checkNewObjectName,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLCheckMathComputationOrFuncOnIndex,
			Desc:       "禁止对索引列进行数学运算和使用函数",
			Annotation: "对索引列进行数学运算和使用函数会导致索引失效，从而导致全表扫描，影响查询性能。",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeIndexInvalidation,
		},
		AllowOffline: false,
		Message:      "禁止对索引列进行数学运算和使用函数",
		Func:         checkMathComputationOrFuncOnIndex,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLDisableDropStatement,
			Desc:       "禁止除索引外的DROP操作",
			Annotation: "DROP是DDL，数据变更不会写入日志，无法进行回滚；建议开启此规则，避免误删除操作",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeUsageSuggestion,
		},
		Message:      "禁止除索引外的DROP操作",
		AllowOffline: true,
		Func:         disableDropStmt,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLCheckScanRows,
			Desc:       "扫描行数超过阈值，筛选条件必须带上主键或者索引",
			Annotation: "筛选条件必须带上主键或索引可降低数据库查询的时间复杂度，提高查询效率。",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   DefaultSingleParamKeyName,
					Value: "10",
					Desc:  "扫描行数量（万）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		AllowOffline: false,
		Message:      "扫描行数超过阈值，筛选条件必须带上主键或者索引",
		Func:         checkScanRows,
	},
	{
		Rule: driverV2.Rule{
			Name:       DMLCheckWhereIsInvalid,
			Desc:       "禁止使用没有WHERE条件或者WHERE条件恒为TRUE的SQL",
			Annotation: "SQL缺少WHERE条件在执行时会进行全表扫描产生额外开销，建议在大数据量高并发环境下开启，避免影响数据库查询性能",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeDMLConvention,
		},
		Message:      "禁止使用没有WHERE条件或者WHERE条件恒为TRUE的SQL",
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: driverV2.Rule{
			Name:       DDLCheckColumnWithoutDefault,
			Desc:       "除了自增列及大字段列之外，每个列都必须添加默认值",
			Annotation: "列添加默认值，可避免列为NULL值时对查询的影响",
			Level:      driverV2.RuleLevelError,
			Category:   RuleTypeDDLConvention,
		},
		Message:      "除了自增列及大字段列之外，每个列都必须添加默认值",
		AllowOffline: true,
		Func:         checkColumnWithoutDefault,
	},
}
