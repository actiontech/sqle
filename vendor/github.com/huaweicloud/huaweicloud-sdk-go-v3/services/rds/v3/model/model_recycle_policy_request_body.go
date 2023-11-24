package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RecyclePolicyRequestBody struct {
	RecyclePolicy *RecyclePolicy `json:"recycle_policy"`
}

func (o RecyclePolicyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecyclePolicyRequestBody struct{}"
	}

	return strings.Join([]string{"RecyclePolicyRequestBody", string(data)}, " ")
}
