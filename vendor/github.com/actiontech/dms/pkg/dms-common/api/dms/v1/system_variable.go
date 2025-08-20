package v1

import base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

// swagger:model UpdateSystemVariablesReqV1
type UpdateSystemVariablesReqV1 struct {
	Url                         *string `json:"url" form:"url" example:"http://10.186.61.32:8080" validate:"url"`
	OperationRecordExpiredHours *int    `json:"operation_record_expired_hours" form:"operation_record_expired_hours" example:"2160"`
	CbOperationLogsExpiredHours *int    `json:"cb_operation_logs_expired_hours" form:"cb_operation_logs_expired_hours" example:"2160"`
}

// swagger:model GetSystemVariablesReply
type GetSystemVariablesReply struct {
	// Generic reply
	base.GenericResp
	Data    SystemVariablesResV1 `json:"data"`
}

type SystemVariablesResV1 struct {
	Url                         string `json:"url"`
	OperationRecordExpiredHours int    `json:"operation_record_expired_hours"`
	CbOperationLogsExpiredHours int    `json:"cb_operation_logs_expired_hours"`
}
