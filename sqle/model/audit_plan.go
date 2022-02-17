package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/pkg/params"
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

	CreateUser    *User             `gorm:"foreignkey:CreateUserId"`
	Instance      *Instance         `gorm:"foreignkey:InstanceName;association_foreignkey:Name"`
	AuditPlanSQLs []*AuditPlanSQLV2 `gorm:"foreignkey:AuditPlanID"`
}

type AuditPlanSQLV2 struct {
	Model

	// add unique index on fingerprint and audit_plan_id
	// it's done by AutoMigrate() because gorm can't create index on TEXT column directly by tag.
	AuditPlanID uint   `json:"audit_plan_id" gorm:"not null"`
	Fingerprint string `json:"fingerprint" gorm:"type:text;not null"`
	SQLContent  string `json:"sql" gorm:"type:text;not null"`
	Info        JSON   `gorm:"type:json"`
}

func (a AuditPlanSQLV2) TableName() string {
	return "audit_plan_sqls_v2"
}

func (s *Storage) GetAuditPlans() ([]*AuditPlan, error) {
	var aps []*AuditPlan
	err := s.db.Model(AuditPlan{}).Preload("Instance").Find(&aps).Error
	return aps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanByName(name string) (*AuditPlan, bool, error) {
	ap := &AuditPlan{}
	err := s.db.Model(AuditPlan{}).Preload("Instance").Where("name = ?", name).Find(ap).Error
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
		pattern = append(pattern, "(?, ?, ?, ?)")
		args = append(args, ap.ID, sql.Fingerprint, sql.SQLContent, sql.Info)
	}
	raw = fmt.Sprintf("INSERT INTO `audit_plan_sqls_v2` (`audit_plan_id`, `fingerprint`, `sql_content`, `info`) VALUES %s",
		strings.Join(pattern, ", "))
	return
}

func (s *Storage) UpdateAuditPlanByName(name string, attrs map[string]interface{}) error {
	err := s.db.Model(AuditPlan{}).Where("name = ?", name).Update(attrs).Error
	return errors.New(errors.ConnectStorageError, err)
}
