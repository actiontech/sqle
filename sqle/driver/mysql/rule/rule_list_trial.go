//go:build trial
// +build trial

package rule

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

var sourceRuleHandlers = []*SourceHandler{
	{
		Rule: SourceRule{
			Name:       DMLCheckFuzzySearch,
			Desc:       plocale.DMLCheckFuzzySearchDesc,
			Annotation: plocale.DMLCheckFuzzySearchAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
		},
		Message:      plocale.DMLCheckFuzzySearchMessage,
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: SourceRule{
			Name:       DMLCheckJoinFieldType,
			Desc:       plocale.DMLCheckJoinFieldTypeDesc,
			Annotation: plocale.DMLCheckJoinFieldTypeAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		Message:      plocale.DMLCheckJoinFieldTypeMessage,
		AllowOffline: false,
		Func:         checkJoinFieldType,
	},
	{
		Rule: SourceRule{
			Name:       DDLRecommendTableColumnCharsetSame,
			Desc:       plocale.DDLRecommendTableColumnCharsetSameDesc,
			Annotation: plocale.DDLRecommendTableColumnCharsetSameAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDDLConvention,
		},
		Message: plocale.DDLRecommendTableColumnCharsetSameMessage,
		Func:    recommendTableColumnCharsetSame,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckColumnTimestampWithoutDefault,
			Desc:       plocale.DDLCheckColumnTimestampWithoutDefaultDesc,
			Annotation: plocale.DDLCheckColumnTimestampWithoutDefaultAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDDLConvention,
		},
		Message:      plocale.DDLCheckColumnTimestampWithoutDefaultMessage,
		AllowOffline: true,
		Func:         checkColumnTimestampWithoutDefault,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckIndexPrefix,
			Desc:       plocale.DDLCheckIndexPrefixDesc,
			Annotation: plocale.DDLCheckIndexPrefixAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeNamingConvention,
			Params: []*SourceParam{
				{
					Key:   DefaultSingleParamKeyName,
					Value: "idx_",
					Desc:  plocale.DDLCheckUniqueIndexPrefixParams1,
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      plocale.DDLCheckIndexPrefixMessage,
		AllowOffline: true,
		Func:         checkIndexPrefix,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckPKNotExist,
			Desc:       plocale.DDLCheckPKNotExistDesc,
			Annotation: plocale.DDLCheckPKNotExistAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeIndexingConvention,
		},
		Message:                         plocale.DDLCheckPKNotExistMessage,
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                            checkPrimaryKey,
	},
	{
		Rule: SourceRule{
			Name:       ConfigDMLExplainPreCheckEnable,
			Desc:       plocale.ConfigDMLExplainPreCheckEnableDesc,
			Annotation: plocale.ConfigDMLExplainPreCheckEnableAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeGlobalConfig,
		},
		Func: nil,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckIndexCount,
			Desc:       plocale.DDLCheckIndexCountDesc,
			Annotation: plocale.DDLCheckIndexCountAnnotation,
			Level:      driverV2.RuleLevelNotice,
			//Value:    "5",
			Category: plocale.RuleTypeIndexingConvention,
			Params: []*SourceParam{
				{
					Key:   DefaultSingleParamKeyName,
					Value: "5",
					Desc:  plocale.DDLCheckIndexCountParams1,
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:                         plocale.DDLCheckIndexCountMessage,
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                            checkIndex,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckPKWithoutAutoIncrement,
			Desc:       plocale.DDLCheckPKWithoutAutoIncrementDesc,
			Annotation: plocale.DDLCheckPKWithoutAutoIncrementAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeIndexingConvention,
		},
		Message:                         plocale.DDLCheckPKWithoutAutoIncrementMessage,
		AllowOffline:                    true,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                            checkPrimaryKey,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckObjectNameUsingKeyword,
			Desc:       plocale.DDLCheckObjectNameUsingKeywordDesc,
			Annotation: plocale.DDLCheckObjectNameUsingKeywordAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeNamingConvention,
		},
		Message:      plocale.DDLCheckObjectNameUsingKeywordMessage,
		AllowOffline: true,
		Func:         checkNewObjectName,
	},
	{
		Rule: SourceRule{
			Name:       DMLCheckMathComputationOrFuncOnIndex,
			Desc:       plocale.DMLCheckMathComputationOrFuncOnIndexDesc,
			Annotation: plocale.DMLCheckMathComputationOrFuncOnIndexAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeIndexInvalidation,
		},
		AllowOffline: false,
		Message:      plocale.DMLCheckMathComputationOrFuncOnIndexMessage,
		Func:         checkMathComputationOrFuncOnIndex,
	},
	{
		Rule: SourceRule{
			Name:       DDLDisableDropStatement,
			Desc:       plocale.DDLDisableDropStatementDesc,
			Annotation: plocale.DDLDisableDropStatementAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeUsageSuggestion,
		},
		Message:      plocale.DDLDisableDropStatementMessage,
		AllowOffline: true,
		Func:         disableDropStmt,
	},
	{
		Rule: SourceRule{
			Name:       DMLCheckScanRows,
			Desc:       plocale.DMLCheckScanRowsDesc,
			Annotation: plocale.DMLCheckScanRowsAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
			Params: []*SourceParam{
				{
					Key:   DefaultSingleParamKeyName,
					Value: "10",
					Desc:  plocale.DMLCheckScanRowsParams1,
					Type:  params.ParamTypeInt,
				},
			},
		},
		AllowOffline: false,
		Message:      plocale.DMLCheckScanRowsMessage,
		Func:         checkScanRows,
	},
	{
		Rule: SourceRule{
			Name:       DMLCheckWhereIsInvalid,
			Desc:       plocale.DMLCheckWhereIsInvalidDesc,
			Annotation: plocale.DMLCheckWhereIsInvalidAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
		},
		Message:      plocale.DMLCheckWhereIsInvalidMessage,
		AllowOffline: true,
		Func:         checkSelectWhere,
	},
	{
		Rule: SourceRule{
			Name:       DDLCheckColumnWithoutDefault,
			Desc:       plocale.DDLCheckColumnWithoutDefaultDesc,
			Annotation: plocale.DDLCheckColumnWithoutDefaultAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDDLConvention,
		},
		Message:      plocale.DDLCheckColumnWithoutDefaultMessage,
		AllowOffline: true,
		Func:         checkColumnWithoutDefault,
	},
}
