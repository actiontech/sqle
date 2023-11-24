package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListErrorlogForLtsResponse Response Object
type ListErrorlogForLtsResponse struct {

	// 日志数据集合。
	ErrorLogList   *[]ErrorLogItem `json:"error_log_list,omitempty"`
	HttpStatusCode int             `json:"-"`
}

func (o ListErrorlogForLtsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListErrorlogForLtsResponse struct{}"
	}

	return strings.Join([]string{"ListErrorlogForLtsResponse", string(data)}, " ")
}
