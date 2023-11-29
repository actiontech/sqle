package model

import (
	"time"
)

type SqlManage struct {
	Model
	SqlFingerprint            string       `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	ProjFpSourceInstSchemaMd5 string       `json:"proj_fp_source_inst_schema_md5" gorm:"unique_index:proj_fp_source_inst_schema_md5;not null"`
	SqlText                   string       `json:"sql_text" gorm:"type:mediumtext;not null"`
	Source                    string       `json:"source"`
	AuditLevel                string       `json:"audit_level"`
	AuditResults              AuditResults `json:"audit_results" gorm:"type:json"`
	FpCount                   uint64       `json:"fp_count"`
	FirstAppearTimestamp      *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp      *time.Time   `json:"last_receive_timestamp"`
	InstanceName              string       `json:"instance_name"`
	SchemaName                string       `json:"schema_name"`

	Assignees []*User `gorm:"many2many:sql_manage_assignees;"`
	Status    string  `json:"status" gorm:"default:\"unhandled\""`
	Remark    string  `json:"remark" gorm:"type:varchar(4000)"`

	ProjectId uint     `json:"project_id"`
	Project   *Project `gorm:"foreignkey:ProjectId"`

	AuditPlanId uint       `json:"audit_plan_id"`
	AuditPlan   *AuditPlan `gorm:"foreignkey:AuditPlanId"`
}

type SqlManageSqlAuditRecord struct {
	ProjFpSourceInstSchemaMd5 string `json:"proj_fp_source_inst_schema_md5" gorm:"primary_key;auto_increment:false;"`
	SqlAuditRecordId          uint   `json:"sql_audit_record_id" gorm:"primary_key;auto_increment:false;"`
}

func (sm SqlManageSqlAuditRecord) TableName() string {
	return "sql_manage_sql_audit_records"
}
