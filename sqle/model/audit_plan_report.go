package model

type AuditPlanReportV2 struct {
	Model
	AuditPlanID uint `json:"audit_plan_id" gorm:"index"`
	TaskID      uint `json:"task_id" gorm:"index"`

	AuditPlan           *AuditPlan              `gorm:"foreignkey:AuditPlanID"`
	AuditPlanReportSQLs []*AuditPlanReportSQLV2 `gorm:"foreignkey:AuditPlanReportID"`
	Task                *Task                   `gorm:"foreignkey:TaskID"`
}

func (a AuditPlanReportV2) TableName() string {
	return "audit_plan_reports_v2"
}

type AuditPlanReportSQLV2 struct {
	Model
	AuditPlanReportID uint   `json:"audit_plan_report_id" gorm:"index"`
	SQL               string `json:"sql" gorm:"type:text;not null"`
	AuditResult       string `json:"audit_result" gorm:"type:text"`

	AuditPlanReport *AuditPlanReportV2 `gorm:"foreignkey:AuditPlanReportID"`
}

func (a AuditPlanReportSQLV2) TableName() string {
	return "audit_plan_report_sqls_v2"
}
