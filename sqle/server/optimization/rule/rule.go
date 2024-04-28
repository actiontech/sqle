package optimization

import (
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var OptimizationRuleMap map[string][]OptimizationRuleHandler // ruleCode与plugin的rule name映射关系

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

// mysql sql优化的ruleCode
const (
	RuleAddOrderByNullRewrite             = "RuleAddOrderByNullRewrite"
	RuleCntGtThanZeroRewrite              = "RuleCntGtThanZeroRewrite"
	RuleDelete2TruncateRewrite            = "RuleDelete2TruncateRewrite"
	RuleDiffDataTypeInPredicateWrite      = "RuleDiffDataTypeInPredicateWrite"
	RuleDiffOrderingSpecTypeWarning       = "RuleDiffOrderingSpecTypeWarning"
	RuleFuncWithColumnInPredicate         = "RuleFuncWithColumnInPredicate"
	RuleHavingCond2WhereCondRewrite       = "RuleHavingCond2WhereCondRewrite"
	RuleUseEqual4NullRewrite              = "RuleUseEqual4NullRewrite"
	RuleInSubqueryRewrite                 = "RuleInSubqueryRewrite"
	RuleNotInNullableSubQueryRewrite      = "RuleNotInNullableSubQueryRewrite"
	RuleNoWildcardInPredicateLikeWarning  = "RuleNoWildcardInPredicateLikeWarning"
	RuleUseNonstandardNotEqualOperator    = "RuleUseNonstandardNotEqualOperator"
	RuleLargeOffset                       = "RuleLargeOffset"
	RuleDistinctEliminationRewrite        = "RuleDistinctEliminationRewrite"
	RuleExists2JoinRewrite                = "RuleExists2JoinRewrite"
	RuleFilterPredicatePushDownRewrite    = "RuleFilterPredicatePushDownRewrite"
	RuleGroupingFromDiffTablesRewrite     = "RuleGroupingFromDiffTablesRewrite"
	RuleJoinEliminationRewrite            = "RuleJoinEliminationRewrite"
	RuleLimitClausePushDownRewrite        = "RuleLimitClausePushDownRewrite"
	RuleMaxMinAggRewrite                  = "RuleMaxMinAggRewrite"
	RuleMoveOrder2LeadingRewrite          = "RuleMoveOrder2LeadingRewrite"
	RuleOrCond4SelectRewrite              = "RuleOrCond4SelectRewrite"
	RuleOrCond4UpDeleteRewrite            = "RuleOrCond4UpDeleteRewrite"
	RuleOrderEliminationInSubqueryRewrite = "RuleOrderEliminationInSubqueryRewrite"
	RuleOrderingFromDiffTablesRewrite     = "RuleOrderingFromDiffTablesRewrite"
	RuleOuter2InnerConversionRewrite      = "RuleOuter2InnerConversionRewrite"
	RuleProjectionPushdownRewrite         = "RuleProjectionPushdownRewrite"
	RuleQualifierSubQueryRewrite          = "RuleQualifierSubQueryRewrite"
	RuleQueryFoldingRewrite               = "RuleQueryFoldingRewrite"
	RuleSATTCRewrite                      = "RuleSATTCRewrite"
)

type OptimizationRuleHandler struct {
	Rule     driverV2.Rule
	RuleCode string // sql优化规则的ruleCode
}

func init() {
	OptimizationRuleMap = make(map[string][]OptimizationRuleHandler)
	OptimizationRuleMap["MySQL"] = BaseOptimizationRuleHandler
}

// 通过sql优化规则的ruleCode和dbType获取插件规则的name
func GetPluginRuleNameByOptimizationRule(ruleCode string, dbType string) (string, bool) {
	rules := OptimizationRuleMap[dbType]
	if len(rules) > 0 {
		for _, rule := range rules {
			if rule.RuleCode == ruleCode {
				return rule.Rule.Name, true
			}
		}
	}
	return "", false
}