package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListMsdtcHostsResponse Response Object
type ListMsdtcHostsResponse struct {

	// host总数
	TotalCount *int32 `json:"total_count,omitempty"`

	// host列表
	Hosts          *[]DbsInstanceHostInfoResult `json:"hosts,omitempty"`
	HttpStatusCode int                          `json:"-"`
}

func (o ListMsdtcHostsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMsdtcHostsResponse struct{}"
	}

	return strings.Join([]string{"ListMsdtcHostsResponse", string(data)}, " ")
}
