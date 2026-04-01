package v1

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getTestdataDir returns the absolute path to the testdata/rar directory.
func getTestdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller information")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "rar")
}

// openTestRar opens a RAR file from the testdata/rar directory.
func openTestRar(t *testing.T, name string) *os.File {
	t.Helper()
	path := filepath.Join(getTestdataDir(t), name)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open test RAR file %s: %v", name, err)
	}
	return f
}

func TestProcessRarContent(t *testing.T) {
	cases := map[string]struct {
		rarFile          string
		expectSQLCount   int
		expectXMLCount   int
		expectErr        bool
		errContains      string
		checkSQLContains []string // check that SQL results contain these substrings
		checkSortOrder   []string // if non-empty, check that FilePath order matches
	}{
		"normal RAR with .sql file": {
			rarFile:          "sql_only.rar",
			expectSQLCount:   1,
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT 1;"},
		},
		"mixed format RAR (.sql + .xml + .txt)": {
			rarFile:        "normal.rar",
			expectSQLCount: 2, // query.sql + data.txt
			expectXMLCount: 1, // mapper.xml
			expectErr:      false,
			checkSQLContains: []string{
				"SELECT * FROM users;",    // from query.sql
				"SELECT count(*) FROM products;", // from data.txt
			},
		},
		"nested archive RAR (contains .zip)": {
			rarFile:          "nested.rar",
			expectSQLCount:   1, // query.sql only, inner.zip is skipped
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT * FROM nested_test;"},
		},
		"empty RAR (no auditable files)": {
			rarFile:        "empty.rar",
			expectSQLCount: 0,
			expectXMLCount: 0,
			expectErr:      false,
		},
		"RAR with unsupported formats (.sql + .png + .jpg)": {
			rarFile:          "unsupported.rar",
			expectSQLCount:   1, // only query.sql
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT * FROM unsupported_test;"},
		},
		"RAR with only unsupported formats": {
			rarFile:        "only_unsupported.rar",
			expectSQLCount: 0,
			expectXMLCount: 0,
			expectErr:      false,
		},
		"RAR files are naturally sorted": {
			rarFile:        "sorted_test.rar",
			expectSQLCount: 3,
			expectXMLCount: 0,
			expectErr:      false,
			checkSortOrder: []string{"file1.sql", "file2.sql", "file11.sql"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := openTestRar(t, tc.rarFile)
			defer f.Close()

			sqlFiles, xmlFiles, err := processRarContent(f)

			// Check error
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("expected error containing %q, got: %v", tc.errContains, err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
				return
			}

			// Check SQL file count
			if len(sqlFiles) != tc.expectSQLCount {
				t.Errorf("expected %d SQL files, got %d", tc.expectSQLCount, len(sqlFiles))
				for i, sf := range sqlFiles {
					t.Logf("  SQL[%d]: FilePath=%q, SQLs=%q", i, sf.FilePath, sf.SQLs)
				}
			}

			// Check XML file count
			if len(xmlFiles) != tc.expectXMLCount {
				t.Errorf("expected %d XML files, got %d", tc.expectXMLCount, len(xmlFiles))
				for i, xf := range xmlFiles {
					t.Logf("  XML[%d]: FilePath=%q, SQL=%q", i, xf.FilePath, xf.SQL)
				}
			}

			// Check SQL content contains expected substrings
			if len(tc.checkSQLContains) > 0 {
				allSQLs := ""
				for _, sf := range sqlFiles {
					allSQLs += sf.SQLs + "\n"
				}
				for _, expected := range tc.checkSQLContains {
					if !strings.Contains(allSQLs, expected) {
						t.Errorf("expected SQL results to contain %q, but not found in:\n%s", expected, allSQLs)
					}
				}
			}

			// Check natural sort order
			if len(tc.checkSortOrder) > 0 {
				if len(sqlFiles) != len(tc.checkSortOrder) {
					t.Errorf("sort order check: expected %d files, got %d", len(tc.checkSortOrder), len(sqlFiles))
				} else {
					for i, expectedPath := range tc.checkSortOrder {
						if sqlFiles[i].FilePath != expectedPath {
							t.Errorf("sort order check: position %d expected %q, got %q", i, expectedPath, sqlFiles[i].FilePath)
						}
					}
				}
			}
		})
	}
}

func TestProcessRarContentSizeLimit(t *testing.T) {
	// Test that processRarContent enforces size limits.
	// We generate RAR content in-memory with a large payload that exceeds the limit.
	// Since we can't easily generate valid RAR in Go without a writer,
	// we test the size check logic via the archiveConfig directly.

	// Instead, test with a reader that simulates a RAR file where totalSize exceeds limit.
	// This is covered by the archiveConfig unit tests, but we verify the integration here
	// by checking that the error message format is correct.
	cfg := archiveConfig{
		MaxTotalSize:    100,
		MaxFileCount:    10,
		MaxNestingDepth: 1,
	}

	err := cfg.checkSize(0, 200)
	if err == nil {
		t.Error("expected size limit error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected error message to contain 'exceeds limit', got: %v", err)
	}
}

func TestProcessRarContentFileCountLimit(t *testing.T) {
	// Test file count limit via archiveConfig
	cfg := archiveConfig{
		MaxTotalSize:    10 * 1024 * 1024,
		MaxFileCount:    5,
		MaxNestingDepth: 1,
	}

	err := cfg.checkFileCount(6)
	if err == nil {
		t.Error("expected file count limit error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds limit") {
		t.Errorf("expected error message to contain 'exceeds limit', got: %v", err)
	}
}

func TestProcessRarContentInvalidRar(t *testing.T) {
	// Test with invalid RAR data
	invalidData := bytes.NewReader([]byte("this is not a rar file"))
	_, _, err := processRarContent(invalidData)
	if err == nil {
		t.Error("expected error for invalid RAR data, got nil")
	}
}
