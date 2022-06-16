package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

type AuditPlanReportV2 struct {
	Model
	AuditPlanID uint `json:"audit_plan_id" gorm:"index"`

	AuditPlan           *AuditPlan              `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReportSQLs []*AuditPlanReportSQLV2 `gorm:"foreignkey:AuditPlanReportID"`
	PassRate            float64                 `json:"pass_rate"`
	Score               int32                   `json:"score"`
	AuditLevel          string                  `json:"audit_level"`
}

func (a AuditPlanReportV2) TableName() string {
	return "audit_plan_reports_v2"
}

type AuditPlanReportSQLV2 struct {
	Model
	AuditPlanReportID uint   `json:"audit_plan_report_id" gorm:"index"`
	SQL               string `json:"sql" gorm:"type:text;not null"`
	Number            uint   `json:"number"`
	AuditResult       string `json:"audit_result" gorm:"type:text"`

	AuditPlanReport *AuditPlanReportV2 `gorm:"foreignkey:AuditPlanReportID"`
}

func (a AuditPlanReportSQLV2) TableName() string {
	return "audit_plan_report_sqls_v2"
}

func (s *Storage) GetAuditPlanReportSQLV2ByReportIDAndNumber(reportId, number uint) (*AuditPlanReportSQLV2, bool, error) {
	auditPlanReportSQLV2 := &AuditPlanReportSQLV2{}
	err := s.db.Where("audit_plan_report_id = ? and number = ?", reportId, number).Find(auditPlanReportSQLV2).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlanReportSQLV2, false, nil
	}

	return auditPlanReportSQLV2, true, errors.New(errors.ConnectStorageError, err)
}
