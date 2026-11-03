package auditplan

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/model"
)

const (
	tencentCollectCursorSchemaVersion = 1
	tencentCollectOverlapDefault      = 5 * time.Minute
	tencentMySQLRetentionHours        = 30 * 24
	tencentPGRetentionHours           = 7 * 24
)

// tencentCollectCursor 持久化游标（权威在 DB，禁止仅内存 lastEndTime）。
type tencentCollectCursor struct {
	SchemaVersion        int    `json:"schema_version"`
	LastSuccessWindowEnd string `json:"last_success_window_end"` // RFC3339 UTC，半开区间右端
	LastSuccessAt        string `json:"last_success_at"`         // 游标写入墙钟 RFC3339 UTC
}

func parseTencentCollectCursor(raw model.JSON) (*tencentCollectCursor, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var c tencentCollectCursor
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("parse collect_cursor: %w", err)
	}
	if c.LastSuccessWindowEnd == "" {
		return nil, nil
	}
	return &c, nil
}

func marshalTencentCollectCursor(c *tencentCollectCursor) (model.JSON, error) {
	if c == nil {
		return nil, nil
	}
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return model.JSON(b), nil
}

func loadTencentCollectCursor(persist *model.Storage, auditPlanID uint) (*tencentCollectCursor, error) {
	info, exist, err := persist.GetAuditPlanTaskInfoByAPID(auditPlanID)
	if err != nil {
		return nil, err
	}
	if !exist || info == nil {
		return nil, nil
	}
	return parseTencentCollectCursor(info.CollectCursor)
}

func saveTencentCollectCursor(persist *model.Storage, auditPlanID uint, windowEnd, now time.Time) error {
	c := &tencentCollectCursor{
		SchemaVersion:        tencentCollectCursorSchemaVersion,
		LastSuccessWindowEnd: windowEnd.UTC().Format(time.RFC3339),
		LastSuccessAt:        now.UTC().Format(time.RFC3339),
	}
	raw, err := marshalTencentCollectCursor(c)
	if err != nil {
		return err
	}
	return persist.UpdateAuditPlanInfoByAPID(auditPlanID, map[string]interface{}{
		"collect_cursor": raw,
	})
}

// calcTencentCollectWindow 计算本轮 [start,end]。
// 首次：now - firstHours（再按保留期裁剪）；增量：last_success_window_end - overlap。
func calcTencentCollectWindow(now time.Time, cursor *tencentCollectCursor, firstHours, retentionHours int, overlap time.Duration) (start, end time.Time, err error) {
	if firstHours <= 0 {
		firstHours = 24
	}
	if retentionHours <= 0 {
		retentionHours = firstHours
	}
	if overlap <= 0 {
		overlap = tencentCollectOverlapDefault
	}
	end = now.UTC()
	if cursor == nil || cursor.LastSuccessWindowEnd == "" {
		start = end.Add(-time.Duration(firstHours) * time.Hour)
	} else {
		lastEnd, parseErr := time.Parse(time.RFC3339, cursor.LastSuccessWindowEnd)
		if parseErr != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid last_success_window_end: %w", parseErr)
		}
		start = lastEnd.UTC().Add(-overlap)
	}
	earliest := end.Add(-time.Duration(retentionHours) * time.Hour)
	if start.Before(earliest) {
		start = earliest
	}
	if !start.Before(end) {
		// 空/退化窗：仍返回合法半开语义，调用方按空窗推进
		start = end.Add(-time.Second)
	}
	return start, end, nil
}
