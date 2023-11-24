package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDrReplicaStatusRequest Request Object
type ShowDrReplicaStatusRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowDrReplicaStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDrReplicaStatusRequest struct{}"
	}

	return strings.Join([]string{"ShowDrReplicaStatusRequest", string(data)}, " ")
}
