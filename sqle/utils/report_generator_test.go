package utils

import (
	"testing"
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
