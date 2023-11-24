package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRecycleInstancesResponse Response Object
type ListRecycleInstancesResponse struct {

	// 回收站数据条数
	TotalCount *int32 `json:"total_count,omitempty"`

	// 回收站信息
	Instances      *[]RecycleInstsanceV3 `json:"instances,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o ListRecycleInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRecycleInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListRecycleInstancesResponse", string(data)}, " ")
}
