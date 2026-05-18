package dms

import "testing"

// Test_parseInstanceDBType 覆盖 DMS → SQLE 透传 db_type 的规范化映射。
// 该映射是 sqle-ee #2877 修复 R5 「DMS 透传字面值与 SQLE 后端契约不一致」
// 的唯一入口；任何映射条目改动都必须经过此测试。
func Test_parseInstanceDBType(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected string
	}{
		"GaussDB / openGauss alias normalized to GaussDB": {
			input:    "GaussDB / openGauss",
			expected: "GaussDB",
		},
		"MySQL passthrough untouched": {
			input:    "MySQL",
			expected: "MySQL",
		},
		"PostgreSQL passthrough untouched": {
			input:    "PostgreSQL",
			expected: "PostgreSQL",
		},
		"Empty string returns empty (no panic, no fallback)": {
			input:    "",
			expected: "",
		},
		"Unknown DBType keeps original literal": {
			input:    "TiDB",
			expected: "TiDB",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := parseInstanceDBType(tc.input)
			if got != tc.expected {
				t.Fatalf("parseInstanceDBType(%q): expected %q, got %q", tc.input, tc.expected, got)
			}
		})
	}
}
