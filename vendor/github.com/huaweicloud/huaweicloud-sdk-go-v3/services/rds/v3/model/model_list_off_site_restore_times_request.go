package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListOffSiteRestoreTimesRequest Request Object
type ListOffSiteRestoreTimesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 所需查询的日期，为yyyy-mm-dd字符串格式，时区为UTC。
	Date *string `json:"date,omitempty"`
}

func (o ListOffSiteRestoreTimesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOffSiteRestoreTimesRequest struct{}"
	}

	return strings.Join([]string{"ListOffSiteRestoreTimesRequest", string(data)}, " ")
}
