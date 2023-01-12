package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/actiontech/sqle/sqle/errors"
)

const (
	SyncInstanceStatusSuccess = "success"
	SyncInstanceStatusFailed  = "failed"
)

type SyncInstanceTask struct {
	Model
	Source         string `json:"source" gorm:"not null"`
	Version        string `json:"version" gorm:"not null"`
	URL            string `json:"url" gorm:"not null"`
	DbType         string `json:"db_type" gorm:"not null"`
	RuleTemplateID uint   `json:"rule_template_id" gorm:"not null"`
	// Cron表达式
	SyncInstanceInterval string     `json:"sync_instance_interval" gorm:"not null"`
	LastSyncStatus       string     `json:"last_sync_status" gorm:"default:\"initialized\""`
	LastSyncSuccessTime  *time.Time `json:"last_sync_success_time"`

	// 关系表
	RuleTemplate *RuleTemplate `gorm:"foreignKey:RuleTemplateID"`
}

func (s *Storage) GetAllSyncTasks() ([]SyncInstanceTask, error) {
	var syncTasks []SyncInstanceTask
	if err := s.db.Model(&SyncInstanceTask{}).Preload("RuleTemplate").Find(&syncTasks).Error; err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	return syncTasks, nil
}

func (s *Storage) GetSyncTaskById(id uint) (*SyncInstanceTask, bool, error) {
	syncInstTask := new(SyncInstanceTask)
	err := s.db.Model(&SyncInstanceTask{}).Where("id = ?", id).First(&syncInstTask).Error
	if err == gorm.ErrRecordNotFound {
		return syncInstTask, false, errors.ConnectStorageErrWrapper(err)
	}

	return syncInstTask, true, nil
}
