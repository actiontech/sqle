package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetDatabaseUserPrivilegeRequest Request Object
type SetDatabaseUserPrivilegeRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *SetDatabaseUserPrivilegeReqV3 `json:"body,omitempty"`
}

func (o SetDatabaseUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetDatabaseUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"SetDatabaseUserPrivilegeRequest", string(data)}, " ")
}
