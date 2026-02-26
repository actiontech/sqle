package model

import (
	"time"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
)

// OperationRecord 用于调用 DMS 保存操作记录
type OperationRecord struct {
	OperationTime        time.Time       `gorm:"column:operation_time;type:datetime;" json:"operation_time"`
	OperationUserName    string          `gorm:"column:operation_user_name;type:varchar(255);not null" json:"operation_user_name"`
	OperationReqIP       string          `gorm:"column:operation_req_ip; type:varchar(255)" json:"operation_req_ip"`
	OperationTypeName    string          `gorm:"column:operation_type_name; type:varchar(255)" json:"operation_type_name"`
	OperationAction      string          `gorm:"column:operation_action; type:varchar(255)" json:"operation_action"`
	OperationProjectName string          `gorm:"column:operation_project_name; type:varchar(255)" json:"operation_project_name"`
	OperationStatus      string          `gorm:"column:operation_status; type:varchar(255)" json:"operation_status"`
	OperationI18nContent i18nPkg.I18nStr `gorm:"column:operation_i18n_content; type:json" json:"operation_i18n_content"`
}
