package model

type SQLAuditRecordTags struct {
	Tags []string `json:"tags"`
}

type SQLAuditRecord struct {
	Model
	ProjectId     uint   `gorm:"index;not null"`
	CreatorID     uint   `gorm:"not null"`
	AuditRecordID string `gorm:"unique;not null"`
	Tags          JSON

	TaskId uint  `gorm:"not null"`
	Task   *Task `gorm:"foreignkey:TaskId"`
}
