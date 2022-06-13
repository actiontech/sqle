package model

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
