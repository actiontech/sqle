package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMsdtcHostsRequest Request Object
type ListMsdtcHostsRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 分页参数 最小为0
	Offset *int32 `json:"offset,omitempty"`

	// 分页参数  取值范围为 1-100
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListMsdtcHostsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMsdtcHostsRequest struct{}"
	}

	return strings.Join([]string{"ListMsdtcHostsRequest", string(data)}, " ")
}
