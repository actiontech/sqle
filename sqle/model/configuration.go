package model

import (
	"fmt"
	"strconv"

	e "errors"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/jinzhu/gorm"
)

const globalConfigurationTablePrefix = "global_configuration"
const (
	SystemVariableWorkflowExpiredHours        = "system_variable_workflow_expired_hours"
	SystemVariableSqleUrl                     = "system_variable_sqle_url"
	SystemVariableOperationRecordExpiredHours = "system_variable_operation_record_expired_hours"
)

const (
	DefaultOperationRecordExpiredHours = 90 * 24
)

// SystemVariable store misc K-V.
type SystemVariable struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

func (s *Storage) PathSaveSystemVariables(systemVariables []SystemVariable) error {
	return s.Tx(func(tx *gorm.DB) error {
		for _, v := range systemVariables {
			if err := tx.Save(&v).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) GetAllSystemVariables() (map[string]SystemVariable, error) {
	var svs []SystemVariable
	if s.db.Find(&svs).Error != nil {
		return nil, errors.New(errors.ConnectStorageError, s.db.Error)
	}

	sysVariables := make(map[string] /*system variable key*/ SystemVariable, len(svs))
	for _, sv := range svs {
		sysVariables[sv.Key] = sv
	}

	if _, ok := sysVariables[SystemVariableWorkflowExpiredHours]; !ok {
		sysVariables[SystemVariableWorkflowExpiredHours] = SystemVariable{
			Key:   SystemVariableWorkflowExpiredHours,
			Value: strconv.Itoa(30 * 24),
		}
	}

	if _, ok := sysVariables[SystemVariableOperationRecordExpiredHours]; !ok {
		sysVariables[SystemVariableOperationRecordExpiredHours] = SystemVariable{
			Key:   SystemVariableOperationRecordExpiredHours,
			Value: strconv.Itoa(DefaultOperationRecordExpiredHours),
		}
	}

	return sysVariables, nil
}

func (s *Storage) GetSqleUrl() (string, error) {
	sys, err := s.GetAllSystemVariables()
	if err != nil {
		return "", err
	}
	return sys[SystemVariableSqleUrl].Value, nil
}

func (s *Storage) GetWorkflowExpiredHoursOrDefault() (int64, error) {
	var svs []SystemVariable
	err := s.db.Find(&svs).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}

	for _, sv := range svs {
		if sv.Key == SystemVariableWorkflowExpiredHours {
			wfExpiredHs, err := strconv.ParseInt(sv.Value, 10, 64)
			if err != nil {
				return 0, err
			}
			return wfExpiredHs, nil
		}
	}

	return 30 * 24, nil
}

const (
	ImTypeDingTalk    = "dingTalk"
	ImTypeFeishu      = "feishu"
	ImTypeFeishuAudit = "feishu_audit"
)

type IM struct {
	Model
	AppKey           string `json:"app_key" gorm:"column:app_key"`
	AppSecret        string `json:"-" gorm:"-"`
	IsEnable         bool   `json:"is_enable" gorm:"column:is_enable"`
	ProcessCode      string `json:"process_code" gorm:"column:process_code"`
	EncryptAppSecret string `json:"encrypt_app_secret" gorm:"column:encrypt_app_secret"`
	// 类型唯一
	Type string `json:"type" gorm:"unique"`
}

func (i *IM) TableName() string {
	return fmt.Sprintf("%v_im", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create.
func (i *IM) BeforeSave() error {
	return i.encryptAppSecret()
}

func (i *IM) encryptAppSecret() error {
	if i == nil {
		return nil
	}
	data, err := dmsCommonAes.AesEncrypt(i.AppSecret)
	if err != nil {
		return err
	}
	i.EncryptAppSecret = data
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db.
func (i *IM) AfterFind() error {
	err := i.decryptAppSecret()
	if err != nil {
		log.NewEntry().Errorf("decrypt app secret for IM server configuration failed, error: %v", err)
	}
	return nil
}

func (i *IM) decryptAppSecret() error {
	if i == nil {
		return nil
	}
	if i.AppSecret == "" {
		data, err := dmsCommonAes.AesDecrypt(i.EncryptAppSecret)
		if err != nil {
			return err
		}
		i.AppSecret = data
	}
	return nil
}

func (s *Storage) GetImConfigByType(imType string) (*IM, bool, error) {
	im := new(IM)
	err := s.db.Where("type = ?", imType).First(&im).Error
	if err == gorm.ErrRecordNotFound {
		return im, false, nil
	}
	return im, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAllIMConfig() ([]IM, error) {
	var ims []IM
	err := s.db.Find(&ims).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return ims, nil
}

func (s *Storage) UpdateImConfigById(id uint, m map[string]interface{}) error {
	err := s.db.Model(&IM{}).Where("id = ?", id).Updates(m).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}

const (
	ApproveStatusInitialized = "initialized"
	ApproveStatusAgree       = "agree"
	ApproveStatusRefuse      = "refuse"
	ApproveStatusCancel      = "canceled"
)

type DingTalkInstance struct {
	Model
	ApproveInstanceCode string `json:"approve_instance" gorm:"column:approve_instance"`
	WorkflowId          string `json:"workflow_id" gorm:"column:workflow_id"`
	// 审批实例 taskID
	TaskID int64  `json:"task_id" gorm:"column:task_id"`
	Status string `json:"status" gorm:"default:\"initialized\""`
}

func (s *Storage) GetDingTalkInstanceByWorkflowID(workflowId string) (*DingTalkInstance, bool, error) {
	dti := new(DingTalkInstance)
	err := s.db.Where("workflow_id = ?", workflowId).Last(&dti).Error
	if err == gorm.ErrRecordNotFound {
		return dti, false, nil
	}
	return dti, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetDingTalkInstanceListByWorkflowIDs(workflowIds []string) ([]DingTalkInstance, error) {
	var dingTalkInstances []DingTalkInstance
	err := s.db.Model(&DingTalkInstance{}).Where("workflow_id IN (?)", workflowIds).Find(&dingTalkInstances).Error
	if err != nil {
		return nil, err
	}
	return dingTalkInstances, nil
}

// batch updates ding_talk_instances'status into input status by workflow_ids, the status should be like ApproveStatusXXX in model package.
func (s *Storage) BatchUpdateStatusOfDingTalkInstance(workflowIds []string, status string) error {
	err := s.db.Model(&DingTalkInstance{}).Where("workflow_id IN (?)", workflowIds).Updates(map[string]interface{}{"status": status}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetDingTalkInstByStatus(status string) ([]DingTalkInstance, error) {
	var dingTalkInstances []DingTalkInstance
	err := s.db.Where("status = ?", status).Find(&dingTalkInstances).Error
	if err != nil {
		return nil, err
	}
	return dingTalkInstances, nil
}

const (
	FeishuAuditStatusInitialized = "INITIALIZED"
	FeishuAuditStatusApprove     = "APPROVED"
	FeishuAuditStatusRejected    = "REJECTED"
)

type FeishuInstance struct {
	Model
	ApproveInstanceCode string `json:"approve_instance" gorm:"column:approve_instance"`
	WorkflowId          string `json:"workflow_id" gorm:"column:workflow_id"`
	// 审批实例 taskID
	TaskID string `json:"task_id" gorm:"column:task_id"`
	Status string `json:"status" gorm:"default:\"INITIALIZED\""`
}

func (s *Storage) GetFeishuInstanceListByWorkflowIDs(workflowIds []string) ([]FeishuInstance, error) {
	var feishuInstList []FeishuInstance
	err := s.db.Model(&FeishuInstance{}).Where("workflow_id IN (?)", workflowIds).Find(&feishuInstList).Error
	if err != nil {
		return nil, err
	}
	return feishuInstList, nil
}

func (s *Storage) BatchUpdateStatusOfFeishuInstance(workflowIds []string, status string) error {
	err := s.db.Model(&FeishuInstance{}).Where("workflow_id IN (?)", workflowIds).Updates(map[string]interface{}{"status": status}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetFeishuInstanceByWorkflowID(workflowId string) (*FeishuInstance, bool, error) {
	fi := new(FeishuInstance)
	err := s.db.Where("workflow_id = ?", workflowId).Last(&fi).Error
	if e.Is(err, gorm.ErrRecordNotFound) {
		return fi, false, nil
	}
	return fi, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetFeishuInstByStatus(status string) ([]FeishuInstance, error) {
	var feishuInst []FeishuInstance
	err := s.db.Where("status = ?", status).Find(&feishuInst).Error
	if err != nil {
		return nil, err
	}
	return feishuInst, nil
}
