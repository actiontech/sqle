package workflowexport

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
)

func TestMaxExportWorkflowCount(t *testing.T) {
	if MaxExportWorkflowCount != 10000 {
		t.Fatalf("MaxExportWorkflowCount = %d, want 10000", MaxExportWorkflowCount)
	}
}

func TestExceedLimitError(t *testing.T) {
	err := ExceedLimitError(10001)
	if err == nil {
		t.Fatal("expected error")
	}
	ce, ok := err.(*errors.CodeError)
	if !ok {
		t.Fatalf("expected CodeError, got %T", err)
	}
	if ce.Code() != int(errors.DataInvalid) {
		t.Fatalf("code = %d, want %d", ce.Code(), errors.DataInvalid)
	}
	msg := err.Error()
	if !strings.Contains(msg, "10000") || !strings.Contains(msg, "10001") || !strings.Contains(msg, "收窄筛选") {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestNormalizeMaxAuditNodes(t *testing.T) {
	if normalizeMaxAuditNodes(0) != 1 {
		t.Fatalf("want 1 for 0")
	}
	if normalizeMaxAuditNodes(2) != 2 {
		t.Fatalf("want 2")
	}
}

func TestBuildHeaderForLayoutGlobalSQLRelease(t *testing.T) {
	ctx := context.Background()
	header := BuildHeaderForLayout(ctx, LayoutGlobalSQLRelease, 2)
	joined := strings.Join(header, "|")
	if strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOperator)) {
		t.Fatalf("global sql_release must not contain operator: %v", header)
	}
	if strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionTime)) {
		t.Fatalf("global sql_release must not contain execution time: %v", header)
	}
	if header[0] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportProjectName) {
		t.Fatalf("first col want project name, got %s", header[0])
	}
	if header[len(header)-1] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContent) {
		t.Fatalf("last col want SQL content, got %s", header[len(header)-1])
	}
	auditor1 := fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditorTpl), 1)
	auditor2 := fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditorTpl), 2)
	if !strings.Contains(joined, auditor1) || !strings.Contains(joined, auditor2) {
		t.Fatalf("want dynamic nodes 1 and 2 in %v", header)
	}
	if strings.Contains(joined, fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditorTpl), 3)) {
		t.Fatalf("must not contain node 3 when maxN=2: %v", header)
	}
}

func TestBuildHeaderForLayoutProjectFrozen(t *testing.T) {
	ctx := context.Background()
	header := BuildHeader(ctx, false)
	joined := strings.Join(header, "|")
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOperator)) {
		t.Fatalf("project layout must keep operator")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionTime)) {
		t.Fatalf("project layout must keep execution time")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNode4Auditor)) {
		t.Fatalf("project layout must keep fixed 4 nodes")
	}
	if header[len(header)-1] == locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContent) {
		t.Fatalf("project layout SQL content should not be last column")
	}
}

func TestBuildGlobalCommonHeader(t *testing.T) {
	ctx := context.Background()
	header := BuildHeaderForLayout(ctx, LayoutGlobalCommon, 0)
	if len(header) != 8 {
		t.Fatalf("common header len=%d want 8: %v", len(header), header)
	}
	if header[1] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowType) {
		t.Fatalf("col2 want workflow type, got %s", header[1])
	}
}

func TestBuildGlobalDataExportHeader(t *testing.T) {
	ctx := context.Background()
	header := BuildHeaderForLayout(ctx, LayoutGlobalDataExport, 1)
	joined := strings.Join(header, "|")
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportExecTime)) {
		t.Fatalf("want export exec time col")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportResult)) {
		t.Fatalf("want export result col")
	}
}

func TestMapSQLWorkflowStatusToUnified(t *testing.T) {
	cases := []struct {
		native string
		want   string
	}{
		{"wait_for_audit", "pending_approval"},
		{"wait_for_execution", "pending_action"},
		{"rejected", "rejected"},
		{"canceled", "cancelled"},
		{"executing", "in_progress"},
		{"exec_failed", "failed"},
		{"finished", "completed"},
		{"", "unknown"},
	}
	for _, c := range cases {
		if got := mapSQLWorkflowStatus(c.native); got != c.want {
			t.Fatalf("mapSQLWorkflowStatus(%q)=%q want %q", c.native, got, c.want)
		}
	}
}

func TestCommonColumnStatusVocabularyUnified(t *testing.T) {
	ctx := context.Background()
	nativeWaitAudit := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WorkflowStatusWaitForAudit)
	unifiedPendingApproval := LocalizeUnifiedStatus(ctx, "pending_approval")
	if nativeWaitAudit == unifiedPendingApproval {
		t.Fatalf("fixture broken: native and unified labels should differ (%q)", nativeWaitAudit)
	}

	sqlReleaseStatus := LocalizeUnifiedStatus(ctx, mapSQLWorkflowStatus("wait_for_audit"))
	dataExportStatus := LocalizeUnifiedStatus(ctx, mapDataExportStatus("wait_for_approve"))
	if sqlReleaseStatus != unifiedPendingApproval {
		t.Fatalf("SQL release common status = %q, want unified %q (not native %q)",
			sqlReleaseStatus, unifiedPendingApproval, nativeWaitAudit)
	}
	if dataExportStatus != unifiedPendingApproval {
		t.Fatalf("data export common status = %q, want same unified %q", dataExportStatus, unifiedPendingApproval)
	}
	if sqlReleaseStatus != dataExportStatus {
		t.Fatalf("same-column vocabulary split: sql=%q data_export=%q", sqlReleaseStatus, dataExportStatus)
	}

	_, rows := BuildGlobalCommonExport(ctx, []CommonExportRow{
		{ProjectName: "p", WorkflowType: "SQL", WorkflowID: "1", Status: sqlReleaseStatus},
		{ProjectName: "p", WorkflowType: "DE", WorkflowID: "2", Status: dataExportStatus},
	})
	statusCol := 6
	seen := map[string]bool{}
	for _, row := range rows {
		seen[row[statusCol]] = true
		if row[statusCol] == nativeWaitAudit {
			t.Fatalf("common export must not emit native wait-for-audit label %q", nativeWaitAudit)
		}
	}
	if len(seen) != 1 || !seen[unifiedPendingApproval] {
		t.Fatalf("want single unified status vocab in status column, got %v", seen)
	}
}
