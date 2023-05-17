package webhook

var WorkflowCfg *WebHookConfig = &WebHookConfig{}

type WebHookConfig struct {
	Enable               bool
	MaxRetryTimes        int
	RetryIntervalSeconds int
	URL                  string
	Token                string
}

func UpdateWorkflowConfig(enable bool,
	maxRetryTimes, retryIntervalSeconds int, url string, token string) {

	WorkflowCfg.Enable = enable
	WorkflowCfg.MaxRetryTimes = maxRetryTimes
	WorkflowCfg.RetryIntervalSeconds = retryIntervalSeconds
	WorkflowCfg.URL = url
	WorkflowCfg.Token = token

}
