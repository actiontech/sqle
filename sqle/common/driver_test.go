package common

import (
	"testing"

	"github.com/actiontech/sqle/sqle/driver"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

// 5 个独立 Test 函数覆盖 Issue #2868 / fix-002 / compat-RISK-1 中
// sqle-ee NewDriverManagerWithoutAudit 在调 OpenPlugin 前对 db_type
// 做归一化的关键不变量。
//
// 为了不真正启动 plugin manager (会触发 go-plugin handshake + 子进程 spawn)，
// 测试通过替换 openPluginFunc 这个包内函数变量来捕获实际传给 OpenPlugin 的
// db_type 参数，并断言其字面量。
//
// 与 compat_risks.md compat-RISK-1 字段中 unit_tests 列表名字逐字对齐：
//   - TestNewDriverManagerWithoutAudit_GaussDBForMySQL_Routed
//   - TestNewDriverManagerWithoutAudit_OpenGauss_Routed
//   - TestNewDriverManagerWithoutAudit_GaussDB_Routed
//   - TestNewDriverManagerWithoutAudit_PostgreSQL_Unaffected
//   - TestNewDriverManagerWithoutAudit_UnknownDbType_PassThrough

// withMockOpenPlugin 替换 openPluginFunc，捕获 dbType 入参 + 返回 nil/nil；
// 返回值 restore 函数用于测试结束 defer 还原，避免污染其他 Test。
func withMockOpenPlugin(t *testing.T) (captured *string, restore func()) {
	t.Helper()
	var cap string
	original := openPluginFunc
	openPluginFunc = func(l *logrus.Entry, dbType string, cfg *v2.Config) (driver.Plugin, error) {
		cap = dbType
		return nil, nil
	}
	return &cap, func() { openPluginFunc = original }
}

// fakeInstance 构造最小可用 model.Instance；NewDSN 只校验非 nil，不真正连数据源。
func fakeInstance(dbType string) *model.Instance {
	return &model.Instance{
		Host:     "127.0.0.1",
		Port:     "5432",
		User:     "root",
		Password: "irrelevant",
		DbType:   dbType,
	}
}

// TestNewDriverManagerWithoutAudit_GaussDBForMySQL_Routed
// inst.DbType="GaussDB for MySQL" → OpenPlugin 收到 "GaussDB"（compat-RISK-1 决策 B + fix-002 a-2）。
func TestNewDriverManagerWithoutAudit_GaussDBForMySQL_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutAudit(l, fakeInstance("GaussDB for MySQL"), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutAudit_OpenGauss_Routed
// inst.DbType="openGauss" → OpenPlugin 收到 "GaussDB"。
func TestNewDriverManagerWithoutAudit_OpenGauss_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutAudit(l, fakeInstance("openGauss"), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutAudit_GaussDB_Routed
// inst.DbType="GaussDB" → OpenPlugin 收到 "GaussDB"（基线 / canonical 自身回环）。
func TestNewDriverManagerWithoutAudit_GaussDB_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutAudit(l, fakeInstance("GaussDB"), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutAudit_PostgreSQL_Unaffected
// inst.DbType="PostgreSQL" → OpenPlugin 收到 "PostgreSQL"（基线行为不破）。
func TestNewDriverManagerWithoutAudit_PostgreSQL_Unaffected(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutAudit(l, fakeInstance("PostgreSQL"), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "PostgreSQL" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"PostgreSQL\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutAudit_UnknownDbType_PassThrough
// inst.DbType="UnknownDb" → OpenPlugin 收到原字符串 "UnknownDb"（未知 db_type
// 不动，由 OpenPlugin 自己报 plugin not found，不掩盖真实错误）。
func TestNewDriverManagerWithoutAudit_UnknownDbType_PassThrough(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutAudit(l, fakeInstance("UnknownDb"), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "UnknownDb" {
		t.Errorf("OpenPlugin 接收到的 db_type 应原样透传 \"UnknownDb\"，got=%q", *got)
	}
}

// ---------------------------------------------------------------------------
// fix-004（Issue #2868 / Bug F）：NewDriverManagerWithoutCfg 在调 OpenPlugin
// 前对 dbType 做归一化的对仗单测，覆盖 ADR-004 历史命名 / 同义别名 / 展示名
// 三类 db_type 输入；mock hook openPluginFunc 捕获实际传给 OpenPlugin 的字面量。
// 与 compat_risks.md compat-RISK-1 字段中 unit_tests 列表名字逐字对齐：
//   - TestNewDriverManagerWithoutCfg_GaussDBForMySQL_Routed
//   - TestNewDriverManagerWithoutCfg_GaussDBLower_Routed
//   - TestNewDriverManagerWithoutCfg_OpenGaussLower_Routed
//   - TestNewDriverManagerWithoutCfg_GaussDBCanonical_Routed
//   - TestNewDriverManagerWithoutCfg_UnknownTypePassthrough_Routed
//
// 详见 docs/test/decision_round4_bug_f.md §2.3 修复点 2 与 §2.3 修复点 4。
// ---------------------------------------------------------------------------

// TestNewDriverManagerWithoutCfg_GaussDBForMySQL_Routed
// dbType="GaussDB for MySQL" → OpenPlugin 收到 "GaussDB"（ADR-004 历史命名）。
func TestNewDriverManagerWithoutCfg_GaussDBForMySQL_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutCfg(l, "GaussDB for MySQL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutCfg_GaussDBLower_Routed
// dbType="gaussdb" → OpenPlugin 收到 "GaussDB"（大小写不敏感归一化）。
func TestNewDriverManagerWithoutCfg_GaussDBLower_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutCfg(l, "gaussdb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutCfg_OpenGaussLower_Routed
// dbType="opengauss" → OpenPlugin 收到 "GaussDB"（同义别名 + 小写）。
func TestNewDriverManagerWithoutCfg_OpenGaussLower_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutCfg(l, "opengauss")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutCfg_GaussDBCanonical_Routed
// dbType="GaussDB" → OpenPlugin 收到 "GaussDB"（基线 / canonical 自身回环）。
func TestNewDriverManagerWithoutCfg_GaussDBCanonical_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutCfg(l, "GaussDB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "GaussDB" {
		t.Errorf("OpenPlugin 接收到的 db_type 应为 \"GaussDB\"，got=%q", *got)
	}
}

// TestNewDriverManagerWithoutCfg_UnknownTypePassthrough_Routed
// dbType="UnknownDb" → OpenPlugin 收到原字符串 "UnknownDb"（未知 db_type
// 不动，由 OpenPlugin 自己报 plugin not found，不掩盖真实错误）。
func TestNewDriverManagerWithoutCfg_UnknownTypePassthrough_Routed(t *testing.T) {
	got, restore := withMockOpenPlugin(t)
	defer restore()

	l := logrus.NewEntry(logrus.New())
	_, err := NewDriverManagerWithoutCfg(l, "UnknownDb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if *got != "UnknownDb" {
		t.Errorf("OpenPlugin 接收到的 db_type 应原样透传 \"UnknownDb\"，got=%q", *got)
	}
}
