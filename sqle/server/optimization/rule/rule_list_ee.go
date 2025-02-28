//go:build enterprise
// +build enterprise

package optimization

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

type SourceOptimizationRuleHandler struct {
	Rule     rulepkg.SourceRule
	RuleCode string
}

var MySQLOptimizationRuleHandler = generateOptimizationRuleHandlers(mySQLOptimizationRuleHandlerSource, driverV2.DriverTypeMySQL)

var OracleOptimizationRuleHandler = generateOptimizationRuleHandlers(oracleOptimizationRuleHandlerSource, driverV2.DriverTypeOracle)

func generateOptimizationRuleHandlers(sources []SourceOptimizationRuleHandler, dbType string) []OptimizationRuleHandler {
	result := make([]OptimizationRuleHandler, len(sources))
	for k, v := range sources {
		result[k] = OptimizationRuleHandler{
			Rule:     *rulepkg.ConvertSourceRule(plocale.Bundle, &v.Rule, dbType),
			RuleCode: v.RuleCode,
		}
	}
	return result
}

var mySQLOptimizationRuleHandlerSource = []SourceOptimizationRuleHandler{
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLHintGroupByRequiresConditions,
			Desc:       plocale.OptDMLHintGroupByRequiresConditionsDesc,
			Annotation: plocale.OptDMLHintGroupByRequiresConditionsAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleAddOrderByNullRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLCheckWhereExistScalarSubquery,
			Desc:       plocale.OptDMLCheckWhereExistScalarSubqueryDesc,
			Annotation: plocale.OptDMLCheckWhereExistScalarSubqueryAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleCntGtThanZeroRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLHintUseTruncateInsteadOfDelete,
			Desc:       plocale.OptDMLHintUseTruncateInsteadOfDeleteDesc,
			Annotation: plocale.OptDMLHintUseTruncateInsteadOfDeleteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDelete2TruncateRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLCheckWhereExistImplicitConversion,
			Desc:       plocale.OptDMLCheckWhereExistImplicitConversionDesc,
			Annotation: plocale.OptDMLCheckWhereExistImplicitConversionAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDiffDataTypeInPredicateWrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLCheckMathComputationOrFuncOnIndex,
			Desc:       plocale.OptDMLCheckMathComputationOrFuncOnIndexDesc,
			Annotation: plocale.OptDMLCheckMathComputationOrFuncOnIndexAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeIndexInvalidation,
		},
		RuleCode: RuleFuncWithColumnInPredicate,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLNotRecommendHaving,
			Desc:       plocale.OptDMLNotRecommendHavingDesc,
			Annotation: plocale.OptDMLNotRecommendHavingAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleHavingCond2WhereCondRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLWhereExistNull,
			Desc:       plocale.OptDMLWhereExistNullDesc,
			Annotation: plocale.OptDMLWhereExistNullAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleUseEqual4NullRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLNotRecommendIn,
			Desc:       plocale.OptDMLNotRecommendInDesc,
			Annotation: plocale.OptDMLNotRecommendInAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleInSubqueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLHintInNullOnlyFalse,
			Desc:       plocale.OptDMLHintInNullOnlyFalseDesc,
			Annotation: plocale.OptDMLHintInNullOnlyFalseAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNotInNullableSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLNotRecommendNotWildcardLike,
			Desc:       plocale.OptDMLNotRecommendNotWildcardLikeDesc,
			Annotation: plocale.OptDMLNotRecommendNotWildcardLikeAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNoWildcardInPredicateLikeWarning,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLCheckNotEqualSymbol,
			Desc:       plocale.OptDMLCheckNotEqualSymbolDesc,
			Annotation: plocale.OptDMLCheckNotEqualSymbolAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleUseNonstandardNotEqualOperator,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       rulepkg.DMLCheckLimitOffsetNum,
			Desc:       plocale.OptDMLCheckLimitOffsetNumDesc,
			Annotation: plocale.OptDMLCheckLimitOffsetNumAnnotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleLargeOffset,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleNPERewrite,
			Desc:       plocale.OptDMLRuleNPERewriteDesc,
			Annotation: plocale.OptDMLRuleNPERewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNPERewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleAllSubqueryRewrite,
			Desc:       plocale.OptDMLRuleAllSubqueryRewriteDesc,
			Annotation: plocale.OptDMLRuleAllSubqueryRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleAllQualifierSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleDiffOrderingSpecTypeWarning,
			Desc:       plocale.OptDMLRuleDiffOrderingSpecTypeWarningDesc,
			Annotation: plocale.OptDMLRuleDiffOrderingSpecTypeWarningAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDiffOrderingSpecTypeWarning,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleDistinctEliminationRewrite,
			Desc:       plocale.OptDMLRuleDistinctEliminationRewriteDesc,
			Annotation: plocale.OptDMLRuleDistinctEliminationRewriteAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDistinctEliminationRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleExists2JoinRewrite,
			Desc:       plocale.OptDMLRuleExists2JoinRewriteDesc,
			Annotation: plocale.OptDMLRuleExists2JoinRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleExists2JoinRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleFilterPredicatePushDownRewrite,
			Desc:       plocale.OptDMLRuleFilterPredicatePushDownRewriteDesc,
			Annotation: plocale.OptDMLRuleFilterPredicatePushDownRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleFilterPredicatePushDownRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleGroupingFromDiffTablesRewrite,
			Desc:       plocale.OptDMLRuleGroupingFromDiffTablesRewriteDesc,
			Annotation: plocale.OptDMLRuleGroupingFromDiffTablesRewriteAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleGroupingFromDiffTablesRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleJoinEliminationRewrite,
			Desc:       plocale.OptDMLRuleJoinEliminationRewriteDesc,
			Annotation: plocale.OptDMLRuleJoinEliminationRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleJoinEliminationRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleLimitClausePushDownRewrite,
			Desc:       plocale.OptDMLRuleLimitClausePushDownRewriteDesc,
			Annotation: plocale.OptDMLRuleLimitClausePushDownRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
			Params: []*rulepkg.SourceParam{
				{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "1000",
					Desc:  plocale.OptDMLRuleLimitClausePushDownRewriteParams1,
					Type:  params.ParamTypeInt,
				},
			},
		},
		RuleCode: RuleLimitClausePushDownRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleMaxMinAggRewrite,
			Desc:       plocale.OptDMLRuleMaxMinAggRewriteDesc,
			Annotation: plocale.OptDMLRuleMaxMinAggRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleMaxMinAggRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleMoveOrder2LeadingRewrite,
			Desc:       plocale.OptDMLRuleMoveOrder2LeadingRewriteDesc,
			Annotation: plocale.OptDMLRuleMoveOrder2LeadingRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleMoveOrder2LeadingRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleOrCond4SelectRewrite,
			Desc:       plocale.OptDMLRuleOrCond4SelectRewriteDesc,
			Annotation: plocale.OptDMLRuleOrCond4SelectRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrCond4SelectRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleOrCond4UpDeleteRewrite,
			Desc:       plocale.OptDMLRuleOrCond4UpDeleteRewriteDesc,
			Annotation: plocale.OptDMLRuleOrCond4UpDeleteRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrCond4UpDeleteRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleOrderEliminationInSubqueryRewrite,
			Desc:       plocale.OptDMLRuleOrderEliminationInSubqueryRewriteDesc,
			Annotation: plocale.OptDMLRuleOrderEliminationInSubqueryRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrderEliminationInSubqueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleOrderingFromDiffTablesRewrite,
			Desc:       plocale.OptDMLRuleOrderingFromDiffTablesRewriteDesc,
			Annotation: plocale.OptDMLRuleOrderingFromDiffTablesRewriteAnnotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrderingFromDiffTablesRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleOuter2InnerConversionRewrite,
			Desc:       plocale.OptDMLRuleOuter2InnerConversionRewriteDesc,
			Annotation: plocale.OptDMLRuleOuter2InnerConversionRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOuter2InnerConversionRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleProjectionPushdownRewrite,
			Desc:       plocale.OptDMLRuleProjectionPushdownRewriteDesc,
			Annotation: plocale.OptDMLRuleProjectionPushdownRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleProjectionPushdownRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleQualifierSubQueryRewrite,
			Desc:       plocale.OptDMLRuleQualifierSubQueryRewriteDesc,
			Annotation: plocale.OptDMLRuleQualifierSubQueryRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleQualifierSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleQueryFoldingRewrite,
			Desc:       plocale.OptDMLRuleQueryFoldingRewriteDesc,
			Annotation: plocale.OptDMLRuleQueryFoldingRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleQueryFoldingRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       DMLRuleSATTCRewrite,
			Desc:       plocale.OptDMLRuleSATTCRewriteDesc,
			Annotation: plocale.OptDMLRuleSATTCRewriteAnnotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleSATTCRewrite,
	},
}

var oracleOptimizationRuleHandlerSource = []SourceOptimizationRuleHandler{
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_500",
			Desc:       plocale.OptOracle500Desc,
			Annotation: plocale.OptOracle500Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNPERewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_501",
			Desc:       plocale.OptOracle501Desc,
			Annotation: plocale.OptOracle501Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleAllQualifierSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_502",
			Desc:       plocale.OptOracle502Desc,
			Annotation: plocale.OptOracle502Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleCntGtThanZeroRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_503",
			Desc:       plocale.OptOracle503Desc,
			Annotation: plocale.OptOracle503Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDelete2TruncateRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_504",
			Desc:       plocale.OptOracle504Desc,
			Annotation: plocale.OptOracle504Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDiffDataTypeInPredicateWrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_505",
			Desc:       plocale.OptOracle505Desc,
			Annotation: plocale.OptOracle505Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDDLConvention,
		},
		RuleCode: RuleDiffOrderingSpecTypeWarning,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_506",
			Desc:       plocale.OptOracle506Desc,
			Annotation: plocale.OptOracle506Annotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeIndexInvalidation,
		},
		RuleCode: RuleFuncWithColumnInPredicate,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_507",
			Desc:       plocale.OptOracle507Desc,
			Annotation: plocale.OptOracle507Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleHavingCond2WhereCondRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_508",
			Desc:       plocale.OptOracle508Desc,
			Annotation: plocale.OptOracle508Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleUseEqual4NullRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_509",
			Desc:       plocale.OptOracle509Desc,
			Annotation: plocale.OptOracle509Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleInSubqueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_510",
			Desc:       plocale.OptOracle510Desc,
			Annotation: plocale.OptOracle510Annotation,
			Level:      driverV2.RuleLevelError,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNotInNullableSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_511",
			Desc:       plocale.OptOracle511Desc,
			Annotation: plocale.OptOracle511Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleNoWildcardInPredicateLikeWarning,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_512",
			Desc:       plocale.OptOracle512Desc,
			Annotation: plocale.OptOracle512Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleUseNonstandardNotEqualOperator,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_513",
			Desc:       plocale.OptOracle513Desc,
			Annotation: plocale.OptOracle513Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleDistinctEliminationRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_514",
			Desc:       plocale.OptOracle514Desc,
			Annotation: plocale.OptOracle514Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleExists2JoinRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_515",
			Desc:       plocale.OptOracle515Desc,
			Annotation: plocale.OptOracle515Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleFilterPredicatePushDownRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_516",
			Desc:       plocale.OptOracle516Desc,
			Annotation: plocale.OptOracle516Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleGroupingFromDiffTablesRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_517",
			Desc:       plocale.OptOracle517Desc,
			Annotation: plocale.OptOracle517Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleJoinEliminationRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_518",
			Desc:       plocale.OptOracle518Desc,
			Annotation: plocale.OptOracle518Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleMaxMinAggRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_519",
			Desc:       plocale.OptOracle519Desc,
			Annotation: plocale.OptOracle519Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleMoveOrder2LeadingRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_520",
			Desc:       plocale.OptOracle520Desc,
			Annotation: plocale.OptOracle520Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrCond4SelectRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_521",
			Desc:       plocale.OptOracle521Desc,
			Annotation: plocale.OptOracle521Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrCond4UpDeleteRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_522",
			Desc:       plocale.OptOracle522Desc,
			Annotation: plocale.OptOracle522Annotation,
			Level:      driverV2.RuleLevelWarn,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOrderingFromDiffTablesRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_523",
			Desc:       plocale.OptOracle523Desc,
			Annotation: plocale.OptOracle523Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleOuter2InnerConversionRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_524",
			Desc:       plocale.OptOracle524Desc,
			Annotation: plocale.OptOracle524Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleProjectionPushdownRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_525",
			Desc:       plocale.OptOracle525Desc,
			Annotation: plocale.OptOracle525Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleQualifierSubQueryRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_526",
			Desc:       plocale.OptOracle526Desc,
			Annotation: plocale.OptOracle526Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleQueryFoldingRewrite,
	},
	{
		Rule: rulepkg.SourceRule{
			Name:       "Oracle_527",
			Desc:       plocale.OptOracle527Desc,
			Annotation: plocale.OptOracle527Annotation,
			Level:      driverV2.RuleLevelNotice,
			Category:   plocale.RuleTypeDMLConvention,
		},
		RuleCode: RuleSATTCRewrite,
	},
}
