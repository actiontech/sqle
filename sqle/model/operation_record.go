package model

import "time"

type OperationRecord struct {
	Model
	OperationTime        time.Time `gorm:"column:operation_time;type:datetime;" json:"operation_time"`
	OperationUserName    string    `gorm:"column:operation_user_name;type:varchar(255);not null" json:"operation_user_name"`
	OperationReqIP       string    `gorm:"column:operation_req_ip" json:"operation_req_ip"`
	OperationTypeName    string    `gorm:"column:operation_type_name" json:"operation_type_name"`
	OperationAction      string    `gorm:"column:operation_action" json:"operation_action"`
	OperationContent     string    `gorm:"column:operation_content" json:"operation_content"`
	OperationProjectName string    `gorm:"column:operation_project_name" json:"operation_project_name"`
	OperationStatus      string    `gorm:"column:operation_status" json:"operation_status"`
}
