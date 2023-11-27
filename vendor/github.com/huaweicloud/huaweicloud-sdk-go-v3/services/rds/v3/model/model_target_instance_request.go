package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TargetInstanceRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`
}

func (o TargetInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetInstanceRequest struct{}"
	}

	return strings.Join([]string{"TargetInstanceRequest", string(data)}, " ")
}
