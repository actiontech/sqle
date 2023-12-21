package model

// import (
// 	"time"
// )

// type SyncInstanceTask struct {
// 	Model
// 	Source         string `json:"source" gorm:"not null"`
// 	Version        string `json:"version" gorm:"not null"`
// 	URL            string `json:"url" gorm:"not null"`
// 	DbType         string `json:"db_type" gorm:"not null"`
// 	RuleTemplateID uint   `json:"rule_template_id" gorm:"not null"`
// 	// Cron表达式
// 	SyncInstanceInterval string     `json:"sync_instance_interval" gorm:"not null"`
// 	LastSyncStatus       string     `json:"last_sync_status" gorm:"default:\"initialized\""`
// 	LastSyncSuccessTime  *time.Time `json:"last_sync_success_time"`

// 	// 关系表
// 	RuleTemplate *RuleTemplate `gorm:"foreignKey:RuleTemplateID"`
// }
