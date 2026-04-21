package v1

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/auditreport"
	"golang.org/x/text/language"
)

func TestToLevelCounts(t *testing.T) {
	testCases := map[string]struct {
		input       map[string]int
		wantLen     int
		wantFirst   string // expected first level (highest priority)
		wantLast    string // expected last level (lowest priority)
		description string
	}{
		"mixed levels in fixed order": {
			input: map[string]int{
				"normal": 5,
				"error":  2,
				"warn":   3,
				"notice": 1,
			},
			wantLen:     4,
			wantFirst:   "normal",
			wantLast:    "error",
			description: "should list normal, notice, warn, error when all present",
		},
		"empty map returns empty slice": {
			input:       map[string]int{},
			wantLen:     0,
			description: "empty input should return empty result",
		},
		"single level": {
			input: map[string]int{
				"warn": 10,
			},
			wantLen:     1,
			wantFirst:   "warn",
			wantLast:    "warn",
			description: "single level should still work",
		},
		"only normal level": {
			input: map[string]int{
				"normal": 100,
			},
			wantLen:     1,
			wantFirst:   "normal",
			wantLast:    "normal",
			description: "all normal should return single entry",
		},
		"unknown level after standard levels": {
			input: map[string]int{
				"normal":  1,
				"error":   1,
				"unknown": 1,
			},
			wantLen:     3,
			wantFirst:   "normal",
			wantLast:    "unknown",
			description: "unknown levels should follow standard levels, sorted by name",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := toLevelCounts(tc.input)

			if len(result) != tc.wantLen {
				t.Fatalf("toLevelCounts() returned %d items, want %d", len(result), tc.wantLen)
			}

			if tc.wantLen == 0 {
				return
			}

			// Verify first (highest priority) element
			if result[0].Level != tc.wantFirst {
				t.Errorf("first level = %q, want %q (%s)", result[0].Level, tc.wantFirst, tc.description)
			}

			// Verify last (lowest priority) element
			if result[len(result)-1].Level != tc.wantLast {
				t.Errorf("last level = %q, want %q (%s)", result[len(result)-1].Level, tc.wantLast, tc.description)
			}

			// Verify counts match input
			for _, lc := range result {
				expectedCount, ok := tc.input[lc.Level]
				if !ok {
					t.Errorf("unexpected level %q in result", lc.Level)
					continue
				}
				if lc.Count != expectedCount {
					t.Errorf("level %q count = %d, want %d", lc.Level, lc.Count, expectedCount)
				}
			}
		})
	}
}

func TestToRuleHits(t *testing.T) {
	testCases := map[string]struct {
		input       map[string]int
		wantLen     int
		wantFirst   string // expected first rule (highest hit count)
		description string
	}{
		"multiple rules sorted by hit count descending": {
			input: map[string]int{
				"no_select_all": 5,
				"no_drop_table": 10,
				"add_index":     3,
			},
			wantLen:     3,
			wantFirst:   "no_drop_table",
			description: "should sort by hit count descending",
		},
		"empty map returns empty slice": {
			input:       map[string]int{},
			wantLen:     0,
			description: "empty input should return empty result",
		},
		"single rule": {
			input: map[string]int{
				"no_select_all": 7,
			},
			wantLen:     1,
			wantFirst:   "no_select_all",
			description: "single rule should work correctly",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := toRuleHits(tc.input)

			if len(result) != tc.wantLen {
				t.Fatalf("toRuleHits() returned %d items, want %d", len(result), tc.wantLen)
			}

			if tc.wantLen == 0 {
				return
			}

			// Verify first element has highest hit count
			if result[0].RuleName != tc.wantFirst {
				t.Errorf("first rule = %q, want %q (%s)", result[0].RuleName, tc.wantFirst, tc.description)
			}

			// Verify descending order
			for i := 1; i < len(result); i++ {
				if result[i].HitCount > result[i-1].HitCount {
					t.Errorf("rule at index %d (count=%d) > rule at index %d (count=%d), expected descending order",
						i, result[i].HitCount, i-1, result[i-1].HitCount)
				}
			}

			// Verify counts match input
			for _, rh := range result {
				expectedCount, ok := tc.input[rh.RuleName]
				if !ok {
					t.Errorf("unexpected rule %q in result", rh.RuleName)
					continue
				}
				if rh.HitCount != expectedCount {
					t.Errorf("rule %q hit count = %d, want %d", rh.RuleName, rh.HitCount, expectedCount)
				}
			}
		})
	}
}

func TestExtractRuleInfo(t *testing.T) {
	ctx := context.Background()

	testCases := map[string]struct {
		auditResults    model.AuditResults
		wantRuleName    string
		wantHasRuleName bool
		wantHasSugg     bool
		description     string
	}{
		"empty audit results": {
			auditResults:    model.AuditResults{},
			wantRuleName:    "",
			wantHasRuleName: false,
			wantHasSugg:     false,
			description:     "empty results should return empty strings",
		},
		"single rule hit": {
			auditResults: model.AuditResults{
				{
					Level:    "warn",
					RuleName: "no_select_all",
					I18nAuditResultInfo: model.I18nAuditResultInfo{
						language.Chinese: model.AuditResultInfo{Message: "should not use SELECT *"},
					},
				},
			},
			wantRuleName:    "no_select_all",
			wantHasRuleName: true,
			wantHasSugg:     true,
			description:     "single rule should return its name and message",
		},
		"multiple rule hits": {
			auditResults: model.AuditResults{
				{
					Level:    "warn",
					RuleName: "no_select_all",
					I18nAuditResultInfo: model.I18nAuditResultInfo{
						language.Chinese: model.AuditResultInfo{Message: "avoid SELECT *"},
					},
				},
				{
					Level:    "error",
					RuleName: "no_drop_table",
					I18nAuditResultInfo: model.I18nAuditResultInfo{
						language.Chinese: model.AuditResultInfo{Message: "DROP TABLE not allowed"},
					},
				},
			},
			wantRuleName:    "no_select_all, no_drop_table",
			wantHasRuleName: true,
			wantHasSugg:     true,
			description:     "multiple rules should be comma-separated",
		},
		"rule with empty name is skipped": {
			auditResults: model.AuditResults{
				{
					Level:    "notice",
					RuleName: "",
					I18nAuditResultInfo: model.I18nAuditResultInfo{
						language.Chinese: model.AuditResultInfo{Message: "some notice"},
					},
				},
			},
			wantRuleName:    "",
			wantHasRuleName: false,
			wantHasSugg:     true,
			description:     "rule with empty name should be skipped in rule names but message still included",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ruleName, suggestion := extractRuleInfo(tc.auditResults, ctx)

			if tc.wantHasRuleName {
				if ruleName != tc.wantRuleName {
					t.Errorf("ruleName = %q, want %q (%s)", ruleName, tc.wantRuleName, tc.description)
				}
			} else {
				if ruleName != "" {
					t.Errorf("ruleName = %q, want empty (%s)", ruleName, tc.description)
				}
			}

			if tc.wantHasSugg {
				if suggestion == "" {
					t.Errorf("suggestion should not be empty (%s)", tc.description)
				}
			} else {
				if suggestion != "" {
					t.Errorf("suggestion = %q, want empty (%s)", suggestion, tc.description)
				}
			}
		})
	}
}

func TestBuildReportLabels(t *testing.T) {
	ctx := context.Background()

	testCases := map[string]struct {
		description string
	}{
		"default context returns non-empty labels": {
			description: "all label fields should be non-empty with default locale",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			labels := buildReportLabels(ctx)

			// Verify all label fields are non-empty
			fieldChecks := map[string]string{
				"AuditSummary":      labels.AuditSummary,
				"ResultStatistics":  labels.ResultStatistics,
				"ProblemSQLList":    labels.ProblemSQLList,
				"RuleHitStatistics": labels.RuleHitStatistics,
				"AuditTime":         labels.AuditTime,
				"DataSource":        labels.DataSource,
				"Schema":            labels.Schema,
				"TotalSQL":          labels.TotalSQL,
				"PassRate":          labels.PassRate,
				"Score":             labels.Score,
				"AuditLevel":        labels.AuditLevel,
				"Number":            labels.Number,
				"SQL":               labels.SQL,
				"AuditStatus":       labels.AuditStatus,
				"AuditResult":       labels.AuditResult,
				"ExecStatus":        labels.ExecStatus,
				"ExecResult":        labels.ExecResult,
				"RollbackSQL":       labels.RollbackSQL,
				"RuleName":          labels.RuleName,
				"Description":       labels.Description,
				"Suggestion":        labels.Suggestion,
				"Count":             labels.Count,
				"HitCount":          labels.HitCount,
			}

			for field, value := range fieldChecks {
				if value == "" {
					t.Errorf("%s: label field %q is empty (%s)", tc.description, field, tc.description)
				}
			}
		})
	}
}

// TestLevelCountsPreserveAllLevels verifies that toLevelCounts preserves
// the count data correctly for all four standard audit levels.
func TestLevelCountsPreserveAllLevels(t *testing.T) {
	input := map[string]int{
		"error":  3,
		"warn":   5,
		"notice": 2,
		"normal": 10,
	}

	result := toLevelCounts(input)

	if len(result) != 4 {
		t.Fatalf("expected 4 levels, got %d", len(result))
	}

	// Build a lookup map from result
	resultMap := make(map[string]int)
	for _, lc := range result {
		resultMap[lc.Level] = lc.Count
	}

	// Verify each level count matches
	for level, expectedCount := range input {
		if count, ok := resultMap[level]; !ok {
			t.Errorf("level %q missing from result", level)
		} else if count != expectedCount {
			t.Errorf("level %q: count = %d, want %d", level, count, expectedCount)
		}
	}

	// Verify ordering: normal, notice, warn, error
	expectedOrder := []string{"normal", "notice", "warn", "error"}
	for i, expected := range expectedOrder {
		if result[i].Level != expected {
			t.Errorf("position %d: level = %q, want %q", i, result[i].Level, expected)
		}
	}
}

// TestToRuleHitsStableSortForEqualCounts verifies that toRuleHits handles
// rules with equal hit counts without error.
func TestToRuleHitsStableSortForEqualCounts(t *testing.T) {
	input := map[string]int{
		"rule_a": 5,
		"rule_b": 5,
		"rule_c": 5,
	}

	result := toRuleHits(input)

	if len(result) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(result))
	}

	// All should have count 5
	for _, rh := range result {
		if rh.HitCount != 5 {
			t.Errorf("rule %q: hit count = %d, want 5", rh.RuleName, rh.HitCount)
		}
	}
}

// TestToLevelCountsNilMap verifies toLevelCounts handles nil map gracefully.
func TestToLevelCountsNilMap(t *testing.T) {
	result := toLevelCounts(nil)
	if result == nil {
		t.Fatal("toLevelCounts(nil) should return empty slice, not nil")
	}
	if len(result) != 0 {
		t.Errorf("toLevelCounts(nil) returned %d items, want 0", len(result))
	}
}

// TestToRuleHitsNilMap verifies toRuleHits handles nil map gracefully.
func TestToRuleHitsNilMap(t *testing.T) {
	result := toRuleHits(nil)
	if result == nil {
		t.Fatal("toRuleHits(nil) should return empty slice, not nil")
	}
	if len(result) != 0 {
		t.Errorf("toRuleHits(nil) returned %d items, want 0", len(result))
	}
}

// TestExtractRuleInfoNilResults verifies extractRuleInfo handles nil AuditResults.
func TestExtractRuleInfoNilResults(t *testing.T) {
	ctx := context.Background()
	ruleName, suggestion := extractRuleInfo(nil, ctx)
	if ruleName != "" {
		t.Errorf("extractRuleInfo(nil) ruleName = %q, want empty", ruleName)
	}
	if suggestion != "" {
		t.Errorf("extractRuleInfo(nil) suggestion = %q, want empty", suggestion)
	}
}

// TestCSVHeaders verifies that CSVHeaders returns the correct number of columns
// based on the report labels.
func TestCSVHeaders(t *testing.T) {
	data := &auditreport.AuditReportData{
		Labels: auditreport.ReportLabels{
			Number:      "Number",
			SQL:         "SQL",
			AuditStatus: "Audit Status",
			AuditResult: "Audit Result",
			ExecStatus:  "Exec Status",
			ExecResult:  "Exec Result",
			RollbackSQL: "Rollback SQL",
			Description: "Description",
		},
	}

	gen := auditreport.NewCSVReportGenerator()
	headers := gen.CSVHeaders(data)
	if len(headers) != 8 {
		t.Errorf("CSVHeaders() returned %d columns, want 8", len(headers))
	}

	expectedHeaders := []string{"Number", "SQL", "Audit Status", "Audit Result", "Exec Status", "Exec Result", "Rollback SQL", "Description"}
	for i, h := range headers {
		if h != expectedHeaders[i] {
			t.Errorf("CSVHeaders()[%d] = %q, want %q", i, h, expectedHeaders[i])
		}
	}
}
