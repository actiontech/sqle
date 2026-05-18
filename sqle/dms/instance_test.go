package dms

import (
	"testing"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsV2 "github.com/actiontech/dms/pkg/dms-common/api/dms/v2"
	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
)

// Test_ParseInstanceDBType 覆盖 DMS → SQLE 透传 db_type 的规范化映射。
// 该映射是 sqle-ee #2877 修复 R5 「DMS 透传字面值与 SQLE 后端契约不一致」
// 的唯一入口；任何映射条目改动都必须经过此测试。
func Test_ParseInstanceDBType(t *testing.T) {
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
			got := ParseInstanceDBType(tc.input)
			if got != tc.expected {
				t.Fatalf("ParseInstanceDBType(%q): expected %q, got %q", tc.input, tc.expected, got)
			}
		})
	}
}

// TestConvertInstance_EmptyRuleTemplateID 覆盖 #2877 bug-B：
// DMS 透传的 ListDBService.SQLEConfig.RuleTemplateID 在"未配规则模板"场景下
// 始终是空字符串，convertInstance 不得再因为 strconv.ParseInt("") 报错；
// 应把空字符串映射为 0，并返回合法的 *model.Instance（DbType 经 ParseInstanceDBType 规范化）。
func TestConvertInstance_EmptyRuleTemplateID(t *testing.T) {
	// 用 AesEncrypt 构造合法的加密密码字符串，避免 convertInstance 在 AES 解密阶段就先报错。
	encryptedPwd, err := dmsCommonAes.AesEncrypt("test-password")
	if err != nil {
		t.Fatalf("AesEncrypt setup failed: %v", err)
	}

	inst := &dmsV2.ListDBService{
		DBServiceUid: "100",
		Name:         "opengauss_test",
		DBType:       "GaussDB / openGauss",
		Host:         "127.0.0.1",
		Port:         "5432",
		User:         "gaussdb",
		Password:     encryptedPwd,
		ProjectUID:   "700300",
		SQLEConfig: &dmsCommonV1.SQLEConfig{
			AuditEnabled:     false,
			RuleTemplateID:   "", // 关键：未配规则模板
			RuleTemplateName: "",
			SQLQueryConfig:   &dmsCommonV1.SQLQueryConfig{}, // convertInstance 内部解引用 SQLQueryConfig，需要非 nil
		},
	}

	got, err := convertInstance(inst)
	if err != nil {
		t.Fatalf("convertInstance with empty RuleTemplateID returned error: %v (expected nil after #2877 bug-B fix)", err)
	}
	if got == nil {
		t.Fatalf("convertInstance returned nil instance with empty RuleTemplateID")
	}
	if got.RuleTemplateId != 0 {
		t.Fatalf("convertInstance: RuleTemplateId = %d, want 0 for empty RuleTemplateID", got.RuleTemplateId)
	}
	// 同时验证 DbType 仍经过 ParseInstanceDBType 规范化（GaussDB / openGauss → GaussDB）
	if got.DbType != "GaussDB" {
		t.Fatalf("convertInstance: DbType = %q, want %q (ParseInstanceDBType normalization)", got.DbType, "GaussDB")
	}
}

// TestConvertInstance_NonEmptyRuleTemplateID 回归：非空合法数字字符串仍按 ParseInt 处理。
func TestConvertInstance_NonEmptyRuleTemplateID(t *testing.T) {
	encryptedPwd, err := dmsCommonAes.AesEncrypt("test-password")
	if err != nil {
		t.Fatalf("AesEncrypt setup failed: %v", err)
	}

	inst := &dmsV2.ListDBService{
		DBServiceUid: "100",
		Name:         "mysql_with_rule",
		DBType:       "MySQL",
		Host:         "127.0.0.1",
		Port:         "3306",
		User:         "root",
		Password:     encryptedPwd,
		ProjectUID:   "700300",
		SQLEConfig: &dmsCommonV1.SQLEConfig{
			AuditEnabled:     true,
			RuleTemplateID:   "42",
			RuleTemplateName: "my_template",
			SQLQueryConfig:   &dmsCommonV1.SQLQueryConfig{},
		},
	}

	got, err := convertInstance(inst)
	if err != nil {
		t.Fatalf("convertInstance with valid RuleTemplateID failed: %v", err)
	}
	if got.RuleTemplateId != 42 {
		t.Fatalf("convertInstance: RuleTemplateId = %d, want 42", got.RuleTemplateId)
	}
}

// TestConvertInstance_InvalidRuleTemplateID 回归：非法非空字符串（不是数字）仍报错。
func TestConvertInstance_InvalidRuleTemplateID(t *testing.T) {
	encryptedPwd, err := dmsCommonAes.AesEncrypt("test-password")
	if err != nil {
		t.Fatalf("AesEncrypt setup failed: %v", err)
	}

	inst := &dmsV2.ListDBService{
		DBServiceUid: "100",
		Name:         "mysql_invalid_rule",
		DBType:       "MySQL",
		Host:         "127.0.0.1",
		Port:         "3306",
		User:         "root",
		Password:     encryptedPwd,
		ProjectUID:   "700300",
		SQLEConfig: &dmsCommonV1.SQLEConfig{
			RuleTemplateID: "not-a-number",
			SQLQueryConfig: &dmsCommonV1.SQLQueryConfig{},
		},
	}

	if _, err := convertInstance(inst); err == nil {
		t.Fatalf("convertInstance: expected ParseInt error for invalid RuleTemplateID, got nil")
	}
}
