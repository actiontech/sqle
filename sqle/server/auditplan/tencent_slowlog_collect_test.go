package auditplan

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type stubMySQLAPI struct {
	mu       sync.Mutex
	lastReq  tencentMySQLSlowLogReq
	reqs     []tencentMySQLSlowLogReq
	resp     *tencentMySQLSlowLogResp
	err      error
	callCnt  int
	// handler 优先：按 Offset 返回多页 / 注入第 N 页失败
	handler func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error)
}

func (s *stubMySQLAPI) DescribeSlowLogData(req tencentMySQLSlowLogReq) (*tencentMySQLSlowLogResp, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callCnt++
	s.lastReq = req
	s.reqs = append(s.reqs, req)
	if s.handler != nil {
		return s.handler(req, s.callCnt)
	}
	if s.err != nil {
		return nil, s.err
	}
	if s.resp == nil {
		return &tencentMySQLSlowLogResp{}, nil
	}
	return s.resp, nil
}

type stubPGAPI struct {
	mu      sync.Mutex
	lastReq tencentPGSlowLogReq
	reqs    []tencentPGSlowLogReq
	resp    *tencentPGSlowLogResp
	err     error
	callCnt int
	handler func(req tencentPGSlowLogReq, callIdx int) (*tencentPGSlowLogResp, error)
}

func (s *stubPGAPI) DescribeSlowQueryList(req tencentPGSlowLogReq) (*tencentPGSlowLogResp, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callCnt++
	s.lastReq = req
	s.reqs = append(s.reqs, req)
	if s.handler != nil {
		return s.handler(req, s.callCnt)
	}
	if s.err != nil {
		return nil, s.err
	}
	if s.resp == nil {
		return &tencentPGSlowLogResp{}, nil
	}
	return s.resp, nil
}

func testTencentAPParams() params.Params {
	return []*params.Param{
		{Key: paramKeyAccessKeyId, Value: "AKIDtest"},
		{Key: paramKeyAccessKeySecret, Value: "secrettest"},
		{Key: paramKeyRegion, Value: "ap-guangzhou"},
		{Key: paramKeyDBInstanceId, Value: "cdb-mysql-1"},
		{Key: paramKeyFirstSqlsScrappedInLastPeriodHours, Value: "24"},
		{Key: paramKeyCollectIntervalMinute, Value: "60"},
	}
}

func TestExtractTencentMySQLSlowLog_FirstWindowAndEnqueue(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{
		resp: &tencentMySQLSlowLogResp{
			TotalCount: 1,
			Items: []tencentSlowLogItem{{
				SQLText:            "SELECT 1",
				SchemaName:         "db1",
				ExecutionStartTime: now.Add(-10 * time.Minute),
				QueryTimeSec:       tencentFloat64(1.5),
			}},
		},
	}
	var savedCursor *tencentCollectCursor
	var enqueued []*SQLV2
	var memCursor *tencentCollectCursor

	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now: func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) {
				return memCursor, nil
			},
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				savedCursor = &tencentCollectCursor{
					SchemaVersion:        1,
					LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt:        at.UTC().Format(time.RFC3339),
				}
				memCursor = savedCursor
				return nil
			},
			enqueue: func(sqls []*SQLV2) error {
				enqueued = append(enqueued, sqls...)
				return nil
			},
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{
		ID:                  101,
		Type:                TypeTencentCdbMySQLSlowLog,
		InstanceAuditPlanId: 11,
		ProjectId:           "700300",
		InstanceID:          "inst-mysql",
		Params:              testTencentAPParams(),
	}
	logger := logrus.NewEntry(logrus.New())
	out, err := task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, out, "scheme A must return empty slice")
	assert.Equal(t, 1, api.callCnt)
	assert.Equal(t, now.Add(-24*time.Hour).Unix(), api.lastReq.StartTime)
	assert.Equal(t, now.Unix(), api.lastReq.EndTime)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, "SELECT 1", enqueued[0].SQLContent)
	}
	if assert.NotNil(t, savedCursor) {
		assert.Equal(t, now.Format(time.RFC3339), savedCursor.LastSuccessWindowEnd)
	}

	// 第二轮：start 来自游标（含 overlap），不得再次首次 24h
	enqueued = nil
	out, err = task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, out)
	expectedStart := now.Add(-5 * time.Minute).Unix() // lastEnd=now, overlap=5m；但 now 未变，第二轮 end 仍为 now
	// 第二轮 now 相同 → start = now - 5m
	assert.Equal(t, expectedStart, api.lastReq.StartTime)
	assert.NotEqual(t, now.Add(-24*time.Hour).Unix(), api.lastReq.StartTime)
}

func TestExtractTencentPGSlowLog_FirstWindowAndEnqueue(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubPGAPI{
		resp: &tencentPGSlowLogResp{
			TotalCount: 1,
			Items: []tencentSlowLogItem{{
				SQLText:            "SELECT 2",
				SchemaName:         "public",
				ExecutionStartTime: now.Add(-5 * time.Minute),
				DurationRaw:        "12.3",
			}},
		},
	}
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2

	ps := testTencentAPParams()
	ps = append(ps, &params.Param{Key: paramKeyDatabaseName, Value: "appdb"})
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"

	deps := tencentCollectDeps{
		now: func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) {
			return memCursor, nil
		},
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion:        1,
				LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt:        at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = sqls
			return nil
		},
		pgAPI: api,
	}
	ap := &AuditPlan{
		ID:                  202,
		Type:                "tencent_postgresql_slow_log",
		InstanceAuditPlanId: 22,
		ProjectId:           "700300",
		InstanceID:          "inst-pg",
		Params:              ps,
	}
	logger := logrus.NewEntry(logrus.New())
	out, err := extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Equal(t, 1, api.callCnt)
	assert.Equal(t, formatTencentPGTime(now.Add(-24*time.Hour)), api.lastReq.StartTime)
	assert.Equal(t, formatTencentPGTime(now), api.lastReq.EndTime)
	assert.Equal(t, "appdb", api.lastReq.DatabaseName)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, "SELECT 2", enqueued[0].SQLContent)
		assert.Equal(t, "12.3", enqueued[0].Info.Get(metricTencentPGDurationRaw).String())
		assert.Equal(t, tencentPGDurationUnitPending, enqueued[0].Info.Get(metricTencentPGDurationUnit).String())
		assert.Nil(t, enqueued[0].Info.Get(MetricNameQueryTimeAvg), "PG uncalibrated must not write query_time_*")
	}

	// 续采
	enqueued = nil
	_, err = extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	assert.Equal(t, formatTencentPGTime(now.Add(-5*time.Minute)), api.lastReq.StartTime)
}

func TestExtractTencentMySQLSlowLog_EnqueueFailDoesNotAdvanceCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{
		resp: &tencentMySQLSlowLogResp{
			Items: []tencentSlowLogItem{{SQLText: "SELECT 1", ExecutionStartTime: now}},
		},
	}
	saved := false
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return nil, nil },
			saveCur: func(uint, time.Time, time.Time) error {
				saved = true
				return nil
			},
			enqueue:  func([]*SQLV2) error { return assert.AnError },
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{ID: 1, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.Error(t, err)
	assert.False(t, saved)
}

func TestTencentPageExhausted(t *testing.T) {
	assert.True(t, tencentPageExhausted(0, 2, 0, 10))
	assert.True(t, tencentPageExhausted(1, 2, 0, 10))  // 不足 Limit
	assert.False(t, tencentPageExhausted(2, 2, 0, 5)) // 满页且未达 Total
	assert.True(t, tencentPageExhausted(2, 2, 2, 4))  // offset+len >= TotalCount
	assert.False(t, tencentPageExhausted(2, 2, 0, 0)) // TotalCount 未知，满页继续
}

// AC-009：MySQL ≥2 页翻页至结束，入队条数 = 各页 Items 合计
func TestExtractTencentMySQLSlowLog_MultiPageUntilExhausted(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	pageLimit := uint64(2)
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, _ int) (*tencentMySQLSlowLogResp, error) {
			assert.Equal(t, pageLimit, req.Limit)
			switch req.Offset {
			case 0:
				return &tencentMySQLSlowLogResp{
					TotalCount: 5,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT p1a", ExecutionStartTime: now.Add(-3 * time.Minute)},
						{SQLText: "SELECT p1b", ExecutionStartTime: now.Add(-2 * time.Minute)},
					},
				}, nil
			case 2:
				return &tencentMySQLSlowLogResp{
					TotalCount: 5,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT p2a", ExecutionStartTime: now.Add(-90 * time.Second)},
						{SQLText: "SELECT p2b", ExecutionStartTime: now.Add(-60 * time.Second)},
					},
				}, nil
			case 4:
				return &tencentMySQLSlowLogResp{
					TotalCount: 5,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT p3a", ExecutionStartTime: now.Add(-30 * time.Second)},
					},
				}, nil
			default:
				t.Fatalf("unexpected offset %d", req.Offset)
				return nil, nil
			}
		},
	}
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:            func() time.Time { return now },
			mysqlPageLimit: pageLimit,
			loadCur:        func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				memCursor = &tencentCollectCursor{
					SchemaVersion:        1,
					LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt:        at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue: func(sqls []*SQLV2) error {
				enqueued = append(enqueued, sqls...)
				return nil
			},
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{
		ID: 301, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 31,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	out, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Equal(t, 3, api.callCnt)
	assert.Equal(t, []uint64{0, 2, 4}, []uint64{api.reqs[0].Offset, api.reqs[1].Offset, api.reqs[2].Offset})
	assert.Len(t, enqueued, 5)
	assert.Equal(t, now.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-009：PG ≥2 页
func TestExtractTencentPGSlowLog_MultiPageUntilExhausted(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	pageLimit := uint64(2)
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, _ int) (*tencentPGSlowLogResp, error) {
			assert.Equal(t, pageLimit, req.Limit)
			switch req.Offset {
			case 0:
				return &tencentPGSlowLogResp{
					TotalCount: 3,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT g1", ExecutionStartTime: now.Add(-2 * time.Minute), DurationRaw: "1"},
						{SQLText: "SELECT g2", ExecutionStartTime: now.Add(-1 * time.Minute), DurationRaw: "2"},
					},
				}, nil
			case 2:
				return &tencentPGSlowLogResp{
					TotalCount: 3,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT g3", ExecutionStartTime: now.Add(-30 * time.Second), DurationRaw: "3"},
					},
				}, nil
			default:
				t.Fatalf("unexpected offset %d", req.Offset)
				return nil, nil
			}
		},
	}
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	ps := testTencentAPParams()
	ps = append(ps, &params.Param{Key: paramKeyDatabaseName, Value: "appdb"})
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"

	deps := tencentCollectDeps{
		now:         func() time.Time { return now },
		pgPageLimit: pageLimit,
		loadCur:     func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion:        1,
				LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt:        at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = sqls
			return nil
		},
		pgAPI: api,
	}
	ap := &AuditPlan{
		ID: 302, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 32,
		ProjectId: "700300", InstanceID: "inst-pg", Params: ps,
	}
	out, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Equal(t, 2, api.callCnt)
	assert.Len(t, enqueued, 3)
	assert.NotNil(t, memCursor)
}

// AC-010：第 N 页失败 → 游标不推进；下一轮成功可自旧游标恢复
func TestExtractTencentMySQLSlowLog_PageFailKeepsCursor_ThenRecover(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-1 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion:        1,
		LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt:        oldEnd.UTC().Format(time.RFC3339),
	}
	pageLimit := uint64(2)
	failOnSecond := true
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error) {
			if failOnSecond && callIdx == 2 {
				return nil, assert.AnError
			}
			if req.Offset == 0 {
				return &tencentMySQLSlowLogResp{
					TotalCount: 3,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT a", ExecutionStartTime: now.Add(-10 * time.Minute)},
						{SQLText: "SELECT b", ExecutionStartTime: now.Add(-9 * time.Minute)},
					},
				}, nil
			}
			return &tencentMySQLSlowLogResp{
				TotalCount: 3,
				Items: []tencentSlowLogItem{
					{SQLText: "SELECT c", ExecutionStartTime: now.Add(-8 * time.Minute)},
				},
			}, nil
		},
	}
	saveCnt := 0
	var enqueued []*SQLV2
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:            func() time.Time { return now },
			mysqlPageLimit: pageLimit,
			loadCur:        func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				saveCnt++
				memCursor = &tencentCollectCursor{
					SchemaVersion:        1,
					LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt:        at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue: func(sqls []*SQLV2) error {
				enqueued = append(enqueued, sqls...)
				return nil
			},
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{
		ID: 303, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 33,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	logger := logrus.NewEntry(logrus.New())

	_, err := task.ExtractSQL(logger, ap, nil)
	assert.Error(t, err)
	assert.Equal(t, 0, saveCnt)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	assert.Empty(t, enqueued)

	// 下轮：同旧游标起点恢复，翻页成功并推进
	failOnSecond = false
	api.callCnt = 0
	api.reqs = nil
	enqueued = nil
	_, err = task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, saveCnt)
	assert.Equal(t, now.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	assert.Len(t, enqueued, 3)
	// 续采起点 = 旧游标 − overlap
	assert.Equal(t, oldEnd.Add(-5*time.Minute).Unix(), api.reqs[0].StartTime)
}

// AC-010：入队失败不推进（多页拉完后）
func TestExtractTencentPGSlowLog_EnqueueFailAfterPagesKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-2 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion:        1,
		LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt:        oldEnd.UTC().Format(time.RFC3339),
	}
	pageLimit := uint64(2)
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, _ int) (*tencentPGSlowLogResp, error) {
			if req.Offset == 0 {
				return &tencentPGSlowLogResp{
					TotalCount: 3,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT x", ExecutionStartTime: now},
						{SQLText: "SELECT y", ExecutionStartTime: now},
					},
				}, nil
			}
			return &tencentPGSlowLogResp{
				TotalCount: 3,
				Items:      []tencentSlowLogItem{{SQLText: "SELECT z", ExecutionStartTime: now}},
			}, nil
		},
	}
	saved := false
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	deps := tencentCollectDeps{
		now:         func() time.Time { return now },
		pgPageLimit: pageLimit,
		loadCur:     func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(uint, time.Time, time.Time) error {
			saved = true
			return nil
		},
		enqueue: func([]*SQLV2) error { return assert.AnError },
		pgAPI:   api,
	}
	ap := &AuditPlan{ID: 304, Type: "tencent_postgresql_slow_log", Params: ps}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.False(t, saved)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	assert.Equal(t, 2, api.callCnt)
}

func TestTencentShrinkHelpers(t *testing.T) {
	assert.True(t, tencentMySQLNeedsShrinkByTotal(100, 20000))
	assert.False(t, tencentMySQLNeedsShrinkByTotal(100, 500))
	assert.True(t, tencentMySQLNextOffsetOverMax(9900, 100)) // 10000 > 9999
	assert.False(t, tencentMySQLNextOffsetOverMax(9800, 100))
	assert.True(t, tencentPGNeedsShrinkByTotal(10000))
	assert.False(t, tencentPGNeedsShrinkByTotal(9999))
	assert.True(t, isTencentResultOverLimitErr(fmt.Errorf("tencent api error code=ResultSetOverLimit request_id=x")))
	assert.True(t, isTencentResultOverLimitErr(fmt.Errorf("code=LogContentOverLimit")))
	assert.False(t, isTencentResultOverLimitErr(fmt.Errorf("code=AuthFailure")))

	start := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	assert.True(t, tencentWindowAtMin(start, start.Add(time.Second)))
	assert.False(t, tencentWindowAtMin(start, start.Add(2*time.Second)))
	mid, ok := splitTencentTimeWindow(start, start.Add(4*time.Second))
	assert.True(t, ok)
	assert.Equal(t, start.Add(2*time.Second), mid)
	_, ok = splitTencentTimeWindow(start, start.Add(time.Second))
	assert.False(t, ok)
}

// AC-012：MySQL TotalCount 使 Offset 将越 9999 → 缩窗后采完子窗并推进游标
func TestExtractTencentMySQLSlowLog_ShrinkOnOffsetCapThenComplete(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	pageLimit := uint64(100)
	var windowStarts []int64
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, _ int) (*tencentMySQLSlowLogResp, error) {
			windowStarts = append(windowStarts, req.StartTime)
			dur := req.EndTime - req.StartTime
			// 大窗：TotalCount 迫使末页 Offset>9999 → 缩窗
			if dur > 6*3600 {
				return &tencentMySQLSlowLogResp{
					TotalCount: 20000,
					Items: []tencentSlowLogItem{
						{SQLText: "SELECT big", ExecutionStartTime: time.Unix(req.StartTime, 0)},
					},
				}, nil
			}
			return &tencentMySQLSlowLogResp{
				TotalCount: 1,
				Items: []tencentSlowLogItem{
					{SQLText: fmt.Sprintf("SELECT w%d", req.StartTime), ExecutionStartTime: time.Unix(req.StartTime, 0)},
				},
			}, nil
		},
	}
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:            func() time.Time { return now },
			mysqlPageLimit: pageLimit,
			loadCur:        func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				memCursor = &tencentCollectCursor{
					SchemaVersion:        1,
					LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt:        at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue: func(sqls []*SQLV2) error {
				enqueued = append(enqueued, sqls...)
				return nil
			},
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{
		ID: 401, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 41,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	out, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Greater(t, api.callCnt, 2, "must issue sub-window calls after shrink")
	assert.GreaterOrEqual(t, len(enqueued), 2)
	assert.Equal(t, now.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	// 子窗 StartTime 不完全相同
	uniq := map[int64]struct{}{}
	for _, s := range windowStarts {
		uniq[s] = struct{}{}
	}
	assert.GreaterOrEqual(t, len(uniq), 2)
}

// AC-012：PG TotalCount≥10000 → 缩窗继续
func TestExtractTencentPGSlowLog_ShrinkOnTotalCountThenComplete(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	pageLimit := uint64(50)
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, _ int) (*tencentPGSlowLogResp, error) {
			start, _ := time.Parse("2006-01-02 15:04:05", req.StartTime)
			end, _ := time.Parse("2006-01-02 15:04:05", req.EndTime)
			if end.Sub(start) > 6*time.Hour {
				return &tencentPGSlowLogResp{TotalCount: 10000, Items: nil}, nil
			}
			return &tencentPGSlowLogResp{
				TotalCount: 1,
				Items: []tencentSlowLogItem{
					{SQLText: "SELECT pg-" + req.StartTime, ExecutionStartTime: start, DurationRaw: "1"},
				},
			}, nil
		},
	}
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	deps := tencentCollectDeps{
		now:         func() time.Time { return now },
		pgPageLimit: pageLimit,
		loadCur:     func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion:        1,
				LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt:        at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = sqls
			return nil
		},
		pgAPI: api,
	}
	ap := &AuditPlan{
		ID: 402, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 42,
		ProjectId: "700300", InstanceID: "inst-pg", Params: ps,
	}
	out, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Greater(t, api.callCnt, 2)
	assert.GreaterOrEqual(t, len(enqueued), 2)
	assert.Equal(t, now.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-012：MySQL 结果集超限错误走缩窗（非盲重试同窗）
func TestExtractTencentMySQLSlowLog_ShrinkOnResultSetOverLimit(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	seenLarge := 0
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, _ int) (*tencentMySQLSlowLogResp, error) {
			dur := req.EndTime - req.StartTime
			if dur > 6*3600 {
				seenLarge++
				return nil, fmt.Errorf("tencent api error code=ResultSetOverLimit request_id=r1")
			}
			return &tencentMySQLSlowLogResp{
				TotalCount: 1,
				Items:      []tencentSlowLogItem{{SQLText: "SELECT ok", ExecutionStartTime: time.Unix(req.StartTime, 0)}},
			}, nil
		},
	}
	var memCursor *tencentCollectCursor
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return nil, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				memCursor = &tencentCollectCursor{
					SchemaVersion:        1,
					LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt:        at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue:  func([]*SQLV2) error { return nil },
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{ID: 403, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, seenLarge, 1)
	assert.NotNil(t, memCursor)
	// 缩窗后应对更小窗发请求，而非同窗盲重试：存在更短时间窗请求
	var hasSmall bool
	for _, r := range api.reqs {
		if r.EndTime-r.StartTime <= 6*3600 {
			hasSmall = true
			break
		}
	}
	assert.True(t, hasSmall)
}

// AC-013：最小窗仍超限 → 异常且游标不推进（MySQL）
func TestExtractTencentMySQLSlowLog_MinWindowStillOverLimitKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-30 * time.Minute)
	memCursor := &tencentCollectCursor{
		SchemaVersion:        1,
		LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt:        oldEnd.UTC().Format(time.RFC3339),
	}
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, _ int) (*tencentMySQLSlowLogResp, error) {
			// 任意窗均触顶，逼近最小窗后失败
			return &tencentMySQLSlowLogResp{TotalCount: 20000, Items: nil}, nil
		},
	}
	saveCnt := 0
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:            func() time.Time { return now },
			mysqlPageLimit: 100,
			loadCur:        func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(uint, time.Time, time.Time) error {
				saveCnt++
				return nil
			},
			enqueue:  func([]*SQLV2) error { return nil },
			mysqlAPI: api,
		},
	}
	ap := &AuditPlan{ID: 404, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errTencentMinWindowExceeded))
	assert.Equal(t, 0, saveCnt)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-013：最小窗仍超限 → 异常且游标不推进（PG）
func TestExtractTencentPGSlowLog_MinWindowStillOverLimitKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-1 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion:        1,
		LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt:        oldEnd.UTC().Format(time.RFC3339),
	}
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, _ int) (*tencentPGSlowLogResp, error) {
			return &tencentPGSlowLogResp{TotalCount: 10000, Items: nil}, nil
		},
	}
	saved := false
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	deps := tencentCollectDeps{
		now:         func() time.Time { return now },
		pgPageLimit: 50,
		loadCur:     func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(uint, time.Time, time.Time) error {
			saved = true
			return nil
		},
		enqueue: func([]*SQLV2) error { return nil },
		pgAPI:   api,
	}
	ap := &AuditPlan{ID: 405, Type: "tencent_postgresql_slow_log", Params: ps}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errTencentMinWindowExceeded))
	assert.False(t, saved)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// memEventDedup 进程内幂等存储（fixture）
type memEventDedup struct {
	mu   sync.Mutex
	keys map[string]struct{}
}

func newMemEventDedup() *memEventDedup {
	return &memEventDedup{keys: map[string]struct{}{}}
}

func (m *memEventDedup) seen(_ uint, hash string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.keys[hash]
	return ok, nil
}

func (m *memEventDedup) commit(_ uint, hashes []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range hashes {
		m.keys[h] = struct{}{}
	}
	return nil
}

func (m *memEventDedup) len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.keys)
}

func sameMySQLEvent(now time.Time) tencentSlowLogItem {
	return tencentSlowLogItem{
		SQLText:            "SELECT dup FROM t",
		SchemaName:         "db1",
		UserName:           "u1",
		UserHost:           "1.2.3.4",
		ExecutionStartTime: now.Add(-10 * time.Minute),
		Md5:                "md5-same-event",
		QueryTimeSec:       tencentFloat64(1.0),
	}
}

// AC-011：同批重复事件 + 跨轮重叠 → 仅入队/计数一次（MySQL）
func TestExtractTencentMySQLSlowLog_EventDedupSkipDuplicate(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	ev := sameMySQLEvent(now)
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error) {
			// 同页返回两次同一事件，模拟分页/重叠边界重复
			return &tencentMySQLSlowLogResp{
				TotalCount: 2,
				Items:      []tencentSlowLogItem{ev, ev},
			}, nil
		},
	}
	dedup := newMemEventDedup()
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	deps := &tencentCollectDeps{
		now:         func() time.Time { return now },
		loadCur:     func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur:     func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = append(enqueued, sqls...)
			return nil
		},
		mysqlAPI:    api,
		dedupSeen:   dedup.seen,
		dedupCommit: dedup.commit,
		sleep:       func(time.Duration) {},
	}
	task := &MySQLSlowLogTencentTaskV2{collectDeps: deps}
	ap := &AuditPlan{ID: 501, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 51, ProjectId: "p", InstanceID: "i", Params: testTencentAPParams()}
	logger := logrus.NewEntry(logrus.New())

	_, err := task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, int64(1), enqueued[0].Info.Get(MetricNameCounter).Int())
	}
	assert.Equal(t, 1, dedup.len())

	// 第二轮再拉同一事件：跳过入队，counter 不再增加；游标仍推进
	enqueued = nil
	_, err = task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, enqueued)
	assert.Equal(t, 1, dedup.len())
	assert.NotNil(t, memCursor)
}

// AC-011：PG 同事件重复不增 counter
func TestExtractTencentPGSlowLog_EventDedupSkipDuplicate(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	ev := tencentSlowLogItem{
		SQLText:            "SELECT 1 FROM pg",
		SchemaName:         "public",
		UserName:           "pguser",
		UserHost:           "10.0.0.1",
		ExecutionStartTime: now.Add(-3 * time.Minute),
		SessionId:          "sess-1",
		ProcessId:          "1234",
		DurationRaw:        "1.2",
	}
	api := &stubPGAPI{
		resp: &tencentPGSlowLogResp{TotalCount: 2, Items: []tencentSlowLogItem{ev, ev}},
	}
	dedup := newMemEventDedup()
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	deps := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = append(enqueued, sqls...)
			return nil
		},
		pgAPI:       api,
		dedupSeen:   dedup.seen,
		dedupCommit: dedup.commit,
		sleep:       func(time.Duration) {},
	}
	ap := &AuditPlan{ID: 502, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 52, ProjectId: "p", InstanceID: "i", Params: ps}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.NoError(t, err)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, int64(1), enqueued[0].Info.Get(MetricNameCounter).Int())
	}

	enqueued = nil
	_, err = extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.NoError(t, err)
	assert.Empty(t, enqueued)
}

// AC-014：可重试错误有限退避后成功（MySQL）；游标推进
func TestExtractTencentMySQLSlowLog_RetryableBackoffThenSuccess(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error) {
			if callIdx < 3 {
				return nil, fmt.Errorf("tencent api error code=FailedOperation.TimeoutError request_id=r%d", callIdx)
			}
			return &tencentMySQLSlowLogResp{
				TotalCount: 1,
				Items: []tencentSlowLogItem{{
					SQLText: "SELECT after retry", ExecutionStartTime: now.Add(-1 * time.Minute), Md5: "m1",
				}},
			}, nil
		},
	}
	var sleeps []time.Duration
	var memCursor *tencentCollectCursor
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return nil, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				memCursor = &tencentCollectCursor{
					SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt: at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue:        func([]*SQLV2) error { return nil },
			mysqlAPI:       api,
			sleep:          func(d time.Duration) { sleeps = append(sleeps, d) },
			maxAPIAttempts: 5,
		},
	}
	ap := &AuditPlan{ID: 503, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, api.callCnt)
	assert.Len(t, sleeps, 2)
	for _, d := range sleeps {
		assert.True(t, d >= tencentAPIBackoffBase/2 && d <= tencentAPIBackoffCap)
	}
	assert.NotNil(t, memCursor)
}

// AC-014：鉴权类立即失败，不进入退避；游标不推进（MySQL）
func TestExtractTencentMySQLSlowLog_AuthFailureNoRetryKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-40 * time.Minute)
	memCursor := &tencentCollectCursor{
		SchemaVersion: 1, LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt: oldEnd.UTC().Format(time.RFC3339),
	}
	api := &stubMySQLAPI{
		err: fmt.Errorf("tencent api error code=AuthFailure.SecretIdNotFound request_id=x"),
	}
	var sleeps []time.Duration
	saveCnt := 0
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(uint, time.Time, time.Time) error { saveCnt++; return nil },
			enqueue: func([]*SQLV2) error { return nil },
			mysqlAPI: api,
			sleep:    func(d time.Duration) { sleeps = append(sleeps, d) },
		},
	}
	ap := &AuditPlan{ID: 504, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AuthFailure")
	assert.Equal(t, 1, api.callCnt, "auth must not blind-retry")
	assert.Empty(t, sleeps)
	assert.Equal(t, 0, saveCnt)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-014：PG CAM 拒绝立即失败；可重试限流后成功
func TestExtractTencentPGSlowLog_CamDeniedNoRetry_And_LimitRetry(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	oldEnd := now.Add(-2 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion: 1, LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
		LastSuccessAt: oldEnd.UTC().Format(time.RFC3339),
	}

	camAPI := &stubPGAPI{err: fmt.Errorf("tencent api error code=OperationDenied.CamDeniedError request_id=c1")}
	var sleeps []time.Duration
	deps := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(uint, time.Time, time.Time) error { return errors.New("should not save") },
		enqueue: func([]*SQLV2) error { return nil },
		pgAPI:   camAPI,
		sleep:   func(d time.Duration) { sleeps = append(sleeps, d) },
	}
	ap := &AuditPlan{ID: 505, Type: "tencent_postgresql_slow_log", Params: ps}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CamDeniedError")
	assert.Equal(t, 1, camAPI.callCnt)
	assert.Empty(t, sleeps)

	limitAPI := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, callIdx int) (*tencentPGSlowLogResp, error) {
			if callIdx == 1 {
				return nil, fmt.Errorf("tencent api error code=RequestLimitExceeded request_id=l1")
			}
			return &tencentPGSlowLogResp{
				TotalCount: 1,
				Items: []tencentSlowLogItem{{
					SQLText: "SELECT ok", ExecutionStartTime: now.Add(-time.Minute),
					SessionId: "s", ProcessId: "1",
				}},
			}, nil
		},
	}
	var saved *tencentCollectCursor
	deps2 := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return nil, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			saved = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func([]*SQLV2) error { return nil },
		pgAPI:   limitAPI,
		sleep:   func(time.Duration) {},
	}
	_, err = extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps2)
	assert.NoError(t, err)
	assert.Equal(t, 2, limitAPI.callCnt)
	assert.NotNil(t, saved)
}

func TestTencentAPIErrorClassifyAndBackoff(t *testing.T) {
	assert.True(t, isTencentFatalAPIErr(fmt.Errorf("code=AuthFailure.SignatureFailure")))
	assert.True(t, isTencentFatalAPIErr(fmt.Errorf("code=ResourceNotFound.InstanceNotFoundError")))
	assert.True(t, isTencentRetryableAPIErr(fmt.Errorf("code=FailedOperation.TimeoutError")))
	assert.True(t, isTencentRetryableAPIErr(fmt.Errorf("tencent api http status 503")))
	assert.True(t, isTencentRetryableAPIErr(fmt.Errorf("tencent api request failed: i/o timeout")))
	assert.False(t, isTencentRetryableAPIErr(fmt.Errorf("code=AuthFailure.SecretIdNotFound")))
	assert.False(t, isTencentRetryableAPIErr(fmt.Errorf("code=ResultSetOverLimit")))

	d := tencentAPIBackoffDelay(0, tencentAPIBackoffBase, tencentAPIBackoffCap, rand.New(rand.NewSource(1)))
	assert.True(t, d >= tencentAPIBackoffBase/2 && d <= tencentAPIBackoffBase)

	attempts := 0
	err := withTencentAPIRetry(func(time.Duration) {}, 3, nil, func() error {
		attempts++
		return fmt.Errorf("tencent api error code=RequestLimitExceeded")
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "retry exhausted")
	assert.Equal(t, 3, attempts)

	attempts = 0
	err = withTencentAPIRetry(func(time.Duration) {}, 5, nil, func() error {
		attempts++
		return fmt.Errorf("tencent api error code=CamSigAndAuthError")
	})
	assert.Error(t, err)
	assert.Equal(t, 1, attempts)

	assert.NotEqual(t, "", tencentMySQLEventKeyHash("cdb-1", sameMySQLEvent(time.Now())))
	k1 := tencentPGEventKeyHash(nil, "pg-1", tencentSlowLogItem{
		SQLText: "SELECT 1", ExecutionStartTime: time.Unix(1, 0), SessionId: "a", ProcessId: "b",
	})
	k2 := tencentPGEventKeyHash(nil, "pg-1", tencentSlowLogItem{
		SQLText: "SELECT 1", ExecutionStartTime: time.Unix(1, 0), SessionId: "a", ProcessId: "b",
	})
	assert.Equal(t, k1, k2)
}

// AC-015：分页中途禁用 → 之后无新 API，游标不推进。
func TestExtractTencentMySQLSlowLog_DisableMidPaginationStopsAPIKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-2 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion: 1, LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
	}
	disabled := false
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error) {
			if callIdx == 1 {
				disabled = true // 模拟第一页后任务被禁用
				return &tencentMySQLSlowLogResp{
					TotalCount: 3,
					Items: []tencentSlowLogItem{{
						SQLText: "SELECT page1", SchemaName: "db1",
						ExecutionStartTime: now.Add(-30 * time.Minute), QueryTimeSec: tencentFloat64(1),
					}},
				}, nil
			}
			t.Fatalf("unexpected API call after disable: callIdx=%d", callIdx)
			return nil, nil
		},
	}
	saveCnt := 0
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:     func() time.Time { return now },
			loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(uint, time.Time, time.Time) error { saveCnt++; return nil },
			enqueue: func([]*SQLV2) error { t.Fatal("must not enqueue after disable"); return nil },
			mysqlAPI: api,
			mysqlPageLimit: 1,
			shouldStop:     func() bool { return disabled },
			sleep:          func(time.Duration) {},
		},
	}
	ap := &AuditPlan{
		ID: 201, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 21,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, 1, api.callCnt, "no new API after disable")
	assert.Equal(t, 0, saveCnt)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-015：PG 退避循环中禁用 → 停止新 API，游标不推进。
func TestExtractTencentPGSlowLog_DisableDuringRetryStopsAPIKeepsCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	oldEnd := now.Add(-3 * time.Hour)
	memCursor := &tencentCollectCursor{
		SchemaVersion: 1, LastSuccessWindowEnd: oldEnd.UTC().Format(time.RFC3339),
	}
	disabled := false
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, callIdx int) (*tencentPGSlowLogResp, error) {
			if callIdx == 1 {
				disabled = true
				return nil, fmt.Errorf("tencent api error code=RequestLimitExceeded")
			}
			t.Fatalf("unexpected API call after disable: callIdx=%d", callIdx)
			return nil, nil
		},
	}
	saveCnt := 0
	pgParams := testTencentAPParams()
	pgParams = append(pgParams, &params.Param{Key: paramKeyDatabaseName, Value: "postgres"})
	deps := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(uint, time.Time, time.Time) error { saveCnt++; return nil },
		enqueue: func([]*SQLV2) error { return nil },
		pgAPI:      api,
		shouldStop: func() bool { return disabled },
		sleep:      func(time.Duration) {},
	}
	ap := &AuditPlan{
		ID: 202, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 22,
		ProjectId: "700300", InstanceID: "inst-pg", Params: pgParams,
	}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, 1, api.callCnt)
	assert.Equal(t, 0, saveCnt)
	assert.Equal(t, oldEnd.UTC().Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
}

// AC-015：Extract 入口已禁用 → 零次 API。
func TestExtractTencentMySQLSlowLog_AlreadyDisabledZeroAPI(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{resp: &tencentMySQLSlowLogResp{TotalCount: 0}}
	saveCnt := 0
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:        func() time.Time { return now },
			loadCur:    func(uint) (*tencentCollectCursor, error) { return nil, nil },
			saveCur:    func(uint, time.Time, time.Time) error { saveCnt++; return nil },
			enqueue:    func([]*SQLV2) error { return nil },
			mysqlAPI:   api,
			shouldStop: func() bool { return true },
		},
	}
	ap := &AuditPlan{
		ID: 203, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 23,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, 0, api.callCnt)
	assert.Equal(t, 0, saveCnt)
}

// AC-016：连续成功轮次游标推进到声明 end，重叠窗下覆盖无间隙。
func TestExtractTencentMySQLSlowLog_SuccessRoundsCursorNoGap(t *testing.T) {
	clock := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	var memCursor *tencentCollectCursor
	var declaredEnds []time.Time
	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, _ int) (*tencentMySQLSlowLogResp, error) {
			return &tencentMySQLSlowLogResp{TotalCount: 0, Items: nil}, nil
		},
	}
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now: func() time.Time { return clock },
			loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
			saveCur: func(_ uint, windowEnd, at time.Time) error {
				declaredEnds = append(declaredEnds, windowEnd)
				memCursor = &tencentCollectCursor{
					SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
					LastSuccessAt: at.UTC().Format(time.RFC3339),
				}
				return nil
			},
			enqueue:  func([]*SQLV2) error { return nil },
			mysqlAPI: api,
			overlap:  5 * time.Minute,
		},
	}
	ap := &AuditPlan{
		ID: 301, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 31,
		ProjectId: "700300", InstanceID: "inst-mysql", Params: testTencentAPParams(),
	}
	logger := logrus.NewEntry(logrus.New())

	// round1：首次 → cursor = end1
	_, err := task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	end1 := clock
	assert.Equal(t, end1.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	assert.Equal(t, end1.Add(-24*time.Hour).Unix(), api.reqs[0].StartTime)
	assert.Equal(t, end1.Unix(), api.reqs[0].EndTime)

	// round2：now 前进 1h；start = end1 - overlap，cursor → end2；与上一轮声明范围无间隙
	clock = clock.Add(time.Hour)
	api.reqs = nil
	_, err = task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	end2 := clock
	assert.Equal(t, end2.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	start2 := end1.Add(-5 * time.Minute)
	assert.Equal(t, start2.Unix(), api.reqs[0].StartTime)
	assert.Equal(t, end2.Unix(), api.reqs[0].EndTime)
	// 无间隙：上一轮上界 end1 落在本轮 [start2, end2] 内（重叠窗）
	assert.False(t, end1.Before(start2), "gap: previous end before next start")
	assert.True(t, !end1.After(end2))
	assert.Equal(t, []time.Time{end1, end2}, declaredEnds)
}

// AC-016：PG 成功推进上界；禁用中断保持旧游标后再成功续采无间隙。
func TestExtractTencentPGSlowLog_CursorAdvanceNoGap_AfterDisableInterrupt(t *testing.T) {
	clock := time.Date(2026, 11, 2, 8, 0, 0, 0, time.UTC)
	var memCursor *tencentCollectCursor
	pgParams := testTencentAPParams()
	pgParams = append(pgParams, &params.Param{Key: paramKeyDatabaseName, Value: "postgres"})
	api := &stubPGAPI{
		handler: func(req tencentPGSlowLogReq, _ int) (*tencentPGSlowLogResp, error) {
			return &tencentPGSlowLogResp{TotalCount: 0}, nil
		},
	}
	disabled := false
	deps := tencentCollectDeps{
		now:     func() time.Time { return clock },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue:    func([]*SQLV2) error { return nil },
		pgAPI:      api,
		shouldStop: func() bool { return disabled },
		overlap:    5 * time.Minute,
	}
	ap := &AuditPlan{
		ID: 302, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 32,
		ProjectId: "700300", InstanceID: "inst-pg", Params: pgParams,
	}
	logger := logrus.NewEntry(logrus.New())

	_, err := extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	end1 := clock
	assert.Equal(t, end1.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)

	// 中断轮：禁用 → 游标不变
	clock = clock.Add(30 * time.Minute)
	disabled = true
	api.callCnt = 0
	_, err = extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, 0, api.callCnt)
	assert.Equal(t, end1.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)

	// 恢复：从旧游标续采，推进到新 end，与 end1 无间隙
	disabled = false
	clock = clock.Add(30 * time.Minute)
	api.reqs = nil
	_, err = extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	end3 := clock
	assert.Equal(t, end3.Format(time.RFC3339), memCursor.LastSuccessWindowEnd)
	start3 := end1.Add(-5 * time.Minute)
	assert.Equal(t, formatTencentPGTime(start3), api.reqs[0].StartTime)
	assert.False(t, end1.Before(start3))
}

func TestWithTencentAPIRetry_StopsOnShouldStopBetweenAttempts(t *testing.T) {
	attempts := 0
	stopAfter := false
	err := withTencentAPIRetry(func(time.Duration) { stopAfter = true }, 5, func() bool { return stopAfter }, func() error {
		attempts++
		return fmt.Errorf("tencent api error code=RequestLimitExceeded")
	})
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, 1, attempts, "must not retry after stop")
}

// AC-021：两事件同指纹 → AggregateSQL counter=2；dedup 行数=2；同事件重放不二次入队、counter 不膨胀。
func TestExtractTencentMySQLSlowLog_AC021_TwoEventsSameFingerprintCounterAndReplay(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	sameSQL := "SELECT id FROM users WHERE name = 'ac021'"
	ev1 := tencentSlowLogItem{
		SQLText:            sameSQL,
		SchemaName:         "db1",
		UserName:           "u1",
		UserHost:           "1.2.3.4",
		ExecutionStartTime: now.Add(-20 * time.Minute),
		Md5:                "md5-event-a",
		QueryTimeSec:       tencentFloat64(0.5),
	}
	ev2 := tencentSlowLogItem{
		SQLText:            sameSQL, // 同指纹，不同事件键（Md5 + 时间）
		SchemaName:         "db1",
		UserName:           "u1",
		UserHost:           "1.2.3.4",
		ExecutionStartTime: now.Add(-10 * time.Minute),
		Md5:                "md5-event-b",
		QueryTimeSec:       tencentFloat64(0.8),
	}
	assert.NotEqual(t, tencentMySQLEventKeyHash("cdb-mysql-1", ev1), tencentMySQLEventKeyHash("cdb-mysql-1", ev2))

	api := &stubMySQLAPI{
		handler: func(req tencentMySQLSlowLogReq, callIdx int) (*tencentMySQLSlowLogResp, error) {
			return &tencentMySQLSlowLogResp{TotalCount: 2, Items: []tencentSlowLogItem{ev1, ev2}}, nil
		},
	}
	dedup := newMemEventDedup()
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	deps := &tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = append(enqueued, sqls...)
			return nil
		},
		mysqlAPI:    api,
		dedupSeen:   dedup.seen,
		dedupCommit: dedup.commit,
		sleep:       func(time.Duration) {},
	}
	task := &MySQLSlowLogTencentTaskV2{collectDeps: deps}
	ap := &AuditPlan{ID: 921, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 91, ProjectId: "p", InstanceID: "i", Params: testTencentAPParams()}
	logger := logrus.NewEntry(logrus.New())

	_, err := task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	if assert.Len(t, enqueued, 1, "two events same fingerprint must merge to one SQLV2") {
		assert.Equal(t, int64(2), enqueued[0].Info.Get(MetricNameCounter).Int(), "fingerprint layer counter must reflect 2 real executions")
	}
	assert.Equal(t, 2, dedup.len(), "dedup rows must equal distinct event keys")

	// 同事件重放 / 重叠窗：幂等命中 → 不二次入队 → counter 不膨胀；dedup 行数不变
	prevCounter := enqueued[0].Info.Get(MetricNameCounter).Int()
	enqueued = nil
	_, err = task.ExtractSQL(logger, ap, nil)
	assert.NoError(t, err)
	assert.Empty(t, enqueued, "replay must not enqueue again")
	assert.Equal(t, 2, dedup.len())
	assert.Equal(t, int64(2), prevCounter, "counter must not inflate on event replay")
}

// AC-021：PG 两事件同指纹 counter=2；重放不增。
func TestExtractTencentPGSlowLog_AC021_TwoEventsSameFingerprintCounterAndReplay(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	sameSQL := "SELECT * FROM orders WHERE id = $1"
	ev1 := tencentSlowLogItem{
		SQLText:            sameSQL,
		SchemaName:         "public",
		UserName:           "pguser",
		UserHost:           "10.0.0.1",
		ExecutionStartTime: now.Add(-8 * time.Minute),
		SessionId:          "sess-a",
		ProcessId:          "1001",
		DurationRaw:        "1.1",
	}
	ev2 := tencentSlowLogItem{
		SQLText:            sameSQL,
		SchemaName:         "public",
		UserName:           "pguser",
		UserHost:           "10.0.0.1",
		ExecutionStartTime: now.Add(-4 * time.Minute),
		SessionId:          "sess-b",
		ProcessId:          "1002",
		DurationRaw:        "2.2",
	}
	logger := logrus.NewEntry(logrus.New())
	assert.NotEqual(t, tencentPGEventKeyHash(logger, "postgres-1", ev1), tencentPGEventKeyHash(logger, "postgres-1", ev2))

	api := &stubPGAPI{
		resp: &tencentPGSlowLogResp{TotalCount: 2, Items: []tencentSlowLogItem{ev1, ev2}},
	}
	dedup := newMemEventDedup()
	var memCursor *tencentCollectCursor
	var enqueued []*SQLV2
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"
	deps := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return memCursor, nil },
		saveCur: func(_ uint, windowEnd, at time.Time) error {
			memCursor = &tencentCollectCursor{
				SchemaVersion: 1, LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
				LastSuccessAt: at.UTC().Format(time.RFC3339),
			}
			return nil
		},
		enqueue: func(sqls []*SQLV2) error {
			enqueued = append(enqueued, sqls...)
			return nil
		},
		pgAPI:       api,
		dedupSeen:   dedup.seen,
		dedupCommit: dedup.commit,
		sleep:       func(time.Duration) {},
	}
	ap := &AuditPlan{ID: 922, Type: "tencent_postgresql_slow_log", InstanceAuditPlanId: 92, ProjectId: "p", InstanceID: "i", Params: ps}

	_, err := extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	if assert.Len(t, enqueued, 1) {
		assert.Equal(t, int64(2), enqueued[0].Info.Get(MetricNameCounter).Int())
	}
	assert.Equal(t, 2, dedup.len())

	enqueued = nil
	_, err = extractTencentPGSlowLog(logger, ap, nil, deps)
	assert.NoError(t, err)
	assert.Empty(t, enqueued)
	assert.Equal(t, 2, dedup.len())
}

// AC-021：filterAndAggregate 层直接断言（同批两事件同指纹 + 同事件跳过）。
func TestFilterAndAggregateTencentItems_AC021_CounterVsDedup(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	sameSQL := "SELECT x FROM t WHERE a = 1"
	items := []tencentSlowLogItem{
		{SQLText: sameSQL, SchemaName: "db", UserName: "u", UserHost: "h", ExecutionStartTime: now.Add(-3 * time.Minute), Md5: "k1", QueryTimeSec: tencentFloat64(1)},
		{SQLText: sameSQL, SchemaName: "db", UserName: "u", UserHost: "h", ExecutionStartTime: now.Add(-2 * time.Minute), Md5: "k2", QueryTimeSec: tencentFloat64(1)},
		{SQLText: sameSQL, SchemaName: "db", UserName: "u", UserHost: "h", ExecutionStartTime: now.Add(-3 * time.Minute), Md5: "k1", QueryTimeSec: tencentFloat64(1)}, // 同事件重放
	}
	dedup := newMemEventDedup()
	cache := NewSQLV2Cache()
	ap := &AuditPlan{ID: 923, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 93, ProjectId: "p", InstanceID: "i"}
	pending, err := filterAndAggregateTencentItems(
		logrus.NewEntry(logrus.New()), ap, "cdb-1", items, true,
		tencentCollectDeps{dedupSeen: dedup.seen, dedupCommit: dedup.commit},
		cache, &DefaultTaskV2{},
	)
	assert.NoError(t, err)
	assert.Len(t, pending, 2, "third item is same event key → skipped before aggregate")
	sqls := cache.GetSQLs()
	if assert.Len(t, sqls, 1) {
		assert.Equal(t, int64(2), sqls[0].Info.Get(MetricNameCounter).Int())
	}
	assert.NoError(t, dedup.commit(ap.ID, pending))
	assert.Equal(t, 2, dedup.len())
}
