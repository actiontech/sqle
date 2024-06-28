package model

import (
	"time"
)

type SQLDevRecord struct {
	Model
	ProjectId                 string       `json:"project_id" gorm:"type:varchar(255)"`
	SqlFingerprint            string       `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	ProjFpSourceInstSchemaMd5 string       `json:"proj_fp_source_inst_schema_md5" gorm:"unique_index:creator_proj_fp_source_inst_schema_md5;not null;type:varchar(255)"`
	SqlText                   string       `json:"sql_text" gorm:"type:mediumtext;not null"`
	Source                    string       `json:"source" gorm:"type:varchar(255)"`
	AuditLevel                string       `json:"audit_level" gorm:"type:varchar(255)"`
	AuditResults              AuditResults `json:"audit_results" gorm:"type:json"`
	FpCount                   uint64       `json:"fp_count"`
	FirstAppearTimestamp      *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp      *time.Time   `json:"last_receive_timestamp"`
	InstanceName              string       `json:"instance_name" gorm:"type:varchar(255)"`
	SchemaName                string       `json:"schema_name" gorm:"type:varchar(255)"`
	Creator                   string       `json:"creator"  gorm:"unique_index:creator_proj_fp_source_inst_schema_md5;not null;type:varchar(255)"`
}

func (sm SQLDevRecord) TableName() string {
	return "sql_dev_records"
}
