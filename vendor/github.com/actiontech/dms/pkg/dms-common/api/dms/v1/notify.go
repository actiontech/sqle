package v1

import base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"

// swagger:parameters Notification
type NotificationReq struct {
	// notification
	// in:body
	Notification *Notification `json:"notification" validate:"required"`
}

type Notification struct {
	NotificationSubject string   `json:"notification_subject"`
	NotificationBody    string   `json:"notification_body"`
	UserUids            []string `json:"user_uids"`
}

// swagger:model NotificationReply
type NotificationReply struct {
	// Generic reply
	base.GenericResp
}
