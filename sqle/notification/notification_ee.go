//go:build enterprise
// +build enterprise

package notification

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
)

func (n *AuditPlanNotifier) sendWebHook(notification Notification, webHookUrl, webHookTemplate string) error {
	message, err := n.splicingWebHookMessage(webHookTemplate, notification)
	if err != nil {
		return err
	}
	return n.sendJsonReq(webHookUrl, message)
}

func (n *AuditPlanNotifier) splicingWebHookMessage(webHookTemplate string, notification Notification) (string, error) {
	temp, err := template.New("webhook").Parse(webHookTemplate)
	if err != nil {
		return "", err
	}
	// Newline characters are not allowed in json message and need to be replaced with '\n' string
	data := map[string]string{
		"subject": strings.ReplaceAll(notification.NotificationSubject(), "\n", "\\n"),
		"body":    strings.ReplaceAll(notification.NotificationBody(), "\n", "\\n"),
	}
	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		return "", err
	}

	return buff.String(), err
}

func (n *AuditPlanNotifier) sendJsonReq(webHookUrl string, message string) error {

	req, err := http.NewRequest(http.MethodPost, webHookUrl, bytes.NewBufferString(message))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("request failed, status code: %v, webhook return value read failed: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("request failed, status code: %v, webhook return value: %v", resp.StatusCode, string(body))
	}
	return resp.Body.Close()
}
