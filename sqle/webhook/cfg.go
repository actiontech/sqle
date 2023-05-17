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
