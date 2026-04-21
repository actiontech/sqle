//go:build !enterprise

package auditreport

import (
	"strings"
	"testing"
)

func TestExportAuditReport_CEEdition_PDFBlocked(t *testing.T) {
	testCases := map[string]struct {
		format        ExportFormat
		wantErrSubstr string
	}{
		"PDF format is blocked in CE edition": {
			format:        ExportFormatPDF,
			wantErrSubstr: "enterprise edition",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := buildTestReportData()
			result, err := ExportAuditReport(tc.format, data)
			if err == nil {
				t.Fatal("ExportAuditReport(PDF) should return error in CE edition, got nil")
			}
			if result != nil {
				t.Errorf("ExportAuditReport(PDF) should return nil result in CE edition, got %+v", result)
			}
			if !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Errorf("error message %q does not contain %q", err.Error(), tc.wantErrSubstr)
			}
			if !strings.Contains(err.Error(), string(tc.format)) {
				t.Errorf("error message %q does not contain format name %q", err.Error(), tc.format)
			}
		})
	}
}

func TestExportAuditReport_CEEdition_WORDBlocked(t *testing.T) {
	testCases := map[string]struct {
		format        ExportFormat
		wantErrSubstr string
	}{
		"WORD format is blocked in CE edition": {
			format:        ExportFormatWORD,
			wantErrSubstr: "enterprise edition",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := buildTestReportData()
			result, err := ExportAuditReport(tc.format, data)
			if err == nil {
				t.Fatal("ExportAuditReport(WORD) should return error in CE edition, got nil")
			}
			if result != nil {
				t.Errorf("ExportAuditReport(WORD) should return nil result in CE edition, got %+v", result)
			}
			if !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Errorf("error message %q does not contain %q", err.Error(), tc.wantErrSubstr)
			}
			if !strings.Contains(err.Error(), string(tc.format)) {
				t.Errorf("error message %q does not contain format name %q", err.Error(), tc.format)
			}
		})
	}
}

func TestExportAuditReport_DefaultCSV(t *testing.T) {
	testCases := map[string]struct {
		format        ExportFormat
		wantErrSubstr string
	}{
		"invalid format returns error": {
			format:        ExportFormat("invalid"),
			wantErrSubstr: "unsupported export format",
		},
		"empty format returns error": {
			format:        ExportFormat(""),
			wantErrSubstr: "unsupported export format",
		},
		"unknown format returns error": {
			format:        ExportFormat("json"),
			wantErrSubstr: "unsupported export format",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := buildTestReportData()
			result, err := ExportAuditReport(tc.format, data)
			if err == nil {
				t.Fatalf("ExportAuditReport(%q) expected error, got nil", tc.format)
			}
			if result != nil {
				t.Errorf("ExportAuditReport(%q) expected nil result, got %+v", tc.format, result)
			}
			if !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Errorf("error message %q does not contain %q", err.Error(), tc.wantErrSubstr)
			}
			if !strings.Contains(err.Error(), string(tc.format)) {
				t.Errorf("error message %q does not contain format name %q", err.Error(), tc.format)
			}
		})
	}
}
