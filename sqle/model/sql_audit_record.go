package model

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/actiontech/sqle/sqle/errors"
)

type SQLAuditRecordTags struct {
	Tags []string `json:"tags"`
}

type SQLAuditRecord struct {
	Model
	ProjectId     string `gorm:"index;not null"`
	CreatorId     string `gorm:"not null"`
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
		return errors.New(errors.DataInvalid, fmt.Errorf("marshal tags failed: %v", err))
	}

	db := s.db.Model(SQLAuditRecord{}).Where("audit_record_id = ?", SQLAuditRecordId)
	if len(data.Tags) == 0 {
		err = db.Update("tags", nil).Error
	} else {
		err = db.Update("tags", tags).Error
	}

	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsSQLAuditRecordBelongToCurrentUser(userId, projectId string, SQLAuditRecordId string) (bool, error) {
	count := 0
	if err := s.db.Table("sql_audit_records").
		Where("audit_record_id = ?", SQLAuditRecordId).
		Where("creator_id = ?", userId).
		Count(&count).Error; err != nil {
		return false, errors.New(errors.ConnectStorageError, fmt.Errorf("check creator failed: %v", err))
	}
	return count == 1, nil
}

func (s *Storage) GetSQLAuditRecordById(projectId string, SQLAuditRecordId string) (record *SQLAuditRecord, exist bool, err error) {
	record = &SQLAuditRecord{}
	if err = s.db.Preload("Task").Preload("Task.ExecuteSQLs").
		Where("project_id = ?", projectId).Where("audit_record_id = ?", SQLAuditRecordId).
		Find(&record).Error; err != nil && err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return record, true, nil
}
