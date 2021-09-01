package model

import (
	"actiontech.cloud/sqle/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

type AuditPlan struct {
	Model
	Name             string `json:"name" gorm:"not null"`
	CronExpression   string `json:"cron_expression" gorm:"not null"`
	DBType           string `json:"db_type" gorm:"not null"`
	Token            string `json:"token" gorm:"not null"`
	InstanceName     string `json:"instance_name"`
	CreateUserID     uint
	InstanceDatabase string `json:"instance_database"`

	CreateUser       *User              `gorm:"foreignkey:CreateUserId"`
	Instance         *Instance          `gorm:"foreignkey:InstanceName;association_foreignkey:Name"`
	AuditPlanSQLs    []*AuditPlanSQL    `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReports []*AuditPlanReport `gorm:"foreignkey:AuditPlanID"`
}

type AuditPlanSQL struct {
	Model
	AuditPlanID string `json:"audit_plan_id" gorm:"index"`

	Fingerprint          string `json:"fingerprint" gorm:"not null"`
	Counter              string `json:"counter" gorm:"not null"`
	LastSQLText          string `json:"last_sql" gorm:"not null"`
	LastReceiveTimestamp string `json:"last_receive_timestamp" gorm:"not null"`
}

type AuditPlanReport struct {
	Model
	AuditPlanID string `json:"audit_plan_id" gorm:"index"`

	AuditPlan           *AuditPlan            `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReportSQLs []*AuditPlanReportSQL `gorm:"foreignkey:AuditPlanReportID"`
}

type AuditPlanReportSQL struct {
	Model
	AuditResult string `json:"audit_result" gorm:"type:text"`

	AuditPlanSQLID    string `json:"audit_plan_sql_id" gorm:"index"`
	AuditPlanReportID string `json:"audit_plan_report_id" gorm:"index"`

	AuditPlanSQL    *AuditPlanSQL    `gorm:"foreignkey:AuditPlanSQLID"`
	AuditPlanReport *AuditPlanReport `gorm:"foreignkey:AuditPlanReportID"`
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

func (s *Storage) GetAuditPlanSQLs(name string) ([]*AuditPlanSQL, error) {
	ap, exist, err := s.GetAuditPlanByName(name)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, gorm.ErrRecordNotFound
	}

	var sqls []*AuditPlanSQL
	err = s.db.Model(AuditPlanSQL{}).Where("audit_plan_id = ?", ap.ID).Find(&sqls).Error
	return sqls, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateAuditPlanByName(name string, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}
