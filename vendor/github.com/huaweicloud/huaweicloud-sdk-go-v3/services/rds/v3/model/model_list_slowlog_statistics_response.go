package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSlowlogStatisticsResponse Response Object
type ListSlowlogStatisticsResponse struct {

	// 当前页码
	PageNumber *int32 `json:"pageNumber,omitempty"`

	// 每页条数
	PageRecord *int32 `json:"pageRecord,omitempty"`

	// 慢日志列表
	SlowLogList *[]SlowLogStatistics `json:"slowLogList,omitempty"`

	// 总条数
	TotalRecord *int32 `json:"totalRecord,omitempty"`

	// 开始时间
	StartTime *int64 `json:"startTime,omitempty"`

	// 结束时间
	EndTime        *int64 `json:"endTime,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListSlowlogStatisticsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowlogStatisticsResponse struct{}"
	}

	return strings.Join([]string{"ListSlowlogStatisticsResponse", string(data)}, " ")
}
