package v2

import (
	"context"
	"testing"

	"github.com/actiontech/sqle/sqle/model"
)

func Test_convertTaskResultToAuditResV2_returnsSkippedAuditResult(t *testing.T) {
	task := &model.Task{
		ExecuteSQLs: []*model.ExecuteSQL{
			{
				BaseSQL: model.BaseSQL{Number: 1, Content: "select count(*) from t"},
				AuditResults: model.AuditResults{
					{Level: "warn", RuleName: "rule_kept"},
				},
				SkippedAuditResults: model.AuditResults{
					{Level: "error", RuleName: "rule_skipped"},
				},
				AuditLevel: "warn",
			},
		},
	}

	out := convertTaskResultToAuditResV2(context.Background(), task)
	if len(out.SQLResults) != 1 {
		t.Fatalf("unexpected sql result length: %d", len(out.SQLResults))
	}
	if len(out.SQLResults[0].SkippedAuditResult) != 1 {
		t.Fatalf("unexpected skipped result length: %d", len(out.SQLResults[0].SkippedAuditResult))
	}
	if out.SQLResults[0].SkippedAuditResult[0].RuleName != "rule_skipped" {
		t.Fatalf("unexpected skipped rule: %+v", out.SQLResults[0].SkippedAuditResult[0])
	}
}
