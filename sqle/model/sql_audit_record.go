package model

import (
	"encoding/json"
	"fmt"
)

type SQLAuditRecordTags struct {
	Tags []string `json:"tags"`
}

type SQLAuditRecord struct {
	Model
	ProjectId     uint   `gorm:"index;not null"`
	CreatorId     uint   `gorm:"not null"`
	AuditRecordId string `gorm:"unique;not null"`
	Tags          JSON

	TaskId uint  `gorm:"not null"`
	Task   *Task `gorm:"foreignkey:TaskId"`
}

type SQLAuditRecordUpdateData struct {
	Tags []string
}

func (s *Storage) UpdateSQLAuditRecordById(SQLAuditRecordId string, data SQLAuditRecordUpdateData) error {
	tags, err := json.Marshal(data.Tags)
	if err != nil {
		return fmt.Errorf("marshal tags failed: %v", err)
	}

	db := s.db.Model(SQLAuditRecord{}).Where("audit_record_id = ?", SQLAuditRecordId)
	if len(data.Tags) == 0 {
		return db.Update("tags", nil).Error
	}
	return db.Update("tags", tags).Error
}

func (s *Storage) IsUserCanUpdateSQLAuditRecord(userId, projectId uint, SQLAuditRecordId string) (bool, error) {
	isManager, err := s.IsProjectManagerByID(userId, projectId)
	if err != nil {
		return false, fmt.Errorf("check project manager failed: %v", err)
	}
	if isManager {
		return true, nil
	}

	count := 0
	if err := s.db.Table("sql_audit_records").
		Where("audit_record_id = ?", SQLAuditRecordId).
		Where("creator_id = ?", userId).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check creator failed: %v", err)
	}
	return count == 1, nil
}
