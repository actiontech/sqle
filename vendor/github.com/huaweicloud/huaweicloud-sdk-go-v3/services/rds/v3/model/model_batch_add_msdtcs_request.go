package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAddMsdtcsRequest Request Object
type BatchAddMsdtcsRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *AddMsdtcRequestBody `json:"body,omitempty"`
}

func (o BatchAddMsdtcsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAddMsdtcsRequest struct{}"
	}

	return strings.Join([]string{"BatchAddMsdtcsRequest", string(data)}, " ")
}
