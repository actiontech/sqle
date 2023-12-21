package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTdeStatusRequest Request Object
type UpdateTdeStatusRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateTdeStatusRequestBody `json:"body,omitempty"`
}

func (o UpdateTdeStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTdeStatusRequest struct{}"
	}

	return strings.Join([]string{"UpdateTdeStatusRequest", string(data)}, " ")
}
