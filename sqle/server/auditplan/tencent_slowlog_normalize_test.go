package auditplan

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// AC-018：MySQL QueryTime/LockTime 按秒写入 avg/max；显式 0 可写；缺字段无假 0 键。
func TestBuildTencentSQLV2_MySQLUnitsNoFakeZero(t *testing.T) {
	ap := &AuditPlan{
		ID:                  701,
		Type:                TypeTencentCdbMySQLSlowLog,
		InstanceAuditPlanId: 71,
		ProjectId:           "p",
		InstanceID:          "i",
	}
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)

	full := buildTencentSQLV2(ap, tencentSlowLogItem{
		SQLText:            "SELECT a FROM t WHERE id=1",
		SchemaName:         "db1",
		UserName:           "u1",
		UserHost:           "10.0.0.1",
		ExecutionStartTime: now,
		QueryTimeSec:       tencentFloat64(1.5),
		LockTimeSec:        tencentFloat64(0.2),
		RowsExamined:       tencentInt64(100),
		Md5:                "abc",
	}, true)
	assert.Equal(t, TypeTencentCdbMySQLSlowLog, full.Source)
	assert.InDelta(t, 1.5, full.Info.Get(MetricNameQueryTimeAvg).Float(), 1e-9)
	assert.InDelta(t, 1.5, full.Info.Get(MetricNameQueryTimeMax).Float(), 1e-9)
	assert.InDelta(t, 0.2, full.Info.Get(MetricNameLockTimeAvg).Float(), 1e-9)
	assert.InDelta(t, 0.2, full.Info.Get(MetricNameLockTimeMax).Float(), 1e-9)
	assert.InDelta(t, 100.0, full.Info.Get(MetricNameRowExaminedAvg).Float(), 1e-9)
	assert.Equal(t, "abc", full.Info.Get(MetricNameChecksum).String())

	// 显式云侧 0：允许写入 0（非伪造）
	zeroQT := buildTencentSQLV2(ap, tencentSlowLogItem{
		SQLText:      "SELECT 0",
		QueryTimeSec: tencentFloat64(0),
		LockTimeSec:  tencentFloat64(0),
	}, true)
	assert.NotNil(t, zeroQT.Info.Get(MetricNameQueryTimeAvg))
	assert.Equal(t, 0.0, zeroQT.Info.Get(MetricNameQueryTimeAvg).Float())
	assert.NotNil(t, zeroQT.Info.Get(MetricNameLockTimeAvg))
	assert.Equal(t, 0.0, zeroQT.Info.Get(MetricNameLockTimeAvg).Float())

	// 缺可选字段：Metrics 无对应键（禁止 JSON 缺省 0 → 假 0）
	missing := buildTencentSQLV2(ap, tencentSlowLogItem{
		SQLText: "SELECT missing",
	}, true)
	assert.Nil(t, missing.Info.Get(MetricNameQueryTimeAvg))
	assert.Nil(t, missing.Info.Get(MetricNameQueryTimeMax))
	assert.Nil(t, missing.Info.Get(MetricNameLockTimeAvg))
	assert.Nil(t, missing.Info.Get(MetricNameLockTimeMax))
	assert.Nil(t, missing.Info.Get(MetricNameRowExaminedAvg))
}

// AC-017：MySQL source 入队后等于类型常量，且不被覆盖为其它云类型。
func TestExtractTencentSlowLog_SourceConstants_MySQL(t *testing.T) {
	now := time.Date(2026, 11, 2, 15, 0, 0, 0, time.UTC)
	logger := logrus.NewEntry(logrus.New())

	var enqueued []*SQLV2
	var memCursor *tencentCollectCursor
	api := &stubMySQLAPI{
		resp: &tencentMySQLSlowLogResp{
			TotalCount: 1,
			Items: []tencentSlowLogItem{{
				SQLText:            "SELECT src_mysql",
				ExecutionStartTime: now.Add(-time.Minute),
				QueryTimeSec:       tencentFloat64(0.5),
			}},
		},
	}
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				memCursor = &tencentCollectCursor{
					SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt: at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue:  func(sqls []*SQLV2) error { enqueued = sqls; return nil },
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{
		ID: 711, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 711,
		ProjectId: "p", InstanceID: "i", Params: testTencentAPParams(),
	}
	out, err := task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, out)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, TypeTencentCdbMySQLSlowLog, enqueued[0].Source)
		assert.NotEqual(t, "aliyun_rds_mysql_slow_log", enqueued[0].Source)
		assert.NotEqual(t, "tencent_postgresql_slow_log", enqueued[0].Source)
		assert.InDelta(t, 0.5, enqueued[0].Info.Get(MetricNameQueryTimeAvg).Float(), 1e-9)
	}
}

func TestTencentMySQLTaskMetricsDeclareQueryTime(t *testing.T) {
	mysqlM := (&MySQLSlowLogTencentTaskV2{}).Metrics()
	assert.Contains(t, mysqlM, MetricNameQueryTimeAvg)
	assert.Contains(t, mysqlM, MetricNameLockTimeAvg)
	assert.Contains(t, mysqlM, MetricNameRowExaminedAvg)
}
