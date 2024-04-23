package optimization

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
)

var RuleMapping map[string]string // ruleCode与plugin的rule name映射关系

// DML规则
const (
	DMLRuleDistinctEliminationRewrite        = "dml_rule_distinct_elimination_rewrite"
	DMLRuleExists2JoinRewrite                = "dml_rule_exists_2_join_rewrite"
	DMLRuleFilterPredicatePushDownRewrite    = "dml_rule_filter_predicate_push_down_rewrite"
	DMLRuleGroupingFromDiffTablesRewrite     = "dml_rule_grouping_from_diff_tables_rewrite"
	DMLRuleJoinEliminationRewrite            = "dml_rule_join_elimination_rewrite"
	DMLRuleLimitClausePushDownRewrite        = "dml_rule_limit_clause_push_down_rewrite"
	DMLRuleMaxMinAggRewrite                  = "dml_rule_max_min_agg_rewrite"
	DMLRuleMoveOrder2LeadingRewrite          = "dml_rule_move_order_2_leading_rewrite"
	DMLRuleOrCond4SelectRewrite              = "dml_rule_or_cond_4_select_rewrite"
	DMLRuleOrCond4UpDeleteRewrite            = "dml_rule_or_cond_4_up_delete_rewrite"
	DMLRuleOrderEliminationInSubqueryRewrite = "dml_rule_order_elimination_in_subquery_rewrite"
	DMLRuleOrderingFromDiffTablesRewrite     = "dml_rule_ordering_from_diff_tables_rewrite"
	DMLRuleOuter2InnerConversionRewrite      = "dml_rule_outer_2_inner_conversion_rewrite"
	DMLRuleProjectionPushdownRewrite         = "dml_rule_projection_pushdown_rewrite"
	DMLRuleQualifierSubQueryRewrite          = "dml_rule_qualifier_sub_query_rewrite"
	DMLRuleQueryFoldingRewrite               = "dml_rule_query_folding_rewrite"
	DMLRuleSATTCRewrite                      = "dml_rule_sattc_rewrite"
)

func init() {
	RuleMapping = make(map[string]string)

	// 有审核能力的重写规则
	RuleMapping[rulepkg.DMLHintGroupByRequiresConditions] = "RuleAddOrderByNullRewrite"
	RuleMapping[rulepkg.DMLCheckWhereExistScalarSubquery] = "RuleCntGtThanZeroRewrite"
	RuleMapping[rulepkg.DMLHintUseTruncateInsteadOfDelete] = "RuleDelete2TruncateRewrite"
	RuleMapping[rulepkg.DMLCheckWhereExistImplicitConversion] = "RuleDiffDataTypeInPredicateWrite"
	RuleMapping[rulepkg.DDLCheckDatabaseCollation] = "RuleDiffOrderingSpecTypeWarning"
	RuleMapping[rulepkg.DMLCheckMathComputationOrFuncOnIndex] = "RuleFuncWithColumnInPredicate"
	RuleMapping[rulepkg.DMLNotRecommendHaving] = "RuleHavingCond2WhereCondRewrite"
	RuleMapping[rulepkg.DMLNotRecommendIn] = "RuleInSubqueryRewrite"
	RuleMapping[rulepkg.DMLCheckLimitOffsetNum] = "RuleLargeOffset"
	RuleMapping[rulepkg.DMLHintInNullOnlyFalse] = "RuleNotInNullableSubQueryRewrite"
	RuleMapping[rulepkg.DMLNotRecommendNotWildcardLike] = "RuleNoWildcardInPredicateLikeWarning"
	RuleMapping[rulepkg.DMLWhereExistNull] = "RuleUseEqual4NullRewrite"
	RuleMapping[rulepkg.DMLCheckNotEqualSymbol] = "RuleUseNonstandardNotEqualOperator"

	// 仅有重写能力的规则
	RuleMapping[DMLRuleDistinctEliminationRewrite] = "RuleDistinctEliminationRewrite"
	RuleMapping[DMLRuleExists2JoinRewrite] = "RuleExists2JoinRewrite"
	RuleMapping[DMLRuleFilterPredicatePushDownRewrite] = "RuleFilterPredicatePushDownRewrite"
	RuleMapping[DMLRuleGroupingFromDiffTablesRewrite] = "RuleGroupingFromDiffTablesRewrite"
	RuleMapping[DMLRuleJoinEliminationRewrite] = "RuleJoinEliminationRewrite"
	RuleMapping[DMLRuleLimitClausePushDownRewrite] = "RuleLimitClausePushDownRewrite"
	RuleMapping[DMLRuleMaxMinAggRewrite] = "RuleMaxMinAggRewrite"
	RuleMapping[DMLRuleMoveOrder2LeadingRewrite] = "RuleMoveOrder2LeadingRewrite"
	RuleMapping[DMLRuleOrCond4SelectRewrite] = "RuleOrCond4SelectRewrite"
	RuleMapping[DMLRuleOrCond4UpDeleteRewrite] = "RuleOrCond4UpDeleteRewrite"
	RuleMapping[DMLRuleOrderEliminationInSubqueryRewrite] = "RuleOrderEliminationInSubqueryRewrite"
	RuleMapping[DMLRuleOrderingFromDiffTablesRewrite] = "RuleOrderingFromDiffTablesRewrite"
	RuleMapping[DMLRuleOuter2InnerConversionRewrite] = "RuleOuter2InnerConversionRewrite"
	RuleMapping[DMLRuleProjectionPushdownRewrite] = "RuleProjectionPushdownRewrite"
	RuleMapping[DMLRuleQualifierSubQueryRewrite] = "RuleQualifierSubQueryRewrite"
	RuleMapping[DMLRuleQueryFoldingRewrite] = "RuleQueryFoldingRewrite"
	RuleMapping[DMLRuleSATTCRewrite] = "RuleSATTCRewrite"
}

// 通过重写规则的ruleCode获取插件规则的name
func GetPluginNameByRuleCode(ruleCode string) (string, bool) {
	for key, value := range RuleMapping {
		if value == ruleCode {
			return key, true
		}
	}
	return "", false
}
