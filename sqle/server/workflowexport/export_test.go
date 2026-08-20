package workflowexport

import (
	"context"
	"regexp"
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
	if header[3] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription) {
		t.Fatalf("col4 want description, got %s", header[3])
	}
	if header[4] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType) {
		t.Fatalf("col5 want ops type after description, got %s", header[4])
	}
	if header[5] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource) {
		t.Fatalf("col6 want data source after ops type, got %s", header[5])
	}
	if header[len(header)-1] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContent) {
		t.Fatalf("last col want SQL content, got %s", header[len(header)-1])
	}
	auditCol := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportAuditRecord)
	if !strings.Contains(joined, auditCol) {
		t.Fatalf("want single audit record col in %v", header)
	}
	nodeRe := regexp.MustCompile(`\[节点\d+\]`)
	if nodeRe.MatchString(joined) {
		t.Fatalf("global sql_release must not contain [节点N] headers: %v", header)
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
	if header[2] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription) {
		t.Fatalf("col3 want description, got %s", header[2])
	}
	if header[3] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType) {
		t.Fatalf("col4 want ops type after description, got %s", header[3])
	}
}

func TestBuildGlobalCommonHeader(t *testing.T) {
	ctx := context.Background()
	header := BuildHeaderForLayout(ctx, LayoutGlobalCommon, 0)
	if len(header) != 11 {
		t.Fatalf("common header len=%d want 11: %v", len(header), header)
	}
	if header[1] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowType) {
		t.Fatalf("col2 want workflow type, got %s", header[1])
	}
	if header[4] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription) {
		t.Fatalf("col5 want description, got %s", header[4])
	}
	if header[7] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTaskOrderStatus) {
		t.Fatalf("col8 want status, got %s", header[7])
	}
	if header[8] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType) {
		t.Fatalf("col9 want ops type after status, got %s", header[8])
	}
	if header[9] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource) {
		t.Fatalf("col10 want data source after ops type, got %s", header[9])
	}
	if header[10] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContentPlain) {
		t.Fatalf("last col want SQL content plain, got %s", header[10])
	}
}

func TestBuildGlobalDataExportHeader(t *testing.T) {
	ctx := context.Background()
	header := BuildHeaderForLayout(ctx, LayoutGlobalDataExport, 1)
	joined := strings.Join(header, "|")
	if header[3] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription) {
		t.Fatalf("col4 want description, got %s", header[3])
	}
	if header[4] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType) {
		t.Fatalf("col5 want ops type after description, got %s", header[4])
	}
	if header[5] != locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource) {
		t.Fatalf("col6 want data source after ops type, got %s", header[5])
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportExecTime)) {
		t.Fatalf("want export exec time col")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportResult)) {
		t.Fatalf("want export result col")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContentPlain)) {
		t.Fatalf("want SQL content col")
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportAuditRecord)) {
		t.Fatalf("want audit record col")
	}
	nodeRe := regexp.MustCompile(`\[节点\d+\]`)
	if nodeRe.MatchString(joined) {
		t.Fatalf("data_export must not contain [节点N]: %v", header)
	}
}

func TestPruneEmptyColumns(t *testing.T) {
	header := []string{"a", "b", "c"}
	rows := [][]string{
		{"1", "", "x"},
		{"2", "", "y"},
	}
	h, r := pruneEmptyColumns(header, rows)
	if len(h) != 2 || h[0] != "a" || h[1] != "c" {
		t.Fatalf("header=%v", h)
	}
	if len(r) != 2 || r[0][0] != "1" || r[0][1] != "x" {
		t.Fatalf("rows=%v", r)
	}
	h0, r0 := pruneEmptyColumns(header, nil)
	if len(h0) != 3 || r0 != nil {
		t.Fatalf("zero rows should keep baseline header: %v %v", h0, r0)
	}
}

func TestFormatAuditRecordParts(t *testing.T) {
	got := FormatAuditRecordParts([]string{"张三 2026-08-01 10:00 通过", "  ", "李四 2026-08-01 11:20 驳回"})
	want := "张三 2026-08-01 10:00 通过；李四 2026-08-01 11:20 驳回"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBuildGlobalDataExportPrunesEmptyOptionalCols(t *testing.T) {
	ctx := context.Background()
	header, rows := BuildGlobalDataExportExport(ctx, []DataExportExportRecord{
		{
			ProjectName:    "p",
			WorkflowID:     "1",
			WorkflowName:   "n",
			Description:    "d",
			DBServiceNames: []string{"db"},
			CreatedAt:      "2026-01-01 00:00:00",
			CreatorName:    "u",
			UnifiedStatus:  "completed",
			SQLContent:     "SELECT 1",
			// OpsTypeName / AuditRecord / ExportExecTime / ExportResult empty → pruned
		},
	})
	joined := strings.Join(header, "|")
	if strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType)) {
		t.Fatalf("empty ops type col should be pruned: %v", header)
	}
	if strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportAuditRecord)) {
		t.Fatalf("empty audit col should be pruned: %v", header)
	}
	if strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportExecTime)) {
		t.Fatalf("empty exec time should be pruned: %v", header)
	}
	if !strings.Contains(joined, locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContentPlain)) {
		t.Fatalf("SQL content must remain: %v", header)
	}
	if len(rows) != 1 {
		t.Fatalf("rows=%d", len(rows))
	}
}

func TestBuildGlobalDataExportKeepsOpsTypeWhenPresent(t *testing.T) {
	ctx := context.Background()
	opsType := "共享运维类型-导出验收"
	header, rows := BuildGlobalDataExportExport(ctx, []DataExportExportRecord{
		{
			ProjectName:    "p",
			WorkflowID:     "1",
			WorkflowName:   "n",
			Description:    "d",
			OpsTypeName:    opsType,
			DBServiceNames: []string{"db"},
			CreatedAt:      "2026-01-01 00:00:00",
			CreatorName:    "u",
			UnifiedStatus:  "completed",
			SQLContent:     "SELECT 1",
		},
		{
			ProjectName:    "p2",
			WorkflowID:     "2",
			WorkflowName:   "n2",
			Description:    "d2",
			OpsTypeName:    "",
			DBServiceNames: []string{"db2"},
			CreatedAt:      "2026-01-02 00:00:00",
			CreatorName:    "u2",
			UnifiedStatus:  "completed",
			SQLContent:     "SELECT 2",
		},
	})
	opsCol := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOpsType)
	descCol := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription)
	dsCol := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource)
	opsIdx, descIdx, dsIdx := -1, -1, -1
	for i, h := range header {
		switch h {
		case opsCol:
			opsIdx = i
		case descCol:
			descIdx = i
		case dsCol:
			dsIdx = i
		}
	}
	if opsIdx < 0 {
		t.Fatalf("ops type col must remain when any row has value: %v", header)
	}
	if !(descIdx >= 0 && opsIdx == descIdx+1 && dsIdx == opsIdx+1) {
		t.Fatalf("want description, ops type, data source consecutive, got desc=%d ops=%d ds=%d header=%v",
			descIdx, opsIdx, dsIdx, header)
	}
	if rows[0][opsIdx] != opsType {
		t.Fatalf("row0 ops=%q want %q", rows[0][opsIdx], opsType)
	}
	if rows[1][opsIdx] != "" {
		t.Fatalf("row1 empty ops must be empty cell, got %q", rows[1][opsIdx])
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
		{ProjectName: "p", WorkflowType: "SQL", WorkflowID: "1", Description: "d", Status: sqlReleaseStatus, SQLContent: "s1"},
		{ProjectName: "p", WorkflowType: "DE", WorkflowID: "2", Description: "d", Status: dataExportStatus, SQLContent: "s2"},
	})
	statusCol := -1
	header := BuildHeaderForLayout(ctx, LayoutGlobalCommon, 0)
	// After prune, status column index may shift; locate by building without prune via known order.
	// Rows still follow candidate order before prune... BuildGlobalCommonExport prunes.
	// Find status by matching known value.
	for i, cell := range rows[0] {
		if cell == sqlReleaseStatus {
			statusCol = i
			break
		}
	}
	if statusCol < 0 {
		t.Fatalf("status col not found in row %v header candidate %v", rows[0], header)
	}
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
