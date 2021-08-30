package model

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
