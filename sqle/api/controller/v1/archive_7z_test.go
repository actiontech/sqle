package v1

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// get7zTestdataDir returns the absolute path to the testdata/7z directory.
func get7zTestdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller information")
	}
	return filepath.Join(filepath.Dir(filename), "testdata", "7z")
}

// openTest7z opens a 7z file from the testdata/7z directory and returns its content as bytes.Reader.
// sevenzip.NewReader requires io.ReaderAt + size, so we read the entire file into memory.
func openTest7z(t *testing.T, name string) (*bytes.Reader, int64) {
	t.Helper()
	path := filepath.Join(get7zTestdataDir(t), name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test 7z file %s: %v", name, err)
	}
	return bytes.NewReader(data), int64(len(data))
}

func TestProcess7zContent(t *testing.T) {
	cases := map[string]struct {
		szFile           string
		expectSQLCount   int
		expectXMLCount   int
		expectErr        bool
		errContains      string
		checkSQLContains []string // check that SQL results contain these substrings
		checkSortOrder   []string // if non-empty, check that FilePath order matches
	}{
		"normal 7z with .sql file": {
			szFile:           "sql_only.7z",
			expectSQLCount:   1,
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT 1;"},
		},
		"mixed format 7z (.sql + .xml + .txt)": {
			szFile:         "normal.7z",
			expectSQLCount: 2, // query.sql + data.txt
			expectXMLCount: 1, // mapper.xml
			expectErr:      false,
			checkSQLContains: []string{
				"SELECT * FROM users;",              // from query.sql
				"SELECT count(*) FROM products;",    // from data.txt
			},
		},
		"nested archive 7z (contains .zip)": {
			szFile:           "nested.7z",
			expectSQLCount:   1, // query.sql only, inner.zip is skipped
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT * FROM nested_test;"},
		},
		"empty 7z (no auditable files)": {
			szFile:         "empty.7z",
			expectSQLCount: 0,
			expectXMLCount: 0,
			expectErr:      false,
		},
		"7z with unsupported formats (.sql + .png + .jpg)": {
			szFile:           "unsupported.7z",
			expectSQLCount:   1, // only query.sql
			expectXMLCount:   0,
			expectErr:        false,
			checkSQLContains: []string{"SELECT * FROM unsupported_test;"},
		},
		"7z with only unsupported formats": {
			szFile:         "only_unsupported.7z",
			expectSQLCount: 0,
			expectXMLCount: 0,
			expectErr:      false,
		},
		"7z files are naturally sorted": {
			szFile:         "sorted_test.7z",
			expectSQLCount: 3,
			expectXMLCount: 0,
			expectErr:      false,
			checkSortOrder: []string{"file1.sql", "file2.sql", "file11.sql"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r, size := openTest7z(t, tc.szFile)

			sqlFiles, xmlFiles, _, err := process7zContent(r, size)

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

func TestProcess7zContentSizeLimit(t *testing.T) {
	// Test that the size limit check logic works correctly via archiveConfig.
	// Integration with process7zContent is ensured by the shared archiveConfig mechanism.
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

func TestProcess7zContentFileCountLimit(t *testing.T) {
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

func TestProcess7zContentInvalid7z(t *testing.T) {
	// Test with invalid 7z data
	invalidData := bytes.NewReader([]byte("this is not a 7z file"))
	_, _, _, err := process7zContent(invalidData, int64(len("this is not a 7z file")))
	if err == nil {
		t.Error("expected error for invalid 7z data, got nil")
	}
}
