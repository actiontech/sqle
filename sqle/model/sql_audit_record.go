package model

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/actiontech/sqle/sqle/errors"
)

type SQLAuditRecordTags struct {
	Tags []string `json:"tags"`
}

type SQLAuditRecord struct {
	Model
	ProjectId     string `gorm:"index;not null;type:varchar(255)"`
	CreatorId     string `gorm:"not null;type:varchar(255)"`
	AuditRecordId string `gorm:"not null;type:varchar(255);index:unique"`
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
	var count int64
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
		First(&record).Error; err != nil && err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return record, true, nil
}

func (s *Storage) GetSQLAuditRecordProjectIdByTaskId(taskId uint) (projectId string, err error) {
	var record = &SQLAuditRecord{}
	if err = s.db.Where("task_id = ?", taskId).Select("project_id").First(&record).Error; err != nil && err == gorm.ErrRecordNotFound {
		return "", nil
	} else if err != nil {
		return "", errors.New(errors.ConnectStorageError, err)
	}
	return record.ProjectId, nil
}
