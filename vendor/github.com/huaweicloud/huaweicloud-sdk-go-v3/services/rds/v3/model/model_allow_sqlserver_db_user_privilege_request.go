package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowSqlserverDbUserPrivilegeRequest Request Object
type AllowSqlserverDbUserPrivilegeRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *SqlserverGrantRequest `json:"body,omitempty"`
}

func (o AllowSqlserverDbUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowSqlserverDbUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"AllowSqlserverDbUserPrivilegeRequest", string(data)}, " ")
}
