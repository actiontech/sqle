package optimization

import (
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
)

var OptimizationRuleMap map[string][]OptimizationRuleHandler // ruleCode与plugin重写规则的映射关系

// DML规则
const (
	DMLRuleNPERewrite                        = "dml_rule_npe_rewrite"
	DMLRuleAllSubqueryRewrite                = "dml_rule_all_subquery_rewrite"
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

// SQL优化规则的ruleCode
const (
	RuleAddOrderByNullRewrite             = "RuleAddOrderByNullRewrite"
	RuleAllQualifierSubQueryRewrite       = "RuleAllQualifierSubQueryRewrite"
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
	RuleNPERewrite                        = "RuleNPERewrite"
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
	RuleCode string // SQL优化规则的ruleCode
}

func init() {
	OptimizationRuleMap = make(map[string][]OptimizationRuleHandler)
	OptimizationRuleMap["MySQL"] = MySQLOptimizationRuleHandler
	OptimizationRuleMap["Oracle"] = OracleOptimizationRuleHandler

	// SQL优化规则知识库
	defaultRulesKnowledge, err := getDefaultRulesKnowledge()
	if err != nil {
		log.NewEntry().Errorf("get default rules knowledge failed: %v", err)
		return
	}
	for _, optimizationRule := range OptimizationRuleMap {
		for i, rule := range optimizationRule {
			if knowledge, ok := defaultRulesKnowledge[rule.RuleCode]; ok {
				rule.Rule.Knowledge = driverV2.RuleKnowledge{Content: knowledge}
				optimizationRule[i] = rule
			}
		}
	}
}

// GetOptimizationRuleByRuleCode 通过pawsql的ruleCode和dbType获取重写规则
func GetOptimizationRuleByRuleCode(ruleCode string, dbType string) (*driverV2.Rule, bool) {
	rules := OptimizationRuleMap[dbType]
	if len(rules) > 0 {
		for _, rule := range rules {
			if rule.RuleCode == ruleCode {
				return &rule.Rule, true
			}
		}
	}
	return nil, false
}

// CanOptimizeDbType SQL优化是否支持该数据源类型
func CanOptimizeDbType(dt string) bool {
	_, exist := OptimizationRuleMap[dt]
	return exist
}
