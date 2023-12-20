package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstanceParamHistoriesRequest Request Object
type ListInstanceParamHistoriesRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 分页参数
	Offset *int32 `json:"offset,omitempty"`

	// 分页参数
	Limit *int32 `json:"limit,omitempty"`

	// 开始时间 默认为当前时间的前7天 格式如 2020-09-01T18:50:20Z
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间 默认为当前时间 格式如 2020-09-01T18:50:20Z
	EndTime *string `json:"end_time,omitempty"`

	// 参数名称
	ParamName *string `json:"param_name,omitempty"`
}

func (o ListInstanceParamHistoriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceParamHistoriesRequest struct{}"
	}

	return strings.Join([]string{"ListInstanceParamHistoriesRequest", string(data)}, " ")
}
