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

const (
	ActiveStatusNormal   = "normal"
	ActiveStatusDisabled = "disabled"
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

// GetLatestAuditPlanIds 获取所有变更过的记录，包括删除
// 采集时会更新last_collection_time会同步更新updated_at，此处获取updated_at > last_collection_time的任务，即为配置变更过的任务
// 影响：会查出所有被删除的任务，在syncTask时做一次额外的删除操作
func (s *Storage) GetLatestAuditPlanRecordsV2() ([]*AuditPlanDetail, error) {
	var aps []*AuditPlanDetail
	err := s.db.Unscoped().Model(AuditPlanV2{}).Select("audit_plans_v2.id, audit_plans_v2.updated_at,audit_plans_v2.last_collection_time").
		Where("(audit_plans_v2.updated_at > audit_plans_v2.last_collection_time OR last_collection_time IS NULL)").Order("updated_at").Scan(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

type AuditPlanV2 struct {
	Model

	InstanceAuditPlanID     uint                      `json:"instance_audit_plan_id" gorm:"not null"`
	Type                    string                    `json:"type" gorm:"type:varchar(255)"`
	RuleTemplateName        string                    `json:"rule_template_name" gorm:"type:varchar(255)"`
	Params                  params.Params             `json:"params" gorm:"type:varchar(1000)"`
	HighPriorityParams      params.ParamsWithOperator `json:"high_priority_params" gorm:"type:varchar(1000)"`
	NeedMarkHighPrioritySQL bool                      `json:"need_mark_high_priority_sql"`
	ActiveStatus            string                    `json:"active_status" gorm:"type:varchar(255)"`
	LastCollectionTime      *time.Time                `json:"last_collection_time" gorm:"type:datetime(3)"`

	AuditPlanSQLs []*SQLManageRecord `gorm:"foreignKey:SourceId"`
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

type SQLManageRecord struct {
	Model

	Source         string         `json:"source" gorm:"type:varchar(255)"`
	SourceId       uint           `json:"source_id" gorm:"type:varchar(255)"`
	ProjectId      string         `json:"project_id" gorm:"type:varchar(255)"`
	InstanceID     string         `json:"instance_id" gorm:"type:varchar(255)"`
	SchemaName     string         `json:"schema_name" gorm:"type:varchar(255)"`
	SqlFingerprint string         `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	SqlText        string         `json:"sql_text" gorm:"type:mediumtext;not null"`
	Info           JSON           `gorm:"type:json"` // 慢日志的 执行时间等特殊属性
	AuditLevel     string         `json:"audit_level" gorm:"type:varchar(255)"`
	AuditResults   AuditResults   `json:"audit_results" gorm:"type:json"`
	EndPoint       string         `json:"endpoint" gorm:"type:varchar(255)"`
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

func (s *Storage) GetManageSQLById(sqlId string) (*SQLManageRecord, bool, error) {
	sql := &SQLManageRecord{}

	err := s.db.Where("id = ?", sqlId).First(sql).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return sql, true, nil
}

func (s *Storage) GetManagerSQLListByAuditPlanId(apId uint) ([]*SQLManageRecord, error) {
	sqls := []*SQLManageRecord{}
	err := s.db.Where("source_id = ?", apId).Find(&sqls).Error
	return sqls, err
}

func (s *Storage) GetManagerSqlSchemaNameByAuditPlan(auditPlanId uint) ([]string, error) {
	var metricValueTips []string
	err := s.db.Table("sql_manage_records").Where("sql_manage_records.source_id = ?", auditPlanId).
		Select("DISTINCT sql_manage_records.schema_name as schema_name").
		Scan(&metricValueTips).Error
	return metricValueTips, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetManagerSqlMetricTipsByAuditPlan(auditPlanId uint, metricName string) ([]string, error) {
	var metricValueTips []string
	err := s.db.Table("sql_manage_records").Where("sql_manage_records.source_id = ?", auditPlanId).
		Select(fmt.Sprintf("DISTINCT sql_manage_records.info->>'$.%s' as metric_value", metricName)).
		Scan(&metricValueTips).Error
	return metricValueTips, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetManagerSqlRuleTipsByAuditPlan(auditPlanId uint) ([]*SqlManageRuleTips, error) {
	sqlManageRuleTips := make([]*SqlManageRuleTips, 0)
	err := s.db.Table("sql_manage_records smr").
		Joins("LEFT JOIN audit_plans_v2 ap ON ap.id = smr.source_id").
		Joins("LEFT JOIN instance_audit_plans iap ON iap.id = ap.instance_audit_plan_id").
		Joins("LEFT JOIN rules ON rules.db_type = iap.db_type").
		Where("smr.audit_results LIKE CONCAT('%' , rules.name , '%') AND smr.source_id = ?", auditPlanId).
		Select("DISTINCT iap.db_type, rules.name as rule_name, rules.desc").
		Scan(&sqlManageRuleTips).Error
	return sqlManageRuleTips, errors.New(errors.ConnectStorageError, err)
}

type SQLManageRecordProcess struct {
	Model

	SQLManageRecordID *uint `json:"sql_manage_record_id" gorm:"unique;not null"`
	// 任务属性字段
	Assignees string `json:"assignees" gorm:"type:varchar(255)"`
	Status    string `json:"status" gorm:"default:\"unhandled\""`
	Remark    string `json:"remark" gorm:"type:varchar(4000)"`
}

func (s *Storage) GetAuditPlanByID(auditPlanID int) (*AuditPlanV2, bool, error) {
	auditPlan := &AuditPlanV2{}
	err := s.db.Model(AuditPlanV2{}).
		Where("audit_plans_v2.id = ?", auditPlanID).
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
	err := s.db.Model(InstanceAuditPlan{}).Where("id = ?", instanceAuditPlanID).Preload("AuditPlans").First(&instanceAuditPlan).Error
	if err == gorm.ErrRecordNotFound {
		return instanceAuditPlan, false, nil
	}
	return instanceAuditPlan, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanTotalSQL(sourceID uint) (int64, error) {
	var count int64
	err := s.db.Model(&SQLManageRecord{}).Where("source_id = ?", sourceID).Count(&count).Error
	return count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) BatchSaveAuditPlans(auditPlans []*AuditPlanV2) error {
	return s.Tx(func(txDB *gorm.DB) error {
		for _, auditPlan := range auditPlans {
			if err := txDB.Save(auditPlan).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) DeleteInstanceAuditPlan(instanceAuditPlanId string) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除队列表中数据
		err := txDB.Exec(`DELETE FROM sql_manage_queues USING sql_manage_queues
		JOIN audit_plans_v2 ap ON ap.id=sql_manage_queues.source_id
		JOIN instance_audit_plans iap ON iap.id = ap.instance_audit_plan_id
		WHERE iap.ID = ?`, instanceAuditPlanId).Error
		if err != nil {
			return err
		}
		err = txDB.Exec(`UPDATE instance_audit_plans iap 
		LEFT JOIN audit_plans_v2 ap ON iap.id = ap.instance_audit_plan_id
		LEFT JOIN sql_manage_records oms ON oms.source_id = ap.id
		LEFT JOIN sql_manage_record_processes sm ON sm.sql_manage_record_id = oms.id
		SET iap.deleted_at = now(),
		ap.deleted_at = now(),
		oms.deleted_at = now(),
		sm.deleted_at = now()
		WHERE iap.ID = ?`, instanceAuditPlanId).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) DeleteAuditPlan(auditPlanID int) error {
	return s.Tx(func(txDB *gorm.DB) error {
		// 删除队列表中数据
		err := txDB.Exec(`DELETE FROM sql_manage_queues USING sql_manage_queues
		JOIN audit_plans_v2 ap ON ap.id=sql_manage_queues.source_id
		WHERE ap.id = ?`, auditPlanID).Error
		if err != nil {
			return err
		}
		err = txDB.Exec(`UPDATE audit_plans_v2 ap 
		LEFT JOIN sql_manage_records oms ON oms.source_id = ap.id
		LEFT JOIN sql_manage_record_processes sm ON sm.sql_manage_record_id = oms.id
		SET ap.deleted_at = now(),
		oms.deleted_at = now(),
		sm.deleted_at = now()
		WHERE  ap.id = ?`, auditPlanID).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) GetAuditPlanDetailByType(InstanceAuditPlanId, auditPlanType string) (*AuditPlanDetail, bool, error) {
	var auditPlanDetail *AuditPlanDetail
	err := s.db.Model(AuditPlanV2{}).Joins("JOIN instance_audit_plans ON instance_audit_plans.id = audit_plans_v2.instance_audit_plan_id").
		Where("instance_audit_plans.id = ? AND audit_plans_v2.type = ?", InstanceAuditPlanId, auditPlanType).
		Select("audit_plans_v2.*,instance_audit_plans.project_id,instance_audit_plans.db_type,instance_audit_plans.token,instance_audit_plans.instance_id,instance_audit_plans.create_user_id").
		Scan(&auditPlanDetail).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlanDetail, false, nil
	}
	return auditPlanDetail, true, errors.New(errors.ConnectStorageError, err)
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
	SourceId       uint   `json:"source_id" gorm:"type:varchar(255)"`
	ProjectId      string `json:"project_id" gorm:"type:varchar(255)"`
	InstanceID     string `json:"instance_id" gorm:"type:varchar(255)"`
	SchemaName     string `json:"schema_name" gorm:"type:varchar(255)"`
	SqlFingerprint string `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	SqlText        string `json:"sql_text" gorm:"type:mediumtext;not null"`
	Info           JSON   `gorm:"type:json"` // 慢日志的 执行时间等特殊属性
	EndPoint       string `json:"endpoint" gorm:"type:varchar(255)"`

	SQLID string `json:"sql_id" gorm:"type:varchar(255);not null"`
}

func (s *Storage) PushSQLToManagerSQLQueue(sqls []*SQLManageQueue) error {
	if sqls == nil || len(sqls) == 0 {
		return nil
	}
	return s.db.Create(sqls).Error
}

func (s *Storage) PullSQLFromManagerSQLQueue() ([]*SQLManageQueue, error) {
	sqls := []*SQLManageQueue{}
	err := s.db.Find(&sqls).Limit(100).Error
	return sqls, err
}

func (s *Storage) RemoveSQLFromQueue(sql *SQLManageQueue) error {
	return s.db.Unscoped().Delete(sql).Error
}

func (s *Storage) UpdateManagerSQL(sql *SQLManageRecord) error {
	const query = "INSERT INTO `sql_manage_records` (`sql_id`,`source`,`source_id`,`project_id`,`instance_id`,`schema_name`,`sql_fingerprint`, `sql_text`, `info`, `audit_level`, `audit_results`,`priority`) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `source` = VALUES(source),`source_id` = VALUES(source_id),`project_id` = VALUES(project_id), `instance_id` = VALUES(instance_id), `priority` = VALUES(priority), " +
		"`schema_name` = VALUES(`schema_name`), `sql_text` = VALUES(sql_text), `sql_fingerprint` = VALUES(sql_fingerprint), `info`= VALUES(info), `audit_level`= VALUES(audit_level), `audit_results`= VALUES(audit_results)"
	return s.db.Exec(query, sql.SQLID, sql.Source, sql.SourceId, sql.ProjectId, sql.InstanceID, sql.SchemaName, sql.SqlFingerprint, sql.SqlText, sql.Info, sql.AuditLevel, sql.AuditResults, sql.Priority).Error
}

func (s *Storage) UpdateManagerSQLStatus(sql *SQLManageRecord) error {
	const query = `	INSERT INTO sql_manage_record_processes (sql_manage_record_id)
	SELECT oms.id FROM sql_manage_records oms WHERE oms.sql_id = ?
	ON DUPLICATE KEY UPDATE sql_manage_record_id = VALUES(sql_manage_record_id);`
	return s.db.Exec(query, sql.SQLID).Error
}

func (s *Storage) UpdateAuditPlanLastCollectionTime(auditPlanID uint, collectionTime time.Time) error {
	const query = `UPDATE audit_plans_v2 SET last_collection_time = now(3) WHERE id = ?;`
	return s.db.Exec(query, auditPlanID).Error
}
