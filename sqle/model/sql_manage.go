package model

import (
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

type SqlManage struct {
	Model
	SqlFingerprint            string       `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	ProjFpSourceInstSchemaMd5 string       `json:"proj_fp_source_inst_schema_md5" gorm:"unique_index:proj_fp_source_inst_schema_md5;not null;type:varchar(255)"`
	SqlText                   string       `json:"sql_text" gorm:"type:mediumtext;not null"`
	Source                    string       `json:"source" gorm:"type:varchar(255)"`
	AuditLevel                string       `json:"audit_level" gorm:"type:varchar(255)"`
	AuditResults              AuditResults `json:"audit_results" gorm:"type:json"`
	FpCount                   uint64       `json:"fp_count"`
	FirstAppearTimestamp      *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp      *time.Time   `json:"last_receive_timestamp"`
	InstanceName              string       `json:"instance_name" gorm:"type:varchar(255)"`
	SchemaName                string       `json:"schema_name" gorm:"type:varchar(255)"`

	Assignees string `json:"assignees" gorm:"type:varchar(255)"`
	Status    string `json:"status" gorm:"default:\"unhandled\"; type:varchar(255)"`
	Remark    string `json:"remark" gorm:"type:varchar(4000)"`

	ProjectId string `json:"project_id" gorm:"type:varchar(255)"`

	AuditPlanId uint       `json:"audit_plan_id"`
	AuditPlan   *AuditPlan `gorm:"foreignkey:AuditPlanId"`
}

type SqlManageSqlAuditRecord struct {
	SQLID            string `json:"sql_id" gorm:"primary_key;auto_increment:false;type:varchar(255)"`
	SqlAuditRecordId string `json:"sql_audit_record_id" gorm:"primary_key;auto_increment:false;type:varchar(255)"`
}

func (sm SqlManageSqlAuditRecord) TableName() string {
	return "sql_manage_sql_audit_records"
}

type SqlManageEndpoint struct {
	Model
	ProjFpSourceInstSchemaMd5 string `json:"proj_fp_source_inst_schema_md5" gorm:"unique_index:uniq_md5_endpoint;type:varchar(255)"`
	Endpoint                  string `json:"endpoint" gorm:"unique_index:uniq_md5_endpoint;type:varchar(255)"`
}

func (sm SqlManageEndpoint) TableName() string {
	return "sql_manage_endpoints"
}

type SqlManageRuleTips struct {
	DbType       string                `json:"db_type"`
	RuleName     string                `json:"rule_name"`
	I18nRuleInfo driverV2.I18nRuleInfo `json:"i18n_rule_info"`
}

const PriorityHigh = "high"
