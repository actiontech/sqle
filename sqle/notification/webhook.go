package notification

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
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
}

type httpBodyPayload struct {
	Workflow *workflowPayload `json:"workflow"`
}

// func TestWorkflowConfig() (err error) {
// 	return workflowSendRequest("create",
// 		"test_project", "1658637666259832832", "test_workflow", "wait_for_audit")
// }

func workflowSendRequest(action,
	projectName, workflowID, workflowSubject, workflowStatus string) (err error) {

	reqBody := &webHookRequestBody{
		Event:     "workflow",
		Action:    action,
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: &httpBodyPayload{
			Workflow: &workflowPayload{
				ProjectName:     projectName,
				WorkflowID:      workflowID,
				WorkflowSubject: workflowSubject,
				WorkflowStatus:  workflowStatus,
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
