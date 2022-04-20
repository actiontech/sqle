//go:build !enterprise
// +build !enterprise

package notification

func (n *AuditPlanNotifier) sendWebHook(notification Notification, webHookUrl, webHookTemplate string) error {
	return nil
}
