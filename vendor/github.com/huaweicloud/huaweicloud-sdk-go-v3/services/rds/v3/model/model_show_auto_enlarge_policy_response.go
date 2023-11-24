package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAutoEnlargePolicyResponse Response Object
type ShowAutoEnlargePolicyResponse struct {

	// 是否已开启自动扩容，true为开启
	SwitchOption *bool `json:"switch_option,omitempty"`

	// 扩容上限，单位GB
	LimitSize *int32 `json:"limit_size,omitempty"`

	// 可用空间百分比，小于等于此值或者10GB时触发扩容
	TriggerThreshold *int32 `json:"trigger_threshold,omitempty"`
	HttpStatusCode   int    `json:"-"`
}

func (o ShowAutoEnlargePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAutoEnlargePolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowAutoEnlargePolicyResponse", string(data)}, " ")
}
