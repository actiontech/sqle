package auditplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalcTencentCollectWindow_First24h(t *testing.T) {
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	start, end, err := calcTencentCollectWindow(now, nil, 24, tencentMySQLRetentionHours, tencentCollectOverlapDefault)
	assert.NoError(t, err)
	assert.Equal(t, now.UTC(), end)
	assert.Equal(t, now.Add(-24*time.Hour).UTC(), start)
}

func TestCalcTencentCollectWindow_IncrementalFromCursor(t *testing.T) {
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	lastEnd := now.Add(-1 * time.Hour)
	cursor := &tencentCollectCursor{
		SchemaVersion:        1,
		LastSuccessWindowEnd: lastEnd.Format(time.RFC3339),
	}
	start, end, err := calcTencentCollectWindow(now, cursor, 24, tencentMySQLRetentionHours, 5*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, now.UTC(), end)
	assert.Equal(t, lastEnd.Add(-5*time.Minute).UTC(), start)
	// 不得再按首次 24h 从零重拉
	assert.True(t, start.After(now.Add(-2*time.Hour)))
}

func TestCalcTencentCollectWindow_RetentionClip(t *testing.T) {
	now := time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)
	// firstHours 超过 PG 保留 7 天 → 裁剪
	start, end, err := calcTencentCollectWindow(now, nil, 24*30, tencentPGRetentionHours, tencentCollectOverlapDefault)
	assert.NoError(t, err)
	assert.Equal(t, now.UTC(), end)
	assert.Equal(t, now.Add(-time.Duration(tencentPGRetentionHours)*time.Hour).UTC(), start)
}
