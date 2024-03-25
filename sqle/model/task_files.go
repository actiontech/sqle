package model

import (
	"time"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
)

type File struct {
	Model
	TaskId     uint   `json:"-" gorm:"index"`
	UniqueName string `json:"unique_name" gorm:"type:varchar(255)"`
	FileHost   string `json:"file_host" gorm:"type:varchar(255)"`
	NickName   string `json:"nick_name" gorm:"type:varchar(255)"`
}

const FixFilePath string = "audit_files/"

func NewFileRecord(taskID uint, nickName, uniqueName string) *File {
	return &File{
		TaskId:     taskID,
		UniqueName: uniqueName,
		FileHost:   config.GetOptions().SqleOptions.ReportHost,
		NickName:   nickName,
	}
}
func DefaultFilePath(fileName string) string {
	return FixFilePath + fileName
}

func GenUniqueFileName() string {
	return time.Now().Format("2006-01-02") + "_" + utils.GenerateRandomString(5)
}

func (s *Storage) GetFileByTaskId(taskId string) ([]*File, error) {
	files := []*File{}
	err := s.db.Where("task_id = ?", taskId).Find(&files).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return files, errors.New(errors.ConnectStorageError, err)
}
