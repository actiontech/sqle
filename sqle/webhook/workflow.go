package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	workflow EventType = "workflow"

	create ActionType = "create" // 创建工单
)

var workflowCfg *WebHookConfig = &WebHookConfig{}

func UpdateWorkflowConfig(enable bool,
	maxRetryTimes, retryIntervalSeconds int, url string, token string) {

	workflowCfg.enable = enable
	workflowCfg.maxRetryTimes = maxRetryTimes
	workflowCfg.retryIntervalSeconds = retryIntervalSeconds
	workflowCfg.url = url
	workflowCfg.token = token

}

type workflowPayload struct {
	ProjectName     string `json:"project_name"`
	WorkflowID      string `json:"workflow_id"`
	WorkflowSubject string `json:"workflow_subject"`
	WorkflowStatus  string `json:"workflow_status"`
}

func TestWorkflowConfig() (err error) {
	if workflowCfg == nil {
		return fmt.Errorf("workflow webhook config missing")
	}

	if workflowCfg.url == "" {
		return fmt.Errorf("url is missing, please check webhook config")
	}

	testReqBody := &httpRequestBody{
		Event:     workflow,
		Action:    create,
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

	req, err := http.NewRequest(http.MethodPost, workflowCfg.url, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	if workflowCfg.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", workflowCfg.token))
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
