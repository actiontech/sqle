package webhook

const (
	workflow EventType = "workflow"

	create ActionType = "create" // 创建工单
	audit  ActionType = "audit"  // 审核工单
	exec   ActionType = "exec"   // 上线工单
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
