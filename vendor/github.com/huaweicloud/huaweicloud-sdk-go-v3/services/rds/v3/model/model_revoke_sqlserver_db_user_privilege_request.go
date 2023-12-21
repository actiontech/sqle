package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokeSqlserverDbUserPrivilegeRequest Request Object
type RevokeSqlserverDbUserPrivilegeRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *SqlserverRevokeRequest `json:"body,omitempty"`
}

func (o RevokeSqlserverDbUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeSqlserverDbUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"RevokeSqlserverDbUserPrivilegeRequest", string(data)}, " ")
}
