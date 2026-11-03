package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm/clause"
)

// AuditPlanCollectEventDedup 云慢日志采集事件幂等表（S2 AC-011）。
// 唯一键：audit_plan_id + event_key_hash；无存量回填，由 AutoMigrate 建表。
type AuditPlanCollectEventDedup struct {
	Model
	AuditPlanID  uint   `json:"audit_plan_id" gorm:"not null;uniqueIndex:uk_ap_event_key;index:idx_ap_collect_dedup_ap"`
	EventKeyHash string `json:"event_key_hash" gorm:"type:char(64);not null;uniqueIndex:uk_ap_event_key"`
}

func (AuditPlanCollectEventDedup) TableName() string {
	return "audit_plan_collect_event_dedup"
}

// HasAuditPlanCollectEventDedup 判断某审计计划是否已处理过该事件键。
func (s *Storage) HasAuditPlanCollectEventDedup(auditPlanID uint, eventKeyHash string) (bool, error) {
	var count int64
	err := s.db.Model(&AuditPlanCollectEventDedup{}).
		Where("audit_plan_id = ? AND event_key_hash = ?", auditPlanID, eventKeyHash).
		Limit(1).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

// SaveAuditPlanCollectEventDedups 批量写入幂等键；冲突忽略（幂等）。
func (s *Storage) SaveAuditPlanCollectEventDedups(auditPlanID uint, eventKeyHashes []string) error {
	if len(eventKeyHashes) == 0 {
		return nil
	}
	rows := make([]*AuditPlanCollectEventDedup, 0, len(eventKeyHashes))
	seen := make(map[string]struct{}, len(eventKeyHashes))
	for _, h := range eventKeyHashes {
		if h == "" {
			continue
		}
		if _, ok := seen[h]; ok {
			continue
		}
		seen[h] = struct{}{}
		rows = append(rows, &AuditPlanCollectEventDedup{
			AuditPlanID:  auditPlanID,
			EventKeyHash: h,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "audit_plan_id"}, {Name: "event_key_hash"}},
		DoNothing: true,
	}).CreateInBatches(rows, 200).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}
