package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/notification/webhook"
	"github.com/actiontech/sqle/sqle/utils/retry"
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

func TestWorkflowConfig() (err error) {
	return workflowSendRequest("workflow", "create",
		"test_project", "1658637666259832832", "test_workflow", "wait_for_audit")
}

func workflowSendRequest(event, action,
	projectName, workflowID, workflowSubject, workflowStatus string) (err error) {
	cfg := webhook.WorkflowCfg
	if cfg == nil {
		return fmt.Errorf("workflow webhook config missing")
	}

	if !cfg.Enable {
		return nil
	}

	if cfg.URL == "" {
		return fmt.Errorf("url is missing, please check webhook config")
	}

	reqBody := &webHookRequestBody{
		Event:     event,
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

	req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	if cfg.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.Token))
	}

	doneChan := make(chan struct{}, 0)
	return retry.Do(func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("response status_code(%v) body(%s)", resp.StatusCode, respBytes)
	}, doneChan,
		retry.Delay(time.Duration(cfg.RetryIntervalSeconds)),
		retry.Attempts(uint(cfg.MaxRetryTimes)))

}
