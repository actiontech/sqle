package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUpgradeHistoriesRequest Request Object
type ListUpgradeHistoriesRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int32 `json:"offset,omitempty"`

	// 查询记录数。默认为10，不能为负数，最小值为1，最大值为100。
	Limit *int32 `json:"limit,omitempty"`

	// 排序方式。 DESC，降序。 ASC，升序。 默认降序。
	Order *string `json:"order,omitempty"`

	// 排序字段。 start_time 开始时间。 end_time 结束时间。 默认开始时间。
	SortField *string `json:"sort_field,omitempty"`

	// 语言。默认en-us。
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ListUpgradeHistoriesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUpgradeHistoriesRequest struct{}"
	}

	return strings.Join([]string{"ListUpgradeHistoriesRequest", string(data)}, " ")
}
