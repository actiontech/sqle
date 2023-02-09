package model

import "time"

const (
	OperationRecordPlatform = "--"

	// operation record type
	OperationRecordTypeProjectManage = "project_manage"

	// Content operation record content
	OperationRecordContentCreateProject = "create_project"

	// Status operation record status
	OperationRecordStatusSuccess = "success"
	OperationRecordStatusFail    = "fail"
)

type OperationRecord struct {
	Model
	OperationTime        time.Time `gorm:"column:operation_time;type:datetime;" json:"operation_time"`
	OperationUserName    string    `gorm:"column:operation_user_name;type:varchar(255);not null" json:"operation_user_name"`
	OperationReqIP       string    `gorm:"column:operation_req_ip" json:"operation_req_ip"`
	OperationTypeName    string    `gorm:"column:operation_type_name" json:"operation_type_name"`
	OperationAction      string    `gorm:"column:operation_action" json:"operation_action"`
	OperationObjectName  string    `gorm:"column:operation_object_name" json:"operation_object_name"`
	OperationProjectName string    `gorm:"column:operation_project_name" json:"operation_project_name"`
	OperationStatus      string    `gorm:"column:operation_status" json:"operation_status"`
}
