package model

import (
	"time"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
)

type AuditFile struct {
	Model
	TaskId     uint   `json:"-" gorm:"index"`
	UniqueName string `json:"unique_name" gorm:"type:varchar(255)"`
	FileHost   string `json:"file_host" gorm:"type:varchar(255)"`
	FileName   string `json:"file_name" gorm:"type:varchar(255)"`
}

const FixFilePath string = "audit_files/"

func NewFileRecord(taskID uint, nickName, uniqueName string) *AuditFile {
	return &AuditFile{
		TaskId:     taskID,
		UniqueName: uniqueName,
		FileHost:   config.GetOptions().SqleOptions.ReportHost,
		FileName:   nickName,
	}
}
func DefaultFilePath(fileName string) string {
	return FixFilePath + fileName
}

func GenUniqueFileName() string {
	return time.Now().Format("2006-01-02") + "_" + utils.GenerateRandomString(5)
}

func (s *Storage) GetFileByTaskId(taskId string) ([]*AuditFile, error) {
	files := []*AuditFile{}
	err := s.db.Where("task_id = ?", taskId).Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return files, errors.New(errors.ConnectStorageError, err)
}

// expired time 24hour
func (s *Storage) GetExpiredFileWithNoWorkflow() ([]AuditFile, error) {
	files := []AuditFile{}
	err := s.db.Model(&AuditFile{}).
		Joins("LEFT JOIN workflow_instance_records wir ON files.task_id = wir.task_id").
		Where("wir.task_id IS NULL AND files.deleted_at IS NULL"). // 删除没有提交为工单的SQL文件
		Where("files.file_host = ?", config.GetOptions().SqleOptions.ReportHost).
		Where("files.created_at < ?", time.Now().Add(-24*time.Hour)). // 减少提交前文件就被删除的几率
		Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return files, errors.New(errors.ConnectStorageError, err)
}

// expired time 7*24hour
func (s *Storage) GetExpiredFile() ([]AuditFile, error) {
	files := []AuditFile{}
	err := s.db.Model(&AuditFile{}).
		Where("files.file_host = ?", config.GetOptions().SqleOptions.ReportHost).
		Where("files.created_at < ?", time.Now().Add(-7*24*time.Hour)). // 减少提交前文件就被删除的几率
		Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return files, errors.New(errors.ConnectStorageError, err)
}
