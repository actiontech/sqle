package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

type AuditPlan struct {
	Model
	ProjectId        uint   `gorm:"index; not null"`
	Name             string `json:"name" gorm:"not null;index"`
	CronExpression   string `json:"cron_expression" gorm:"not null"`
	DBType           string `json:"db_type" gorm:"not null"`
	Token            string `json:"token" gorm:"not null"`
	InstanceName     string `json:"instance_name"`
	CreateUserID     uint
	InstanceDatabase string        `json:"instance_database"`
	Type             string        `json:"type"`
	RuleTemplateName string        `json:"rule_template_name"`
	Params           params.Params `json:"params" gorm:"type:varchar(1000)"`

	NotifyInterval      int    `json:"notify_interval" gorm:"default:10"`
	NotifyLevel         string `json:"notify_level" gorm:"default:'warn'"`
	EnableEmailNotify   bool   `json:"enable_email_notify"`
	EnableWebHookNotify bool   `json:"enable_web_hook_notify"`
	WebHookURL          string `json:"web_hook_url"`
	WebHookTemplate     string `json:"web_hook_template"`

	CreateUser    *User             `gorm:"foreignkey:CreateUserId"`
	Instance      *Instance         `gorm:"foreignkey:InstanceName;association_foreignkey:Name"`
	AuditPlanSQLs []*AuditPlanSQLV2 `gorm:"foreignkey:AuditPlanID"`
}

type AuditPlanSQLV2 struct {
	Model

	// add unique index on fingerprint and audit_plan_id
	// it's done by AutoMigrate() because gorm can't create index on TEXT column directly by tag.
	AuditPlanID    uint   `json:"audit_plan_id" gorm:"not null"`
	Fingerprint    string `json:"fingerprint" gorm:"type:text;not null"`
	FingerprintMD5 string `json:"fingerprint_md5" gorm:"column:fingerprint_md5;not null"`
	SQLContent     string `json:"sql" gorm:"type:mediumtext;not null"`
	Info           JSON   `gorm:"type:json"`
	Schema         string `json:"schema" gorm:"type:varchar(512);not null"`
}

func (a AuditPlanSQLV2) TableName() string {
	return "audit_plan_sqls_v2"
}

func (a *AuditPlanSQLV2) GetFingerprintMD5() string {
	if a.FingerprintMD5 != "" {
		return a.FingerprintMD5
	}
	return utils.Md5String(a.Fingerprint)
}

// BeforeSave is a hook implement gorm model before exec create.
func (a *AuditPlanSQLV2) BeforeSave() error {
	a.FingerprintMD5 = a.GetFingerprintMD5()
	return nil
}

func (s *Storage) GetAuditPlans() ([]*AuditPlan, error) {
	var aps []*AuditPlan
	err := s.db.Model(AuditPlan{}).Find(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetActiveAuditPlans() ([]*AuditPlan, error) {
	var aps []*AuditPlan
	err := s.db.Model(AuditPlan{}).
		Joins("LEFT JOIN projects ON projects.id = audit_plans.project_id").
		Where(fmt.Sprintf("projects.status = '%v'", ProjectStatusActive)).
		Find(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanByName(name string) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanById(id uint) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).Where("id = ?", id).Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetActiveAuditPlanById(id uint) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).
		Joins("LEFT JOIN projects ON projects.id = audit_plans.project_id").
		Where(fmt.Sprintf("projects.status = '%v'", ProjectStatusActive)).
		Where("audit_plans.id = ?", id).Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanFromProjectByName(projectName, AuditPlanName string) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).Joins("LEFT JOIN projects ON projects.id = audit_plans.project_id").
		Where("projects.name = ? AND audit_plans.name = ?", projectName, AuditPlanName).Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanReportByID(auditPlanId, id uint) (*AuditPlanReportV2, bool, error) {
	ap := &AuditPlanReportV2{}
	err := s.db.Model(AuditPlanReportV2{}).Where("id = ? AND audit_plan_id = ?", id, auditPlanId).Preload("AuditPlan").Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanSQLs(auditPlanId uint) ([]*AuditPlanSQLV2, error) {
	var sqls []*AuditPlanSQLV2
	err := s.db.Model(AuditPlanSQLV2{}).Where("audit_plan_id = ?", auditPlanId).Find(&sqls).Error
	return sqls, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetLatestStartTimeAuditPlanSQL(auditPlanId uint) (string, error) {
	var info = struct {
		StartTime string `gorm:"column:max_start_time"`
	}{}
	err := s.db.Raw(`SELECT MAX(STR_TO_DATE(JSON_UNQUOTE(JSON_EXTRACT(info, '$.start_time')), '%Y-%m-%dT%H:%i:%s.%fZ')) 
					AS max_start_time FROM audit_plan_sqls_v2 WHERE audit_plan_id = ?`, auditPlanId).Scan(&info).Error
	return info.StartTime, err
}

func (s *Storage) OverrideAuditPlanSQLs(auditPlanId uint, sqls []*AuditPlanSQLV2) error {
	err := s.db.Unscoped().
		Model(AuditPlanSQLV2{}).
		Where("audit_plan_id = ?", auditPlanId).
		Delete(&AuditPlanSQLV2{}).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	raw, args := getBatchInsertRawSQL(auditPlanId, sqls)
	return errors.New(errors.ConnectStorageError, s.db.Exec(fmt.Sprintf("%v;", raw), args...).Error)
}

func (s *Storage) UpdateDefaultAuditPlanSQLs(auditPlanId uint, sqls []*AuditPlanSQLV2) error {
	raw, args := getBatchInsertRawSQL(auditPlanId, sqls)
	// counter column is a accumulate value when update.
	raw += `
ON DUPLICATE KEY UPDATE sql_content = VALUES(sql_content),
                        info        = JSON_SET(COALESCE(info, '{}'),
                                              '$.counter', CAST(COALESCE(JSON_EXTRACT(values(info), '$.counter'), 0) +
                                                                COALESCE(JSON_EXTRACT(info, '$.counter'), 0) AS SIGNED),
                                              '$.last_receive_timestamp',
                                              JSON_EXTRACT(values(info), '$.last_receive_timestamp'));`

	return errors.New(errors.ConnectStorageError, s.db.Exec(raw, args...).Error)
}

func (s *Storage) UpdateSlowLogCollectAuditPlanSQLs(auditPlanId uint, sqls []*AuditPlanSQLV2) error {
	raw, args := getBatchInsertRawSQL(auditPlanId, sqls)
	// counter column is a accumulate value when update.
	raw += `
ON DUPLICATE KEY UPDATE sql_content = VALUES(sql_content),
                        info        = JSON_SET(COALESCE(info, '{}'),
											  '$.counter', CAST(COALESCE(JSON_EXTRACT(values(info), '$.counter'), 0) +
                                                                COALESCE(JSON_EXTRACT(info, '$.counter'), 0) AS SIGNED),
                                              '$.last_receive_timestamp',
                                              JSON_EXTRACT(values(info), '$.last_receive_timestamp'),
											  '$.average_query_time',
											  CAST(
												((JSON_EXTRACT(info, '$.average_query_time') + 0) * (JSON_EXTRACT(info, '$.counter'))
												+ (JSON_EXTRACT(VALUES(info), '$.average_query_time') + 0) * (JSON_EXTRACT(VALUES(info), '$.counter')))
												/ (JSON_EXTRACT(info, '$.counter') + JSON_EXTRACT(VALUES(info), '$.counter'))
												AS UNSIGNED
											  ),
											  '$.start_time',
											  JSON_EXTRACT(values(info), '$.start_time'));`

	return errors.New(errors.ConnectStorageError, s.db.Exec(raw, args...).Error)
}

func getBatchInsertRawSQL(auditPlanId uint, sqls []*AuditPlanSQLV2) (raw string, args []interface{}) {
	pattern := make([]string, 0, len(sqls))
	for _, sql := range sqls {
		pattern = append(pattern, "(?, ?, ?, ?, ?, ?)")
		args = append(args, auditPlanId, sql.GetFingerprintMD5(), sql.Fingerprint, sql.SQLContent, sql.Info, sql.Schema)
	}
	raw = fmt.Sprintf("INSERT INTO `audit_plan_sqls_v2` (`audit_plan_id`,`fingerprint_md5`, `fingerprint`, `sql_content`, `info`, `schema`) VALUES %s",
		strings.Join(pattern, ", "))
	return
}

func (s *Storage) UpdateAuditPlanByName(name string, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateAuditPlanById(id uint, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("id = ?", id).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanTotalByProjectName(projectName string) (uint64, error) {
	var count uint64
	err := s.db.
		Table("audit_plans").
		Joins("LEFT JOIN projects ON audit_plans.project_id = projects.id").
		Where("projects.name = ?", projectName).
		Where("audit_plans.deleted_at IS NULL").
		Count(&count).
		Error
	return count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetAuditPlanIDsByProjectName(projectName string) ([]uint, error) {
	ids := []struct {
		ID uint `json:"id"`
	}{}
	err := s.db.Table("audit_plans").
		Select("audit_plans.id").
		Joins("LEFT JOIN projects ON projects.id = audit_plans.project_id").
		Where("projects.name = ?", projectName).
		Find(&ids).Error

	resp := []uint{}
	for _, id := range ids {
		resp = append(resp, id.ID)
	}

	return resp, errors.ConnectStorageErrWrapper(err)
}

// GetLatestAuditPlanIds 获取所有变更过的记录，包括删除
func (s *Storage) GetLatestAuditPlanRecords(after time.Time) ([]*AuditPlan, error) {
	var aps []*AuditPlan
	err := s.db.Unscoped().Model(AuditPlan{}).Select("id, updated_at").Where("updated_at > ?", after).Order("updated_at").Find(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

type RiskAuditPlan struct {
	ReportId       uint       `json:"report_id"`
	AuditPlanName  string     `json:"audit_plan_name"`
	ReportCreateAt *time.Time `json:"report_create_at"`
	RiskSqlCOUNT   uint       `json:"risk_sql_count"`
}

func (s *Storage) GetRiskAuditPlan(projectName string) ([]*RiskAuditPlan, error) {
	var RiskAuditPlans []*RiskAuditPlan
	err := s.db.Model(AuditPlan{}).
		Select(`reports.id report_id, audit_plans.name audit_plan_name, reports.created_at report_create_at, 
				count(case when JSON_TYPE(report_sqls.audit_results)<>'NULL' then 1 else null end) risk_sql_count`).
		Joins("left join audit_plan_reports_v2 reports on audit_plans.id=reports.audit_plan_id").
		Joins("left join audit_plan_report_sqls_v2 report_sqls on report_sqls.audit_plan_report_id=reports.id").
		Joins("left join projects on projects.id=audit_plans.project_id").
		Where("reports.score<60 and projects.name=? and audit_plans.deleted_at is NULL", projectName).
		Group("audit_plans.name, reports.created_at, audit_plans.created_at, reports.id").
		Order("reports.created_at desc").Scan(&RiskAuditPlans).Error

	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}
	return RiskAuditPlans, nil

}

// 使用子查询获取最新的report时间，然后再获取最新report的sql数量和触发规则的sql数量
func (s *Storage) GetAuditPlanSQLCountAndTriggerRuleCountByProject(projectName string) (SqlCountAndTriggerRuleCount, error) {
	sqlCountAndTriggerRuleCount := SqlCountAndTriggerRuleCount{}
	subQuery := s.db.Model(&AuditPlan{}).
		Select("audit_plans.id as audit_plan_id, MAX(audit_plan_reports_v2.created_at) as latest_created_at").
		Joins("left join audit_plan_reports_v2 on audit_plan_reports_v2.audit_plan_id=audit_plans.id").
		Joins("left join projects on audit_plans.project_id=projects.id").
		Where("projects.name=? and audit_plans.deleted_at is null and audit_plan_reports_v2.id is not null", projectName).
		Group("audit_plans.id").
		SubQuery()

	err := s.db.Model(&AuditPlan{}).
		Select("count(report_sqls.id) sql_count, count(case when JSON_TYPE(report_sqls.audit_results)<>'NULL' then 1 else null end) trigger_rule_count").
		Joins("left join audit_plan_reports_v2 reports on reports.audit_plan_id=audit_plans.id").
		Joins("left join audit_plan_report_sqls_v2 report_sqls on report_sqls.audit_plan_report_id=reports.id").
		Joins("left join projects on audit_plans.project_id=projects.id").
		Joins("join (?) as sq on audit_plans.id=sq.audit_plan_id and reports.created_at=sq.latest_created_at", subQuery).
		Where("projects.name=? and audit_plans.deleted_at is null", projectName).
		Scan(&sqlCountAndTriggerRuleCount).Error

	return sqlCountAndTriggerRuleCount, errors.ConnectStorageErrWrapper(err)
}

type DBTypeAuditPlanCount struct {
	DbType         string `json:"db_type"`
	Type           string `json:"type"`
	AuditPlanCount uint   `json:"audit_plan_count"`
}

func (s *Storage) GetDBTypeAuditPlanCountByProject(projectName string) ([]*DBTypeAuditPlanCount, error) {
	dBTypeAuditPlanCounts := []*DBTypeAuditPlanCount{}
	err := s.db.Model(AuditPlan{}).
		Select("audit_plans.db_type, audit_plans.type, count(1) audit_plan_count").
		Joins("left join projects on audit_plans.project_id=projects.id").
		Where("projects.name=?", projectName).
		Group("audit_plans.db_type, audit_plans.type").Scan(&dBTypeAuditPlanCounts).Error
	return dBTypeAuditPlanCounts, errors.New(errors.ConnectStorageError, err)
}
