package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// FailoverRequest 强制倒换请求参数对象。
type FailoverRequest struct {

	// 是否强制倒换；true：强制倒换；false和默认null为不强制。
	Force *bool `json:"force,omitempty"`
}

func (o FailoverRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "FailoverRequest struct{}"
	}

	return strings.Join([]string{"FailoverRequest", string(data)}, " ")
}
