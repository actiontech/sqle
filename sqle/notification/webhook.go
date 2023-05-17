package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/notification/webhook"
)

type httpRequestBody struct {
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

	cfg := webhook.WorkflowCfg

	if cfg == nil {
		return fmt.Errorf("workflow webhook config missing")
	}

	if cfg.URL == "" {
		return fmt.Errorf("url is missing, please check webhook config")
	}

	testReqBody := &httpRequestBody{
		Event:     "workflow",
		Action:    "create",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: &httpBodyPayload{
			Workflow: &workflowPayload{
				ProjectName:     "test_project",
				WorkflowID:      "1658637666259832832",
				WorkflowSubject: "test_workflow",
				WorkflowStatus:  "wait_for_audit",
			},
		},
	}
	b, err := json.Marshal(testReqBody)
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

	resp, err := http.DefaultClient.Do(req) // test request no need response
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return fmt.Errorf("test request response: %s", respBytes)
}
