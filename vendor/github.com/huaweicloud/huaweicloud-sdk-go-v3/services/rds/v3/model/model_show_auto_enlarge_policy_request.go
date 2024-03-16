package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutoEnlargePolicyRequest Request Object
type ShowAutoEnlargePolicyRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowAutoEnlargePolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutoEnlargePolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowAutoEnlargePolicyRequest", string(data)}, " ")
}
