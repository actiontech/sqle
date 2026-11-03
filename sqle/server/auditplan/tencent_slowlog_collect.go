package auditplan

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/sirupsen/logrus"
)

// tencentEventDedupSeen / Commit：可注入；默认走 audit_plan_collect_event_dedup。
type tencentEventDedupSeen func(auditPlanID uint, eventKeyHash string) (bool, error)
type tencentEventDedupCommit func(auditPlanID uint, eventKeyHashes []string) error

const (
	tencentMySQLPreferredPageLimit uint64 = 100 // 硬上限 800，本轮 preferred=100
	tencentPGPreferredPageLimit    uint64 = 50  // 硬上限 100

	tencentMySQLMaxOffset    uint64 = 9999  // S2：Offset≤9999
	tencentPGMaxTotalCount   int64  = 10000 // S2：TotalCount≤10000
	tencentMinWindowDuration        = time.Second
)

// errTencentMinWindowExceeded 最小时间窗仍超上限；调用方不得推进游标、不得忙循环。
var errTencentMinWindowExceeded = errors.New("tencent slowlog: min window still over limit; cursor preserved")

// errTencentCollectStopped 任务已禁用/收到 stop：停止新的腾讯云 API，本批失败不推进游标（AC-015）。
var errTencentCollectStopped = errors.New("tencent slowlog: collect stopped (disabled); cursor preserved")

// tencentSlowLogEnqueuer 方案 A：入队成功边界（可注入以便单测）。
type tencentSlowLogEnqueuer func(sqls []*SQLV2) error

func defaultTencentEnqueuer(persist *model.Storage) tencentSlowLogEnqueuer {
	return func(sqls []*SQLV2) error {
		if len(sqls) == 0 {
			return nil
		}
		queues := make([]*model.SQLManageQueue, 0, len(sqls))
		for _, sql := range sqls {
			queues = append(queues, ConvertSQLV2ToMangerSQLQueue(sql))
		}
		return persist.PushSQLToManagerSQLQueue(queues)
	}
}

type tencentCollectDeps struct {
	now      func() time.Time
	loadCur  func(auditPlanID uint) (*tencentCollectCursor, error)
	saveCur  func(auditPlanID uint, windowEnd, now time.Time) error
	enqueue  tencentSlowLogEnqueuer
	mysqlAPI tencentMySQLSlowLogAPI
	pgAPI    tencentPGSlowLogAPI
	overlap  time.Duration
	// pageLimit 覆盖 preferred（单测注入小页；0 用默认）。仍受引擎硬上限约束。
	mysqlPageLimit uint64
	pgPageLimit    uint64
	// 事件幂等（AC-011）
	dedupSeen   tencentEventDedupSeen
	dedupCommit tencentEventDedupCommit
	// API 退避（AC-014）
	sleep          func(time.Duration)
	maxAPIAttempts int
	// 禁用停请求（AC-015）：为 true 时停止新 API / 不推进游标。默认查 ActiveStatus。
	shouldStop func() bool
}

func (d *tencentCollectDeps) withDefaults(persist *model.Storage) {
	if d.now == nil {
		d.now = time.Now
	}
	if d.loadCur == nil {
		d.loadCur = func(id uint) (*tencentCollectCursor, error) {
			return loadTencentCollectCursor(persist, id)
		}
	}
	if d.saveCur == nil {
		d.saveCur = func(id uint, windowEnd, now time.Time) error {
			return saveTencentCollectCursor(persist, id, windowEnd, now)
		}
	}
	if d.enqueue == nil {
		d.enqueue = defaultTencentEnqueuer(persist)
	}
	if d.dedupSeen == nil && persist != nil {
		d.dedupSeen = func(id uint, hash string) (bool, error) {
			return persist.HasAuditPlanCollectEventDedup(id, hash)
		}
	}
	if d.dedupSeen == nil {
		d.dedupSeen = func(uint, string) (bool, error) { return false, nil }
	}
	if d.dedupCommit == nil && persist != nil {
		d.dedupCommit = func(id uint, hashes []string) error {
			return persist.SaveAuditPlanCollectEventDedups(id, hashes)
		}
	}
	if d.dedupCommit == nil {
		d.dedupCommit = func(uint, []string) error { return nil }
	}
	if d.sleep == nil {
		d.sleep = time.Sleep
	}
	if d.maxAPIAttempts <= 0 {
		d.maxAPIAttempts = tencentAPIMaxAttempts
	}
	if d.overlap <= 0 {
		d.overlap = tencentCollectOverlapDefault
	}
	if d.mysqlPageLimit == 0 {
		d.mysqlPageLimit = tencentMySQLPreferredPageLimit
	}
	if d.mysqlPageLimit > 800 {
		d.mysqlPageLimit = 800
	}
	if d.pgPageLimit == 0 {
		d.pgPageLimit = tencentPGPreferredPageLimit
	}
	if d.pgPageLimit > 100 {
		d.pgPageLimit = 100
	}
}

// tencentPageExhausted 翻页终止：本页不足 Limit，或已覆盖 TotalCount。
func tencentPageExhausted(pageItemCount int, pageLimit, offset uint64, totalCount int64) bool {
	if pageItemCount <= 0 {
		return true
	}
	if uint64(pageItemCount) < pageLimit {
		return true
	}
	if totalCount > 0 && offset+uint64(pageItemCount) >= uint64(totalCount) {
		return true
	}
	return false
}

func tencentWindowAtMin(start, end time.Time) bool {
	return !end.After(start.Add(tencentMinWindowDuration))
}

// defaultTencentShouldStop：非 active（禁用/删除）→ 停止新请求。DB 查询错误不误停。
func defaultTencentShouldStop(persist *model.Storage, auditPlanID uint) func() bool {
	return func() bool {
		if persist == nil {
			return false
		}
		_, exist, err := persist.GetActiveAuditPlanDetail(auditPlanID)
		if err != nil {
			return false
		}
		return !exist
	}
}

func checkTencentCollectStopped(deps tencentCollectDeps) error {
	if deps.shouldStop != nil && deps.shouldStop() {
		return errTencentCollectStopped
	}
	return nil
}

func splitTencentTimeWindow(start, end time.Time) (mid time.Time, ok bool) {
	if tencentWindowAtMin(start, end) {
		return time.Time{}, false
	}
	mid = start.Add(end.Sub(start) / 2)
	if !mid.After(start) {
		mid = start.Add(tencentMinWindowDuration)
	}
	if !end.After(mid) {
		return time.Time{}, false
	}
	return mid, true
}

// isTencentResultOverLimitErr 结果集超限类错误 → 缩窗，禁止盲重试同窗。
func isTencentResultOverLimitErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, code := range []string{"ResultSetOverLimit", "LogContentOverLimit", "ResultOverLimit"} {
		if strings.Contains(msg, code) {
			return true
		}
	}
	return false
}

// tencentMySQLNeedsShrinkByTotal：末页所需 Offset 将越过 9999。
func tencentMySQLNeedsShrinkByTotal(pageLimit uint64, totalCount int64) bool {
	if totalCount <= 0 || pageLimit == 0 {
		return false
	}
	lastOffset := ((uint64(totalCount) - 1) / pageLimit) * pageLimit
	return lastOffset > tencentMySQLMaxOffset
}

func tencentMySQLNextOffsetOverMax(offset, pageLimit uint64) bool {
	return offset+pageLimit > tencentMySQLMaxOffset
}

func tencentPGNeedsShrinkByTotal(totalCount int64) bool {
	return totalCount >= tencentPGMaxTotalCount
}

func readTencentCommonParams(ap *AuditPlan) (secretId, secretKey, region, instanceId string, firstHours int, err error) {
	secretId = ap.Params.GetParam(paramKeyAccessKeyId).String()
	if secretId == "" {
		return "", "", "", "", 0, fmt.Errorf("access_key_id is not configured")
	}
	secretKey = ap.Params.GetParam(paramKeyAccessKeySecret).String()
	if secretKey == "" {
		return "", "", "", "", 0, fmt.Errorf("access_key_secret is not configured")
	}
	region = ap.Params.GetParam(paramKeyRegion).String()
	if region == "" {
		return "", "", "", "", 0, fmt.Errorf("region is not configured")
	}
	instanceId = ap.Params.GetParam(paramKeyDBInstanceId).String()
	if instanceId == "" {
		return "", "", "", "", 0, fmt.Errorf("db_instance_id is not configured")
	}
	firstHours = ap.Params.GetParam(paramKeyFirstSqlsScrappedInLastPeriodHours).Int()
	if firstHours == 0 {
		firstHours = 24
	}
	return secretId, secretKey, region, instanceId, firstHours, nil
}

// extractTencentMySQLSlowLog 方案 A：拉取（含缩窗）→入队→写游标→返回空切片。
func extractTencentMySQLSlowLog(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage, deps tencentCollectDeps) (sqls []*SQLV2, err error) {
	defer func() { err = wrapTencentCollectError(err) }()
	deps.withDefaults(persist)
	if deps.shouldStop == nil {
		deps.shouldStop = defaultTencentShouldStop(persist, ap.ID)
	}
	secretId, secretKey, region, instanceId, firstHours, err := readTencentCommonParams(ap)
	if err != nil {
		return nil, err
	}
	if deps.mysqlAPI == nil {
		deps.mysqlAPI = newTencentMySQLOpenAPIClient(secretId, secretKey, region)
	}
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, err
	}

	cursor, err := deps.loadCur(ap.ID)
	if err != nil {
		return nil, fmt.Errorf("load collect cursor: %w", err)
	}
	now := deps.now()
	start, end, err := calcTencentCollectWindow(now, cursor, firstHours, tencentMySQLRetentionHours, deps.overlap)
	if err != nil {
		return nil, err
	}
	logger.Infof("tencent mysql slowlog window start=%s end=%s first=%v", start.Format(time.RFC3339), end.Format(time.RFC3339), cursor == nil)

	allItems, pages, err := collectTencentMySQLWindow(logger, deps, deps.mysqlAPI, instanceId, start, end, deps.mysqlPageLimit)
	if err != nil {
		return nil, err
	}
	// 禁用后丢弃已拉到的 in-flight 结果，不入队、不推进游标、不开收尾 API
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, err
	}

	cache := NewSQLV2Cache()
	aggregator := &DefaultTaskV2{}
	pendingKeys, err := filterAndAggregateTencentItems(logger, ap, instanceId, allItems, true, deps, cache, aggregator)
	if err != nil {
		return nil, err
	}
	enqueued := cache.GetSQLs()
	if err := deps.enqueue(enqueued); err != nil {
		return nil, fmt.Errorf("enqueue sql failed: %w", err)
	}
	if err := deps.dedupCommit(ap.ID, pendingKeys); err != nil {
		return nil, fmt.Errorf("commit event dedup failed: %w", err)
	}
	// AC-016：全部子窗分页+入队成功后才推进到本轮声明 end
	if err := deps.saveCur(ap.ID, end, now); err != nil {
		return nil, fmt.Errorf("save collect cursor failed: %w", err)
	}
	logger.Infof("tencent mysql slowlog collected items=%d pages=%d enqueued=%d dedup_new=%d cursor_end=%s",
		len(allItems), pages, len(enqueued), len(pendingKeys), end.Format(time.RFC3339))
	return []*SQLV2{}, nil
}

// collectTencentMySQLWindow 对 [start,end] 完整分页；触顶则二分缩窗。返回本窗（含子窗）条目与 API 调用次数。
func collectTencentMySQLWindow(logger *logrus.Entry, deps tencentCollectDeps, api tencentMySQLSlowLogAPI, instanceId string, start, end time.Time, pageLimit uint64) ([]tencentSlowLogItem, int, error) {
	var allItems []tencentSlowLogItem
	var offset uint64
	var pages int
	for {
		if err := checkTencentCollectStopped(deps); err != nil {
			return nil, pages, err
		}
		var resp *tencentMySQLSlowLogResp
		err := withTencentAPIRetry(deps.sleep, deps.maxAPIAttempts, deps.shouldStop, func() error {
			var callErr error
			resp, callErr = api.DescribeSlowLogData(tencentMySQLSlowLogReq{
				InstanceId: instanceId,
				StartTime:  start.Unix(),
				EndTime:    end.Unix(),
				Offset:     offset,
				Limit:      pageLimit,
			})
			return callErr
		})
		pages++
		if err != nil {
			if errors.Is(err, errTencentCollectStopped) {
				return nil, pages, err
			}
			if isTencentResultOverLimitErr(err) {
				logger.Infof("tencent mysql shrink: result over limit offset=%d window=[%s,%s]", offset, start.Format(time.RFC3339), end.Format(time.RFC3339))
				return shrinkCollectTencentMySQL(logger, deps, api, instanceId, start, end, pageLimit, pages)
			}
			return nil, pages, fmt.Errorf("DescribeSlowLogData failed (offset=%d): %w", offset, err)
		}
		if tencentMySQLNeedsShrinkByTotal(pageLimit, resp.TotalCount) {
			logger.Infof("tencent mysql shrink: total_count=%d exceeds offset max window=[%s,%s]", resp.TotalCount, start.Format(time.RFC3339), end.Format(time.RFC3339))
			return shrinkCollectTencentMySQL(logger, deps, api, instanceId, start, end, pageLimit, pages)
		}
		allItems = append(allItems, resp.Items...)
		if tencentPageExhausted(len(resp.Items), pageLimit, offset, resp.TotalCount) {
			return allItems, pages, nil
		}
		if tencentMySQLNextOffsetOverMax(offset, pageLimit) {
			logger.Infof("tencent mysql shrink: next offset would exceed %d (cur=%d limit=%d) window=[%s,%s]",
				tencentMySQLMaxOffset, offset, pageLimit, start.Format(time.RFC3339), end.Format(time.RFC3339))
			return shrinkCollectTencentMySQL(logger, deps, api, instanceId, start, end, pageLimit, pages)
		}
		offset += pageLimit
	}
}

func shrinkCollectTencentMySQL(logger *logrus.Entry, deps tencentCollectDeps, api tencentMySQLSlowLogAPI, instanceId string, start, end time.Time, pageLimit uint64, pagesSoFar int) ([]tencentSlowLogItem, int, error) {
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, pagesSoFar, err
	}
	mid, ok := splitTencentTimeWindow(start, end)
	if !ok {
		return nil, pagesSoFar, fmt.Errorf("DescribeSlowLogData over limit at min window [%s,%s]: %w",
			start.Format(time.RFC3339), end.Format(time.RFC3339), errTencentMinWindowExceeded)
	}
	logger.Infof("tencent mysql shrink split [%s,%s) + [%s,%s]", start.Format(time.RFC3339), mid.Format(time.RFC3339), mid.Format(time.RFC3339), end.Format(time.RFC3339))
	left, p1, err := collectTencentMySQLWindow(logger, deps, api, instanceId, start, mid, pageLimit)
	pagesSoFar += p1
	if err != nil {
		return nil, pagesSoFar, err
	}
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, pagesSoFar, err
	}
	right, p2, err := collectTencentMySQLWindow(logger, deps, api, instanceId, mid, end, pageLimit)
	pagesSoFar += p2
	if err != nil {
		return nil, pagesSoFar, err
	}
	return append(left, right...), pagesSoFar, nil
}

// extractTencentPGSlowLog 方案 A：同上。
func extractTencentPGSlowLog(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage, deps tencentCollectDeps) (sqls []*SQLV2, err error) {
	defer func() { err = wrapTencentCollectError(err) }()
	deps.withDefaults(persist)
	if deps.shouldStop == nil {
		deps.shouldStop = defaultTencentShouldStop(persist, ap.ID)
	}
	secretId, secretKey, region, instanceId, firstHours, err := readTencentCommonParams(ap)
	if err != nil {
		return nil, err
	}
	dbName := ""
	if p := ap.Params.GetParam(paramKeyDatabaseName); p != nil {
		dbName = p.String()
	}
	if deps.pgAPI == nil {
		deps.pgAPI = newTencentPGOpenAPIClient(secretId, secretKey, region)
	}
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, err
	}

	cursor, err := deps.loadCur(ap.ID)
	if err != nil {
		return nil, fmt.Errorf("load collect cursor: %w", err)
	}
	now := deps.now()
	start, end, err := calcTencentCollectWindow(now, cursor, firstHours, tencentPGRetentionHours, deps.overlap)
	if err != nil {
		return nil, err
	}
	logger.Infof("tencent pg slowlog window start=%s end=%s first=%v", start.Format(time.RFC3339), end.Format(time.RFC3339), cursor == nil)

	allItems, pages, err := collectTencentPGWindow(logger, deps, deps.pgAPI, instanceId, dbName, start, end, deps.pgPageLimit)
	if err != nil {
		return nil, err
	}
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, err
	}

	cache := NewSQLV2Cache()
	aggregator := &DefaultTaskV2{}
	pendingKeys, err := filterAndAggregateTencentItems(logger, ap, instanceId, allItems, false, deps, cache, aggregator)
	if err != nil {
		return nil, err
	}
	enqueued := cache.GetSQLs()
	if err := deps.enqueue(enqueued); err != nil {
		return nil, fmt.Errorf("enqueue sql failed: %w", err)
	}
	if err := deps.dedupCommit(ap.ID, pendingKeys); err != nil {
		return nil, fmt.Errorf("commit event dedup failed: %w", err)
	}
	if err := deps.saveCur(ap.ID, end, now); err != nil {
		return nil, fmt.Errorf("save collect cursor failed: %w", err)
	}
	logger.Infof("tencent pg slowlog collected items=%d pages=%d enqueued=%d dedup_new=%d cursor_end=%s",
		len(allItems), pages, len(enqueued), len(pendingKeys), end.Format(time.RFC3339))
	return []*SQLV2{}, nil
}

func collectTencentPGWindow(logger *logrus.Entry, deps tencentCollectDeps, api tencentPGSlowLogAPI, instanceId, dbName string, start, end time.Time, pageLimit uint64) ([]tencentSlowLogItem, int, error) {
	var allItems []tencentSlowLogItem
	var offset uint64
	var pages int
	for {
		if err := checkTencentCollectStopped(deps); err != nil {
			return nil, pages, err
		}
		var resp *tencentPGSlowLogResp
		err := withTencentAPIRetry(deps.sleep, deps.maxAPIAttempts, deps.shouldStop, func() error {
			var callErr error
			resp, callErr = api.DescribeSlowQueryList(tencentPGSlowLogReq{
				DBInstanceId: instanceId,
				StartTime:    formatTencentPGTime(start),
				EndTime:      formatTencentPGTime(end),
				Offset:       offset,
				Limit:        pageLimit,
				DatabaseName: dbName,
			})
			return callErr
		})
		pages++
		if err != nil {
			if errors.Is(err, errTencentCollectStopped) {
				return nil, pages, err
			}
			if isTencentResultOverLimitErr(err) {
				logger.Infof("tencent pg shrink: result over limit offset=%d window=[%s,%s]", offset, start.Format(time.RFC3339), end.Format(time.RFC3339))
				return shrinkCollectTencentPG(logger, deps, api, instanceId, dbName, start, end, pageLimit, pages)
			}
			return nil, pages, fmt.Errorf("DescribeSlowQueryList failed (offset=%d): %w", offset, err)
		}
		if tencentPGNeedsShrinkByTotal(resp.TotalCount) {
			logger.Infof("tencent pg shrink: total_count=%d >= %d window=[%s,%s]", resp.TotalCount, tencentPGMaxTotalCount, start.Format(time.RFC3339), end.Format(time.RFC3339))
			return shrinkCollectTencentPG(logger, deps, api, instanceId, dbName, start, end, pageLimit, pages)
		}
		allItems = append(allItems, resp.Items...)
		if tencentPageExhausted(len(resp.Items), pageLimit, offset, resp.TotalCount) {
			return allItems, pages, nil
		}
		offset += pageLimit
	}
}

func shrinkCollectTencentPG(logger *logrus.Entry, deps tencentCollectDeps, api tencentPGSlowLogAPI, instanceId, dbName string, start, end time.Time, pageLimit uint64, pagesSoFar int) ([]tencentSlowLogItem, int, error) {
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, pagesSoFar, err
	}
	mid, ok := splitTencentTimeWindow(start, end)
	if !ok {
		return nil, pagesSoFar, fmt.Errorf("DescribeSlowQueryList over limit at min window [%s,%s]: %w",
			start.Format(time.RFC3339), end.Format(time.RFC3339), errTencentMinWindowExceeded)
	}
	logger.Infof("tencent pg shrink split [%s,%s) + [%s,%s]", start.Format(time.RFC3339), mid.Format(time.RFC3339), mid.Format(time.RFC3339), end.Format(time.RFC3339))
	left, p1, err := collectTencentPGWindow(logger, deps, api, instanceId, dbName, start, mid, pageLimit)
	pagesSoFar += p1
	if err != nil {
		return nil, pagesSoFar, err
	}
	if err := checkTencentCollectStopped(deps); err != nil {
		return nil, pagesSoFar, err
	}
	right, p2, err := collectTencentPGWindow(logger, deps, api, instanceId, dbName, mid, end, pageLimit)
	pagesSoFar += p2
	if err != nil {
		return nil, pagesSoFar, err
	}
	return append(left, right...), pagesSoFar, nil
}

// filterAndAggregateTencentItems 事件幂等过滤后聚合；返回待提交的新 event_key_hash。
func filterAndAggregateTencentItems(
	logger *logrus.Entry,
	ap *AuditPlan,
	instanceId string,
	items []tencentSlowLogItem,
	mysqlFingerprint bool,
	deps tencentCollectDeps,
	cache SQLV2Cacher,
	aggregator *DefaultTaskV2,
) ([]string, error) {
	pending := make([]string, 0, len(items))
	seenBatch := make(map[string]struct{}, len(items))
	skippedDedup := 0
	skippedEmpty := 0
	skippedOversized := 0
	skippedAggregate := 0
	for _, item := range items {
		// AC-020：空文本 / 超长 → 跳过该条，不阻断整批，不写成功 dedup。
		sqlLen := len(item.SQLText)
		if strings.TrimSpace(item.SQLText) == "" {
			skippedEmpty++
			if logger != nil {
				logger.Warnf("tencent slowlog skip empty sql audit_plan_id=%d event_key_prefix=%s sql_len=%d",
					ap.ID, tencentItemEventKeyPrefix(instanceId, item, mysqlFingerprint, logger), sqlLen)
			}
			continue
		}
		if sqlLen > tencentSQLMaxLen {
			skippedOversized++
			if logger != nil {
				logger.Warnf("tencent slowlog skip oversized sql audit_plan_id=%d event_key_prefix=%s sql_len=%d max=%d",
					ap.ID, tencentItemEventKeyPrefix(instanceId, item, mysqlFingerprint, logger), sqlLen, tencentSQLMaxLen)
			}
			continue
		}

		var key string
		if mysqlFingerprint {
			key = tencentMySQLEventKeyHash(instanceId, item)
		} else {
			key = tencentPGEventKeyHash(logger, instanceId, item)
		}
		if _, ok := seenBatch[key]; ok {
			skippedDedup++
			continue
		}
		exists, err := deps.dedupSeen(ap.ID, key)
		if err != nil {
			return nil, fmt.Errorf("check event dedup: %w", err)
		}
		if exists {
			skippedDedup++
			continue
		}
		seenBatch[key] = struct{}{}

		sqlV2, buildOK := safeBuildTencentSQLV2(logger, ap, item, mysqlFingerprint)
		if !buildOK || sqlV2 == nil {
			skippedAggregate++
			continue
		}
		if err := aggregator.AggregateSQL(cache, sqlV2); err != nil {
			skippedAggregate++
			if logger != nil {
				logger.Warnf("tencent slowlog skip aggregate failed audit_plan_id=%d event_key_prefix=%s err=%v",
					ap.ID, keyPrefix(key), err)
			}
			continue
		}
		pending = append(pending, key)
	}
	if logger != nil && (skippedDedup > 0 || skippedEmpty > 0 || skippedOversized > 0 || skippedAggregate > 0) {
		logger.Infof("tencent slowlog item isolation dedup_skipped=%d empty=%d oversized=%d aggregate_or_fp=%d new=%d",
			skippedDedup, skippedEmpty, skippedOversized, skippedAggregate, len(pending))
	}
	return pending, nil
}

func keyPrefix(key string) string {
	if len(key) <= 12 {
		return key
	}
	return key[:12]
}

func tencentItemEventKeyPrefix(instanceId string, item tencentSlowLogItem, mysqlFingerprint bool, logger *logrus.Entry) string {
	var key string
	if mysqlFingerprint {
		key = tencentMySQLEventKeyHash(instanceId, item)
	} else {
		key = tencentPGEventKeyHash(logger, instanceId, item)
	}
	return keyPrefix(key)
}

// safeBuildTencentSQLV2 捕获指纹/构建 panic，失败时跳过该条（AC-020）。
func safeBuildTencentSQLV2(logger *logrus.Entry, ap *AuditPlan, item tencentSlowLogItem, mysqlFingerprint bool) (sql *SQLV2, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			sql = nil
			if logger != nil {
				logger.Warnf("tencent slowlog skip fingerprint/build panic audit_plan_id=%d recover=%v sql_len=%d",
					ap.ID, r, len(item.SQLText))
			}
		}
	}()
	return buildTencentSQLV2(ap, item, mysqlFingerprint), true
}

// buildTencentSQLV2 按 S3 AC-017/018 标准化单条：Source=ap.Type；MySQL 耗时按秒；
// PG Duration 未校准前只写诊断键，禁止定死写入 query_time_*；缺字段不伪造 0。
func buildTencentSQLV2(ap *AuditPlan, item tencentSlowLogItem, mysqlFingerprint bool) *SQLV2 {
	fp := item.SQLText
	if mysqlFingerprint {
		if f, err := util.Fingerprint(item.SQLText, true); err == nil && f != "" {
			fp = f
		}
	}
	info := NewMetrics()
	info.SetInt(MetricNameCounter, 1)
	info.SetString(MetricNameLastReceiveTimestamp, time.Now().UTC().Format(time.RFC3339))
	if !item.ExecutionStartTime.IsZero() {
		info.SetString(MetricNameStartTimeOfLastScrapedSQL, item.ExecutionStartTime.UTC().Format(time.RFC3339))
	}
	if item.UserName != "" {
		info.SetString(MetricNameDBUser, item.UserName)
	}
	if item.UserHost != "" {
		info.SetString(MetricNameHost, item.UserHost)
	}

	if mysqlFingerprint {
		// MySQL：官方 QueryTime/LockTime 为秒；禁止再 ÷1000 / ×1000。
		if item.QueryTimeSec != nil {
			info.SetFloat(MetricNameQueryTimeAvg, *item.QueryTimeSec)
			info.SetFloat(MetricNameQueryTimeMax, *item.QueryTimeSec)
		}
		if item.LockTimeSec != nil {
			info.SetFloat(MetricNameLockTimeAvg, *item.LockTimeSec)
			info.SetFloat(MetricNameLockTimeMax, *item.LockTimeSec)
		}
		if item.RowsExamined != nil {
			info.SetFloat(MetricNameRowExaminedAvg, float64(*item.RowsExamined))
		}
		if item.Md5 != "" {
			info.SetString(MetricNameChecksum, item.Md5)
		}
	} else {
		// PG：S2 校准结论未冻结 → 只保留 raw + unit=pending，不得写入 query_time_*。
		if item.DurationRaw != "" {
			info.SetString(metricTencentPGDurationRaw, item.DurationRaw)
			info.SetString(metricTencentPGDurationUnit, tencentPGDurationUnitPending)
		}
	}

	sqlV2 := &SQLV2{
		Source:      ap.Type, // AC-017：硬等于扫描类型常量，禁止改写成其它云 type
		SourceId:    strconv.FormatUint(uint64(ap.InstanceAuditPlanId), 10),
		AuditPlanId: strconv.FormatUint(uint64(ap.ID), 10),
		ProjectId:   ap.ProjectId,
		InstanceID:  ap.InstanceID,
		SchemaName:  item.SchemaName,
		SQLContent:  item.SQLText,
		Fingerprint: fp,
		Info:        info,
	}
	sqlV2.GenSQLId()
	return sqlV2
}
