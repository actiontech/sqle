package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListErrorLogsResponse Response Object
type ListErrorLogsResponse struct {
	ErrorLogList *[]ErrorLog `json:"error_log_list,omitempty"`

	// 总记录数。
	TotalRecord    *int32 `json:"total_record,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListErrorLogsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListErrorLogsResponse struct{}"
	}

	return strings.Join([]string{"ListErrorLogsResponse", string(data)}, " ")
}
