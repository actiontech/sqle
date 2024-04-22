package optimization

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var ruleMapping map[string]string // ruleCode与plugin的rule name映射关系

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
	ruleMapping = make(map[string]string)

	// 有审核能力的重写规则
	ruleMapping[rulepkg.DMLHintGroupByRequiresConditions] = "RuleAddOrderByNullRewrite"
	ruleMapping[rulepkg.DMLCheckWhereExistScalarSubquery] = "RuleCntGtThanZeroRewrite"
	ruleMapping[rulepkg.DMLHintUseTruncateInsteadOfDelete] = "RuleDelete2TruncateRewrite"
	ruleMapping[rulepkg.DMLCheckWhereExistImplicitConversion] = "RuleDiffDataTypeInPredicateWrite"
	ruleMapping[rulepkg.DDLCheckDatabaseCollation] = "RuleDiffOrderingSpecTypeWarning"
	ruleMapping[rulepkg.DMLCheckMathComputationOrFuncOnIndex] = "RuleFuncWithColumnInPredicate"
	ruleMapping[rulepkg.DMLNotRecommendHaving] = "RuleHavingCond2WhereCondRewrite"
	ruleMapping[rulepkg.DMLNotRecommendIn] = "RuleInSubqueryRewrite"
	ruleMapping[rulepkg.DMLCheckLimitOffsetNum] = "RuleLargeOffset"
	ruleMapping[rulepkg.DMLHintInNullOnlyFalse] = "RuleNotInNullableSubQueryRewrite"
	ruleMapping[rulepkg.DMLNotRecommendNotWildcardLike] = "RuleNoWildcardInPredicateLikeWarning"
	ruleMapping[rulepkg.DMLWhereExistNull] = "RuleUseEqual4NullRewrite"
	ruleMapping[rulepkg.DMLCheckNotEqualSymbol] = "RuleUseNonstandardNotEqualOperator"

	// 仅有重写能力的规则
	ruleMapping[DMLRuleDistinctEliminationRewrite] = "RuleDistinctEliminationRewrite"
	ruleMapping[DMLRuleExists2JoinRewrite] = "RuleExists2JoinRewrite"
	ruleMapping[DMLRuleFilterPredicatePushDownRewrite] = "RuleFilterPredicatePushDownRewrite"
	ruleMapping[DMLRuleGroupingFromDiffTablesRewrite] = "RuleGroupingFromDiffTablesRewrite"
	ruleMapping[DMLRuleJoinEliminationRewrite] = "RuleJoinEliminationRewrite"
	ruleMapping[DMLRuleLimitClausePushDownRewrite] = "RuleLimitClausePushDownRewrite"
	ruleMapping[DMLRuleMaxMinAggRewrite] = "RuleMaxMinAggRewrite"
	ruleMapping[DMLRuleMoveOrder2LeadingRewrite] = "RuleMoveOrder2LeadingRewrite"
	ruleMapping[DMLRuleOrCond4SelectRewrite] = "RuleOrCond4SelectRewrite"
	ruleMapping[DMLRuleOrCond4UpDeleteRewrite] = "RuleOrCond4UpDeleteRewrite"
	ruleMapping[DMLRuleOrderEliminationInSubqueryRewrite] = "RuleOrderEliminationInSubqueryRewrite"
	ruleMapping[DMLRuleOrderingFromDiffTablesRewrite] = "RuleOrderingFromDiffTablesRewrite"
	ruleMapping[DMLRuleOuter2InnerConversionRewrite] = "RuleOuter2InnerConversionRewrite"
	ruleMapping[DMLRuleProjectionPushdownRewrite] = "RuleProjectionPushdownRewrite"
	ruleMapping[DMLRuleQualifierSubQueryRewrite] = "RuleQualifierSubQueryRewrite"
	ruleMapping[DMLRuleQueryFoldingRewrite] = "RuleQueryFoldingRewrite"
	ruleMapping[DMLRuleSATTCRewrite] = "RuleSATTCRewrite"
}

// 通过规则的ruleCode获取插件规则的name
func GetPluginNameByRuleCode(ruleCode string) (string, bool) {
	for key, value := range ruleMapping {
		if value == ruleCode {
			return key, true
		}
	}
	return "", false
}

// 整合规则，并赋予规则审核、重写能力
func MergeRulesAndPower(pluginRules []*driverV2.Rule) []*driverV2.Rule {
	allRules := []*driverV2.Rule{}
	rulesMap := make(map[string]bool)
	// 只有审核能力或审核、重写能力都有的规则
	for _, pluginRule := range pluginRules {
		pluginRule.AuditPower = "true"
		if _, ok := ruleMapping[pluginRule.Name]; ok {
			pluginRule.RewritePower = "true"
		} else {
			pluginRule.RewritePower = "false"
		}
		rulesMap[pluginRule.Name] = true
		allRules = append(allRules, pluginRule)
	}
	// 仅有重写能力的规则
	for _, rewriteRule := range RuleHandler {
		rewriteRule.Rule.RewritePower = "true"
		rewriteRule.Rule.AuditPower = "false"
		if _, exist := rulesMap[rewriteRule.Rule.Name]; !exist {
			rule := rewriteRule.Rule
			allRules = append(allRules, &rule)
		}
	}
	return allRules
}
