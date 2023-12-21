package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDbUserPrivilegeRequest Request Object
type UpdateDbUserPrivilegeRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *DbUserPrivilegeRequest `json:"body,omitempty"`
}

func (o UpdateDbUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDbUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"UpdateDbUserPrivilegeRequest", string(data)}, " ")
}
