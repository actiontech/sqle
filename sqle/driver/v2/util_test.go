package driverV2

import (
	"testing"
)

// TestDriverTypeOpenGauss_Const_EE 断言 sqle-ee 副本中 DriverTypeOpenGauss 常量字面值
// 严格为 "openGauss"（首字母小写 o，G 大写），并回归断言 DriverTypeGaussDB == "GaussDB"，
// 保证两个 GaussDB / openGauss 字面值与 design §4.4 严格匹配契约一致。
//
// 任何对该字面值的归一化（strings.ToLower / trim / 大小写改写）都会让本测试 fail，
// 提示 code_review 阶段同步落地的 sqle CE 副本必须保持完全一致。
func TestDriverTypeOpenGauss_Const_EE(t *testing.T) {
	cases := map[string]struct {
		got  string
		want string
	}{
		"openGauss literal":          {got: DriverTypeOpenGauss, want: "openGauss"},
		"GaussDB literal regression": {got: DriverTypeGaussDB, want: "GaussDB"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("constant literal drift: got %q, want %q", tc.got, tc.want)
			}
		})
	}
}
