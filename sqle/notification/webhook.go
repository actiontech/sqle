package notification

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
)

type webHookRequestBody struct {
	Event     string           `json:"event"`
	Action    string           `json:"action"`
	Timestamp string           `json:"timestamp"` // time.RFC3339
	Payload   *httpBodyPayload `json:"payload"`
}

type workflowPayload struct {
	ProjectName     string `json:"project_name"`
	WorkflowID      string `json:"workflow_id"`
	WorkflowSubject string `json:"workflow_subject"`
	WorkflowStatus  string `json:"workflow_status"`

	ThirdPartyUserInfo string          `json:"third_party_user_info"`
	CurrentStepID      uint            `json:"current_step_info"`
	WorkflowTaskID     uint            `json:"workflow_task_id"`
	InstanceInfo       []InstanceInfo  `json:"instanceInfo"`
	WorkflowDesc       string          `json:"workflow_desc"`
	SqlTypeMap         map[string]bool `json:"sql_type_map"` // DDL DML DQL
}

type InstanceInfo struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	Port   string `json:"port"`
	Desc   string `json:"desc"`
}

type httpBodyPayload struct {
	Workflow *workflowPayload `json:"workflow"`
}

// func TestWorkflowConfig() (err error) {
// 	return workflowSendRequest("create",
// 		"test_project", "1658637666259832832", "test_workflow", "wait_for_audit")
// }

func workflowSendRequest(action string, workflow *model.Workflow) (err error) {
	user, err := dms.GetUser(context.TODO(), workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	reqBody := &webHookRequestBody{
		Event:     "workflow",
		Action:    action,
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: &httpBodyPayload{
			Workflow: &workflowPayload{
				ProjectName:        string(workflow.ProjectId),
				WorkflowID:         workflow.WorkflowId,
				WorkflowSubject:    workflow.Subject,
				WorkflowStatus:     workflow.Record.Status,
				ThirdPartyUserInfo: user.ThirdPartyUserInfo,
				CurrentStepID:      workflow.CurrentStep().ID,
				WorkflowDesc:       workflow.Desc,
				SqlTypeMap: map[string]bool{
					driverV2.SQLTypeDDL: false,
					driverV2.SQLTypeDML: false,
					driverV2.SQLTypeDQL: false,
				},
			},
		},
	}
	for _, record := range workflow.Record.InstanceRecords {
		if record.Instance == nil {
			continue
		}
		info := InstanceInfo{
			Host: record.Instance.Host,
			Port: record.Instance.Port,
			Desc: record.Instance.Desc,
		}
		if record.Task == nil {
			reqBody.Payload.Workflow.InstanceInfo = append(reqBody.Payload.Workflow.InstanceInfo, info)
			continue
		}
		info.Schema = record.Task.Schema
		reqBody.Payload.Workflow.WorkflowTaskID = record.Task.ID
		for _, executeSql := range record.Task.ExecuteSQLs {
			if executeSql.SQLType != "" {
				reqBody.Payload.Workflow.SqlTypeMap[executeSql.SQLType] = true
			}
		}
		reqBody.Payload.Workflow.InstanceInfo = append(reqBody.Payload.Workflow.InstanceInfo, info)
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	return dmsobject.WebHookSendMessage(context.TODO(), controller.GetDMSServerAddress(), &v1.WebHookSendMessageReq{
		WebHookMessage: &v1.WebHooksMessage{
			Message:          string(b),
			TriggerEventType: v1.TriggerEventTypeWorkflow,
		},
	})

}
