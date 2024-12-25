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
			Name:         DMLCheckFuzzySearch,
			Desc:         plocale.DMLCheckFuzzySearchDesc,
			Annotation:   plocale.DMLCheckFuzzySearchAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeDMLConvention,
			AllowOffline: true,
		},
		Message: plocale.DMLCheckFuzzySearchMessage,
		Func:    checkSelectWhere,
	},
	{
		Rule: SourceRule{
			Name:         DMLCheckJoinFieldType,
			Desc:         plocale.DMLCheckJoinFieldTypeDesc,
			Annotation:   plocale.DMLCheckJoinFieldTypeAnnotation,
			Level:        driverV2.RuleLevelWarn,
			Category:     plocale.RuleTypeDMLConvention,
			AllowOffline: false,
		},
		Message: plocale.DMLCheckJoinFieldTypeMessage,
		Func:    checkJoinFieldType,
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
			Name:         DDLCheckColumnTimestampWithoutDefault,
			Desc:         plocale.DDLCheckColumnTimestampWithoutDefaultDesc,
			Annotation:   plocale.DDLCheckColumnTimestampWithoutDefaultAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeDDLConvention,
			AllowOffline: true,
		},
		Message: plocale.DDLCheckColumnTimestampWithoutDefaultMessage,
		Func:    checkColumnTimestampWithoutDefault,
	},
	{
		Rule: SourceRule{
			Name:         DDLCheckIndexPrefix,
			Desc:         plocale.DDLCheckIndexPrefixDesc,
			Annotation:   plocale.DDLCheckIndexPrefixAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeNamingConvention,
			AllowOffline: true,
			Params: []*SourceParam{
				{
					Key:   DefaultSingleParamKeyName,
					Value: "idx_",
					Desc:  plocale.DDLCheckUniqueIndexPrefixParams1,
					Type:  params.ParamTypeString,
				},
			},
		},
		Message: plocale.DDLCheckIndexPrefixMessage,
		Func:    checkIndexPrefix,
	},
	{
		Rule: SourceRule{
			Name:         DDLCheckPKNotExist,
			Desc:         plocale.DDLCheckPKNotExistDesc,
			Annotation:   plocale.DDLCheckPKNotExistAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeIndexingConvention,
			AllowOffline: true,
		},
		Message:                         plocale.DDLCheckPKNotExistMessage,
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
			Category:     plocale.RuleTypeIndexingConvention,
			AllowOffline: true,
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
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}, &ast.CreateIndexStmt{}},
		Func:                            checkIndex,
	},
	{
		Rule: SourceRule{
			Name:         DDLCheckPKWithoutAutoIncrement,
			Desc:         plocale.DDLCheckPKWithoutAutoIncrementDesc,
			Annotation:   plocale.DDLCheckPKWithoutAutoIncrementAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeIndexingConvention,
			AllowOffline: true,
		},
		Message:                         plocale.DDLCheckPKWithoutAutoIncrementMessage,
		NotAllowOfflineStmts:            []ast.Node{&ast.AlterTableStmt{}},
		NotSupportExecutedSQLAuditStmts: []ast.Node{&ast.AlterTableStmt{}},
		Func:                            checkPrimaryKey,
	},
	{
		Rule: SourceRule{
			Name:         DDLCheckObjectNameUsingKeyword,
			Desc:         plocale.DDLCheckObjectNameUsingKeywordDesc,
			Annotation:   plocale.DDLCheckObjectNameUsingKeywordAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeNamingConvention,
			AllowOffline: true,
		},
		Message: plocale.DDLCheckObjectNameUsingKeywordMessage,
		Func:    checkNewObjectName,
	},
	{
		Rule: SourceRule{
			Name:         DMLCheckMathComputationOrFuncOnIndex,
			Desc:         plocale.DMLCheckMathComputationOrFuncOnIndexDesc,
			Annotation:   plocale.DMLCheckMathComputationOrFuncOnIndexAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeIndexInvalidation,
			AllowOffline: false,
		},
		Message: plocale.DMLCheckMathComputationOrFuncOnIndexMessage,
		Func:    checkMathComputationOrFuncOnIndex,
	},
	{
		Rule: SourceRule{
			Name:         DDLDisableDropStatement,
			Desc:         plocale.DDLDisableDropStatementDesc,
			Annotation:   plocale.DDLDisableDropStatementAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeUsageSuggestion,
			AllowOffline: true,
		},
		Message: plocale.DDLDisableDropStatementMessage,
		Func:    disableDropStmt,
	},
	{
		Rule: SourceRule{
			Name:         DMLCheckScanRows,
			Desc:         plocale.DMLCheckScanRowsDesc,
			Annotation:   plocale.DMLCheckScanRowsAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeDMLConvention,
			AllowOffline: false,
			Params: []*SourceParam{
				{
					Key:   DefaultSingleParamKeyName,
					Value: "10",
					Desc:  plocale.DMLCheckScanRowsParams1,
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: plocale.DMLCheckScanRowsMessage,
		Func:    checkScanRows,
	},
	{
		Rule: SourceRule{
			Name:         DMLCheckWhereIsInvalid,
			Desc:         plocale.DMLCheckWhereIsInvalidDesc,
			Annotation:   plocale.DMLCheckWhereIsInvalidAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeDMLConvention,
			AllowOffline: true,
		},
		Message: plocale.DMLCheckWhereIsInvalidMessage,
		Func:    checkSelectWhere,
	},
	{
		Rule: SourceRule{
			Name:         DDLCheckColumnWithoutDefault,
			Desc:         plocale.DDLCheckColumnWithoutDefaultDesc,
			Annotation:   plocale.DDLCheckColumnWithoutDefaultAnnotation,
			Level:        driverV2.RuleLevelError,
			Category:     plocale.RuleTypeDDLConvention,
			AllowOffline: true,
		},
		Message: plocale.DDLCheckColumnWithoutDefaultMessage,
		Func:    checkColumnWithoutDefault,
	},
}
