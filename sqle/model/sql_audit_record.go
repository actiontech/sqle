package model

type SQLAuditRecordTags struct {
	Tags []string `json:"tags"`
}

type SQLAuditRecord struct {
	Model
	ProjectId     uint `gorm:"index"`
	CreatorID     uint
	AuditRecordID string `gorm:"unique"`
	Tags          JSON

	TaskId uint
	Task   *Task `gorm:"foreignkey:TaskId"`
}
