package model

import "time"

const (
	OperationRecordPlatform = "--"

	// operation record type
	OperationRecordProjectManageType = "project_manage"

	// Content operation record content
	OperationRecordCreateProjectContent = "create_project"

	// Status operation record status
	OperationRecordSuccessStatus = "success"
	OperationRecordFailStatus    = "fail"
)

type OperationRecord struct {
	Model
	OperationTime time.Time `gorm:"column:operation_time;type:datetime;not null" json:"operation_time"`
	UserName      string    `gorm:"column:user_name;type:varchar(255);not null" json:"user_name"`
	IP            string    `gorm:"column:ip;type:varchar(255);not null" json:"ip"`
	TypeName      string    `gorm:"column:type_name;type:varchar(255);not null" json:"type_name"`
	Content       string    `gorm:"column:content;type:varchar(255);not null" json:"content"`
	ObjectName    string    `gorm:"column:object_name;type:varchar(255);not null" json:"object_name"`
	ProjectName   string    `gorm:"column:project_name;type:varchar(255);not null" json:"project_name"`
	Status        string    `gorm:"column:status;type:varchar(255);not null" json:"status"`
}
