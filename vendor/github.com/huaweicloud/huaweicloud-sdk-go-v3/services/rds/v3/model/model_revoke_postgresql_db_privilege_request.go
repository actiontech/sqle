package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokePostgresqlDbPrivilegeRequest Request Object
type RevokePostgresqlDbPrivilegeRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *RevokePostgresqlDbPrivilegeRequestBody `json:"body,omitempty"`
}

func (o RevokePostgresqlDbPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokePostgresqlDbPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"RevokePostgresqlDbPrivilegeRequest", string(data)}, " ")
}
