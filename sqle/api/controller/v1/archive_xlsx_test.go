package v1

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getXlsxTestdataDir returns the absolute path to the testdata/xlsx directory.
func getXlsxTestdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller information")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "xlsx")
}

// openTestXlsx opens an XLSX file from testdata/xlsx and returns its bytes as io.Reader.
func openTestXlsx(t *testing.T, name string) *bytes.Reader {
	t.Helper()
	path := filepath.Join(getXlsxTestdataDir(t), name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test XLSX file %s: %v", name, err)
	}
	return bytes.NewReader(data)
}

func TestProcessXlsxContent(t *testing.T) {
	cases := map[string]struct {
		xlsxFile        string
		expectSQL       string   // expected joined SQL output
		expectErr       bool
		errContains     string
		checkContains   []string // check that result contains these substrings
		expectSQLCount  int      // expected number of SQL statements (separated by ";\n")
		expectEmpty     bool     // expect empty string result (no error)
	}{
		"standard template with 序号/SQL/备注 header": {
			xlsxFile:       "standard_template.xlsx",
			expectErr:      false,
			expectSQLCount: 3,
			checkContains: []string{
				"SELECT 1",
				"SELECT 2",
				"INSERT INTO t1 VALUES (1)",
			},
		},
		"SQL column name lowercase 'sql'": {
			xlsxFile:       "sql_lowercase.xlsx",
			expectErr:      false,
			expectSQLCount: 1,
			checkContains:  []string{"SELECT 'lowercase'"},
		},
		"SQL column name uppercase 'SQL'": {
			xlsxFile:       "sql_uppercase.xlsx",
			expectErr:      false,
			expectSQLCount: 1,
			checkContains:  []string{"SELECT 'uppercase'"},
		},
		"SQL column name variant 'Sql Statement'": {
			xlsxFile:       "sql_mixed_name.xlsx",
			expectErr:      false,
			expectSQLCount: 1,
			checkContains:  []string{"SELECT 'mixed'"},
		},
		"no SQL column returns error": {
			xlsxFile:    "no_sql_column.xlsx",
			expectErr:   true,
			errContains: "no column containing \"SQL\" found",
		},
		"with empty rows skips blank lines": {
			xlsxFile:       "with_empty_rows.xlsx",
			expectErr:      false,
			expectSQLCount: 2,
			checkContains: []string{
				"SELECT 1",
				"SELECT 3",
			},
		},
		"multi sheet reads only the first sheet": {
			xlsxFile:       "multi_sheet.xlsx",
			expectErr:      false,
			expectSQLCount: 1,
			checkContains:  []string{"SELECT 'first_sheet'"},
		},
		"empty file returns empty content": {
			xlsxFile:    "empty_file.xlsx",
			expectErr:   false,
			expectEmpty: true,
		},
		"invalid xlsx data returns error": {
			xlsxFile:    "",
			expectErr:   true,
			errContains: "open xlsx file failed",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var r *bytes.Reader
			if tc.xlsxFile == "" {
				// Test with invalid data
				r = bytes.NewReader([]byte("this is not a valid xlsx file"))
			} else {
				r = openTestXlsx(t, tc.xlsxFile)
			}

			result, err := processXlsxContent(r)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("expected error containing %q, got: %v", tc.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.expectEmpty {
				if result != "" {
					t.Fatalf("expected empty result, got: %q", result)
				}
				return
			}

			// Check SQL count
			if tc.expectSQLCount > 0 {
				parts := strings.Split(result, ";\n")
				if len(parts) != tc.expectSQLCount {
					t.Fatalf("expected %d SQL statements, got %d. Result: %q", tc.expectSQLCount, len(parts), result)
				}
			}

			// Check that result contains expected substrings
			for _, substr := range tc.checkContains {
				if !strings.Contains(result, substr) {
					t.Errorf("expected result to contain %q, got: %q", substr, result)
				}
			}
		})
	}
}

// TestProcessXlsxContentSQLJoinFormat verifies the SQL join format uses ";\n" separator.
func TestProcessXlsxContentSQLJoinFormat(t *testing.T) {
	r := openTestXlsx(t, "standard_template.xlsx")
	result, err := processXlsxContent(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the join format is ";\n"
	expectedJoined := "SELECT 1;\nSELECT 2;\nINSERT INTO t1 VALUES (1)"
	if result != expectedJoined {
		t.Fatalf("expected exact output %q, got %q", expectedJoined, result)
	}
}

// TestProcessXlsxContentMultiSheetDoesNotReadSecondSheet verifies that only
// the first sheet is read and second sheet content is excluded.
func TestProcessXlsxContentMultiSheetDoesNotReadSecondSheet(t *testing.T) {
	r := openTestXlsx(t, "multi_sheet.xlsx")
	result, err := processXlsxContent(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "second_sheet") {
		t.Fatalf("result should not contain content from second sheet, got: %q", result)
	}
}
