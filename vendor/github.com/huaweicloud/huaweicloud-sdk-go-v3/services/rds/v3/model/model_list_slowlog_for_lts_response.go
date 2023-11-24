package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSlowlogForLtsResponse Response Object
type ListSlowlogForLtsResponse struct {

	// 日志数据集合。
	SlowLogList *[]MysqlSlowLogDetailsItem `json:"slow_log_list,omitempty"`

	// 当前慢日志阈值时间。
	LongQueryTime  *string `json:"long_query_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ListSlowlogForLtsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowlogForLtsResponse struct{}"
	}

	return strings.Join([]string{"ListSlowlogForLtsResponse", string(data)}, " ")
}
