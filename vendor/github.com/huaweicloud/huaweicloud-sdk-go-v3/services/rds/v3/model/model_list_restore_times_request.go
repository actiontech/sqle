package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListRestoreTimesRequest Request Object
type ListRestoreTimesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 所需查询的日期，为yyyy-mm-dd字符串格式，时区为UTC。
	Date *string `json:"date,omitempty"`
}

func (o ListRestoreTimesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRestoreTimesRequest struct{}"
	}

	return strings.Join([]string{"ListRestoreTimesRequest", string(data)}, " ")
}
