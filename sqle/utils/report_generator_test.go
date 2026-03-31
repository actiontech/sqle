package utils

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeExportFormatStr(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected ExportFormat
	}{
		"empty string defaults to csv": {
			input:    "",
			expected: CsvExportFormat,
		},
		"csv returns csv": {
			input:    "csv",
			expected: CsvExportFormat,
		},
		"CSV uppercase returns csv": {
			input:    "CSV",
			expected: CsvExportFormat,
		},
		"excel returns excel": {
			input:    "excel",
			expected: ExcelExportFormat,
		},
		"xlsx returns excel": {
			input:    "xlsx",
			expected: ExcelExportFormat,
		},
		"html returns html": {
			input:    "html",
			expected: ExportFormatHTML,
		},
		"HTML uppercase returns html": {
			input:    "HTML",
			expected: ExportFormatHTML,
		},
		"pdf returns pdf": {
			input:    "pdf",
			expected: ExportFormatPDF,
		},
		"PDF uppercase returns pdf": {
			input:    "PDF",
			expected: ExportFormatPDF,
		},
		"word returns word": {
			input:    "word",
			expected: ExportFormatWORD,
		},
		"WORD uppercase returns word": {
			input:    "WORD",
			expected: ExportFormatWORD,
		},
		"docx returns word": {
			input:    "docx",
			expected: ExportFormatWORD,
		},
		"DOCX uppercase returns word": {
			input:    "DOCX",
			expected: ExportFormatWORD,
		},
		"invalid value defaults to csv": {
			input:    "invalid",
			expected: CsvExportFormat,
		},
		"unknown format defaults to csv": {
			input:    "json",
			expected: CsvExportFormat,
		},
		"whitespace-only defaults to csv": {
			input:    "   ",
			expected: CsvExportFormat,
		},
		"leading and trailing spaces are trimmed": {
			input:    "  pdf  ",
			expected: ExportFormatPDF,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := NormalizeExportFormatStr(tc.input)
			if result != tc.expected {
				t.Errorf("NormalizeExportFormatStr(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// buildTestReportData 构建测试用的 AuditReportData
func buildTestReportData() *AuditReportData {
	return &AuditReportData{
		TaskID:       1001,
		Title:        "SQL Audit Report",
		InstanceName: "test-mysql",
		Schema:       "test_db",
		GeneratedAt:  time.Now(),
		Lang:         "en-US",
		Summary: AuditSummary{
			AuditTime:    "2026-03-31 10:00:00",
			InstanceName: "test-mysql",
			Schema:       "test_db",
			TotalSQL:     3,
			PassRate:     66.7,
			Score:        70,
			AuditLevel:   "warn",
		},
		Statistics: AuditStatistics{
			LevelDistribution: []LevelCount{
				{Level: "normal", Count: 1},
				{Level: "warn", Count: 1},
				{Level: "error", Count: 1},
			},
			RuleHits: []RuleHit{
				{RuleName: "no_select_all", HitCount: 1},
				{RuleName: "no_drop_table", HitCount: 1},
			},
		},
		SQLList: []AuditSQLItem{
			{
				Number:      1,
				SQL:         "SELECT * FROM users",
				AuditLevel:  "warn",
				AuditStatus: "finished",
				AuditResult: "should not use SELECT *",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "",
				Description: "query all users",
				RuleName:    "no_select_all",
				Suggestion:  "specify column names",
			},
			{
				Number:      2,
				SQL:         "DROP TABLE test",
				AuditLevel:  "error",
				AuditStatus: "finished",
				AuditResult: "DROP TABLE is prohibited",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "",
				Description: "drop test table",
				RuleName:    "no_drop_table",
				Suggestion:  "do not use DROP TABLE",
			},
			{
				Number:      3,
				SQL:         "INSERT INTO t VALUES(1)",
				AuditLevel:  "normal",
				AuditStatus: "finished",
				AuditResult: "",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "DELETE FROM t WHERE id=1",
				Description: "insert a row",
			},
		},
		ProblemSQLs: []AuditSQLItem{
			{Number: 1, SQL: "SELECT * FROM users", AuditLevel: "warn", RuleName: "no_select_all"},
			{Number: 2, SQL: "DROP TABLE test", AuditLevel: "error", RuleName: "no_drop_table"},
		},
		Labels: ReportLabels{
			AuditSummary:      "Audit Summary",
			ResultStatistics:  "Audit Result Statistics",
			ProblemSQLList:    "Problem SQL List",
			RuleHitStatistics: "Rule Hit Statistics",
			AuditTime:         "Audit Time",
			DataSource:        "Data Source",
			Schema:            "Schema",
			TotalSQL:          "Total SQL",
			PassRate:          "Pass Rate",
			Score:             "Score",
			AuditLevel:        "Audit Level",
			Number:            "Number",
			SQL:               "SQL",
			AuditStatus:       "Audit Status",
			AuditResult:       "Audit Result",
			ExecStatus:        "Exec Status",
			ExecResult:        "Exec Result",
			RollbackSQL:       "Rollback SQL",
			RuleName:          "Rule Name",
			Description:       "Description",
			Suggestion:        "Suggestion",
			Count:             "Count",
			HitCount:          "Hit Count",
		},
	}
}

func TestCSVReportGenerator_Normal(t *testing.T) {
	testCases := map[string]struct {
		data            *AuditReportData
		wantContentType string
		wantFilePrefix  string
		wantFileSuffix  string
		wantBOM         bool
		wantHeaders     []string
		wantDataRows    int
	}{
		"normal data generates valid CSV report": {
			data:            buildTestReportData(),
			wantContentType: "text/csv",
			wantFilePrefix:  "SQL_audit_report_test-mysql_1001",
			wantFileSuffix:  ".csv",
			wantBOM:         true,
			wantHeaders:     []string{"Number", "SQL", "Audit Status", "Audit Result", "Exec Status", "Exec Result", "Rollback SQL", "Description"},
			wantDataRows:    3,
		},
	}

	gen := NewCSVReportGenerator()

	// Verify the generator implements ReportGenerator interface
	var _ ReportGenerator = gen

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := gen.Generate(tc.data)
			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("Generate() returned nil result")
			}

			// Verify ContentType
			if result.ContentType != tc.wantContentType {
				t.Errorf("ContentType = %q, want %q", result.ContentType, tc.wantContentType)
			}

			// Verify FileName contains InstanceName and TaskID, ends with .csv
			if !strings.Contains(result.FileName, tc.data.InstanceName) {
				t.Errorf("FileName %q does not contain InstanceName %q", result.FileName, tc.data.InstanceName)
			}
			if !strings.Contains(result.FileName, "1001") {
				t.Errorf("FileName %q does not contain TaskID", result.FileName)
			}
			if !strings.HasSuffix(result.FileName, tc.wantFileSuffix) {
				t.Errorf("FileName %q does not end with %q", result.FileName, tc.wantFileSuffix)
			}
			if !strings.HasPrefix(result.FileName, tc.wantFilePrefix) {
				t.Errorf("FileName %q does not start with %q", result.FileName, tc.wantFilePrefix)
			}

			// Verify UTF-8 BOM
			content := string(result.Content)
			if tc.wantBOM {
				if !strings.HasPrefix(content, "\xEF\xBB\xBF") {
					t.Error("Content does not start with UTF-8 BOM")
				}
			}

			// Verify headers exist in content
			for _, h := range tc.wantHeaders {
				if !strings.Contains(content, h) {
					t.Errorf("Content does not contain header %q", h)
				}
			}

			// Verify the number of data rows (excluding BOM and header line)
			contentWithoutBOM := strings.TrimPrefix(content, "\xEF\xBB\xBF")
			lines := strings.Split(strings.TrimRight(contentWithoutBOM, "\n"), "\n")
			// First line is header, remaining are data rows
			dataLineCount := len(lines) - 1
			if dataLineCount != tc.wantDataRows {
				t.Errorf("data row count = %d, want %d", dataLineCount, tc.wantDataRows)
			}
		})
	}
}

func TestCSVReportGenerator_EmptyData(t *testing.T) {
	testCases := map[string]struct {
		data         *AuditReportData
		wantDataRows int
	}{
		"empty SQL list produces header only": {
			data: &AuditReportData{
				TaskID:       2002,
				InstanceName: "empty-instance",
				SQLList:      []AuditSQLItem{},
				Labels: ReportLabels{
					Number:      "Number",
					SQL:         "SQL",
					AuditStatus: "Audit Status",
					AuditResult: "Audit Result",
					ExecStatus:  "Exec Status",
					ExecResult:  "Exec Result",
					RollbackSQL: "Rollback SQL",
					Description: "Description",
				},
			},
			wantDataRows: 0,
		},
	}

	gen := NewCSVReportGenerator()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := gen.Generate(tc.data)
			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("Generate() returned nil result")
			}

			content := string(result.Content)

			// Verify BOM is present
			if !strings.HasPrefix(content, "\xEF\xBB\xBF") {
				t.Error("Content does not start with UTF-8 BOM")
			}

			// Verify only header line, no data rows
			contentWithoutBOM := strings.TrimPrefix(content, "\xEF\xBB\xBF")
			lines := strings.Split(strings.TrimRight(contentWithoutBOM, "\n"), "\n")
			// Should have exactly 1 line (header only)
			if len(lines) != 1 {
				t.Errorf("expected 1 line (header only), got %d lines", len(lines))
			}

			// Verify headers are present
			headerLine := lines[0]
			for _, h := range []string{"Number", "SQL", "Audit Status"} {
				if !strings.Contains(headerLine, h) {
					t.Errorf("header line does not contain %q", h)
				}
			}

			// Verify FileName
			if result.FileName != "SQL_audit_report_empty-instance_2002.csv" {
				t.Errorf("FileName = %q, want %q", result.FileName, "SQL_audit_report_empty-instance_2002.csv")
			}
		})
	}
}

func TestCSVReportGenerator_SpecialChars(t *testing.T) {
	testCases := map[string]struct {
		sqlItem     AuditSQLItem
		wantInRow   string
		description string
	}{
		"SQL with comma is quoted": {
			sqlItem: AuditSQLItem{
				Number:      1,
				SQL:         "SELECT a, b FROM users",
				AuditStatus: "finished",
				AuditResult: "ok",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "",
				Description: "",
			},
			wantInRow:   `"SELECT a, b FROM users"`,
			description: "field containing comma should be wrapped in double quotes",
		},
		"SQL with double quote is escaped": {
			sqlItem: AuditSQLItem{
				Number:      2,
				SQL:         `SELECT "name" FROM users`,
				AuditStatus: "finished",
				AuditResult: "ok",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "",
				Description: "",
			},
			wantInRow:   `"SELECT ""name"" FROM users"`,
			description: "double quotes within a field should be escaped as two double quotes",
		},
		"SQL with newline is quoted": {
			sqlItem: AuditSQLItem{
				Number:      3,
				SQL:         "SELECT *\nFROM users",
				AuditStatus: "finished",
				AuditResult: "ok",
				ExecStatus:  "initialized",
				ExecResult:  "",
				RollbackSQL: "",
				Description: "",
			},
			wantInRow:   "\"SELECT *\nFROM users\"",
			description: "field containing newline should be wrapped in double quotes",
		},
	}

	gen := NewCSVReportGenerator()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := &AuditReportData{
				TaskID:       3003,
				InstanceName: "special-instance",
				SQLList:      []AuditSQLItem{tc.sqlItem},
				Labels: ReportLabels{
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

			result, err := gen.Generate(data)
			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("Generate() returned nil result")
			}

			content := string(result.Content)

			// Verify that the special character handling is correct
			if !strings.Contains(content, tc.wantInRow) {
				t.Errorf("%s\ncontent does not contain expected substring %q\nfull content:\n%s",
					tc.description, tc.wantInRow, content)
			}
		})
	}
}

func TestCSVReportGenerator_Format(t *testing.T) {
	gen := NewCSVReportGenerator()
	if gen.Format() != CsvExportFormat {
		t.Errorf("Format() = %q, want %q", gen.Format(), CsvExportFormat)
	}
}
