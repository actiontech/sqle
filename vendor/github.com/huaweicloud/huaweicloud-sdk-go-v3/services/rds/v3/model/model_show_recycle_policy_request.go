package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRecyclePolicyRequest Request Object
type ShowRecyclePolicyRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowRecyclePolicyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRecyclePolicyRequest struct{}"
	}

	return strings.Join([]string{"ShowRecyclePolicyRequest", string(data)}, " ")
}
