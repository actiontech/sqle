package auditplan

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestClassifyTencentCollectError_StableCodes(t *testing.T) {
	cases := []struct {
		name string
		err  error
		code string
	}{
		{"auth_secret", fmt.Errorf("tencent api error code=AuthFailure.SecretIdNotFound request_id=x"), TencentErrCodeAuthInvalid},
		{"auth_sig", fmt.Errorf("tencent api error code=AuthFailure.SignatureFailure request_id=x"), TencentErrCodeAuthInvalid},
		{"auth_cam_sig", fmt.Errorf("tencent api error code=CamSigAndAuthError request_id=x"), TencentErrCodeAuthInvalid},
		{"cam_denied", fmt.Errorf("tencent api error code=OperationDenied.CamDeniedError request_id=c1"), TencentErrCodeCamDenied},
		{"cam_instance_access", fmt.Errorf("tencent api error code=OperationDenied.InstanceAccessDeniedError"), TencentErrCodeCamDenied},
		{"region", fmt.Errorf("tencent api error code=InvalidParameterValue.Region request_id=r1"), TencentErrCodeRegionMismatch},
		{"instance_cdb", fmt.Errorf("tencent api error code=ResourceNotFound.CdbInstanceNotFoundError"), TencentErrCodeInstanceNotFound},
		{"instance_pg", fmt.Errorf("tencent api error code=ResourceNotFound.InstanceNotFoundError"), TencentErrCodeInstanceNotFound},
		{"rate", fmt.Errorf("tencent api error code=RequestLimitExceeded"), TencentErrCodeRateLimited},
		{"config", fmt.Errorf("access_key_secret is not configured"), TencentErrCodeConfigInvalid},
		{"page", errTencentMinWindowExceeded, TencentErrCodePageOrParse},
		{"over_limit", fmt.Errorf("code=ResultSetOverLimit"), TencentErrCodePageOrParse},
		{"network", fmt.Errorf("tencent api request failed: i/o timeout"), TencentErrCodeNetworkOrService},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg := classifyTencentCollectError(c.err)
			assert.Equal(t, c.code, code)
			assert.NotEmpty(t, msg)
			// 允许文案提及参数名 SecretId/SecretKey；禁止密钥值 / 签名串 / Authorization 头
			assert.NotContains(t, msg, "SecretKey=")
			assert.NotContains(t, msg, "Authorization")
			assert.NotContains(t, msg, "TC3-HMAC")
			assert.NotContains(t, msg, "request_id")
		})
	}
}

func TestWrapTencentCollectError_ContainsCodeNoSecrets(t *testing.T) {
	raw := fmt.Errorf("tencent api error code=AuthFailure.SecretIdNotFound request_id=abc SecretKey=AKSECRET123 Authorization=TC3-HMAC-SHA256/xyz")
	wrapped := wrapTencentCollectError(raw)
	assert.Error(t, wrapped)
	assert.Equal(t, TencentErrCodeAuthInvalid, tencentCollectErrorCode(wrapped))
	s := wrapped.Error()
	assert.Contains(t, s, "["+TencentErrCodeAuthInvalid+"]")
	assert.Contains(t, s, "鉴权失败")
	assert.NotContains(t, s, "AKSECRET123")
	assert.NotContains(t, s, "TC3-HMAC-SHA256/xyz")
	// unwrap 仍可匹配致命类重试逻辑所用原文码
	assert.True(t, errors.Is(wrapped, raw) || strings.Contains(errors.Unwrap(wrapped).Error(), "AuthFailure"))
}

func TestWrapTencentCollectError_PreservesStoppedSentinel(t *testing.T) {
	err := wrapTencentCollectError(errTencentCollectStopped)
	assert.True(t, errors.Is(err, errTencentCollectStopped))
	assert.Equal(t, errTencentCollectStopped, err)
}

func TestExtractTencentMySQL_WrapsAuthInvalid(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{
		err: fmt.Errorf("tencent api error code=AuthFailure.SecretIdNotFound request_id=x"),
	}
	task := &MySQLSlowLogTencentTaskV2{
		collectDeps: &tencentCollectDeps{
			now:      func() time.Time { return now },
			loadCur:  func(uint) (*tencentCollectCursor, error) { return nil, nil },
			saveCur:  func(uint, time.Time, time.Time) error { return errors.New("should not save") },
			enqueue:  func([]*SQLV2) error { return nil },
			mysqlAPI: api,
			sleep:    func(time.Duration) {},
		},
	}
	ap := &AuditPlan{ID: 801, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := task.ExtractSQL(logrus.NewEntry(logrus.New()), ap, nil)
	assert.Error(t, err)
	assert.Equal(t, TencentErrCodeAuthInvalid, tencentCollectErrorCode(err))
	assert.Contains(t, err.Error(), TencentErrCodeAuthInvalid)
	assert.NotContains(t, err.Error(), "SecretKey=")
}

func TestExtractTencentPG_WrapsCamAndInstance(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	ps := testTencentAPParams()
	ps.GetParam(paramKeyDBInstanceId).Value = "postgres-1"

	camAPI := &stubPGAPI{err: fmt.Errorf("tencent api error code=OperationDenied.CamDeniedError request_id=c1")}
	deps := tencentCollectDeps{
		now:     func() time.Time { return now },
		loadCur: func(uint) (*tencentCollectCursor, error) { return nil, nil },
		saveCur: func(uint, time.Time, time.Time) error { return errors.New("no") },
		enqueue: func([]*SQLV2) error { return nil },
		pgAPI:   camAPI,
		sleep:   func(time.Duration) {},
	}
	ap := &AuditPlan{ID: 802, Type: "tencent_postgresql_slow_log", Params: ps}
	_, err := extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.Equal(t, TencentErrCodeCamDenied, tencentCollectErrorCode(err))

	nfAPI := &stubPGAPI{err: fmt.Errorf("tencent api error code=ResourceNotFound.InstanceNotFoundError request_id=n1")}
	deps.pgAPI = nfAPI
	_, err = extractTencentPGSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.Equal(t, TencentErrCodeInstanceNotFound, tencentCollectErrorCode(err))
}

func TestExtractTencentMySQL_RegionMismatch(t *testing.T) {
	now := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	api := &stubMySQLAPI{
		err: fmt.Errorf("tencent api error code=InvalidParameterValue.Region request_id=r"),
	}
	deps := tencentCollectDeps{
		now:      func() time.Time { return now },
		loadCur:  func(uint) (*tencentCollectCursor, error) { return nil, nil },
		saveCur:  func(uint, time.Time, time.Time) error { return errors.New("no") },
		enqueue:  func([]*SQLV2) error { return nil },
		mysqlAPI: api,
		sleep:    func(time.Duration) {},
	}
	ap := &AuditPlan{ID: 803, Type: TypeTencentCdbMySQLSlowLog, Params: testTencentAPParams()}
	_, err := extractTencentMySQLSlowLog(logrus.NewEntry(logrus.New()), ap, nil, deps)
	assert.Error(t, err)
	assert.Equal(t, TencentErrCodeRegionMismatch, tencentCollectErrorCode(err))
}

func TestSanitizeTencentErrorText(t *testing.T) {
	in := "boom SecretKey=abc123 Authorization: TC3-HMAC-SHA256 Credential=x Signature=deadbeef"
	out := sanitizeTencentErrorText(in)
	assert.NotContains(t, out, "abc123")
	assert.NotContains(t, out, "deadbeef")
	assert.Contains(t, out, "***")
}

// AC-020：超长 + 空 SQL 跳过，后续正常条目仍入队；跳过不进 dedup pending。
func TestFilterAndAggregateTencentItems_IsolationOversizedAndEmpty(t *testing.T) {
	var buf bytes.Buffer
	log := logrus.New()
	log.SetOutput(&buf)
	log.SetLevel(logrus.InfoLevel)
	logger := logrus.NewEntry(log)

	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	longSQL := strings.Repeat("x", tencentSQLMaxLen+10)
	items := []tencentSlowLogItem{
		{SQLText: "", ExecutionStartTime: now, Md5: "e1"},
		{SQLText: longSQL, ExecutionStartTime: now.Add(time.Second), Md5: "e2"},
		{SQLText: "SELECT id FROM t_a WHERE id=1", ExecutionStartTime: now.Add(2 * time.Second), Md5: "e3", UserName: "u", SchemaName: "db"},
		{SQLText: "SELECT name FROM t_b WHERE name='x'", ExecutionStartTime: now.Add(3 * time.Second), Md5: "e4", UserName: "u", SchemaName: "db"},
	}

	deps := tencentCollectDeps{
		dedupSeen:   func(uint, string) (bool, error) { return false, nil },
		dedupCommit: func(uint, []string) error { return nil },
	}
	deps.withDefaults(nil)

	ap := &AuditPlan{ID: 900, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 1, ProjectId: "p1"}
	cache := NewSQLV2Cache()
	pending, err := filterAndAggregateTencentItems(logger, ap, "cdb-1", items, true, deps, cache, &DefaultTaskV2{})
	assert.NoError(t, err)
	assert.Len(t, pending, 2, "empty/oversized must not enter pending dedup; two valid events remain")
	assert.GreaterOrEqual(t, len(cache.GetSQLs()), 1, "subsequent valid items still aggregated")

	logOut := buf.String()
	assert.Contains(t, logOut, "skip empty")
	assert.Contains(t, logOut, "skip oversized")
	assert.Contains(t, logOut, "empty=1")
	assert.Contains(t, logOut, "oversized=1")
	assert.NotContains(t, logOut, longSQL[:64], "must not log oversized SQL body")
}

func TestFilterAndAggregateTencentItems_SkipDoesNotBlockBatch(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	now := time.Now().UTC()
	items := []tencentSlowLogItem{
		{SQLText: "   ", ExecutionStartTime: now, Md5: "bad"},
		{SQLText: "SELECT ok FROM t", ExecutionStartTime: now.Add(time.Second), Md5: "ok1"},
	}
	deps := tencentCollectDeps{}
	deps.withDefaults(nil)
	ap := &AuditPlan{ID: 901, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 2}
	cache := NewSQLV2Cache()
	pending, err := filterAndAggregateTencentItems(logger, ap, "cdb-x", items, true, deps, cache, &DefaultTaskV2{})
	assert.NoError(t, err)
	assert.Len(t, pending, 1)
	assert.Equal(t, "SELECT ok FROM t", cache.GetSQLs()[0].SQLContent)
}

func TestSafeBuildTencentSQLV2_FingerprintFallback(t *testing.T) {
	ap := &AuditPlan{ID: 1, Type: TypeTencentCdbMySQLSlowLog, InstanceAuditPlanId: 1}
	sql, ok := safeBuildTencentSQLV2(logrus.NewEntry(logrus.New()), ap, tencentSlowLogItem{
		SQLText: "SELECT /* odd */ 1", QueryTimeSec: tencentFloat64(0.1),
	}, true)
	assert.True(t, ok)
	assert.NotNil(t, sql)
	assert.NotEmpty(t, sql.Fingerprint)
}
