package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAuditlogsRequest Request Object
type ListAuditlogsRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 查询开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartTime string `json:"start_time"`

	// 查询结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”，且大于查询开始时间，时间跨度不超过30天。  其中，T指某个时间的开始，Z指时区偏移量，例如北京时间偏移显示为+0800。
	EndTime string `json:"end_time"`

	// 索引位置，偏移量。  从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset int32 `json:"offset"`

	// 查询记录数。取值范围[1, 50]。
	Limit int32 `json:"limit"`
}

func (o ListAuditlogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAuditlogsRequest struct{}"
	}

	return strings.Join([]string{"ListAuditlogsRequest", string(data)}, " ")
}
