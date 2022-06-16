package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
)

type AuditPlan struct {
	Model
	Name             string `json:"name" gorm:"not null;index"`
	CronExpression   string `json:"cron_expression" gorm:"not null"`
	DBType           string `json:"db_type" gorm:"not null"`
	Token            string `json:"token" gorm:"not null"`
	InstanceName     string `json:"instance_name"`
	CreateUserID     uint
	InstanceDatabase string        `json:"instance_database"`
	Type             string        `json:"type"`
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
	SQLContent     string `json:"sql" gorm:"type:text;not null"`
	Info           JSON   `gorm:"type:json"`
}

func (a AuditPlanSQLV2) TableName() string {
	return "audit_plan_sqls_v2"
}

func (a *AuditPlanSQLV2) GetFingerprintMD5() string {
	if a.FingerprintMD5 != "" {
		return a.GetFingerprintMD5()
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

func (s *Storage) GetAuditPlanByName(name string) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanReportByID(id uint) (*AuditPlanReportV2, bool, error) {
	ap := &AuditPlanReportV2{}
	err := s.db.Model(AuditPlanReportV2{}).Where("id = ?", id).Preload("AuditPlan").Find(ap).Error
	if err == gorm.ErrRecordNotFound {
		return ap, false, nil
	}
	return ap, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanSQLs(name string) ([]*AuditPlanSQLV2, error) {
	ap, exist, err := s.GetAuditPlanByName(name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, gorm.ErrRecordNotFound
	}

	var sqls []*AuditPlanSQLV2
	err = s.db.Model(AuditPlanSQLV2{}).Where("audit_plan_id = ?", ap.ID).Find(&sqls).Error
	return sqls, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) OverrideAuditPlanSQLs(apName string, sqls []*AuditPlanSQLV2) error {
	ap, _, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}

	err = s.db.Unscoped().
		Model(AuditPlanSQLV2{}).
		Where("audit_plan_id = ?", ap.ID).
		Delete(&AuditPlanSQLV2{}).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}
	raw, args := getBatchInsertRawSQL(ap, sqls)
	return errors.New(errors.ConnectStorageError, s.db.Exec(fmt.Sprintf("%v;", raw), args...).Error)
}

func (s *Storage) UpdateDefaultAuditPlanSQLs(apName string, sqls []*AuditPlanSQLV2) error {
	ap, _, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}

	raw, args := getBatchInsertRawSQL(ap, sqls)
	// counter column is a accumulate value when update.
	raw += `ON DUPLICATE KEY UPDATE sql_content = VALUES(sql_content), info = JSON_SET(COALESCE(info, '{}'), 
'$.counter', COALESCE(JSON_EXTRACT(values(info), '$.counter'), 0)+COALESCE(JSON_EXTRACT(info, '$.counter'), 0),
'$.last_receive_timestamp', JSON_EXTRACT(values(info), '$.last_receive_timestamp'));`
	return errors.New(errors.ConnectStorageError, s.db.Exec(raw, args...).Error)
}

func getBatchInsertRawSQL(ap *AuditPlan, sqls []*AuditPlanSQLV2) (raw string, args []interface{}) {
	pattern := make([]string, 0, len(sqls))
	for _, sql := range sqls {
		pattern = append(pattern, "(?, ?, ?, ?, ?)")
		args = append(args, ap.ID, sql.GetFingerprintMD5(), sql.Fingerprint, sql.SQLContent, sql.Info)
	}
	raw = fmt.Sprintf("INSERT INTO `audit_plan_sqls_v2` (`audit_plan_id`,`fingerprint_md5`, `fingerprint`, `sql_content`, `info`) VALUES %s",
		strings.Join(pattern, ", "))
	return
}

func (s *Storage) UpdateAuditPlanByName(name string, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) CheckUserCanCreateAuditPlan(user *User, instName, dbType string) (bool, error) {
	if user.Name == DefaultAdminUser {
		return true, nil
	}
	instances, err := s.GetUserCanOpInstances(user, []uint{OP_AUDIT_PLAN_SAVE})
	if err != nil {
		return false, err
	}
	for _, instance := range instances {
		if instName == instance.Name {
			return true, nil
		}
	}
	return false, nil
}
