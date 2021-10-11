package model

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
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
	InstanceDatabase string `json:"instance_database"`

	CreateUser       *User              `gorm:"foreignkey:CreateUserId"`
	Instance         *Instance          `gorm:"foreignkey:InstanceName;association_foreignkey:Name"`
	AuditPlanSQLs    []*AuditPlanSQL    `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReports []*AuditPlanReport `gorm:"foreignkey:AuditPlanID"`
}

type AuditPlanSQL struct {
	Model
	AuditPlanID          uint   `json:"audit_plan_id" gorm:"UNIQUE_INDEX:idx_audit_plan_sqls_audit_plan_id_fingerprint;not null"`
	Fingerprint          string `json:"fingerprint" gorm:"UNIQUE_INDEX:idx_audit_plan_sqls_audit_plan_id_fingerprint;not null"`
	Counter              int    `json:"counter" gorm:"not null"`
	LastSQL              string `json:"last_sql" gorm:"not null"`
	LastReceiveTimestamp string `json:"last_receive_timestamp" gorm:"not null"`
}

type AuditPlanReport struct {
	Model
	AuditPlanID uint `json:"audit_plan_id" gorm:"index"`

	AuditPlan           *AuditPlan            `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReportSQLs []*AuditPlanReportSQL `gorm:"foreignkey:AuditPlanReportID"`
}

type AuditPlanReportSQL struct {
	Model
	AuditResult string `json:"audit_result" gorm:"type:text"`

	AuditPlanSQLID    uint `json:"audit_plan_sql_id" gorm:"index"`
	AuditPlanReportID uint `json:"audit_plan_report_id" gorm:"index"`

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

func (s *Storage) OverrideAuditPlanSQLs(apName string, sqls []*AuditPlanSQL) error {
	ap, _, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}

	err = s.db.Unscoped().
		Model(AuditPlanSQL{}).
		Where("audit_plan_id = ?", ap.ID).
		Delete(&AuditPlanSQL{}).Error
	if err != nil {
		return errors.New(errors.ConnectStorageError, err)
	}

	raw, args := getBatchInsertRawSQL(ap, sqls)
	return errors.New(errors.ConnectStorageError, s.db.Exec(fmt.Sprintf("%v;", raw), args...).Error)
}

func (s *Storage) UpdateAuditPlanSQLs(apName string, sqls []*AuditPlanSQL) error {
	ap, _, err := s.GetAuditPlanByName(apName)
	if err != nil {
		return err
	}

	raw, args := getBatchInsertRawSQL(ap, sqls)
	// counter column is a accumulate value when update.
	raw += " ON DUPLICATE KEY UPDATE `counter` = VALUES(`counter`) + `counter`, `last_sql` = VALUES(`last_sql`), `last_receive_timestamp` = VALUES(`last_receive_timestamp`);"
	return errors.New(errors.ConnectStorageError, s.db.Exec(raw, args...).Error)
}

func getBatchInsertRawSQL(ap *AuditPlan, sqls []*AuditPlanSQL) (raw string, args []interface{}) {
	raw = "INSERT INTO `audit_plan_sqls` (`audit_plan_id`, `fingerprint`, `counter`, `last_sql`, `last_receive_timestamp`) VALUES "
	for i, sql := range sqls {
		if i == len(sqls)-1 {
			raw += "(?, ?, ?, ?, ?) "
		} else {
			raw += "(?, ?, ?, ?, ?), "
		}
		args = append(args, ap.ID, sql.Fingerprint, sql.Counter, sql.LastSQL, sql.LastReceiveTimestamp)
	}
	return
}

func (s *Storage) UpdateAuditPlanByName(name string, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}
