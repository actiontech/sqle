package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SecurityGroupRequest struct {

	// - 安全组ID。
	SecurityGroupId string `json:"security_group_id"`
}

func (o SecurityGroupRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SecurityGroupRequest struct{}"
	}

	return strings.Join([]string{"SecurityGroupRequest", string(data)}, " ")
}
