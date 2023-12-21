//go:build enterprise
// +build enterprise

package model

// import (
// 	"time"

// 	"github.com/actiontech/sqle/sqle/errors"
// 	"github.com/jinzhu/gorm"
// )

// const (
// 	SyncTaskSourceActiontechDmp = "actiontech-dmp"

// 	SyncInstanceStatusSucceeded = "succeeded"
// 	SyncInstanceStatusFailed    = "failed"
// )

// func (s *Storage) GetAllSyncInstanceTasks() ([]SyncInstanceTask, error) {
// 	var syncTasks []SyncInstanceTask
// 	if err := s.db.Model(&SyncInstanceTask{}).Preload("RuleTemplate").Find(&syncTasks).Error; err != nil {
// 		return nil, errors.ConnectStorageErrWrapper(err)
// 	}

// 	return syncTasks, nil
// }

// func (s *Storage) GetSyncInstanceTaskById(id uint) (*SyncInstanceTask, bool, error) {
// 	syncInstTask := new(SyncInstanceTask)
// 	err := s.db.Model(&SyncInstanceTask{}).Preload("RuleTemplate").Where("id = ?", id).First(&syncInstTask).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return syncInstTask, false, errors.ConnectStorageErrWrapper(err)
// 	}

// 	return syncInstTask, true, nil
// }

// func (s *Storage) UpdateSyncInstanceTaskById(id uint, syncTask map[string]interface{}) error {
// 	if err := s.db.Model(&SyncInstanceTask{}).Where("id = ?", id).Updates(syncTask).Error; err != nil {
// 		return errors.ConnectStorageErrWrapper(err)
// 	}
// 	return nil
// }

// func (s *Storage) GetLatestSyncInstanceTask(after time.Time) (*SyncInstanceTask, bool, error) {
// 	task := new(SyncInstanceTask)
// 	err := s.db.Unscoped().Model(&SyncInstanceTask{}).Select("id, updated_at").Where("updated_at > ?", after).Order("updated_at DESC").Limit(1).Find(&task).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return nil, false, nil
// 	}
// 	return task, true, errors.New(errors.ConnectStorageError, err)
// }
