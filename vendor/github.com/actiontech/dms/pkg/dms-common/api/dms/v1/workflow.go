package v1

import (
	"time"

	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
)

type FilterGlobalDataExportWorkflowReq struct {
	FilterStatusList                 []DataExportWorkflowStatus `json:"filter_status_list" query:"filter_status_list"`
	FilterDBServiceUid               string                     `json:"filter_db_service_uid" query:"filter_db_service_uid"`
	FilterProjectUids                []string                   `json:"filter_project_uids" query:"filter_project_uids"`
	FilterCurrentStepAssigneeUserId  string                     `json:"filter_current_step_assignee_user_id" query:"filter_current_step_assignee_user_id"`
	FilterProjectUid                 string                     `json:"filter_project_uid" query:"filter_project_uid"`
	PageSize                         uint32                     `query:"page_size" json:"page_size" validate:"required"`
	PageIndex                        uint32                     `query:"page_index" json:"page_index"`
	FilterByStatus                   DataExportWorkflowStatus   `query:"filter_by_status" json:"filter_by_status"`
	FilterByCreateUserUid            string                     `json:"filter_by_create_user_uid" query:"filter_by_create_user_uid"`
	FilterCurrentStepAssigneeUserUid string                     `json:"filter_current_step_assignee_user_uid" query:"filter_current_step_assignee_user_uid"`
	FilterByDBServiceUid             string                     `json:"filter_by_db_service_uid" query:"filter_by_db_service_uid"`
	FuzzyKeyword                     string                     `json:"fuzzy_keyword" query:"fuzzy_keyword"`
}

type DataExportWorkflowStatus string

const (
	DataExportWorkflowStatusWaitForApprove   DataExportWorkflowStatus = "wait_for_approve"
	DataExportWorkflowStatusWaitForExport    DataExportWorkflowStatus = "wait_for_export"
	DataExportWorkflowStatusWaitForExporting DataExportWorkflowStatus = "exporting"
	DataExportWorkflowStatusRejected         DataExportWorkflowStatus = "rejected"
	DataExportWorkflowStatusCancel           DataExportWorkflowStatus = "cancel"
	DataExportWorkflowStatusFailed           DataExportWorkflowStatus = "failed"
	DataExportWorkflowStatusFinish           DataExportWorkflowStatus = "finish"
)

type ListDataExportWorkflow struct {
	WorkflowID               string                      `json:"workflow_uid"`                    // 数据导出工单ID
	WorkflowName             string                      `json:"workflow_name"`                   // 数据导出工单的名称
	Description              string                      `json:"desc"`                            // 数据导出工单的描述
	Creater                  UidWithName                 `json:"creater"`                         // 数据导出工单的创建人
	CreatedAt                time.Time                   `json:"created_at"`                      // 数据导出工单的创建时间
	Status                   DataExportWorkflowStatus    `json:"status"`                          // 数据导出工单的状态
	CurrentStepAssigneeUsers []UidWithName               `json:"current_step_assignee_user_list"` // 工单待操作人
	DBServiceInfos           []*DBServiceUidWithNameInfo `json:"db_service_info,omitempty"`       // 所属数据源信息
	ProjectInfo              *ProjectInfo                `json:"project_info,omitempty"`          // 所属项目信息
}

type ListDataExportWorkflowsReply struct {
	Data  []*ListDataExportWorkflow `json:"data"`
	Total int64                     `json:"total_nums"`
	base.GenericResp
}
