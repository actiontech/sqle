package model

import (
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"gorm.io/gorm"
)

type ReportPushConfig struct {
	Model
	ProjectId         string    `json:"project_id" gorm:"type:varchar(255)"`
	Type              string    `json:"type" gorm:"type:varchar(255)"`
	TriggerType       string    `json:"trigger_type"  gorm:"type:varchar(255)"`
	PushFrequencyCron string    `json:"cron" gorm:"type:varchar(255)"`
	PushUserType      string    `json:"push_user_Type" gorm:"type:varchar(255)"`
	PushUserList      Strings   `json:"push_user_list"`
	LastPushTime      time.Time `json:"last_push_time" gorm:"type:datetime(3)"`
	Enabled           bool      `json:"enabled" gorm:"type:varchar(255)"`
	UpdateTime        time.Time `json:"update_time" gorm:"type:datetime(3)"`
}

func (s Storage) GetReportPushConfigListInProject(projectID string) ([]ReportPushConfig, error) {
	reportPushConfigs := make([]ReportPushConfig, 0)
	err := s.db.Model(ReportPushConfig{}).Where("project_id = ?", projectID).Find(&reportPushConfigs).Error
	if err != nil {
		return nil, err
	}
	return reportPushConfigs, nil
}

const (
	// 推送报告类型
	TypeWorkflow  = "workflow"
	TypeSQLManage = "sql_manage"

	// 推送报告触发类型
	TriggerTypeImmediately = "immediately"
	TriggerTypeTiming      = "timing"

	// 推送报告指定用户类型
	PushUserTypeFixed           = "fixed"
	PushUserTypePermissionMatch = "permission_match"
)

// 新增项目需要新增的配置
func (s Storage) InitReportPushConfigInProject(projectID string) error {
	var defaultPushConfigs = []ReportPushConfig{
		{
			ProjectId:         projectID,
			Type:              TypeWorkflow,
			TriggerType:       TriggerTypeImmediately,
			PushFrequencyCron: "",
			PushUserType:      PushUserTypePermissionMatch,
			PushUserList:      []string{},
			Enabled:           true,
			LastPushTime:      time.Now(),
			UpdateTime:        time.Now(),
		}, {
			ProjectId:         projectID,
			Type:              TypeSQLManage,
			TriggerType:       TriggerTypeTiming,
			PushFrequencyCron: "",
			PushUserType:      PushUserTypeFixed,
			PushUserList:      []string{},
			Enabled:           false,
			LastPushTime:      time.Now(),
			UpdateTime:        time.Now(),
		},
	}
	err := s.db.Save(defaultPushConfigs).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetReportPushConfigByProjectId(projectId ProjectUID) (*ReportPushConfig, bool, error) {
	ReportPushConfig := &ReportPushConfig{}
	err := s.db.Where("project_id = ?", projectId).First(ReportPushConfig).Error
	if err == gorm.ErrRecordNotFound {
		return ReportPushConfig, false, nil
	}
	return ReportPushConfig, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetReportPushConfigById(id uint) (*ReportPushConfig, bool, error) {
	ReportPushConfig := &ReportPushConfig{}
	err := s.db.Where("id = ?", id).First(ReportPushConfig).Error
	if err == gorm.ErrRecordNotFound {
		return ReportPushConfig, false, nil
	}
	return ReportPushConfig, true, errors.New(errors.ConnectStorageError, err)
}

func (s Storage) GetLastUpdateReportPushConfig(lastSyncTime time.Time) ([]*ReportPushConfig, error) {
	rpcList := make([]*ReportPushConfig, 0)
	err := s.db.Model(&ReportPushConfig{}).Where("update_time > ? AND trigger_type = 'timing'", lastSyncTime).Find(&rpcList).Error
	if err != nil {
		return nil, err
	}
	return rpcList, nil
}
