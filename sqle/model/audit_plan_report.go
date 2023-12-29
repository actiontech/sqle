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

func (s *Storage) GetAuditPlanReportByProjectName(projectUid string) ([]*AuditPlanReportV2, error) {
	auditPlanReportV2Slice := []*AuditPlanReportV2{}
	err := s.db.Model(&AuditPlanReportV2{}).
		Joins("left join audit_plans on audit_plans.id=audit_plan_reports_v2.audit_plan_id").
		Where("audit_plans.project_id=? and audit_plans.deleted_at is NULL", projectUid).Find(&auditPlanReportV2Slice).Error
	return auditPlanReportV2Slice, errors.ConnectStorageErrWrapper(err)
}

type LatestAuditPlanReportScore struct {
	DbType       string `json:"db_type"`
	InstanceName string `json:"instance_name"`
	Score        uint   `json:"score"`
}

// 使用子查询获取最新的report的生成时间，再去获取report相关信息
func (s *Storage) GetLatestAuditPlanReportScoreFromInstanceByProject(projectUid string, instanceNames []string) ([]*LatestAuditPlanReportScore, error) {
	if len(instanceNames) == 0 {
		return nil, nil
	}

	var latestAuditPlanReportScore []*LatestAuditPlanReportScore
	subQuery := s.db.Model(&AuditPlanReportV2{}).
		Select("audit_plans.db_type, audit_plans.instance_name, MAX(audit_plan_reports_v2.created_at) as latest_created_at").
		Joins("left join audit_plans on audit_plan_reports_v2.audit_plan_id=audit_plans.id").
		Where("audit_plans.project_id=?", projectUid).
		Where("audit_plans.instance_name in (?)", instanceNames).
		Group("audit_plans.db_type, audit_plans.instance_name").
		SubQuery()

	err := s.db.Model(&AuditPlanReportV2{}).
		Select("audit_plans.db_type, audit_plans.instance_name, audit_plan_reports_v2.score").
		Joins("left join audit_plans on audit_plan_reports_v2.audit_plan_id=audit_plans.id").
		Joins("join (?) as sq on audit_plans.db_type=sq.db_type and audit_plans.instance_name=sq.instance_name and audit_plan_reports_v2.created_at=sq.latest_created_at", subQuery).
		Where("audit_plans.project_id=? and audit_plans.deleted_at is null", projectUid).
		Scan(&latestAuditPlanReportScore).Error

	return latestAuditPlanReportScore, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetReportWithAuditPlanByReportID(reportId int) (auditPlanReportV2 *AuditPlanReportV2, exist bool, err error) {
	auditPlanReportV2 = &AuditPlanReportV2{}
	err = s.db.Preload("AuditPlan").Preload("AuditPlanReportSQLs").Where("id=?", reportId).Find(auditPlanReportV2).Error
	if err == gorm.ErrRecordNotFound {
		return auditPlanReportV2, false, nil
	}
	return auditPlanReportV2, true, errors.New(errors.ConnectStorageError, err)
}
