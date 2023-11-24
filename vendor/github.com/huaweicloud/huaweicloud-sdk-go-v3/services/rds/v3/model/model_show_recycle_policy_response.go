package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowRecyclePolicyResponse Response Object
type ShowRecyclePolicyResponse struct {

	// 回收站实例保留天数
	RetentionPeriodInDays *int32 `json:"retention_period_in_days,omitempty"`
	HttpStatusCode        int    `json:"-"`
}

func (o ShowRecyclePolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRecyclePolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowRecyclePolicyResponse", string(data)}, " ")
}
