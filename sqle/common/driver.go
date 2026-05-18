package common

import (
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// dbTypeAliasForPlugin 把 DBService 表中可能存在的历史 / 显示用 db_type 字面量
// 归一化为 sqle-pg-plugin 注册的 PluginName（来自 SQLE_PG_PLUGIN_DB_TYPE 环境
// 变量）。
//
// 背景（Issue #2868 / fix-002 / compat-RISK-1）：
//   - dms-ee `internal/dms/pkg/constant/const.go: ParseDBType` 与
//     `internal/dataQuery/pkg/db/db_conn_ee.go: NewDBConnection` 中已分别落地
//     GaussDB 系列的别名归一化，但 sqle-ee `NewDriverManagerWithoutAudit` 调
//     `OpenPlugin(l, inst.DbType, cfg)` 时**直接传 inst.DbType 原始字符串**，
//     当 db_type 是 `GaussDB for MySQL`（ADR-004 历史命名）/ `openGauss` /
//     `GaussDB / openGauss`（展示名）时，plugin 注册名只有 `GaussDB` 与
//     `PostgreSQL` 两个，OpenPlugin lookup 失败 → 整条 schemas / audit / data
//     export 链路 `plugin not found`。
//   - 本 map 在 sqle-ee 内**自实现最小别名集合**而**不**新增对 dms-ee constant
//     包的跨仓库依赖（依赖 go.mod / vendor 调整会破坏本轮 fix 的零依赖硬约束）。
//   - 别名集合与 dms-ee ParseDBType (commit dms-ee@f00ee113 之后版本) 中
//     `case "GaussDB", "openGauss", "GaussDB / openGauss", "GaussDB for MySQL"`
//     的字面量集合**保持一致**；任何一侧扩展时另一侧也要同步。
//   - 归一化策略：先 `strings.ToLower(strings.TrimSpace(...))` 再查表，覆盖
//     `gaussdb for mysql` / `gaussdb` / `opengauss` / `gaussdb / opengauss`
//     四个等价键；查表未命中则**原样透传**，由 OpenPlugin 自己报错（不掩盖
//     真正的未知 db_type）。
//
// 未来 vendor 同步 dms-ee 新版 ParseDBType 后，可考虑改为直接 import
// `github.com/actiontech/dms/internal/dms/pkg/constant.ParseDBType` 并删除
// 本 map，但需保持 plugin 注册名（`GaussDB`）的目标值不变。
var dbTypeAliasForPlugin = map[string]string{
	"gaussdb for mysql":   "GaussDB", // ADR-004 历史命名
	"gaussdb":             "GaussDB", // canonical 自身回环
	"opengauss":           "GaussDB", // 同义别名
	"gaussdb / opengauss": "GaussDB", // dms-ee DBTypeGaussDB 常量值
	"postgresql":          "PostgreSQL", // 基线自回环（不影响行为）
}

// NormalizeDbTypeForPluginLookup 在调 plugin manager `OpenPlugin(name, cfg)` 前
// 把 db_type 字面量转成 plugin 注册名。查表未命中时**原样返回**输入字符串，
// 保证未知 / 未来新增 db_type 仍由 OpenPlugin 报真实的 `plugin not found`。
//
// fix-002（commit be6ce279）仅在 NewDriverManagerWithoutAudit 内调用，函数命名
// 为 unexported（同包私有）。fix-004（Issue #2868 / Bug F）补齐 audit 路径
// （sqle-ee/sqle/server/sqled.go: newDriverManagerWithAudit）与 offline 路径
// （本文件 NewDriverManagerWithoutCfg）的归一化时，audit 路径在 `server` 包
// 需要跨包调用本函数，故 export 首字母为大写。
//
// 详见 `docs/test/decision_round2_bugs.md §3` /
// `docs/test/decision_round4_bug_f.md §2.3` /
// `docs/test/fix-002_bug_a_round2_bug_c.md` /
// `docs/test/fix-004_bug_f_complete_normalization.md`。
func NormalizeDbTypeForPluginLookup(dbType string) string {
	key := strings.ToLower(strings.TrimSpace(dbType))
	if canonical, ok := dbTypeAliasForPlugin[key]; ok {
		return canonical
	}
	return dbType
}

// openPluginFunc 是 driver.GetPluginManager().OpenPlugin 的可测试间接层。
// 生产路径走默认实现；测试可在 setup 阶段替换为 mock 以验证 db_type 归一化
// 行为，无需真正启动 plugin manager。
var openPluginFunc = func(l *logrus.Entry, dbType string, cfg *v2.Config) (driver.Plugin, error) {
	return driver.GetPluginManager().OpenPlugin(l, dbType, cfg)
}

func NewDriverManagerWithoutAudit(l *logrus.Entry, inst *model.Instance, database string) (driver.Plugin, error) {
	if inst == nil {
		return nil, errors.Errorf("instance is nil")
	}

	dsn, err := NewDSN(inst, database)
	if err != nil {
		return nil, errors.Wrap(err, "new dsn")
	}

	cfg := &v2.Config{
		DSN: dsn,
	}
	// 把 inst.DbType 通过 NormalizeDbTypeForPluginLookup 归一化到 plugin 注册
	// 名再 lookup，覆盖 ADR-004 历史命名 "GaussDB for MySQL" / 展示名
	// "GaussDB / openGauss" / 同义别名 "openGauss" 三类历史 db_type，
	// 避免 OpenPlugin 直接报 `plugin not found`（Issue #2868 / compat-RISK-1 /
	// fix-002，详见 docs/test/decision_round2_bugs.md §3）。
	pluginName := NormalizeDbTypeForPluginLookup(inst.DbType)
	plugin, err := openPluginFunc(l, pluginName, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "open plugin")
	}

	return plugin, nil
}

func NewDriverManagerWithoutCfg(l *logrus.Entry, dbType string) (driver.Plugin, error) {
	// 与 NewDriverManagerWithoutAudit 对称：在调 OpenPlugin 前对 dbType 做归一
	// 化，覆盖 ADR-004 历史命名 "GaussDB for MySQL" / 展示名 "GaussDB /
	// openGauss" / 同义别名 "openGauss" 三类历史 db_type，避免 OpenPlugin 直接
	// 报 `plugin not found`（Issue #2868 / compat-RISK-1 / fix-004，详见
	// docs/test/decision_round4_bug_f.md §2.3 修复点 2）。
	//
	// 复用 openPluginFunc 包级 hook 与 fix-002 的 mock 模式一致，便于单测捕获
	// 实际传给 OpenPlugin 的字面量。
	pluginName := NormalizeDbTypeForPluginLookup(dbType)
	return openPluginFunc(l, pluginName, &v2.Config{})
}

func NewDSN(instance *model.Instance, database string) (*v2.DSN, error) {
	if instance == nil {
		return nil, errors.Errorf("instance is nil")
	}

	return &v2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     database,
	}, nil
}
