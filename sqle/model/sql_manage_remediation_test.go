package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRemediationResultRuleDiff(t *testing.T) {
	result := CalculateRemediationResult(
		AuditResults{
			{RuleName: "rule_resolved", Level: "warn"},
			{RuleName: "rule_unchanged", Level: "warn"},
		},
		AuditResults{
			{RuleName: "rule_unchanged", Level: "error"},
			{RuleName: "rule_new", Level: "warn"},
		},
		false,
	)

	assert.Equal(t, RemediationStatusDeteriorated, result.Status)
	assert.Equal(t, []string{"rule_resolved"}, auditResultRuleNames(result.Resolved))
	assert.Equal(t, []string{"rule_new"}, auditResultRuleNames(result.New))
	assert.Equal(t, []string{"rule_unchanged"}, auditResultRuleNames(result.Unchanged))
}

func TestCalculateRemediationResultStatus(t *testing.T) {
	testCases := []struct {
		name   string
		first  AuditResults
		latest AuditResults
		status string
	}{
		{name: "resolved", first: auditResults("rule_a"), latest: nil, status: RemediationStatusResolved},
		{name: "partially fixed", first: auditResults("rule_a", "rule_b"), latest: auditResults("rule_b"), status: RemediationStatusPartiallyFixed},
		{name: "unchanged", first: auditResults("rule_a"), latest: auditResults("rule_a"), status: RemediationStatusUnchanged},
		{name: "deteriorated", first: auditResults("rule_a"), latest: auditResults("rule_a", "rule_b"), status: RemediationStatusDeteriorated},
		{name: "newly discovered", first: nil, latest: auditResults("rule_a"), status: RemediationStatusNewlyDiscovered},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := CalculateRemediationResult(testCase.first, testCase.latest, false)
			assert.Equal(t, testCase.status, result.Status)
		})
	}
}

func TestCalculateRemediationResultFirstAuditMissing(t *testing.T) {
	result := CalculateRemediationResult(nil, auditResults("rule_a"), true)

	assert.True(t, result.FirstAuditMissing)
	assert.Equal(t, RemediationStatusNewlyDiscovered, result.Status)
	assert.Equal(t, []string{"rule_a"}, auditResultRuleNames(result.New))
}

func TestAuditResultsScanNull(t *testing.T) {
	var results AuditResults

	assert.NoError(t, results.Scan(nil))
	assert.Nil(t, results)
	assert.NoError(t, results.Scan([]byte("null")))
	assert.Nil(t, results)
}

func TestCalculateAuditResultsScore(t *testing.T) {
	assert.Equal(t, int32(100), CalculateAuditResultsScore([]AuditResults{nil, {}}))
	assert.Equal(t, int32(0), CalculateAuditResultsScore([]AuditResults{auditResultsWithLevel(auditLevelError)}))
	assert.Greater(t, CalculateAuditResultsScore([]AuditResults{auditResultsWithLevel(auditLevelWarn)}), CalculateAuditResultsScore([]AuditResults{auditResultsWithLevel(auditLevelError)}))
}

func TestCalculateSqlManageRemediationOverview(t *testing.T) {
	records := []*SQLManageRecord{
		{FirstAuditResults: auditResults("rule_a"), AuditResults: nil},
		{FirstAuditResults: auditResults("rule_b", "rule_c"), AuditResults: auditResultsPtr("rule_c")},
		{FirstAuditResults: auditResults("rule_d"), AuditResults: auditResultsPtr("rule_d", "rule_e")},
	}

	overview := CalculateSqlManageRemediationOverview("project-id", "plan-id", "default", records)

	assert.Equal(t, "project-id", overview.ProjectID)
	assert.Equal(t, "plan-id", overview.InstanceAuditPlanID)
	assert.Equal(t, "default", overview.AuditPlanType)
	assert.Equal(t, uint64(3), overview.SqlTotalNum)
	assert.Equal(t, uint64(1), overview.RemediationStatusCount.Resolved)
	assert.Equal(t, uint64(1), overview.RemediationStatusCount.PartiallyFixed)
	assert.Equal(t, uint64(1), overview.RemediationStatusCount.Deteriorated)
	assert.InDelta(t, 2.0/3.0, overview.RemediationRate, 0.001)
}

func auditResultsPtr(ruleNames ...string) *AuditResults {
	results := auditResults(ruleNames...)
	return &results
}

func auditResults(ruleNames ...string) AuditResults {
	results := make(AuditResults, 0, len(ruleNames))
	for _, ruleName := range ruleNames {
		results = append(results, AuditResult{RuleName: ruleName, Level: auditLevelWarn})
	}
	return results
}

func auditResultsWithLevel(level string) AuditResults {
	return AuditResults{{RuleName: level + "_rule", Level: level}}
}

func auditResultRuleNames(auditResults AuditResults) []string {
	ruleNames := make([]string, 0, len(auditResults))
	for _, auditResult := range auditResults {
		ruleNames = append(ruleNames, auditResult.RuleName)
	}
	return ruleNames
}
