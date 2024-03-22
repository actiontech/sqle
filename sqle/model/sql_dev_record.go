package model

import (
	"time"
)

type SQLDevRecord struct {
	Model
	ProjectId                 string       `json:"project_id"`
	SqlFingerprint            string       `json:"sql_fingerprint" gorm:"type:mediumtext;not null"`
	ProjFpSourceInstSchemaMd5 string       `json:"proj_fp_source_inst_schema_md5" gorm:"unique_index:creator_proj_fp_source_inst_schema_md5;not null"`
	SqlText                   string       `json:"sql_text" gorm:"type:mediumtext;not null"`
	Source                    string       `json:"source"`
	AuditLevel                string       `json:"audit_level"`
	AuditResults              AuditResults `json:"audit_results" gorm:"type:json"`
	FpCount                   uint64       `json:"fp_count"`
	FirstAppearTimestamp      *time.Time   `json:"first_appear_timestamp"`
	LastReceiveTimestamp      *time.Time   `json:"last_receive_timestamp"`
	InstanceName              string       `json:"instance_name"`
	SchemaName                string       `json:"schema_name"`
	Creator                   string       `json:"creator"  gorm:"unique_index:creator_proj_fp_source_inst_schema_md5;not null"`
}

func (sm SQLDevRecord) TableName() string {
	return "sql_dev_records"
}
