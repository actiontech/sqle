package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowDbUserPrivilegeRequest Request Object
type AllowDbUserPrivilegeRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *GrantRequest `json:"body,omitempty"`
}

func (o AllowDbUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowDbUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"AllowDbUserPrivilegeRequest", string(data)}, " ")
}
