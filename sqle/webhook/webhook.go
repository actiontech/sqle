package webhook

type WebHookConfig struct {
	enable               bool
	maxRetryTimes        int
	retryIntervalSeconds int
	url                  string
	token                string
}

type (
	EventType  string
	ActionType string
)

type httpRequestBody struct {
	Event     EventType        `json:"event"`
	Action    ActionType       `json:"action"`
	Timestamp string           `json:"timestamp"` // time.RFC3339
	Payload   *httpBodyPayload `json:"payload"`
}

type httpBodyPayload struct {
	Workflow *workflowPayload `json:"workflow"`
}
