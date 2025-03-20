package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"gorm.io/gorm"
)

type InstanceAuditPlan struct {
	Model
	ProjectId    ProjectUID `gorm:"index; not null"`
	InstanceID   uint64     `json:"instance_id"`
	DBType       string     `json:"db_type" gorm:"type:varchar(255) not null"`
	Token        string     `json:"token" gorm:"not null"`
	CreateUserID string     `json:"create_user_id" gorm:"type:varchar(255)"`
	ActiveStatus string     `json:"active_status" gorm:"type:varchar(255)"`

	AuditPlans []*AuditPlanV2
}

// 扫描任务状态
const (
	ActiveStatusNormal   = "normal"
	ActiveStatusDisabled = "disabled"
)

// 上一次采集执行状态
const (
	LastCollectionNormal   = "normal"
	LastCollectionAbnormal = "abnormal"
)

// TODO 推送配置
type NotifyConfig struct {
	// NotifyInterval      int    `json:"notify_interval" gorm:"default:10"`
	// NotifyLevel         string `json:"notify_level" gorm:"default:'warn'"`
	// EnableEmailNotify   bool   `json:"enable_email_notify"`
	// EnableWebHookNotify bool   `json:"enable_web_hook_notify"`
	// WebHookURL          string `json:"web_hook_url"`
	// WebHookTemplate     string `json:"web_hook_template"`
}

type AuditPlanDetail struct {
	AuditPlanV2
	ProjectId    string `json:"project_id"`
	DBType       string `json:"db_type"`
	Token        string `json:"token" `
	InstanceID   string `json:"instance_id"`
	CreateUserID string `json:"create_user_id"`

	Instance *Instance `gorm:"-"`
}

func (a AuditPlanDetail) GetInstanceName() string {
	if a.Instance == nil {
		return ""
	}
	return a.Instance.Name
}

func (s *Storage) ListActiveAuditPlanDetail() ([]*AuditPlanDetail, error) {
	var aps []*AuditPlanDetail
	err := s.db.Model(AuditPlanV2{}).Joins("JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id").
		Where("audit_plans_v2.active_status = ? AND instance_audit_plans.active_status = ?", ActiveStatusNormal, ActiveStatusNormal).
		Select("audit_plans_v2.*,instance_audit_plans.project_id,instance_audit_plans.db_type,instance_audit_plans.token,instance_audit_plans.instance_id,instance_audit_plans.create_user_id").
		Scan(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanDetailByID(id uint) (*AuditPlanDetail, error) {
	ap, exist, err := s.GetAuditPlanDetailByIDExist(id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("cant find audit plan by id %d", id)
	}
	return ap, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetActiveAuditPlanDetail(id uint) (*AuditPlanDetail, bool, error) {
	return s.getAuditPlanDetailByID(id, ActiveStatusNormal)
}

func (s *Storage) GetAuditPlanDetailByIDExist(id uint) (*AuditPlanDetail, bool, error) {
	return s.getAuditPlanDetailByID(id, "")
}

func (s *Storage) getAuditPlanDetailByID(id uint, status string) (*AuditPlanDetail, bool, error) {
	var ap *AuditPlanDetail
	query := s.db.Model(AuditPlanV2{}).Joins("JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id").
		Where("audit_plans_v2.id = ?", id).
		Select("audit_plans_v2.*,instance_audit_plans.project_id,instance_audit_plans.db_type,instance_audit_plans.token,instance_audit_plans.instance_id,instance_audit_plans.create_user_id")

	if status != "" {
		query = query.Where("audit_plans_v2.active_status = ?", status)
	}

	err := query.Scan(&ap).Error
	if err != nil {
		return ap, false, errors.New(errors.ConnectStorageError, err)
	}
	if ap == nil {
		return nil, false, nil
	}

	return ap, true, nil
}

func (s *Storage) UpdateInstanceAuditPlanByID(id uint, attrs map[string]interface{}) error {
	err := s.db.Model(InstanceAuditPlan{}).Where("id = ?", id).Updates(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

// 获取所有变更过的记录，包括删除
func (s *Storage) GetLatestAuditPlanRecordsV2(after time.Time) ([]*AuditPlanDetail, error) {
	var aps []*AuditPlanDetail
	err := s.db.Unscoped().Model(AuditPlanV2{}).Select("id, updated_at").Where("updated_at > ?", after).Order("updated_at").Find(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

type AuditPlanV2 struct {
	Model

	InstanceAuditPlanID     uint                      `json:"instance_audit_plan_id" gorm:"not null"`
	Type                    string                    `json:"type" gorm:"type:varchar(255)"`
	RuleTemplateName        string                    `json:"rule_template_name" gorm:"type:varchar(255)"`
	Params                  params.Params             `json:"params" gorm:"type:json"`
	HighPriorityParams      params.ParamsWithOperator `json:"high_priority_params" gorm:"type:json"`
	NeedMarkHighPrioritySQL bool                      `json:"need_mark_high_priority_sql"`
	ActiveStatus            string                    `json:"active_status" gorm:"type:varchar(255)"`

	AuditPlanSQLs     []*SQLManageRecord `gorm:"-"`
	AuditPlanTaskInfo *AuditPlanTaskInfo `gorm:"foreignkey:AuditPlanID"`
}

type AuditPlanTaskInfo struct {
	Model
	AuditPlanID          uint       `json:"audit_plan_id" gorm:"not null"`
	LastCollectionTime   *time.Time `json:"last_collection_time" gorm:"type:datetime(3)"`
	LastCollectionStatus string     `json:"last_collection_status" gorm:"type:varchar(25)"`
}

func (a AuditPlanV2) TableName() string {
	return "audit_plans_v2"
}

type AuditPlaner interface {
	GetBaseInfo() BaseAuditPlan
}

type BaseAuditPlan struct {
	ID               uint       `json:"id" gorm:"primary_key" example:"1"`
	ProjectId        ProjectUID `gorm:"index; not null"`
	Name             string     `json:"name" gorm:"not null;index"`
	CronExpression   string     `json:"cron_expression" gorm:"not null"`
	Type             string     `json:"type"`
	RuleTemplateName string     `json:"rule_template_name"`
}

func (a *AuditPlanV2) GetBaseInfo() BaseAuditPlan {
	return BaseAuditPlan{
		ID:               a.ID,
		Type:             a.Type,
		RuleTemplateName: a.RuleTemplateName,
	}
}

func (s *Storage) GetLatestStartTimeAuditPlanSQLV2(sourceId uint) (string, error) {
	var info = struct {
		StartTime string `gorm:"column:max_start_time"`
	}{}
	err := s.db.Raw(`SELECT MAX(STR_TO_DATE(JSON_UNQUOTE(JSON_EXTRACT(info, '$.start_time_of_last_scraped_sql')), '%Y-%m-%dT%H:%i:%s.%f')) 
					AS max_start_time FROM sql_manage_records WHERE source_id = ?`, sourceId).Scan(&info).Error
	return info.StartTime, err
}

// 此表对于来源是扫描任务的相关sql, 目前仅在采集和审核时会更新, 如有其他场景更新此表, 需要考虑更新后会触发审核影响
// 如有其他sql业务相关字段补充, 可新增至SQLManageRecordProcess中
type SQLManageRecord struct {
	Model

	Source         string         `json:"source" gorm:"type:varchar(255);index:idx_source_id_source"`
	SourceId       string         `json:"source_id" gorm:"type:varchar(255);index:idx_source_id_source"`
	ProjectId      string         `json:"project_id" gorm:"type:varchar(255)"`
	InstanceID     string         `json:"instance_id" gorm:"type:varchar(255)"`
	SchemaName     string         `json:"schema_name" gorm:"type:varchar(255)"`
	SqlFingerprint string         `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	SqlText        string         `json:"sql_text" gorm:"type:mediumtext;not null"`
	Info           JSON           `gorm:"type:json"` // 慢日志的 执行时间等特殊属性
	AuditLevel     string         `json:"audit_level" gorm:"type:varchar(255)"`
	AuditResults   *AuditResults  `json:"audit_results" gorm:"type:json"`
	SQLID          string         `json:"sql_id" gorm:"type:varchar(255);unique;not null"`
	Priority       sql.NullString `json:"priority" gorm:"type:varchar(255)"`

	SQLManager SQLManageRecordProcess
}

func (o SQLManageRecord) GetFingerprintMD5() string {
	if o.SQLID != "" {
		return o.SQLID
	}
	// 为了区分具有相同Fingerprint但Schema不同的SQL，在这里加入Schema信息进行区分
	sqlIdentityJSON, _ := json.Marshal(
		struct {
			Fingerprint string
			Schema      string
		}{
			Fingerprint: o.SqlFingerprint,
			Schema:      o.SchemaName,
		},
	)
	return utils.Md5String(string(sqlIdentityJSON))
}

type DataLock struct {
	Model
	AuditPlanId             uint64     `json:"audit_plan_id" gorm:"type:bigint unsigned;not null"`
	InstanceAuditPlanId     uint64     `json:"instance_audit_plan_id" gorm:"type:bigint unsigned;not null"`
	Engine                  string     `json:"engine" gorm:"type:varchar(255);not null"`
	DbUser                  string     `json:"db_user" gorm:"type:varchar(255)"`
	Host                    string     `json:"host" gorm:"type:varchar(255)"`
	DatabaseName            string     `json:"database_name" gorm:"type:varchar(255)"`
	ObjectName              string     `json:"object_name" gorm:"type:varchar(255)"`
	IndexType               string     `json:"index_type" gorm:"type:varchar(255)"`
	GrantedLockId           string     `json:"granted_lock_id" gorm:"type:varchar(255);not null"`
	WaitingLockId           string     `json:"waiting_lock_id" gorm:"type:varchar(255);not null"`
	LockType                string     `json:"lock_type" gorm:"type:varchar(255);not null"`
	LockMode                string     `json:"lock_mode" gorm:"type:varchar(255);not null"`
	GrantedLockConnectionId int64      `json:"granted_lock_connection_id" gorm:"type:bigint"`
	WaitingLockConnectionId int64      `json:"waiting_lock_connection_id" gorm:"type:bigint"`
	GrantedLockTrxId        int64      `json:"granted_lock_trx_id" gorm:"type:bigint"`
	WaitingLockTrxId        int64      `json:"waiting_lock_trx_id" gorm:"type:bigint"`
	GrantedLockSql          string     `json:"granted_lock_sql" gorm:"type:longtext"`
	WaitingLockSql          string     `json:"waiting_lock_sql" gorm:"type:longtext"`
	TrxStarted              *time.Time `json:"trx_started" gorm:"type:datetime"`
	TrxWaitStarted          *time.Time `json:"trx_wait_started" gorm:"type:datetime"`
}

const (
	LockType   = "lock_type"
	ObjectName = "object_name"
	Database   = "database_name"
)

func (s *Storage) SelectDistinctLockType(auditPlanId uint, instanceAuditPlanId uint) ([]string, error) {
	var lockTypes []string
	err := s.db.Model(&DataLock{}).Distinct("lock_type").
		Where("audit_plan_id = ? AND instance_audit_plan_id = ?", auditPlanId, instanceAuditPlanId).
		Find(&lockTypes).Error
	if err != nil {
		return lockTypes, err
	}
	return lockTypes, nil
}

func (s *Storage) SelectDistinctDatabase(auditPlanId uint, instanceAuditPlanId uint) ([]string, error) {
	var schemas []string
	err := s.db.Model(&DataLock{}).Distinct("database_name").
		Where("audit_plan_id = ? AND instance_audit_plan_id = ?", auditPlanId, instanceAuditPlanId).
		Find(&schemas).Error
	if err != nil {
		return schemas, err
	}
	return schemas, nil
}

func (s *Storage) SelectDistinctObjectName(auditPlanId uint, instanceAuditPlanId uint) ([]string, error) {
	var objectNames []string
	err := s.db.Model(&DataLock{}).Distinct("object_name").
		Where("audit_plan_id = ? AND instance_audit_plan_id = ?", auditPlanId, instanceAuditPlanId).
		Find(&objectNames).Error
	if err != nil {
		return objectNames, err
	}
	return objectNames, nil
}

func (s *Storage) PushSQLToDataLock(dataLocks []*DataLock) error {
	// 新增前删除，保证数据的实时性
	s.db.Exec("DELETE FROM data_locks")
	if len(dataLocks) == 0 {
		return nil
	}
	return s.db.Create(dataLocks).Error
}

func (s *Storage) GetDataLockList(filters map[string]string, limit, offset int, auditPlanId uint, instanceAuditPlanId uint) ([]*DataLock, error) {
	dataLocks := []*DataLock{}
	query := s.db.Model(&DataLock{}).Limit(limit).
		Where("audit_plan_id = ? AND instance_audit_plan_id = ?", auditPlanId, instanceAuditPlanId).
		Offset(offset)
	if filters[LockType] != "" {
		query = query.Where("lock_type = ?", filters[LockType])
	}
	if filters[ObjectName] != "" {
		query = query.Where("object_name = ?", filters[ObjectName])
	}
	if filters[Database] != "" {
		query = query.Where("database_name = ?", filters[Database])
	}

	err := query.Find(&dataLocks).Error
	if err != nil {
		return nil, err
	}
	return dataLocks, nil
}

func (s *Storage) CountDataLock(filters map[string]string, auditPlanId uint, instanceAuditPlanId uint) (int64, error) {
	var count int64
	query := s.db.Model(&DataLock{}).Where("audit_plan_id = ? AND instance_audit_plan_id = ?", auditPlanId, instanceAuditPlanId)
	if filters[LockType] != "" {
		query = query.Where("lock_type = ?", filters[LockType])
	}
	if filters[ObjectName] != "" {
		query = query.Where("object_name = ?", filters[ObjectName])
	}
	if filters[Database] != "" {
		query = query.Where("database_name = ?", filters[Database])
	}
	err := query.Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *Storage) GetManageSQLBySQLId(sqlId string) (*SQLManageRecord, bool, error) {
	sql := &SQLManageRecord{}

	err := s.db.Where("sql_id = ?", sqlId).First(sql).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return sql, true, nil
}

func (s *Storage) GetManageSQLById(id string) (*SQLManageRecord, bool, error) {
	sql := &SQLManageRecord{}

	err := s.db.Where("id = ?", id).First(sql).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return sql, true, nil
}

// 获取指定扫描任务下的所有SQL
func (s *Storage) GetManagerSQLListByAuditPlanId(auditPlanID uint) ([]*SQLManageRecord, error) {
	sqls := []*SQLManageRecord{}
	err := s.db.Joins(`
		JOIN audit_plans_v2 ON sql_manage_records.source_id = CONCAT(audit_plans_v2.instance_audit_plan_id, '')
		AND sql_manage_records.source = audit_plans_v2.type`).
		Where("audit_plans_v2.id = ?", auditPlanID).
		Find(&sqls).Error
	if err != nil {
		return nil, err
	}
	return sqls, nil
}

// 获取指定扫描任务下的所有Schema
func (s *Storage) GetManagerSqlSchemaNameByAuditPlan(auditPlanId uint) ([]string, error) {
	var metricValueTips []string
	err := s.db.Table("sql_manage_records").
		Select("DISTINCT sql_manage_records.schema_name as schema_name").
		Joins(`
		JOIN audit_plans_v2 ON sql_manage_records.source = audit_plans_v2.type 
		AND sql_manage_records.source_id = CONCAT(audit_plans_v2.instance_audit_plan_id, '')`).
		Where("audit_plans_v2.id = ?", auditPlanId).
		Scan(&metricValueTips).Error
	return metricValueTips, errors.New(errors.ConnectStorageError, err)
}

// 获取指定扫描任务下的所有SQL的指定指标
func (s *Storage) GetManagerSqlMetricTipsByAuditPlan(auditPlanId uint, metricName string) ([]string, error) {
	var metricValueTips []string
	err := s.db.Table("sql_manage_records").
		Select(fmt.Sprintf("DISTINCT sql_manage_records.info->>'$.%s' as metric_value", metricName)).
		Joins(`
			JOIN audit_plans_v2 ON sql_manage_records.source = audit_plans_v2.type 
			AND sql_manage_records.source_id = CONCAT(audit_plans_v2.instance_audit_plan_id, '')`).
		Where("audit_plans_v2.id = ? AND sql_manage_records.info->>'$.db_user' IS NOT NULL", auditPlanId).
		Scan(&metricValueTips).Error
	return metricValueTips, errors.New(errors.ConnectStorageError, err)
}

/*
TODO 优先级:高 目的: 优化该方法的SQL性能
 1. 该方法的SQL性能差，本地数据量约三四千条，SQL的响应时间为(2.56 sec)
 2. 该方法主要性能下降在rules(type ALL filtered 10.00)表和smr表(type ALL filtered 1.11)
*/
func (s *Storage) GetManagerSqlRuleTipsByAuditPlan(auditPlanId uint) ([]*SqlManageRuleTips, error) {
	sqlManageRuleTips := make([]*SqlManageRuleTips, 0)
	err := s.db.Table("sql_manage_records smr").
		Select("DISTINCT iap.db_type, rules.name as rule_name, rules.i18n_rule_info").
		Joins("JOIN audit_plans_v2 ap ON ap.instance_audit_plan_id = smr.source_id AND ap.type = smr.source").
		Joins("JOIN instance_audit_plans iap ON iap.id = ap.instance_audit_plan_id").
		Joins("LEFT JOIN rules ON rules.db_type = iap.db_type").
		Where("smr.audit_results LIKE CONCAT('%' , rules.name , '%') AND ap.id = ?", auditPlanId).
		Scan(&sqlManageRuleTips).Error
	return sqlManageRuleTips, errors.New(errors.ConnectStorageError, err)
}

type SQLManageRecordProcess struct {
	Model

	SQLManageRecordID *uint      `json:"sql_manage_record_id" gorm:"unique;not null"`
	LastAuditTime     *time.Time `json:"last_audit_time" gorm:"type:datetime(3)"`
	// 任务属性字段
	Assignees string        `json:"assignees" gorm:"type:varchar(2000)"`
	Status    ProcessStatus `json:"status" gorm:"default:\"unhandled\""`
	Remark    string        `json:"remark" gorm:"type:varchar(4000)"`
}

type ProcessStatus string

const (
	ProcessStatusUnhandled     = "unhandled"
	ProcessStatusSolved        = "solved"
	ProcessStatusIgnored       = "ignored"
	ProcessStatusManualAudited = "manual_audited"
	ProcessStatusSent          = "sent"
)

func (s *Storage) GetSQLManageRecordProcess(sqlManageRecordID uint) (*SQLManageRecordProcess, error) {
	sqlManageRecordProcess := &SQLManageRecordProcess{}
	err := s.db.Model(SQLManageRecordProcess{}).
		Where("sql_manage_record_id = ?", sqlManageRecordID).
		First(sqlManageRecordProcess).Error
	if err != nil {
		return nil, err
	}
	return sqlManageRecordProcess, nil
}

func (s *Storage) GetAuditPlanByID(auditPlanID int) (*AuditPlanV2, bool, error) {
	auditPlan := &AuditPlanV2{}
	err := s.db.Model(AuditPlanV2{}).
		Where("audit_plans_v2.id = ?", auditPlanID).Preload("AuditPlanTaskInfo").
		First(auditPlan).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlan, false, nil
	}
	return auditPlan, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanByInstanceIdAndType(instanceAuditPlanID string, auditPlanType string) (*AuditPlanV2, bool, error) {
	auditPlan := &AuditPlanV2{}
	err := s.db.Model(AuditPlanV2{}).
		Where("audit_plans_v2.instance_audit_plan_id = ?", instanceAuditPlanID).
		Where("audit_plans_v2.type = ?", auditPlanType).
		First(auditPlan).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlan, false, nil
	}
	return auditPlan, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceAuditPlanDetail(instanceAuditPlanID string) (*InstanceAuditPlan, bool, error) {
	instanceAuditPlan := &InstanceAuditPlan{}
	err := s.db.Model(InstanceAuditPlan{}).Where("id = ?", instanceAuditPlanID).Preload("AuditPlans").Preload("AuditPlans.AuditPlanTaskInfo").First(&instanceAuditPlan).Error
	if err == gorm.ErrRecordNotFound {
		return instanceAuditPlan, false, nil
	}
	return instanceAuditPlan, true, errors.New(errors.ConnectStorageError, err)
}

// 获取指定扫描任务下的所有SQL的数量
func (s *Storage) GetAuditPlanTotalSQL(auditPlanID uint) (int64, error) {
	var count int64
	err := s.db.Model(&SQLManageRecord{}).
		Joins(`
			JOIN audit_plans_v2 ON sql_manage_records.source_id = CONCAT(audit_plans_v2.instance_audit_plan_id, '')
			AND sql_manage_records.source = audit_plans_v2.type`).
		Where("audit_plans_v2.id = ?", auditPlanID).
		Count(&count).Error
	return count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) SaveInstanceAuditPlan(instAuditPlans *InstanceAuditPlan) error {
	return s.Tx(func(txDB *gorm.DB) error {
		if err := txDB.Save(instAuditPlans).Error; err != nil {
			return err
		}
		apTaskInfos := make([]*AuditPlanTaskInfo, 0, len(instAuditPlans.AuditPlans))
		for _, auditPlan := range instAuditPlans.AuditPlans {
			apTaskInfos = append(apTaskInfos, &AuditPlanTaskInfo{
				AuditPlanID: auditPlan.ID,
			})
		}
		if err := txDB.Save(apTaskInfos).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) BatchSaveAuditPlans(auditPlans []*AuditPlanV2) error {
	return s.Tx(func(txDB *gorm.DB) error {
		for _, auditPlan := range auditPlans {
			// 新增的扫描任务类型需要保存audit task info
			if auditPlan.ID == 0 {
				if err := txDB.Save(auditPlan).Error; err != nil {
					return err
				}
				apTaskInfo := &AuditPlanTaskInfo{
					AuditPlanID: auditPlan.ID,
				}
				if err := txDB.Save(apTaskInfo).Error; err != nil {
					return err
				}
			} else {
				if err := txDB.Save(auditPlan).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Storage) DeleteInstanceAuditPlan(instanceAuditPlanId string) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除队列表中数据
		err := txDB.Exec(`DELETE FROM sql_manage_queues USING sql_manage_queues
		JOIN instance_audit_plans iap ON iap.id = sql_manage_queues.source_id
		WHERE iap.ID = ?`, instanceAuditPlanId).Error
		if err != nil {
			return err
		}
		err = txDB.Exec(`UPDATE instance_audit_plans iap 
		LEFT JOIN audit_plans_v2 ap ON iap.id = ap.instance_audit_plan_id
		LEFT JOIN audit_plan_task_infos apti ON apti.audit_plan_id = ap.id
		LEFT JOIN sql_manage_records smr ON smr.source_id = ap.instance_audit_plan_id AND smr.source = ap.type
		LEFT JOIN sql_manage_record_processes smrp ON smrp.sql_manage_record_id = smr.id
		SET iap.deleted_at = now(),
		ap.deleted_at = now(),
		smr.deleted_at = now(),
		smrp.deleted_at = now(),
		apti.deleted_at = now()
		WHERE iap.ID = ?`, instanceAuditPlanId).Error
		if err != nil {
			return err
		}
		return nil
	})
}

/*
TODO 优先级:中 目的: 优化该方法的SQL性能
 1. 该SQL存在隐式转换(ap.instance_audit_plan_id=sql_manage_queues.source_id)导致性能下降，且执行计划存在大表(sql_manage_records)的全表扫描，
 2. 由于操作是删除操作，使用频率较低，优化优先级降低
*/
func (s *Storage) DeleteAuditPlan(auditPlanID int) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除队列表中数据
		err := txDB.Exec(`DELETE FROM sql_manage_queues USING sql_manage_queues
		JOIN audit_plans_v2 ap ON ap.instance_audit_plan_id=sql_manage_queues.source_id
		WHERE ap.id = ?`, auditPlanID).Error
		if err != nil {
			return err
		}
		err = txDB.Exec(`UPDATE audit_plans_v2 ap 
		LEFT JOIN audit_plan_task_infos apti ON apti.audit_plan_id = ap.id
		LEFT JOIN sql_manage_records smr ON smr.source_id = ap.instance_audit_plan_id AND smr.source = ap.type
		LEFT JOIN sql_manage_record_processes smrp ON smrp.sql_manage_record_id = smr.id
		SET ap.deleted_at = now(),
		smr.deleted_at = now(),
		smrp.deleted_at = now(),
		apti.deleted_at = now()
		WHERE  ap.id = ?`, auditPlanID).Error
		if err != nil {
			return err
		}
		return nil
	})
}

var ErrAuditPlanNotFound = errors.New(errors.DataNotExist, fmt.Errorf("cant find audit plan"))

func (s *Storage) GetAuditPlanDetailByInstAuditPlanIdAndType(instAuditPlanId string, auditPlanType string) (*AuditPlanDetail, error) {
	ap, exist, err := s.GetAuditPlanDetailByType(instAuditPlanId, auditPlanType)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, ErrAuditPlanNotFound
	}
	return ap, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanDetailByType(InstanceAuditPlanId, auditPlanType string) (*AuditPlanDetail, bool, error) {
	var auditPlanDetail *AuditPlanDetail
	err := s.db.Model(AuditPlanV2{}).Joins("JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id").
		Where("instance_audit_plans.id = ? AND audit_plans_v2.type = ?", InstanceAuditPlanId, auditPlanType).
		Select("audit_plans_v2.*,instance_audit_plans.project_id,instance_audit_plans.db_type,instance_audit_plans.token,instance_audit_plans.instance_id,instance_audit_plans.create_user_id").
		Scan(&auditPlanDetail).Error
	if err != nil {
		return auditPlanDetail, false, errors.New(errors.ConnectStorageError, err)
	}
	if auditPlanDetail == nil {
		return nil, false, nil
	}
	return auditPlanDetail, true, nil
}

func (s *Storage) GetInstanceAuditPlanByInstanceID(instanceID int64) (*InstanceAuditPlan, bool, error) {
	instanceAuditPlan := &InstanceAuditPlan{}
	err := s.db.Model(InstanceAuditPlan{}).Where("instance_id = ?", instanceID).First(&instanceAuditPlan).Error
	if err == gorm.ErrRecordNotFound {
		return instanceAuditPlan, false, nil
	}
	return instanceAuditPlan, true, errors.New(errors.ConnectStorageError, err)
}

type SQLManageQueue struct {
	Model

	Source         string `json:"source" gorm:"type:varchar(255)"` // 智能扫描SQL/快速审核SQL/IDE审核SQL/CB审核SQL
	SourceId       string `json:"source_id" gorm:"type:varchar(255)"`
	ProjectId      string `json:"project_id" gorm:"type:varchar(255)"`
	InstanceID     string `json:"instance_id" gorm:"type:varchar(255)"`
	SchemaName     string `json:"schema_name" gorm:"type:varchar(255)"`
	SqlFingerprint string `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	SqlText        string `json:"sql_text" gorm:"type:mediumtext;not null"`
	Info           JSON   `gorm:"type:json"` // 慢日志的 执行时间等特殊属性
	SQLID          string `json:"sql_id" gorm:"type:varchar(255);not null"`
}

func (s *Storage) PushSQLToManagerSQLQueue(sqls []*SQLManageQueue) error {
	if len(sqls) == 0 {
		return nil
	}
	return s.db.Create(sqls).Error
}

type SqlManageMetricRecord struct {
	Model
	SQLID                             string                              `gorm:"type:varchar(255);not null;index:idx_sql_id"`
	ExecutionCount                    int                                 `gorm:"not null;default:1;comment:该时间范围内的执行次数"`
	RecordBeginAt                     time.Time                           `gorm:"not null;index:idx_time_range,priority:1;comment:统计该度量值的起始时间"`
	RecordEndAt                       time.Time                           `gorm:"not null;index:idx_time_range,priority:2;comment:统计该度量值的终止时间"`
	MetricValues                      []*SqlManageMetricValue             `gorm:"foreignKey:SqlManageMetricRecordID;references:ID"`
	SqlManageMetricExecutePlanRecords []*SqlManageMetricExecutePlanRecord `gorm:"foreignKey:SqlManageMetricRecordID"`
}

func (SqlManageMetricRecord) TableName() string {
	return "sql_manage_metric_records"
}

type SqlManageMetricValue struct {
	SqlManageMetricRecordID uint    `gorm:"not null;index:idx_metric_record_id"`
	MetricName              string  `gorm:"not null;index:idx_metric_name;comment:统计信息类型的名称"`
	MetricValue             float64 `gorm:"type:decimal(20,4);not null;default:0.0000;comment:存储数值数据，可以是 INT, FLOAT, TIME 等"`
}

func (SqlManageMetricValue) TableName() string {
	return "sql_manage_metric_values"
}

type SqlManageMetricExecutePlanRecord struct {
	Model
	SqlManageMetricRecordID uint    `gorm:"not null"`
	SelectId                int     `gorm:"not null"`
	SelectType              string  `gorm:"type:varchar(50);not null"`
	Table                   string  `gorm:"type:varchar(255);not null"`
	Partitions              string  `gorm:"type:varchar(255)"`
	Type                    string  `gorm:"type:varchar(255);not null"`
	PossibleKeys            string  `gorm:"type:varchar(255)"`
	Key                     string  `gorm:"type:varchar(255)"`
	KeyLen                  int     `gorm:"type:int"`
	Ref                     string  `gorm:"type:varchar(50)"`
	Rows                    int     `gorm:"type:int"`
	Filtered                float64 `gorm:"type:decimal(20,4)"`
	Extra                   string  `gorm:"type:varchar(255)"`
}

func (SqlManageMetricExecutePlanRecord) TableName() string {
	return "sql_manage_metric_execute_plan_records"
}

func (s *Storage) GetSqlManageMetricRecordsByTime(sqlId string, metricName string, timeBegin, timeEnd time.Time) ([]SqlManageMetricRecord, error) {
	var records []SqlManageMetricRecord
	err := s.db.Preload("MetricValues", func(db *gorm.DB) *gorm.DB {
		return db.Where("metric_name = ?", metricName)
	}).
		Preload("SqlManageMetricExecutePlanRecords").
		Where("sql_id = ? AND record_begin_at >= ? AND record_end_at <= ?", sqlId, timeBegin, timeEnd).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *Storage) GetMaxValueSqlManageMetricRecordByTime(sqlId string, metricName string, timeBegin, timeEnd time.Time) (*SqlManageMetricRecord, bool, error) {
	var record SqlManageMetricRecord
	err := s.db.Model(&SqlManageMetricRecord{}).
		Preload("MetricValues").Preload("SqlManageMetricExecutePlanRecords").
		Joins("left join sql_manage_metric_values ON sql_manage_metric_records.id = sql_manage_metric_values.sql_manage_metric_record_id").
		Where("sql_manage_metric_records.sql_id = ? AND sql_manage_metric_records.record_begin_at >= ? AND sql_manage_metric_records.record_end_at <= ? AND sql_manage_metric_values.metric_name = ?", sqlId, timeBegin, timeEnd, metricName).
		Order("sql_manage_metric_values.metric_value DESC").
		First(&record).Error
	if err == gorm.ErrRecordNotFound {
		return &record, false, nil
	}
	return &record, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) PullSQLFromManagerSQLQueue() ([]*SQLManageQueue, error) {
	sqls := []*SQLManageQueue{}
	err := s.db.Limit(1000).Find(&sqls).Error
	return sqls, err
}

func (s *Storage) RemoveSQLFromQueue(txDB *gorm.DB, sql *SQLManageQueue) error {
	return txDB.Unscoped().Delete(sql).Error
}

func (s *Storage) DeleteSQLManageRecordBySourceId(sourceId string) error {
	return s.Tx(func(txDB *gorm.DB) error {
		err := txDB.Exec(`UPDATE sql_manage_record smr
							LEFT JOIN sql_manage_record_processes smrp ON smrp.sql_manage_record_id = smr.id
							SET smr.deleted_at = now(),
							smrp.deleted_at = now()
							WHERE smr.source_id = ?`, sourceId).Error

		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) SaveManagerSQL(txDB *gorm.DB, sql *SQLManageRecord) error {
	const query = "INSERT INTO `sql_manage_records` (`sql_id`,`source`,`source_id`,`project_id`,`instance_id`,`schema_name`,`sql_fingerprint`, `sql_text`, `info`, `audit_level`, `audit_results`,`priority`) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `source` = VALUES(source),`source_id` = VALUES(source_id),`project_id` = VALUES(project_id), `instance_id` = VALUES(instance_id), `priority` = VALUES(priority), " +
		"`schema_name` = VALUES(`schema_name`), `sql_text` = VALUES(sql_text), `sql_fingerprint` = VALUES(sql_fingerprint), `info`= VALUES(info), `audit_level`= VALUES(audit_level), `audit_results`= VALUES(audit_results)"
	return txDB.Exec(query, sql.SQLID, sql.Source, sql.SourceId, sql.ProjectId, sql.InstanceID, sql.SchemaName, sql.SqlFingerprint, sql.SqlText, sql.Info, sql.AuditLevel, sql.AuditResults, sql.Priority).Error
}

func (s *Storage) UpdateManagerSQLStatus(txDB *gorm.DB, sql *SQLManageRecord) error {
	const query = `	INSERT INTO sql_manage_record_processes (sql_manage_record_id)
	SELECT smr.id FROM sql_manage_records smr WHERE smr.sql_id = ?
	ON DUPLICATE KEY UPDATE sql_manage_record_id = VALUES(sql_manage_record_id);`
	return txDB.Exec(query, sql.SQLID).Error
}

func (s *Storage) UpdateManagerSQLBySqlId(sqlId string, sqlManageMap map[string]interface{}) error {
	err := s.db.Model(&SQLManageRecord{}).Where("sql_id = ?", sqlId).
		Updates(sqlManageMap).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateAuditPlanLastCollectionTime(auditPlanID uint, collectionTime time.Time) error {
	err := s.db.Model(AuditPlanTaskInfo{}).Where("audit_plan_id = ?", auditPlanID).Update("last_collection_time", collectionTime).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateAuditPlanInfoByAPID(id uint, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlanTaskInfo{}).Where("audit_plan_id = ?", id).Updates(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlansByProjectId(projectID string) ([]*InstanceAuditPlan, error) {
	instanceAuditPlan := []*InstanceAuditPlan{}
	err := s.db.Model(InstanceAuditPlan{}).Where("project_id = ?", projectID).Find(&instanceAuditPlan).Error
	return instanceAuditPlan, err
}

func (s *Storage) GetNormalAuditPlanInstancesByLastCollectionStatus(projectID, status string) ([]*InstanceAuditPlan, error) {
	instanceAuditPlan := []*InstanceAuditPlan{}
	err := s.db.Model(InstanceAuditPlan{}).
		Distinct("instance_audit_plans.instance_id").
		Select("instance_audit_plans.*").
		Joins("JOIN audit_plans_v2 ap ON instance_audit_plans.id = ap.instance_audit_plan_id").
		Joins("JOIN audit_plan_task_infos apti ON ap.id = apti.audit_plan_id").
		Where("instance_audit_plans.project_id = ? AND apti.last_collection_status = ? AND ap.active_status = ?", projectID, status, ActiveStatusNormal).Find(&instanceAuditPlan).Error
	return instanceAuditPlan, err
}

/*
TODO 优先级:低 目的: 优化该方法的SQL性能
	1. 该SQL存在隐式转换(sql_manage_records.source_id = apv.instance_audit_plan_id)导致索引过滤率下降
	2. 由于数据量不是很大，优化优先级降低
*/
// 获取需要审核的sql，
// 当更新时间大于最后审核时间或最后审核时间为空时需要重新审核（采集或重新采集到的sql）
// 需要注意：当前仅在采集和审核时会更sql_manage_records中扫描任务相关的sql，所以使用了updated_at > last_audit_time条件。
func (s *Storage) GetSQLsToAuditFromManage() ([]*SQLManageRecord, error) {
	manageRecords := []*SQLManageRecord{}
	err := s.db.Limit(1000).Model(SQLManageRecord{}).
		Joins("JOIN audit_plans_v2 apv ON sql_manage_records.source_id = apv.instance_audit_plan_id AND sql_manage_records.source = apv.type AND apv.deleted_at IS NULL").
		Joins("JOIN sql_manage_record_processes smrp ON sql_manage_records.id =smrp.sql_manage_record_id").
		Where("sql_manage_records.updated_at > smrp.last_audit_time OR smrp.last_audit_time IS NULL").
		Find(&manageRecords).Error
	if err == gorm.ErrRecordNotFound {
		return manageRecords, nil
	}
	return manageRecords, err
}

func (s *Storage) UpdateManageSQLProcessByManageIDs(ids []uint, attrs map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}
	err := s.db.Model(SQLManageRecordProcess{}).Where("sql_manage_record_id IN (?)", ids).Updates(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}
