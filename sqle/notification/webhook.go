package notification

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
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

	ThirdPartyUserInfo string `json:"third_party_user_info"`
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
			},
		},
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
