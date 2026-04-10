package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNormalizeExportFormatStr(t *testing.T) {
	testCases := map[string]struct {
		input       string
		expected    ExportFormat
		expectError bool
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
		"invalid value returns error": {
			input:       "invalid",
			expectError: true,
		},
		"unknown format returns error": {
			input:       "json",
			expectError: true,
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
			result, err := NormalizeExportFormatStr(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("NormalizeExportFormatStr(%q) expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("NormalizeExportFormatStr(%q) unexpected error: %v", tc.input, err)
				return
			}
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

func TestHTMLReportGenerator_Normal(t *testing.T) {
	gen, err := NewHTMLReportGenerator()
	if err != nil {
		t.Fatalf("NewHTMLReportGenerator() returned unexpected error: %v", err)
	}

	// Verify the generator implements ReportGenerator interface
	var _ ReportGenerator = gen

	testCases := map[string]struct {
		data             *AuditReportData
		wantContentType  string
		wantFilePrefix   string
		wantFileSuffix   string
		wantHTMLTags     []string
		wantSQLContents  []string
		wantLabels       []string
	}{
		"normal data generates valid HTML report": {
			data:            buildTestReportData(),
			wantContentType: "text/html",
			wantFilePrefix:  "SQL_audit_report_test-mysql_1001",
			wantFileSuffix:  ".html",
			wantHTMLTags:    []string{"<table>", "<h1>", "<h2>", "<pre>", "<!DOCTYPE html>", "</html>"},
			wantSQLContents: []string{"SELECT * FROM users", "DROP TABLE test"},
			wantLabels:      []string{"Audit Summary", "Audit Result Statistics", "Problem SQL List", "Rule Hit Statistics"},
		},
	}

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

			// Verify FileName format
			if !strings.HasPrefix(result.FileName, tc.wantFilePrefix) {
				t.Errorf("FileName %q does not start with %q", result.FileName, tc.wantFilePrefix)
			}
			if !strings.HasSuffix(result.FileName, tc.wantFileSuffix) {
				t.Errorf("FileName %q does not end with %q", result.FileName, tc.wantFileSuffix)
			}
			if !strings.Contains(result.FileName, tc.data.InstanceName) {
				t.Errorf("FileName %q does not contain InstanceName %q", result.FileName, tc.data.InstanceName)
			}
			if !strings.Contains(result.FileName, "1001") {
				t.Errorf("FileName %q does not contain TaskID", result.FileName)
			}

			// Verify HTML content contains key HTML tags
			content := string(result.Content)
			for _, tag := range tc.wantHTMLTags {
				if !strings.Contains(content, tag) {
					t.Errorf("Content does not contain expected HTML tag %q", tag)
				}
			}

			// Verify SQL contents exist in the output
			for _, sql := range tc.wantSQLContents {
				if !strings.Contains(content, sql) {
					t.Errorf("Content does not contain expected SQL %q", sql)
				}
			}

			// Verify i18n labels are rendered
			for _, label := range tc.wantLabels {
				if !strings.Contains(content, label) {
					t.Errorf("Content does not contain expected label %q", label)
				}
			}

			// Verify Format() returns ExportFormatHTML
			if gen.Format() != ExportFormatHTML {
				t.Errorf("Format() = %q, want %q", gen.Format(), ExportFormatHTML)
			}
		})
	}
}

func TestHTMLReportGenerator_XSSPrevention(t *testing.T) {
	gen, err := NewHTMLReportGenerator()
	if err != nil {
		t.Fatalf("NewHTMLReportGenerator() returned unexpected error: %v", err)
	}

	testCases := map[string]struct {
		maliciousSQL    string
		wantAbsent      []string
		wantDescription string
	}{
		"script tag in SQL is escaped": {
			maliciousSQL: "<script>alert('xss')</script>",
			wantAbsent:   []string{"<script>", "</script>"},
			wantDescription: "script tags should be HTML-escaped by html/template",
		},
		"img onerror in SQL is escaped": {
			maliciousSQL: `<img src=x onerror="alert('xss')">`,
			wantAbsent:   []string{`onerror="alert`},
			wantDescription: "event handler attributes should be HTML-escaped",
		},
		"script tag in description is escaped": {
			maliciousSQL: "SELECT 1",
			wantAbsent:   []string{"<script>"},
			wantDescription: "script tags in other fields are also escaped",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := buildTestReportData()

			if name == "script tag in description is escaped" {
				// Put malicious content in description field
				data.ProblemSQLs[0].Description = "<script>alert('xss')</script>"
			} else {
				// Put malicious content in SQL field
				data.ProblemSQLs[0].SQL = tc.maliciousSQL
				data.SQLList[0].SQL = tc.maliciousSQL
			}

			result, err := gen.Generate(data)
			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}

			content := string(result.Content)

			// Verify that raw malicious content is NOT present (it should be escaped)
			for _, absent := range tc.wantAbsent {
				if strings.Contains(content, absent) {
					snippetLen := 500
					if len(content) < snippetLen {
						snippetLen = len(content)
					}
					t.Errorf("%s: content contains unescaped %q\ncontent snippet: %s",
						tc.wantDescription, absent, content[:snippetLen])
				}
			}
		})
	}
}

func TestHTMLReportGenerator_EmptyData(t *testing.T) {
	gen, err := NewHTMLReportGenerator()
	if err != nil {
		t.Fatalf("NewHTMLReportGenerator() returned unexpected error: %v", err)
	}

	testCases := map[string]struct {
		data         *AuditReportData
		wantHTMLTags []string
	}{
		"empty SQL list renders without error": {
			data: &AuditReportData{
				TaskID:       2002,
				Title:        "Empty Report",
				InstanceName: "empty-instance",
				Schema:       "empty_db",
				GeneratedAt:  time.Now(),
				Lang:         "en-US",
				Summary: AuditSummary{
					AuditTime:    "2026-03-31 10:00:00",
					InstanceName: "empty-instance",
					Schema:       "empty_db",
					TotalSQL:     0,
					PassRate:     100.0,
					Score:        100,
				},
				Statistics: AuditStatistics{
					LevelDistribution: []LevelCount{},
					RuleHits:          []RuleHit{},
				},
				SQLList:     []AuditSQLItem{},
				ProblemSQLs: []AuditSQLItem{},
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
					AuditResult:       "Audit Result",
					RuleName:          "Rule Name",
					Description:       "Description",
					Suggestion:        "Suggestion",
					Count:             "Count",
					HitCount:          "Hit Count",
				},
			},
			wantHTMLTags: []string{"<!DOCTYPE html>", "</html>", "<table>", "<h1>", "<h2>"},
		},
	}

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
			if result.ContentType != "text/html" {
				t.Errorf("ContentType = %q, want %q", result.ContentType, "text/html")
			}

			// Verify the output is valid HTML with required tags
			content := string(result.Content)
			for _, tag := range tc.wantHTMLTags {
				if !strings.Contains(content, tag) {
					t.Errorf("Content does not contain expected HTML tag %q", tag)
				}
			}

			// Verify FileName
			expectedFileName := "SQL_audit_report_empty-instance_2002.html"
			if result.FileName != expectedFileName {
				t.Errorf("FileName = %q, want %q", result.FileName, expectedFileName)
			}

			// Verify the title is rendered
			if !strings.Contains(content, "Empty Report") {
				t.Error("Content does not contain the report title")
			}
		})
	}
}

func TestHTMLReportGenerator_LargeData(t *testing.T) {
	gen, err := NewHTMLReportGenerator()
	if err != nil {
		t.Fatalf("NewHTMLReportGenerator() returned unexpected error: %v", err)
	}

	testCases := map[string]struct {
		sqlCount int
	}{
		"10000 SQL items generates successfully": {
			sqlCount: 10000,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data := buildTestReportData()

			// Build large SQL list and problem SQL list
			sqlList := make([]AuditSQLItem, 0, tc.sqlCount)
			problemSQLs := make([]AuditSQLItem, 0, tc.sqlCount/2)
			for i := 0; i < tc.sqlCount; i++ {
				item := AuditSQLItem{
					Number:      uint(i + 1),
					SQL:         fmt.Sprintf("SELECT * FROM table_%d WHERE id = %d", i, i),
					AuditLevel:  "warn",
					AuditStatus: "finished",
					AuditResult: fmt.Sprintf("audit result for SQL #%d", i),
					RuleName:    "no_select_all",
					Description: fmt.Sprintf("description for SQL #%d", i),
					Suggestion:  "specify column names",
				}
				sqlList = append(sqlList, item)
				if i%2 == 0 {
					problemSQLs = append(problemSQLs, item)
				}
			}
			data.SQLList = sqlList
			data.ProblemSQLs = problemSQLs
			data.Summary.TotalSQL = tc.sqlCount

			start := time.Now()
			result, err := gen.Generate(data)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("Generate() returned unexpected error: %v", err)
			}
			if result == nil {
				t.Fatal("Generate() returned nil result")
			}

			// Verify content is not empty
			if len(result.Content) == 0 {
				t.Error("Generate() returned empty content")
			}

			// Verify ContentType
			if result.ContentType != "text/html" {
				t.Errorf("ContentType = %q, want %q", result.ContentType, "text/html")
			}

			// Log the elapsed time (informational, not a hard failure since dev machines vary)
			t.Logf("Generated HTML report with %d SQLs in %v, output size: %d bytes",
				tc.sqlCount, elapsed, len(result.Content))

			// Soft check: warn if it takes more than 5 seconds
			if elapsed > 5*time.Second {
				t.Logf("WARNING: Generation took %v which exceeds 5s target", elapsed)
			}
		})
	}
}
