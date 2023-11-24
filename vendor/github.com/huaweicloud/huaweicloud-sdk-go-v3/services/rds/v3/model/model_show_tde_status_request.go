package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTdeStatusRequest Request Object
type ShowTdeStatusRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`
}

func (o ShowTdeStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTdeStatusRequest struct{}"
	}

	return strings.Join([]string{"ShowTdeStatusRequest", string(data)}, " ")
}
