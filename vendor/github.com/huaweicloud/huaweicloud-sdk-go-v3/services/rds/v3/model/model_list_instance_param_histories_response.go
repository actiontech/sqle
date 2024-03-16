package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstanceParamHistoriesResponse Response Object
type ListInstanceParamHistoriesResponse struct {

	// 历史记录总数
	TotalCount *int32 `json:"total_count,omitempty"`

	// host列表
	Histories      *[]ParamGroupHistoryResult `json:"histories,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListInstanceParamHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceParamHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListInstanceParamHistoriesResponse", string(data)}, " ")
}
