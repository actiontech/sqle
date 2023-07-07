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
	AuditPlanReportID uint         `json:"audit_plan_report_id" gorm:"index"`
	SQL               string       `json:"sql" gorm:"type:mediumtext;not null"`
	Number            uint         `json:"number"`
	AuditResults      AuditResults `json:"audit_results" gorm:"type:json"`
	Schema            string       `json:"schema" gorm:"type:varchar(512);not null"`

	AuditPlanReport *AuditPlanReportV2 `gorm:"foreignkey:AuditPlanReportID"`
}

func (a AuditPlanReportSQLV2) TableName() string {
	return "audit_plan_report_sqls_v2"
}

func (s *Storage) GetAuditPlanReportSQLV2ByReportIDAndNumber(reportId, number uint) (
	auditPlanReportSQLV2 *AuditPlanReportSQLV2, exist bool, err error) {

	auditPlanReportSQLV2 = &AuditPlanReportSQLV2{}
	err = s.db.Where("audit_plan_report_id = ? and number = ?", reportId, number).Find(auditPlanReportSQLV2).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlanReportSQLV2, false, nil
	}

	return auditPlanReportSQLV2, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetAuditPlanReportByProjectName(projectName string) ([]*AuditPlanReportV2, error) {
	auditPlanReportV2Slice := []*AuditPlanReportV2{}
	err := s.db.Model(&AuditPlanReportV2{}).
		Joins("left join audit_plans on audit_plans.id=audit_plan_reports_v2.audit_plan_id").
		Joins("left join projects on projects.id=audit_plans.project_id").
		Where("projects.name=? and audit_plans.deleted_at is NULL", projectName).Find(&auditPlanReportV2Slice).Error
	return auditPlanReportV2Slice, errors.ConnectStorageErrWrapper(err)
}
