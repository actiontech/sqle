package model

import "time"

const PlatformOperation = 0

type OperationRecord struct {
	Model
	OperationTime *time.Time `gorm:"column:operation_time;type:datetime;not null" json:"operation_time"`
	UserID        int64      `gorm:"column:user_id;type:bigint(20);not null" json:"user_id"`
	TypeName      string     `gorm:"column:type_name;type:varchar(255);not null" json:"type_name"`
	Content       string     `gorm:"column:content;type:varchar(255);not null" json:"content"`
	ObjectName    string     `gorm:"column:object_name;type:varchar(255);not null" json:"object_name"`
	// ProjectID is 0 means platform operation
	ProjectID int64  `gorm:"column:project_id;type:bigint(20);not null" json:"project_id"`
	Status    string `gorm:"column:status;type:varchar(255);not null" json:"status"`

	User    *User    `gorm:"foreignKey:UserID" json:"user"`
	Project *Project `gorm:"foreignKey:ProjectID" json:"project"`
}
