package model

import (
	"fmt"
	"strconv"

	e "errors"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"gorm.io/gorm"
)

const globalConfigurationTablePrefix = "global_configuration"
const (
	SystemVariableSqlManageRawExpiredHours    = "system_variable_sql_manage_raw_expired_hours"
	SystemVariableWorkflowExpiredHours        = "system_variable_workflow_expired_hours"
	SystemVariableSqleUrl                     = "system_variable_sqle_url"
	SystemVariableOperationRecordExpiredHours = "system_variable_operation_record_expired_hours"
	SystemVariableCbOperationLogsExpiredHours = "system_variable_cb_operation_logs_expired_hours"
	SystemVariableSSHPrimaryKey               = "system_variable_ssh_primary_key"
)

const (
	DefaultSqlManageRawExpiredHours    = 30 * 24
	DefaultOperationRecordExpiredHours = 90 * 24
	DefaultCbOperationLogsExpiredHours = 90 * 24
)

// SystemVariable store misc K-V.
type SystemVariable struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null;type:text"`
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

	// if _, ok := sysVariables[SystemVariableWorkflowExpiredHours]; !ok {
	// 	sysVariables[SystemVariableWorkflowExpiredHours] = SystemVariable{
	// 		Key:   SystemVariableWorkflowExpiredHours,
	// 		Value: strconv.Itoa(30 * 24),
	// 	}
	// }

	if _, ok := sysVariables[SystemVariableOperationRecordExpiredHours]; !ok {
		sysVariables[SystemVariableOperationRecordExpiredHours] = SystemVariable{
			Key:   SystemVariableOperationRecordExpiredHours,
			Value: strconv.Itoa(DefaultOperationRecordExpiredHours),
		}
	}

	if _, ok := sysVariables[SystemVariableCbOperationLogsExpiredHours]; !ok {
		sysVariables[SystemVariableCbOperationLogsExpiredHours] = SystemVariable{
			Key:   SystemVariableCbOperationLogsExpiredHours,
			Value: strconv.Itoa(DefaultCbOperationLogsExpiredHours),
		}
	}

	return sysVariables, nil
}

// GetSystemVariableByKey retrieves a system variable by its key.
// Returns the system variable, a boolean indicating if it was found, and any error that occurred.
func (s *Storage) GetSystemVariableByKey(key string) (SystemVariable, bool, error) {
	var systemVariable SystemVariable

	err := s.db.Where("`key` = ?", key).First(&systemVariable).Error

	if err == gorm.ErrRecordNotFound {
		return systemVariable, false, nil
	}
	if err != nil {
		return systemVariable, false, errors.New(errors.ConnectStorageError, err)
	}

	return systemVariable, true, nil
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
func (s *Storage) GetSqlManageRawSqlExpiredHoursOrDefault() (int64, error) {
	var svs []SystemVariable
	err := s.db.Find(&svs).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}

	for _, sv := range svs {
		if sv.Key == SystemVariableSqlManageRawExpiredHours {
			expiredHs, err := strconv.ParseInt(sv.Value, 10, 64)
			if err != nil {
				return 0, err
			}
			return expiredHs, nil
		}
	}

	return DefaultSqlManageRawExpiredHours, nil
}

const (
	ImTypeDingTalk    = "dingTalk"
	ImTypeFeishuAudit = "feishu_audit"
	ImTypeWechatAudit = "wechat_audit"
	ImTypeCoding      = "coding"
)

type IM struct {
	Model
	AppKey           string `json:"app_key" gorm:"column:app_key; type:varchar(255)"`
	AppSecret        string `json:"-" gorm:"-"`
	IsEnable         bool   `json:"is_enable" gorm:"column:is_enable"`
	ProcessCode      string `json:"process_code" gorm:"column:process_code; type:varchar(255)"`
	EncryptAppSecret string `json:"encrypt_app_secret" gorm:"column:encrypt_app_secret; type:varchar(255)"`
	// 类型唯一
	Type string `json:"type" gorm:"index:unique; type:varchar(255)"`
}

func (i *IM) TableName() string {
	return fmt.Sprintf("%v_im", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create.
func (i *IM) BeforeSave(tx *gorm.DB) error {
	return i.encryptAppSecret(tx)
}

func (i *IM) encryptAppSecret(tx *gorm.DB) error {
	if i == nil {
		return nil
	}
	data, err := dmsCommonAes.AesEncrypt(i.AppSecret)
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("EncryptAppSecret", data)
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db.
func (i *IM) AfterFind(tx *gorm.DB) error {
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
	im := new(IM)
	// 使用First将IM实例化，避免BeforeSave加密空的AppSecret
	err := s.db.Where("id = ?", id).First(&im).Updates(m).Error
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
	ApproveInstanceCode string `json:"approve_instance" gorm:"column:approve_instance; type:varchar(255)"`
	WorkflowId          string `json:"workflow_id" gorm:"column:workflow_id; type:varchar(255)"`
	// 审批实例 taskID
	TaskID int64  `json:"task_id" gorm:"column:task_id"`
	Status string `json:"status" gorm:"default:\"initialized\"; type:varchar(255)"`
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
	ApproveInstanceCode string `json:"approve_instance" gorm:"column:approve_instance; type:varchar(255)"`
	WorkflowId          string `json:"workflow_id" gorm:"column:workflow_id; type:varchar(255)"`
	// 审批实例 taskID
	TaskID string `json:"task_id" gorm:"column:task_id; type:varchar(255)"`
	Status string `json:"status" gorm:"default:\"INITIALIZED\"; type:varchar(255)"`
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

type WechatOAStatus int

const (
	INITIALIZED WechatOAStatus = 1
	APPROVED    WechatOAStatus = 2
	REJECTED    WechatOAStatus = 3
)

type WechatRecord struct {
	Model
	TaskId   uint   `json:"task_id" gorm:"column:task_id"`
	OaResult string `json:"oa_result" gorm:"column:oa_result;default:\"INITIALIZED\";type:varchar(255)"`
	SpNo     string `json:"sp_no" gorm:"column:sp_no;type:varchar(255)"`

	Task *Task `gorm:"foreignkey:TaskId"`
}

func (s *Storage) GetWechatRecordByStatus(status string) ([]WechatRecord, error) {
	var wcRecords []WechatRecord
	err := s.db.Where("oa_result = ?", status).Preload("Task").Find(&wcRecords).Error
	if err != nil {
		return nil, err
	}
	return wcRecords, nil
}

func (s *Storage) WechatCancelScheduledTask(w WechatRecord) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.RejectScheduledInstanceRecord(w.TaskId)
		if err != nil {
			return err
		}

		w.OaResult = ApproveStatusRefuse
		err = s.Save(&w)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) WechatAgreeScheduledTask(w WechatRecord) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.AgreeScheduledInstanceRecord(w.TaskId)
		if err != nil {
			return err
		}

		w.OaResult = ApproveStatusAgree
		err = s.Save(&w)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) GetWechatRecordsByTaskIds(taskIds []uint) ([]*WechatRecord, error) {
	var wcRecords []*WechatRecord
	err := s.db.Where("task_id in (?)", taskIds).Find(&wcRecords).Error
	if err != nil {
		return nil, err
	}
	return wcRecords, nil
}

func (s *Storage) UpdateWechatRecordByTaskId(taskId uint, m map[string]interface{}) error {
	err := s.db.Model(&WechatRecord{}).Where("task_id = ?", taskId).Updates(m).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}

func (s *Storage) DeleteWechatRecordByTaskId(taskId uint) error {
	return s.db.Where("task_id = ?", taskId).Delete(&WechatRecord{}).Error
}

type FeishuScheduledRecord struct {
	Model
	TaskId              uint   `json:"task_id" gorm:"column:task_id"`
	OaResult            string `json:"oa_result" gorm:"column:oa_result;default:\"INITIALIZED\";type:varchar(255)"`
	ApproveInstanceCode string `json:"approve_instance_code" gorm:"column:approve_instance_code;type:varchar(255)"`

	Task *Task `gorm:"foreignkey:TaskId"`
}

func (s *Storage) UpdateFeishuScheduledByTaskId(taskId uint, m map[string]interface{}) error {
	err := s.db.Model(&FeishuScheduledRecord{}).Where("task_id = ?", taskId).Updates(m).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	return nil
}

func (s *Storage) GetFeishuScheduledByStatus(status string) ([]FeishuScheduledRecord, error) {
	var fsRecords []FeishuScheduledRecord
	err := s.db.Where("oa_result = ?", status).Preload("Task").Find(&fsRecords).Error
	if err != nil {
		return nil, err
	}
	return fsRecords, nil
}

func (s *Storage) FeishuCancelScheduledTask(f FeishuScheduledRecord) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.RejectScheduledInstanceRecord(f.TaskId)
		if err != nil {
			return err
		}

		f.OaResult = FeishuAuditStatusRejected
		err = s.Save(&f)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) FeishuAgreeScheduledTask(f FeishuScheduledRecord) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.AgreeScheduledInstanceRecord(f.TaskId)
		if err != nil {
			return err
		}

		f.OaResult = FeishuAuditStatusApprove
		err = s.Save(&f)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) GetFeishuRecordsByTaskIds(taskIds []uint) ([]*FeishuScheduledRecord, error) {
	var fsRecords []*FeishuScheduledRecord
	err := s.db.Where("task_id in (?)", taskIds).Find(&fsRecords).Error
	if err != nil {
		return nil, err
	}
	return fsRecords, nil
}

func (s *Storage) DeleteFeishuRecordByTaskId(taskId uint) error {
	return s.db.Where("task_id = ?", taskId).Delete(&FeishuScheduledRecord{}).Error
}

const (
	NotifyTypeWechat = "wechat"
	NotifyTypeFeishu = "feishu"
)

func (s *Storage) CreateNotifyRecord(notifyType string, curTaskRecord *WorkflowInstanceRecord) error {
	switch notifyType {
	case NotifyTypeWechat:
		record := WechatRecord{
			TaskId: curTaskRecord.TaskId,
		}
		if err := s.Save(&record); err != nil {
			return nil
		}
	case NotifyTypeFeishu:
		record := FeishuScheduledRecord{
			TaskId: curTaskRecord.TaskId,
		}
		if err := s.Save(&record); err != nil {
			return nil
		}
	default:
		return nil
	}
	err := s.UpdateWorkflowInstanceRecordById(curTaskRecord.ID, map[string]interface{}{"need_scheduled_task_notify": true})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CancelNotify(taskId uint) error {
	ir, err := s.GetWorkInstanceRecordByTaskId(fmt.Sprint(taskId))
	if err != nil {
		return err
	}
	// 定时上线原本不需要发送通知，就不需要再删除record记录
	if !ir.NeedScheduledTaskNotify {
		return nil
	}

	ir.NeedScheduledTaskNotify = false
	if err := s.Save(ir); err != nil {
		return err
	}

	// wechat
	{
		records, err := s.GetWechatRecordsByTaskIds([]uint{taskId})
		if err != nil {
			return err
		}
		if len(records) > 0 {
			return s.DeleteWechatRecordByTaskId(taskId)
		}
	}
	// feishu
	{
		records, err := s.GetFeishuRecordsByTaskIds([]uint{taskId})
		if err != nil {
			return err
		}
		if len(records) > 0 {
			return s.DeleteFeishuRecordByTaskId(taskId)
		}
	}
	return nil
}
