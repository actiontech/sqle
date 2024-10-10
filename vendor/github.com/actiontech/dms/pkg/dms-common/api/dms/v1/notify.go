package v1

import (
	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
)

// swagger:parameters Notification
type NotificationReq struct {
	// notification
	// in:body
	Notification *Notification `json:"notification" validate:"required"`
}

type Notification struct {
	NotificationSubject i18nPkg.I18nStr `json:"notification_subject"`
	NotificationBody    i18nPkg.I18nStr `json:"notification_body"`
	UserUids            []string        `json:"user_uids"`
}

// swagger:model NotificationReply
type NotificationReply struct {
	// Generic reply
	base.GenericResp
}
