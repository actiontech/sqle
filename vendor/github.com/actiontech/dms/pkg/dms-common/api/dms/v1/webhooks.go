package v1

import base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

type TriggerEventType string

const (
	TriggerEventTypeWorkflow TriggerEventType = "workflow"
	TriggerEventAuditPlan    TriggerEventType = "auditplan"
)

// swagger:model
type WebHookSendMessageReq struct {
	WebHookMessage *WebHooksMessage `json:"webhook_message" validate:"required"`
}

type WebHooksMessage struct {
	Message          string           `json:"message"`
	TriggerEventType TriggerEventType `json:"trigger_event_type"`
}

// swagger:model WebHookSendMessageReply
type WebHooksSendMessageReply struct {
	// Generic reply
	base.GenericResp
}
